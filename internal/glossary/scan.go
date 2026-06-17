package glossary

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Candidate is one discovered term from a project scan.
type Candidate struct {
	Term     string `json:"term"`
	Count    int    `json:"count"`
	Category string `json:"category,omitempty"`
}

// scanTermRE matches PascalCase / CamelCase identifiers of length >= 4.
var scanTermRE = regexp.MustCompile(`[A-Z][a-zA-Z]{3,}`)

// scanExclusions are common language keywords and primitives excluded from
// discovery results.
var scanExclusions = map[string]bool{
	"true": true, "false": true, "null": true, "error": true, "string": true,
	"int": true, "return": true, "func": true, "function": true, "class": true,
	"interface": true, "struct": true, "package": true, "import": true,
	"var": true, "const": true, "let": true, "this": true, "self": true,
	"static": true, "public": true, "private": true,
}

// ScanOptions configures a project term scan.
type ScanOptions struct {
	Max int
}

// Scan walks root for candidate project-specific terms, excluding known glossary entries. Ranked by frequency.
func Scan(root string, known []Entry, opts ScanOptions) ([]Candidate, error) {
	knownSet := make(map[string]bool, len(known))
	for _, e := range known {
		knownSet[strings.ToLower(e.Term)] = true
		for _, a := range e.Aliases {
			knownSet[strings.ToLower(a)] = true
		}
	}

	counts := make(map[string]int)
	add := func(text string) {
		for _, m := range scanTermRE.FindAllString(text, -1) {
			lower := strings.ToLower(m)
			if scanExclusions[lower] {
				continue
			}
			if knownSet[lower] {
				continue
			}
			counts[m]++
		}
	}

	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "dist" || name == "build" || name == "target" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		name := d.Name()
		if !shouldScan(name) {
			return nil
		}
		data, err := readFile(p)
		if err != nil {
			return nil
		}
		add(string(data))
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := make([]Candidate, 0, len(counts))
	for term, count := range counts {
		out = append(out, Candidate{Term: term, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Term < out[j].Term
	})
	if opts.Max > 0 && len(out) > opts.Max {
		out = out[:opts.Max]
	}
	return out, nil
}

func shouldScan(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".md", ".markdown", ".go", ".ts", ".tsx", ".js", ".jsx", ".py",
		".rs", ".rb", ".java", ".kt", ".swift", ".yaml", ".yml", ".json",
		".toml", ".txt":
		return true
	}
	return false
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
