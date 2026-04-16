package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type FileLoader struct {
	Path string
}

func (l FileLoader) LoadOrDefault() (Config, error) {
	if l.Path == "" {
		return Default(), nil
	}

	b, err := os.ReadFile(l.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("invalid weave.yaml: %w", err)
	}

	if cfg.Skills == nil {
		cfg.Skills = []Asset{}
	}

	if cfg.Commands == nil {
		cfg.Commands = []Asset{}
	}

	if cfg.Providers == nil {
		cfg.Providers = []Provider{}
	}

	return cfg, nil
}
