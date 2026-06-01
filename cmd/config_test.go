package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joshk418/flexspec/internal/config"
)

func TestConfig_missingFile(t *testing.T) {
	root := t.TempDir()
	configJSON = false

	var buf bytes.Buffer
	configCmd.SetOut(&buf)
	configCmd.SetErr(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	err := configCmd.RunE(configCmd, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "flexspec init") {
		t.Errorf("error = %v, want mention of flexspec init", err)
	}
}

func TestConfig_human(t *testing.T) {
	root := t.TempDir()
	writeConfigYAML(t, root, "specs_dir: my-specs\nalways_one_shot: true\nspec_template: expanded\n")

	configJSON = false
	var buf bytes.Buffer
	configCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := configCmd.RunE(configCmd, nil); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	for _, want := range []string{
		"KEY",
		"VALUE",
		"specs_dir",
		"my-specs",
		"always_one_shot",
		"true",
		"spec_template",
		"expanded",
		"config:",
		".flexspec" + string(os.PathSeparator) + "config.yaml",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestConfig_human_emptyTemplate(t *testing.T) {
	root := t.TempDir()
	writeConfigYAML(t, root, "specs_dir: specs\nalways_one_shot: false\nspec_template:\n")

	configJSON = false
	var buf bytes.Buffer
	configCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := configCmd.RunE(configCmd, nil); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "spec_template") || !strings.Contains(out, "-") {
		t.Errorf("expected dash for empty spec_template\n%s", out)
	}
}

func TestConfig_json(t *testing.T) {
	root := t.TempDir()
	writeConfigYAML(t, root, "specs_dir: specs\nalways_one_shot: false\nspec_template:\n")

	configJSON = true
	var buf bytes.Buffer
	configCmd.SetOut(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := configCmd.RunE(configCmd, nil); err != nil {
		t.Fatal(err)
	}

	var doc config.JSONDocument
	if err := json.Unmarshal(buf.Bytes(), &doc); err != nil {
		t.Fatal(err)
	}
	if doc.SpecsDir != "specs" || doc.AlwaysOneShot != false || doc.SpecTemplate != "" {
		t.Errorf("got %+v", doc)
	}
}

func TestConfig_inRootHelp(t *testing.T) {
	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "config") {
		t.Error("root help should list config subcommand")
	}
}

func writeConfigYAML(t *testing.T, root, content string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(base, "config.yaml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
