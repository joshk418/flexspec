package migrate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

type typeBackfillMigration struct{}

func (m *typeBackfillMigration) ID() string { return "type-backfill" }

func (m *typeBackfillMigration) Description() string {
	return "Backfill spec frontmatter `type` field for existing specs"
}

func (m *typeBackfillMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, false)
}

func (m *typeBackfillMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	return m.scan(root, cfg, true)
}

func (m *typeBackfillMigration) scan(root string, cfg config.Config, apply bool) ([]Change, error) {
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

		inferred, hasSections45, err := m.inferTypeFor(readme)
		if err != nil {
			return nil, err
		}
		if inferred == "" {
			continue
		}

		rel, _ := filepath.Rel(root, readme)
		rel = filepath.ToSlash(rel)
		detail := fmt.Sprintf("type: %s (inferred from %s)", inferred, inferenceSource(hasSections45))
		changes = append(changes, Change{
			Migration: m.ID(),
			Path:      rel,
			Kind:      KindRewrite,
			Detail:    detail,
		})
		if apply {
			if err := m.writeType(readme, inferred); err != nil {
				return nil, err
			}
		}
	}
	return changes, nil
}

// inferTypeFor returns the inferred type when frontmatter `type` is missing; "" when already set.
func (m *typeBackfillMigration) inferTypeFor(readme string) (inferred string, hasSections45 bool, err error) {
	meta, err := spec.ParseSpecMeta(readme)
	if err != nil {
		return "", false, nil
	}
	if spec.NormalizeType(meta.Type) != "" {
		return "", false, nil
	}

	parts, perr := spec.ReadFileParts(readme)
	if perr != nil {
		return "", false, perr
	}
	hasSections45 = hasBugSections(parts.Body)
	if hasSections45 {
		return "bug", true, nil
	}
	return spec.DefaultType, false, nil
}

func inferenceSource(hasSections45 bool) string {
	if hasSections45 {
		return "filled Sections 4/5"
	}
	return "default"
}

func yamlUnmarshal(data []byte, out any) error {
	return yaml.Unmarshal(data, out)
}

func yamlMarshal(in any) ([]byte, error) {
	return yaml.Marshal(in)
}

// hasBugSections reports whether Section 4/5 has bug-headed content (not "Not applicable").
func hasBugSections(body string) bool {
	lines := strings.Split(body, "\n")
	var capturing []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## 4.") || strings.HasPrefix(trimmed, "## 5.") {
			headingLower := strings.ToLower(trimmed)
			if strings.Contains(headingLower, "expected result") ||
				strings.Contains(headingLower, "actual result") ||
				strings.Contains(headingLower, "bug result") {
				capturing = append(capturing, "")
				continue
			}
			capturing = nil
			continue
		}
		if capturing != nil && strings.HasPrefix(trimmed, "## ") {
			capturing = nil
			continue
		}
		if capturing != nil {
			capturing = append(capturing, line)
		}
	}
	if len(capturing) == 0 {
		return false
	}
	joined := strings.ToLower(strings.Join(capturing, "\n"))
	if strings.Contains(joined, "not applicable") {
		return false
	}
	return strings.TrimSpace(strings.Join(capturing, "\n")) != ""
}

// writeType inserts or updates the `type` frontmatter field via yaml round-trip.
func (m *typeBackfillMigration) writeType(readme, typeValue string) error {
	parts, err := spec.ReadFileParts(readme)
	if err != nil {
		return err
	}
	var fields map[string]any
	if err := yamlUnmarshal([]byte(parts.Frontmatter), &fields); err != nil {
		return err
	}
	if fields == nil {
		fields = map[string]any{}
	}
	fields["type"] = typeValue
	updated, err := yamlMarshal(fields)
	if err != nil {
		return err
	}
	full := "---\n" + string(updated) + "---\n" + parts.Body
	return os.WriteFile(readme, []byte(full), 0o644)
}
