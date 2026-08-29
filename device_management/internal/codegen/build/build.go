// Package build flattens parsed specs into fully-resolved view models. It
// is the only codegen stage that inspects the spec model; every naming and
// type decision funnels through here (via the naming package) so templates
// stay decision-free.
//
// Presence rules: required keys become value fields, optional keys become
// pointers (or nil-able slices/maps) tagged omitempty. Constraints
// (rangelist, range, format, repetition, subtypes) become statements in the
// generated Validate method.
package build

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen/naming"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen/view"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/spec"
)

const commentWidth = 96

// Kind selects the wire-identifier method a spec's main struct carries.
type Kind int

const (
	KindCommand     Kind = iota // RequestType()
	KindProfile                 // PayloadType()
	KindDeclaration             // DeclarationType()
	KindPlain                   // no identifier method
)

// SharedTypes tracks every named nested struct emitted in one package,
// keyed by name with a shape signature. It serves two purposes: subkeytype
// reuse across specs (first definition wins, later references reuse), and
// shape-aware dedup — identical (name, shape) pairs collapse to one type,
// while a name collision with a different shape gets a numeric suffix.
type SharedTypes struct {
	byName map[string]string // type name -> shape signature
	// enumMembers records the constant names of every emitted enum type so
	// a reusing spec can reference them without re-emitting.
	enumMembers map[string][]string
}

// NewSharedTypes returns an empty per-package shared-type registry.
func NewSharedTypes() *SharedTypes {
	return &SharedTypes{byName: map[string]string{}, enumMembers: map[string][]string{}}
}

// File builds the view model for one spec.
func File(s *spec.Spec, pkg string, kind Kind, shared *SharedTypes, usedConsts map[string]bool) (*view.File, error) {
	f := &view.File{
		Name:        naming.FileName(s.Name),
		PackageName: pkg,
	}
	b := &builder{file: f, shared: shared, usedConsts: usedConsts}

	mainName := claimName(shared, naming.StructName(s.TypeIdentifier(), s.Name), "main:"+s.Category+"/"+s.Name)
	method, value := "", ""
	switch kind {
	case KindCommand:
		method, value = "RequestType", s.Payload.RequestType
	case KindProfile:
		method, value = "PayloadType", s.Payload.PayloadType
	case KindDeclaration:
		method, value = "DeclarationType", s.Payload.DeclarationType
	}

	main := b.buildStruct(mainName, specComment(s), s.PayloadKeys, method, value)
	f.MainStruct = main.Name
	f.Structs = append(f.Structs, main)
	if len(s.ResponseKeys) > 0 {
		respName := claimName(shared, mainName+"Response", "resp:"+s.Category+"/"+s.Name)
		comment := []string{respName + " models the device response to the " + mainName + " command."}
		resp := b.buildStruct(respName, comment, s.ResponseKeys, "", "")
		f.Structs = append(f.Structs, resp)
	}
	// Nested subkey structs render after the payloads they support.
	f.Structs = append(f.Structs, b.nested...)
	return f, nil
}

type builder struct {
	file       *view.File
	shared     *SharedTypes
	usedConsts map[string]bool
	nested     []view.Struct
}

// buildStruct builds one struct; nested subkey structs it discovers are
// accumulated on the builder for the caller to place.
func (b *builder) buildStruct(name string, comment []string, keys []spec.Key, typeMethod, typeValue string) view.Struct {
	st := view.Struct{
		Name:         name,
		CommentLines: comment,
		TypeMethod:   typeMethod,
		TypeValue:    typeValue,
	}
	usedFields := map[string]bool{}
	var validation []string

	// The wire-identifier method must not collide with a spec key of the
	// same name (profile specs declare a literal PayloadType key).
	if typeMethod != "" {
		usedFields[typeMethod] = true
	}

	for i := range keys {
		k := &keys[i]
		fieldName := b.fieldName(k.Key, usedFields)
		enumType, enumMembers := b.enumFor(name, fieldName, k)
		goType, nested := b.fieldType(name, fieldName, k, enumType)
		st.Fields = append(st.Fields, view.Field{
			Name:         fieldName,
			Type:         goType,
			Tag:          tag(k),
			CommentLines: keyComment(k),
		})
		validation = append(validation, b.checks(fieldName, goType, k, nested, enumType, enumMembers)...)
	}
	st.ValidateLines = validation
	if len(validation) > 0 {
		b.file.NeedsValidate = b.file.NeedsValidate || usesValidate(validation)
	}
	return st
}

