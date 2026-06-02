package migrate

import (
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
