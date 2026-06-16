package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/clioutput"
	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
	"github.com/joshk418/flexspec/internal/ui"
)

var listJSON bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all specs in the project",
	Long: `List all specs in the configured specs directory.

Reads specs_dir from .flexspec/config.yaml, then for each spec folder
(NNN-slug/README.md) prints the directory name, status, and task count
from YAML frontmatter (task_count, or computed from simple task rows/bullets / tasks/ files).
Use --json for full spec and task details.`,
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
		if listJSON {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(ui.EncodeSpecsForCLI(entries))
		}

		if len(entries) == 0 {
			if _, err := fmt.Fprintf(out, "No specs in %s\n", cfg.SpecsDir); err != nil {
				return err
			}

			return nil
		}

		rows := make([][]string, 0, len(entries))
		for _, e := range entries {
			rows = append(rows, []string{
				e.Dir,
				displayOrDash(e.Meta.Status),
				strconv.Itoa(e.TaskCount),
				e.Meta.SpecType,
			})
		}
		return clioutput.WriteTable(out,
			[]string{"IDENTIFIER", "STATUS", "TASKS", "TEMPLATE"},
			rows,
		)
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
	listCmd.Flags().BoolVar(&listJSON, "json", false, "Output specs as JSON")
}
