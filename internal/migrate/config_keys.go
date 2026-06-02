package migrate

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/joshk418/flexspec/internal/config"
)

var knownConfigKeys = map[string]struct{}{
	"specs_dir":       {},
	"always_one_shot": {},
	"spec_template":   {},
}

type configKeysMigration struct{}

func (m *configKeysMigration) ID() string { return "config-keys" }

func (m *configKeysMigration) Description() string {
	return "Reconcile .flexspec/config.yaml keys"
}

func (m *configKeysMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, false)
}

func (m *configKeysMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, true)
}

func (m *configKeysMigration) scan(root string, cfg config.Config, apply bool) ([]Change, error) {
	path := config.ConfigPath(root)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	rel := filepath.ToSlash(filepath.Join(".flexspec", "config.yaml"))

	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return nil, err
	}
	if rootNode.Kind != yaml.DocumentNode || len(rootNode.Content) == 0 {
		return nil, nil
	}
	mapNode := rootNode.Content[0]
	if mapNode.Kind != yaml.MappingNode {
		return nil, nil
	}

	var changes []Change
	present := map[string]bool{}
	for i := 0; i < len(mapNode.Content); i += 2 {
		key := mapNode.Content[i].Value
		present[key] = true
		if _, ok := knownConfigKeys[key]; !ok {
			changes = append(changes, Change{
				Migration: m.ID(),
				Path:      rel,
				Kind:      KindReport,
				Detail:    "unknown config key: " + key,
			})
		}
	}
	if !present["spec_template"] {
		changes = append(changes, Change{
			Migration: m.ID(),
			Path:      rel,
			Kind:      KindRewrite,
			Detail:    "add spec_template (default empty)",
		})
		if apply {
			updated, err := config.ApplyUpdate(cfg, "spec_template", "")
			if err != nil {
				return nil, err
			}
			if err := config.Save(root, updated); err != nil {
				return nil, err
			}
		}
	}
	return changes, nil
}
