package migrate

import (
	"bytes"
	"strings"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

type fakeMigration struct {
	id      string
	detect  []Change
	apply   []Change
	detectN int
	applyN  int
}

func (f *fakeMigration) ID() string          { return f.id }
func (f *fakeMigration) Description() string { return "fake" }

func (f *fakeMigration) Detect(root string, cfg config.Config) ([]Change, error) {
	f.detectN++
	return f.detect, nil
}

func (f *fakeMigration) Apply(root string, cfg config.Config) ([]Change, error) {
	f.applyN++
	return f.apply, nil
}

func TestPlan_noWrites(t *testing.T) {
	m := &fakeMigration{
		id: "fake",
		detect: []Change{{
			Migration: "fake",
			Path:      "x",
			Kind:      KindRewrite,
			Detail:    "test",
		}},
	}
	changes, err := Plan(".", config.Config{SpecsDir: "specs"}, []Migration{m})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 {
		t.Fatalf("got %d changes", len(changes))
	}
	if m.applyN != 0 {
		t.Fatal("Apply should not run during Plan")
	}
}

func TestSelect_unknownID(t *testing.T) {
	migs := []Migration{&fakeMigration{id: "a"}}
	_, err := Select(migs, []string{"missing"})
	if err == nil {
		t.Fatal("expected error for unknown migration id")
	}
}

func TestWriteChanges(t *testing.T) {
	tests := []struct {
		name   string
		input  []Change
		want   []string
		absent []string
	}{
		{
			name:   "empty",
			input:  nil,
			want:   []string{"0 pending change(s)"},
			absent: []string{"MIGRATION"},
		},
		{
			name: "with changes",
			input: []Change{{
				Migration: "status-rename",
				Path:      "specs/001/README.md",
				Kind:      KindRewrite,
				Detail:    "planned",
			}},
			want: []string{"MIGRATION", "PATH", "KIND", "DETAIL", "1 change(s), 1 pending"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := WriteChanges(&buf, tt.input); err != nil {
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
