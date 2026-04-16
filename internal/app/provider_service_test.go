package app

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestProviderService_AddProvider_EnablesMultipleProviders(t *testing.T) {
	t.Parallel()

	executor := &providerExecutorSpy{}
	writer := &providerWriterSpy{}

	sut := ProviderService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
		BinaryResolver: providerBinaryResolverStub{paths: map[string]string{
			"claude":   "/bin/claude",
			"opencode": "/bin/opencode",
		}},
	}

	registry := providerRegistryStub{adapters: map[string]ProviderAdapter{
		"claude-code": providerAdapterStub{
			name:             "claude-code",
			requiredBinaries: []string{"claude"},
			setupOps:         []fsops.Operation{{Type: fsops.OpCreateLink, Path: "/tmp/proj/.claude/CLAUDE.md", Target: "../.agents/AGENTS.md"}},
		},
		"opencode": providerAdapterStub{
			name:             "opencode",
			requiredBinaries: []string{"opencode"},
			setupOps:         []fsops.Operation{{Type: fsops.OpCreateLink, Path: "/tmp/proj/.opencode/AGENTS.md", Target: "../.agents/AGENTS.md"}},
		},
	}}

	cfg := config.Default()

	_, err := sut.AddProvider(context.Background(), cfg, registry, "claude-code")
	if err != nil {
		t.Fatalf("unexpected add provider error: %v", err)
	}

	if !writer.called {
		t.Fatalf("expected config write after provider add")
	}

	if len(writer.lastCfg.Providers) != 1 || writer.lastCfg.Providers[0].Name != "claude-code" {
		t.Fatalf("expected claude-code provider in config, got %+v", writer.lastCfg.Providers)
	}

	writer.called = false
	_, err = sut.AddProvider(context.Background(), writer.lastCfg, registry, "opencode")
	if err != nil {
		t.Fatalf("unexpected second add provider error: %v", err)
	}

	if !writer.called {
		t.Fatalf("expected config write after second provider add")
	}

	if len(writer.lastCfg.Providers) != 2 {
		t.Fatalf("expected two providers enabled, got %+v", writer.lastCfg.Providers)
	}
}

func TestProviderService_AddProvider_MissingBinaryReturnsActionableError(t *testing.T) {
	t.Parallel()

	executor := &providerExecutorSpy{}
	writer := &providerWriterSpy{}

	sut := ProviderService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
		BinaryResolver:      providerBinaryResolverStub{},
	}

	registry := providerRegistryStub{adapters: map[string]ProviderAdapter{
		"claude-code": providerAdapterStub{
			name:             "claude-code",
			requiredBinaries: []string{"claude"},
			setupOps:         []fsops.Operation{{Type: fsops.OpCreateLink, Path: "/tmp/proj/.claude/CLAUDE.md", Target: "../.agents/AGENTS.md"}},
		},
	}}

	_, err := sut.AddProvider(context.Background(), config.Default(), registry, "claude-code")
	if err == nil {
		t.Fatalf("expected add provider to fail when binary is missing")
	}

	if !errors.Is(err, ErrMissingProviderBinaries) {
		t.Fatalf("expected ErrMissingProviderBinaries, got %v", err)
	}

	msg := err.Error()
	if !strings.Contains(msg, "claude") || !strings.Contains(msg, "Install the missing binaries") {
		t.Fatalf("expected actionable missing binary message, got %q", msg)
	}

	if !strings.Contains(msg, DocsPathProviders) || !strings.Contains(msg, DocsURL(DocsPathProviders)) {
		t.Fatalf("expected docs references in missing binary message, got %q", msg)
	}

	if executor.called {
		t.Fatalf("expected fs executor not called when prerequisites are missing")
	}

	if writer.called {
		t.Fatalf("expected config writer not called when prerequisites are missing")
	}
}

