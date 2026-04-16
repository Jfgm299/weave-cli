package cli

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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
