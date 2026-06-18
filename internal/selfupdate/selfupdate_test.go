package selfupdate

import (
	"strings"
	"testing"
)

func TestPlanSkillsFallback(t *testing.T) {
	a := PlanSkillsFallback()
	if !strings.Contains(a.Command, "npx skills") {
		t.Fatalf("command = %q", a.Command)
	}
	if a.Target != "skills" {
		t.Fatalf("target = %q", a.Target)
	}
}

func TestApplySkillsFallback_invokesRunner(t *testing.T) {
	var args []string
	run := func(name string, a ...string) error {
		args = a
		return nil
	}
	_, err := ApplySkillsFallback(run)
	if err != nil {
		t.Fatal(err)
	}
	if len(args) < 3 || args[0] != "skills" || args[1] != "add" {
		t.Fatalf("args = %v", args)
	}
}

func TestApplySkillsFallback_runnerError(t *testing.T) {
	_, err := ApplySkillsFallback(func(name string, a ...string) error {
		return &execError{msg: "npx failed"}
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

type execError struct{ msg string }

func (e *execError) Error() string { return e.msg }