func (b *builder) fieldName(key string, used map[string]bool) string {
	base := naming.ExportName(key)
	name := base
	for i := 2; used[name]; i++ {
		name = base + strconv.Itoa(i)
	}
	used[name] = true
	return name
}

// fieldType resolves a key to its Go type, emitting nested structs as
// needed. nested reports the named nested struct type ("" when none) so
// checks can recurse into it.
func (b *builder) fieldType(structName, fieldName string, k *spec.Key, enumType string) (goType string, nested string) {
	optional := !k.Required()
	switch k.Type {
	case "<string>":
		if enumType != "" {
			return ptrIf(optional, enumType), ""
		}
		return ptrIf(optional, "string"), ""
	case "<integer>":
		if enumType != "" {
			return ptrIf(optional, enumType), ""
		}
		return ptrIf(optional, "int64"), ""
	case "<real>":
		return ptrIf(optional, "float64"), ""
	case "<boolean>":
		return ptrIf(optional, "bool"), ""
	case "<date>":
		b.file.NeedsTime = true
		return ptrIf(optional, "time.Time"), ""
	case "<data>":
		return "[]byte", ""
	case "<dictionary>":
		if len(k.Subkeys) == 0 {
			return "map[string]any", ""
		}
		n := b.nestedStruct(structName+fieldName, k)
		return ptrIf(optional, n), n
	case "<array>":
		itemType, itemNested := b.itemType(structName, fieldName, k)
		return "[]" + itemType, itemNested
	default: // <any> and unrecognised types transport untyped
		return "any", ""
	}
}

// itemType resolves an array's item type from its subkeys.
func (b *builder) itemType(structName, fieldName string, k *spec.Key) (goType, nested string) {
	if len(k.Subkeys) == 0 {
		return "any", ""
	}
	if len(k.Subkeys) == 1 {
		item := &k.Subkeys[0]
		switch item.Type {
		case "<string>":
			return "string", ""
		case "<integer>":
			return "int64", ""
		case "<real>":
			return "float64", ""
		case "<boolean>":
			return "bool", ""
		case "<date>":
			b.file.NeedsTime = true
			return "time.Time", ""
		case "<data>":
			return "[]byte", ""
		case "<dictionary>":
			if len(item.Subkeys) == 0 {
				return "map[string]any", ""
			}
			itemName := naming.ExportName(item.Key)
			if itemName == "X" {
				itemName = fieldName + "Item"
			}
			n := b.nestedStructNamed(structName+itemName, item)
			return n, n
		case "<array>":
			inner, nestedInner := b.itemType(structName, fieldName+"Item", item)
			return "[]" + inner, nestedInner
		default:
			return "any", ""
		}
	}
	// Multiple subkeys: the items are dictionaries described inline.
	synthetic := spec.Key{Key: k.Key, Type: "<dictionary>", Subkeys: k.Subkeys}
	n := b.nestedStructNamed(structName+fieldName+"Item", &synthetic)
	return n, n
}

// nestedStruct emits (or references) the struct for a dictionary key.
func (b *builder) nestedStruct(defaultName string, k *spec.Key) string {
	return b.nestedStructNamed(defaultName, k)
}

func (b *builder) nestedStructNamed(defaultName string, k *spec.Key) string {
	name := defaultName
	if k.SubkeyType != "" {
		name = naming.ExportName(k.SubkeyType)
	}
	sig := shapeSignature(k.Subkeys)
	for i := 2; ; i++ {
		existing, taken := b.shared.byName[name]
		if !taken {
			break
		}
		if existing == sig {
			return name // identical shape already emitted: reuse it
		}
		name = defaultName + strconv.Itoa(i)
	}
	b.shared.byName[name] = sig

	comment := []string{name + " is the " + k.Key + " dictionary."}
	if c := naming.FirstLine(k.Content); c != "" {
		comment = append(comment, naming.WrapComment(c, commentWidth)...)
	}
	b.nested = append(b.nested, b.buildStruct(name, comment, k.Subkeys, "", ""))
	return name
}

