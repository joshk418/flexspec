package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func TestTypeBackfill_missingType(t *testing.T) {
	root := t.TempDir()
	specsDir := filepath.Join(root, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Spec with no type field -> should backfill feature (default).
	readme := `---
name: Test
description: d
status: planned
spec_type: simple
---

# Test
`
	specDir := filepath.Join(specsDir, "001-test")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &typeBackfillMigration{}
	changes, err := m.Apply(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}

	data, err := os.ReadFile(filepath.Join(specDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "type: feature") {
		t.Errorf("backfilled README should contain type: feature, got %q", string(data))
	}
}

func TestTypeBackfill_bugFromSections(t *testing.T) {
	root := t.TempDir()
	specsDir := filepath.Join(root, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Spec with filled Sections 4/5 -> should infer bug.
	readme := `---
name: Login bug
description: d
status: planned
spec_type: simple
---

# Login bug

## 4. Expected Result (bugs only)

User sees a success toast.

## 5. Actual Result (bugs only)

User sees a blank screen.
`
	specDir := filepath.Join(specsDir, "002-login-bug")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &typeBackfillMigration{}
	changes, err := m.Apply(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}
	if !strings.Contains(changes[0].Detail, "bug") {
		t.Errorf("change detail should mention bug, got %q", changes[0].Detail)
	}

	data, err := os.ReadFile(filepath.Join(specDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "type: bug") {
		t.Errorf("backfilled README should contain type: bug, got %q", string(data))
	}
}

func TestTypeBackfill_skipsWhenTypePresent(t *testing.T) {
	root := t.TempDir()
	specsDir := filepath.Join(root, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	readme := `---
name: Chore
description: d
status: planned
spec_type: simple
type: chore
---

# Chore
`
	specDir := filepath.Join(specsDir, "003-chore")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &typeBackfillMigration{}
	changes, err := m.Detect(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 0 {
		t.Fatalf("got %d changes, want 0 (type already set)", len(changes))
	}
}

func TestTypeBackfill_skipsNotApplicable(t *testing.T) {
	root := t.TempDir()
	specsDir := filepath.Join(root, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Section 4/5 present with bug headings but says "Not applicable" -> feature, not bug.
	readme := `---
name: Feature
description: d
status: planned
spec_type: simple
---

# Feature

## 4. Expected Result (bugs only)

Not applicable - this is not a bug fix.

## 5. Actual Result (bugs only)

Not applicable - this is not a bug fix.
`
	specDir := filepath.Join(specsDir, "004-feature")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &typeBackfillMigration{}
	changes, err := m.Apply(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}

	data, err := os.ReadFile(filepath.Join(specDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "type: feature") {
		t.Errorf("should infer feature for not-applicable sections, got %q", string(data))
	}
}

func TestTypeBackfill_nonBugSection45DefaultsFeature(t *testing.T) {
	root := t.TempDir()
	specsDir := filepath.Join(root, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Older spec format: Section 4/5 used for non-bug content -> feature, not bug.
	readme := `---
name: Old feature
description: d
status: planned
spec_type: simple
---

# Old feature

## 4. Testing Criteria

Must pass go test.

## 5. Other

Notes here.
`
	specDir := filepath.Join(specsDir, "005-old-feature")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &typeBackfillMigration{}
	changes, err := m.Apply(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}

	data, err := os.ReadFile(filepath.Join(specDir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "type: feature") {
		t.Errorf("non-bug Section 4/5 should default to feature, got %q", string(data))
	}
}

func TestTypeBackfill_emptySpecsDir(t *testing.T) {
	root := t.TempDir()
	m := &typeBackfillMigration{}
	changes, err := m.Detect(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 0 {
		t.Fatalf("got %d changes, want 0", len(changes))
	}
}
