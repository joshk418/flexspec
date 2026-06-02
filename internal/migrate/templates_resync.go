package migrate

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/joshk418/flexspec/internal/config"
)

type templatesResyncMigration struct {
	templates fs.FS
	force     bool
}

func (m *templatesResyncMigration) ID() string { return "templates-resync" }

func (m *templatesResyncMigration) Description() string {
	return "Re-sync .flexspec/templates from embedded templates"
}

func (m *templatesResyncMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	return m.walk(root, false)
}

func (m *templatesResyncMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	return m.walk(root, true)
}

func (m *templatesResyncMigration) walk(root string, apply bool) ([]Change, error) {
	if m.templates == nil {
		return nil, nil
	}
	destRoot := filepath.Join(root, ".flexspec", "templates")
	var changes []Change

	err := fs.WalkDir(m.templates, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Base(p) == "charter.md" {
			return nil
		}
		rel := filepath.ToSlash(p)
		if rel == "." {
			return nil
		}
		embedded, err := fs.ReadFile(m.templates, p)
		if err != nil {
			return err
		}
		target := filepath.Join(destRoot, filepath.FromSlash(rel))
		relOut, _ := filepath.Rel(root, target)
		relOut = filepath.ToSlash(relOut)

		if _, err := os.Stat(target); os.IsNotExist(err) {
			changes = append(changes, Change{
				Migration: m.ID(),
				Path:      relOut,
				Kind:      KindCreate,
				Detail:    "missing template file",
			})
			if apply {
				if err := writeFile(target, embedded); err != nil {
					return err
				}
			}
			return nil
		}
		existing, err := os.ReadFile(target)
		if err != nil {
			return err
		}
		if bytes.Equal(existing, embedded) {
			return nil
		}
		kind := KindReport
		detail := "differs from embedded template (use --force to overwrite)"
		if m.force {
			kind = KindRewrite
			detail = "overwrite with embedded template"
		}
		changes = append(changes, Change{
			Migration: m.ID(),
			Path:      relOut,
			Kind:      kind,
			Detail:    detail,
		})
		if apply && m.force {
			if err := writeFile(target, embedded); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
