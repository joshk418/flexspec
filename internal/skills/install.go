// Package skills installs embedded FlexSpec agent skills; npx is the fallback when no agent is detected.
package skills

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Scope selects whether skills are written to per-project or per-user dirs.
type Scope int

const (
	// ScopeGlobal writes to each agent's home-directory skills dir (e.g. ~/.claude/skills/). Default.
	ScopeGlobal Scope = iota
	// ScopeProject writes to each agent's project-local skills dir (e.g. ./.claude/skills/).
	ScopeProject
)

// Agent describes one supported coding agent and where its skills live.
type Agent struct {
	Name       string // human-readable identifier (matches vercel-labs/skills)
	ProjectDir string // project-local skills directory, relative to project root
	GlobalDir  string // user-wide skills directory, relative to home
	RootDir    string // agent config root used for detection (we check this dir exists, not just skills/)
}

// Supported agents. Paths sourced from the vercel-labs/skills README.
var agents = []Agent{
	{Name: "claude-code", ProjectDir: ".claude/skills", GlobalDir: ".claude/skills", RootDir: ".claude"},
	{Name: "cursor", ProjectDir: ".agents/skills", GlobalDir: ".cursor/skills", RootDir: ".cursor"},
	{Name: "codex", ProjectDir: ".agents/skills", GlobalDir: ".codex/skills", RootDir: ".codex"},
	{Name: "opencode", ProjectDir: ".agents/skills", GlobalDir: ".config/opencode/skills", RootDir: ".config/opencode"},
	{Name: "cline", ProjectDir: ".agents/skills", GlobalDir: ".agents/skills", RootDir: ".agents"},
}

// DetectAgents returns supported agents whose config root exists; skills/ is created on install if missing.
func DetectAgents(scope Scope, home, projectRoot string) []Agent {
	var out []Agent
	for _, a := range agents {
		var root string
		if scope == ScopeProject {
			root = filepath.Join(projectRoot, filepath.FromSlash(a.RootDir))
		} else {
			root = filepath.Join(home, filepath.FromSlash(a.RootDir))
		}
		if isDir(root) {
			out = append(out, a)
		}
	}
	return out
}

// Result records what was written for one agent.
type Result struct {
	Agent        string
	Dir          string
	FilesWritten int
}

// Install writes the embedded skills tree into each agent's skills directory, overwriting version-pinned files.
func Install(skillsFS fs.FS, scope Scope, home, projectRoot string, agents []Agent) ([]Result, error) {
	var results []Result
	for _, a := range agents {
		dir := skillsDirFor(a, scope, home, projectRoot)
		count, err := installInto(skillsFS, dir)
		if err != nil {
			return results, fmt.Errorf("install skills for %s: %w", a.Name, err)
		}
		results = append(results, Result{
			Agent:        a.Name,
			Dir:          dir,
			FilesWritten: count,
		})
	}
	return results, nil
}

// skillsDirFor resolves the target skills directory for an agent given the scope.
func skillsDirFor(a Agent, scope Scope, home, projectRoot string) string {
	if scope == ScopeProject {
		return filepath.Join(projectRoot, filepath.FromSlash(a.ProjectDir))
	}
	return filepath.Join(home, filepath.FromSlash(a.GlobalDir))
}

// installInto writes every file from skillsFS into dest, preserving subdirectories and overwriting existing files.
func installInto(skillsFS fs.FS, dest string) (int, error) {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return 0, fmt.Errorf("mkdir %s: %w", dest, err)
	}
	count := 0
	err := fs.WalkDir(skillsFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "." {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		rel := filepath.ToSlash(p)
		target := filepath.Join(dest, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		data, err := fs.ReadFile(skillsFS, p)
		if err != nil {
			return err
		}
		if err := os.WriteFile(target, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", target, err)
		}
		count++
		return nil
	})
	return count, err
}

// isDir reports whether path exists and is a directory.
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FallbackInstruction tells the user how to install skills manually when no agent/npx fallback is available.
const FallbackInstruction = `No supported coding agent detected and npx not found on PATH.

To install FlexSpec skills manually, install Node.js and run:
    npx skills add joshk418/flexspec --global

Supported agents: claude-code, cursor, codex, opencode, cline.
See https://github.com/vercel-labs/skills for the full list.`

// PrintFallbackInstruction writes manual-install instructions to w, ignoring write errors.
func PrintFallbackInstruction(w io.Writer) {
	_, _ = fmt.Fprintln(w, FallbackInstruction)
}

// AgentNames returns supported agent names for display in --dry-run output.
func AgentNames() string {
	names := make([]string, 0, len(agents))
	for _, a := range agents {
		names = append(names, a.Name)
	}
	return strings.Join(names, ", ")
}
