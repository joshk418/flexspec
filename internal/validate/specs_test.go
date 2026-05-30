package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckSpecs_badFrontmatter(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	specDir := filepath.Join(root, "specs", "001-bad")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte("no frontmatter\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings := CheckSpecs(root, testConfig(t, root), Options{})
	var got *Finding
	for i := range findings {
		if findings[i].Rule == "specs.frontmatter" {
			got = &findings[i]
			break
		}
	}
	if got == nil {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckSpecs_badTaskFrontmatter(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	specDir := filepath.Join(root, "specs", "002-expanded")
	tasksDir := filepath.Join(specDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := `---
name: Expanded
description: test
status: planned
spec_type: expanded
---
`
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tasksDir, "T-001-bad.md"), []byte("broken\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	findings := CheckSpecs(root, testConfig(t, root), Options{})
	var got *Finding
	for i := range findings {
		if findings[i].Rule == "specs.task_frontmatter" {
			got = &findings[i]
			break
		}
	}
	if got == nil {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestCheckSpecs_orphanAndDuplicate(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	if err := os.MkdirAll(filepath.Join(root, "specs", "junk"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"001-a", "001-b"} {
		dir := filepath.Join(root, "specs", name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(minimalSpecREADME), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	findings := CheckSpecs(root, testConfig(t, root), Options{})
	var orphan, dup bool
	for _, f := range findings {
		if f.Rule == "specs.orphan_dir" {
			orphan = true
		}
		if f.Rule == "specs.duplicate_sequence" {
			dup = true
		}
	}
	if !orphan || !dup {
		t.Fatalf("findings = %+v", findings)
	}
}

func TestRunAll_missingConfig(t *testing.T) {
	root := t.TempDir()
	findings := RunAll(root, Options{})
	if !HasErrors(findings) {
		t.Fatal("expected errors")
	}
}

func TestRunAll_ok(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	findings := RunAll(root, Options{})
	if HasErrors(findings) {
		t.Fatalf("findings = %+v", findings)
	}
}
