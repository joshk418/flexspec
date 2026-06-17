package spec

import "testing"

func TestIsValidType(t *testing.T) {
	valid := []string{"feature", "bug", "chore", "refactor", "docs", "infra", "spike", "research"}
	for _, v := range valid {
		if !IsValidType(v) {
			t.Errorf("IsValidType(%q) = false, want true", v)
		}
		if !IsValidType(upper(v)) {
			t.Errorf("IsValidType should be case-insensitive for %q", v)
		}
	}
	if IsValidType("hotfix") {
		t.Error("IsValidType(hotfix) = true, want false")
	}
	if IsValidType("") {
		t.Error("IsValidType(\"\") = true, want false")
	}
}

func TestNormalizeType(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Feature", "feature"},
		{"  BUG  ", "bug"},
		{"Spike", "spike"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := NormalizeType(tt.in); got != tt.want {
			t.Errorf("NormalizeType(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestInferType(t *testing.T) {
	tests := []struct {
		name    string
		request string
		want    string
	}{
		{"bug signal", "fix the broken login flow", "bug"},
		{"regression", "regression after v1.2 deploy", "bug"},
		{"spike signal", "spike to evaluate option A vs B", "spike"},
		{"investigate", "investigate why latency increased", "spike"},
		{"research signal", "research auth patterns for the migration", "research"},
		{"chore signal", "bump go version to 1.27", "chore"},
		{"rename", "rename Foo to Bar across the repo", "chore"},
		{"refactor signal", "refactor the export module to use streams", "refactor"},
		{"simplify", "simplify the config loader", "refactor"},
		{"docs signal", "document the onboarding flow in the readme", "docs"},
		{"infra signal", "deploy the new worker via terraform", "infra"},
		{"pipeline", "fix the CI pipeline provisioning", "infra"},
		{"default feature", "add a notifications center", "feature"},
		{"empty", "", "feature"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InferType(tt.request); got != tt.want {
				t.Errorf("InferType(%q) = %q, want %q", tt.request, got, tt.want)
			}
		})
	}
}

func upper(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		out[i] = c
	}
	return string(out)
}
