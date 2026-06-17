package spec

import (
	"slices"
	"strings"
)

// ValidTypes are the accepted spec type values.
var ValidTypes = []string{
	"feature",
	"bug",
	"chore",
	"refactor",
	"docs",
	"infra",
	"spike",
	"research",
}

// DefaultType is used when no type is supplied or inferred.
const DefaultType = "feature"

// IsValidType reports whether v is a recognized spec type.
func IsValidType(v string) bool {
	return slices.Contains(ValidTypes, strings.ToLower(strings.TrimSpace(v)))
}

// NormalizeType lowercases and trims a raw type value.
func NormalizeType(v string) string {
	return strings.ToLower(strings.TrimSpace(v))
}

// InferType derives a spec type from a free-text request, falling back to DefaultType.
func InferType(request string) string {
	s := " " + strings.ToLower(request) + " "
	switch {
	case containsAny(s, "broken", "bug", "wrong result", "regression", "repro", "reproduce", "crash", "failing"):
		return "bug"
	case containsAny(s, "spike", "investigate", "explore", "evaluate", "feasibility"):
		return "spike"
	case containsAny(s, "research", "study", "analyze", "analysis", "survey"):
		return "research"
	case containsAny(s, "rename", "bump", "update dependency", "config change", "rotate key", "chore"):
		return "chore"
	case containsAny(s, "refactor", "restructure", "extract", "simplify", "clean up"):
		return "refactor"
	case containsAny(s, "document", "readme", "guide", "docs", "documentation"):
		return "docs"
	case containsAny(s, "deploy", "migration", "infrastructure", "pipeline", "provision", "terraform", "k8s"):
		return "infra"
	default:
		return DefaultType
	}
}

func containsAny(haystack string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(haystack, n) {
			return true
		}
	}
	return false
}
