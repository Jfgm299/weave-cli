package config

import "fmt"

type Validator struct{}

func (Validator) Validate(cfg Config) error {
	if cfg.Version <= 0 {
		return fmt.Errorf("version is required and must be greater than 0")
	}

	if cfg.Sync.Mode != "symlink" {
		return fmt.Errorf("sync.mode must be 'symlink' in v1")
	}

	return nil
}
