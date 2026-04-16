package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/app"
	"github.com/Jfgm299/weave-cli/internal/cli"
	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestDoctor_Integration_ConfigWriteFailureRollsBackSymlinkInStrictMode(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	source := filepath.Join(repo, "shared", "skills", "sdd-orchestrator", "SKILL.md")
	installed := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(source, []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write source skill: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(installed), 0o755); err != nil {
		t.Fatalf("mkdir installed dir: %v", err)
	}
	if err := os.Symlink(source, installed); err != nil {
		t.Fatalf("seed symlink: %v", err)
	}

	cfg := config.Default()
	cfg.Skills = []config.Asset{{Name: "sdd-orchestrator", Source: source}}

	svc := app.ForgeService{
		ProjectRootDetector: integrationDetector{root: repo},
		ConfigValidator:     integrationValidator{},
		Executor:            successfulExecutor{},
		Writer:              failingWriter{},
	}

	_, err := svc.AddAsset(context.Background(), cfg, app.AddAssetInput{
		Kind:           app.AssetKindSkill,
		Name:           "sdd-orchestrator",
		SourcePath:     source,
		ProjectPath:    installed,
		ConflictPolicy: app.ConflictPolicyOverwrite,
	})
	if err == nil {
		t.Fatalf("expected add asset to fail when config write rename fails")
	}

	if _, err := os.Lstat(installed); !os.IsNotExist(err) {
		t.Fatalf("expected symlink rollback on failed config persistence, got err: %v", err)
	}

	result, err := (app.DoctorService{}).Run(context.Background(), repo, cfg)
	if err != nil {
		t.Fatalf("doctor run failed: %v", err)
	}

	if result.Status != app.DoctorStatusIssuesFound {
		t.Fatalf("expected doctor to report missing link issue after rollback, got %q", result.Status)
	}
}

func TestDoctor_Integration_ExitCodeMappingForIssueState(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	source := filepath.Join(repo, "shared", "commands", "pr-review.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source dir: %v", err)
	}
	if err := os.WriteFile(source, []byte("# cmd"), 0o644); err != nil {
		t.Fatalf("write source command: %v", err)
	}

	configBody := "version: 1\n" +
		"sources:\n" +
		"  skills_dir: ~/.weave/skills\n" +
		"  commands_dir: ~/.weave/commands\n" +
		"sync:\n" +
		"  mode: symlink\n" +
		"skills: []\n" +
		"commands:\n" +
		"  - name: pr-review\n" +
		"    source: " + source + "\n"

	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), []byte(configBody), 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	t.Setenv("WEAVE_WORKDIR", repo)
	code, err := cli.Run(context.Background(), []string{"doctor"})
	if err != nil {
		t.Fatalf("expected nil error for issue-reporting doctor, got %v", err)
	}

	if code != cli.ExitDoctorIssues {
		t.Fatalf("expected ExitDoctorIssues, got %d", code)
	}
}

type successfulExecutor struct{}

func (successfulExecutor) Apply(ctx context.Context, ops []fsops.Operation) error {
	return (fsops.Engine{}).Apply(ctx, ops)
}

type failingWriter struct{}

func (failingWriter) Write(config.Config) error {
	return os.ErrPermission
}
