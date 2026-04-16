package app

import (
	"context"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/config"
)

func TestMigrationService_Migrate_UpgradesOutdatedSchema(t *testing.T) {
	t.Parallel()

	writer := &writerSpy{}
	sut := MigrationService{Writer: writer}

	result, migrated, err := sut.Migrate(context.Background(), config.Config{Version: 0}, false)
	if err != nil {
		t.Fatalf("unexpected migrate error: %v", err)
	}
	if !result.Upgraded || migrated.Version != config.CurrentSchemaVersion {
		t.Fatalf("expected schema to be upgraded, got result=%+v cfg=%+v", result, migrated)
	}
	if !writer.called {
		t.Fatalf("expected writer call for non-dry migration")
	}
}

func TestMigrationService_Migrate_DryRunDoesNotPersist(t *testing.T) {
	t.Parallel()

	writer := &writerSpy{}
	sut := MigrationService{Writer: writer}

	result, _, err := sut.Migrate(context.Background(), config.Config{Version: 0}, true)
	if err != nil {
		t.Fatalf("unexpected migrate error: %v", err)
	}
	if !result.Upgraded || !result.DryRun {
		t.Fatalf("expected dry-run upgrade result, got %+v", result)
	}
	if writer.called {
		t.Fatalf("writer should not run on dry-run migration")
	}
}
