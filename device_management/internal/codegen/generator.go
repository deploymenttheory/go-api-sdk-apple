// Package codegen orchestrates the offline device-management code
// generator: it loads the committed spec snapshots, clears previously
// generated output (identified by the DO-NOT-EDIT header), builds view
// models (build), renders them through the template firewall (render) and
// assembles files (fileasm). Everything is deterministic from
// metadata/specs.
package codegen

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen/build"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen/fileasm"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen/render"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/codegen/view"
	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/spec"
)

const modulePath = "github.com/deploymenttheory/go-sdk-appleservices/device_management"

// family maps an upstream spec category onto a generated package.
type family struct {
	categoryPrefix string
	dir            string // under the device_management root
	pkg            string
	kind           build.Kind
	iface          string // registry interface type
	ifaceImport    string
	mapName        string
	registryDoc    string
}

var families = []family{
	{
		categoryPrefix: "mdm/commands",
		dir:            "mdm/commands", pkg: "commands",
		kind:  build.KindCommand,
		iface: "mdm.CommandPayload", ifaceImport: modulePath + "/mdm",
		mapName:     "ByRequestType",
		registryDoc: "ByRequestType maps MDM RequestType identifiers to payload factories.",
	},
	{
		categoryPrefix: "mdm/profiles",
		dir:            "mdm/profiles", pkg: "profiles",
		kind:  build.KindProfile,
		iface: "mdm.ProfilePayload", ifaceImport: modulePath + "/mdm",
		mapName:     "ByPayloadType",
		registryDoc: "ByPayloadType maps profile PayloadType identifiers to payload factories.",
	},
	{
		categoryPrefix: "declarative/declarations/configurations",
		dir:            "ddm/configurations", pkg: "configurations",
		kind:  build.KindDeclaration,
		iface: "ddm.DeclarationPayload", ifaceImport: modulePath + "/ddm",
		mapName:     "ByDeclarationType",
		registryDoc: "ByDeclarationType maps declaration type identifiers to payload factories.",
	},
	{
		categoryPrefix: "declarative/declarations/assets",
		dir:            "ddm/assets", pkg: "assets",
		kind:  build.KindDeclaration,
		iface: "ddm.DeclarationPayload", ifaceImport: modulePath + "/ddm",
		mapName:     "ByDeclarationType",
		registryDoc: "ByDeclarationType maps declaration type identifiers to payload factories.",
	},
	{
		categoryPrefix: "declarative/declarations/activations",
		dir:            "ddm/activations", pkg: "activations",
		kind:  build.KindDeclaration,
		iface: "ddm.DeclarationPayload", ifaceImport: modulePath + "/ddm",
		mapName:     "ByDeclarationType",
		registryDoc: "ByDeclarationType maps declaration type identifiers to payload factories.",
	},
	{
		categoryPrefix: "declarative/declarations/management",
		dir:            "ddm/management", pkg: "management",
		kind:  build.KindDeclaration,
		iface: "ddm.DeclarationPayload", ifaceImport: modulePath + "/ddm",
		mapName:     "ByDeclarationType",
		registryDoc: "ByDeclarationType maps declaration type identifiers to payload factories.",
	},
}

// Run generates the SDK surface from metadataDir into outDir (the
// device_management root).
func Run(metadataDir, outDir string) error {
	specs, err := loadSpecs(metadataDir)
	if err != nil {
		return err
	}

	total := 0
	for i := range families {
		fam := &families[i]
		dir := filepath.Join(outDir, filepath.FromSlash(fam.dir))
		// Clear before re-emitting: a run never mixes fresh output with
		// leftovers from a previous generation (renamed specs, removed
		// specs, layout changes). Only header-marked files are removed;
		// hand-written files in generated packages survive.
		if err := clearGenerated(dir); err != nil {
			return fmt.Errorf("%s: clear: %w", fam.dir, err)
		}
		n, err := emitFamily(fam, specs, outDir)
		if err != nil {
			return fmt.Errorf("%s: %w", fam.dir, err)
		}
		total += n
	}
	fmt.Printf("generated %d specs across %d families -> %s\n", total, len(families), outDir)
	return nil
}

