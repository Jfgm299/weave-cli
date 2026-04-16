package cli

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/providers"
)

func Run(ctx context.Context, args []string) (int, error) {
	if len(args) == 0 || args[0] == "forge" {
		h := NewDefaultForgeHandler()
		return h.Run(ctx)
	}

	if len(args) >= 3 && args[0] == "skill" && args[1] == "add" {
		return runAssetAdd(ctx, assetKindSkill, args[2], args[3:])
	}

	if len(args) >= 3 && args[0] == "command" && args[1] == "add" {
		return runAssetAdd(ctx, assetKindCommand, args[2], args[3:])
	}

	if len(args) >= 2 && args[0] == "provider" {
		return runProvider(ctx, args[1:])
	}

	return ExitRuntimeError, fmt.Errorf("unsupported command")
}

func runProvider(ctx context.Context, args []string) (int, error) {
	action, name, err := parseProviderAction(args)
	if err != nil {
		return ExitRuntimeError, err
	}

	workdir := resolveWorkdir()
	svc := newProviderService(workdir, nil)
	registry := providers.NewDefaultRegistry()

	names, err := runProviderAction(ctx, svc, registry, action, name)
	if err != nil {
		if errors.Is(err, app.ErrInvalidConfig) {
			return ExitInvalidConfig, err
		}
		if errors.Is(err, app.ErrMissingProviderBinaries) {
			return ExitMissingDependency, err
		}
		return ExitRuntimeError, err
	}

	if action == providerActionList {
		for _, name := range names {
			fmt.Println(name)
		}
		return ExitOK, nil
	}

	return ExitOK, nil
}

func runAssetAdd(ctx context.Context, kind assetKind, name string, rest []string) (int, error) {
	fromFlag, err := parseFromFlag(rest)
	if err != nil {
		return ExitRuntimeError, err
	}

	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		return ExitInvalidConfig, err
	}

	service := newDefaultAssetAddService()
	_, err = service.Add(ctx, kind, name, fromFlag, cfg)
	if err != nil {
		if errors.Is(err, app.ErrInvalidConfig) {
			return ExitInvalidConfig, err
		}
		return ExitRuntimeError, err
	}

	return ExitOK, nil
}

func parseFromFlag(args []string) (string, error) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--from" {
			if i+1 >= len(args) {
				return "", fmt.Errorf("--from requires a value")
			}
			return args[i+1], nil
		}

		if strings.HasPrefix(a, "--from=") {
			return strings.TrimPrefix(a, "--from="), nil
		}
	}

	return "", nil
}

func parseProviderAction(args []string) (providerAction, string, error) {
	if len(args) == 0 {
		return "", "", fmt.Errorf("provider action is required: add|list|remove|repair")
	}

	action := providerAction(args[0])
	switch action {
	case providerActionList:
		return action, "", nil
	case providerActionAdd, providerActionRemove, providerActionRepair:
		if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
			return "", "", fmt.Errorf("provider name is required for %s", action)
		}
		return action, args[1], nil
	default:
		return "", "", fmt.Errorf("unsupported provider action: %s", args[0])
	}
}
