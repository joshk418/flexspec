package validate

import (
	"fmt"
	"io"
	"strings"
)

// WriteFindings prints findings and a summary to w.
func WriteFindings(w io.Writer, findings []Finding) error {
	errCount, warnCount := 0, 0
	for _, f := range findings {
		if f.Severity == SeverityError {
			errCount++
		} else {
			warnCount++
		}
		if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			f.Severity,
			f.Path,
			f.Rule,
			f.Message,
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
