package spec

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func TestNormalizeSpecStatus(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"legacy refined", "refined", "planned"},
		{"legacy initial", "initial", "draft"},
		{"trim and lowercase", "  Planned ", "planned"},
		{"known passthrough", "in_review", "in_review"},
		{"unknown passthrough", "archived", "archived"},
		{"empty", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NormalizeSpecStatus(tc.in); got != tc.want {
				t.Errorf("NormalizeSpecStatus(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestColumnForSpecStatus(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"legacy refined to planned", "refined", "planned"},
		{"legacy initial to draft", "initial", "draft"},
		{"draft", "draft", "draft"},
		{"complete", "complete", "complete"},
		{"unknown to unassigned", "archived", UnassignedColumn},
		{"empty to unassigned", "", UnassignedColumn},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ColumnForSpecStatus(tc.in); got != tc.want {
				t.Errorf("ColumnForSpecStatus(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestSpecStatuses(t *testing.T) {
	want := []string{"draft", "planned", "in_progress", "in_review", "complete"}
	if got := SpecStatuses(); !reflect.DeepEqual(got, want) {
		t.Errorf("SpecStatuses() = %v, want %v", got, want)
	}
	// Mutating the returned slice must not affect the canonical list.
	SpecStatuses()[0] = "mutated"
	if got := SpecStatuses(); got[0] != "draft" {
		t.Errorf("SpecStatuses() not defensively copied: got[0] = %q", got[0])
	}
}

// TestSpecStatusesMatchTypeScript enforces FR-005: the Go canonical list and the
// TypeScript SPEC_COLUMNS mirror must stay in the same order. Fails if they diverge.
func TestSpecStatusesMatchTypeScript(t *testing.T) {
	path := filepath.Join("..", "..", "ui", "src", "api", "status.ts")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	block := regexp.MustCompile(`(?s)SPEC_COLUMNS\s*=\s*\[(.*?)\]`).FindSubmatch(data)
	if block == nil {
		t.Fatalf("SPEC_COLUMNS array not found in %s", path)
	}
	items := regexp.MustCompile(`"([^"]+)"`).FindAllStringSubmatch(string(block[1]), -1)
	var tsStatuses []string
	for _, m := range items {
		tsStatuses = append(tsStatuses, m[1])
	}
	if want := SpecStatuses(); !reflect.DeepEqual(tsStatuses, want) {
		t.Errorf("status.ts SPEC_COLUMNS = %v, want %v (Go and TS lists diverged)", tsStatuses, want)
	}
}