// claimName reserves a package-unique type name for a non-reusable struct
// (payload or response), suffixing on collision.
func claimName(shared *SharedTypes, base, sig string) string {
	name := base
	for i := 2; ; i++ {
		if _, taken := shared.byName[name]; !taken {
			shared.byName[name] = sig
			return name
		}
		name = base + strconv.Itoa(i)
	}
}

// shapeSignature is a stable fingerprint of a subkey tree, used to decide
// whether two same-named nested types are actually the same type.
func shapeSignature(keys []spec.Key) string {
	body, err := json.Marshal(keys)
	if err != nil {
		return fmt.Sprintf("unmarshalable:%v", err)
	}
	return string(body)
}

func ptrIf(optional bool, t string) string {
	if optional {
		return "*" + t
	}
	return t
}

func tag(k *spec.Key) string {
	name := k.Key
	if !k.Required() {
		return fmt.Sprintf("plist:%q json:%q", name+",omitempty", name+",omitempty")
	}
	return fmt.Sprintf("plist:%q json:%q", name, name)
}

// checks produces the fully-resolved Validate statements for one field.
func (b *builder) checks(fieldName, goType string, k *spec.Key, nested, enumType string, enumMembers []string) []string {
	var lines []string
	optional := !k.Required()
	q := strconv.Quote(k.Key)

	// Presence of required non-scalar values (scalars are value fields and
	// always serialized).
	if k.Required() {
		switch {
		case strings.HasPrefix(goType, "[]"), strings.HasPrefix(goType, "map["), goType == "any":
			lines = append(lines,
				fmt.Sprintf("if p.%s == nil {", fieldName),
				fmt.Sprintf("\terrs = append(errs, validate.Required(%s))", q),
				"}")
		}
	}

	// Scalar constraint checks.
	scalar := scalarChecks(valueExpr(fieldName, goType), q, k, enumType, enumMembers)
	if len(scalar) > 0 {
		if optional && strings.HasPrefix(goType, "*") {
			lines = append(lines, fmt.Sprintf("if p.%s != nil {", fieldName))
			lines = append(lines, indent(scalar)...)
			lines = append(lines, "}")
		} else {
			lines = append(lines, scalar...)
		}
	}

	// Nested dictionaries validate recursively.
	if nested != "" && !strings.HasPrefix(goType, "[]") {
		inner := []string{
			fmt.Sprintf("if err := p.%s.Validate(); err != nil {", fieldName),
			fmt.Sprintf("\terrs = append(errs, validate.Nested(%s, err))", q),
			"}",
		}
		if strings.HasPrefix(goType, "*") {
			lines = append(lines, fmt.Sprintf("if p.%s != nil {", fieldName))
			lines = append(lines, indent(inner)...)
			lines = append(lines, "}")
		} else {
			lines = append(lines, inner...)
		}
	}

	// Arrays: cardinality, then per-item validation.
	if strings.HasPrefix(goType, "[]") && goType != "[]byte" {
		if r := k.Repetition; r != nil {
			lines = append(lines,
				fmt.Sprintf("if p.%s != nil {", fieldName),
				fmt.Sprintf("\tif err := validate.Repetition(%s, len(p.%s), %d, %d); err != nil {", q, fieldName, r.Min, r.Max),
				"\t\terrs = append(errs, err)",
				"\t}",
				"}")
		}
		if nested != "" && goType == "[]"+nested {
			lines = append(lines,
				fmt.Sprintf("for i := range p.%s {", fieldName),
				fmt.Sprintf("\tif err := p.%s[i].Validate(); err != nil {", fieldName),
				fmt.Sprintf("\t\terrs = append(errs, validate.Indexed(%s, i, err))", q),
				"\t}",
				"}")
		} else if nested == "" && len(k.Subkeys) == 1 {
			item := &k.Subkeys[0]
			itemChecks := scalarChecks("p."+fieldName+"[i]", q, item, "", nil)
			if len(itemChecks) > 0 && itemScalar(goType) {
				lines = append(lines, fmt.Sprintf("for i := range p.%s {", fieldName))
				lines = append(lines, indent(itemChecks)...)
				lines = append(lines, "}")
			}
		}
	}
	return lines
}

// itemScalar reports whether an array's Go type has scalar items that
// per-item checks can apply to.
func itemScalar(goType string) bool {
	switch strings.TrimPrefix(goType, "[]") {
	case "string", "int64", "float64":
		return true
	}
	return false
}

