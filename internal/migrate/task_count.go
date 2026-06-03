package migrate

import (
	"os"
	"path/filepath"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

type taskCountMigration struct{}

func (m *taskCountMigration) ID() string { return "task-count" }

func (m *taskCountMigration) Description() string {
	return "Backfill spec task_count frontmatter and README metadata header"
}

func (m *taskCountMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, false)
}

func (m *taskCountMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, true)
}

func (m *taskCountMigration) scan(root string, cfg config.Config, apply bool) ([]Change, error) {
	specsPath := filepath.Join(root, cfg.SpecsDir)
	entries, err := os.ReadDir(specsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var changes []Change
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		readme := filepath.Join(specsPath, e.Name(), "README.md")
		if _, err := os.Stat(readme); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		needs, err := spec.NeedsTaskCountSync(readme)
		if err != nil {
			return nil, err
		}
		if !needs {
			continue
		}
		rel, _ := filepath.Rel(root, readme)
		if apply {
			if err := spec.SyncTaskCount(readme); err != nil {
				return nil, err
			}
		}
		changes = append(changes, Change{
			Migration: m.ID(),
			Path:      rel,
			Kind:      KindRewrite,
			Detail:    "sync task_count frontmatter and metadata header",
		})
	}
	return changes, nil
}
