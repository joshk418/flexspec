package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSave(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, flexspecDir)
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := Config{SpecsDir: "specs", AlwaysOneShot: true}
	if err := Save(root, cfg); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !loaded.AlwaysOneShot || loaded.SpecsDir != "specs" {
		t.Fatalf("loaded = %+v", loaded)
	}
}

func TestSave_invalidTemplate(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, flexspecDir)
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	err := Save(root, Config{SpecsDir: "specs", SpecTemplate: "invalid"})
	if err == nil {
		t.Fatal("expected error")
	}
}
