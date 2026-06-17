package selfupdate

import (
	"fmt"
	"os"
	"os/exec"
)

// skillsPackage is the npm package name for the npx fallback (no supported agent detected).
const skillsPackage = "joshk418/flexspec"

// Action describes a CLI or skills self-update step, for dry-run reporting.
type Action struct {
	Target  string // cli | skills
	Command string
	Detail  string
}

// Runner executes an external command; tests inject a fake runner.
type Runner func(name string, args ...string) error

// DefaultRunner runs a command via exec.LookPath.
func DefaultRunner(name string, args ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s not found on PATH: install %s to run this step", name, name)
	}
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PlanSkillsFallback returns the action that would install flexspec skills via npx (fallback only).
func PlanSkillsFallback() Action {
	return Action{
		Target:  "skills",
		Command: "npx skills add " + skillsPackage + " --global",
		Detail:  "install flexspec skills via npx (no supported agent detected)",
	}
}

// ApplySkillsFallback runs npx skills to install flexspec skills globally (fallback path).
func ApplySkillsFallback(run Runner) (Action, error) {
	if run == nil {
		run = DefaultRunner
	}
	action := PlanSkillsFallback()
	if err := run("npx", "skills", "add", skillsPackage, "--global"); err != nil {
		return action, fmt.Errorf("npx skills add: %w", err)
	}
	return action, nil
}
