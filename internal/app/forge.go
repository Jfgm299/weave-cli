package app

import (
	"context"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type ForgeService struct {
	ProjectRootDetector ProjectRootDetector
	ConfigValidator     ConfigValidator
	Planner             ForgePlanner
	Executor            OperationExecutor
	Writer              ConfigWriter
}

type ProjectRootDetector interface {
	Detect(ctx context.Context) (string, error)
}

type ConfigValidator interface {
	Validate(cfg config.Config) error
}

type ForgePlanner interface {
	Plan(cfg config.Config) ([]fsops.Operation, error)
}

type OperationExecutor interface {
	Apply(ctx context.Context, ops []fsops.Operation) error
}

type ConfigWriter interface {
	Write(cfg config.Config) error
}

type ForgeResult struct {
	Root        string
	OpsPlanned  int
	OpsApplied  int
	WasNoOp     bool
	ConfigSaved bool
	DryRun      bool
}

func (s ForgeService) Run(ctx context.Context, cfg config.Config) (ForgeResult, error) {
	return s.RunWithOptions(ctx, cfg, RunOptions{})
}

type RunOptions struct {
	DryRun bool
}

func (s ForgeService) RunWithOptions(ctx context.Context, cfg config.Config, opts RunOptions) (ForgeResult, error) {
	root, err := s.ProjectRootDetector.Detect(ctx)
	if err != nil {
		return ForgeResult{}, ErrNotInProjectRoot
	}

	if err := s.ConfigValidator.Validate(cfg); err != nil {
		return ForgeResult{}, WrapInvalidConfig(err)
	}

	ops, err := s.Planner.Plan(cfg)
	if err != nil {
		return ForgeResult{}, err
	}

	if err := ensureOpsWithinRoot(root, ops); err != nil {
		return ForgeResult{}, err
	}

	result := ForgeResult{Root: root, OpsPlanned: len(ops), WasNoOp: len(ops) == 0, DryRun: opts.DryRun}

	if opts.DryRun {
		return result, nil
	}

	if len(ops) > 0 {
		if err := s.Executor.Apply(ctx, ops); err != nil {
			return ForgeResult{}, err
		}
		result.OpsApplied = len(ops)
	}

	if err := s.Writer.Write(cfg); err != nil {
		return ForgeResult{}, err
	}
	result.ConfigSaved = true

	return result, nil
}
