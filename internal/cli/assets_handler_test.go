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

func TestAssetPathFor_BuildsSkillDestination(t *testing.T) {
	t.Parallel()

	got := assetPathFor(assetKindSkill, "/repo", "sdd-orchestrator")
	if got != filepath.Join("/repo", ".agents", "skills", "sdd-orchestrator", "SKILL.md") {
		t.Fatalf("unexpected skill destination: %s", got)
	}
}

func TestAssetPathFor_BuildsCommandDestination(t *testing.T) {
	t.Parallel()

	got := assetPathFor(assetKindCommand, "/repo", "pr-review")
	if got != filepath.Join("/repo", ".agents", "commands", "pr-review.md") {
		t.Fatalf("unexpected command destination: %s", got)
	}
}

func TestSourcePathFor_BuildsSkillSource(t *testing.T) {
	t.Parallel()

	got := sourcePathFor(assetKindSkill, "/src/skills", "sdd-orchestrator")
	if got != filepath.Join("/src/skills", "sdd-orchestrator", "SKILL.md") {
		t.Fatalf("unexpected skill source: %s", got)
	}
}

func TestSourcePathFor_BuildsCommandSource(t *testing.T) {
	t.Parallel()

	got := sourcePathFor(assetKindCommand, "/src/commands", "pr-review")
	if got != filepath.Join("/src/commands", "pr-review.md") {
		t.Fatalf("unexpected command source: %s", got)
	}
}

func TestAssetAddService_Add_CreatesSymlinkAndPersistsConfig(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	skillsRoot := filepath.Join(repo, "shared-skills")
	skillDir := filepath.Join(skillsRoot, "sdd-orchestrator")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}

	service := newAssetAddService(repo)
	cfg := config.Default()
	cfg.Sources.SkillsDir = skillsRoot
	_, err := service.Add(context.Background(), assetKindSkill, "sdd-orchestrator", "", false, "", cfg)
	if err != nil {
		t.Fatalf("add skill: %v", err)
	}

	installed := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")
	fi, err := os.Lstat(installed)
	if err != nil {
		t.Fatalf("lstat installed skill: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected installed skill to be symlink")
	}

	b, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if !containsString(string(b), "name: sdd-orchestrator") {
		t.Fatalf("expected skill entry in weave.yaml, got: %s", string(b))
	}
}

func TestAssetAddService_Add_StrictModeKeepsConfigUnchangedOnSymlinkFailure(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	seed := []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed config: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repo, ".agents"), []byte("blocked"), 0o644); err != nil {
		t.Fatalf("seed blocking .agents file: %v", err)
	}

	service := newAssetAddService(repo)
	cfg := config.Default()
	_, err := service.Add(context.Background(), assetKindSkill, "missing-skill", "", false, "", cfg)
	if err == nil {
		t.Fatalf("expected add to fail for missing source")
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected config unchanged on symlink failure")
	}
}

func TestAssetAddService_Add_DryRunDoesNotMutateFilesystemOrConfig(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	skillsRoot := filepath.Join(repo, "shared-skills")
	skillDir := filepath.Join(skillsRoot, "sdd-orchestrator")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}

	seed := []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	service := newAssetAddService(repo)
	cfg := config.Default()
	cfg.Sources.SkillsDir = skillsRoot
	result, err := service.Add(context.Background(), assetKindSkill, "sdd-orchestrator", "", true, "", cfg)
	if err != nil {
		t.Fatalf("dry-run add skill: %v", err)
	}

	if !result.DryRun {
		t.Fatalf("expected dry-run result")
	}

	installed := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")
	if _, err := os.Lstat(installed); !os.IsNotExist(err) {
		t.Fatalf("expected no installed symlink on dry-run, got err: %v", err)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged on dry-run")
	}
}

func TestAssetAddService_Add_ConflictSkipPolicyProducesNoMutation(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	skillsRoot := filepath.Join(repo, "shared-skills")
	skillDir := filepath.Join(skillsRoot, "sdd-orchestrator")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}

	installed := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(installed), 0o755); err != nil {
		t.Fatalf("mkdir installed dir: %v", err)
	}
	if err := os.WriteFile(installed, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed installed path: %v", err)
	}

	seed := []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	service := newAssetAddService(repo)
	cfg := config.Default()
	cfg.Sources.SkillsDir = skillsRoot
	result, err := service.Add(context.Background(), assetKindSkill, "sdd-orchestrator", "", false, app.ConflictPolicySkip, cfg)
	if err != nil {
		t.Fatalf("unexpected add result: %v", err)
	}
	if result.OpsPlanned != 0 || result.OpsApplied != 0 {
		t.Fatalf("expected skip conflict to plan/apply no ops, got %+v", result)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged on skip conflict policy")
	}
}

func containsString(s string, token string) bool {
	return strings.Contains(s, token)
}

func newAssetAddService(workdir string) assetAddService {
	validator := config.Validator{}
	return assetAddService{
		Service: app.ForgeService{
			ProjectRootDetector: projectRootDetector{Workdir: workdir},
			ConfigValidator:     validator,
			Executor:            fsops.Engine{},
			Writer:              config.AtomicFileWriter{Path: filepath.Join(workdir, "weave.yaml")},
		},
		Resolver: sourceResolver{
			lookupEnv: func(_ string) (string, bool) { return "", false },
			homeDir:   func() (string, error) { return workdir, nil },
		},
		Workdir: workdir,
	}
}
