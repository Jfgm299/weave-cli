package app

import (
	"context"
	"errors"
	"strings"
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

func TestForgeService_Run_DryRunPlansWithoutApplyingOrPersisting(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Planner:             plannerStub{ops: []fsops.Operation{{Type: fsops.OpEnsureDir, Path: "/tmp/proj/.agents/skills"}}},
		Executor:            executor,
		Writer:              writer,
	}

	result, err := sut.RunWithOptions(context.Background(), config.Default(), RunOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected dry-run result")
	}

	if executor.called {
		t.Fatalf("executor should not run during dry-run")
	}

	if writer.called {
		t.Fatalf("writer should not run during dry-run")
	}
}

func TestForgeService_Run_RejectsOpsOutsideDetectedRoot(t *testing.T) {
	t.Parallel()

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Planner:             plannerStub{ops: []fsops.Operation{{Type: fsops.OpEnsureDir, Path: "/tmp/other/.agents"}}},
		Executor:            &executorSpy{},
		Writer:              &writerSpy{},
	}

	_, err := sut.Run(context.Background(), config.Default())
	if err == nil {
		t.Fatalf("expected unsafe path guard error")
	}

	if !errors.Is(err, ErrUnsafeMutationPath) {
		t.Fatalf("expected ErrUnsafeMutationPath, got %v", err)
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

func TestForgeService_AddAsset_ConfigWriteFailureRollsBackSymlink(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{err: errors.New("disk full")}

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
		SourcePath:  "/tmp/proj/src/skills/sdd-orchestrator/SKILL.md",
		ProjectPath: "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md",
	})
	if err == nil {
		t.Fatalf("expected config write failure")
	}

	if len(executor.calls) != 2 {
		t.Fatalf("expected two executor calls (create link + rollback), got %d", len(executor.calls))
	}

	if len(executor.calls[1]) != 1 || executor.calls[1][0].Type != fsops.OpRemovePath {
		t.Fatalf("expected rollback remove_path call, got %+v", executor.calls[1])
	}

	if executor.calls[1][0].Path != "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md" {
		t.Fatalf("unexpected rollback path: %s", executor.calls[1][0].Path)
	}

	if got := err.Error(); got == "" || !contains(got, "rollback completed so no config or symlink changes were committed") {
		t.Fatalf("expected rollback semantics in error output, got %q", got)
	}
}

func TestForgeService_AddAsset_DryRunDoesNotApplyOrPersist(t *testing.T) {
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

	result, err := sut.AddAssetWithOptions(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindSkill,
		Name:        "sdd-orchestrator",
		SourcePath:  "/tmp/proj/shared/skills/sdd-orchestrator/SKILL.md",
		ProjectPath: "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md",
	}, RunOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected dry-run result")
	}

	if executor.called {
		t.Fatalf("executor should not run during dry-run")
	}

	if writer.called {
		t.Fatalf("writer should not run during dry-run")
	}
}

func TestForgeService_AddAsset_RejectsProjectPathOutsideRoot(t *testing.T) {
	t.Parallel()

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Planner:             plannerStub{},
		Executor:            &executorSpy{},
		Writer:              &writerSpy{},
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  "/tmp/proj/shared/commands/pr-review.md",
		ProjectPath: "/tmp/other/pr-review.md",
	})
	if err == nil {
		t.Fatalf("expected unsafe path guard error")
	}

	if !errors.Is(err, ErrUnsafeMutationPath) {
		t.Fatalf("expected ErrUnsafeMutationPath, got %v", err)
	}
}

func contains(s string, token string) bool {
	return strings.Contains(s, token)
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
			SourcePath:  "/tmp/proj/src/skills/sdd-orchestrator/SKILL.md",
			ProjectPath: "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md",
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
			SourcePath:  "/tmp/proj/src/commands/pr-review.md",
			ProjectPath: "/tmp/proj/.agents/commands/pr-review.md",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(executor.lastOps) != 1 || executor.lastOps[0].Type != fsops.OpCreateLink {
			t.Fatalf("expected create_link operation for command add")
		}
	})
}

func TestForgeService_AddAsset_ConflictPromptRequiredWhenNoPrompter(t *testing.T) {
	t.Parallel()

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            &executorSpy{},
		Writer:              &writerSpy{},
		PathChecker:         pathCheckerStub{exists: true},
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindSkill,
		Name:        "sdd-orchestrator",
		SourcePath:  "/tmp/proj/src/skills/sdd-orchestrator/SKILL.md",
		ProjectPath: "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md",
	})
	if err == nil {
		t.Fatalf("expected prompt-required conflict error")
	}
	if !errors.Is(err, ErrConflictPromptRequired) {
		t.Fatalf("expected ErrConflictPromptRequired, got %v", err)
	}
}

func TestForgeService_AddAsset_ConflictBackupPolicyPlansBackupAndCreateLink(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
		PathChecker:         pathCheckerStub{exists: true},
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:           AssetKindSkill,
		Name:           "sdd-orchestrator",
		SourcePath:     "/tmp/proj/src/skills/sdd-orchestrator/SKILL.md",
		ProjectPath:    "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md",
		ConflictPolicy: ConflictPolicyBackup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(executor.lastOps) != 2 {
		t.Fatalf("expected backup + create link ops, got %+v", executor.lastOps)
	}
	if executor.lastOps[0].Type != fsops.OpBackupPath || executor.lastOps[1].Type != fsops.OpCreateLink {
		t.Fatalf("unexpected conflict operation order: %+v", executor.lastOps)
	}
}

