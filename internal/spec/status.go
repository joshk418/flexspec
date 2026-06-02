package spec

import "strings"

// UnassignedColumn is the board column for specs whose status is empty or not
// a known lifecycle status.
const UnassignedColumn = "unassigned"

// specStatuses is the canonical, ordered list of spec lifecycle statuses.
// Keep in sync with ui/src/api/status.ts (SPEC_COLUMNS).
var specStatuses = []string{
	"draft",
	"planned",
	"in_progress",
	"in_review",
	"complete",
}

// legacyStatusAliases maps removed/renamed statuses to their current value.
var legacyStatusAliases = map[string]string{
	"refined": "planned",
	"initial": "draft",
}

// SpecStatuses returns the canonical spec lifecycle statuses in board order.
func SpecStatuses() []string {
	out := make([]string, len(specStatuses))
	copy(out, specStatuses)
	return out
}

// NormalizeSpecStatus trims and lowercases a raw status and maps legacy values
// (refined → planned, initial → draft) to the current vocabulary. Unknown
// values pass through unchanged so callers can route them to Unassigned.
func NormalizeSpecStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	if mapped, ok := legacyStatusAliases[s]; ok {
		return mapped
	}
	return s
}

// ColumnForSpecStatus returns the board column for a raw status. It normalizes
// legacy values first; anything not in the canonical list maps to Unassigned.
func ColumnForSpecStatus(status string) string {
	s := NormalizeSpecStatus(status)
	for _, known := range specStatuses {
		if s == known {
			return s
		}
	}
	return UnassignedColumn
}
