package app

import (
	"context"
	"errors"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestForgeService_Run_ConvergedProjectReturnsNoOp(t *testing.T) {
	t.Parallel()

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Planner:             plannerStub{ops: nil},
		Executor:            &executorSpy{},
		Writer:              &writerSpy{},
	}

	result, err := sut.Run(context.Background(), config.Config{Version: 1, Sync: config.Sync{Mode: "symlink"}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !result.WasNoOp {
		t.Fatalf("expected no-op result")
	}

	if result.OpsApplied != 0 {
		t.Fatalf("expected 0 applied ops, got %d", result.OpsApplied)
	}

	if !result.ConfigSaved {
		t.Fatalf("expected config to be saved")
	}
}

func TestForgeService_Run_InvalidConfigShortCircuitsBeforeMutation(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{err: errors.New("bad config")},
		Planner:             plannerStub{ops: []fsops.Operation{{Type: fsops.OpEnsureDir, Path: ".agents"}}},
		Executor:            executor,
		Writer:              writer,
	}

	_, err := sut.Run(context.Background(), config.Config{})
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}

	if executor.called {
		t.Fatalf("expected executor not to be called on invalid config")
	}

	if writer.called {
		t.Fatalf("expected writer not to be called on invalid config")
	}
}

func TestForgeService_AddAsset_StrictModeDoesNotPersistOnSymlinkFailure(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{err: errors.New("symlink failed")}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Planner:             plannerStub{},
		Executor:            executor,
		Writer:              writer,
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindSkill,
		Name:        "sdd-orchestrator",
		SourcePath:  "/src/skill/SKILL.md",
		ProjectPath: "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md",
	})
	if err == nil {
		t.Fatalf("expected add asset to fail")
	}

	if writer.called {
		t.Fatalf("expected config writer not called when symlink fails")
	}
}

func TestForgeService_AddAsset_StrictModeDoesNotPersistOnSymlinkFailureForCommand(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{err: errors.New("symlink failed")}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Planner:             plannerStub{},
		Executor:            executor,
		Writer:              writer,
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  "/src/command/pr-review.md",
		ProjectPath: "/tmp/proj/.agents/commands/pr-review.md",
	})
	if err == nil {
		t.Fatalf("expected add asset to fail")
	}

	if writer.called {
		t.Fatalf("expected config writer not called when command symlink fails")
	}
}

func TestForgeService_AddAsset_PersistsConfigAfterSymlinkSuccess(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Planner:             plannerStub{},
		Executor:            executor,
		Writer:              writer,
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  "/src/cmd/pr-review.md",
		ProjectPath: "/tmp/proj/.agents/commands/pr-review.md",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !writer.called {
		t.Fatalf("expected config writer called after fs operation success")
	}
}

func TestForgeService_AddAsset_UsesCreateLinkForSkillAndCommand(t *testing.T) {
	t.Parallel()

	t.Run("skill", func(t *testing.T) {
		t.Parallel()

		executor := &executorSpy{}
		writer := &writerSpy{}

		sut := ForgeService{
			ProjectRootDetector: detectorStub{root: "/tmp/proj"},
			ConfigValidator:     validatorStub{},
			Planner:             plannerStub{},
			Executor:            executor,
			Writer:              writer,
		}

		_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
			Kind:        AssetKindSkill,
			Name:        "sdd-orchestrator",
			SourcePath:  "/src/skills/sdd-orchestrator/SKILL.md",
			ProjectPath: "/repo/.agents/skills/sdd-orchestrator/SKILL.md",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(executor.lastOps) != 1 || executor.lastOps[0].Type != fsops.OpCreateLink {
			t.Fatalf("expected create_link operation for skill add")
		}
	})

	t.Run("command", func(t *testing.T) {
		t.Parallel()

		executor := &executorSpy{}
		writer := &writerSpy{}

		sut := ForgeService{
			ProjectRootDetector: detectorStub{root: "/tmp/proj"},
			ConfigValidator:     validatorStub{},
			Planner:             plannerStub{},
			Executor:            executor,
			Writer:              writer,
		}

		_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
			Kind:        AssetKindCommand,
			Name:        "pr-review",
			SourcePath:  "/src/commands/pr-review.md",
			ProjectPath: "/repo/.agents/commands/pr-review.md",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(executor.lastOps) != 1 || executor.lastOps[0].Type != fsops.OpCreateLink {
			t.Fatalf("expected create_link operation for command add")
		}
	})
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
	called  bool
	err     error
	lastOps []fsops.Operation
}

func (e *executorSpy) Apply(_ context.Context, ops []fsops.Operation) error {
	e.called = true
	e.lastOps = ops
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
