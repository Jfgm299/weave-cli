package fsops

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestEngine_Apply_CreatesSymlinkOperations(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	target := filepath.Join(tmp, "source.txt")
	link := filepath.Join(tmp, "dir", "link.txt")

	if err := os.WriteFile(target, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write target: %v", err)
	}

	err := (Engine{}).Apply(context.Background(), []Operation{{
		Type:   OpCreateLink,
		Path:   link,
		Target: target,
	}})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	fi, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("lstat: %v", err)
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected symlink at %s", link)
	}
}

func TestEngine_Apply_BackupPathRenamesExistingPath(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	src := filepath.Join(root, ".agents", "commands", "pr-review.md")
	backup := src + ".bak"
	if err := os.MkdirAll(filepath.Dir(src), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(src, []byte("old"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	err := (Engine{}).Apply(context.Background(), []Operation{{Type: OpBackupPath, Path: src, Target: backup}})
	if err != nil {
		t.Fatalf("unexpected backup error: %v", err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatalf("expected original source moved, got err: %v", err)
	}
	if b, err := os.ReadFile(backup); err != nil {
		t.Fatalf("read backup: %v", err)
	} else if string(b) != "old" {
		t.Fatalf("expected backed up content, got %q", string(b))
	}
}
