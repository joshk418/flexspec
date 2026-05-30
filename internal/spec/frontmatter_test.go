package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetFileStatus_preservesBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	content := "---\nname: Test\nstatus: initial\n---\n\n# Hello\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := SetFileStatus(path, "planned"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "status: planned") {
		t.Fatalf("expected planned status in %q", s)
	}
	if !strings.Contains(s, "# Hello") {
		t.Fatalf("expected body preserved in %q", s)
	}
	if strings.Contains(s, "status: initial") {
		t.Fatalf("old status should be gone: %q", s)
	}
}
