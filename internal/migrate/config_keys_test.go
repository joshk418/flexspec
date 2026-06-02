package migrate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func TestConfigKeys_addSpecTemplate(t *testing.T) {
	root := t.TempDir()
	flexDir := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(flexDir, 0o755); err != nil {
		t.Fatal(err)
	}
	cfgContent := "specs_dir: specs\nalways_one_shot: false\n"
	if err := os.WriteFile(filepath.Join(flexDir, "config.yaml"), []byte(cfgContent), 0o644); err != nil {
		t.Fatal(err)
	}

	m := &configKeysMigration{}
	cfg := config.Config{SpecsDir: "specs", AlwaysOneShot: false}
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
	data, err := os.ReadFile(filepath.Join(flexDir, "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "spec_template") {
		t.Fatalf("expected spec_template in %q", string(data))
	}
}
