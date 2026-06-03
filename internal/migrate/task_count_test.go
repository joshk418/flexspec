package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func TestTaskCount_detectAndApply_simple(t *testing.T) {
	root := t.TempDir()
	specDir := filepath.Join(root, "specs", "001-test")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `---
name: Test
description: d
status: planned
spec_type: simple
---

# Test

> **Status**: planned · **Priority**: medium · **Created**: 2026-01-01

### 3.2 Task List

- **T-001** — First
- **T-002** — Second
`
	readme := filepath.Join(specDir, "README.md")
	if err := os.WriteFile(readme, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &taskCountMigration{}
	cfg := config.Config{SpecsDir: "specs"}
	changes, err := m.Detect(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}
	if _, err := m.Apply(root, cfg); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "task_count: 2") {
		t.Fatalf("expected task_count: 2 in %q", s)
	}
	if !strings.Contains(s, "**Tasks**: 2") {
		t.Fatalf("expected metadata Tasks: 2 in %q", s)
	}
}
