/*
Copyright © 2026 Josh Kyte
*/
package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// TemplatesFS holds the embedded template tree. It is populated by main before
// Execute runs so the init command can scaffold templates into a project.
var TemplatesFS embed.FS

const (
	flexspecDir      = ".flexspec"
	configFile       = "config.yaml"
	charterFile      = "charter.md"
	templatesDir     = "templates"
	embedRootDir     = "templates"
	embedCharterPath = "templates/charter.md"
	defaultSpecs     = "specs"
	dirPerm          = 0o755
	filePerm         = 0o644
)

var (
	specsDir      string
	alwaysOneShot bool
	forceInit     bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create all of the files and directories for a project to get started with FlexSpec.",
	Long: `Create all of the files and directories for a project to get started with FlexSpec.
This will create the following files and directories:
- .flexspec/
- .flexspec/config.yaml
- .flexspec/charter.md

- .flexspec/templates/
- .flexspec/templates/README.md
- .flexspec/templates/flexspec-simple.md

- .flexspec/templates/expanded/
- .flexspec/templates/expanded/flexspec-expanded.md
- .flexspec/templates/expanded/flexspec-expanded-task.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		base := filepath.Join(root, flexspecDir)
		if _, err := os.Stat(base); err == nil && !forceInit {
			return fmt.Errorf("%s already exists; re-run with --force to overwrite", flexspecDir)
		}

		if err := os.MkdirAll(base, dirPerm); err != nil {
			return fmt.Errorf("create %s: %w", flexspecDir, err)
		}

		if err := writeConfig(filepath.Join(base, configFile)); err != nil {
			return fmt.Errorf("write %s: %w", configFile, err)
		}

		if err := writeCharter(filepath.Join(base, charterFile)); err != nil {
			return fmt.Errorf("write %s: %w", charterFile, err)
		}

		if err := copyTemplates(filepath.Join(base, templatesDir)); err != nil {
			return fmt.Errorf("scaffold templates: %w", err)
		}

		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Initialized FlexSpec in %s\n", base); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  specs directory: %s\n  always_one_shot: %t\n", specsDir, alwaysOneShot); err != nil {
			return err
		}
		return nil
	},
}

// writeConfig renders the project config.yaml with the resolved flag values.
func writeConfig(path string) error {
	content := fmt.Sprintf(`# FlexSpec project configuration.

# Directory (relative to the project root) where specs are created by
# `+"`flexspec new`"+`. The directory is created on demand if it does not exist.
specs_dir: %s

# When true, `+"`/flexspec`"+` runs all lifecycle phases back-to-back without
# stopping between them, as if --one-shot were always passed.
always_one_shot: %t

# Force which template `+"`/flexspec`"+` uses: `+"`simple`"+` or `+"`expanded`"+`.
# Intentionally has no default — leave it blank so the default workflow infers
# whether to create a simple or expanded spec. The --template flag overrides it.
spec_template:
`, specsDir, alwaysOneShot)
	return os.WriteFile(path, []byte(content), filePerm)
}

// writeCharter copies the embedded charter template into .flexspec/charter.md.
func writeCharter(path string) error {
	data, err := TemplatesFS.ReadFile(embedCharterPath)
	if err != nil {
		return fmt.Errorf("read embedded charter: %w", err)
	}
	return writeFileIfAbsent(path, data)
}

// copyTemplates walks the embedded template tree and writes it into dest,
// preserving the directory structure.
func copyTemplates(dest string) error {
	return fs.WalkDir(TemplatesFS, embedRootDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(embedRootDir, p)
		if err != nil {
			return err
		}
		if rel == charterFile {
			return nil
		}
		target := filepath.Join(dest, rel)

		if d.IsDir() {
			return os.MkdirAll(target, dirPerm)
		}

		data, err := TemplatesFS.ReadFile(p)
		if err != nil {
			return err
		}
		return writeFileIfAbsent(target, data)
	})
}

// writeFileIfAbsent writes data to path, skipping existing files unless --force
// is set so a plain re-init does not clobber user edits.
func writeFileIfAbsent(path string, data []byte) error {
	if _, err := os.Stat(path); err == nil && !forceInit {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return err
	}
	return os.WriteFile(path, data, filePerm)
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVar(&specsDir, "specs-dir", defaultSpecs, "Directory where specs are created by `flexspec new`")
	initCmd.Flags().BoolVar(&alwaysOneShot, "always-one-shot", false, "Run all /flexspec phases back-to-back without stopping")
	initCmd.Flags().BoolVar(&forceInit, "force", false, "Overwrite existing .flexspec files")
}
