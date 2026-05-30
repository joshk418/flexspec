package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	flexspecDir = ".flexspec"
	configFile  = "config.yaml"
)

// Config holds project settings from .flexspec/config.yaml.
type Config struct {
	SpecsDir      string `yaml:"specs_dir"`
	AlwaysOneShot bool   `yaml:"always_one_shot"`
	// SpecTemplate optionally forces the template /flexspec uses ("simple" or
	// "expanded"). It has no default: an empty value means /flexspec infers the
	// template from the work's size.
	SpecTemplate string `yaml:"spec_template"`
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
	if cfg.SpecsDir == "" {
		return Config{}, fmt.Errorf("%s: specs_dir must be set", path)
	}
	return cfg, nil
}
