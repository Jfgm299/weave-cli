//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestForge_E2E_NoGitRoot_NonInteractiveShowsActionableGuidance(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()

	out, err := runCLI(repo, root, []string{"forge"}, []string{"WEAVE_NON_INTERACTIVE=1"})
	if err == nil {
		t.Fatalf("expected forge to fail without .git in non-interactive mode")
	}

	if !strings.Contains(out, "project root not detected") {
		t.Fatalf("expected root detection error, got: %s", out)
	}

	if !strings.Contains(out, "Run `git init` and retry") {
		t.Fatalf("expected actionable git init guidance, got: %s", out)
	}
}
