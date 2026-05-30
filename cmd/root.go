/*
Copyright © 2026 Josh Kyte
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "0.1.0" // x-release-please-version

var template string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flexspec",
	Short: "Spec-driven development CLI for generating and tracking feature specs",
	Long: `Flexspec is a spec-driven development CLI that helps teams generate,
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&template, "template", "t", "", "Template to use for spec generation")
}
