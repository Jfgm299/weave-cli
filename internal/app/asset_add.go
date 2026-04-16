package app

import (
	"context"
	"fmt"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type AssetKind string

const (
	AssetKindSkill   AssetKind = "skill"
	AssetKindCommand AssetKind = "command"
)

type AddAssetInput struct {
	Kind        AssetKind
	Name        string
	SourcePath  string
	ProjectPath string
}

type AddAssetResult struct {
	Root        string
	ConfigSaved bool
	OpsPlanned  int
	OpsApplied  int
	DryRun      bool
}

func (s ForgeService) AddAsset(ctx context.Context, cfg config.Config, input AddAssetInput) (AddAssetResult, error) {
	return s.AddAssetWithOptions(ctx, cfg, input, RunOptions{})
}

func (s ForgeService) AddAssetWithOptions(ctx context.Context, cfg config.Config, input AddAssetInput, opts RunOptions) (AddAssetResult, error) {
	root, err := s.ProjectRootDetector.Detect(ctx)
	if err != nil {
		return AddAssetResult{}, ErrNotInProjectRoot
	}

	if err := s.ConfigValidator.Validate(cfg); err != nil {
		return AddAssetResult{}, WrapInvalidConfig(err)
	}

	op := fsops.Operation{
		Type:   fsops.OpCreateLink,
		Path:   input.ProjectPath,
		Target: input.SourcePath,
	}

	if err := ensureOpsWithinRoot(root, []fsops.Operation{op}); err != nil {
		return AddAssetResult{}, err
	}

	result := AddAssetResult{Root: root, OpsPlanned: 1, DryRun: opts.DryRun}

	if opts.DryRun {
		return result, nil
	}

	if err := s.Executor.Apply(ctx, []fsops.Operation{op}); err != nil {
		return AddAssetResult{}, err
	}
	result.OpsApplied = 1

	nextCfg := cfg
	switch input.Kind {
	case AssetKindSkill:
		nextCfg.Skills = upsertAsset(nextCfg.Skills, config.Asset{Name: input.Name, Source: input.SourcePath})
	case AssetKindCommand:
		nextCfg.Commands = upsertAsset(nextCfg.Commands, config.Asset{Name: input.Name, Source: input.SourcePath})
	}

	if err := s.Writer.Write(nextCfg); err != nil {
		rollbackOp := fsops.Operation{
			Type: fsops.OpRemovePath,
			Path: input.ProjectPath,
		}
		if guardErr := ensureOpsWithinRoot(root, []fsops.Operation{rollbackOp}); guardErr != nil {
			return AddAssetResult{}, guardErr
		}

		rollbackErr := s.Executor.Apply(ctx, []fsops.Operation{{
			Type: fsops.OpRemovePath,
			Path: input.ProjectPath,
		}})
		if rollbackErr != nil {
			return AddAssetResult{}, fmt.Errorf("failed to persist weave.yaml after symlink apply; rollback failed so project may be partially modified: %w; rollback failed: %v. Run `weave doctor` and then the suggested repair command. See %s (%s)", err, rollbackErr, DocsPathTransactions, DocsURL(DocsPathTransactions))
		}
		return AddAssetResult{}, fmt.Errorf("failed to persist weave.yaml after symlink apply; rollback completed so no config or symlink changes were committed: %w. Re-run the command after fixing the config write issue. See %s (%s)", err, DocsPathTransactions, DocsURL(DocsPathTransactions))
	}

	result.ConfigSaved = true
	return result, nil
}

func upsertAsset(in []config.Asset, asset config.Asset) []config.Asset {
	for i := range in {
		if in[i].Name == asset.Name {
			in[i] = asset
			return in
		}
	}

	out := make([]config.Asset, 0, len(in)+1)
	out = append(out, in...)
	return append(out, asset)
}
