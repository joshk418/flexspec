package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.3.6" // x-release-please-version

var template string

// SkillsFS holds the embedded skills tree mounted by main before Execute runs.
var SkillsFS fs.FS

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flexspec",
	Short: "Spec-driven development CLI for generating and tracking feature specs",
	Long: `FlexSpec is a spec-driven development CLI that helps teams generate,
organize, and keep track of feature specifications.

Write specs as simple markdown files from built-in templates, or connect
adapters for Jira, Shortcut, GitHub Issues, and other issue trackers.`,
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		if template != "" {
			fmt.Println("template: not yet available")
			return nil
		}
		return cmd.Help()
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&template, "template", "t", "", "Template to use for spec generation")
}
