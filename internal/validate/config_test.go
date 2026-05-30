package validate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func TestLoadConfig_missing(t *testing.T) {
	root := t.TempDir()
	_, findings, ok := LoadConfig(root)
	if ok {
		t.Fatal("expected config load failure")
	}
	if len(findings) != 1 || findings[0].Rule != "config.missing" {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestLoadConfig_invalidYAML(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	path := filepath.Join(root, ".flexspec", "config.yaml")
	if err := os.WriteFile(path, []byte("specs_dir: [\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, findings, ok := LoadConfig(root)
	if ok {
		t.Fatal("expected failure")
	}
	if len(findings) != 1 || findings[0].Rule != "config.parse" {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestLoadConfig_emptySpecsDir(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	path := filepath.Join(root, ".flexspec", "config.yaml")
	if err := os.WriteFile(path, []byte("specs_dir: \"\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, findings, ok := LoadConfig(root)
	if ok {
		t.Fatal("expected failure")
	}
	if len(findings) != 1 || findings[0].Rule != "config.specs_dir" {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestLoadConfig_badSpecTemplate(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	path := filepath.Join(root, ".flexspec", "config.yaml")
	if err := os.WriteFile(path, []byte("specs_dir: specs\nspec_template: huge\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, findings, ok := LoadConfig(root)
	if ok {
		t.Fatal("expected failure")
	}
	if len(findings) != 1 || findings[0].Rule != "config.spec_template" {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckConfig_ok(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	cfg, _, ok := LoadConfig(root)
	if !ok {
		t.Fatal("expected valid config")
	}
	findings := CheckConfig(root, cfg, Options{})
	if len(findings) != 0 {
		t.Fatalf("unexpected findings: %+v", findings)
	}
}

func testConfig(t *testing.T, root string) config.Config {
	t.Helper()
	cfg, _, ok := LoadConfig(root)
	if !ok {
		t.Fatal("expected valid config")
	}
	return cfg
}
