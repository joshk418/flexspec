package migrate

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joshk418/flexspec/internal/config"
)

var placeholderRE = regexp.MustCompile(`\{[^{}]+\}`)

type charterCheckMigration struct {
	templates fs.FS
}

func (m *charterCheckMigration) ID() string { return "charter-check" }

func (m *charterCheckMigration) Description() string {
	return "Report missing charter sections and template placeholders"
}

func (m *charterCheckMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	return m.check(root)
}

func (m *charterCheckMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	// Report-only: Apply returns the same findings without writing.
	return m.check(root)
}

func (m *charterCheckMigration) check(root string) ([]Change, error) {
	charterPath := filepath.Join(root, ".flexspec", "charter.md")
	rel := filepath.ToSlash(charterPath)
	data, err := os.ReadFile(charterPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Change{{
				Migration: m.ID(),
				Path:      rel,
				Kind:      KindReport,
				Detail:    "charter file missing",
			}}, nil
		}
		return nil, err
	}
	content := string(data)
	var changes []Change

	expected, err := m.expectedHeadings()
	if err != nil {
		return nil, err
	}
	for _, heading := range expected {
		if !strings.Contains(content, heading) {
			changes = append(changes, Change{
				Migration: m.ID(),
				Path:      rel,
				Kind:      KindReport,
				Detail:    "missing section: " + heading,
			})
		}
	}
	if strings.Contains(content, "<!--") {
		changes = append(changes, Change{
			Migration: m.ID(),
			Path:      rel,
			Kind:      KindReport,
			Detail:    "contains HTML guidance comments",
		})
	}
	if placeholderRE.MatchString(content) {
		changes = append(changes, Change{
			Migration: m.ID(),
			Path:      rel,
			Kind:      KindReport,
			Detail:    "contains template placeholders",
		})
	}
	return changes, nil
}

func (m *charterCheckMigration) expectedHeadings() ([]string, error) {
	if m.templates == nil {
		return nil, nil
	}
	data, err := fs.ReadFile(m.templates, "charter.md")
	if err != nil {
		return nil, err
	}
	var headings []string
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "## ") {
			headings = append(headings, line)
		}
	}
	return headings, sc.Err()
}
