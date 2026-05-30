package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/joshk418/flexspec/internal/ui"
)

func TestListJSON(t *testing.T) {
	root := t.TempDir()
	writeListProject(t, root)

	listJSON = true
	defer func() { listJSON = false }()

	var buf bytes.Buffer
	listCmd.SetOut(&buf)
	listCmd.SetErr(&buf)
	listCmd.SetArgs(nil)
	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := listCmd.RunE(listCmd, nil); err != nil {
		t.Fatal(err)
	}
	var specs []ui.SpecJSON
	if err := json.Unmarshal(buf.Bytes(), &specs); err != nil {
		t.Fatal(err)
	}
	if len(specs) != 1 || specs[0].Name != "Test" {
		t.Fatalf("specs = %+v", specs)
	}
}

func writeListProject(t *testing.T, root string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "config.yaml"), []byte("specs_dir: specs\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	specDir := filepath.Join(root, "specs", "001-test")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := "---\nname: Test\ndescription: d\nstatus: planned\nspec_type: simple\n---\n"
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
}
