package clioutput

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// WriteTable prints an aligned table with a header row and data rows.
func WriteTable(w io.Writer, headers []string, rows [][]string) error {
	if len(headers) == 0 {
		return nil
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, strings.Join(headers, "\t")); err != nil {
		return err
	}
	for _, row := range rows {
		if _, err := fmt.Fprintln(tw, strings.Join(row, "\t")); err != nil {
			return err
		}
	}
	return tw.Flush()
}
