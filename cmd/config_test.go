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

func TestConfigSet_missingConfig(t *testing.T) {
	root := t.TempDir()

	var buf bytes.Buffer
	configSetCmd.SetOut(&buf)
	configSetCmd.SetErr(&buf)

	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	err := configSetCmd.RunE(configSetCmd, []string{"specs_dir", "specs"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "flexspec init") {
		t.Errorf("error = %v, want mention of flexspec init", err)
	}
}

func TestConfigSet_wrongArgs(t *testing.T) {
	err := configSetCmd.RunE(configSetCmd, []string{"specs_dir"})
	if err == nil {
		t.Fatal("expected error for missing value arg")
	}
}

func TestConfigSet_updates(t *testing.T) {
	tests := []struct {
		name    string
		before  string
		key     string
		value   string
		wantOut []string
		verify  func(t *testing.T, root string)
	}{
		{
			name:    "specs_dir",
			before:  "specs_dir: specs\nalways_one_shot: false\nspec_template:\n",
			key:     "specs_dir",
			value:   "custom-specs",
			wantOut: []string{"custom-specs"},
			verify: func(t *testing.T, root string) {
				cfg, err := config.Load(root)
				if err != nil {
					t.Fatal(err)
				}
				if cfg.SpecsDir != "custom-specs" {
					t.Errorf("SpecsDir = %q", cfg.SpecsDir)
				}
			},
		},
		{
			name:    "always_one_shot true",
			before:  "specs_dir: specs\nalways_one_shot: false\nspec_template:\n",
			key:     "always_one_shot",
			value:   "true",
			wantOut: []string{"always_one_shot", "true"},
			verify: func(t *testing.T, root string) {
				cfg, err := config.Load(root)
				if err != nil {
					t.Fatal(err)
				}
				if !cfg.AlwaysOneShot {
					t.Error("AlwaysOneShot = false, want true")
				}
			},
		},
		{
			name:    "spec_template expanded",
			before:  "specs_dir: specs\nalways_one_shot: false\nspec_template:\n",
			key:     "spec_template",
			value:   "expanded",
			wantOut: []string{"spec_template", "expanded"},
			verify: func(t *testing.T, root string) {
				cfg, err := config.Load(root)
				if err != nil {
					t.Fatal(err)
				}
				if cfg.SpecTemplate != "expanded" {
					t.Errorf("SpecTemplate = %q", cfg.SpecTemplate)
				}
			},
		},
		{
			name:    "spec_template empty",
			before:  "specs_dir: specs\nalways_one_shot: false\nspec_template: simple\n",
			key:     "spec_template",
			value:   "",
			wantOut: []string{"spec_template", "-"},
			verify: func(t *testing.T, root string) {
				cfg, err := config.Load(root)
				if err != nil {
					t.Fatal(err)
				}
				if cfg.SpecTemplate != "" {
					t.Errorf("SpecTemplate = %q, want empty", cfg.SpecTemplate)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeConfigYAML(t, root, tt.before)

			var buf bytes.Buffer
			configSetCmd.SetOut(&buf)

			oldWd, _ := os.Getwd()
			if err := os.Chdir(root); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.Chdir(oldWd) }()

			if err := configSetCmd.RunE(configSetCmd, []string{tt.key, tt.value}); err != nil {
				t.Fatal(err)
			}

			out := buf.String()
			for _, want := range tt.wantOut {
				if !strings.Contains(out, want) {
					t.Errorf("output missing %q\n%s", want, out)
				}
			}
			if !strings.Contains(out, "KEY") || !strings.Contains(out, "config:") {
				t.Errorf("expected table output\n%s", out)
			}
			tt.verify(t, root)
		})
	}
}

func TestConfigSet_invalidKeyOrValue(t *testing.T) {
	tests := []struct {
		name   string
		before string
		key    string
		value  string
	}{
		{"unknown key", "specs_dir: specs\nalways_one_shot: false\n", "unknown", "x"},
		{"invalid bool", "specs_dir: specs\nalways_one_shot: false\n", "always_one_shot", "maybe"},
		{"invalid template", "specs_dir: specs\nalways_one_shot: false\n", "spec_template", "bad"},
		{"empty specs_dir", "specs_dir: specs\nalways_one_shot: false\n", "specs_dir", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			writeConfigYAML(t, root, tt.before)

			oldWd, _ := os.Getwd()
			if err := os.Chdir(root); err != nil {
				t.Fatal(err)
			}
			defer func() { _ = os.Chdir(oldWd) }()

			err := configSetCmd.RunE(configSetCmd, []string{tt.key, tt.value})
			if err == nil {
				t.Fatal("expected error")
			}

			cfg, loadErr := config.Load(root)
			if loadErr != nil {
				t.Fatal(loadErr)
			}
			if tt.key == "specs_dir" && cfg.SpecsDir != "specs" {
				t.Errorf("config was modified: SpecsDir = %q", cfg.SpecsDir)
			}
			if tt.key == "always_one_shot" && cfg.AlwaysOneShot {
				t.Error("config was modified")
			}
			if tt.key == "spec_template" && cfg.SpecTemplate != "" {
				t.Errorf("config was modified: SpecTemplate = %q", cfg.SpecTemplate)
			}
		})
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
