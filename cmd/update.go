package cmd

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/clioutput"
	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/migrate"
	"github.com/joshk418/flexspec/internal/selfupdate"
)

var (
	updateCLI     bool
	updateSkills  bool
	updateMigrate bool
	updateDryRun  bool
	updateCheck   bool
	updateForce   bool
	updateOnly    []string
	updateRunner  selfupdate.Runner
)

// updateCmd upgrades flexspec, skills, and in-project files.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update flexspec CLI, skills, and project files",
	Long: `Update brings your environment and project up to date.

By default (no step flags), flexspec update:
  1. Upgrades the flexspec CLI via go install
  2. Reinstalls flexspec skills via npx
  3. Runs in-project migrations (spec statuses, templates, config, charter checks)

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
	dryPlan := updateDryRun || updateCheck

	if updateCheck {
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

	runner := updateRunner
	if dryPlan {
		runner = nil
	} else if runner == nil {
		runner = selfupdate.DefaultRunner
	}

	var selfUpdateActions []selfupdate.Action

	if doCLI {
		action := selfupdate.PlanCLI(version)
		if apply && runner != nil {
			action, err = selfupdate.ApplyCLI(version, runner)
			if err != nil {
				return err
			}
		}
		selfUpdateActions = append(selfUpdateActions, action)
	}

	if doCLI && apply && (doSkills || doMigrate) {
		if err := writeSelfUpdateActions(cmd.OutOrStdout(), selfUpdateActions, apply); err != nil {
			return err
		}
		if err := selfupdate.ApplyLatestUpdate(runner, latestUpdateArgs(doSkills, doMigrate)...); err != nil {
			return err
		}
		return nil
	}

	if doSkills {
		action := selfupdate.PlanSkills()
		if apply && runner != nil {
			action, err = selfupdate.ApplySkills(runner)
			if err != nil {
				return err
			}
		}
		selfUpdateActions = append(selfUpdateActions, action)
	}

	if len(selfUpdateActions) > 0 {
		if err := writeSelfUpdateActions(cmd.OutOrStdout(), selfUpdateActions, apply); err != nil {
			return err
		}
	}

	if doMigrate {
		var migrationChanges []migrate.Change
		cfg, migs, err := loadUpdateMigrations(root)
		if err != nil {
			return err
		}
		if apply {
			migrationChanges, err = migrate.Apply(root, cfg, migs)
		} else {
			migrationChanges, err = migrate.Plan(root, cfg, migs)
		}
		if err != nil {
			return err
		}
		if err := migrate.WriteChanges(cmd.OutOrStdout(), migrationChanges); err != nil {
			return err
		}
	}

	return nil
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

func latestUpdateArgs(doSkills, doMigrate bool) []string {
	var args []string
	if doSkills {
		args = append(args, "--skills")
	}
	if doMigrate {
		args = append(args, "--migrate")
		if updateForce {
			args = append(args, "--force")
		}
		for _, id := range updateOnly {
			args = append(args, "--only", id)
		}
	}
	return args
}

func resolveUpdateSteps() (cli, skills, migrateStep bool) {
	if !updateCLI && !updateSkills && !updateMigrate {
		return true, true, true
	}
	return updateCLI, updateSkills, updateMigrate
}

func writeSelfUpdateActions(w io.Writer, actions []selfupdate.Action, applying bool) error {
	actionLabel := "plan"
	if applying {
		actionLabel = "exec"
	}
	rows := make([][]string, len(actions))
	for i, a := range actions {
		rows[i] = []string{a.Target, a.Command, actionLabel, a.Detail}
	}
	return clioutput.WriteTable(w,
		[]string{"TARGET", "COMMAND", "ACTION", "DETAIL"},
		rows,
	)
}

func embeddedTemplatesFS() (fs.FS, error) {
	if TemplatesFS == nil {
		return nil, nil
	}
	if _, err := fs.ReadFile(TemplatesFS, filepath.Join(embedRootDir, "flexspec-simple.md")); err != nil {
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
	updateCmd.Flags().BoolVar(&updateCLI, "cli", false, "Update the installed flexspec CLI (go install)")
	updateCmd.Flags().BoolVar(&updateSkills, "skills", false, "Reinstall flexspec agent skills (npx)")
	updateCmd.Flags().BoolVar(&updateMigrate, "migrate", false, "Run in-project migrations only")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Preview changes without writing or executing")
	updateCmd.Flags().BoolVar(&updateCheck, "check", false, "Exit 1 if migrations are pending (detect only)")
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "Overwrite differing template files on migrate")
	updateCmd.Flags().StringSliceVar(&updateOnly, "only", nil, "Run only these migration ids (repeatable)")
}
