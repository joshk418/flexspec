package selfupdate

import (
	"strings"
	"testing"
)

func TestPlanCLI(t *testing.T) {
	a := PlanCLI("0.2.1")
	if !strings.Contains(a.Command, "go install") {
		t.Fatalf("command = %q", a.Command)
	}
	if !strings.Contains(a.Detail, "0.2.1") {
		t.Fatalf("detail = %q", a.Detail)
	}
}

func TestApplyCLI_invokesRunner(t *testing.T) {
	var called bool
	var name string
	var args []string
	run := func(n string, a ...string) error {
		called = true
		name = n
		args = a
		return nil
	}
	_, err := ApplyCLI("0.2.1", run)
	if err != nil {
		t.Fatal(err)
	}
	if !called || name != "go" {
		t.Fatalf("runner not called correctly: %v %q %v", called, name, args)
	}
	if len(args) != 2 || args[0] != "install" || args[1] != cliModule {
		t.Fatalf("args = %v", args)
	}
}

func TestApplyCLI_runnerError(t *testing.T) {
	_, err := ApplyCLI("0.2.1", func(name string, args ...string) error {
		if name != "go" {
			t.Fatalf("name = %q", name)
		}
		return &execError{msg: "install failed"}
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

type execError struct{ msg string }

func (e *execError) Error() string { return e.msg }

func TestApplySkills_invokesRunner(t *testing.T) {
	var args []string
	run := func(name string, a ...string) error {
		args = a
		return nil
	}
	_, err := ApplySkills(run)
	if err != nil {
		t.Fatal(err)
	}
	if len(args) < 3 || args[0] != "skills" || args[1] != "add" {
		t.Fatalf("args = %v", args)
	}
}
