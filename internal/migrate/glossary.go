package migrate

import (
	"os"
	"path/filepath"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/glossary"
)

type glossaryMigration struct{}

func (m *glossaryMigration) ID() string { return "glossary" }

func (m *glossaryMigration) Description() string {
	return "Create .flexspec/glossary.yaml when missing"
}

func (m *glossaryMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	return m.run(root, false)
}

func (m *glossaryMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	return m.run(root, true)
}

func (m *glossaryMigration) run(root string, apply bool) ([]Change, error) {
	path := filepath.Join(root, ".flexspec", "glossary.yaml")
	rel, _ := filepath.Rel(root, path)
	rel = filepath.ToSlash(rel)

	if _, err := os.Stat(path); err == nil {
		return nil, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if apply {
		if err := glossary.Save(root, glossary.DefaultDocument()); err != nil {
			return nil, err
		}
	}
	return []Change{{
		Migration: m.ID(),
		Path:      rel,
		Kind:      KindCreate,
		Detail:    "missing glossary file",
	}}, nil
}
