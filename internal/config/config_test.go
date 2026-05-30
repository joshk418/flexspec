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
