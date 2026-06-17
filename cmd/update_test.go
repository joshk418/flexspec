package cmd

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/joshk418/flexspec/internal/selfupdate"
	"github.com/joshk418/flexspec/internal/skills"
)

// resetUpdateFlags zeroes all update command package-level state between tests.
func resetUpdateFlags() {
	updateCLI = false
	updateSkills = false
	updateMigrate = false
	updateDryRun = false
	updateCheck = false
	updateForce = false
	updateOnly = nil
	updateResume = ""
	updateNoReexec = false
	updateSkillsMethod = ""
	updateSkillsProject = false
	updateApplyBinary = nil
	updateReexec = nil
	updateLatestRelease = nil
	updateSkillsRunner = nil
	exitAfterReexec = func(int) {} // tests must not actually exit
}

func chdirTemp(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(cwd) })
	return root
}

// fakeSkillsFS mounts a tiny skills FS for update command tests.
func fakeSkillsFS(t *testing.T) {
	t.Helper()
	orig := SkillsFS
	t.Cleanup(func() { SkillsFS = orig })
	SkillsFS = fstest.MapFS{
		"flexspec/SKILL.md":         &fstest.MapFile{Data: []byte("# flexspec\n")},
		"flexspec-charter/SKILL.md": &fstest.MapFile{Data: []byte("# charter\n")},
	}
}

// fakeHomeWithAgents creates a temp HOME with the given agent root dirs.
func fakeHomeWithAgents(t *testing.T, agentRoots ...string) string {
	t.Helper()
	home := t.TempDir()
	for _, root := range agentRoots {
		if err := os.MkdirAll(filepath.Join(home, root), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return home
}

func TestResolveUpdateSteps_defaultAll(t *testing.T) {
	resetUpdateFlags()
	updateCLI, updateSkills, updateMigrate = false, false, false
	cli, sk, mig := resolveUpdateSteps()
	if !cli || !sk || !mig {
		t.Fatalf("want all true, got cli=%v skills=%v migrate=%v", cli, sk, mig)
	}
}

func TestResolveUpdateSteps_singleFlag(t *testing.T) {
	resetUpdateFlags()
	updateCLI, updateSkills, updateMigrate = true, false, false
	cli, sk, mig := resolveUpdateSteps()
	if !cli || sk || mig {
		t.Fatalf("want only cli, got cli=%v skills=%v migrate=%v", cli, sk, mig)
	}
}

func TestUpdateCmd_dryRunMigrateOnly(t *testing.T) {
	resetUpdateFlags()
	root := chdirTemp(t)
	writeValidateFixture(t, root)

	updateMigrate = true
	updateDryRun = true
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		t.Fatal("ApplyBinary should not be called in dry-run")
		return selfupdate.ApplyResult{}, nil
	}
	updateReexec = func(_ ...string) error {
		t.Fatal("reexec should not be called in dry-run")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	for _, want := range []string{"MIGRATION", "PATH", "KIND", "DETAIL"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q\n%s", want, out.String())
		}
	}
}

func TestUpdateCmd_dryRunCLIPlan(t *testing.T) {
	resetUpdateFlags()
	root := chdirTemp(t)
	_ = root

	updateCLI = true
	updateDryRun = true
	updateLatestRelease = func(_ context.Context) (selfupdate.Release, error) {
		return selfupdate.Release{Tag: "v0.3.5", Version: "0.3.5", Assets: []selfupdate.Asset{{Name: "flexspec_0.3.5_linux_amd64.tar.gz"}}}, nil
	}
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		t.Fatal("ApplyBinary should not be called in dry-run")
		return selfupdate.ApplyResult{}, nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	for _, want := range []string{"TARGET", "COMMAND", "ACTION", "DETAIL", "v0.3.5", "plan"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output missing %q\n%s", want, out.String())
		}
	}
}

