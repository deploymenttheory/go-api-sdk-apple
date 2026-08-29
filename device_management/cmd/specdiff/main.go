// Command specdiff renders a semantic, human-readable markdown report of
// what changed between two spec snapshot trees — the review surface for
// Apple's moving schema target. The weekly spec-update workflow uses it as
// the pull-request body.
//
//	go run ./device_management/cmd/specdiff -old <dir> -new <dir>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/spec"
)

func main() {
	oldDir := flag.String("old", "", "previous snapshot directory")
	newDir := flag.String("new", filepath.Join("device_management", "metadata", "specs"), "current snapshot directory")
	flag.Parse()
	if *oldDir == "" {
		fmt.Fprintln(os.Stderr, "specdiff: -old is required")
		os.Exit(2)
	}
	report, changed, err := diff(*oldDir, *newDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "specdiff:", err)
		os.Exit(1)
	}
	fmt.Print(report)
	if !changed {
		os.Exit(0)
	}
}

func diff(oldDir, newDir string) (string, bool, error) {
	oldSpecs, err := load(oldDir)
	if err != nil {
		return "", false, err
	}
	newSpecs, err := load(newDir)
	if err != nil {
		return "", false, err
	}

	var b strings.Builder
	b.WriteString("## Device-management spec changes\n\n")
	changed := false

	var added, removed, common []string
	for k := range newSpecs {
		if _, ok := oldSpecs[k]; ok {
			common = append(common, k)
		} else {
			added = append(added, k)
		}
	}
	for k := range oldSpecs {
		if _, ok := newSpecs[k]; !ok {
			removed = append(removed, k)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(common)

	if len(added) > 0 {
		changed = true
		b.WriteString("### New specs\n\n")
		for _, k := range added {
			fmt.Fprintf(&b, "- `%s` — %s\n", k, newSpecs[k].Title)
		}
		b.WriteString("\n")
	}
	if len(removed) > 0 {
		changed = true
		b.WriteString("### Removed specs\n\n")
		for _, k := range removed {
			fmt.Fprintf(&b, "- `%s` — %s\n", k, oldSpecs[k].Title)
		}
		b.WriteString("\n")
	}

	var modified []string
	details := map[string][]string{}
	for _, k := range common {
		lines := diffSpec(oldSpecs[k], newSpecs[k])
		if len(lines) > 0 {
			modified = append(modified, k)
			details[k] = lines
		}
	}
	if len(modified) > 0 {
		changed = true
		b.WriteString("### Changed specs\n\n")
		for _, k := range modified {
			fmt.Fprintf(&b, "#### `%s`\n\n", k)
			for _, l := range details[k] {
				fmt.Fprintf(&b, "- %s\n", l)
			}
			b.WriteString("\n")
		}
	}
	if !changed {
		b.WriteString("No spec changes.\n")
	}
	return b.String(), changed, nil
}

func load(dir string) (map[string]*spec.Spec, error) {
	out := map[string]*spec.Spec{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
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
		out[s.Category+"/"+s.Name] = &s
		return nil
	})
	return out, err
}

func diffSpec(oldS, newS *spec.Spec) []string {
	var lines []string
	if oldS.Payload.Apply != newS.Payload.Apply {
		lines = append(lines, fmt.Sprintf("apply: `%s` → `%s`", orNone(oldS.Payload.Apply), orNone(newS.Payload.Apply)))
	}
	if id0, id1 := oldS.TypeIdentifier(), newS.TypeIdentifier(); id0 != id1 {
		lines = append(lines, fmt.Sprintf("type identifier: `%s` → `%s`", id0, id1))
	}
	lines = append(lines, diffOS("payload", oldS.Payload.SupportedOS, newS.Payload.SupportedOS)...)
	lines = append(lines, diffKeys("", oldS.PayloadKeys, newS.PayloadKeys)...)
	lines = append(lines, diffKeys("response:", oldS.ResponseKeys, newS.ResponseKeys)...)
	return lines
}

func diffOS(scope string, o, n map[string]spec.OSSupport) []string {
	var lines []string
	oses := map[string]bool{}
	for k := range o {
		oses[k] = true
	}
	for k := range n {
		oses[k] = true
	}
	names := make([]string, 0, len(oses))
	for k := range oses {
		names = append(names, k)
	}
	sort.Strings(names)
	for _, os := range names {
		a, b := o[os], n[os]
		if a.Introduced != b.Introduced {
			lines = append(lines, fmt.Sprintf("%s %s introduced: `%s` → `%s`", scope, os, orNone(a.Introduced), orNone(b.Introduced)))
		}
		if a.Deprecated != b.Deprecated {
			lines = append(lines, fmt.Sprintf("%s %s deprecated: `%s` → `%s`", scope, os, orNone(a.Deprecated), orNone(b.Deprecated)))
		}
		if a.Removed != b.Removed {
			lines = append(lines, fmt.Sprintf("%s %s removed: `%s` → `%s`", scope, os, orNone(a.Removed), orNone(b.Removed)))
		}
	}
	return lines
}

func diffKeys(prefix string, o, n []spec.Key) []string {
	var lines []string
	oldByKey := index(o)
	newByKey := index(n)

	var names []string
	seen := map[string]bool{}
	for _, k := range o {
		if !seen[k.Key] {
			seen[k.Key] = true
			names = append(names, k.Key)
		}
	}
	for _, k := range n {
		if !seen[k.Key] {
			seen[k.Key] = true
			names = append(names, k.Key)
		}
	}
	sort.Strings(names)

	for _, name := range names {
		path := prefix + name
		ok, oldK := oldByKey[name] != nil, oldByKey[name]
		nk, newK := newByKey[name] != nil, newByKey[name]
		switch {
		case ok && !nk:
			lines = append(lines, fmt.Sprintf("removed key `%s` (%s)", path, oldK.Type))
		case !ok && nk:
			presence := newK.Presence
			if presence == "" {
				presence = "optional"
			}
			lines = append(lines, fmt.Sprintf("added key `%s` (%s, %s)", path, newK.Type, presence))
		default:
			lines = append(lines, diffKey(path, oldK, newK)...)
		}
	}
	return lines
}

func index(keys []spec.Key) map[string]*spec.Key {
	out := map[string]*spec.Key{}
	for i := range keys {
		if _, dup := out[keys[i].Key]; !dup {
			out[keys[i].Key] = &keys[i]
		}
	}
	return out
}

func diffKey(path string, o, n *spec.Key) []string {
	var lines []string
	if o.Type != n.Type {
		lines = append(lines, fmt.Sprintf("key `%s` type: `%s` → `%s`", path, o.Type, n.Type))
	}
	if o.Presence != n.Presence {
		lines = append(lines, fmt.Sprintf("key `%s` presence: `%s` → `%s`", path, orNone(o.Presence), orNone(n.Presence)))
	}
	if !reflect.DeepEqual(o.RangeList, n.RangeList) {
		lines = append(lines, fmt.Sprintf("key `%s` allowed values: `%v` → `%v`", path, o.RangeList, n.RangeList))
	}
	if !reflect.DeepEqual(o.Range, n.Range) {
		lines = append(lines, fmt.Sprintf("key `%s` range changed", path))
	}
	if o.Format != n.Format {
		lines = append(lines, fmt.Sprintf("key `%s` format: `%s` → `%s`", path, orNone(o.Format), orNone(n.Format)))
	}
	if !reflect.DeepEqual(o.Default, n.Default) {
		lines = append(lines, fmt.Sprintf("key `%s` default: `%v` → `%v`", path, o.Default, n.Default))
	}
	lines = append(lines, diffKeys(path+"/", o.Subkeys, n.Subkeys)...)
	return lines
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}
