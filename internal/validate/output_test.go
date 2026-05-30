package validate

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteFindings_summary(t *testing.T) {
	var buf bytes.Buffer
	findings := []Finding{
		{Severity: SeverityError, Path: "a", Rule: "r1", Message: "m1"},
		{Severity: SeverityWarning, Path: "b", Rule: "r2", Message: "m2"},
	}
	if err := WriteFindings(&buf, findings); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "1 error(s)") || !strings.Contains(out, "1 warning(s)") {
		t.Fatalf("output = %q", out)
	}
}
