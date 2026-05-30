/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new spec file",
	Long: `Create a new spec file in the user defined specs directory (defined in .flexspec/config.yaml).
	This will create the following an autoincrementing sequence number spec directory with the spec template:
	Simple:
	- <specs_dir>/001-<spec_name>/
	- <specs_dir>/001-<spec_name>/README.md
	Expanded:
	- <specs_dir>/002-<spec_name>/
	- <specs_dir>/002-<spec_name>/README.md
	- <specs_dir>/002-<spec_name>/tasks/
	- <specs_dir>/002-<spec_name>/tasks/README.md`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("new called")
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