// valueExpr is the dereferenced access expression for a scalar field.
func valueExpr(fieldName, goType string) string {
	if strings.HasPrefix(goType, "*") {
		return "*p." + fieldName
	}
	return "p." + fieldName
}

// scalarChecks builds rangelist/range/format/subtype checks for a scalar
// key against expr. Enum-typed keys check membership against the typed
// constants; other string checks convert back to plain string.
func scalarChecks(expr, q string, k *spec.Key, enumType string, enumMembers []string) []string {
	var lines []string
	appendErr := func(call string) {
		lines = append(lines,
			"if err := "+call+"; err != nil {",
			"\terrs = append(errs, err)",
			"}")
	}
	switch k.Type {
	case "<string>":
		strExpr := expr
		if enumType != "" {
			strExpr = "string(" + expr + ")"
			appendErr(fmt.Sprintf("validate.InList(%s, %s, []%s{%s})", q, expr, enumType, strings.Join(enumMembers, ", ")))
		} else if list := stringList(k.RangeList); len(list) > 0 {
			appendErr(fmt.Sprintf("validate.InList(%s, %s, []string{%s})", q, expr, strings.Join(list, ", ")))
		}
		if k.Format != "" {
			appendErr(fmt.Sprintf("validate.Format(%s, %s, %s)", q, strExpr, strconv.Quote(k.Format)))
		}
		switch k.Subtype {
		case "<url>":
			appendErr(fmt.Sprintf("validate.URL(%s, %s)", q, strExpr))
		case "<hostname>":
			appendErr(fmt.Sprintf("validate.Hostname(%s, %s)", q, strExpr))
		case "<email>":
			appendErr(fmt.Sprintf("validate.Email(%s, %s)", q, strExpr))
		}
	case "<integer>":
		if enumType != "" {
			appendErr(fmt.Sprintf("validate.InList(%s, %s, []%s{%s})", q, expr, enumType, strings.Join(enumMembers, ", ")))
		} else if list := intList(k.RangeList); len(list) > 0 {
			appendErr(fmt.Sprintf("validate.InList(%s, %s, []int64{%s})", q, expr, strings.Join(list, ", ")))
		}
		if r := k.Range; r != nil && (r.Min != nil || r.Max != nil) {
			intExpr := expr
			if enumType != "" {
				intExpr = "int64(" + expr + ")"
			}
			appendErr(fmt.Sprintf("validate.IntRange(%s, %s, %s, %s)", q, intExpr, intPtrLit(r.Min), intPtrLit(r.Max)))
		}
	case "<real>":
		if r := k.Range; r != nil && (r.Min != nil || r.Max != nil) {
			appendErr(fmt.Sprintf("validate.FloatRange(%s, %s, %s, %s)", q, expr, floatPtrLit(r.Min), floatPtrLit(r.Max)))
		}
	}
	return lines
}

func stringList(vals []any) []string {
	var out []string
	for _, v := range vals {
		s, ok := v.(string)
		if !ok {
			return nil // mixed-type rangelist: skip the check
		}
		out = append(out, strconv.Quote(s))
	}
	return out
}

func intList(vals []any) []string {
	var out []string
	for _, v := range vals {
		switch n := v.(type) {
		case int64:
			out = append(out, strconv.FormatInt(n, 10))
		case float64:
			if n != float64(int64(n)) {
				return nil
			}
			out = append(out, strconv.FormatInt(int64(n), 10))
		default:
			return nil
		}
	}
	return out
}

func intPtrLit(v *float64) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("ptr.To(int64(%d))", int64(*v))
}

func floatPtrLit(v *float64) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("ptr.To(float64(%s))", strconv.FormatFloat(*v, 'g', -1, 64))
}

func indent(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = "\t" + l
	}
	return out
}

func usesValidate(lines []string) bool {
	for _, l := range lines {
		if strings.Contains(l, "validate.") {
			return true
		}
	}
	return false
}

