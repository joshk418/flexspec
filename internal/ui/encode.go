package ui

import (
	"github.com/joshk418/flexspec/internal/spec"
)

// EncodeSpecsForCLI exports list entries as JSON API objects.
func EncodeSpecsForCLI(entries []spec.SpecEntry) []SpecJSON {
	return encodeSpecs(entries)
}

func encodeSpecs(entries []spec.SpecEntry) []SpecJSON {
	out := make([]SpecJSON, 0, len(entries))
	for _, e := range entries {
		out = append(out, encodeSpec(e))
	}
	return out
}

func encodeSpec(e spec.SpecEntry) SpecJSON {
	s := SpecJSON{
		ID:          e.ID,
		Dir:         e.Dir,
		Name:        e.Meta.Name,
		Description: e.Meta.Description,
		Status:      e.Meta.Status,
		SpecType:    e.Meta.SpecType,
	}
	for _, t := range e.Tasks {
		s.Tasks = append(s.Tasks, TaskJSON{
			ID:     t.Meta.ID,
			File:   t.File,
			Name:   t.Meta.Name,
			Status: t.Meta.Status,
		})
	}
	return s
}
