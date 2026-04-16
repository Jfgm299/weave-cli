package integration

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestAddAsset_Integration_SuccessCreatesSymlinkAndConfigEntry(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	source := filepath.Join(repo, "shared", "skills", "sdd-orchestrator", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(source, []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	svc := app.ForgeService{
		ProjectRootDetector: integrationDetector{root: repo},
		ConfigValidator:     integrationValidator{},
		Executor:            fsops.Engine{},
		Writer:              config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")},
	}

	_, err := svc.AddAsset(context.Background(), config.Default(), app.AddAssetInput{
		Kind:        app.AssetKindSkill,
		Name:        "sdd-orchestrator",
		SourcePath:  source,
		ProjectPath: filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md"),
	})
	if err != nil {
		t.Fatalf("add asset: %v", err)
	}

	installed := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")
	if fi, err := os.Lstat(installed); err != nil {
		t.Fatalf("lstat installed: %v", err)
	} else if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected installed asset to be symlink")
	}

	b, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if !contains(string(b), "name: sdd-orchestrator") {
		t.Fatalf("expected config inventory update, got: %s", string(b))
	}
}

func TestAddAsset_Integration_CommandSuccessCreatesSymlinkAndConfigEntry(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	source := filepath.Join(repo, "shared", "commands", "pr-review.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(source, []byte("# command"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	svc := app.ForgeService{
		ProjectRootDetector: integrationDetector{root: repo},
		ConfigValidator:     integrationValidator{},
		Executor:            fsops.Engine{},
		Writer:              config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")},
	}

	_, err := svc.AddAsset(context.Background(), config.Default(), app.AddAssetInput{
		Kind:        app.AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  source,
		ProjectPath: filepath.Join(repo, ".agents", "commands", "pr-review.md"),
	})
	if err != nil {
		t.Fatalf("add asset: %v", err)
	}

	installed := filepath.Join(repo, ".agents", "commands", "pr-review.md")
	if fi, err := os.Lstat(installed); err != nil {
		t.Fatalf("lstat installed: %v", err)
	} else if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected installed asset to be symlink")
	}

	b, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if !contains(string(b), "name: pr-review") {
		t.Fatalf("expected command inventory update, got: %s", string(b))
	}
}

func TestAddAsset_Integration_SymlinkFailureLeavesConfigUnchanged(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	if err := os.WriteFile(filepath.Join(repo, ".agents"), []byte("blocked"), 0o644); err != nil {
		t.Fatalf("seed blocking .agents file: %v", err)
	}
	seed := []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	svc := app.ForgeService{
		ProjectRootDetector: integrationDetector{root: repo},
		ConfigValidator:     integrationValidator{},
		Executor:            fsops.Engine{},
		Writer:              config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")},
	}

	_, err := svc.AddAsset(context.Background(), config.Default(), app.AddAssetInput{
		Kind:        app.AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  filepath.Join(repo, "missing", "pr-review.md"),
		ProjectPath: filepath.Join(repo, ".agents", "commands", "pr-review.md"),
	})
	if err == nil {
		t.Fatalf("expected add failure")
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged on symlink failure")
	}
}

type integrationDetector struct {
	root string
}

func (d integrationDetector) Detect(_ context.Context) (string, error) {
	if d.root == "" {
		return "", errors.New("missing root")
	}
	return d.root, nil
}

type integrationValidator struct{}

func (integrationValidator) Validate(cfg config.Config) error {
	return (config.Validator{}).Validate(cfg)
}

func contains(s string, token string) bool {
	return strings.Contains(s, token)
}
