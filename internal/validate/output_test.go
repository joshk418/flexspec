package validate

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteFindings(t *testing.T) {
	tests := []struct {
		name     string
		findings []Finding
		want     []string
		absent   []string
	}{
		{
			name:     "clean",
			findings: nil,
			want:     []string{"0 error(s), 0 warning(s)"},
			absent:   []string{"SEVERITY"},
		},
		{
			name: "errors and warnings",
			findings: []Finding{
				{Severity: SeverityError, Path: "a", Rule: "test.err", Message: "err"},
				{Severity: SeverityWarning, Path: "b", Rule: "test.warn", Message: "warn"},
			},
			want: []string{"SEVERITY", "PATH", "RULE", "MESSAGE", "1 error(s)", "1 warning(s)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := WriteFindings(&buf, tt.findings); err != nil {
				t.Fatal(err)
			}
			out := buf.String()
			for _, want := range tt.want {
				if !strings.Contains(out, want) {
					t.Errorf("output missing %q\n%s", want, out)
				}
			}
			for _, absent := range tt.absent {
				if strings.Contains(out, absent) {
					t.Errorf("output should not contain %q\n%s", absent, out)
				}
			}
		})
	}
}
