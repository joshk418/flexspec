package cmd

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/clioutput"
	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/migrate"
	"github.com/joshk418/flexspec/internal/selfupdate"
	"github.com/joshk418/flexspec/internal/skills"
)

var (
	updateCLI           bool
	updateSkills        bool
	updateMigrate       bool
	updateDryRun        bool
	updateCheck         bool
	updateForce         bool
	updateOnly          []string
	updateResume        string
	updateNoReexec      bool
	updateSkillsMethod  string
	updateSkillsProject bool

	// updateApplyBinary is the test seam for selfupdate.ApplyBinary.
	updateApplyBinary func(ctx context.Context, currentVersion string, opts selfupdate.ApplyOpts) (selfupdate.ApplyResult, error)
	// updateReexec is the test seam for selfupdate.ReexecSelf.
	updateReexec func(args ...string) error
	// updateLatestRelease is the test seam for selfupdate.LatestRelease (dry-run only).
	updateLatestRelease func(ctx context.Context) (selfupdate.Release, error)
	// updateSkillsRunner is the npx fallback runner (tests inject a fake).
	updateSkillsRunner selfupdate.Runner
	// exitAfterReexec is the seam over os.Exit(0) after a successful Windows re-exec; tests override it.
	exitAfterReexec = os.Exit
)

// updateCmd upgrades flexspec, skills, and in-project files.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update flexspec CLI, skills, and project files",
	Long: `Update brings your environment and project up to date.

By default (no step flags), flexspec update:
  1. Checks the latest GitHub release and compares it to the installed CLI
  2. If newer: downloads the matching prebuilt binary, verifies its SHA256
     against checksums.txt, atomically swaps the running executable, and
     re-execs into the new binary to finish the update under the new code
  3. Installs the embedded agent skills into each detected coding agent's
     skills directory (claude-code, cursor, codex, opencode, cline). If no
     supported agent is detected, falls back to 'npx skills add --global'.
  4. Runs in-project migrations (spec statuses, templates, config, charter,
     glossary, task counts, type backfill) - only when inside a .flexspec dir

No Go toolchain or Node install is required for the binary update. Skills
install uses the embedded skills tree (no Node needed) unless no agent is
detected, in which case the npx fallback requires Node.

Use --dry-run to preview without writing or executing external commands.
Use --check for a CI gate: detect-only, exit 1 when migrations are pending.

Step flags (--cli, --skills, --migrate) restrict which steps run.`,
	RunE: runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) error {
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working directory: %w", err)
	}

	doCLI, doSkills, doMigrate := resolveUpdateSteps()
	if doMigrate && !isFlexspecDir(root) {
		doMigrate = false
	}

	apply := !updateDryRun && !updateCheck

	// --check: unchanged CI gate (migration detect-only, exit 1 if pending).
	if updateCheck {
		return runUpdateCheck(cmd, root, doMigrate)
	}

	// Resume path: this process is the freshly-swapped binary, so skip binary download.
	if updateResume != "" {
		return runUpdateResume(cmd, root, updateResume, doSkills, doMigrate)
	}

	// 1. BINARY: check + download + swap + (optional) re-exec.
	if doCLI {
		if updateDryRun {
			if err := printBinaryPlan(cmd, version); err != nil {
				return err
			}
		} else if apply {
			applied, err := runBinaryUpdate(cmd.Context(), version, cmd.OutOrStdout())
			if err != nil {
				return err
			}
			if applied && !updateNoReexec {
				// Re-exec into the new binary; this process ends here.
				resumeArgs := selfupdate.ResumeArgs(version, doSkills, doMigrate, updateForce, updateOnly)
				if err := reexecSelf(resumeArgs...); err != nil {
					bestEffortFprintf(cmd.ErrOrStderr(), "Updated binary to latest but could not re-exec: %v\n", err)
					bestEffortFprintf(cmd.ErrOrStderr(), "Re-run `flexspec update --skills --migrate` to finish.\n")
					return nil
				}
				// Unix exec does not return; Windows spawn does, so exit the parent.
				exitAfterReexec(0)
			}
			// Already latest or --no-reexec: run skills + migrate in this process.
		}
	}

	// Run skills + migrate here when no re-exec happened.
	if err := runSkillsStep(cmd, doSkills, apply); err != nil {
		return err
	}
	if err := runMigrateStep(cmd, root, doMigrate, apply); err != nil {
		return err
	}
	return nil
}

// runBinaryUpdate applies the latest binary when needed and reports whether a swap occurred.
func runBinaryUpdate(ctx context.Context, currentVersion string, w io.Writer) (bool, error) {
	opts := selfupdate.ApplyOpts{Force: updateForce, Progress: w}
	if updateApplyBinary != nil {
		res, err := updateApplyBinary(ctx, currentVersion, opts)
		if err != nil {
			return false, err
		}
		return res.Applied, nil
	}
	res, err := selfupdate.ApplyBinary(ctx, currentVersion, opts)
	if err != nil {
		return false, err
	}
	return res.Applied, nil
}

