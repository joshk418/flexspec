package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const brainstormTestTemplate = "---\nname: '{name}'\n---\n# {name}\n"

func writeBrainstormTemplate(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, ".flexspec", "templates")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "brainstorm.md"), []byte(brainstormTestTemplate), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBrainstormNew(t *testing.T) {
	root := t.TempDir()
	writeBrainstormTemplate(t, root)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	targetPath := filepath.Join(root, ".flexspec", "brainstorms", "demo-feature.md")

	t.Run("creates doc and prints path", func(t *testing.T) {
		brainstormForce = false
		var buf bytes.Buffer
		brainstormNewCmd.SetOut(&buf)

		if err := runBrainstormNew(brainstormNewCmd, []string{"demo feature"}); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(buf.String(), targetPath) {
			t.Errorf("output = %q, want it to contain %q", buf.String(), targetPath)
		}
		if _, err := os.Stat(targetPath); err != nil {
			t.Fatalf("expected file at %s: %v", targetPath, err)
		}
	})

	t.Run("second run without force errors", func(t *testing.T) {
		brainstormForce = false
		var buf bytes.Buffer
		brainstormNewCmd.SetOut(&buf)

		err := runBrainstormNew(brainstormNewCmd, []string{"demo feature"})
		if err == nil {
			t.Fatal("expected error for existing doc without --force")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Errorf("error = %q, want it to mention already exists", err.Error())
		}
	})

	t.Run("second run with force overwrites", func(t *testing.T) {
		brainstormForce = true
		defer func() { brainstormForce = false }()
		var buf bytes.Buffer
		brainstormNewCmd.SetOut(&buf)

		if err := runBrainstormNew(brainstormNewCmd, []string{"demo feature"}); err != nil {
			t.Fatal(err)
		}
	})
}

func TestBrainstormNew_requiresName(t *testing.T) {
	if err := brainstormNewCmd.Args(brainstormNewCmd, []string{}); err == nil {
		t.Fatal("expected an error when no name is given")
	}
}
