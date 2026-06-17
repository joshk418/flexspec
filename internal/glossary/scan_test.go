package glossary

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan_excludesKnown(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "spec.md"), []byte("FlexSpec and Cobra are used. FlexSpec again."), 0o644); err != nil {
		t.Fatal(err)
	}

	known := []Entry{{Term: "FlexSpec", Definition: "the tool"}}
	candidates, err := Scan(root, known, ScanOptions{Max: 50})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range candidates {
		if c.Term == "FlexSpec" {
			t.Errorf("known term FlexSpec should be excluded, got candidate %+v", c)
		}
	}
}

func TestScan_excludesKeywords(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "code.go"), []byte("func Return True False Null String"), 0o644); err != nil {
		t.Fatal(err)
	}

	candidates, err := Scan(root, nil, ScanOptions{Max: 50})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range candidates {
		switch c.Term {
		case "Return", "True", "False", "Null", "String":
			t.Errorf("keyword %q should be excluded, got candidate %+v", c.Term, c)
		}
	}
}

func TestScan_ranksByCount(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.md"), []byte("Alpha Alpha Alpha Beta"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "b.md"), []byte("Beta"), 0o644); err != nil {
		t.Fatal(err)
	}

	candidates, err := Scan(root, nil, ScanOptions{Max: 50})
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) < 2 {
		t.Fatalf("got %d candidates, want >= 2", len(candidates))
	}
	if candidates[0].Term != "Alpha" || candidates[0].Count != 3 {
		t.Errorf("top candidate = %+v, want Alpha count 3", candidates[0])
	}
}

func TestScan_respectsMax(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.md"), []byte("Alpha Beta Gamma Delta Epsilon"), 0o644); err != nil {
		t.Fatal(err)
	}

	candidates, err := Scan(root, nil, ScanOptions{Max: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(candidates) > 2 {
		t.Errorf("got %d candidates, want <= 2", len(candidates))
	}
}

func TestScan_skipsVendoredDirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "node_modules"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "node_modules", "lib.js"), []byte("IgnoreThisTerm"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "app.md"), []byte("KeepThisTerm"), 0o644); err != nil {
		t.Fatal(err)
	}

	candidates, err := Scan(root, nil, ScanOptions{Max: 50})
	if err != nil {
		t.Fatal(err)
	}
	terms := map[string]bool{}
	for _, c := range candidates {
		terms[c.Term] = true
	}
	if terms["IgnoreThisTerm"] {
		t.Error("node_modules content should be skipped")
	}
	if !terms["KeepThisTerm"] {
		t.Error("root content should be scanned")
	}
}
