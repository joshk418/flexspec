package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckFlexspec_missingCharter(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	if err := os.Remove(filepath.Join(root, ".flexspec", "charter.md")); err != nil {
		t.Fatal(err)
	}
	findings := CheckFlexspec(root, testConfig(t, root), Options{})
	if len(findings) != 1 || findings[0].Rule != "charter.missing" {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckFlexspec_missingTemplate(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	if err := os.Remove(filepath.Join(root, ".flexspec", "templates", "flexspec-simple.md")); err != nil {
		t.Fatal(err)
	}
	findings := CheckFlexspec(root, testConfig(t, root), Options{})
	if len(findings) != 1 || findings[0].Rule != "templates.missing" {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckFlexspec_ok(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	findings := CheckFlexspec(root, testConfig(t, root), Options{})
	if len(findings) != 0 {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckFlexspec_malformedGlossary(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	path := filepath.Join(root, ".flexspec", "glossary.yaml")
	if err := os.WriteFile(path, []byte("not: yaml: [broken"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings := CheckFlexspec(root, testConfig(t, root), Options{})
	if len(findings) != 1 || findings[0].Rule != "glossary.invalid" {
		t.Fatalf("findings = %+v", findings)
	}
}
