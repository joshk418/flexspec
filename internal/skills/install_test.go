package skills

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

// fakeSkillsFS builds an in-memory skills FS with two skills and references.
func fakeSkillsFS(t *testing.T) fs.FS {
	t.Helper()
	return fstest.MapFS{
		"flexspec/SKILL.md":                &fstest.MapFile{Data: []byte("# flexspec\n")},
		"flexspec-migrate/SKILL.md":        &fstest.MapFile{Data: []byte("# migrate\n")},
		"flexspec-migrate/references/x.md": &fstest.MapFile{Data: []byte("# x\n")},
		"flexspec-migrate/references/y.md": &fstest.MapFile{Data: []byte("# y\n")},
	}
}

func TestDetectAgents_globalScope(t *testing.T) {
	tmp := t.TempDir()
	// Create ~/.claude and ~/.cursor roots, but NOT ~/.codex.
	for _, d := range []string{".claude", ".cursor"} {
		if err := os.MkdirAll(filepath.Join(tmp, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	got := DetectAgents(ScopeGlobal, tmp, "/irrelevant")
	if len(got) != 2 {
		t.Fatalf("got %d agents: %+v", len(got), got)
	}
	names := []string{got[0].Name, got[1].Name}
	if names[0] != "claude-code" && names[1] != "claude-code" {
		t.Errorf("claude-code missing: %v", names)
	}
	if names[0] != "cursor" && names[1] != "cursor" {
		t.Errorf("cursor missing: %v", names)
	}
}

func TestDetectAgents_projectScope(t *testing.T) {
	project := t.TempDir()
	// Create ./.claude in the project root.
	if err := os.MkdirAll(filepath.Join(project, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	// Home has .codex but project doesn't, so codex shouldn't be detected.
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".codex"), 0o755); err != nil {
		t.Fatal(err)
	}
	got := DetectAgents(ScopeProject, home, project)
	if len(got) != 1 {
		t.Fatalf("got %d agents: %+v", len(got), got)
	}
	if got[0].Name != "claude-code" {
		t.Errorf("got %q, want claude-code", got[0].Name)
	}
}

func TestDetectAgents_noneDetected(t *testing.T) {
	tmp := t.TempDir()
	got := DetectAgents(ScopeGlobal, tmp, ".")
	if len(got) != 0 {
		t.Fatalf("got %d agents, want 0", len(got))
	}
}

func TestInstall_writesAllFilesToAllAgents(t *testing.T) {
	skillsFS := fakeSkillsFS(t)
	home := t.TempDir()

	// Pretend claude-code and opencode are installed.
	for _, d := range []string{".claude", ".config/opencode"} {
		if err := os.MkdirAll(filepath.Join(home, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	detected := DetectAgents(ScopeGlobal, home, ".")
	if len(detected) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(detected))
	}

	results, err := Install(skillsFS, ScopeGlobal, home, ".", detected)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}

	// Each agent should have 4 files.
	for _, r := range results {
		if r.FilesWritten != 4 {
			t.Errorf("%s: FilesWritten = %d, want 4", r.Agent, r.FilesWritten)
		}
		// Verify the files actually landed.
		skillPath := filepath.Join(r.Dir, "flexspec", "SKILL.md")
		data, err := os.ReadFile(skillPath)
		if err != nil {
			t.Errorf("%s: read %s: %v", r.Agent, skillPath, err)
			continue
		}
		if string(data) != "# flexspec\n" {
			t.Errorf("%s: SKILL.md content = %q", r.Agent, string(data))
		}
		refPath := filepath.Join(r.Dir, "flexspec-migrate", "references", "x.md")
		if _, err := os.Stat(refPath); err != nil {
			t.Errorf("%s: references/x.md missing: %v", r.Agent, err)
		}
	}
}

func TestInstall_overwritesExisting(t *testing.T) {
	skillsFS := fakeSkillsFS(t)
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0o755); err != nil {
		t.Fatal(err)
	}

	// Pre-write an old SKILL.md to verify overwrite.
	target := filepath.Join(home, ".claude", "skills", "flexspec", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, []byte("OLD CONTENT"), 0o644); err != nil {
		t.Fatal(err)
	}

	detected := DetectAgents(ScopeGlobal, home, ".")
	_, err := Install(skillsFS, ScopeGlobal, home, ".", detected)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "# flexspec\n" {
		t.Fatalf("file was not overwritten: %q", string(got))
	}
}

func TestInstall_projectScope(t *testing.T) {
	skillsFS := fakeSkillsFS(t)
	project := t.TempDir()
	if err := os.MkdirAll(filepath.Join(project, ".agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	// cline uses .agents root for detection, .agents/skills for project install.
	detected := DetectAgents(ScopeProject, "/irrelevant/home", project)
	if len(detected) == 0 {
		t.Fatal("expected at least one agent (.agents dir should match cline)")
	}
	results, err := Install(skillsFS, ScopeProject, "/irrelevant/home", project, detected)
	if err != nil {
		t.Fatal(err)
	}
	// Find the cline result.
	var clineResult *Result
	for i := range results {
		if results[i].Agent == "cline" {
			clineResult = &results[i]
		}
	}
	if clineResult == nil {
		t.Fatalf("cline not in results: %+v", results)
	}
	expectedDir := filepath.Join(project, ".agents", "skills")
	if clineResult.Dir != expectedDir {
		t.Errorf("dir = %q, want %q", clineResult.Dir, expectedDir)
	}
	if _, err := os.Stat(filepath.Join(expectedDir, "flexspec", "SKILL.md")); err != nil {
		t.Errorf("file not written to project dir: %v", err)
	}
}

func TestInstall_emptyAgentList(t *testing.T) {
	results, err := Install(fakeSkillsFS(t), ScopeGlobal, t.TempDir(), ".", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("got %d results, want 0", len(results))
	}
}

func TestAgentNames(t *testing.T) {
	got := AgentNames()
	for _, want := range []string{"claude-code", "cursor", "codex", "opencode", "cline"} {
		if !strings.Contains(got, want) {
			t.Errorf("AgentNames() = %q, missing %q", got, want)
		}
	}
}
