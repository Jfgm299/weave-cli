package cli

import (
	"context"
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

type ForgeRunResult struct {
	ServiceResult app.ForgeResult
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
	code, _, err := h.RunWithOptions(ctx, false)
	return code, err
}

func (h ForgeHandler) RunWithOptions(ctx context.Context, dryRun bool) (int, ForgeRunResult, error) {
	workdir := resolveWorkdir()
	loader := config.FileLoader{Path: filepath.Join(workdir, "weave.yaml")}
	cfg, err := loader.LoadOrDefault()
	if err != nil {
		return ExitInvalidConfig, ForgeRunResult{}, err
	}

	result, err := h.Service.RunWithOptions(ctx, cfg, app.RunOptions{DryRun: dryRun})
	if err != nil {
		return exitCodeForError(err), ForgeRunResult{}, err
	}

	return ExitOK, ForgeRunResult{ServiceResult: result}, nil
}

type projectRootDetector struct {
	Workdir string
}

func (d projectRootDetector) Detect(_ context.Context) (string, error) {
	workdir := d.Workdir
	if workdir == "" {
		workdir = resolveWorkdir()
	}
	root, err := detectProjectRootFrom(workdir)
	if err != nil {
		return "", fmt.Errorf("project root not detected: %w", err)
	}
	return root, nil
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
	root, err := detectProjectRootFrom(workdir)
	if err != nil {
		return nil, fmt.Errorf("project root not detected: %w", err)
	}

	candidates := []string{filepath.Join(root, ".agents/skills"), filepath.Join(root, ".agents/commands"), filepath.Join(root, ".agents/docs")}
	ops := make([]fsops.Operation, 0, len(candidates))
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			continue
		} else if !os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot inspect %s: %w", path, err)
		}
		ops = append(ops, fsops.Operation{Type: fsops.OpEnsureDir, Path: path})
	}

	return ops, nil
}

func resolveWorkdir() string {
	if wd := os.Getenv("WEAVE_WORKDIR"); wd != "" {
		if root, err := detectProjectRootFrom(wd); err == nil {
			return root
		}
		return wd
	}
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	if root, err := detectProjectRootFrom(wd); err == nil {
		return root
	}
	return wd
}

func detectProjectRootFrom(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(current, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", os.ErrNotExist
		}
		current = parent
	}
}
