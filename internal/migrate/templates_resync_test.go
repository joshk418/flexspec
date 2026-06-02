package migrate

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/joshk418/flexspec/internal/config"
)

func TestTemplatesResync_missingAndDiffers(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, ".flexspec", "templates")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatal(err)
	}
	existing := filepath.Join(dest, "flexspec-simple.md")
	if err := os.WriteFile(existing, []byte("user-edited"), 0o644); err != nil {
		t.Fatal(err)
	}

	tmpl := fstest.MapFS{
		"flexspec-simple.md":   {Data: []byte("embedded")},
		"README.md":            {Data: []byte("# templates")},
		"expanded/expanded.md": {Data: []byte("expanded")},
	}
	m := &templatesResyncMigration{templates: tmpl, force: false}
	cfg := config.Config{SpecsDir: "specs"}

	changes, err := m.Detect(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) < 2 {
		t.Fatalf("got %d changes, want at least 2 (differs + missing)", len(changes))
	}
	applied, err := m.Apply(root, cfg)
	if err != nil {
		t.Fatal(err)
	}
	// Missing file created; differing file not overwritten without force.
	data, _ := os.ReadFile(existing)
	if string(data) != "user-edited" {
		t.Fatalf("expected user-edited, got %q", string(data))
	}
	created := filepath.Join(dest, "expanded", "expanded.md")
	if _, err := os.Stat(created); err != nil {
		t.Fatalf("expected created file: %v", err)
	}
	if len(applied) < 2 {
		t.Fatalf("applied %d, want at least 2", len(applied))
	}

	m.force = true
	_, _ = m.Apply(root, cfg)
	data, _ = os.ReadFile(existing)
	if string(data) != "embedded" {
		t.Fatalf("expected embedded after force, got %q", string(data))
	}
}

var _ fs.FS = fstest.MapFS{}