// reexecSelf calls the test seam if set, otherwise selfupdate.ReexecSelf.
func reexecSelf(args ...string) error {
	if updateReexec != nil {
		return updateReexec(args...)
	}
	return selfupdate.ReexecSelf(args...)
}

// printBinaryPlan previews the binary step by querying the releases API without downloading assets.
func printBinaryPlan(cmd *cobra.Command, currentVersion string) error {
	out := cmd.OutOrStdout()
	var rel selfupdate.Release
	var err error
	if updateLatestRelease != nil {
		rel, err = updateLatestRelease(cmd.Context())
	} else {
		rel, err = selfupdate.LatestRelease(cmd.Context())
	}
	if err != nil {
		return fmt.Errorf("check latest release: %w", err)
	}
	actionLabel := "plan"
	rows := [][]string{
		{"cli", "github:releases/" + rel.Tag, actionLabel, "current v" + currentVersion + " -> latest v" + rel.Version},
	}
	return clioutput.WriteTable(out,
		[]string{"TARGET", "COMMAND", "ACTION", "DETAIL"},
		rows,
	)
}

// runSkillsStep installs embedded skills (auto/embedded) or falls back to npx.
func runSkillsStep(cmd *cobra.Command, doSkills, apply bool) error {
	if !doSkills {
		return nil
	}
	out := cmd.OutOrStdout()

	method := updateSkillsMethod
	if method == "" {
		method = "auto"
	}
	switch method {
	case "embedded":
		return installEmbeddedSkills(cmd, out, apply)
	case "npx":
		return installSkillsViaNpx(cmd, out, apply)
	case "auto":
		fallthrough
	default:
		// auto = embedded if any agent detected, else npx fallback.
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("resolve home directory: %w", err)
		}
		scope := skills.ScopeGlobal
		if updateSkillsProject {
			scope = skills.ScopeProject
		}
		detected := skills.DetectAgents(scope, home, ".")
		if len(detected) > 0 {
			return installEmbeddedSkillsFor(cmd, out, apply, detected, scope, home)
		}
		return installSkillsViaNpx(cmd, out, apply)
	}
}

// installEmbeddedSkills installs via the embedded FS, detecting agents itself.
func installEmbeddedSkills(cmd *cobra.Command, out io.Writer, apply bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("resolve home directory: %w", err)
	}
	scope := skills.ScopeGlobal
	if updateSkillsProject {
		scope = skills.ScopeProject
	}
	detected := skills.DetectAgents(scope, home, ".")
	return installEmbeddedSkillsFor(cmd, out, apply, detected, scope, home)
}

// installEmbeddedSkillsFor writes the embedded skills to the given agents.
func installEmbeddedSkillsFor(cmd *cobra.Command, out io.Writer, apply bool, detected []skills.Agent, scope skills.Scope, home string) error {
	if SkillsFS == nil {
		return fmt.Errorf("embedded skills FS not mounted (binary built without skills/ embed)")
	}
	if len(detected) == 0 {
		bestEffortFprintln(out, "No supported coding agent detected.")
		skills.PrintFallbackInstruction(out)
		return nil
	}
	if !apply {
		rows := make([][]string, 0, len(detected))
		for _, a := range detected {
			dir := skillsDirForDisplay(a, scope, home)
			rows = append(rows, []string{"skills", "embedded -> " + a.Name, "plan", dir})
		}
		return clioutput.WriteTable(out,
			[]string{"TARGET", "COMMAND", "ACTION", "DETAIL"},
			rows,
		)
	}
	results, err := skills.Install(SkillsFS, scope, home, ".", detected)
	if err != nil {
		return err
	}
	rows := make([][]string, len(results))
	for i, r := range results {
		rows[i] = []string{r.Agent, r.Dir, fmt.Sprintf("%d skills", r.FilesWritten), "installed"}
	}
	return clioutput.WriteTable(out,
		[]string{"AGENT", "PATH", "SKILLS", "STATUS"},
		rows,
	)
}

// skillsDirForDisplay returns the target dir string for dry-run display.
func skillsDirForDisplay(a skills.Agent, scope skills.Scope, home string) string {
	if scope == skills.ScopeProject {
		return a.ProjectDir
	}
	return home + "/" + a.GlobalDir
}

// installSkillsViaNpx shells out to the npx skills CLI (legacy fallback).
func installSkillsViaNpx(cmd *cobra.Command, out io.Writer, apply bool) error {
	action := selfupdate.PlanSkillsFallback()
	if !apply {
		rows := [][]string{{action.Target, action.Command, "plan", action.Detail}}
		return clioutput.WriteTable(out,
			[]string{"TARGET", "COMMAND", "ACTION", "DETAIL"},
			rows,
		)
	}
	runner := updateSkillsRunner
	if runner == nil {
		runner = selfupdate.DefaultRunner
	}
	_, err := selfupdate.ApplySkillsFallback(runner)
	return err
}

