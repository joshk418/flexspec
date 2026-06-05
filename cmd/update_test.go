package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveUpdateSteps_defaultAll(t *testing.T) {
	updateCLI, updateSkills, updateMigrate = false, false, false
	cli, skills, mig := resolveUpdateSteps()
	if !cli || !skills || !mig {
		t.Fatalf("want all true, got cli=%v skills=%v migrate=%v", cli, skills, mig)
	}
}

func TestResolveUpdateSteps_singleFlag(t *testing.T) {
	updateCLI, updateSkills, updateMigrate = true, false, false
	cli, skills, mig := resolveUpdateSteps()
	if !cli || skills || mig {
		t.Fatalf("want only cli, got cli=%v skills=%v migrate=%v", cli, skills, mig)
	}
}

func TestUpdateCmd_dryRunNoRunner(t *testing.T) {
	root := t.TempDir()
	writeValidateFixture(t, root)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	updateCLI, updateSkills, updateMigrate = false, false, true
	updateDryRun, updateCheck, updateForce = true, false, false
	updateOnly = []string{"status-rename"}
	updateOnly = nil
	updateRunner = func(name string, args ...string) error {
		t.Fatal("runner should not be invoked in dry-run")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if out.Len() == 0 {
		t.Fatal("expected output")
	}
	for _, want := range []string{"MIGRATION", "PATH", "KIND", "DETAIL"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q\n%s", want, out.String())
		}
	}
}

func TestUpdateCmd_dryRunSelfUpdateHeaders(t *testing.T) {
	root := t.TempDir()
	writeValidateFixture(t, root)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	updateCLI, updateSkills, updateMigrate = true, true, false
	updateDryRun, updateCheck, updateForce = true, false, false
	updateOnly = nil
	updateRunner = func(name string, args ...string) error {
		t.Fatal("runner should not be invoked in dry-run")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	for _, want := range []string{"TARGET", "COMMAND", "ACTION", "DETAIL", "plan"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q\n%s", want, out.String())
		}
	}
	assertBefore(t, out.String(), "go install", "npx skills")
}

func TestUpdateCmd_defaultDryRunReportsSelfUpdateBeforeMigrations(t *testing.T) {
	root := t.TempDir()
	writeValidateFixture(t, root)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	updateCLI, updateSkills, updateMigrate = false, false, false
	updateDryRun, updateCheck, updateForce = true, false, false
	updateOnly = nil
	updateRunner = func(name string, args ...string) error {
		t.Fatal("runner should not be invoked in dry-run")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}

	got := out.String()
	for _, want := range []string{"TARGET", "COMMAND", "ACTION", "DETAIL", "go install", "npx skills", "MIGRATION"} {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\n%s", want, got)
		}
	}
	assertBefore(t, got, "go install", "npx skills")
	assertBefore(t, got, "npx skills", "MIGRATION")
}

func TestUpdateCmd_defaultApplyRunsCLIThenSkillsBeforeMigrations(t *testing.T) {
	root := t.TempDir()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	updateCLI, updateSkills, updateMigrate = false, false, false
	updateDryRun, updateCheck, updateForce = false, false, false
	updateOnly = nil
	var calls []string
	updateRunner = func(name string, args ...string) error {
		calls = append(calls, name+" "+strings.Join(args, " "))
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}

	if len(calls) != 2 {
		t.Fatalf("runner calls = %v, want 2 calls", calls)
	}
	if !strings.HasPrefix(calls[0], "go install github.com/joshk418/flexspec@latest") {
		t.Fatalf("first runner call = %q, want go install", calls[0])
	}
	if got, want := calls[1], "go run github.com/joshk418/flexspec@latest update --skills --migrate"; got != want {
		t.Fatalf("runner order = %q, want %q", got, want)
	}
	got := out.String()
	if !strings.Contains(got, "go install") {
		t.Fatalf("output missing CLI action:\n%s", got)
	}
	if strings.Contains(got, "npx skills") || strings.Contains(got, "MIGRATION") {
		t.Fatalf("parent process should delegate remaining output to latest CLI:\n%s", got)
	}
}

func TestLatestUpdateArgs_preservesMigrationFlags(t *testing.T) {
	updateForce = true
	updateOnly = []string{"config-keys", "task-count"}
	t.Cleanup(func() {
		updateForce = false
		updateOnly = nil
	})

	got := strings.Join(latestUpdateArgs(true, true), " ")
	want := "--skills --migrate --force --only config-keys --only task-count"
	if got != want {
		t.Fatalf("latest update args = %q, want %q", got, want)
	}
}

func TestUpdateCmd_checkClean(t *testing.T) {
	root := t.TempDir()
	writeValidateFixture(t, root)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	updateCLI, updateSkills, updateMigrate = false, false, false
	updateDryRun, updateCheck, updateForce = false, true, false
	updateOnly = []string{"status-rename"}
	updateRunner = func(name string, args ...string) error {
		t.Fatal("runner should not be invoked in check")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("check on clean project: %v\n%s", err, out.String())
	}
	if strings.Contains(out.String(), "go install") || strings.Contains(out.String(), "npx skills") {
		t.Fatalf("check should not print self-update actions:\n%s", out.String())
	}
}

func TestUpdateCmd_checkPending(t *testing.T) {
	root := t.TempDir()
	flexDir := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(flexDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeValidateFixture(t, root)
	// Config without spec_template triggers pending migration.
	_ = os.WriteFile(filepath.Join(flexDir, "config.yaml"), []byte("specs_dir: specs\nalways_one_shot: false\n"), 0o644)

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })

	updateCLI, updateSkills, updateMigrate = false, false, false
	updateDryRun, updateCheck, updateForce = false, true, false
	updateOnly = []string{"config-keys"}
	updateRunner = func(name string, args ...string) error {
		t.Fatal("runner should not be invoked in check")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err == nil {
		t.Fatalf("expected pending error, output:\n%s", out.String())
	}
	if strings.Contains(out.String(), "go install") || strings.Contains(out.String(), "npx skills") {
		t.Fatalf("check should not print self-update actions:\n%s", out.String())
	}
}

func assertBefore(t *testing.T, got, first, second string) {
	t.Helper()
	firstAt := strings.Index(got, first)
	if firstAt < 0 {
		t.Fatalf("output missing %q\n%s", first, got)
	}
	secondAt := strings.Index(got, second)
	if secondAt < 0 {
		t.Fatalf("output missing %q\n%s", second, got)
	}
	if firstAt >= secondAt {
		t.Fatalf("expected %q before %q\n%s", first, second, got)
	}
}
