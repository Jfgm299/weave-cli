package config

import "testing"

func TestDefault_ReturnsV1SymlinkConfigWithExplicitInventory(t *testing.T) {
	t.Parallel()

	cfg := Default()

	if cfg.Version != 1 {
		t.Fatalf("expected version 1, got %d", cfg.Version)
	}

	if cfg.Sync.Mode != "symlink" {
		t.Fatalf("expected symlink mode, got %q", cfg.Sync.Mode)
	}

	if cfg.Sources.SkillsDir != "~/.weave/skills" {
		t.Fatalf("expected default skills source, got %q", cfg.Sources.SkillsDir)
	}

	if cfg.Sources.CommandsDir != "~/.weave/commands" {
		t.Fatalf("expected default commands source, got %q", cfg.Sources.CommandsDir)
	}

	if cfg.Skills == nil {
		t.Fatalf("expected explicit skills inventory slice")
	}

	if cfg.Commands == nil {
		t.Fatalf("expected explicit commands inventory slice")
	}

	if cfg.Providers == nil {
		t.Fatalf("expected explicit providers inventory slice")
	}
}
