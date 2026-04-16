package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestForgeHandler_Run_ReturnsInvalidConfigExitCode(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	t.Setenv("WEAVE_WORKDIR", tmp)
	if err := os.WriteFile(filepath.Join(tmp, "weave.yaml"), []byte("version: 1\nsync:\n  mode: symlink\nskills: []\ncommands: []\n"), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	h := ForgeHandler{Service: app.ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp"},
		ConfigValidator:     validatorStub{err: errors.New("bad")},
		Planner:             plannerStub{},
		Executor:            &executorSpy{},
		Writer:              &writerSpy{},
	}}

	code, err := h.Run(context.Background())
	if code != ExitInvalidConfig {
		t.Fatalf("expected ExitInvalidConfig, got %d", code)
	}
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProjectRootDetector_FailsWhenNoGitAncestorExists(t *testing.T) {
	workdir := t.TempDir()
	promptCalled := false
	runGitInitCalled := false
	t.Cleanup(func() {
		isInteractiveSession = defaultIsInteractiveSession
		promptGitInit = defaultPromptGitInit
		runGitInit = defaultRunGitInit
	})
	isInteractiveSession = func() bool { return false }
	promptGitInit = func(string) (bool, error) {
		promptCalled = true
		return false, nil
	}
	runGitInit = func(string) error {
		runGitInitCalled = true
		return nil
	}
	detector := projectRootDetector{Workdir: workdir}

	_, err := detector.Detect(context.Background())
	if err == nil {
		t.Fatalf("expected root detection failure without .git")
	}

	if !strings.Contains(err.Error(), "project root not detected") {
		t.Fatalf("expected actionable root detection error, got %q", err.Error())
	}

	if !strings.Contains(err.Error(), "Run `git init` and retry") {
		t.Fatalf("expected git init guidance in error, got %q", err.Error())
	}

	if promptCalled {
		t.Fatalf("promptGitInit should not be called in non-interactive mode")
	}

	if runGitInitCalled {
		t.Fatalf("runGitInit should not be called in non-interactive mode")
	}
}

func TestDefaultIsInteractiveSession_RespectsNonInteractiveEnvOverride(t *testing.T) {
	t.Setenv("WEAVE_NON_INTERACTIVE", "1")

	if defaultIsInteractiveSession() {
		t.Fatalf("expected non-interactive override to force false")
	}
}

func TestProjectRootDetector_InteractivePromptDeclinedFails(t *testing.T) {
	workdir := t.TempDir()
	promptCalled := false
	t.Cleanup(func() {
		isInteractiveSession = defaultIsInteractiveSession
		promptGitInit = defaultPromptGitInit
		runGitInit = defaultRunGitInit
	})
	isInteractiveSession = func() bool { return true }
	promptGitInit = func(string) (bool, error) {
		promptCalled = true
		return false, nil
	}
	runGitInit = func(string) error {
		t.Fatalf("runGitInit should not be called when prompt is declined")
		return nil
	}

	detector := projectRootDetector{Workdir: workdir}
	_, err := detector.Detect(context.Background())
	if err == nil {
		t.Fatalf("expected detect to fail when git init is declined")
	}

	if !strings.Contains(err.Error(), "initialization was declined") {
		t.Fatalf("expected declined message, got %q", err.Error())
	}

	if !promptCalled {
		t.Fatalf("expected interactive flow to call promptGitInit")
	}
}

func TestProjectRootDetector_InteractivePromptAcceptedRunsGitInit(t *testing.T) {
	workdir := t.TempDir()
	t.Cleanup(func() {
		isInteractiveSession = defaultIsInteractiveSession
		promptGitInit = defaultPromptGitInit
		runGitInit = defaultRunGitInit
	})
	isInteractiveSession = func() bool { return true }
	promptGitInit = func(string) (bool, error) { return true, nil }
	runGitInit = func(dir string) error {
		if dir != workdir {
			t.Fatalf("expected git init dir %q, got %q", workdir, dir)
		}
		if err := os.Mkdir(filepath.Join(workdir, ".git"), 0o755); err != nil {
			return err
		}
		return nil
	}

	detector := projectRootDetector{Workdir: workdir}
	root, err := detector.Detect(context.Background())
	if err != nil {
		t.Fatalf("expected detect to succeed after git init, got %v", err)
	}

	if root != workdir {
		t.Fatalf("expected root %q, got %q", workdir, root)
	}
}

func TestProjectRootDetector_InteractivePromptAcceptedGitInitFails(t *testing.T) {
	workdir := t.TempDir()
	t.Cleanup(func() {
		isInteractiveSession = defaultIsInteractiveSession
		promptGitInit = defaultPromptGitInit
		runGitInit = defaultRunGitInit
	})
	isInteractiveSession = func() bool { return true }
	promptGitInit = func(string) (bool, error) { return true, nil }
	runGitInit = func(string) error { return errors.New("git init failed") }

	detector := projectRootDetector{Workdir: workdir}
	_, err := detector.Detect(context.Background())
	if err == nil {
		t.Fatalf("expected detect to fail when git init fails")
	}

	if !strings.Contains(err.Error(), "failed to run `git init`") {
		t.Fatalf("expected git init failure guidance, got %q", err.Error())
	}
}

func TestForgeHandler_Run_ReturnsExitOKWhenSuccessful(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	t.Setenv("WEAVE_WORKDIR", tmp)

	h := NewDefaultForgeHandler()

	code, err := h.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if code != ExitOK {
		t.Fatalf("expected ExitOK, got %d", code)
	}

	if _, err := os.Stat(filepath.Join(tmp, "weave.yaml")); err != nil {
		t.Fatalf("expected weave.yaml written: %v", err)
	}
}

func TestProjectRootDetector_DetectsRepositoryRootFromNestedWorkdir(t *testing.T) {
	repo := t.TempDir()
	nested := filepath.Join(repo, "a", "b", "c")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	detector := projectRootDetector{Workdir: nested}
	root, err := detector.Detect(context.Background())
	if err != nil {
		t.Fatalf("detect root: %v", err)
	}

	if root != repo {
		t.Fatalf("expected detected root %q, got %q", repo, root)
	}
}

func TestForgeHandler_RunWithOptions_DryRunDoesNotMutate(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	t.Setenv("WEAVE_WORKDIR", repo)

	h := NewDefaultForgeHandler()
	code, result, err := h.RunWithOptions(context.Background(), true)
	if err != nil {
		t.Fatalf("dry-run forge failed: %v", err)
	}
	if code != ExitOK {
		t.Fatalf("expected ExitOK, got %d", code)
	}

	if !result.ServiceResult.DryRun {
		t.Fatalf("expected dry-run result")
	}

	if result.ServiceResult.OpsApplied != 0 {
		t.Fatalf("expected no applied ops in dry-run, got %d", result.ServiceResult.OpsApplied)
	}

	if _, err := os.Stat(filepath.Join(repo, ".agents", "skills")); !os.IsNotExist(err) {
		t.Fatalf("expected no .agents/skills created on dry-run, got err: %v", err)
	}

	if _, err := os.Stat(filepath.Join(repo, "weave.yaml")); !os.IsNotExist(err) {
		t.Fatalf("expected no weave.yaml write on dry-run, got err: %v", err)
	}
}

type detectorStub struct {
	root string
	err  error
}

func (d detectorStub) Detect(_ context.Context) (string, error) {
	if d.err != nil {
		return "", d.err
	}
	return d.root, nil
}

type validatorStub struct {
	err error
}

func (v validatorStub) Validate(_ config.Config) error {
	return v.err
}

type plannerStub struct {
	ops []fsops.Operation
	err error
}

func (p plannerStub) Plan(_ config.Config) ([]fsops.Operation, error) {
	return p.ops, p.err
}

type executorSpy struct {
	called bool
	err    error
}

func (e *executorSpy) Apply(_ context.Context, _ []fsops.Operation) error {
	e.called = true
	return e.err
}

type writerSpy struct {
	called bool
	err    error
}

func (w *writerSpy) Write(_ config.Config) error {
	w.called = true
	return w.err
}
