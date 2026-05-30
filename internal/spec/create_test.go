package spec

import (
	"path/filepath"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

const testRoot = "project"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "spaces", input: "User Auth", want: "user-auth"},
		{name: "underscores", input: "USER_AUTH", want: "user-auth"},
		{name: "collapse hyphens", input: "  foo---bar  ", want: "foo-bar"},
		{name: "invalid", input: "!!!", wantErr: true},
		{name: "empty", input: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Slugify(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNextSequence(t *testing.T) {
	t.Run("empty dir", func(t *testing.T) {
		fsys := newMemFS()
		specsDir := filepath.Join(testRoot, "specs")
		if err := fsys.MkdirAll(specsDir, dirPerm); err != nil {
			t.Fatal(err)
		}
		seq, err := nextSequenceWithFS(fsys, specsDir)
		if err != nil {
			t.Fatal(err)
		}
		if seq != 1 {
			t.Errorf("seq = %d, want 1", seq)
		}
	})

	t.Run("missing dir", func(t *testing.T) {
		fsys := newMemFS()
		seq, err := nextSequenceWithFS(fsys, filepath.Join(testRoot, "specs"))
		if err != nil {
			t.Fatal(err)
		}
		if seq != 1 {
			t.Errorf("seq = %d, want 1", seq)
		}
	})

	t.Run("existing numbered dirs", func(t *testing.T) {
		fsys := newMemFS()
		dir := filepath.Join(testRoot, "specs")
		for _, name := range []string{"001-a", "010-b", "not-a-spec"} {
			if err := fsys.MkdirAll(filepath.Join(dir, name), dirPerm); err != nil {
				t.Fatal(err)
			}
		}
		seq, err := nextSequenceWithFS(fsys, dir)
		if err != nil {
			t.Fatal(err)
		}
		if seq != 11 {
			t.Errorf("seq = %d, want 11", seq)
		}
	})
}

func TestCreate_simple(t *testing.T) {
	fsys := newMemFS()
	setupProjectMem(fsys, testRoot)
	cfg := config.Config{SpecsDir: "specs"}

	result, err := createWithFS(fsys, testRoot, cfg, "user-auth", "simple")
	if err != nil {
		t.Fatal(err)
	}
	if result.DirName != "001-user-auth" {
		t.Errorf("DirName = %q, want 001-user-auth", result.DirName)
	}

	readme := filepath.Join(result.SpecPath, "README.md")
	data, err := fsys.ReadFile(readme)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "---\nspec_type: simple\n---\n# simple\n" {
		t.Errorf("README content = %q", string(data))
	}
}

func TestCreate_expanded(t *testing.T) {
	fsys := newMemFS()
	setupProjectMem(fsys, testRoot)
	cfg := config.Config{SpecsDir: "specs"}

	result, err := createWithFS(fsys, testRoot, cfg, "billing-export", "expanded")
	if err != nil {
		t.Fatal(err)
	}
	if result.DirName != "001-billing-export" {
		t.Errorf("DirName = %q, want 001-billing-export", result.DirName)
	}

	readme := filepath.Join(result.SpecPath, "README.md")
	if _, err := fsys.Stat(readme); err != nil {
		t.Fatalf("README.md: %v", err)
	}
	tasksDir := filepath.Join(result.SpecPath, "tasks")
	info, err := fsys.Stat(tasksDir)
	if err != nil {
		t.Fatalf("tasks dir: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("tasks is not a directory")
	}
}

func TestCreate_increments(t *testing.T) {
	fsys := newMemFS()
	setupProjectMem(fsys, testRoot)
	cfg := config.Config{SpecsDir: "specs"}

	first, err := createWithFS(fsys, testRoot, cfg, "first", "simple")
	if err != nil {
		t.Fatal(err)
	}
	if first.DirName != "001-first" {
		t.Errorf("first DirName = %q", first.DirName)
	}

	second, err := createWithFS(fsys, testRoot, cfg, "second", "simple")
	if err != nil {
		t.Fatal(err)
	}
	if second.DirName != "002-second" {
		t.Errorf("second DirName = %q, want 002-second", second.DirName)
	}
}

func TestCreate_existingDir(t *testing.T) {
	fsys := newMemFS()
	setupProjectMem(fsys, testRoot)
	cfg := config.Config{SpecsDir: "specs"}
	specsPath := filepath.Join(testRoot, "specs")
	for _, name := range []string{"001-dummy", "002-other"} {
		if err := fsys.MkdirAll(filepath.Join(specsPath, name), dirPerm); err != nil {
			t.Fatal(err)
		}
	}
	if err := fsys.WriteFile(filepath.Join(specsPath, "003-user-auth"), []byte("block"), filePerm); err != nil {
		t.Fatal(err)
	}

	_, err := createWithFS(fsys, testRoot, cfg, "user-auth", "simple")
	if err == nil {
		t.Fatal("expected error when spec directory already exists")
	}
}

func TestCreate_invalidTemplate(t *testing.T) {
	fsys := newMemFS()
	setupProjectMem(fsys, testRoot)
	cfg := config.Config{SpecsDir: "specs"}

	_, err := createWithFS(fsys, testRoot, cfg, "feature", "invalid")
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestCreate_missingTemplate(t *testing.T) {
	fsys := newMemFS()
	cfg := config.Config{SpecsDir: "specs"}

	_, err := createWithFS(fsys, testRoot, cfg, "feature", "simple")
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}