func loadSpecs(metadataDir string) ([]*spec.Spec, error) {
	var specs []*spec.Spec
	err := filepath.WalkDir(metadataDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(p, ".json") || d.Name() == "PROVENANCE.json" {
			return err
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		var s spec.Spec
		if err := json.Unmarshal(body, &s); err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		specs = append(specs, &s)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(specs) == 0 {
		return nil, fmt.Errorf("no spec snapshots in %s", metadataDir)
	}
	sort.Slice(specs, func(i, j int) bool {
		if specs[i].Category != specs[j].Category {
			return specs[i].Category < specs[j].Category
		}
		return specs[i].Name < specs[j].Name
	})
	return specs, nil
}

func emitFamily(fam *family, specs []*spec.Spec, outDir string) (int, error) {
	dir := filepath.Join(outDir, filepath.FromSlash(fam.dir))
	shared := build.NewSharedTypes()
	usedConsts := map[string]bool{}
	usedFiles := map[string]bool{}
	registered := map[string]bool{}
	var reg view.Registry

	count := 0
	for _, s := range specs {
		if s.Category != fam.categoryPrefix && !strings.HasPrefix(s.Category, fam.categoryPrefix+"/") {
			continue
		}
		f, err := build.File(s, fam.pkg, fam.kind, shared, usedConsts)
		if err != nil {
			return count, fmt.Errorf("%s: %w", s.Name, err)
		}

		// One spec emits up to three files, separating construct kinds:
		// <name>.go (struct declarations), <name>_functions.go (methods),
		// <name>_enums.go (allowed-value constants).
		var decls, funcs, enums strings.Builder
		for i := range f.Structs {
			frag, err := render.StructDecl(&f.Structs[i])
			if err != nil {
				return count, fmt.Errorf("%s: %w", s.Name, err)
			}
			decls.WriteString(frag)
			frag, err = render.StructFuncs(&f.Structs[i])
			if err != nil {
				return count, fmt.Errorf("%s: %w", s.Name, err)
			}
			funcs.WriteString(frag)
		}
		for i := range f.Enums {
			frag, err := render.EnumBlock(&f.Enums[i])
			if err != nil {
				return count, fmt.Errorf("%s: %w", s.Name, err)
			}
			enums.WriteString(frag)
		}

		fileName := f.Name
		for i := 2; usedFiles[fileName]; i++ {
			fileName = fmt.Sprintf("%s%d", f.Name, i)
		}
		usedFiles[fileName] = true

		if err := writeGen(filepath.Join(dir, fileName+".go"), fam.pkg, imports(decls.String()), decls.String()); err != nil {
			return count, fmt.Errorf("%s: %w", s.Name, err)
		}
		if err := writeGen(filepath.Join(dir, fileName+"_functions.go"), fam.pkg, imports(funcs.String()), funcs.String()); err != nil {
			return count, fmt.Errorf("%s: %w", s.Name, err)
		}
		if enums.Len() > 0 {
			if err := writeGen(filepath.Join(dir, fileName+"_enums.go"), fam.pkg, imports(enums.String()), enums.String()); err != nil {
				return count, fmt.Errorf("%s: %w", s.Name, err)
			}
		}

		// A few upstream specs share a wire identifier (the com.apple.MCX
		// profile family); the first spec in sorted order wins the registry
		// slot, later ones remain reachable as plain structs.
		if id := s.TypeIdentifier(); id != "" && fam.kind != build.KindPlain && !registered[id] {
			registered[id] = true
			reg.Entries = append(reg.Entries, view.RegistryEntry{
				Identifier: id,
				StructName: f.MainStruct,
			})
		}
		count++
	}
	if count == 0 {
		return 0, nil
	}

	// Family registry.
	reg.PackageName = fam.pkg
	reg.MapName = fam.mapName
	reg.IfaceType = fam.iface
	reg.CommentLines = []string{fam.registryDoc, "Factories return zero payloads ready to populate or decode into."}
	sort.Slice(reg.Entries, func(i, j int) bool { return reg.Entries[i].Identifier < reg.Entries[j].Identifier })
	regBody, err := render.Registry(&reg)
	if err != nil {
		return count, err
	}
	regImports := []fileasm.Import{{Path: fam.ifaceImport}}
	if err := writeGen(filepath.Join(dir, "registry.go"), fam.pkg, regImports, regBody); err != nil {
		return count, err
	}
	return count, nil
}

// imports derives a generated file's imports from its rendered body.
func imports(body string) []fileasm.Import {
	var out []fileasm.Import
	if strings.Contains(body, "errors.Join") {
		out = append(out, fileasm.Import{Path: "errors"})
	}
	if strings.Contains(body, "fmt.") {
		out = append(out, fileasm.Import{Path: "fmt"})
	}
	if strings.Contains(body, "time.Time") {
		out = append(out, fileasm.Import{Path: "time"})
	}
	if strings.Contains(body, "validate.") {
		out = append(out, fileasm.Import{Path: modulePath + "/validate"})
	}
	if strings.Contains(body, "ptr.To(") {
		out = append(out, fileasm.Import{Path: modulePath + "/ptr"})
	}
	return out
}

func writeGen(path, pkg string, imps []fileasm.Import, body string) error {
	return fileasm.WriteFile(path, pkg, "", imps, body)
}

// clearGenerated removes every generated file (identified by the header
// marker) under dir, then removes emptied directories. Hand-written files
// are never touched.
func clearGenerated(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	var dirs []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			dirs = append(dirs, path)
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		head := make([]byte, len(fileasm.Header))
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		n, _ := f.Read(head)
		f.Close()
		if string(head[:n]) == fileasm.Header {
			return os.Remove(path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, d := range dirs {
		if entries, err := os.ReadDir(d); err == nil && len(entries) == 0 {
			os.Remove(d)
		}
	}
	return nil
}