// enumFor emits (or reuses) the typed allowed-values enum for a string or
// integer rangelist key, returning the type name and its constant names.
// Identical (name, value-set) pairs collapse to one shared type across
// specs; a name collision with different values gets a numeric suffix.
func (b *builder) enumFor(structName, fieldName string, k *spec.Key) (string, []string) {
	if len(k.RangeList) < 2 {
		return "", nil
	}
	var baseType string
	// memberSuffix and literal per value; nil on a mixed-type rangelist.
	var suffixes, literals []string
	switch k.Type {
	case "<string>":
		baseType = "string"
		for _, v := range k.RangeList {
			s, ok := v.(string)
			if !ok {
				return "", nil
			}
			suffixes = append(suffixes, naming.ConstName(s))
			literals = append(literals, strconv.Quote(s))
		}
	case "<integer>":
		baseType = "int64"
		for _, v := range k.RangeList {
			n, ok := intValue(v)
			if !ok {
				return "", nil
			}
			lit := strconv.FormatInt(n, 10)
			suffix := "Value" + lit
			if n < 0 {
				suffix = "ValueNeg" + strconv.FormatInt(-n, 10)
			}
			suffixes = append(suffixes, suffix)
			literals = append(literals, lit)
		}
	default:
		return "", nil
	}

	base := structName + fieldName
	sig := "enum:" + baseType + ":" + strings.Join(literals, "\x00")
	name := base
	for i := 2; ; i++ {
		existing, taken := b.shared.byName[name]
		if !taken {
			break
		}
		if existing == sig {
			return name, b.shared.enumMembers[name] // reuse the shared type
		}
		name = base + strconv.Itoa(i)
	}
	b.shared.byName[name] = sig

	block := view.EnumBlock{
		TypeName: name,
		BaseType: baseType,
		Comment:  "allowed values for " + structName + "." + fieldName + ".",
	}
	members := make([]string, 0, len(literals))
	seenLit := map[string]bool{}
	for i, lit := range literals {
		member := name + suffixes[i]
		for j := 2; b.usedConsts[member]; j++ {
			member = fmt.Sprintf("%s%s%d", name, suffixes[i], j)
		}
		b.usedConsts[member] = true
		members = append(members, member)
		block.Members = append(block.Members, view.EnumMember{Name: member, Literal: lit, Dup: seenLit[lit]})
		seenLit[lit] = true
	}
	b.shared.enumMembers[name] = members
	b.file.Enums = append(b.file.Enums, block)
	return name, members
}

// intValue extracts an integer rangelist entry.
func intValue(v any) (int64, bool) {
	switch n := v.(type) {
	case int64:
		return n, true
	case float64:
		if n == float64(int64(n)) {
			return int64(n), true
		}
	}
	return 0, false
}

// specComment builds the main struct's doc comment from the spec metadata.
func specComment(s *spec.Spec) []string {
	name := naming.StructName(s.TypeIdentifier(), s.Name)
	lines := []string{name + " — " + s.Title + "."}
	desc := s.Description
	if desc == "" {
		desc = naming.FirstLine(s.Payload.Content)
	}
	if d := naming.FirstLine(desc); d != "" {
		lines = append(lines, naming.WrapComment(d, commentWidth)...)
	}
	if avail := availability(s.Payload.SupportedOS); avail != "" {
		lines = append(lines, "", avail)
	}
	if s.Payload.Beta {
		lines = append(lines, "", "Beta: this payload may change incompatibly before final release.")
	}
	return lines
}

// availability renders "Supported: iOS 4.0+, macOS 10.7+." in fixed OS
// order, skipping unsupported OSes.
func availability(oses map[string]spec.OSSupport) string {
	order := []string{"iOS", "macOS", "tvOS", "visionOS", "watchOS"}
	var parts []string
	for _, os := range order {
		s, ok := oses[os]
		if !ok || s.Introduced == "" || s.Introduced == "n/a" {
			continue
		}
		p := os + " " + s.Introduced + "+"
		if s.Removed != "" {
			p = os + " " + s.Introduced + "–" + s.Removed
		}
		parts = append(parts, p)
	}
	if len(parts) == 0 {
		return ""
	}
	return "Supported: " + strings.Join(parts, ", ") + "."
}

// keyComment builds a field's doc comment.
func keyComment(k *spec.Key) []string {
	var lines []string
	if c := naming.FirstLine(k.Content); c != "" {
		lines = append(lines, naming.WrapComment(c, commentWidth)...)
	}
	var meta []string
	if k.Default != nil {
		meta = append(meta, fmt.Sprintf("Default: %v.", k.Default))
	}
	if len(k.AssetTypes) > 0 {
		meta = append(meta, "Allowed asset types: "+strings.Join(k.AssetTypes, ", ")+".")
	}
	if len(meta) > 0 {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, meta...)
	}
	return lines
}
