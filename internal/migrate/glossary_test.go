package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func writeMinimalProject(t *testing.T, root string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(filepath.Join(base, "templates", "expanded"), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := "specs_dir: specs\nalways_one_shot: false\nspec_template:\n"
	if err := os.WriteFile(filepath.Join(base, "config.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	charter := "---\nname: Charter\ndescription: test\nstatus: active\nspec_type: simple\n---\n"
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
		if err := os.WriteFile(path, []byte("# template\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestGlossaryMigration_detectMissing(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)

	m := &glossaryMigration{}
	changes, err := m.Detect(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Kind != KindCreate {
		t.Fatalf("expected KindCreate, got %s", changes[0].Kind)
	}
}

func TestGlossaryMigration_applyCreatesFile(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)

	m := &glossaryMigration{}
	changes, err := m.Apply(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}

	path := filepath.Join(root, ".flexspec", "glossary.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected glossary.yaml to be created")
	}
}

func TestGlossaryMigration_existingFile(t *testing.T) {
	root := t.TempDir()
	writeMinimalProject(t, root)
	writeGlossaryYAML(t, root, "version: \"1.0\"\nterms: []\n")

	m := &glossaryMigration{}
	changes, err := m.Detect(root, config.Config{SpecsDir: "specs"})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes for existing file, got %d", len(changes))
	}
}

func writeGlossaryYAML(t *testing.T, root, content string) {
	t.Helper()
	path := filepath.Join(root, ".flexspec", "glossary.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
