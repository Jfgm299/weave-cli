package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunDoctor_HealthyProjectReturnsExitOKAndTextStatus(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	source := filepath.Join(repo, "skills-src", "sdd-orchestrator", "SKILL.md")
	installed := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(source, []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write source skill: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(installed), 0o755); err != nil {
		t.Fatalf("mkdir installed dir: %v", err)
	}
	if err := os.Symlink(source, installed); err != nil {
		t.Fatalf("seed installed symlink: %v", err)
	}

	configBody := "version: 1\n" +
		"sources:\n" +
		"  skills_dir: ~/.weave/skills\n" +
		"  commands_dir: ~/.weave/commands\n" +
		"sync:\n" +
		"  mode: symlink\n" +
		"skills:\n" +
		"  - name: sdd-orchestrator\n" +
		"    source: " + source + "\n" +
		"commands: []\n"

	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte(configBody), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	t.Setenv("WEAVE_WORKDIR", repo)
	out := captureStdout(t, func() {
		code, err := runDoctor(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected runDoctor error: %v", err)
		}
		if code != ExitOK {
			t.Fatalf("expected ExitOK for healthy project, got %d", code)
		}
	})

	if !strings.Contains(out, "Doctor status: healthy") {
		t.Fatalf("expected healthy text output, got: %s", out)
	}
}

func TestRunDoctor_MissingAssetReturnsExitDoctorIssuesAndRepairPath(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	source := filepath.Join(repo, "skills-src", "sdd-orchestrator", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(source, []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write source skill: %v", err)
	}

	configBody := "version: 1\n" +
		"sources:\n" +
		"  skills_dir: ~/.weave/skills\n" +
		"  commands_dir: ~/.weave/commands\n" +
		"sync:\n" +
		"  mode: symlink\n" +
		"skills:\n" +
		"  - name: sdd-orchestrator\n" +
		"    source: " + source + "\n" +
		"commands: []\n"

	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte(configBody), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	t.Setenv("WEAVE_WORKDIR", repo)
	out := captureStdout(t, func() {
		code, err := runDoctor(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected runDoctor error: %v", err)
		}
		if code != ExitDoctorIssues {
			t.Fatalf("expected ExitDoctorIssues, got %d", code)
		}
	})

	if !strings.Contains(out, "Suggested repair path") || !strings.Contains(out, "weave skill add sdd-orchestrator") {
		t.Fatalf("expected repair path output, got: %s", out)
	}

	if !strings.Contains(out, "Docs: docs/reference/doctor.md") {
		t.Fatalf("expected docs reference output, got: %s", out)
	}
}

func TestRunDoctor_JSONOutput(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n"), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	t.Setenv("WEAVE_WORKDIR", repo)
	out := captureStdout(t, func() {
		code, err := runDoctor(context.Background(), []string{"--json"})
		if err != nil {
			t.Fatalf("unexpected runDoctor error: %v", err)
		}
		if code != ExitOK {
			t.Fatalf("expected ExitOK for healthy json output, got %d", code)
		}
	})

	var payload struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("expected valid json output, got error: %v\noutput: %s", err, out)
	}

	if payload.Status != "healthy" {
		t.Fatalf("expected healthy json status, got %q", payload.Status)
	}
}

func TestRunDoctor_StaleProviderIntegrationIsReported(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	configBody := "version: 1\n" +
		"providers:\n" +
		"  - name: claude-code\n" +
		"    enabled: true\n" +
		"sources:\n" +
		"  skills_dir: ~/.weave/skills\n" +
		"  commands_dir: ~/.weave/commands\n" +
		"sync:\n" +
		"  mode: symlink\n" +
		"skills: []\n" +
		"commands: []\n"

	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte(configBody), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	t.Setenv("WEAVE_WORKDIR", repo)
	out := captureStdout(t, func() {
		code, err := runDoctor(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected runDoctor error: %v", err)
		}
		if code != ExitDoctorIssues {
			t.Fatalf("expected ExitDoctorIssues for stale provider integration, got %d", code)
		}
	})

	if !strings.Contains(out, "stale_provider_integration") {
		t.Fatalf("expected stale provider issue in doctor output, got: %s", out)
	}
	if !strings.Contains(out, "weave provider repair claude-code") {
		t.Fatalf("expected provider repair guidance, got: %s", out)
	}
}

func TestRunDoctor_RejectsUnknownFlagDeterministically(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)

	code, err := runDoctor(context.Background(), []string{"--yaml"})
	if err == nil {
		t.Fatalf("expected unknown flag error")
	}
	if code != ExitRuntimeError {
		t.Fatalf("expected ExitRuntimeError for unsupported flag, got %d", code)
	}
}

func captureStdout(t *testing.T, run func()) string {
	t.Helper()

	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = w

	run()

	if err := w.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	os.Stdout = orig

	var b bytes.Buffer
	if _, err := io.Copy(&b, r); err != nil {
		t.Fatalf("read stdout buffer: %v", err)
	}
	_ = r.Close()

	return b.String()
}
