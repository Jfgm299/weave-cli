package app

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

var (
	ErrUnsupportedProvider     = errors.New("unsupported provider")
	ErrMissingProviderBinaries = errors.New("missing provider binaries")
)

type ProviderAdapter interface {
	Name() string
	RequiredBinaries() []string
	PlanSetup(projectRoot string) ([]fsops.Operation, error)
	PlanRepair(projectRoot string) ([]fsops.Operation, error)
	PlanRemove(projectRoot string) ([]fsops.Operation, error)
}

type ProviderRegistry interface {
	Get(name string) (ProviderAdapter, bool)
	SupportedNames() []string
}

type BinaryResolver interface {
	LookPath(file string) (string, error)
}

type ProviderService struct {
	ProjectRootDetector ProjectRootDetector
	ConfigValidator     ConfigValidator
	Executor            OperationExecutor
	Writer              ConfigWriter
	BinaryResolver      BinaryResolver
}

type ProviderAddResult struct {
	Root        string
	Provider    string
	Added       bool
	OpsApplied  int
	ConfigSaved bool
}

func (s ProviderService) AddProvider(ctx context.Context, cfg config.Config, registry ProviderRegistry, providerName string) (ProviderAddResult, error) {
	root, err := s.ProjectRootDetector.Detect(ctx)
	if err != nil {
		return ProviderAddResult{}, ErrNotInProjectRoot
	}

	if err := s.ConfigValidator.Validate(cfg); err != nil {
		return ProviderAddResult{}, WrapInvalidConfig(err)
	}

	adapter, ok := registry.Get(providerName)
	if !ok {
		supported := strings.Join(sortedStrings(registry.SupportedNames()), ", ")
		return ProviderAddResult{}, fmt.Errorf("%w: %s. Supported providers: %s", ErrUnsupportedProvider, providerName, supported)
	}

	if missing := s.missingBinaries(adapter.RequiredBinaries()); len(missing) > 0 {
		return ProviderAddResult{}, fmt.Errorf("%w for provider %q: %s. Install the missing binaries and run `weave provider repair %s`.", ErrMissingProviderBinaries, adapter.Name(), strings.Join(missing, ", "), adapter.Name())
	}

	ops, err := adapter.PlanSetup(root)
	if err != nil {
		return ProviderAddResult{}, err
	}

	if len(ops) > 0 {
		if err := s.Executor.Apply(ctx, ops); err != nil {
			return ProviderAddResult{}, fmt.Errorf("failed to apply provider setup operations for %q: %w. Run `weave provider repair %s` after fixing the filesystem issue", adapter.Name(), err, adapter.Name())
		}
	}

	nextCfg := cfg
	nextCfg.Providers = upsertProvider(nextCfg.Providers, config.Provider{Name: adapter.Name(), Enabled: true})

	if err := s.Writer.Write(nextCfg); err != nil {
		return ProviderAddResult{}, fmt.Errorf("provider setup applied but failed to persist weave.yaml: %w. Re-run `weave provider repair %s`", err, adapter.Name())
	}

	return ProviderAddResult{
		Root:        root,
		Provider:    adapter.Name(),
		Added:       true,
		OpsApplied:  len(ops),
		ConfigSaved: true,
	}, nil
}

func (s ProviderService) RemoveProvider(ctx context.Context, cfg config.Config, registry ProviderRegistry, providerName string) (ProviderAddResult, error) {
	root, err := s.ProjectRootDetector.Detect(ctx)
	if err != nil {
		return ProviderAddResult{}, ErrNotInProjectRoot
	}

	if err := s.ConfigValidator.Validate(cfg); err != nil {
		return ProviderAddResult{}, WrapInvalidConfig(err)
	}

	adapter, ok := registry.Get(providerName)
	if !ok {
		supported := strings.Join(sortedStrings(registry.SupportedNames()), ", ")
		return ProviderAddResult{}, fmt.Errorf("%w: %s. Supported providers: %s", ErrUnsupportedProvider, providerName, supported)
	}

	ops, err := adapter.PlanRemove(root)
	if err != nil {
		return ProviderAddResult{}, err
	}

	if len(ops) > 0 {
		if err := s.Executor.Apply(ctx, ops); err != nil {
			return ProviderAddResult{}, fmt.Errorf("failed to remove provider %q operations: %w. Run `weave provider repair %s` to reconcile", adapter.Name(), err, adapter.Name())
		}
	}

	nextCfg := cfg
	nextCfg.Providers = removeProvider(nextCfg.Providers, adapter.Name())

	if err := s.Writer.Write(nextCfg); err != nil {
		return ProviderAddResult{}, fmt.Errorf("provider remove applied but failed to persist weave.yaml: %w. Re-run `weave provider repair %s`", err, adapter.Name())
	}

	return ProviderAddResult{Root: root, Provider: adapter.Name(), OpsApplied: len(ops), ConfigSaved: true}, nil
}

