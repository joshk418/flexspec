package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/clioutput"
	"github.com/joshk418/flexspec/internal/config"
)

var configJSON bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or update project configuration from .flexspec/config.yaml",
	Long: `Print or update FlexSpec project settings from .flexspec/config.yaml.

Default output is a compact KEY / VALUE table. Use --json for machine-readable
output (scripts and agents). Use "config set" to update a single key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		cfg, err := config.Load(root)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if configJSON {
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			return enc.Encode(config.JSONDocumentFromConfig(cfg))
		}

		return printConfigTable(out, root, cfg)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Update a project configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("resolve working directory: %w", err)
		}

		cfg, err := config.Load(root)
		if err != nil {
			return err
		}

		updated, err := config.ApplyUpdate(cfg, args[0], args[1])
		if err != nil {
			return err
		}
		if err := config.Save(root, updated); err != nil {
			return err
		}

		return printConfigTable(cmd.OutOrStdout(), root, updated)
	},
}

func printConfigTable(out io.Writer, root string, cfg config.Config) error {
	entries := config.DisplayEntries(cfg)
	rows := make([][]string, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, []string{e.Key, e.Value})
	}
	if err := clioutput.WriteTable(out, []string{"KEY", "VALUE"}, rows); err != nil {
		return err
	}

	_, err := fmt.Fprintf(out, "config: %s\n", config.ConfigPath(root))
	return err
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.Flags().BoolVar(&configJSON, "json", false, "Output config as JSON")
}
