package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestListHuman(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, root string)
		wantSubstr []string
		wantAbsent []string
	}{
		{
			name: "empty specs dir",
			setup: func(t *testing.T, root string) {
				t.Helper()
				writeListConfig(t, root)
			},
			wantSubstr: []string{"No specs in specs\n"},
		},
		{
			name: "simple spec task_count frontmatter",
			setup: func(t *testing.T, root string) {
				t.Helper()
				writeListConfig(t, root)
				writeSimpleSpecWithTaskCount(t, root, "001-test", "planned", 3)
			},
			wantSubstr: []string{
				"IDENTIFIER",
				"STATUS",
				"TASKS",
				"001-test",
				"planned",
				"3",
			},
		},
		{
			name: "simple spec computed fallback",
			setup: func(t *testing.T, root string) {
				t.Helper()
				writeListConfig(t, root)
				writeSimpleSpecWithBullets(t, root, "001-test", "planned", 2)
			},
			wantSubstr: []string{
				"001-test",
				"2",
			},
		},
		{
			name: "expanded spec with tasks",
			setup: func(t *testing.T, root string) {
				t.Helper()
				writeListConfig(t, root)
				writeExpandedSpec(t, root, "002-feature", "in_progress", 3)
			},
			wantSubstr: []string{
				"002-feature",
				"in_progress",
				"3",
			},
			wantAbsent: []string{"T-001", "T-002", "T-003"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			tt.setup(t, root)

			listJSON = false
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

			out := buf.String()
			for _, want := range tt.wantSubstr {
				if !strings.Contains(out, want) {
					t.Errorf("output missing %q\n%s", want, out)
				}
			}
			for _, absent := range tt.wantAbsent {
				if strings.Contains(out, absent) {
					t.Errorf("output should not contain %q\n%s", absent, out)
				}
			}
		})
	}
}

func writeListConfig(t *testing.T, root string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "config.yaml"), []byte("specs_dir: specs\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeSimpleSpecWithTaskCount(t *testing.T, root, dirName, status string, taskCount int) {
	t.Helper()
	specDir := filepath.Join(root, "specs", dirName)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := "---\nname: Test\ndescription: d\nstatus: " + status + "\nspec_type: simple\ntask_count: " + strconv.Itoa(taskCount) + "\n---\n"
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeSimpleSpecWithBullets(t *testing.T, root, dirName, status string, bullets int) {
	t.Helper()
	specDir := filepath.Join(root, "specs", dirName)
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	b.WriteString("---\nname: Test\ndescription: d\nstatus: ")
	b.WriteString(status)
	b.WriteString("\nspec_type: simple\n---\n\n")
	for i := 1; i <= bullets; i++ {
		b.WriteString("- **T-00")
		b.WriteString(strconv.Itoa(i))
		b.WriteString("** — task\n")
	}
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func writeExpandedSpec(t *testing.T, root, dirName, status string, taskCount int) {
	t.Helper()
	specDir := filepath.Join(root, "specs", dirName)
	tasksDir := filepath.Join(specDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := "---\nname: Feature\ndescription: d\nstatus: " + status + "\nspec_type: expanded\n---\n"
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= taskCount; i++ {
		id := "T-00" + string(rune('0'+i))
		file := id + "-task.md"
		fm := "---\nid: " + id + "\nname: Task\nstatus: todo\n---\n"
		if err := os.WriteFile(filepath.Join(tasksDir, file), []byte(fm), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}
