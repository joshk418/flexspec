package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func TestStatusRename_detectAndApply(t *testing.T) {
	root := t.TempDir()
	specDir := filepath.Join(root, "specs", "001-test")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\nname: Test\ndescription: d\nstatus: refined\nspec_type: simple\n---\n\n# Body\n"
	readme := filepath.Join(specDir, "README.md")
	if err := os.WriteFile(readme, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &statusRenameMigration{}
	cfg := config.Config{SpecsDir: "specs"}
	changes, err := m.Detect(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes, want 1", len(changes))
	}
	applied, err := m.Apply(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(applied) != 1 {
		t.Fatalf("applied %d changes", len(applied))
	}
	data, err := os.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "status: planned") {
		t.Fatalf("expected planned status in %q", string(data))
	}
	if !strings.Contains(string(data), "# Body") {
		t.Fatal("body should be preserved")
	}
}
