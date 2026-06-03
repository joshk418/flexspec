package validate

import (
	"fmt"
	"io"
	"strings"

	"github.com/joshk418/flexspec/internal/clioutput"
)

// WriteFindings prints findings and a summary to w.
func WriteFindings(w io.Writer, findings []Finding) error {
	errCount, warnCount := 0, 0
	if len(findings) > 0 {
		rows := make([][]string, len(findings))
		for i, f := range findings {
			if f.Severity == SeverityError {
				errCount++
			} else {
				warnCount++
			}
			rows[i] = []string{string(f.Severity), f.Path, f.Rule, f.Message}
		}
		if err := clioutput.WriteTable(w,
			[]string{"SEVERITY", "PATH", "RULE", "MESSAGE"},
			rows,
		); err != nil {
			return err
		}
	}

	var parts []string
	if errCount > 0 {
		parts = append(parts, fmt.Sprintf("%d error(s)", errCount))
	}
	if warnCount > 0 {
		parts = append(parts, fmt.Sprintf("%d warning(s)", warnCount))
	}
	if len(parts) == 0 {
		if _, err := fmt.Fprintln(w, "0 error(s), 0 warning(s)"); err != nil {
			return err
		}
		return nil
	}
	if _, err := fmt.Fprintln(w, strings.Join(parts, ", ")); err != nil {
		return err
	}
	return nil
}
