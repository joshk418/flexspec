package validate

import (
	"sort"

	"github.com/joshk418/flexspec/internal/config"
)

// Severity classifies a validation finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Finding is one validation issue.
type Finding struct {
	Severity Severity
	Path     string
	Rule     string
	Message  string
}

// Options configures validation behavior.
type Options struct {
	Strict bool
}

// Check validates one aspect of a FlexSpec project.
type Check func(root string, cfg config.Config, opts Options) []Finding

// RunAll runs config, flexspec, and specs checks when config loads successfully.
func RunAll(root string, opts Options) []Finding {
	cfg, findings, ok := LoadConfig(root)
	if !ok {
		return sortFindings(findings)
	}
	findings = append(findings, CheckFlexspec(root, cfg, opts)...)
	findings = append(findings, CheckSpecs(root, cfg, opts)...)
	return sortFindings(findings)
}

// Run executes checks and returns combined findings.
func Run(root string, cfg config.Config, opts Options, checks ...Check) []Finding {
	var out []Finding
	for _, check := range checks {
		out = append(out, check(root, cfg, opts)...)
	}
	return sortFindings(out)
}

// HasErrors reports whether any finding has error severity.
func HasErrors(findings []Finding) bool {
	for _, f := range findings {
		if f.Severity == SeverityError {
			return true
		}
	}
	return false
}

func sortFindings(findings []Finding) []Finding {
	sort.Slice(findings, func(i, j int) bool {
		si, sj := severityRank(findings[i].Severity), severityRank(findings[j].Severity)
		if si != sj {
			return si < sj
		}
		if findings[i].Path != findings[j].Path {
			return findings[i].Path < findings[j].Path
		}
		if findings[i].Rule != findings[j].Rule {
			return findings[i].Rule < findings[j].Rule
		}
		return findings[i].Message < findings[j].Message
	})
	return findings
}

func severityRank(s Severity) int {
	if s == SeverityError {
		return 0
	}
	return 1
}
