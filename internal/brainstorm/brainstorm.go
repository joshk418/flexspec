// Package brainstorm scaffolds pre-spec exploration docs under .flexspec/brainstorms/.
package brainstorm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joshk418/flexspec/internal/spec"
)

const (
	flexspecDir        = ".flexspec"
	templatesDir       = "templates"
	brainstormTemplate = "brainstorm.md"
	brainstormsDir     = "brainstorms"
	dirPerm            = 0o755
	filePerm           = 0o644
)

// Result is the outcome of scaffolding a new brainstorm doc.
type Result struct {
	Path string
}

// Create scaffolds .flexspec/brainstorms/<slug>.md from the brainstorm template
func Create(root, name string, force bool) (Result, error) {
	slug, err := spec.Slugify(name)
	if err != nil {
		return Result{}, err
	}

	templatePath := filepath.Join(root, flexspecDir, templatesDir, brainstormTemplate)
	if _, err := os.Stat(templatePath); err != nil {
		if os.IsNotExist(err) {
			return Result{}, fmt.Errorf("brainstorm template not found at %s; run `flexspec init` (new projects) or `flexspec update --migrate` (existing projects) first", templatePath)
		}
		return Result{}, fmt.Errorf("stat %s: %w", templatePath, err)
	}

	data, err := os.ReadFile(templatePath)
	if err != nil {
		return Result{}, fmt.Errorf("read template %s: %w", templatePath, err)
	}

	brainstormsPath := filepath.Join(root, flexspecDir, brainstormsDir)
	if err := os.MkdirAll(brainstormsPath, dirPerm); err != nil {
		return Result{}, fmt.Errorf("create brainstorms directory %s: %w", brainstormsPath, err)
	}

	targetPath := filepath.Join(brainstormsPath, slug+".md")
	if _, err := os.Stat(targetPath); err == nil {
		if !force {
			return Result{}, fmt.Errorf("brainstorm doc %s already exists; use --force to overwrite", targetPath)
		}
	} else if !os.IsNotExist(err) {
		return Result{}, fmt.Errorf("stat %s: %w", targetPath, err)
	}

	if err := os.WriteFile(targetPath, data, filePerm); err != nil {
		return Result{}, fmt.Errorf("write %s: %w", targetPath, err)
	}

	return Result{Path: targetPath}, nil
}