func TestForgeService_AddAsset_ConflictSkipPolicyReturnsWithoutWrites(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
		PathChecker:         pathCheckerStub{exists: true},
	}

	result, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:           AssetKindSkill,
		Name:           "sdd-orchestrator",
		SourcePath:     "/tmp/proj/src/skills/sdd-orchestrator/SKILL.md",
		ProjectPath:    "/tmp/proj/.agents/skills/sdd-orchestrator/SKILL.md",
		ConflictPolicy: ConflictPolicySkip,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OpsPlanned != 0 || result.OpsApplied != 0 {
		t.Fatalf("expected no operations for skip policy, got %+v", result)
	}
	if executor.called {
		t.Fatalf("executor should not run for skip policy")
	}
	if writer.called {
		t.Fatalf("writer should not run for skip policy")
	}
}

func TestForgeService_AddAsset_WithAdditionalOperations_AppliesAllAndPersistsConfig(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  "/tmp/proj/src/commands/pr-review.md",
		ProjectPath: "/tmp/proj/.agents/commands/pr-review.md",
		AdditionalOps: []fsops.Operation{{
			Type:   fsops.OpCreateLink,
			Path:   "/tmp/proj/.codex/commands/__weave_commands__/pr-review/SKILL.md",
			Target: "/tmp/proj/src/commands/pr-review.md",
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(executor.lastOps) != 2 {
		t.Fatalf("expected shared + provider projection ops, got %+v", executor.lastOps)
	}
	if !writer.called {
		t.Fatalf("expected config writer called")
	}
}

func TestForgeService_AddAsset_ConfigWriteFailureRollsBackAdditionalOperations(t *testing.T) {
	t.Parallel()

	executor := &executorSpy{}
	writer := &writerSpy{err: errors.New("disk full")}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  "/tmp/proj/src/commands/pr-review.md",
		ProjectPath: "/tmp/proj/.agents/commands/pr-review.md",
		AdditionalOps: []fsops.Operation{{
			Type:   fsops.OpCreateLink,
			Path:   "/tmp/proj/.codex/commands/__weave_commands__/pr-review/SKILL.md",
			Target: "/tmp/proj/src/commands/pr-review.md",
		}},
	})
	if err == nil {
		t.Fatalf("expected config write failure")
	}

	if len(executor.calls) != 2 {
		t.Fatalf("expected apply + rollback calls, got %d", len(executor.calls))
	}
	if len(executor.calls[1]) != 2 {
		t.Fatalf("expected rollback for both command targets, got %+v", executor.calls[1])
	}
}

func TestForgeService_AddAsset_ApplyFailureRollsBackAllPlannedOperations(t *testing.T) {
	t.Parallel()

	executor := &failFirstApplyExecutor{}
	writer := &writerSpy{}

	sut := ForgeService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
	}

	_, err := sut.AddAsset(context.Background(), config.Default(), AddAssetInput{
		Kind:        AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  "/tmp/proj/src/commands/pr-review.md",
		ProjectPath: "/tmp/proj/.agents/commands/pr-review.md",
		AdditionalOps: []fsops.Operation{{
			Type:   fsops.OpCreateLink,
			Path:   "/tmp/proj/.codex/commands/__weave_commands__/pr-review/SKILL.md",
			Target: "/tmp/proj/src/commands/pr-review.md",
		}},
	})
	if err == nil {
		t.Fatalf("expected apply failure")
	}
	if !strings.Contains(err.Error(), "rollback completed so no config or symlink changes were committed") {
		t.Fatalf("expected full rollback semantics in error, got %q", err.Error())
	}
	if len(executor.calls) != 2 {
		t.Fatalf("expected apply + rollback calls, got %d", len(executor.calls))
	}
	if len(executor.calls[1]) != 2 {
		t.Fatalf("expected rollback for both shared and provider targets, got %+v", executor.calls[1])
	}
	if !allRemovePaths(executor.calls[1]) {
		t.Fatalf("expected rollback to contain remove_path ops only, got %+v", executor.calls[1])
	}
	if writer.called {
		t.Fatalf("config writer must not run when apply failed")
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
	called  bool
	err     error
	lastOps []fsops.Operation
	calls   [][]fsops.Operation
}

func (e *executorSpy) Apply(_ context.Context, ops []fsops.Operation) error {
	e.called = true
	e.lastOps = append([]fsops.Operation(nil), ops...)
	e.calls = append(e.calls, append([]fsops.Operation(nil), ops...))
	return e.err
}

type writerSpy struct {
	called bool
	err    error
}

type failFirstApplyExecutor struct {
	calls [][]fsops.Operation
	seen  bool
}

func (e *failFirstApplyExecutor) Apply(_ context.Context, ops []fsops.Operation) error {
	e.calls = append(e.calls, append([]fsops.Operation(nil), ops...))
	if !e.seen {
		e.seen = true
		return errors.New("apply failed")
	}
	return nil
}

type pathCheckerStub struct {
	exists bool
	err    error
}

func (p pathCheckerStub) Exists(_ string) (bool, error) {
	if p.err != nil {
		return false, p.err
	}
	return p.exists, nil
}

func (w *writerSpy) Write(_ config.Config) error {
	w.called = true
	return w.err
}

func allRemovePaths(ops []fsops.Operation) bool {
	for _, op := range ops {
		if op.Type != fsops.OpRemovePath {
			return false
		}
	}
	return true
}
