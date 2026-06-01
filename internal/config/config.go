package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	flexspecDir = ".flexspec"
	configFile  = "config.yaml"
)

// Config holds project settings from .flexspec/config.yaml.
type Config struct {
	SpecsDir      string `yaml:"specs_dir" json:"specs_dir"`
	AlwaysOneShot bool   `yaml:"always_one_shot" json:"always_one_shot"`
	// SpecTemplate optionally forces the template /flexspec uses ("simple" or
	// "expanded"). It has no default: an empty value means /flexspec infers the
	// template from the work's size.
	SpecTemplate string `yaml:"spec_template" json:"spec_template"`
}

// Entry is one row for human-readable config output (flexspec config).
type Entry struct {
	Key   string
	Value string
}

// JSONDocument is the machine-readable shape for flexspec config --json.
type JSONDocument struct {
	SpecsDir      string `json:"specs_dir"`
	AlwaysOneShot bool   `json:"always_one_shot"`
	SpecTemplate  string `json:"spec_template"`
}

// DisplayEntries returns known config keys in fixed order for table output.
func DisplayEntries(cfg Config) []Entry {
	return []Entry{
		{Key: "specs_dir", Value: cfg.SpecsDir},
		{Key: "always_one_shot", Value: strconv.FormatBool(cfg.AlwaysOneShot)},
		{Key: "spec_template", Value: displayOrDash(cfg.SpecTemplate)},
	}
}

// JSONDocumentFromConfig builds the --json payload for cfg.
func JSONDocumentFromConfig(cfg Config) JSONDocument {
	return JSONDocument(cfg)
}

func displayOrDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "-"
	}
	return s
}

// ApplyUpdate returns cfg with one known key updated. Value is parsed by key type.
func ApplyUpdate(cfg Config, key, value string) (Config, error) {
	switch key {
	case "specs_dir":
		if strings.TrimSpace(value) == "" {
			return Config{}, fmt.Errorf("specs_dir must be set")
		}
		cfg.SpecsDir = value
	case "always_one_shot":
		b, err := strconv.ParseBool(value)
		if err != nil {
			return Config{}, fmt.Errorf("always_one_shot: invalid bool %q", value)
		}
		cfg.AlwaysOneShot = b
	case "spec_template":
		if value != "" && value != "simple" && value != "expanded" {
			return Config{}, fmt.Errorf("spec_template must be simple, expanded, or empty")
		}
		cfg.SpecTemplate = value
	default:
		return Config{}, fmt.Errorf("unknown config key %q", key)
	}
	return cfg, nil
}

// Load reads .flexspec/config.yaml under root.
func Load(root string) (Config, error) {
	path := filepath.Join(root, flexspecDir, configFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("%s not found; run `flexspec init` first", filepath.Join(flexspecDir, configFile))
		}
		return Config{}, fmt.Errorf("read %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validate(&cfg, path); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Save writes config to .flexspec/config.yaml under root.
func Save(root string, cfg Config) error {
	if err := validate(&cfg, filepath.Join(root, flexspecDir, configFile)); err != nil {
		return err
	}
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	path := filepath.Join(root, flexspecDir, configFile)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// ConfigPath returns the absolute path to config.yaml.
func ConfigPath(root string) string {
	return filepath.Join(root, flexspecDir, configFile)
}

func validate(cfg *Config, path string) error {
	if cfg.SpecsDir == "" {
		return fmt.Errorf("%s: specs_dir must be set", path)
	}
	if cfg.SpecTemplate != "" &&
		cfg.SpecTemplate != "simple" && cfg.SpecTemplate != "expanded" {
		return fmt.Errorf("%s: spec_template must be simple, expanded, or empty", path)
	}
	return nil
}
