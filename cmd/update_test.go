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

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("check on clean project: %v\n%s", err, out.String())
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

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err == nil {
		t.Fatalf("expected pending error, output:\n%s", out.String())
	}
}
