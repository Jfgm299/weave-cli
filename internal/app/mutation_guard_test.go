package app

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestEnsureOpsWithinRoot_AllowsPathsInsideRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	ops := []fsops.Operation{{Type: fsops.OpEnsureDir, Path: filepath.Join(root, ".agents", "skills")}}

	if err := ensureOpsWithinRoot(root, ops); err != nil {
		t.Fatalf("expected no guard error, got %v", err)
	}
}

func TestEnsureOpsWithinRoot_RejectsPathOutsideRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	ops := []fsops.Operation{{Type: fsops.OpEnsureDir, Path: filepath.Join(root, "..", "outside")}}

	err := ensureOpsWithinRoot(root, ops)
	if err == nil {
		t.Fatalf("expected unsafe mutation guard error")
	}

	if !errors.Is(err, ErrUnsafeMutationPath) {
		t.Fatalf("expected ErrUnsafeMutationPath, got %v", err)
	}
}
