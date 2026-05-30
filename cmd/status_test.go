package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatusSet(t *testing.T) {
	root := t.TempDir()
	writeListProject(t, root)

	statusSetStatus = "in_progress"
	statusSetTask = ""
	defer func() {
		statusSetStatus = ""
		statusSetTask = ""
	}()

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := runStatusSet(statusSetCmd, []string{"001-test"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "specs", "001-test", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "status: in_progress") {
		t.Fatalf("content = %s", data)
	}
}
