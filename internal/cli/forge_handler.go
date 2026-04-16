package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type ForgeHandler struct {
	Service app.ForgeService
}

func NewDefaultForgeHandler() ForgeHandler {
	validator := config.Validator{}
	workdir := resolveWorkdir()

	service := app.ForgeService{
		ProjectRootDetector: projectRootDetector{Workdir: workdir},
		ConfigValidator:     validator,
		Planner:             forgePlanner{Workdir: workdir},
		Executor:            fsops.Engine{},
		Writer:              config.FileWriter{Path: filepath.Join(workdir, "weave.yaml")},
	}

	return ForgeHandler{Service: service}
}

func (h ForgeHandler) Run(ctx context.Context) (int, error) {
	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		return ExitInvalidConfig, err
	}

	result, err := h.Service.Run(ctx, cfg)
	if err != nil {
		if errors.Is(err, app.ErrInvalidConfig) {
			return ExitInvalidConfig, err
		}
		return ExitRuntimeError, err
	}

	_ = result
	return ExitOK, nil
}

type projectRootDetector struct {
	Workdir string
}

func (d projectRootDetector) Detect(_ context.Context) (string, error) {
	workdir := d.Workdir
	if workdir == "" {
		workdir = resolveWorkdir()
	}
	if _, err := os.Stat(filepath.Join(workdir, ".git")); err != nil {
		return "", fmt.Errorf("project root not detected: %w", err)
	}
	return workdir, nil
}

type forgePlanner struct {
	Workdir string
}

func (p forgePlanner) Plan(cfg config.Config) ([]fsops.Operation, error) {
	if cfg.Sync.Mode != "symlink" {
		return nil, fmt.Errorf("sync mode must be symlink")
	}

	workdir := p.Workdir
	if workdir == "" {
		workdir = resolveWorkdir()
	}

	candidates := []string{filepath.Join(workdir, ".agents/skills"), filepath.Join(workdir, ".agents/commands"), filepath.Join(workdir, ".agents/docs")}
	ops := make([]fsops.Operation, 0, len(candidates))
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("cannot inspect %s: %w", path, err)
		}
		ops = append(ops, fsops.Operation{Type: fsops.OpEnsureDir, Path: path})
	}

	return ops, nil
}

func resolveWorkdir() string {
	if wd := os.Getenv("WEAVE_WORKDIR"); wd != "" {
		return wd
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
