package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/validate"
)

var validateStrict bool

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate FlexSpec project files for structural problems",
	Long: `Validate FlexSpec project files for structural problems.

Checks .flexspec/config.yaml, charter, templates, and spec directories for
issues that can cause list, new, or agent skills to fail. Exits with code 1
when any error-severity finding is reported; warnings alone do not fail.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		opts := validate.Options{Strict: validateStrict}
		findings := validate.RunAll(root, opts)

		if err := validate.WriteFindings(cmd.OutOrStdout(), findings); err != nil {
			return err
		}
		if validate.HasErrors(findings) {
			return fmt.Errorf("validation failed")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "Enable additional semantic checks (reserved for future use)")
}
