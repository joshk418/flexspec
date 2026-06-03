package migrate

import (
	"fmt"
	"io"
	"sort"

	"github.com/joshk418/flexspec/internal/clioutput"
	"github.com/joshk418/flexspec/internal/config"
)

// ChangeKind classifies a migration finding or applied change.
type ChangeKind string

const (
	KindRewrite ChangeKind = "rewrite"
	KindCreate  ChangeKind = "create"
	KindDelete  ChangeKind = "delete"
	KindReport  ChangeKind = "report"
)

// Change is one pending or applied migration item.
type Change struct {
	Migration string
	Path      string
	Kind      ChangeKind
	Detail    string
}

// Migration detects and applies one class of project upgrades.
type Migration interface {
	ID() string
	Description() string
	Detect(root string, cfg config.Config) ([]Change, error)
	Apply(root string, cfg config.Config) ([]Change, error)
}

// Plan runs Detect on each migration without writing.
func Plan(root string, cfg config.Config, migs []Migration) ([]Change, error) {
	var out []Change
	for _, m := range migs {
		changes, err := m.Detect(root, cfg)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", m.ID(), err)
		}
		out = append(out, changes...)
	}
	return sortChanges(out), nil
}

// Apply runs Apply on each migration in order.
func Apply(root string, cfg config.Config, migs []Migration) ([]Change, error) {
	var out []Change
	for _, m := range migs {
		changes, err := m.Apply(root, cfg)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", m.ID(), err)
		}
		out = append(out, changes...)
	}
	return sortChanges(out), nil
}

// Select returns migrations whose ID is in ids, preserving registry order.
func Select(migs []Migration, ids []string) ([]Migration, error) {
	if len(ids) == 0 {
		return migs, nil
	}
	want := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		want[id] = struct{}{}
	}
	var out []Migration
	for _, m := range migs {
		if _, ok := want[m.ID()]; ok {
			out = append(out, m)
			delete(want, m.ID())
		}
	}
	if len(want) > 0 {
		unknown := make([]string, 0, len(want))
		for id := range want {
			unknown = append(unknown, id)
		}
		sort.Strings(unknown)
		return nil, fmt.Errorf("unknown migration %q", unknown[0])
	}
	return out, nil
}

// HasApplicableChanges reports whether any change would write on apply.
func HasApplicableChanges(changes []Change) bool {
	for _, c := range changes {
		if c.Kind != KindReport {
			return true
		}
	}
	return false
}

// WriteChanges prints migration changes as an aligned table plus summary.
func WriteChanges(w io.Writer, changes []Change) error {
	if len(changes) == 0 {
		_, err := fmt.Fprintln(w, "0 pending change(s)")
		return err
	}

	pending := 0
	rows := make([][]string, len(changes))
	for i, c := range changes {
		if c.Kind != KindReport {
			pending++
		}
		rows[i] = []string{c.Migration, c.Path, string(c.Kind), c.Detail}
	}
	if err := clioutput.WriteTable(w,
		[]string{"MIGRATION", "PATH", "KIND", "DETAIL"},
		rows,
	); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, "%d change(s), %d pending\n", len(changes), pending)
	return err
}

func sortChanges(changes []Change) []Change {
	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Migration != changes[j].Migration {
			return changes[i].Migration < changes[j].Migration
		}
		if changes[i].Path != changes[j].Path {
			return changes[i].Path < changes[j].Path
		}
		return changes[i].Detail < changes[j].Detail
	})
	return changes
}
