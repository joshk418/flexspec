package migrate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

type statusRenameMigration struct{}

func (m *statusRenameMigration) ID() string { return "status-rename" }

func (m *statusRenameMigration) Description() string {
	return "Rename legacy spec/task statuses (refined→planned, initial→draft)"
}

func (m *statusRenameMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, false)
}

func (m *statusRenameMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, true)
}

func (m *statusRenameMigration) scan(root string, cfg config.Config, apply bool) ([]Change, error) {
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
		specDir := filepath.Join(specsPath, e.Name())
		readme := filepath.Join(specDir, "README.md")
		if err := m.checkFile(root, readme, apply, &changes); err != nil {
			return nil, err
		}
		tasksDir := filepath.Join(specDir, "tasks")
		taskEntries, err := os.ReadDir(tasksDir)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, te := range taskEntries {
			if te.IsDir() {
				continue
			}
			name := te.Name()
			if !strings.HasPrefix(name, "T-") || !strings.HasSuffix(name, ".md") {
				continue
			}
			if err := m.checkFile(root, filepath.Join(tasksDir, name), apply, &changes); err != nil {
				return nil, err
			}
		}
	}
	return changes, nil
}

func (m *statusRenameMigration) checkFile(root, path string, apply bool, changes *[]Change) error {
	var status string
	if strings.Contains(filepath.ToSlash(path), "/tasks/") {
		meta, err := spec.ParseTaskMeta(path)
		if err != nil {
			return nil
		}
		status = meta.Status
	} else {
		meta, err := spec.ParseSpecMeta(path)
		if err != nil {
			return nil
		}
		status = meta.Status
	}
	return m.applyStatus(root, path, status, apply, changes)
}

func (m *statusRenameMigration) applyStatus(root, path, raw string, apply bool, changes *[]Change) error {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	normalized := spec.NormalizeSpecStatus(raw)
	if normalized == strings.ToLower(raw) {
		return nil
	}
	rel, _ := filepath.Rel(root, path)
	if rel == "" {
		rel = path
	}
	rel = filepath.ToSlash(rel)
	*changes = append(*changes, Change{
		Migration: m.ID(),
		Path:      rel,
		Kind:      KindRewrite,
		Detail:    raw + " → " + normalized,
	})
	if apply {
		if err := spec.SetFileStatus(path, normalized); err != nil {
			return err
		}
	}
	return nil
}
