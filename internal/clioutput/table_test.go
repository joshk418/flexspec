package clioutput

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteTable(t *testing.T) {
	tests := []struct {
		name    string
		headers []string
		rows    [][]string
		want    []string
	}{
		{
			name:    "headers and rows",
			headers: []string{"KEY", "VALUE"},
			rows: [][]string{
				{"specs_dir", "specs"},
				{"always_one_shot", "false"},
			},
			want: []string{"KEY", "VALUE", "specs_dir", "always_one_shot"},
		},
		{
			name:    "empty headers",
			headers: nil,
			rows:    [][]string{{"a", "b"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := WriteTable(&buf, tt.headers, tt.rows); err != nil {
				t.Fatal(err)
			}
			out := buf.String()
			for _, want := range tt.want {
				if !strings.Contains(out, want) {
					t.Errorf("output missing %q\n%s", want, out)
				}
			}
			if tt.headers != nil && len(tt.rows) > 1 {
				keyIdx := strings.Index(out, "KEY")
				specIdx := strings.Index(out, "specs_dir")
				alwaysIdx := strings.Index(out, "always_one_shot")
				if keyIdx < 0 || specIdx < 0 || alwaysIdx < 0 {
					t.Fatalf("alignment check failed\n%s", out)
				}
				if specIdx <= keyIdx || alwaysIdx <= specIdx {
					t.Fatalf("expected header before rows\n%s", out)
				}
			}
		})
	}
}
