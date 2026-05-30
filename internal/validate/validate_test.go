package validate

import (
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func TestRun_stubChecks(t *testing.T) {
	root := t.TempDir()
	stub := func(root string, cfg config.Config, opts Options) []Finding {
		return []Finding{
			{Severity: SeverityWarning, Path: "a", Rule: "test.warn", Message: "warn"},
			{Severity: SeverityError, Path: "b", Rule: "test.err", Message: "err"},
		}
	}
	findings := Run(root, config.Config{}, Options{}, stub)
	if !HasErrors(findings) {
		t.Fatal("expected errors")
	}
	if len(findings) != 2 {
		t.Fatalf("got %d findings", len(findings))
	}
	if findings[0].Severity != SeverityError {
		t.Fatalf("first finding severity = %s, want error", findings[0].Severity)
	}
}

func TestHasErrors(t *testing.T) {
	if HasErrors([]Finding{{Severity: SeverityWarning}}) {
		t.Fatal("warnings alone should not count as errors")
	}
	if !HasErrors([]Finding{{Severity: SeverityError}}) {
		t.Fatal("expected errors")
	}
}
