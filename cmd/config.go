package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/config"
)

var configJSON bool

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show project configuration from .flexspec/config.yaml",
	Long: `Print FlexSpec project settings from .flexspec/config.yaml.

Default output is a compact KEY / VALUE table. Use --json for machine-readable
output (scripts and agents).`,
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

		w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintln(w, "KEY\tVALUE"); err != nil {
			return err
		}
		for _, e := range config.DisplayEntries(cfg) {
			if _, err := fmt.Fprintf(w, "%s\t%s\n", e.Key, e.Value); err != nil {
				return err
			}
		}
		if err := w.Flush(); err != nil {
			return err
		}

		_, err = fmt.Fprintf(out, "config: %s\n", config.ConfigPath(root))
		return err
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVar(&configJSON, "json", false, "Output config as JSON")
}
