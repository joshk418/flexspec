//go:build windows

package selfupdate

import (
	"os"
	"strings"
	"testing"
)

func TestReexecPlatformPropagatesChildFailure(t *testing.T) {
	if os.Getenv("FLEXSPEC_REEXEC_CHILD") == "1" {
		os.Exit(7)
	}

	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("FLEXSPEC_REEXEC_CHILD", "1")
	err = reexecPlatform(exe, []string{"-test.run=TestReexecPlatformPropagatesChildFailure"})
	if err == nil {
		t.Fatal("expected child failure to be returned")
	}
	if !strings.Contains(err.Error(), "exit status 7") {
		t.Fatalf("error = %v, want child exit status", err)
	}
}
