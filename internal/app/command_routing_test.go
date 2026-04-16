package app

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/config"
)

func TestCommandInstallPlanner_Resolve_DefaultUsesEnabledProvidersFromStrategyLayer(t *testing.T) {
	t.Parallel()

	planner := CommandInstallPlanner{Registry: commandRoutingRegistryStub{}}
	plan, err := planner.Resolve("", config.Config{Providers: []config.Provider{{Name: "opencode", Enabled: true}, {Name: "claude-code", Enabled: true}, {Name: "codex", Enabled: true}}}, false, nil)
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	if !plan.SharedInstall {
		t.Fatalf("expected shared install plan")
	}
	if len(plan.ProviderTargets) != 3 || plan.ProviderTargets[0] != "claude-code" || plan.ProviderTargets[1] != "codex" || plan.ProviderTargets[2] != "opencode" {
		t.Fatalf("unexpected provider targets: %+v", plan.ProviderTargets)
	}
}

func TestCommandInstallPlanner_Resolve_UnsupportedProviderFailsInPlanningLayer(t *testing.T) {
	t.Parallel()

	planner := CommandInstallPlanner{Registry: commandRoutingRegistryStub{}}
	_, err := planner.Resolve("future-provider", config.Default(), false, nil)
	if err == nil {
		t.Fatalf("expected unsupported provider error")
	}
	if !strings.Contains(err.Error(), "unsupported provider") || !strings.Contains(err.Error(), "Supported providers") {
		t.Fatalf("expected actionable planner error, got %q", err.Error())
	}
}

func TestBuildCommandProjectionOps_UnsupportedProviderFails(t *testing.T) {
	t.Parallel()

	_, err := BuildCommandProjectionOps("/repo", "pr-review", "/src/pr-review.md", CommandInstallPlan{ProviderTargets: []string{"future-provider"}, SharedInstall: false})
	if err == nil {
		t.Fatalf("expected unsupported projection strategy error")
	}
	if !strings.Contains(err.Error(), "unsupported provider projection strategy") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBuildCommandProjectionOps_DefaultSharedModeProjectsCodexOnly(t *testing.T) {
	t.Parallel()

	ops, err := BuildCommandProjectionOps("/repo", "pr-review", "/src/pr-review.md", CommandInstallPlan{ProviderTargets: []string{"claude-code", "codex", "opencode"}, SharedInstall: true})
	if err != nil {
		t.Fatalf("unexpected projection planning error: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected codex-only projection in shared mode, got %+v", ops)
	}
	if ops[0].Path != filepath.Join("/repo", ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md") {
		t.Fatalf("unexpected projection path: %s", ops[0].Path)
	}
}

type commandRoutingRegistryStub struct{}

func (commandRoutingRegistryStub) Get(name string) (ProviderAdapter, bool) {
	supported := map[string]struct{}{"claude-code": {}, "codex": {}, "opencode": {}}
	_, ok := supported[name]
	if !ok {
		return nil, false
	}
	return doctorAdapterStub{name: name}, true
}

func (commandRoutingRegistryStub) SupportedNames() []string {
	return []string{"claude-code", "codex", "opencode"}
}
