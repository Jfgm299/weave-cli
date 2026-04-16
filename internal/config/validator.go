package config

import (
	"errors"
	"fmt"
)

var (
	ErrOutdatedSchema    = errors.New("outdated schema")
	ErrUnsupportedSchema = errors.New("unsupported schema")
)

const CurrentSchemaVersion = 1

type Validator struct{}

func (Validator) Validate(cfg Config) error {
	if cfg.Version < CurrentSchemaVersion {
		return fmt.Errorf("%w: found version %d, expected %d. Run `weave migrate` to upgrade the project config", ErrOutdatedSchema, cfg.Version, CurrentSchemaVersion)
	}

	if cfg.Version > CurrentSchemaVersion {
		return fmt.Errorf("%w: found version %d, this binary supports up to %d", ErrUnsupportedSchema, cfg.Version, CurrentSchemaVersion)
	}

	if cfg.Sync.Mode != "symlink" {
		return fmt.Errorf("sync.mode must be 'symlink' in v1")
	}

	if cfg.Sync.ConflictPolicy != "" {
		switch cfg.Sync.ConflictPolicy {
		case "prompt", "overwrite", "skip", "backup":
		default:
			return fmt.Errorf("sync.conflict_policy must be one of prompt|overwrite|skip|backup")
		}
	}

	return nil
}
