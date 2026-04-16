//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestForge_E2E_NoGitRoot_ForcedInteractiveDeclineShowsDeclinedGuidance(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()

	// Simulate an interactive user declining git init.
	out, err := runCLI(repo, root, []string{"forge"}, []string{"WEAVE_FORCE_INTERACTIVE=1", "WEAVE_TEST_STDIN=n\n"})
	if err == nil {
		t.Fatalf("expected forge to fail when interactive git-init prompt is declined")
	}

	if !strings.Contains(out, "No git repository detected") {
		t.Fatalf("expected interactive prompt text, got: %s", out)
	}

	if !strings.Contains(out, "initialization was declined") {
		t.Fatalf("expected declined guidance, got: %s", out)
	}
}
