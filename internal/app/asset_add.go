package app

import (
	"context"

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
}

func (s ForgeService) AddAsset(ctx context.Context, cfg config.Config, input AddAssetInput) (AddAssetResult, error) {
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

	if err := s.Executor.Apply(ctx, []fsops.Operation{op}); err != nil {
		return AddAssetResult{}, err
	}

	nextCfg := cfg
	switch input.Kind {
	case AssetKindSkill:
		nextCfg.Skills = upsertAsset(nextCfg.Skills, config.Asset{Name: input.Name, Source: input.SourcePath})
	case AssetKindCommand:
		nextCfg.Commands = upsertAsset(nextCfg.Commands, config.Asset{Name: input.Name, Source: input.SourcePath})
	}

	if err := s.Writer.Write(nextCfg); err != nil {
		return AddAssetResult{}, err
	}

	return AddAssetResult{Root: root, ConfigSaved: true}, nil
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
