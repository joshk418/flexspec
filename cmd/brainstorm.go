package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/brainstorm"
)

var brainstormForce bool

var brainstormCmd = &cobra.Command{
	Use:   "brainstorm",
	Short: "Manage pre-spec brainstorm docs",
}

var brainstormNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new brainstorm doc",
	Long: `Create a new pre-spec exploration doc at .flexspec/brainstorms/<slug>.md
from the project-local brainstorm template.

Brainstorm docs have no status lifecycle and are not returned by
` + "`flexspec list`" + ` or shown in the management UI board.`,
	Args: cobra.MinimumNArgs(1),
	RunE: runBrainstormNew,
}

func runBrainstormNew(cmd *cobra.Command, args []string) error {
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working directory: %w", err)
	}

	name := strings.Join(args, " ")
	result, err := brainstorm.Create(root, name, brainstormForce)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), "Created brainstorm doc"); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "  path: %s\n", result.Path); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(brainstormCmd)
	brainstormCmd.AddCommand(brainstormNewCmd)

	brainstormNewCmd.Flags().BoolVar(&brainstormForce, "force", false, "Overwrite an existing brainstorm doc")
}
