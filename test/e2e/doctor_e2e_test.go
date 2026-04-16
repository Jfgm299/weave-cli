//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctor_E2E_ReportsIssuesAndDeterministicExitCode(t *testing.T) {
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

	source := filepath.Join(repo, "shared", "skills", "sdd-orchestrator", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.WriteFile(source, []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	cfg := "version: 1\n" +
		"sources:\n" +
		"  skills_dir: ~/.weave/skills\n" +
		"  commands_dir: ~/.weave/commands\n" +
		"sync:\n" +
		"  mode: symlink\n" +
		"skills:\n" +
		"  - name: sdd-orchestrator\n" +
		"    source: " + source + "\n" +
		"commands: []\n"

	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"doctor"}, nil)
	if err == nil {
		t.Fatalf("expected doctor to fail with issues")
	}

	if !strings.Contains(out, "exit status 5") {
		t.Fatalf("expected go run output to include deterministic doctor exit status 5, got: %s", out)
	}

	if !strings.Contains(out, "Doctor status: issues_found") {
		t.Fatalf("expected issues status output, got: %s", out)
	}
	if !strings.Contains(out, "weave skill add sdd-orchestrator") {
		t.Fatalf("expected repair command in output, got: %s", out)
	}
	if !strings.Contains(out, "docs/reference/doctor.md") {
		t.Fatalf("expected docs reference in output, got: %s", out)
	}
}

func TestDoctor_E2E_JSONOutputIsScriptFriendly(t *testing.T) {
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

	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n"), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"doctor", "--json"}, nil)
	if err != nil {
		t.Fatalf("expected doctor --json success: %v\n%s", err, out)
	}

	var payload struct {
		Status string `json:"status"`
		Issues []struct {
			DocsPath string `json:"docs_path"`
			DocsURL  string `json:"docs_url"`
		} `json:"issues"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("expected json output, got error: %v\noutput: %s", err, out)
	}

	if payload.Status != "healthy" {
		t.Fatalf("expected healthy status, got %q", payload.Status)
	}
}
