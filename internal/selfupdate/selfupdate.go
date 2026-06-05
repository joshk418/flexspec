package selfupdate

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	cliModule     = "github.com/joshk418/flexspec@latest"
	skillsPackage = "joshk418/flexspec"
)

// Action describes a CLI or skills self-update step.
type Action struct {
	Target  string // cli | skills
	Command string
	Detail  string
}

// Runner executes an external command. Tests inject a fake runner.
type Runner func(name string, args ...string) error

// DefaultRunner runs a command via exec.LookPath.
func DefaultRunner(name string, args ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("%s not found on PATH: install %s to run this step", name, name)
	}
	cmd := exec.Command(path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PlanCLI returns the action that would upgrade the flexspec binary.
func PlanCLI(installedVersion string) Action {
	return Action{
		Target:  "cli",
		Command: "go install " + cliModule,
		Detail:  "installed " + installedVersion,
	}
}

// ApplyCLI runs go install to upgrade flexspec.
func ApplyCLI(installedVersion string, run Runner) (Action, error) {
	if run == nil {
		run = DefaultRunner
	}
	action := PlanCLI(installedVersion)
	if err := run("go", "install", cliModule); err != nil {
		return action, fmt.Errorf("go install %s: %w", cliModule, err)
	}
	return action, nil
}

// ApplyLatestUpdate runs the latest flexspec update command for selected steps.
func ApplyLatestUpdate(run Runner, args ...string) error {
	if run == nil {
		run = DefaultRunner
	}
	runArgs := append([]string{"run", cliModule, "update"}, args...)
	if err := run("go", runArgs...); err != nil {
		return fmt.Errorf("go %s: %w", strings.Join(runArgs, " "), err)
	}
	return nil
}

// PlanSkills returns the action that would reinstall flexspec skills.
func PlanSkills() Action {
	return Action{
		Target:  "skills",
		Command: "npx skills add " + skillsPackage + " --global",
		Detail:  "reinstall flexspec skills",
	}
}

// ApplySkills runs npx skills to reinstall flexspec skills globally.
func ApplySkills(run Runner) (Action, error) {
	if run == nil {
		run = DefaultRunner
	}
	action := PlanSkills()
	if err := run("npx", "skills", "add", skillsPackage, "--global"); err != nil {
		return action, fmt.Errorf("npx skills add: %w", err)
	}
	return action, nil
}