func TestProviderService_AddProvider_UnsupportedProviderReturnsSupportedList(t *testing.T) {
	t.Parallel()

	sut := ProviderService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            &providerExecutorSpy{},
		Writer:              &providerWriterSpy{},
	}

	registry := providerRegistryStub{supported: []string{"claude-code", "opencode"}}

	_, err := sut.AddProvider(context.Background(), config.Default(), registry, "unknown")
	if err == nil {
		t.Fatalf("expected unsupported provider error")
	}

	if !errors.Is(err, ErrUnsupportedProvider) {
		t.Fatalf("expected ErrUnsupportedProvider, got %v", err)
	}

	if !strings.Contains(err.Error(), "Supported providers: claude-code, opencode") {
		t.Fatalf("expected actionable supported providers message, got %q", err.Error())
	}

	if !strings.Contains(err.Error(), DocsPathProviders) || !strings.Contains(err.Error(), DocsURL(DocsPathProviders)) {
		t.Fatalf("expected docs references in unsupported provider error, got %q", err.Error())
	}
}

func TestProviderService_AddProvider_DryRunDoesNotApplyOrPersist(t *testing.T) {
	t.Parallel()

	executor := &providerExecutorSpy{}
	writer := &providerWriterSpy{}

	sut := ProviderService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
		BinaryResolver: providerBinaryResolverStub{paths: map[string]string{
			"claude": "/bin/claude",
		}},
	}

	registry := providerRegistryStub{adapters: map[string]ProviderAdapter{
		"claude-code": providerAdapterStub{
			name:             "claude-code",
			requiredBinaries: []string{"claude"},
			setupOps:         []fsops.Operation{{Type: fsops.OpCreateLink, Path: "/tmp/proj/.claude/CLAUDE.md", Target: "../.agents/AGENTS.md"}},
		},
	}}

	result, err := sut.AddProviderWithOptions(context.Background(), config.Default(), registry, "claude-code", RunOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected dry-run error: %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected dry-run result")
	}

	if executor.called {
		t.Fatalf("executor should not run during dry-run")
	}

	if writer.called {
		t.Fatalf("writer should not run during dry-run")
	}
}

func TestProviderService_AddProvider_RejectsOpsOutsideRoot(t *testing.T) {
	t.Parallel()

	sut := ProviderService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            &providerExecutorSpy{},
		Writer:              &providerWriterSpy{},
		BinaryResolver: providerBinaryResolverStub{paths: map[string]string{
			"claude": "/bin/claude",
		}},
	}

	registry := providerRegistryStub{adapters: map[string]ProviderAdapter{
		"claude-code": providerAdapterStub{
			name:             "claude-code",
			requiredBinaries: []string{"claude"},
			setupOps:         []fsops.Operation{{Type: fsops.OpCreateLink, Path: "/tmp/other/.claude/CLAUDE.md", Target: "../.agents/AGENTS.md"}},
		},
	}}

	_, err := sut.AddProvider(context.Background(), config.Default(), registry, "claude-code")
	if err == nil {
		t.Fatalf("expected unsafe path guard error")
	}

	if !errors.Is(err, ErrUnsafeMutationPath) {
		t.Fatalf("expected ErrUnsafeMutationPath, got %v", err)
	}
}

func TestProviderService_RemoveProvider_RemovesConfigEntryAndAppliesOps(t *testing.T) {
	t.Parallel()

	executor := &providerExecutorSpy{}
	writer := &providerWriterSpy{}

	sut := ProviderService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
	}

	registry := providerRegistryStub{adapters: map[string]ProviderAdapter{
		"claude-code": providerAdapterStub{
			name:      "claude-code",
			removeOps: []fsops.Operation{{Type: fsops.OpRemovePath, Path: "/tmp/proj/.claude"}},
		},
	}}

	cfg := config.Default()
	cfg.Providers = []config.Provider{{Name: "claude-code", Enabled: true}, {Name: "opencode", Enabled: true}}

	_, err := sut.RemoveProvider(context.Background(), cfg, registry, "claude-code")
	if err != nil {
		t.Fatalf("unexpected remove provider error: %v", err)
	}

	if !executor.called {
		t.Fatalf("expected fs operations to run on remove")
	}

	if len(writer.lastCfg.Providers) != 1 || writer.lastCfg.Providers[0].Name != "opencode" {
		t.Fatalf("expected only opencode to remain, got %+v", writer.lastCfg.Providers)
	}
}

