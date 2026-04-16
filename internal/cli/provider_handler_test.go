package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/providers"
)

func TestProviderRunAction_Add_UpdatesConfigWithEnabledProvider(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	cfg := config.Default()
	if err := (config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")}).Write(cfg); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	svc := newProviderService(repo, providerBinaryResolverFunc(func(bin string) (string, error) {
		return "/usr/local/bin/" + bin, nil
	}))

	if _, _, err := runProviderAction(context.Background(), svc, providers.NewDefaultRegistry(), providerActionAdd, "claude-code", false); err != nil {
		t.Fatalf("provider add failed: %v", err)
	}

	loaded, err := (config.FileLoader{Path: filepath.Join(repo, "weave.yaml")}).LoadOrDefault()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if len(loaded.Providers) != 1 || loaded.Providers[0].Name != "claude-code" || !loaded.Providers[0].Enabled {
		t.Fatalf("expected enabled claude provider, got %+v", loaded.Providers)
	}
}

func TestProviderRunAction_Add_MissingBinaryReturnsActionableMessage(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	svc := newProviderService(repo, providerBinaryResolverFunc(func(string) (string, error) {
		return "", os.ErrNotExist
	}))

	_, _, err := runProviderAction(context.Background(), svc, providers.NewDefaultRegistry(), providerActionAdd, "claude-code", false)
	if err == nil {
		t.Fatalf("expected provider add to fail when binary is missing")
	}

	if !strings.Contains(err.Error(), "Install the missing binaries") {
		t.Fatalf("expected actionable missing binary error, got %q", err.Error())
	}
}

func TestProviderRunAction_List_ReturnsSortedEnabledProviders(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	cfg := config.Default()
	cfg.Providers = []config.Provider{{Name: "opencode", Enabled: true}, {Name: "claude-code", Enabled: true}, {Name: "disabled", Enabled: false}}
	if err := (config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")}).Write(cfg); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	svc := newProviderService(repo, providerBinaryResolverFunc(func(bin string) (string, error) {
		return "/usr/local/bin/" + bin, nil
	}))

	_, names, err := runProviderAction(context.Background(), svc, providers.NewDefaultRegistry(), providerActionList, "", false)
	if err != nil {
		t.Fatalf("provider list failed: %v", err)
	}

	if len(names) != 2 || names[0] != "claude-code" || names[1] != "opencode" {
		t.Fatalf("expected sorted enabled providers, got %+v", names)
	}
}

func TestProviderRunAction_Remove_DeletesProjectionAndConfigEntry(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	cfg := config.Default()
	cfg.Providers = []config.Provider{{Name: "claude-code", Enabled: true}}
	if err := (config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")}).Write(cfg); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(repo, ".claude"), 0o755); err != nil {
		t.Fatalf("seed .claude dir: %v", err)
	}

	svc := newProviderService(repo, providerBinaryResolverFunc(func(bin string) (string, error) {
		return "/usr/local/bin/" + bin, nil
	}))

	if _, _, err := runProviderAction(context.Background(), svc, providers.NewDefaultRegistry(), providerActionRemove, "claude-code", false); err != nil {
		t.Fatalf("provider remove failed: %v", err)
	}

	loaded, err := (config.FileLoader{Path: filepath.Join(repo, "weave.yaml")}).LoadOrDefault()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if len(loaded.Providers) != 0 {
		t.Fatalf("expected provider removed from config, got %+v", loaded.Providers)
	}

	if _, err := os.Stat(filepath.Join(repo, ".claude")); !os.IsNotExist(err) {
		t.Fatalf("expected .claude directory removed, got err: %v", err)
	}
}

func TestProviderRunAction_Add_DryRunDoesNotMutateConfig(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	seed := []byte("version: 1\nproviders: []\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	svc := newProviderService(repo, providerBinaryResolverFunc(func(bin string) (string, error) {
		return "/usr/local/bin/" + bin, nil
	}))

	result, _, err := runProviderAction(context.Background(), svc, providers.NewDefaultRegistry(), providerActionAdd, "claude-code", true)
	if err != nil {
		t.Fatalf("provider add dry-run failed: %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected dry-run result")
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}

	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged during dry-run")
	}
}

func TestProviderRunAction_RemoveDryRun_DoesNotMutateConfig(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	seed := []byte("version: 1\nproviders:\n  - name: claude-code\n    enabled: true\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	svc := newProviderService(repo, providerBinaryResolverFunc(func(bin string) (string, error) {
		return "/usr/local/bin/" + bin, nil
	}))

	result, _, err := runProviderAction(context.Background(), svc, providers.NewDefaultRegistry(), providerActionRemove, "claude-code", true)
	if err != nil {
		t.Fatalf("provider remove dry-run failed: %v", err)
	}
	if !result.DryRun {
		t.Fatalf("expected dry-run result")
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged during remove dry-run")
	}
}

func TestProviderRunAction_RepairDryRun_DoesNotMutateConfig(t *testing.T) {
	repo := t.TempDir()
	t.Setenv("WEAVE_WORKDIR", repo)
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	seed := []byte("version: 1\nproviders:\n  - name: claude-code\n    enabled: false\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	svc := newProviderService(repo, providerBinaryResolverFunc(func(bin string) (string, error) {
		return "/usr/local/bin/" + bin, nil
	}))

	result, _, err := runProviderAction(context.Background(), svc, providers.NewDefaultRegistry(), providerActionRepair, "claude-code", true)
	if err != nil {
		t.Fatalf("provider repair dry-run failed: %v", err)
	}
	if !result.DryRun {
		t.Fatalf("expected dry-run result")
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged during repair dry-run")
	}
}
