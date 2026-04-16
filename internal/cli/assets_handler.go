package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type assetAddService struct {
	Service  app.ForgeService
	Resolver sourceResolver
	Workdir  string
}

type stdinConflictPrompter struct {
	in *bufio.Reader
}

func (p stdinConflictPrompter) ResolveConflict(_ context.Context, input app.ConflictPromptInput) (app.ConflictPolicy, error) {
	reader := p.in
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}

	fmt.Printf("Conflict detected for %s %q at %s. Choose action [overwrite/skip/backup]: ", input.Kind, input.Name, input.Path)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	choice := strings.TrimSpace(strings.ToLower(line))
	return app.ConflictPolicy(choice), nil
}

func (s assetAddService) Add(ctx context.Context, kind assetKind, name string, fromFlag string, dryRun bool, conflictPolicy app.ConflictPolicy, cfg config.Config) (app.AddAssetResult, error) {
	return s.AddWithPlan(ctx, kind, name, fromFlag, dryRun, conflictPolicy, cfg, nil)
}

func (s assetAddService) AddWithPlan(ctx context.Context, kind assetKind, name string, fromFlag string, dryRun bool, conflictPolicy app.ConflictPolicy, cfg config.Config, plan *app.CommandInstallPlan) (app.AddAssetResult, error) {
	root := s.Workdir
	if root == "" {
		root = resolveWorkdir()
	}

	fromBase, err := s.Resolver.Resolve(kind, fromFlag, cfg)
	if err != nil {
		return app.AddAssetResult{}, err
	}

	input, err := buildAddAssetInput(kind, root, name, fromBase, conflictPolicy, cfg, plan)
	if err != nil {
		return app.AddAssetResult{}, err
	}

	return s.Service.AddAssetWithOptions(ctx, cfg, input, app.RunOptions{DryRun: dryRun})
}

func buildAddAssetInput(kind assetKind, root string, name string, fromBase string, conflictPolicy app.ConflictPolicy, cfg config.Config, plan *app.CommandInstallPlan) (app.AddAssetInput, error) {
	input := app.AddAssetInput{
		Kind:           mapKind(kind),
		Name:           name,
		SourcePath:     sourcePathFor(kind, fromBase, name),
		ProjectPath:    assetPathFor(kind, root, name),
		ConflictPolicy: conflictPolicy,
	}

	if kind != assetKindCommand {
		return input, nil
	}

	sharedInstall := true
	targets := enabledProviderNames(cfg)
	if plan != nil {
		sharedInstall = plan.SharedInstall
		targets = append([]string(nil), plan.ProviderTargets...)
	}

	if !sharedInstall {
		input.ProjectPath = ""
	}

	if len(targets) == 0 {
		return input, nil
	}

	metaShared := sharedInstall
	input.CommandMeta = &config.CommandMetaV1{ProviderCompat: targets, SharedInstall: &metaShared}
	projectionOps, err := app.BuildCommandProjectionOps(root, name, input.SourcePath, app.CommandInstallPlan{ProviderTargets: targets, SharedInstall: sharedInstall})
	if err != nil {
		return app.AddAssetInput{}, err
	}
	input.AdditionalOps = projectionOps

	return input, nil
}

func enabledProviderNames(cfg config.Config) []string {
	enabledProviders := app.ListEnabledProviders(cfg)
	names := make([]string, 0, len(enabledProviders))
	for _, provider := range enabledProviders {
		names = append(names, provider.Name)
	}
	return names
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
			PathChecker:         appPathChecker{},
			ConflictPrompter:    stdinConflictPrompter{},
		},
		Resolver: newSourceResolver(),
		Workdir:  workdir,
	}
}

type appPathChecker struct{}

func (appPathChecker) Exists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
