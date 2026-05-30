/*
Copyright © 2026 Josh Kyte
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

var newTemplate string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new spec file",
	Long: `Create a new spec file in the user defined specs directory (defined in .flexspec/config.yaml).
This will create an autoincrementing sequence number spec directory with the spec template:

Simple:
- <specs_dir>/001-<spec_name>/
- <specs_dir>/001-<spec_name>/README.md

Expanded:
- <specs_dir>/002-<spec_name>/
- <specs_dir>/002-<spec_name>/README.md
- <specs_dir>/002-<spec_name>/tasks/`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		cfg, err := config.Load(root)
		if err != nil {
			return err
		}

		template, err := resolveNewTemplate(cfg)
		if err != nil {
			return err
		}

		slug, err := spec.Slugify(strings.Join(args, " "))
		if err != nil {
			return err
		}

		result, err := spec.Create(root, cfg, slug, template)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Created spec %s\n", result.DirName); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  path: %s\n", result.SpecPath); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  template: %s\n", template); err != nil {
			return err
		}
		return nil
	},
}

func resolveNewTemplate(cfg config.Config) (string, error) {
	if newTemplate != "" {
		if newTemplate != "simple" && newTemplate != "expanded" {
			return "", fmt.Errorf("invalid template %q; must be simple or expanded", newTemplate)
		}
		return newTemplate, nil
	}
	if cfg.SpecTemplate != "" {
		if cfg.SpecTemplate != "simple" && cfg.SpecTemplate != "expanded" {
			return "", fmt.Errorf("invalid spec_template %q in config; must be simple or expanded", cfg.SpecTemplate)
		}
		return cfg.SpecTemplate, nil
	}
	return "simple", nil
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&newTemplate, "template", "t", "", "Template to use: simple or expanded")
}
