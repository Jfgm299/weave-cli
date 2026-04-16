package cli

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
	"github.com/Jfgm299/weave-cli/internal/providers"
)

type providerAction string

const (
	providerActionAdd    providerAction = "add"
	providerActionList   providerAction = "list"
	providerActionRemove providerAction = "remove"
	providerActionRepair providerAction = "repair"
)

type providerBinaryResolverFunc func(string) (string, error)

func (f providerBinaryResolverFunc) LookPath(file string) (string, error) {
	return f(file)
}

func newProviderService(workdir string, resolver app.BinaryResolver) app.ProviderService {
	if resolver == nil {
		resolver = providerBinaryResolverFunc(exec.LookPath)
	}

	validator := config.Validator{}
	return app.ProviderService{
		ProjectRootDetector: projectRootDetector{Workdir: workdir},
		ConfigValidator:     validator,
		Executor:            fsops.Engine{},
		Writer:              config.AtomicFileWriter{Path: filepath.Join(workdir, "weave.yaml")},
		BinaryResolver:      resolver,
	}
}

func runProviderAction(ctx context.Context, svc app.ProviderService, registry providers.Registry, action providerAction, name string, dryRun bool) (app.ProviderAddResult, []string, error) {
	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		return app.ProviderAddResult{}, nil, err
	}

	switch action {
	case providerActionAdd:
		if name == "" {
			return app.ProviderAddResult{}, nil, fmt.Errorf("provider name is required")
		}
		res, err := svc.AddProviderWithOptions(ctx, cfg, registry, name, app.RunOptions{DryRun: dryRun})
		return res, nil, err
	case providerActionRemove:
		if name == "" {
			return app.ProviderAddResult{}, nil, fmt.Errorf("provider name is required")
		}
		res, err := svc.RemoveProviderWithOptions(ctx, cfg, registry, name, app.RunOptions{DryRun: dryRun})
		return res, nil, err
	case providerActionRepair:
		if name == "" {
			return app.ProviderAddResult{}, nil, fmt.Errorf("provider name is required")
		}
		res, err := svc.RepairProviderWithOptions(ctx, cfg, registry, name, app.RunOptions{DryRun: dryRun})
		return res, nil, err
	case providerActionList:
		enabled := app.ListEnabledProviders(cfg)
		names := make([]string, 0, len(enabled))
		for _, p := range enabled {
			names = append(names, p.Name)
		}
		return app.ProviderAddResult{}, names, nil
	default:
		return app.ProviderAddResult{}, nil, fmt.Errorf("unsupported provider action: %s", action)
	}
}