func TestUpdateCmd_defaultApplyDownloadsAndReexecs(t *testing.T) {
	resetUpdateFlags()
	root := chdirTemp(t)
	writeValidateFixture(t, root)
	fakeSkillsFS(t)

	var applyCalled, reexecCalled bool
	var reexecArgs []string
	updateApplyBinary = func(_ context.Context, cur string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		applyCalled = true
		if cur != version {
			t.Errorf("currentVersion = %q, want %q", cur, version)
		}
		return selfupdate.ApplyResult{Applied: true, ToVersion: "0.3.5"}, nil
	}
	updateReexec = func(args ...string) error {
		reexecCalled = true
		reexecArgs = args
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}

	if !applyCalled {
		t.Fatal("ApplyBinary should be called by default")
	}
	if !reexecCalled {
		t.Fatal("ReexecSelf should be called after a successful apply")
	}
	// Re-exec args should include update --self-update-resume <ver> --skills --migrate.
	joined := strings.Join(reexecArgs, " ")
	if !strings.Contains(joined, "--self-update-resume "+version) {
		t.Errorf("reexec args missing --self-update-resume: %q", joined)
	}
	if !strings.Contains(joined, "--skills") || !strings.Contains(joined, "--migrate") {
		t.Errorf("reexec args should include --skills and --migrate: %q", joined)
	}
}

func TestUpdateCmd_alreadyLatestRunsSkillsAndMigrateInProcess(t *testing.T) {
	resetUpdateFlags()
	root := chdirTemp(t)
	writeValidateFixture(t, root)
	fakeSkillsFS(t)

	// Fake HOME with claude-code so embedded install runs instead of npx fallback.
	home := fakeHomeWithAgents(t, ".claude")
	// Set HOME for Unix-like os.UserHomeDir.
	t.Setenv("HOME", home)
	// On Windows, os.UserHomeDir reads %USERPROFILE%.
	t.Setenv("USERPROFILE", home)

	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		return selfupdate.ApplyResult{Applied: false}, nil // already latest
	}
	updateReexec = func(_ ...string) error {
		t.Fatal("reexec should NOT be called when already latest")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}

	got := out.String()
	// Skills step should run in-process.
	if !strings.Contains(got, "AGENT") && !strings.Contains(got, "claude-code") {
		t.Errorf("output should mention skills install:\n%s", got)
	}
	// Migrations should also run in-process.
	if !strings.Contains(got, "MIGRATION") && !strings.Contains(got, "0 pending change(s)") && !strings.Contains(got, "change(s)") {
		t.Errorf("output should mention migrations:\n%s", got)
	}
}

func TestUpdateCmd_resumePathRunsSkillsAndMigrate(t *testing.T) {
	resetUpdateFlags()
	root := chdirTemp(t)
	writeValidateFixture(t, root)
	fakeSkillsFS(t)

	home := fakeHomeWithAgents(t, ".claude")
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	updateResume = "0.3.4" // simulate being the re-exec'd new binary
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		t.Fatal("ApplyBinary should not be called on the resume path")
		return selfupdate.ApplyResult{}, nil
	}
	updateReexec = func(_ ...string) error {
		t.Fatal("reexec should not be called on the resume path")
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}

	got := out.String()
	if !strings.Contains(got, "Restarted as v"+version+" (from v0.3.4)") {
		t.Errorf("output missing resume banner:\n%s", got)
	}
	if !strings.Contains(got, "Update complete: v0.3.4 -> v"+version) {
		t.Errorf("output missing completion summary:\n%s", got)
	}
}

func TestUpdateCmd_noReexecFallsThroughToSkillsAndMigrate(t *testing.T) {
	resetUpdateFlags()
	root := chdirTemp(t)
	writeValidateFixture(t, root)
	fakeSkillsFS(t)

	home := fakeHomeWithAgents(t, ".claude")
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	var reexecCalled bool
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		return selfupdate.ApplyResult{Applied: true, ToVersion: "0.3.5"}, nil
	}
	updateReexec = func(_ ...string) error {
		reexecCalled = true
		return nil
	}
	updateNoReexec = true

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}

	if reexecCalled {
		t.Fatal("reexec should not be called when --no-reexec is set")
	}
	// Skills + migrate should run in-process under the old binary.
	got := out.String()
	if !strings.Contains(got, "AGENT") && !strings.Contains(got, "claude-code") {
		t.Errorf("output should mention skills install:\n%s", got)
	}
}

