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
