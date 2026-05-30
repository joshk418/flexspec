package validate

import (
	"os"
	"path/filepath"
	"testing"
)

const minimalCharter = `---
name: Charter
description: test
status: active
spec_type: simple
---
`

const minimalSpecREADME = `---
name: Test
description: test
status: planned
spec_type: simple
---
`

func writeMinimalProject(t *testing.T, root string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(filepath.Join(base, "templates", "expanded"), 0o755); err != nil {
		t.Fatal(err)
	}
	config := "specs_dir: specs\nalways_one_shot: false\nspec_template:\n"
	if err := os.WriteFile(filepath.Join(base, "config.yaml"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "charter.md"), []byte(minimalCharter), 0o644); err != nil {
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
		if err := os.WriteFile(path, []byte("# template\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
}
