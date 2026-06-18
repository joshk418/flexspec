//go:build !windows

package selfupdate

import (
	"fmt"
	"os"
	"syscall"
)

// reexecPlatform replaces the current process image via syscall.Exec; does NOT return on success.
func reexecPlatform(exe string, args []string) error {
	if err := syscall.Exec(exe, append([]string{exe}, args...), os.Environ()); err != nil {
		return fmt.Errorf("exec %s: %w", exe, err)
	}
	return nil
}
