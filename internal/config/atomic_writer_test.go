package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAtomicFileWriter_Write_UpdatesExistingConfigAtomically(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "weave.yaml")
	if err := os.WriteFile(path, []byte("version: 1\nsync:\n  mode: symlink\nskills: []\ncommands: []\n"), 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	w := AtomicFileWriter{Path: path}
	err := w.Write(Config{
		Version: 1,
		Sync:    Sync{Mode: "symlink"},
		Skills:  []Asset{{Name: "sdd-orchestrator", Source: "/tmp/src"}},
	})
	if err != nil {
		t.Fatalf("write failed: %v", err)
	}

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	out := string(b)
	if out == "version: 1\nsync:\n  mode: symlink\nskills: []\ncommands: []\n" {
		t.Fatalf("expected config to be updated")
	}
	if !containsAll(out, "sdd-orchestrator", "source: /tmp/src") {
		t.Fatalf("expected persisted skill entry, got: %s", out)
	}
}

func TestAtomicFileWriter_Write_RenameFailureLeavesOriginalUnchanged(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "weave.yaml")
	original := []byte("version: 1\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	w := AtomicFileWriter{
		Path: path,
		renameFn: func(_, _ string) error {
			return errors.New("rename failed")
		},
	}
	err := w.Write(Config{Version: 1, Sync: Sync{Mode: "symlink"}, Skills: []Asset{{Name: "x", Source: "y"}}})
	if err == nil {
		t.Fatalf("expected error")
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if string(after) != string(original) {
		t.Fatalf("expected original config unchanged when rename fails")
	}
}

func containsAll(s string, tokens ...string) bool {
	for _, tk := range tokens {
		if !strings.Contains(s, tk) {
			return false
		}
	}
	return true
}
