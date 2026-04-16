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

func runProviderAction(ctx context.Context, svc app.ProviderService, registry providers.Registry, action providerAction, name string) ([]string, error) {
	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		return nil, err
	}

	switch action {
	case providerActionAdd:
		if name == "" {
			return nil, fmt.Errorf("provider name is required")
		}
		_, err := svc.AddProvider(ctx, cfg, registry, name)
		return nil, err
	case providerActionRemove:
		if name == "" {
			return nil, fmt.Errorf("provider name is required")
		}
		_, err := svc.RemoveProvider(ctx, cfg, registry, name)
		return nil, err
	case providerActionRepair:
		if name == "" {
			return nil, fmt.Errorf("provider name is required")
		}
		_, err := svc.RepairProvider(ctx, cfg, registry, name)
		return nil, err
	case providerActionList:
		enabled := app.ListEnabledProviders(cfg)
		names := make([]string, 0, len(enabled))
		for _, p := range enabled {
			names = append(names, p.Name)
		}
		return names, nil
	default:
		return nil, fmt.Errorf("unsupported provider action: %s", action)
	}
}