// runMigrateStep runs the in-project migrations.
func runMigrateStep(cmd *cobra.Command, root string, doMigrate, apply bool) error {
	if !doMigrate {
		return nil
	}
	cfg, migs, err := loadUpdateMigrations(root)
	if err != nil {
		return err
	}
	var changes []migrate.Change
	if apply {
		changes, err = migrate.Apply(root, cfg, migs)
	} else {
		changes, err = migrate.Plan(root, cfg, migs)
	}
	if err != nil {
		return err
	}
	return migrate.WriteChanges(cmd.OutOrStdout(), changes)
}

// runUpdateCheck is the CI-gate path: detect migrations, exit 1 if pending.
func runUpdateCheck(cmd *cobra.Command, root string, doMigrate bool) error {
	var migrationChanges []migrate.Change
	if doMigrate {
		cfg, migs, err := loadUpdateMigrations(root)
		if err != nil {
			return err
		}
		migrationChanges, err = migrate.Plan(root, cfg, migs)
		if err != nil {
			return err
		}
		if err := migrate.WriteChanges(cmd.OutOrStdout(), migrationChanges); err != nil {
			return err
		}
	}
	if migrate.HasApplicableChanges(migrationChanges) {
		return fmt.Errorf("migrations pending")
	}
	return nil
}

// runUpdateResume skips the binary step because this process is already the new binary.
func runUpdateResume(cmd *cobra.Command, root, prevVersion string, doSkills, doMigrate bool) error {
	out := cmd.OutOrStdout()
	bestEffortFprintf(out, "Restarted as v%s (from v%s). Finishing update.\n", version, prevVersion)

	if err := runSkillsStep(cmd, doSkills, true); err != nil {
		return err
	}
	if err := runMigrateStep(cmd, root, doMigrate, true); err != nil {
		return err
	}
	bestEffortFprintf(out, "Update complete: v%s -> v%s\n", prevVersion, version)
	return nil
}

// bestEffortFprintf writes a formatted status message and ignores write errors.
func bestEffortFprintf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

// bestEffortFprintln writes a line to w, ignoring errors.
func bestEffortFprintln(w io.Writer, args ...any) {
	_, _ = fmt.Fprintln(w, args...)
}

func loadUpdateMigrations(root string) (config.Config, []migrate.Migration, error) {
	cfg, err := config.Load(root)
	if err != nil {
		return config.Config{}, nil, err
	}

	templatesFS, err := embeddedTemplatesFS()
	if err != nil {
		return config.Config{}, nil, fmt.Errorf("mount embedded templates: %w", err)
	}
	migs := migrate.Registry(templatesFS, updateForce)
	if len(updateOnly) > 0 {
		migs, err = migrate.Select(migs, updateOnly)
		if err != nil {
			return config.Config{}, nil, err
		}
	}
	return cfg, migs, nil
}

func resolveUpdateSteps() (cli, skillsStep, migrateStep bool) {
	if !updateCLI && !updateSkills && !updateMigrate {
		return true, true, true
	}
	return updateCLI, updateSkills, updateMigrate
}

func embeddedTemplatesFS() (fs.FS, error) {
	if TemplatesFS == nil {
		return nil, nil
	}
	// fs.FS requires forward slashes; filepath.Join breaks on Windows.
	if _, err := fs.ReadFile(TemplatesFS, path.Join(embedRootDir, "flexspec-simple.md")); err != nil {
		return nil, nil
	}
	return fs.Sub(TemplatesFS, embedRootDir)
}

func isFlexspecDir(root string) bool {
	_, err := os.Stat(config.ConfigPath(root))
	return err == nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&updateCLI, "cli", false, "Update the installed flexspec CLI binary (download from GitHub Releases)")
	updateCmd.Flags().BoolVar(&updateSkills, "skills", false, "Reinstall flexspec agent skills (embedded or npx fallback)")
	updateCmd.Flags().BoolVar(&updateMigrate, "migrate", false, "Run in-project migrations only")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Preview changes without writing or executing")
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "Exit 1 if migrations are pending (detect only)")
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "Force binary re-download and template overwrite on migrate")
	updateCmd.Flags().StringSliceVar(&updateOnly, "only", nil, "Run only these migration ids (repeatable)")
	updateCmd.Flags().StringVar(&updateResume, "self-update-resume", "", "(internal) skip binary step; this process is the freshly-updated binary")
	_ = updateCmd.Flags().MarkHidden("self-update-resume")
	updateCmd.Flags().BoolVar(&updateNoReexec, "no-reexec", false, "Download + swap binary but do not re-exec; user re-runs manually")
	updateCmd.Flags().StringVar(&updateSkillsMethod, "skills-method", "auto", "Skills install method: auto|embedded|npx")
	updateCmd.Flags().BoolVar(&updateSkillsProject, "project", false, "Install skills to project-local agent dirs (default: user-wide)")
}