func TestUpdateCmd_defaultApplySkipsMigrationsOutsideFlexspecDir(t *testing.T) {
	resetUpdateFlags()
	chdirTemp(t) // no .flexspec/
	fakeSkillsFS(t)

	home := fakeHomeWithAgents(t, ".claude")
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	var applyCalled bool
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		applyCalled = true
		return selfupdate.ApplyResult{Applied: true, ToVersion: "0.3.5"}, nil
	}
	updateReexec = func(args ...string) error {
		// Outside a flexspec dir, re-exec should include --skills but not --migrate.
		joined := strings.Join(args, " ")
		if strings.Contains(joined, "--migrate") {
			t.Errorf("reexec args should not include --migrate outside flexspec dir: %q", joined)
		}
		if !strings.Contains(joined, "--skills") {
			t.Errorf("reexec args should include --skills: %q", joined)
		}
		return nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if !applyCalled {
		t.Fatal("ApplyBinary should still be called outside a flexspec dir")
	}
}

func TestUpdateCmd_checkClean(t *testing.T) {
	resetUpdateFlags()
	root := chdirTemp(t)
	writeValidateFixture(t, root)

	updateCheck = true
	// Scope to status-rename because the fixture lacks config/glossary migration inputs.
	updateOnly = []string{"status-rename"}
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		t.Fatal("ApplyBinary should not be called in --check")
		return selfupdate.ApplyResult{}, nil
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
	resetUpdateFlags()
	root := chdirTemp(t)
	flexDir := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(flexDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeValidateFixture(t, root)
	// Config without spec_template triggers a pending migration.
	_ = os.WriteFile(filepath.Join(flexDir, "config.yaml"), []byte("specs_dir: specs\nalways_one_shot: false\n"), 0o644)

	updateCheck = true
	updateOnly = []string{"config-keys"}
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		t.Fatal("ApplyBinary should not be called in --check")
		return selfupdate.ApplyResult{}, nil
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

func TestUpdateCmd_skillsMethodNpx(t *testing.T) {
	resetUpdateFlags()
	chdirTemp(t)
	fakeSkillsFS(t)

	updateSkills = true
	updateSkillsMethod = "npx"
	var npxCalled bool
	var npxArgs []string
	updateSkillsRunner = func(name string, a ...string) error {
		npxCalled = true
		npxArgs = a
		return nil
	}
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		return selfupdate.ApplyResult{Applied: false}, nil
	}
	updateReexec = func(_ ...string) error { return nil }

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if !npxCalled {
		t.Fatal("npx fallback should be called when --skills-method=npx")
	}
	if len(npxArgs) < 3 || npxArgs[0] != "skills" || npxArgs[1] != "add" {
		t.Fatalf("npx args = %v", npxArgs)
	}
}

func TestUpdateCmd_skillsMethodEmbeddedNoAgentsPrintsFallback(t *testing.T) {
	resetUpdateFlags()
	chdirTemp(t)
	fakeSkillsFS(t)

	updateSkills = true
	updateSkillsMethod = "embedded"
	// Empty HOME so DetectAgents finds nothing.
	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())

	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		return selfupdate.ApplyResult{Applied: false}, nil
	}
	updateReexec = func(_ ...string) error { return nil }

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "No supported coding agent detected") {
		t.Errorf("output should print fallback instruction:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "npx skills add") {
		t.Errorf("output should mention manual npx command:\n%s", out.String())
	}
}

func TestUpdateCmd_skillsDryRunPrintsPlan(t *testing.T) {
	resetUpdateFlags()
	chdirTemp(t)
	fakeSkillsFS(t)

	home := fakeHomeWithAgents(t, ".claude")
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	updateSkills = true
	updateDryRun = true
	updateSkillsMethod = "embedded"
	updateApplyBinary = func(_ context.Context, _ string, _ selfupdate.ApplyOpts) (selfupdate.ApplyResult, error) {
		t.Fatal("ApplyBinary should not be called in dry-run")
		return selfupdate.ApplyResult{}, nil
	}

	var out bytes.Buffer
	updateCmd.SetOut(&out)
	updateCmd.SetErr(&out)
	if err := updateCmd.RunE(updateCmd, nil); err != nil {
		t.Fatalf("update: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "embedded -> claude-code") {
		t.Errorf("output should mention embedded install plan:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "plan") {
		t.Errorf("output should be plan not exec:\n%s", out.String())
	}
}

// sanity: the skills package's Agent type is usable from cmd tests.
var _ = skills.ScopeGlobal
