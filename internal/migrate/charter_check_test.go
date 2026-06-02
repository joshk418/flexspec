package migrate

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/joshk418/flexspec/internal/config"
)

func TestCharterCheck_missingSection_noWrite(t *testing.T) {
	root := t.TempDir()
	flexDir := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(flexDir, 0o755); err != nil {
		t.Fatal(err)
	}
	charter := "## 1. Product overview\n\nOnly one section.\n"
	path := filepath.Join(flexDir, "charter.md")
	if err := os.WriteFile(path, []byte(charter), 0o644); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)

	tmpl := fstest.MapFS{
		"charter.md": {Data: []byte("## 1. Product overview\n\n## 2. Vision and goals\n")},
	}
	m := &charterCheckMigration{templates: tmpl}
	cfg := config.Config{SpecsDir: "specs"}

	changes, err := m.Detect(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) == 0 {
		t.Fatal("expected report for missing section")
	}
	if _, err := m.Apply(root, cfg); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Fatal("charter must not be modified by Apply")
	}
}
