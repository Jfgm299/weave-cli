//go:build e2e

package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProviderAddAndList_E2E_MultiProviderFlowAndActionableErrors(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	binDir := filepath.Join(repo, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("mkdir bin dir: %v", err)
	}
	for _, name := range []string{"claude", "opencode"} {
		p := filepath.Join(binDir, name)
		if err := os.WriteFile(p, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("write fake binary %s: %v", name, err)
		}
	}

	out, err := runCLI(repo, root, []string{"provider", "add", "claude-code"}, []string{"PATH=" + binDir + ":" + os.Getenv("PATH")})
	if err != nil {
		t.Fatalf("provider add claude-code failed: %v\n%s", err, out)
	}

	out, err = runCLI(repo, root, []string{"provider", "add", "opencode"}, []string{"PATH=" + binDir + ":" + os.Getenv("PATH")})
	if err != nil {
		t.Fatalf("provider add opencode failed: %v\n%s", err, out)
	}

	out, err = runCLI(repo, root, []string{"provider", "list"}, []string{"PATH=" + binDir + ":" + os.Getenv("PATH")})
	if err != nil {
		t.Fatalf("provider list failed: %v\n%s", err, out)
	}

	if !strings.Contains(out, "claude-code") || !strings.Contains(out, "opencode") {
		t.Fatalf("expected provider list to contain both providers, got: %s", out)
	}

	b, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	cfgOut := string(b)
	if !strings.Contains(cfgOut, "- name: claude-code") || !strings.Contains(cfgOut, "- name: opencode") {
		t.Fatalf("expected both providers in config, got: %s", cfgOut)
	}
}

func TestProviderAdd_E2E_MissingBinaryActionableFailureAndNoConfigMutation(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	seed := []byte("version: 1\nproviders: []\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"provider", "add", "claude-code"}, []string{"PATH=/usr/bin:/bin"})
	if err == nil {
		t.Fatalf("expected provider add to fail when binary is missing")
	}

	if !strings.Contains(out, "Install the missing binaries") || !strings.Contains(out, "weave provider repair claude-code") {
		t.Fatalf("expected actionable failure message, got: %s", out)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged when provider binary is missing")
	}
}

func TestProviderRemoveAndRepair_E2E_ReversibleOperations(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	binDir := filepath.Join(repo, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("mkdir bin dir: %v", err)
	}
	for _, name := range []string{"claude", "opencode"} {
		p := filepath.Join(binDir, name)
		if err := os.WriteFile(p, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
			t.Fatalf("write fake binary %s: %v", name, err)
		}
	}

	pathEnv := []string{"PATH=" + binDir + ":" + os.Getenv("PATH")}

	if out, err := runCLI(repo, root, []string{"provider", "add", "claude-code"}, pathEnv); err != nil {
		t.Fatalf("provider add failed: %v\n%s", err, out)
	}

	if _, err := os.Lstat(filepath.Join(repo, ".claude", "CLAUDE.md")); err != nil {
		t.Fatalf("expected provider link after add: %v", err)
	}

	if out, err := runCLI(repo, root, []string{"provider", "remove", "claude-code"}, pathEnv); err != nil {
		t.Fatalf("provider remove failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(repo, ".claude")); !os.IsNotExist(err) {
		t.Fatalf("expected .claude removed after provider remove, got err: %v", err)
	}

	if out, err := runCLI(repo, root, []string{"provider", "repair", "claude-code"}, pathEnv); err != nil {
		t.Fatalf("provider repair failed: %v\n%s", err, out)
	}

	if _, err := os.Lstat(filepath.Join(repo, ".claude", "CLAUDE.md")); err != nil {
		t.Fatalf("expected provider link restored after repair: %v", err)
	}
}