func TestProviderService_RepairProvider_ReconcilesProjectionAndPersistsConfig(t *testing.T) {
	t.Parallel()

	executor := &providerExecutorSpy{}
	writer := &providerWriterSpy{}

	sut := ProviderService{
		ProjectRootDetector: detectorStub{root: "/tmp/proj"},
		ConfigValidator:     validatorStub{},
		Executor:            executor,
		Writer:              writer,
		BinaryResolver: providerBinaryResolverStub{paths: map[string]string{
			"opencode": "/bin/opencode",
		}},
	}

	registry := providerRegistryStub{adapters: map[string]ProviderAdapter{
		"opencode": providerAdapterStub{
			name:             "opencode",
			requiredBinaries: []string{"opencode"},
			repairOps:        []fsops.Operation{{Type: fsops.OpCreateLink, Path: "/tmp/proj/.opencode/AGENTS.md", Target: "../.agents/AGENTS.md"}},
		},
	}}

	_, err := sut.RepairProvider(context.Background(), config.Default(), registry, "opencode")
	if err != nil {
		t.Fatalf("unexpected repair provider error: %v", err)
	}

	if !executor.called {
		t.Fatalf("expected fs operations to run on repair")
	}

	if len(writer.lastCfg.Providers) != 1 || writer.lastCfg.Providers[0].Name != "opencode" || !writer.lastCfg.Providers[0].Enabled {
		t.Fatalf("expected provider enabled after repair, got %+v", writer.lastCfg.Providers)
	}
}

func TestListEnabledProviders_ReturnsSortedEnabledProviders(t *testing.T) {
	t.Parallel()

	got := ListEnabledProviders(config.Config{Providers: []config.Provider{{Name: "opencode", Enabled: true}, {Name: "claude-code", Enabled: true}, {Name: "disabled", Enabled: false}}})
	if len(got) != 2 {
		t.Fatalf("expected 2 enabled providers, got %d", len(got))
	}

	if got[0].Name != "claude-code" || got[1].Name != "opencode" {
		t.Fatalf("expected sorted providers [claude-code opencode], got %+v", got)
	}
}

type providerAdapterStub struct {
	name             string
	requiredBinaries []string
	setupOps         []fsops.Operation
	repairOps        []fsops.Operation
	removeOps        []fsops.Operation
}

func (s providerAdapterStub) Name() string { return s.name }

func (s providerAdapterStub) RequiredBinaries() []string {
	return append([]string(nil), s.requiredBinaries...)
}

func (s providerAdapterStub) PlanSetup(_ string) ([]fsops.Operation, error) {
	return append([]fsops.Operation(nil), s.setupOps...), nil
}

func (s providerAdapterStub) PlanRepair(_ string) ([]fsops.Operation, error) {
	return append([]fsops.Operation(nil), s.repairOps...), nil
}

func (s providerAdapterStub) PlanRemove(_ string) ([]fsops.Operation, error) {
	return append([]fsops.Operation(nil), s.removeOps...), nil
}

type providerRegistryStub struct {
	adapters  map[string]ProviderAdapter
	supported []string
}

func (s providerRegistryStub) Get(name string) (ProviderAdapter, bool) {
	if s.adapters == nil {
		return nil, false
	}
	a, ok := s.adapters[name]
	return a, ok
}

func (s providerRegistryStub) SupportedNames() []string {
	if len(s.supported) > 0 {
		return append([]string(nil), s.supported...)
	}
	out := make([]string, 0, len(s.adapters))
	for name := range s.adapters {
		out = append(out, name)
	}
	return out
}

type providerBinaryResolverStub struct {
	paths map[string]string
}

func (s providerBinaryResolverStub) LookPath(file string) (string, error) {
	if s.paths != nil {
		if p, ok := s.paths[file]; ok {
			return p, nil
		}
	}
	return "", errors.New("not found")
}

type providerExecutorSpy struct {
	called  bool
	err     error
	lastOps []fsops.Operation
}

func (s *providerExecutorSpy) Apply(_ context.Context, ops []fsops.Operation) error {
	s.called = true
	s.lastOps = append([]fsops.Operation(nil), ops...)
	return s.err
}

type providerWriterSpy struct {
	called  bool
	err     error
	lastCfg config.Config
}

func (s *providerWriterSpy) Write(cfg config.Config) error {
	s.called = true
	s.lastCfg = cfg
	return s.err
}
