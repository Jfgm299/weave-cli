package providers

import (
	"path/filepath"
	"sort"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type Registry struct {
	adapters map[string]app.ProviderAdapter
}

func NewDefaultRegistry() Registry {
	adapters := map[string]app.ProviderAdapter{
		"claude-code": ClaudeCodeAdapter{},
		"codex":       CodexAdapter{},
		"opencode":    OpenCodeAdapter{},
	}

	return Registry{adapters: adapters}
}

func (r Registry) Get(name string) (app.ProviderAdapter, bool) {
	a, ok := r.adapters[name]
	return a, ok
}

func (r Registry) SupportedNames() []string {
	names := make([]string, 0, len(r.adapters))
	for name := range r.adapters {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

type ClaudeCodeAdapter struct{}

func (ClaudeCodeAdapter) Name() string { return "claude-code" }

func (ClaudeCodeAdapter) RequiredBinaries() []string { return []string{"claude"} }

func (ClaudeCodeAdapter) PlanSetup(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{
		link(filepath.Join(projectRoot, ".claude", "CLAUDE.md"), filepath.Join("..", ".agents", "AGENTS.md")),
		link(filepath.Join(projectRoot, ".claude", "commands"), filepath.Join("..", ".agents", "commands")),
		link(filepath.Join(projectRoot, ".claude", "docs"), filepath.Join("..", ".agents", "docs")),
	}, nil
}

func (ClaudeCodeAdapter) PlanRepair(projectRoot string) ([]fsops.Operation, error) {
	return ClaudeCodeAdapter{}.PlanSetup(projectRoot)
}

func (ClaudeCodeAdapter) PlanRemove(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{{Type: fsops.OpRemovePath, Path: filepath.Join(projectRoot, ".claude")}}, nil
}

type OpenCodeAdapter struct{}

func (OpenCodeAdapter) Name() string { return "opencode" }

func (OpenCodeAdapter) RequiredBinaries() []string { return []string{"opencode"} }

func (OpenCodeAdapter) PlanSetup(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{
		link(filepath.Join(projectRoot, ".opencode", "AGENTS.md"), filepath.Join("..", ".agents", "AGENTS.md")),
		link(filepath.Join(projectRoot, ".opencode", "commands"), filepath.Join("..", ".agents", "commands")),
		link(filepath.Join(projectRoot, ".opencode", "docs"), filepath.Join("..", ".agents", "docs")),
	}, nil
}

func (OpenCodeAdapter) PlanRepair(projectRoot string) ([]fsops.Operation, error) {
	return OpenCodeAdapter{}.PlanSetup(projectRoot)
}

func (OpenCodeAdapter) PlanRemove(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{{Type: fsops.OpRemovePath, Path: filepath.Join(projectRoot, ".opencode")}}, nil
}

type CodexAdapter struct{}

func (CodexAdapter) Name() string { return "codex" }

func (CodexAdapter) RequiredBinaries() []string { return []string{"codex"} }

func (CodexAdapter) PlanSetup(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{
		link(filepath.Join(projectRoot, ".codex", "AGENTS.md"), filepath.Join("..", ".agents", "AGENTS.md")),
		link(filepath.Join(projectRoot, ".codex", "commands"), filepath.Join("..", ".agents", "commands")),
		link(filepath.Join(projectRoot, ".codex", "docs"), filepath.Join("..", ".agents", "docs")),
	}, nil
}

func (CodexAdapter) PlanRepair(projectRoot string) ([]fsops.Operation, error) {
	return CodexAdapter{}.PlanSetup(projectRoot)
}

func (CodexAdapter) PlanRemove(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{{Type: fsops.OpRemovePath, Path: filepath.Join(projectRoot, ".codex")}}, nil
}

func link(path string, target string) fsops.Operation {
	return fsops.Operation{Type: fsops.OpCreateLink, Path: path, Target: target}
}
