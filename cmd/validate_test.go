package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateCmd_cleanProject(t *testing.T) {
	root := t.TempDir()
	writeValidateFixture(t, root)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	var out bytes.Buffer
	validateCmd.SetOut(&out)
	validateCmd.SetErr(&out)
	if err := validateCmd.RunE(validateCmd, nil); err != nil {
		t.Fatalf("validate: %v\noutput:\n%s", err, out.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("0 error(s)")) {
		t.Fatalf("output = %q", out.String())
	}
}

func TestValidateCmd_missingConfig(t *testing.T) {
	root := t.TempDir()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	var out bytes.Buffer
	validateCmd.SetOut(&out)
	validateCmd.SetErr(&out)
	if err := validateCmd.RunE(validateCmd, nil); err == nil {
		t.Fatalf("expected error, output:\n%s", out.String())
	}
	if !bytes.Contains(out.Bytes(), []byte("config.missing")) {
		t.Fatalf("output = %q", out.String())
	}
}

func writeValidateFixture(t *testing.T, root string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(filepath.Join(base, "templates", "expanded"), 0o755); err != nil {
		t.Fatal(err)
	}
	config := "specs_dir: specs\n"
	if err := os.WriteFile(filepath.Join(base, "config.yaml"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	charter := "---\nname: C\ndescription: d\nstatus: active\nspec_type: simple\n---\n"
	if err := os.WriteFile(filepath.Join(base, "charter.md"), []byte(charter), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"README.md",
		"flexspec-simple.md",
		filepath.Join("expanded", "flexspec-expanded.md"),
		filepath.Join("expanded", "flexspec-expanded-task.md"),
	} {
		path := filepath.Join(base, "templates", name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("# t\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
}
