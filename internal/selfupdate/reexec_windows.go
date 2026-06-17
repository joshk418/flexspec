//go:build windows

package selfupdate

import (
	"fmt"
	"os"
	"os/exec"
)

// reexecPlatform spawns the new binary as a child because Windows has no syscall.Exec.
func reexecPlatform(exe string, args []string) error {
	cmd := exec.Command(exe, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("spawn %s: %w", exe, err)
	}
	// Detach the child so it survives the parent's imminent os.Exit.
	if err := cmd.Process.Release(); err != nil {
		return fmt.Errorf("release child: %w", err)
	}
	return nil
}
