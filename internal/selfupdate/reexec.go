package selfupdate

import (
	"fmt"
	"os"
)

// reexecFn is the platform-specific exec/spawn function; tests override it to avoid replacing the test binary.
var reexecFn = reexecPlatform

// ReexecSelf restarts the freshly-updated binary (Unix execs in place; Windows spawns a child).
func ReexecSelf(args ...string) error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate current executable: %w", err)
	}
	return reexecFn(exe, args)
}

// ResumeArgs builds argv for the re-exec'd binary so it skips the binary-download step.
func ResumeArgs(prevVersion string, doSkills, doMigrate, force bool, only []string, skillsMethod string, skillsProject bool) []string {
	args := []string{
		"update",
		"--self-update-resume", prevVersion,
	}
	if doSkills {
		args = append(args, "--skills")
		if skillsMethod != "" && skillsMethod != "auto" {
			args = append(args, "--skills-method", skillsMethod)
		}
		if skillsProject {
			args = append(args, "--project")
		}
	}
	if doMigrate {
		args = append(args, "--migrate")
		if force {
			args = append(args, "--force")
		}
		for _, id := range only {
			args = append(args, "--only", id)
		}
	}
	return args
}
