package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

const simpleFrontmatter = `---
name: User auth
description: Login and session flow
status: in_progress
spec_type: simple
---

# User auth
`

const expandedFrontmatter = `---
name: Billing export
description: Multi-step export
status: planned
spec_type: expanded
---
`

const taskFrontmatter = `---
id: T-001
name: Create schema
status: todo
---
`

func TestParseSpecMeta_simple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte(simpleFrontmatter), 0o644); err != nil {
		t.Fatal(err)
	}

	meta, err := ParseSpecMeta(path)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "User auth" {
		t.Errorf("name = %q, want User auth", meta.Name)
	}
	if meta.Description != "Login and session flow" {
		t.Errorf("description = %q", meta.Description)
	}
	if meta.Status != "in_progress" {
		t.Errorf("status = %q", meta.Status)
	}
	if meta.SpecType != "simple" {
		t.Errorf("spec_type = %q", meta.SpecType)
	}
}

func TestParseSpecMeta_expanded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte(expandedFrontmatter), 0o644); err != nil {
		t.Fatal(err)
	}

	meta, err := ParseSpecMeta(path)
	if err != nil {
		t.Fatal(err)
	}
	if meta.SpecType != "expanded" {
		t.Errorf("spec_type = %q", meta.SpecType)
	}
}

func TestParseTaskMeta(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-001-create-schema.md")
	if err := os.WriteFile(path, []byte(taskFrontmatter), 0o644); err != nil {
		t.Fatal(err)
	}

	meta, err := ParseTaskMeta(path)
	if err != nil {
		t.Fatal(err)
	}
	if meta.ID != "T-001" {
		t.Errorf("id = %q", meta.ID)
	}
	if meta.Name != "Create schema" {
		t.Errorf("name = %q", meta.Name)
	}
	if meta.Status != "todo" {
		t.Errorf("status = %q", meta.Status)
	}
}

func TestSplitFrontmatter_missingClose(t *testing.T) {
	_, err := splitFrontmatter("---\nname: x\n")
	if err == nil {
		t.Fatal("expected error for missing closing ---")
	}
}

func TestList_sortOrder(t *testing.T) {
	root := t.TempDir()
	specsDir := filepath.Join(root, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	for _, name := range []string{"010-bar", "001-foo", "002-baz"} {
		dir := filepath.Join(specsDir, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		fm := `---
name: ` + name + `
description: desc
status: initial
spec_type: simple
---
`
		if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(fm), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := List(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Fatalf("got %d entries, want 3", len(entries))
	}
	want := []string{"001", "002", "010"}
	for i, id := range want {
		if entries[i].ID != id {
			t.Errorf("entries[%d].ID = %q, want %q", i, entries[i].ID, id)
		}
	}
}

func TestList_expandedTasks(t *testing.T) {
	root := t.TempDir()
	specDir := filepath.Join(root, "specs", "001-feature")
	tasksDir := filepath.Join(specDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0o755); err != nil {
		t.Fatal(err)
	}

	specFM := `---
name: Feature
description: Big feature
status: planned
spec_type: expanded
---
`
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(specFM), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, file := range []string{"T-002-second.md", "T-001-first.md"} {
		fm := `---
id: ` + file[:5] + `
name: ` + file + `
status: todo
---
`
		if err := os.WriteFile(filepath.Join(tasksDir, file), []byte(fm), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := List(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d entries", len(entries))
	}
	if len(entries[0].Tasks) != 2 {
		t.Fatalf("got %d tasks", len(entries[0].Tasks))
	}
	if entries[0].Tasks[0].Meta.ID != "T-001" {
		t.Errorf("first task id = %q", entries[0].Tasks[0].Meta.ID)
	}
	if entries[0].Tasks[1].Meta.ID != "T-002" {
		t.Errorf("second task id = %q", entries[0].Tasks[1].Meta.ID)
	}
}

func TestList_missingSpecsDir(t *testing.T) {
	root := t.TempDir()
	entries, err := List(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("got %d entries, want 0", len(entries))
	}
}
