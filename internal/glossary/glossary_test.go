package glossary

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_missingFileReturnsEmpty(t *testing.T) {
	tmp := t.TempDir()
	doc, err := Load(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(doc.Terms) != 0 {
		t.Fatalf("expected empty terms, got %d", len(doc.Terms))
	}
	if doc.Version != "1.0" {
		t.Fatalf("expected version 1.0, got %q", doc.Version)
	}
}

func TestLoad_malformedYAML(t *testing.T) {
	tmp := t.TempDir()
	writeGlossary(t, tmp, "not: yaml: [broken")
	_, err := Load(tmp)
	if err == nil {
		t.Fatal("expected error for malformed yaml")
	}
	if !strings.Contains(err.Error(), "glossary.yaml") {
		t.Fatalf("error should mention glossary.yaml, got: %v", err)
	}
}

func TestLoad_invalidEntry(t *testing.T) {
	tmp := t.TempDir()
	writeGlossary(t, tmp, "version: \"1.0\"\nterms:\n  - term: \"\"\n    definition: ok\n")
	_, err := Load(tmp)
	if err == nil {
		t.Fatal("expected error for empty term")
	}
}

func TestSave_deterministicSort(t *testing.T) {
	tmp := t.TempDir()
	doc := Document{
		Version: "1.0",
		Terms: []Entry{
			{Term: "Beta", Definition: "second"},
			{Term: "alpha", Definition: "first"},
			{Term: "Gamma", Definition: "third"},
		},
	}
	if err := Save(tmp, doc); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := Load(tmp)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.Terms) != 3 {
		t.Fatalf("expected 3 terms, got %d", len(loaded.Terms))
	}
	want := []string{"alpha", "Beta", "Gamma"}
	for i, w := range want {
		if loaded.Terms[i].Term != w {
			t.Fatalf("term[%d] = %q, want %q", i, loaded.Terms[i].Term, w)
		}
	}
}

func TestUpsert_requiresTermAndDefinition(t *testing.T) {
	doc := DefaultDocument()
	_, err := Upsert(doc, Entry{Term: "", Definition: ""})
	if err == nil || !strings.Contains(err.Error(), "term") {
		t.Fatalf("expected term required error, got: %v", err)
	}
	_, err = Upsert(doc, Entry{Term: "ok", Definition: ""})
	if err == nil || !strings.Contains(err.Error(), "definition") {
		t.Fatalf("expected definition required error, got: %v", err)
	}
}

func TestUpsert_insertAndUpdate(t *testing.T) {
	doc := DefaultDocument()
	var err error
	doc, err = Upsert(doc, Entry{Term: "Foo", Definition: "bar", Aliases: []string{"old"}, Sources: []string{"initial"}})
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if len(doc.Terms) != 1 {
		t.Fatalf("expected 1 term, got %d", len(doc.Terms))
	}
	if doc.Terms[0].Created == "" {
		t.Fatal("expected created timestamp")
	}

	created := doc.Terms[0].Created
	doc, err = Upsert(doc, Entry{Term: "foo", Definition: "baz", Category: "cat", Aliases: []string{"new", "old"}, Sources: []string{"discovery"}})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if len(doc.Terms) != 1 {
		t.Fatalf("expected 1 term after upsert, got %d", len(doc.Terms))
	}
	if doc.Terms[0].Definition != "baz" {
		t.Fatalf("definition = %q, want baz", doc.Terms[0].Definition)
	}
	if doc.Terms[0].Created != created {
		t.Fatal("expected created timestamp to be preserved")
	}
	if doc.Terms[0].Category != "cat" {
		t.Fatalf("category = %q, want cat", doc.Terms[0].Category)
	}
	if got := strings.Join(doc.Terms[0].Aliases, ","); got != "old,new" {
		t.Fatalf("aliases = %q, want old,new", got)
	}
	if got := strings.Join(doc.Terms[0].Sources, ","); got != "initial,discovery" {
		t.Fatalf("sources = %q, want initial,discovery", got)
	}
}

func TestQuery_ranking(t *testing.T) {
	doc := Document{
		Terms: []Entry{
			{Term: "Alpha", Definition: "first letter", Aliases: []string{"A", "Alpha Alias"}},
			{Term: "Beta", Definition: "second letter"},
			{Term: "Alphabet", Definition: "collection of letters"},
		},
	}

	results := Query(doc, "Alpha")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Term != "Alpha" {
		t.Fatalf("expected exact match first, got %q", results[0].Term)
	}
	if results[1].Term != "Alphabet" {
		t.Fatalf("expected substring match second, got %q", results[1].Term)
	}

	aliasResults := Query(doc, "a")
	if len(aliasResults) != 3 {
		t.Fatalf("expected 3 results for substring 'a', got %d", len(aliasResults))
	}
	if aliasResults[0].Term != "Alpha" {
		t.Fatalf("expected alias match first for 'a', got %q", aliasResults[0].Term)
	}
	aliasSubstringResults := Query(doc, "Alias")
	if len(aliasSubstringResults) != 1 || aliasSubstringResults[0].Term != "Alpha" {
		t.Fatalf("expected alias substring match for Alias, got %+v", aliasSubstringResults)
	}

	if len(Query(doc, "")) != 0 {
		t.Fatal("expected no results for empty query")
	}
	if len(Query(doc, "unknown")) != 0 {
		t.Fatal("expected no results for unknown query")
	}
}

func writeGlossary(t *testing.T, root, content string) {
	t.Helper()
	path := filepath.Join(root, ".flexspec", glossaryFile)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
