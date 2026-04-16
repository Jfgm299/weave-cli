package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileLoader_LoadOrDefault_ReturnsDefaultsWhenMissing(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	loader := FileLoader{Path: filepath.Join(tmp, "weave.yaml")}

	cfg, err := loader.LoadOrDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Version != 1 {
		t.Fatalf("expected default version 1, got %d", cfg.Version)
	}

	if cfg.Sync.Mode != "symlink" {
		t.Fatalf("expected default symlink mode, got %s", cfg.Sync.Mode)
	}
}

func TestFileLoader_LoadOrDefault_LoadsExistingFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "weave.yaml")
	body := []byte("version: 1\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	loader := FileLoader{Path: path}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Version != 1 || cfg.Sync.Mode != "symlink" {
		t.Fatalf("unexpected config loaded: %+v", cfg)
	}

	if cfg.Providers == nil {
		t.Fatalf("expected providers inventory to be initialized")
	}
}
