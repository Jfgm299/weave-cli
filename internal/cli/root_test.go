package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestParseFromFlag_ParsesLongForm(t *testing.T) {
	t.Parallel()

	g, err := parseFromFlag([]string{"--from", "/tmp/src"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if g != "/tmp/src" {
		t.Fatalf("expected /tmp/src, got %q", g)
	}
}

func TestParseFromFlag_ParsesEqualsForm(t *testing.T) {
	t.Parallel()

	g, err := parseFromFlag([]string{"--from=/tmp/src"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if g != "/tmp/src" {
		t.Fatalf("expected /tmp/src, got %q", g)
	}
}

func TestParseFromFlag_MissingValueFails(t *testing.T) {
	t.Parallel()

	_, err := parseFromFlag([]string{"--from"})
	if err == nil {
		t.Fatalf("expected missing value error")
	}
}

func TestParseProviderAction_ListParsesWithoutName(t *testing.T) {
	t.Parallel()

	action, name, dryRun, err := parseProviderAction([]string{"list"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if action != providerActionList {
		t.Fatalf("expected providerActionList, got %s", action)
	}

	if name != "" {
		t.Fatalf("expected empty provider name for list, got %q", name)
	}

	if dryRun {
		t.Fatalf("expected dryRun false for list")
	}
}

func TestParseProviderAction_AddMissingNameFails(t *testing.T) {
	t.Parallel()

	_, _, _, err := parseProviderAction([]string{"add"})
	if err == nil {
		t.Fatalf("expected provider add parse error when name is missing")
	}
}

func TestParseProviderAction_AddParsesDryRun(t *testing.T) {
	t.Parallel()

	action, name, dryRun, err := parseProviderAction([]string{"add", "claude-code", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if action != providerActionAdd || name != "claude-code" || !dryRun {
		t.Fatalf("unexpected parse output: action=%s name=%s dryRun=%v", action, name, dryRun)
	}
}

func TestParseProviderAction_ListRejectsDryRun(t *testing.T) {
	t.Parallel()

	_, _, _, err := parseProviderAction([]string{"list", "--dry-run"})
	if err == nil {
		t.Fatalf("expected list to reject --dry-run")
	}
}

func TestParseAddFlags_RejectsUnknownFlag(t *testing.T) {
	t.Parallel()

	_, _, _, _, err := parseAddFlags([]string{"--unknown"})
	if err == nil {
		t.Fatalf("expected unsupported flag error")
	}
}

func TestParseAddFlags_ParsesConflictPolicyFlags(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		args   []string
		policy app.ConflictPolicy
	}{
		{name: "overwrite", args: []string{"--overwrite"}, policy: app.ConflictPolicyOverwrite},
		{name: "skip", args: []string{"--skip"}, policy: app.ConflictPolicySkip},
		{name: "backup", args: []string{"--backup"}, policy: app.ConflictPolicyBackup},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, _, policy, _, err := parseAddFlags(tc.args)
			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}

			if policy != tc.policy {
				t.Fatalf("expected policy %q, got %q", tc.policy, policy)
			}
		})
	}
}

func TestParseAddFlags_RejectsMultipleConflictPolicyFlags(t *testing.T) {
	t.Parallel()

	_, _, _, _, err := parseAddFlags([]string{"--overwrite", "--backup"})
	if err == nil {
		t.Fatalf("expected conflict policy parsing error")
	}
}

func TestParseAddFlags_ParsesProviderFlag(t *testing.T) {
	t.Parallel()

	_, _, _, provider, err := parseAddFlags([]string{"--provider", "codex"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if provider != "codex" {
		t.Fatalf("expected provider codex, got %q", provider)
	}
}

func TestParseAddFlags_ParsesProviderEqualsFlag(t *testing.T) {
	t.Parallel()

	_, _, _, provider, err := parseAddFlags([]string{"--provider=opencode"})
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if provider != "opencode" {
		t.Fatalf("expected provider opencode, got %q", provider)
	}
}

func TestResolveCommandInstallPlan_DefaultTargetsAllEnabledProviders(t *testing.T) {
	t.Parallel()

	plan, err := resolveCommandInstallPlan(assetKindCommand, "", config.Config{Providers: []config.Provider{
		{Name: "opencode", Enabled: true},
		{Name: "claude-code", Enabled: true},
		{Name: "codex", Enabled: true},
	}})
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	if plan == nil || !plan.SharedInstall {
		t.Fatalf("expected shared install plan, got %+v", plan)
	}
	if len(plan.ProviderTargets) != 3 || plan.ProviderTargets[0] != "claude-code" || plan.ProviderTargets[1] != "codex" || plan.ProviderTargets[2] != "opencode" {
		t.Fatalf("unexpected provider targets: %+v", plan.ProviderTargets)
	}
}

func TestResolveCommandInstallPlan_ProviderFlagCreatesExclusivePlan(t *testing.T) {
	t.Parallel()

	plan, err := resolveCommandInstallPlan(assetKindCommand, "codex", config.Default())
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	if plan == nil || plan.SharedInstall {
		t.Fatalf("expected exclusive plan, got %+v", plan)
	}
	if len(plan.ProviderTargets) != 1 || plan.ProviderTargets[0] != "codex" {
		t.Fatalf("unexpected provider targets: %+v", plan.ProviderTargets)
	}
}

func TestResolveCommandInstallPlanWithDeps_DelegatesUnsupportedProviderToPlanner(t *testing.T) {
	t.Parallel()

	_, err := resolveCommandInstallPlanWithDeps(assetKindCommand, "future-provider", config.Default(), app.CommandInstallPlanner{Registry: commandRoutingRegistryStub{}})
	if err == nil {
		t.Fatalf("expected unsupported provider error from planning layer")
	}
	if !strings.Contains(err.Error(), "unsupported provider") {
		t.Fatalf("expected planning-layer provider validation error, got %q", err.Error())
	}
}

func TestResolveCommandInstallPlan_NoProvidersNonInteractiveFails(t *testing.T) {
	t.Parallel()

	orig := isInteractiveSession
	isInteractiveSession = func() bool { return false }
	t.Cleanup(func() { isInteractiveSession = orig })

	_, err := resolveCommandInstallPlan(assetKindCommand, "", config.Default())
	if err == nil {
		t.Fatalf("expected no-provider non-interactive failure")
	}
	if !strings.Contains(err.Error(), "No providers are currently enabled") || !strings.Contains(err.Error(), "--provider") {
		t.Fatalf("expected actionable no-provider guidance, got %q", err.Error())
	}
}

type commandRoutingRegistryStub struct{}

func (commandRoutingRegistryStub) Get(name string) (app.ProviderAdapter, bool) {
	if name == "claude-code" || name == "codex" || name == "opencode" {
		return rootProviderAdapterStub{name: name}, true
	}
	return nil, false
}

func (commandRoutingRegistryStub) SupportedNames() []string {
	return []string{"claude-code", "codex", "opencode"}
}

type rootProviderAdapterStub struct{ name string }

func (s rootProviderAdapterStub) Name() string                               { return s.name }
func (rootProviderAdapterStub) RequiredBinaries() []string                   { return nil }
func (rootProviderAdapterStub) PlanSetup(string) ([]fsops.Operation, error)  { return nil, nil }
func (rootProviderAdapterStub) PlanRepair(string) ([]fsops.Operation, error) { return nil, nil }
func (rootProviderAdapterStub) PlanRemove(string) ([]fsops.Operation, error) { return nil, nil }

func TestParseDryRunOnly_RejectsUnknownFlag(t *testing.T) {
	t.Parallel()

	_, err := parseDryRunOnly([]string{"--json"}, "forge")
	if err == nil {
		t.Fatalf("expected unsupported flag error")
	}
}

func TestRun_HelpPrintsQuickstart(t *testing.T) {
	t.Parallel()

	out := captureStdout(t, func() {
		code, err := Run(context.Background(), []string{"--help"})
		if err != nil {
			t.Fatalf("unexpected help error: %v", err)
		}
		if code != ExitOK {
			t.Fatalf("expected ExitOK, got %d", code)
		}
	})

	if !strings.Contains(out, "60-second quickstart") {
		t.Fatalf("expected quickstart in help output, got: %s", out)
	}
}

func TestRun_MigrateDryRun(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte("version: 0\nsync:\n  mode: symlink\nskills: []\ncommands: []\n"), 0o644); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	out := captureStdout(t, func() {
		code, err := Run(context.Background(), []string{"migrate", "--dry-run"})
		if err != nil {
			t.Fatalf("unexpected migrate error: %v", err)
		}
		if code != ExitOK {
			t.Fatalf("expected ExitOK, got %d", code)
		}
	})

	if !strings.Contains(out, "would upgrade weave.yaml schema") {
		t.Fatalf("expected dry-run migration summary, got: %s", out)
	}
}
