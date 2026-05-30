/*
Copyright © 2026 Josh Kyte
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

var statusSetStatus string
var statusSetTask string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Read or update spec and task status in frontmatter",
}

var statusSetCmd = &cobra.Command{
	Use:   "set <spec>",
	Short: "Set status on a spec or task README",
	Long: `Set the status field in YAML frontmatter for a spec README.md or task file.

<spec> is the spec directory name (e.g. 002-management-ui) or numeric id (002).`,
	Args: cobra.ExactArgs(1),
	RunE: runStatusSet,
}

func runStatusSet(_ *cobra.Command, args []string) error {
	if strings.TrimSpace(statusSetStatus) == "" {
		return fmt.Errorf("--status is required")
	}
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve working directory: %w", err)
	}
	cfg, err := config.Load(root)
	if err != nil {
		return err
	}
	target := args[0]
	dir, err := resolveSpecDir(root, cfg, target)
	if err != nil {
		return err
	}
	if statusSetTask != "" {
		taskFile := statusSetTask
		if !strings.HasSuffix(taskFile, ".md") {
			taskFile += ".md"
		}
		path := filepath.Join(root, cfg.SpecsDir, dir, "tasks", taskFile)
		if err := spec.SetFileStatus(path, statusSetStatus); err != nil {
			return err
		}
		fmt.Printf("Updated %s status to %s\n", taskFile, statusSetStatus)
		return nil
	}
	path := filepath.Join(root, cfg.SpecsDir, dir, "README.md")
	if err := spec.SetFileStatus(path, statusSetStatus); err != nil {
		return err
	}
	fmt.Printf("Updated spec %s status to %s\n", dir, statusSetStatus)
	return nil
}

func resolveSpecDir(root string, cfg config.Config, target string) (string, error) {
	entries, err := spec.List(root, cfg)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if e.Dir == target || e.ID == target {
			return e.Dir, nil
		}
	}
	return "", fmt.Errorf("spec %q not found", target)
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.AddCommand(statusSetCmd)
	statusSetCmd.Flags().StringVar(&statusSetStatus, "status", "", "New status value")
	statusSetCmd.Flags().StringVar(&statusSetTask, "task", "", "Task filename (e.g. T-001-slug.md)")
	_ = statusSetCmd.MarkFlagRequired("status")
}
