/*
Copyright © 2026 Josh Kyte
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all specs in the project",
	Long: `List all specs in the configured specs directory.

Reads specs_dir from .flexspec/config.yaml, then for each spec folder
(NNN-slug/README.md) prints name, description, status, and spec_type from
YAML frontmatter. For expanded specs, also lists tasks under tasks/ with
id, name, and status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		cfg, err := config.Load(root)
		if err != nil {
			return err
		}

		entries, err := spec.List(root, cfg)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if len(entries) == 0 {
			if _, err := fmt.Fprintf(out, "No specs in %s\n", cfg.SpecsDir); err != nil {
				return err
			}

			return nil
		}

		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "ID\tNAME\tDESCRIPTION\tSTATUS\tTYPE"); err != nil {
			return err
		}

		for _, e := range entries {
			_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				e.ID,
				displayOrDash(e.Meta.Name),
				displayOrDash(e.Meta.Description),
				displayOrDash(e.Meta.Status),
				displayOrDash(e.Meta.SpecType),
			)
			if err != nil {
				return err
			}

			for _, t := range e.Tasks {
				id := displayOrDash(t.Meta.ID)
				if id == "-" && t.File != "" {
					id = strings.TrimSuffix(t.File, ".md")
				}
				_, err := fmt.Fprintf(w, " \t%s\t%s\t%s\t\n",
					id,
					displayOrDash(t.Meta.Name),
					displayOrDash(t.Meta.Status),
				)
				if err != nil {
					return err
				}
			}
		}
		return w.Flush()
	},
}

func displayOrDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

func init() {
	rootCmd.AddCommand(listCmd)
}
