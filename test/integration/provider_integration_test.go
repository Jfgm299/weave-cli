package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
	"github.com/Jfgm299/weave-cli/internal/providers"
)

func TestProviderAdd_Integration_EnablesMultipleProvidersAndCreatesProjectionLinks(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	svc := app.ProviderService{
		ProjectRootDetector: integrationDetector{root: repo},
		ConfigValidator:     integrationValidator{},
		Executor:            fsops.Engine{},
		Writer:              config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")},
		BinaryResolver: providerBinaryResolverStub{paths: map[string]string{
			"claude":   "/usr/local/bin/claude",
			"opencode": "/usr/local/bin/opencode",
		}},
	}

	registry := providers.NewDefaultRegistry()
	cfg := config.Default()

	if _, err := svc.AddProvider(context.Background(), cfg, registry, "claude-code"); err != nil {
		t.Fatalf("add claude-code: %v", err)
	}

	loaded, err := (config.FileLoader{Path: filepath.Join(repo, "weave.yaml")}).LoadOrDefault()
	if err != nil {
		t.Fatalf("load config after first add: %v", err)
	}

	if _, err := svc.AddProvider(context.Background(), loaded, registry, "opencode"); err != nil {
		t.Fatalf("add opencode: %v", err)
	}

	loaded, err = (config.FileLoader{Path: filepath.Join(repo, "weave.yaml")}).LoadOrDefault()
	if err != nil {
		t.Fatalf("load config after second add: %v", err)
	}

	if len(loaded.Providers) != 2 {
		t.Fatalf("expected two enabled providers, got %+v", loaded.Providers)
	}

	for _, p := range []string{
		filepath.Join(repo, ".claude", "CLAUDE.md"),
		filepath.Join(repo, ".claude", "commands"),
		filepath.Join(repo, ".claude", "docs"),
		filepath.Join(repo, ".opencode", "AGENTS.md"),
		filepath.Join(repo, ".opencode", "commands"),
		filepath.Join(repo, ".opencode", "docs"),
	} {
		fi, err := os.Lstat(p)
		if err != nil {
			t.Fatalf("expected provider projection path %s: %v", p, err)
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("expected provider projection %s to be symlink", p)
		}
	}
}

func TestProviderAdd_Integration_MissingBinaryKeepsConfigUnchanged(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	seed := []byte("version: 1\nproviders:\n  - name: opencode\n    enabled: true\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	svc := app.ProviderService{
		ProjectRootDetector: integrationDetector{root: repo},
		ConfigValidator:     integrationValidator{},
		Executor:            fsops.Engine{},
		Writer:              config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")},
		BinaryResolver:      providerBinaryResolverStub{},
	}

	_, err := svc.AddProvider(context.Background(), config.Config{
		Version:   1,
		Providers: []config.Provider{{Name: "opencode", Enabled: true}},
		Sources:   config.Sources{SkillsDir: "~/.weave/skills", CommandsDir: "~/.weave/commands"},
		Sync:      config.Sync{Mode: "symlink"},
		Skills:    []config.Asset{},
		Commands:  []config.Asset{},
	}, providers.NewDefaultRegistry(), "claude-code")
	if err == nil {
		t.Fatalf("expected provider add failure when binary is missing")
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}

	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged when binary dependency is missing")
	}
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
	return "", os.ErrNotExist
}
