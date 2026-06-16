package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountTasks_simpleBullets(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	body := `---
name: Test
description: d
status: planned
spec_type: simple
---

# Test

- **T-001** — First
- **T-002** — Second
- **T-003** — Third
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	meta, err := ParseSpecMeta(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := CountTasks(path, meta)
	if err != nil {
		t.Fatal(err)
	}
	if got != 3 {
		t.Errorf("CountTasks = %d, want 3", got)
	}
}

func TestCountTasks_simpleTaskTableRows(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	body := `---
name: Test
description: d
status: planned
spec_type: simple
---

# Test

| Task | Name | Description | Blocks | Blocked by | Requirements |
| --- | --- | --- | --- | --- | --- |
| **T-001** | First | Do one | T-002 | - | FR-001 |
| T-002 | Second | Do two | - | T-001 | FR-002 |
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	meta, err := ParseSpecMeta(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := CountTasks(path, meta)
	if err != nil {
		t.Fatal(err)
	}
	if got != 2 {
		t.Errorf("CountTasks = %d, want 2", got)
	}
}

func TestCountTasks_expandedFiles(t *testing.T) {
	root := t.TempDir()
	specDir := filepath.Join(root, "spec")
	tasksDir := filepath.Join(specDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := filepath.Join(specDir, "README.md")
	if err := os.WriteFile(readme, []byte(expandedFrontmatter), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"T-001-a.md", "T-002-b.md"} {
		fm := "---\nid: T-001\nname: x\nstatus: todo\n---\n"
		if err := os.WriteFile(filepath.Join(tasksDir, name), []byte(fm), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	meta, err := ParseSpecMeta(readme)
	if err != nil {
		t.Fatal(err)
	}
	got, err := CountTasks(readme, meta)
	if err != nil {
		t.Fatal(err)
	}
	if got != 2 {
		t.Errorf("CountTasks = %d, want 2", got)
	}
}

func TestEffectiveTaskCount_prefersFrontmatter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	count := 9
	body := `---
name: Test
description: d
status: planned
spec_type: simple
task_count: 9
---

# Test
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	meta, err := ParseSpecMeta(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := EffectiveTaskCount(path, meta)
	if err != nil {
		t.Fatal(err)
	}
	if got != count {
		t.Errorf("EffectiveTaskCount = %d, want %d", got, count)
	}
}

func TestEffectiveTaskCount_computedWhenMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	body := `---
name: Test
description: d
status: planned
spec_type: simple
---

# Test

- **T-001** — One
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	meta, err := ParseSpecMeta(path)
	if err != nil {
		t.Fatal(err)
	}
	got, err := EffectiveTaskCount(path, meta)
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Errorf("EffectiveTaskCount = %d, want 1", got)
	}
}

func TestSyncTaskCount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	body := `---
name: Test
description: d
status: planned
spec_type: simple
---

# Test

> **Status**: planned · **Priority**: low · **Created**: 2026-01-01

- **T-001** — One
- **T-002** — Two
`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SyncTaskCount(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "task_count: 2") {
		t.Fatalf("missing task_count: 2 in %q", s)
	}
	if !strings.Contains(s, "**Tasks**: 2") {
		t.Fatalf("missing **Tasks**: 2 in %q", s)
	}
}
