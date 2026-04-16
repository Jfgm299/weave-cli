package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/app"
)

func TestParseFromFlag_ParsesLongForm(t *testing.T) {
	t.Parallel()

	g, err := parseFromFlag([]string{"--from", "/tmp/src"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if g != "/tmp/src" {
		t.Fatalf("expected /tmp/src, got %q", g)
	}
}

func TestParseFromFlag_ParsesEqualsForm(t *testing.T) {
	t.Parallel()

	g, err := parseFromFlag([]string{"--from=/tmp/src"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if g != "/tmp/src" {
		t.Fatalf("expected /tmp/src, got %q", g)
	}
}

func TestParseFromFlag_MissingValueFails(t *testing.T) {
	t.Parallel()

	_, err := parseFromFlag([]string{"--from"})
	if err == nil {
		t.Fatalf("expected missing value error")
	}
}

func TestParseProviderAction_ListParsesWithoutName(t *testing.T) {
	t.Parallel()

	action, name, dryRun, err := parseProviderAction([]string{"list"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if action != providerActionList {
		t.Fatalf("expected providerActionList, got %s", action)
	}

	if name != "" {
		t.Fatalf("expected empty provider name for list, got %q", name)
	}

	if dryRun {
		t.Fatalf("expected dryRun false for list")
	}
}

func TestParseProviderAction_AddMissingNameFails(t *testing.T) {
	t.Parallel()

	_, _, _, err := parseProviderAction([]string{"add"})
	if err == nil {
		t.Fatalf("expected provider add parse error when name is missing")
	}
}

func TestParseProviderAction_AddParsesDryRun(t *testing.T) {
	t.Parallel()

	action, name, dryRun, err := parseProviderAction([]string{"add", "claude-code", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if action != providerActionAdd || name != "claude-code" || !dryRun {
		t.Fatalf("unexpected parse output: action=%s name=%s dryRun=%v", action, name, dryRun)
	}
}

func TestParseProviderAction_ListRejectsDryRun(t *testing.T) {
	t.Parallel()

	_, _, _, err := parseProviderAction([]string{"list", "--dry-run"})
	if err == nil {
		t.Fatalf("expected list to reject --dry-run")
	}
}

func TestParseAddFlags_RejectsUnknownFlag(t *testing.T) {
	t.Parallel()

	_, _, _, err := parseAddFlags([]string{"--unknown"})
	if err == nil {
		t.Fatalf("expected unsupported flag error")
	}
}

func TestParseAddFlags_ParsesConflictPolicyFlags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		args   []string
		policy app.ConflictPolicy
	}{
		{name: "overwrite", args: []string{"--overwrite"}, policy: app.ConflictPolicyOverwrite},
		{name: "skip", args: []string{"--skip"}, policy: app.ConflictPolicySkip},
		{name: "backup", args: []string{"--backup"}, policy: app.ConflictPolicyBackup},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, _, policy, err := parseAddFlags(tc.args)
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}

			if policy != tc.policy {
				t.Fatalf("expected policy %q, got %q", tc.policy, policy)
			}
		})
	}
}

func TestParseAddFlags_RejectsMultipleConflictPolicyFlags(t *testing.T) {
	t.Parallel()

	_, _, _, err := parseAddFlags([]string{"--overwrite", "--backup"})
	if err == nil {
		t.Fatalf("expected conflict policy parsing error")
	}
}

func TestParseDryRunOnly_RejectsUnknownFlag(t *testing.T) {
	t.Parallel()

	_, err := parseDryRunOnly([]string{"--json"}, "forge")
	if err == nil {
		t.Fatalf("expected unsupported flag error")
	}
}

func TestRun_HelpPrintsQuickstart(t *testing.T) {
	t.Parallel()

	out := captureStdout(t, func() {
		code, err := Run(context.Background(), []string{"--help"})
		if err != nil {
			t.Fatalf("unexpected help error: %v", err)
		}
		if code != ExitOK {
			t.Fatalf("expected ExitOK, got %d", code)
		}
	})

	if !strings.Contains(out, "60-second quickstart") {
		t.Fatalf("expected quickstart in help output, got: %s", out)
	}
}

func TestRun_MigrateDryRun(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte("version: 0\nsync:\n  mode: symlink\nskills: []\ncommands: []\n"), 0o644); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	out := captureStdout(t, func() {
		code, err := Run(context.Background(), []string{"migrate", "--dry-run"})
		if err != nil {
			t.Fatalf("unexpected migrate error: %v", err)
		}
		if code != ExitOK {
			t.Fatalf("expected ExitOK, got %d", code)
		}
	})

	if !strings.Contains(out, "would upgrade weave.yaml schema") {
		t.Fatalf("expected dry-run migration summary, got: %s", out)
	}
}
