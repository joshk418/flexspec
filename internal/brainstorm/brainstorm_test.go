package brainstorm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testTemplate = "---\nname: '{name}'\n---\n# Brainstorm: {name}\n"

func writeTemplate(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, flexspecDir, templatesDir)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, brainstormTemplate), []byte(testTemplate), filePerm); err != nil {
		t.Fatal(err)
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		withTemplate bool
		preExisting  string
		slug         string
		force        bool
		wantErr      string
		wantContent  string
	}{
		{
			name:         "scaffolds file and directory from template",
			withTemplate: true,
			slug:         "My Feature",
			wantContent:  testTemplate,
		},
		{
			name:         "missing template returns actionable error",
			withTemplate: false,
			slug:         "my-feature",
			wantErr:      "flexspec init",
		},
		{
			name:         "existing file without force errors and is left untouched",
			withTemplate: true,
			slug:         "existing",
			preExisting:  "old content",
			force:        false,
			wantErr:      "already exists",
			wantContent:  "old content",
		},
		{
			name:         "existing file with force overwrites",
			withTemplate: true,
			slug:         "existing",
			preExisting:  "old content",
			force:        true,
			wantContent:  testTemplate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if tt.withTemplate {
				writeTemplate(t, root)
			}
			if tt.preExisting != "" {
				dir := filepath.Join(root, flexspecDir, brainstormsDir)
				if err := os.MkdirAll(dir, dirPerm); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dir, "existing.md"), []byte(tt.preExisting), filePerm); err != nil {
					t.Fatal(err)
				}
			}

			result, err := Create(root, tt.slug, tt.force)

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var targetPath string
			if err == nil {
				targetPath = result.Path
			} else if tt.preExisting != "" {
				targetPath = filepath.Join(root, flexspecDir, brainstormsDir, "existing.md")
			}

			if tt.wantContent != "" && targetPath != "" {
				got, readErr := os.ReadFile(targetPath)
				if readErr != nil {
					t.Fatalf("read %s: %v", targetPath, readErr)
				}
				if string(got) != tt.wantContent {
					t.Errorf("content = %q, want %q", string(got), tt.wantContent)
				}
			}
		})
	}
}

func TestCreate_missingTemplateMentionsUpdateMigrate(t *testing.T) {
	root := t.TempDir()
	_, err := Create(root, "topic", false)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "flexspec update --migrate") {
		t.Fatalf("error = %q, want it to mention `flexspec update --migrate`", err.Error())
	}
}
