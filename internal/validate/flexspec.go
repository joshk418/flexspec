package validate

import (
	"os"
	"path/filepath"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

var requiredTemplates = []string{
	"README.md",
	"flexspec-simple.md",
	filepath.Join("expanded", "flexspec-expanded.md"),
	filepath.Join("expanded", "flexspec-expanded-task.md"),
}

// CheckFlexspec validates charter and template files under .flexspec/.
func CheckFlexspec(root string, _ config.Config, _ Options) []Finding {
	var findings []Finding

	charterRel := filepath.Join(flexspecDir, "charter.md")
	charterPath := filepath.Join(root, charterRel)
	if _, err := os.Stat(charterPath); err != nil {
		if os.IsNotExist(err) {
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     charterRel,
				Rule:     "charter.missing",
				Message:  "not found; run `flexspec init` first",
			})
		} else {
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     charterRel,
				Rule:     "charter.stat",
				Message:  err.Error(),
			})
		}
	} else if _, err := spec.ParseSpecMeta(charterPath); err != nil {
		findings = append(findings, Finding{
			Severity: SeverityError,
			Path:     charterRel,
			Rule:     "charter.frontmatter",
			Message:  err.Error(),
		})
	}

	templatesBase := filepath.Join(root, flexspecDir, "templates")
	for _, rel := range requiredTemplates {
		relPath := filepath.Join(flexspecDir, "templates", rel)
		path := filepath.Join(templatesBase, rel)
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				findings = append(findings, Finding{
					Severity: SeverityError,
					Path:     relPath,
					Rule:     "templates.missing",
					Message:  "required template file not found",
				})
				continue
			}
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     relPath,
				Rule:     "templates.stat",
				Message:  err.Error(),
			})
		}
	}

	return findings
}
