package validate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joshk418/flexspec/internal/config"
	"gopkg.in/yaml.v3"
)

const (
	flexspecDir = ".flexspec"
	configFile  = "config.yaml"
)

// LoadConfig loads project config or returns findings when load fails.
func LoadConfig(root string) (config.Config, []Finding, bool) {
	path := filepath.Join(root, flexspecDir, configFile)
	relPath := filepath.Join(flexspecDir, configFile)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config.Config{}, []Finding{{
				Severity: SeverityError,
				Path:     relPath,
				Rule:     "config.missing",
				Message:  "not found; run `flexspec init` first",
			}}, false
		}
		return config.Config{}, []Finding{{
			Severity: SeverityError,
			Path:     relPath,
			Rule:     "config.read",
			Message:  err.Error(),
		}}, false
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return config.Config{}, []Finding{{
			Severity: SeverityError,
			Path:     relPath,
			Rule:     "config.parse",
			Message:  err.Error(),
		}}, false
	}

	var findings []Finding
	if cfg.SpecsDir == "" {
		findings = append(findings, Finding{
			Severity: SeverityError,
			Path:     relPath,
			Rule:     "config.specs_dir",
			Message:  "specs_dir must be set",
		})
	}

	st := strings.TrimSpace(cfg.SpecTemplate)
	if st != "" && !strings.EqualFold(st, "simple") && !strings.EqualFold(st, "expanded") {
		findings = append(findings, Finding{
			Severity: SeverityError,
			Path:     relPath,
			Rule:     "config.spec_template",
			Message:  "spec_template must be simple or expanded",
		})
	}

	if len(findings) > 0 {
		return cfg, findings, false
	}
	return cfg, nil, true
}

// CheckConfig validates .flexspec/config.yaml.
func CheckConfig(root string, _ config.Config, _ Options) []Finding {
	_, findings, _ := LoadConfig(root)
	return findings
}
