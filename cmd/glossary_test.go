package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshk418/flexspec/internal/glossary"
)

func TestGlossaryList_empty(t *testing.T) {
	root := t.TempDir()
	glossaryJSON = false

	var buf bytes.Buffer
	glossaryListCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := glossaryListCmd.RunE(glossaryListCmd, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No glossary terms") {
		t.Fatalf("expected empty message, got:\n%s", buf.String())
	}
}

func TestGlossaryList_table(t *testing.T) {
	root := t.TempDir()
	writeGlossaryYAML(t, root, "version: \"1.0\"\nterms:\n  - term: Alpha\n    definition: first\n    category: letters\n  - term: Beta\n    definition: second\n")
	glossaryJSON = false

	var buf bytes.Buffer
	glossaryListCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := glossaryListCmd.RunE(glossaryListCmd, nil); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"TERM", "DEFINITION", "CATEGORY", "Alpha", "Beta", "letters"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestGlossaryList_json(t *testing.T) {
	root := t.TempDir()
	writeGlossaryYAML(t, root, "version: \"1.0\"\nterms:\n  - term: Alpha\n    definition: first\n")
	glossaryJSON = true

	var buf bytes.Buffer
	glossaryListCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := glossaryListCmd.RunE(glossaryListCmd, nil); err != nil {
		t.Fatal(err)
	}
	var doc glossary.Document
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Terms) != 1 || doc.Terms[0].Term != "Alpha" {
		t.Fatalf("unexpected json result: %+v", doc)
	}
}

func TestGlossaryQuery_empty(t *testing.T) {
	root := t.TempDir()
	glossaryJSON = false

	var buf bytes.Buffer
	glossaryQueryCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := glossaryQueryCmd.RunE(glossaryQueryCmd, []string{"foo"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "No matches") {
		t.Fatalf("expected no-matches message, got:\n%s", buf.String())
	}
}

func TestGlossaryQuery_table(t *testing.T) {
	root := t.TempDir()
	writeGlossaryYAML(t, root, "version: \"1.0\"\nterms:\n  - term: Alpha\n    definition: first\n  - term: Beta\n    definition: second\n")
	glossaryJSON = false

	var buf bytes.Buffer
	glossaryQueryCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := glossaryQueryCmd.RunE(glossaryQueryCmd, []string{"Alpha"}); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alpha") || strings.Contains(out, "Beta") {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestGlossaryQuery_missingArg(t *testing.T) {
	err := glossaryQueryCmd.RunE(glossaryQueryCmd, []string{})
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestGlossaryAdd_missingDefinition(t *testing.T) {
	addDefinition = ""
	err := glossaryAddCmd.RunE(glossaryAddCmd, []string{"term"})
	if err == nil || !strings.Contains(err.Error(), "definition") {
		t.Fatalf("expected definition required error, got: %v", err)
	}
}

func TestGlossaryAdd_andQuery(t *testing.T) {
	root := t.TempDir()
	glossaryJSON = false

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	addDefinition = "a test term"
	addAliases = []string{"alias1", "alias2"}
	addCategory = "test"
	addSources = []string{"discovery"}

	var addBuf bytes.Buffer
	glossaryAddCmd.SetOut(&addBuf)
	if err := glossaryAddCmd.RunE(glossaryAddCmd, []string{"TestTerm"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(addBuf.String(), "TestTerm") {
		t.Fatalf("expected term in add output, got:\n%s", addBuf.String())
	}

	var queryBuf bytes.Buffer
	glossaryQueryCmd.SetOut(&queryBuf)
	if err := glossaryQueryCmd.RunE(glossaryQueryCmd, []string{"TestTerm"}); err != nil {
		t.Fatal(err)
	}
	out := queryBuf.String()
	for _, want := range []string{"TestTerm", "a test term", "test"} {
		if !strings.Contains(out, want) {
			t.Errorf("query output missing %q\n%s", want, out)
		}
	}
}

func writeGlossaryYAML(t *testing.T, root, content string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "glossary.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
