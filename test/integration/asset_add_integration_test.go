package integration

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
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

func TestAddAsset_Integration_CommandWithProviderProjections_WritesSharedAndProviderTargets(t *testing.T) {
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
		Kind:          app.AssetKindCommand,
		Name:          "pr-review",
		SourcePath:    source,
		ProjectPath:   filepath.Join(repo, ".agents", "commands", "pr-review.md"),
		CommandMeta:   &config.CommandMetaV1{ProviderCompat: []string{"claude-code", "codex", "opencode"}, SharedInstall: boolPtr(true)},
		AdditionalOps: []fsops.Operation{{Type: fsops.OpCreateLink, Path: filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md"), Target: source}},
	})
	if err != nil {
		t.Fatalf("add command with projections: %v", err)
	}

	for _, p := range []string{
		filepath.Join(repo, ".agents", "commands", "pr-review.md"),
		filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md"),
	} {
		fi, err := os.Lstat(p)
		if err != nil {
			t.Fatalf("lstat %s: %v", p, err)
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("expected %s to be symlink", p)
		}
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if !strings.Contains(string(after), "provider_compat:") || !strings.Contains(string(after), "shared_install: true") {
		t.Fatalf("expected command metadata persisted transactionally, got: %s", string(after))
	}
}

func TestAddAsset_Integration_CommandExclusiveProvider_DoesNotRequireSharedPath(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	source := filepath.Join(repo, "shared", "commands", "pr-review.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(source, []byte("# command"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repo, ".agents"), []byte("blocked"), 0o644); err != nil {
		t.Fatalf("seed blocking .agents path: %v", err)
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
		ProjectPath: "",
		CommandMeta: &config.CommandMetaV1{ProviderCompat: []string{"codex"}},
		AdditionalOps: []fsops.Operation{
			{Type: fsops.OpCreateLink, Path: filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md"), Target: source},
		},
	})
	if err != nil {
		t.Fatalf("exclusive provider add should not require shared .agents path: %v", err)
	}

	if _, err := os.Stat(filepath.Join(repo, ".agents", "commands", "pr-review.md")); err == nil {
		t.Fatalf("expected no shared command install for exclusive provider")
	}

	if fi, err := os.Lstat(filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md")); err != nil {
		t.Fatalf("expected codex exclusive wrapper path: %v", err)
	} else if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected codex exclusive wrapper to be symlink")
	}
}

func TestAddAsset_Integration_ApplyFailureRollsBackSharedAndProviderTargetsAndKeepsConfigUnchanged(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	seed := []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	commandSource := filepath.Join(repo, "shared", "commands", "pr-review.md")
	if err := os.MkdirAll(filepath.Dir(commandSource), 0o755); err != nil {
		t.Fatalf("mkdir command source: %v", err)
	}
	if err := os.WriteFile(commandSource, []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	exec := &rollbackAwareExecutor{}
	svc := app.ForgeService{
		ProjectRootDetector: integrationDetector{root: repo},
		ConfigValidator:     integrationValidator{},
		Executor:            exec,
		Writer:              config.AtomicFileWriter{Path: filepath.Join(repo, "weave.yaml")},
	}

	_, err := svc.AddAsset(context.Background(), config.Default(), app.AddAssetInput{
		Kind:        app.AssetKindCommand,
		Name:        "pr-review",
		SourcePath:  commandSource,
		ProjectPath: filepath.Join(repo, ".agents", "commands", "pr-review.md"),
		AdditionalOps: []fsops.Operation{{
			Type:   fsops.OpCreateLink,
			Path:   filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md"),
			Target: commandSource,
		}},
	})
	if err == nil {
		t.Fatalf("expected apply failure")
	}

	if len(exec.calls) != 2 {
		t.Fatalf("expected apply + rollback executor calls, got %d", len(exec.calls))
	}
	if !samePaths(exec.calls[0], exec.calls[1]) {
		t.Fatalf("expected rollback call to target all applied paths, got apply=%+v rollback=%+v", exec.calls[0], exec.calls[1])
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged after apply rollback")
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

type rollbackAwareExecutor struct {
	calls [][]string
	seen  bool
}

func (e *rollbackAwareExecutor) Apply(ctx context.Context, ops []fsops.Operation) error {
	paths := make([]string, 0, len(ops))
	for _, op := range ops {
		if op.Type == fsops.OpCreateLink || op.Type == fsops.OpRemovePath {
			paths = append(paths, op.Path)
		}
	}
	e.calls = append(e.calls, paths)
	if !e.seen {
		e.seen = true
		if err := (fsops.Engine{}).Apply(ctx, ops); err != nil {
			return err
		}
		return errors.New("injected apply failure")
	}
	return (fsops.Engine{}).Apply(ctx, ops)
}

func samePaths(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	am := map[string]int{}
	bm := map[string]int{}
	for _, v := range a {
		am[v]++
	}
	for _, v := range b {
		bm[v]++
	}
	return reflect.DeepEqual(am, bm)
}

func boolPtr(v bool) *bool { return &v }
