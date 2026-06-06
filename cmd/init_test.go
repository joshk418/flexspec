package cmd

import (
	"bytes"
	"embed"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestInit_createsGlossary(t *testing.T) {
	root := t.TempDir()
	forceInit = true
	TemplatesFS = fstest.MapFS{
		"templates/charter.md":                         {Data: []byte("# charter\n")},
		"templates/glossary.yaml":                      {Data: []byte("version: \"1.0\"\nterms: []\n")},
		"templates/README.md":                          {Data: []byte("# readme\n")},
		"templates/flexspec-simple.md":                 {Data: []byte("# simple\n")},
		"templates/expanded/flexspec-expanded.md":      {Data: []byte("# expanded\n")},
		"templates/expanded/flexspec-expanded-task.md": {Data: []byte("# task\n")},
	}
	defer func() { TemplatesFS = embed.FS{} }()

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	var buf bytes.Buffer
	initCmd.SetOut(&buf)

	if err := initCmd.RunE(initCmd, nil); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(root, ".flexspec", "glossary.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected glossary.yaml to be created")
	}
}

func TestInit_preservesExistingGlossary(t *testing.T) {
	root := t.TempDir()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	existing := "version: \"1.0\"\nterms:\n  - term: KeepMe\n    definition: preserved\n"
	if err := os.WriteFile(filepath.Join(base, "glossary.yaml"), []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}
	forceInit = false

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	var buf bytes.Buffer
	initCmd.SetOut(&buf)

	if err := initCmd.RunE(initCmd, nil); err == nil {
		t.Fatal("expected error for existing .flexspec without --force")
	}

	data, err := os.ReadFile(filepath.Join(base, "glossary.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != existing {
		t.Fatalf("glossary was overwritten:\n%s", string(data))
	}
}