func (s ProviderService) RepairProvider(ctx context.Context, cfg config.Config, registry ProviderRegistry, providerName string) (ProviderAddResult, error) {
	root, err := s.ProjectRootDetector.Detect(ctx)
	if err != nil {
		return ProviderAddResult{}, ErrNotInProjectRoot
	}

	if err := s.ConfigValidator.Validate(cfg); err != nil {
		return ProviderAddResult{}, WrapInvalidConfig(err)
	}

	adapter, ok := registry.Get(providerName)
	if !ok {
		supported := strings.Join(sortedStrings(registry.SupportedNames()), ", ")
		return ProviderAddResult{}, fmt.Errorf("%w: %s. Supported providers: %s", ErrUnsupportedProvider, providerName, supported)
	}

	if missing := s.missingBinaries(adapter.RequiredBinaries()); len(missing) > 0 {
		return ProviderAddResult{}, fmt.Errorf("%w for provider %q: %s. Install the missing binaries and run `weave provider repair %s`.", ErrMissingProviderBinaries, adapter.Name(), strings.Join(missing, ", "), adapter.Name())
	}

	ops, err := adapter.PlanRepair(root)
	if err != nil {
		return ProviderAddResult{}, err
	}

	if len(ops) > 0 {
		if err := s.Executor.Apply(ctx, ops); err != nil {
			return ProviderAddResult{}, fmt.Errorf("failed to repair provider %q operations: %w. Re-run `weave provider repair %s` after fixing the filesystem issue", adapter.Name(), err, adapter.Name())
		}
	}

	nextCfg := cfg
	nextCfg.Providers = upsertProvider(nextCfg.Providers, config.Provider{Name: adapter.Name(), Enabled: true})

	if err := s.Writer.Write(nextCfg); err != nil {
		return ProviderAddResult{}, fmt.Errorf("provider repair applied but failed to persist weave.yaml: %w. Re-run `weave provider repair %s`", err, adapter.Name())
	}

	return ProviderAddResult{Root: root, Provider: adapter.Name(), OpsApplied: len(ops), ConfigSaved: true}, nil
}

func ListEnabledProviders(cfg config.Config) []config.Provider {
	enabled := make([]config.Provider, 0, len(cfg.Providers))
	for _, p := range cfg.Providers {
		if p.Enabled {
			enabled = append(enabled, p)
		}
	}

	sort.Slice(enabled, func(i, j int) bool {
		return enabled[i].Name < enabled[j].Name
	})

	return enabled
}

func upsertProvider(in []config.Provider, provider config.Provider) []config.Provider {
	for i := range in {
		if in[i].Name == provider.Name {
			in[i] = provider
			return in
		}
	}

	out := make([]config.Provider, 0, len(in)+1)
	out = append(out, in...)
	return append(out, provider)
}

func removeProvider(in []config.Provider, name string) []config.Provider {
	out := make([]config.Provider, 0, len(in))
	for _, p := range in {
		if p.Name == name {
			continue
		}
		out = append(out, p)
	}
	return out
}

func (s ProviderService) missingBinaries(required []string) []string {
	resolver := s.BinaryResolver
	if resolver == nil {
		resolver = lookPathFunc(exec.LookPath)
	}

	missing := make([]string, 0)
	for _, bin := range required {
		if strings.TrimSpace(bin) == "" {
			continue
		}
		if _, err := resolver.LookPath(bin); err != nil {
			missing = append(missing, bin)
		}
	}
	return sortedStrings(missing)
}

func sortedStrings(values []string) []string {
	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}

type lookPathFunc func(string) (string, error)

func (f lookPathFunc) LookPath(file string) (string, error) { return f(file) }
