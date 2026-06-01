package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, flexspecDir)
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `specs_dir: my-specs
always_one_shot: true
`
	if err := os.WriteFile(filepath.Join(base, configFile), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.SpecsDir != "my-specs" {
		t.Errorf("SpecsDir = %q", cfg.SpecsDir)
	}
	if !cfg.AlwaysOneShot {
		t.Error("AlwaysOneShot = false, want true")
	}
}

func TestLoad_missingConfig(t *testing.T) {
	root := t.TempDir()
	_, err := Load(root)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDisplayEntries(t *testing.T) {
	entries := DisplayEntries(Config{
		SpecsDir:      "specs",
		AlwaysOneShot: true,
		SpecTemplate:  "",
	})
	if len(entries) != 3 {
		t.Fatalf("len = %d", len(entries))
	}
	if entries[0].Key != "specs_dir" || entries[0].Value != "specs" {
		t.Errorf("specs_dir entry = %+v", entries[0])
	}
	if entries[1].Value != "true" {
		t.Errorf("always_one_shot = %q", entries[1].Value)
	}
	if entries[2].Value != "-" {
		t.Errorf("spec_template = %q, want -", entries[2].Value)
	}
}

func TestJSONDocumentFromConfig(t *testing.T) {
	doc := JSONDocumentFromConfig(Config{
		SpecsDir:      "x",
		AlwaysOneShot: false,
		SpecTemplate:  "simple",
	})
	if doc.SpecsDir != "x" || doc.AlwaysOneShot != false || doc.SpecTemplate != "simple" {
		t.Errorf("got %+v", doc)
	}
}

func TestLoad_emptySpecsDir(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, flexspecDir)
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, configFile), []byte("specs_dir: \"\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(root)
	if err == nil {
		t.Fatal("expected error for empty specs_dir")
	}
}
