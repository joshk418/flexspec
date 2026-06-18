//go:build windows

package selfupdate

import (
	"fmt"
	"os"
	"os/exec"
)

// reexecPlatform runs the new binary as a child because Windows has no syscall.Exec.
func reexecPlatform(exe string, args []string) error {
	cmd := exec.Command(exe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %s: %w", exe, err)
	}
	return nil
}
