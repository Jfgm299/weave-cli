package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type assetAddService struct {
	Service  app.ForgeService
	Resolver sourceResolver
	Workdir  string
}

func (s assetAddService) Add(ctx context.Context, kind assetKind, name string, fromFlag string, cfg config.Config) (app.AddAssetResult, error) {
	root := s.Workdir
	if root == "" {
		root = resolveWorkdir()
	}

	fromBase, err := s.Resolver.Resolve(kind, fromFlag, cfg)
	if err != nil {
		return app.AddAssetResult{}, err
	}

	input := app.AddAssetInput{
		Kind:        mapKind(kind),
		Name:        name,
		SourcePath:  sourcePathFor(kind, fromBase, name),
		ProjectPath: assetPathFor(kind, root, name),
	}

	return s.Service.AddAsset(ctx, cfg, input)
}

func mapKind(kind assetKind) app.AssetKind {
	if kind == assetKindCommand {
		return app.AssetKindCommand
	}
	return app.AssetKindSkill
}

func assetPathFor(kind assetKind, repo string, name string) string {
	if kind == assetKindCommand {
		return filepath.Join(repo, ".agents", "commands", fmt.Sprintf("%s.md", name))
	}
	return filepath.Join(repo, ".agents", "skills", name, "SKILL.md")
}

func sourcePathFor(kind assetKind, fromBase string, name string) string {
	if kind == assetKindCommand {
		return filepath.Join(fromBase, fmt.Sprintf("%s.md", name))
	}
	return filepath.Join(fromBase, name, "SKILL.md")
}

func newDefaultAssetAddService() assetAddService {
	workdir := resolveWorkdir()
	validator := config.Validator{}

	return assetAddService{
		Service: app.ForgeService{
			ProjectRootDetector: projectRootDetector{Workdir: workdir},
			ConfigValidator:     validator,
			Planner:             forgePlanner{Workdir: workdir},
			Executor:            fsops.Engine{},
			Writer:              config.AtomicFileWriter{Path: filepath.Join(workdir, "weave.yaml")},
		},
		Resolver: newSourceResolver(),
		Workdir:  workdir,
	}
}
