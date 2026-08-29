// Command fetchspec is the acquisition stage of the device-management SDK
// pipeline. It downloads Apple's schema repo (apple/device-management) at a
// pinned commit, parses every YAML spec, and writes the committed snapshot
// tree (metadata/specs/<category>/<name>.json) plus provenance.
//
// Codegen (cmd/gendm) is offline and deterministic from the committed
// snapshots; moving to a new upstream commit is a deliberate, reviewed act
// (cmd/specdiff renders what changed).
//
//	go run ./device_management/cmd/fetchspec              # pinned commit
//	go run ./device_management/cmd/fetchspec -discover    # latest upstream commit
//	go run ./device_management/cmd/fetchspec -dir <path>  # offline: parse a local checkout
package main

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/deploymenttheory/go-sdk-appleservices/device_management/internal/spec"
)

// Pinned upstream state. Bump commit (and re-run) to adopt a new drop.
const (
	upstreamRepo  = "apple/device-management"
	defaultBranch = "release"
	defaultCommit = "67045e2fa06f528b196c01edee6a8bf88b844beb"
)

// specDirs are the upstream directories that contain schema specs.
var specDirs = []string{"declarative/", "mdm/", "other/"}

func main() {
	commit := flag.String("commit", defaultCommit, "upstream commit SHA to fetch")
	dir := flag.String("dir", "", "parse a local checkout instead of downloading (offline)")
	discover := flag.Bool("discover", false, "resolve the latest commit on the upstream release branch and fetch it")
	fetched := flag.String("fetched", "", "fetch date YYYY-MM-DD recorded in provenance (default: today UTC)")
	outDir := flag.String("out", filepath.Join("device_management", "metadata", "specs"), "snapshot output directory")
	flag.Parse()

	if *fetched == "" {
		*fetched = time.Now().UTC().Format("2006-01-02")
	}
	if *discover {
		latest, err := latestCommit()
		if err != nil {
			fmt.Fprintln(os.Stderr, "fetchspec: discover:", err)
			os.Exit(1)
		}
		if latest != *commit {
			fmt.Printf("discovered new upstream commit: %s\n", latest)
			*commit = latest
		} else {
			fmt.Println("pinned commit is current")
		}
	}
	if err := run(*commit, *dir, *fetched, *outDir); err != nil {
		fmt.Fprintln(os.Stderr, "fetchspec:", err)
		os.Exit(1)
	}
}

func run(commit, dir, fetched, outDir string) error {
	var specs map[string]*spec.Spec
	var digest string
	var err error
	if dir != "" {
		specs, err = parseDir(dir)
	} else {
		specs, digest, err = download(commit)
	}
	if err != nil {
		return err
	}
	if len(specs) == 0 {
		return fmt.Errorf("no YAML specs found")
	}

	written := map[string]bool{}
	rels := make([]string, 0, len(specs))
	for rel := range specs {
		rels = append(rels, rel)
	}
	sort.Strings(rels)
	for _, rel := range rels {
		p := filepath.Join(outDir, filepath.FromSlash(rel)+".json")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		body, err := json.MarshalIndent(specs[rel], "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(p, append(body, '\n'), 0o644); err != nil {
			return err
		}
		abs, _ := filepath.Abs(p)
		written[abs] = true
	}

	prov := spec.Provenance{
		Source:  "https://github.com/" + upstreamRepo,
		Ref:     defaultBranch,
		Commit:  commit,
		SHA256:  digest,
		Fetched: fetched,
	}
	provBody, err := json.MarshalIndent(prov, "", "  ")
	if err != nil {
		return err
	}
	provPath := filepath.Join(outDir, "PROVENANCE.json")
	if err := os.WriteFile(provPath, append(provBody, '\n'), 0o644); err != nil {
		return err
	}
	abs, _ := filepath.Abs(provPath)
	written[abs] = true

	if err := pruneStale(outDir, written); err != nil {
		return err
	}
	fmt.Printf("fetched %s@%s: %d specs -> %s\n", upstreamRepo, commit[:12], len(specs), outDir)
	return nil
}

// download fetches the repo archive for a commit and parses its specs.
func download(commit string) (map[string]*spec.Spec, string, error) {
	url := fmt.Sprintf("https://codeload.github.com/%s/zip/%s", upstreamRepo, commit)
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("download %s: HTTP %d", url, resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	sum := sha256.Sum256(data)

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, "", fmt.Errorf("open archive: %w", err)
	}
	out := map[string]*spec.Spec{}
	for _, entry := range zr.File {
		if entry.FileInfo().IsDir() {
			continue
		}
		// Strip the "<repo>-<commit>/" archive prefix.
		rel := entry.Name
		if i := strings.IndexByte(rel, '/'); i >= 0 {
			rel = rel[i+1:]
		}
		if !isSpecPath(rel) {
			continue
		}
		rc, err := entry.Open()
		if err != nil {
			return nil, "", err
		}
		body, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, "", err
		}
		if err := addSpec(out, body, rel); err != nil {
			return nil, "", err
		}
	}
	return out, hex.EncodeToString(sum[:]), nil
}

// parseDir parses specs from a local checkout of the upstream repo.
func parseDir(dir string) (map[string]*spec.Spec, error) {
	out := map[string]*spec.Spec{}
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !isSpecPath(rel) {
			return nil
		}
		body, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return addSpec(out, body, rel)
	})
	return out, err
}

func addSpec(out map[string]*spec.Spec, body []byte, rel string) error {
	s, err := spec.Parse(body, rel)
	if err != nil {
		return err
	}
	key := path.Join(s.Category, s.Name)
	if _, clash := out[key]; clash {
		return fmt.Errorf("snapshot collision: %q", key)
	}
	out[key] = s
	return nil
}

// isSpecPath reports whether a repo-relative path is a schema spec.
func isSpecPath(rel string) bool {
	if !strings.HasSuffix(rel, ".yaml") {
		return false
	}
	for _, d := range specDirs {
		if strings.HasPrefix(rel, d) {
			return true
		}
	}
	return false
}

// latestCommit resolves the newest commit SHA on the upstream branch.
func latestCommit() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s", upstreamRepo, defaultBranch)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github.sha")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET %s: HTTP %d", url, resp.StatusCode)
	}
	sha, err := io.ReadAll(io.LimitReader(resp.Body, 128))
	if err != nil {
		return "", err
	}
	s := strings.TrimSpace(string(sha))
	if len(s) != 40 {
		return "", fmt.Errorf("unexpected commit response %q", s)
	}
	return s, nil
}

// pruneStale removes snapshot files this run did not write, then empties
// abandoned directories.
func pruneStale(outDir string, written map[string]bool) error {
	var dirs []string
	err := filepath.WalkDir(outDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			dirs = append(dirs, p)
			return nil
		}
		if !strings.HasSuffix(p, ".json") {
			return nil
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		if !written[abs] {
			return os.Remove(p)
		}
		return nil
	})
	if err != nil {
		return err
	}
	sort.Slice(dirs, func(i, j int) bool { return len(dirs[i]) > len(dirs[j]) })
	for _, dir := range dirs {
		if entries, err := os.ReadDir(dir); err == nil && len(entries) == 0 {
			os.Remove(dir)
		}
	}
	return nil
}
