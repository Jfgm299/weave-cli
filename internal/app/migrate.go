package app

import (
	"context"
	"fmt"

	"github.com/Jfgm299/weave-cli/internal/config"
)

type MigrationService struct {
	Validator config.Validator
	Writer    ConfigWriter
}

type MigrationResult struct {
	Upgraded bool
	From     int
	To       int
	DryRun   bool
}

func (s MigrationService) Migrate(_ context.Context, cfg config.Config, dryRun bool) (MigrationResult, config.Config, error) {
	if cfg.Version >= config.CurrentSchemaVersion {
		return MigrationResult{Upgraded: false, From: cfg.Version, To: config.CurrentSchemaVersion, DryRun: dryRun}, cfg, nil
	}

	next := cfg
	next.Version = config.CurrentSchemaVersion
	if next.Sync.Mode == "" {
		next.Sync.Mode = "symlink"
	}
	if next.Sync.ConflictPolicy == "" {
		next.Sync.ConflictPolicy = "prompt"
	}

	if err := (config.Validator{}).Validate(next); err != nil {
		return MigrationResult{}, cfg, fmt.Errorf("failed to validate migrated config: %w", err)
	}

	result := MigrationResult{Upgraded: true, From: cfg.Version, To: config.CurrentSchemaVersion, DryRun: dryRun}
	if dryRun {
		return result, next, nil
	}

	if s.Writer != nil {
		if err := s.Writer.Write(next); err != nil {
			return MigrationResult{}, cfg, err
		}
	}

	return result, next, nil
}
