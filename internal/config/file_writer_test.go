package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileWriter_Write_RoundtripDeterministic(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "weave.yaml")

	w := FileWriter{Path: path}
	cfg := Config{
		Version: 1,
		Sync:    Sync{Mode: "symlink"},
		Skills: []Asset{
			{Name: "zeta", Source: "z"},
			{Name: "alpha", Source: "a"},
		},
	}

	if err := w.Write(cfg); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	b1, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after first write: %v", err)
	}

	if err := w.Write(cfg); err != nil {
		t.Fatalf("write second time failed: %v", err)
	}

	b2, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after second write: %v", err)
	}

	if string(b1) != string(b2) {
		t.Fatalf("expected deterministic output")
	}
}

func TestFileWriter_Write_DoesNotOverwriteExistingConfig(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "weave.yaml")
	original := []byte("version: 1\nsync:\n  mode: symlink\nskills:\n  - name: keep\n    source: x\ncommands: []\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	w := FileWriter{Path: path}
	cfg := Default()

	if err := w.Write(cfg); err != nil {
		t.Fatalf("write existing should be no-op, got %v", err)
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after write: %v", err)
	}

	if string(after) != string(original) {
		t.Fatalf("expected existing config to remain unchanged")
	}
}
