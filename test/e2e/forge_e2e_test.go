//go:build e2e

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestForge_E2E_IdempotentAndNonDestructive(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	agentsDir := filepath.Join(repo, ".agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatalf("mkdir .agents: %v", err)
	}

	agentsPath := filepath.Join(agentsDir, "AGENTS.md")
	original := []byte("existing-agents-file")
	if err := os.WriteFile(agentsPath, original, 0o644); err != nil {
		t.Fatalf("write AGENTS.md: %v", err)
	}

	if out, err := runCLI(repo, root, nil, nil); err != nil {
		t.Fatalf("first forge failed: %v\n%s", err, out)
	}

	weavePath := filepath.Join(repo, "weave.yaml")
	b1, err := os.ReadFile(weavePath)
	if err != nil {
		t.Fatalf("read weave.yaml after first run: %v", err)
	}

	if out, err := runCLI(repo, root, nil, nil); err != nil {
		t.Fatalf("second forge failed: %v\n%s", err, out)
	}

	b2, err := os.ReadFile(weavePath)
	if err != nil {
		t.Fatalf("read weave.yaml after second run: %v", err)
	}

	if string(b1) != string(b2) {
		t.Fatalf("expected weave.yaml to remain deterministic across runs")
	}

	afterAgents, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("read AGENTS.md after runs: %v", err)
	}
	if string(afterAgents) != string(original) {
		t.Fatalf("expected existing AGENTS.md not overwritten")
	}

	for _, p := range []string{".agents/skills", ".agents/commands", ".agents/docs"} {
		if _, err := os.Stat(filepath.Join(repo, p)); err != nil {
			t.Fatalf("expected %s to exist: %v", p, err)
		}
	}
}

func TestForge_E2E_InvalidConfigModeFailsWithoutMutatingFilesystem(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	seed := []byte("version: 1\nsync:\n  mode: copy\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, nil, nil)
	if err == nil {
		t.Fatalf("expected forge to fail for invalid sync mode")
	}

	if !strings.Contains(out, "sync.mode must be 'symlink' in v1") {
		t.Fatalf("expected actionable mode validation error, got: %s", out)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml to remain unchanged on invalid config")
	}

	if _, err := os.Stat(filepath.Join(repo, ".agents/skills")); err == nil {
		t.Fatalf("expected no mutation to .agents/skills on invalid config")
	}
}

func TestSkillAdd_E2E_StrictFailureKeepsConfigUnchanged(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".agents"), []byte("blocked"), 0o644); err != nil {
		t.Fatalf("seed blocking .agents file: %v", err)
	}

	seed := []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"skill", "add", "does-not-exist"}, nil)
	if err == nil {
		t.Fatalf("expected skill add to fail")
	}
	if !strings.Contains(out, "not a directory") {
		t.Fatalf("expected actionable fs error, got: %s", out)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged on symlink failure")
	}
}

func TestCommandAdd_E2E_StrictFailureKeepsConfigUnchanged(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".agents"), []byte("blocked"), 0o644); err != nil {
		t.Fatalf("seed blocking .agents file: %v", err)
	}

	seed := []byte("version: 1\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("seed weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"command", "add", "does-not-exist"}, nil)
	if err == nil {
		t.Fatalf("expected command add to fail")
	}
	if !strings.Contains(out, "not a directory") {
		t.Fatalf("expected actionable fs error, got: %s", out)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged on symlink failure")
	}
}

func TestAssetAdd_E2E_UsesSourcePrecedenceAndSymlinkOnly(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	cfgSkills := filepath.Join(repo, "cfg-skills")
	envSkills := filepath.Join(repo, "env-skills")
	flagSkills := filepath.Join(repo, "flag-skills")
	cfgCommands := filepath.Join(repo, "cfg-commands")
	envCommands := filepath.Join(repo, "env-commands")

	for _, d := range []string{cfgSkills, envSkills, flagSkills, cfgCommands, envCommands} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", d, err)
		}
	}

	if err := os.MkdirAll(filepath.Join(flagSkills, "sdd-orchestrator"), 0o755); err != nil {
		t.Fatalf("mkdir skill source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(flagSkills, "sdd-orchestrator", "SKILL.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write skill source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(envCommands, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	cfg := []byte("version: 1\nsources:\n  skills_dir: " + cfgSkills + "\n  commands_dir: " + cfgCommands + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	env := []string{"WEAVE_SKILLS_DIR=" + envSkills, "WEAVE_COMMANDS_DIR=" + envCommands}

	out, err := runCLI(repo, root, []string{"skill", "add", "sdd-orchestrator", "--from", flagSkills}, env)
	if err != nil {
		t.Fatalf("skill add failed: %v\n%s", err, out)
	}

	out, err = runCLI(repo, root, []string{"command", "add", "pr-review"}, env)
	if err != nil {
		t.Fatalf("command add failed: %v\n%s", err, out)
	}

	skillLink := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")
	cmdLink := filepath.Join(repo, ".agents", "commands", "pr-review.md")

	if fi, err := os.Lstat(skillLink); err != nil {
		t.Fatalf("lstat skill link: %v", err)
	} else if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected skill install to be symlink")
	}

	if fi, err := os.Lstat(cmdLink); err != nil {
		t.Fatalf("lstat command link: %v", err)
	} else if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected command install to be symlink")
	}

	conf, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	outCfg := string(conf)
	if !strings.Contains(outCfg, filepath.ToSlash(filepath.Join(flagSkills, "sdd-orchestrator", "SKILL.md"))) {
		t.Fatalf("expected skill source from --from flag, got: %s", outCfg)
	}
	if !strings.Contains(outCfg, filepath.ToSlash(filepath.Join(envCommands, "pr-review.md"))) {
		t.Fatalf("expected command source from env override, got: %s", outCfg)
	}
}

func TestForge_E2E_DryRunDoesNotMutateAndPrintsActionableSummary(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	out, err := runCLI(repo, root, []string{"forge", "--dry-run"}, nil)
	if err != nil {
		t.Fatalf("dry-run forge failed: %v\n%s", err, out)
	}

	if !strings.Contains(out, "[dry-run] forge: planned") || !strings.Contains(out, "rerun `weave forge` without --dry-run") {
		t.Fatalf("expected concise actionable dry-run summary, got: %s", out)
	}

	if _, err := os.Stat(filepath.Join(repo, ".agents", "skills")); !os.IsNotExist(err) {
		t.Fatalf("expected no .agents/skills created on dry-run, got err: %v", err)
	}

	if _, err := os.Stat(filepath.Join(repo, "weave.yaml")); !os.IsNotExist(err) {
		t.Fatalf("expected no weave.yaml created on dry-run, got err: %v", err)
	}
}

func TestSkillAdd_E2E_DryRunDoesNotMutateAndPrintsActionableSummary(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	skillsRoot := filepath.Join(repo, "skills-src")
	if err := os.MkdirAll(filepath.Join(skillsRoot, "sdd-orchestrator"), 0o755); err != nil {
		t.Fatalf("mkdir skills source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillsRoot, "sdd-orchestrator", "SKILL.md"), []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}

	cfg := []byte("version: 1\nsources:\n  skills_dir: " + skillsRoot + "\n  commands_dir: ~/.weave/commands\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"skill", "add", "sdd-orchestrator", "--dry-run"}, nil)
	if err != nil {
		t.Fatalf("dry-run skill add failed: %v\n%s", err, out)
	}

	if !strings.Contains(out, "[dry-run] skill add sdd-orchestrator") || !strings.Contains(out, "rerun without --dry-run") {
		t.Fatalf("expected concise actionable dry-run summary, got: %s", out)
	}

	if _, err := os.Lstat(filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("expected no symlink created on dry-run, got err: %v", err)
	}
}

func TestCommandAdd_E2E_DefaultAllEnabledProviders_ProjectsForEachProvider(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	commandsRoot := filepath.Join(repo, "commands-src")
	if err := os.MkdirAll(commandsRoot, 0o755); err != nil {
		t.Fatalf("mkdir commands source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandsRoot, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	cfg := []byte("version: 1\nproviders:\n  - name: claude-code\n    enabled: true\n  - name: codex\n    enabled: true\n  - name: opencode\n    enabled: true\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: " + commandsRoot + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"command", "add", "pr-review"}, nil)
	if err != nil {
		t.Fatalf("command add failed: %v\n%s", err, out)
	}

	for _, p := range []string{
		filepath.Join(repo, ".agents", "commands", "pr-review.md"),
		filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md"),
	} {
		fi, err := os.Lstat(p)
		if err != nil {
			t.Fatalf("expected projection path %s: %v", p, err)
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("expected %s to be symlink", p)
		}
	}
}

func TestCommandAdd_E2E_ProviderExclusiveCodex_DoesNotRequireSharedAgentsPath(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	commandsRoot := filepath.Join(repo, "commands-src")
	if err := os.MkdirAll(commandsRoot, 0o755); err != nil {
		t.Fatalf("mkdir commands source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandsRoot, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repo, ".agents"), []byte("blocked"), 0o644); err != nil {
		t.Fatalf("seed blocking .agents path: %v", err)
	}

	cfg := []byte("version: 1\nproviders: []\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: " + commandsRoot + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"command", "add", "pr-review", "--provider", "codex"}, nil)
	if err != nil {
		t.Fatalf("exclusive command add failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(repo, ".agents", "commands", "pr-review.md")); err == nil {
		t.Fatalf("expected no shared command path for exclusive install")
	}

	if fi, err := os.Lstat(filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md")); err != nil {
		t.Fatalf("expected codex wrapper projection: %v", err)
	} else if fi.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("expected codex wrapper projection to be symlink")
	}
}

func TestCommandAdd_E2E_NoProvidersInteractivePrompt_DefaultNo(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	commandsRoot := filepath.Join(repo, "commands-src")
	if err := os.MkdirAll(commandsRoot, 0o755); err != nil {
		t.Fatalf("mkdir commands source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandsRoot, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	cfg := []byte("version: 1\nproviders: []\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: " + commandsRoot + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"command", "add", "pr-review"}, []string{"WEAVE_FORCE_INTERACTIVE=1", "WEAVE_TEST_STDIN=n\n"})
	if err == nil {
		t.Fatalf("expected interactive no-provider default-no flow to fail")
	}
	if !strings.Contains(out, "No providers are currently enabled. Continue anyway? [y/N]:") {
		t.Fatalf("expected explicit no-provider prompt, got: %s", out)
	}
	if !strings.Contains(out, "command add canceled by user") {
		t.Fatalf("expected canceled-by-user message, got: %s", out)
	}
}

func TestCommandAdd_E2E_NoProvidersInteractivePrompt_EmptyInputDefaultsToNo(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	commandsRoot := filepath.Join(repo, "commands-src")
	if err := os.MkdirAll(commandsRoot, 0o755); err != nil {
		t.Fatalf("mkdir commands source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandsRoot, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	cfg := []byte("version: 1\nproviders: []\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: " + commandsRoot + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"command", "add", "pr-review"}, []string{"WEAVE_FORCE_INTERACTIVE=1", "WEAVE_TEST_STDIN=\n"})
	if err == nil {
		t.Fatalf("expected interactive no-provider empty-input flow to fail")
	}
	if !strings.Contains(out, "No providers are currently enabled. Continue anyway? [y/N]:") {
		t.Fatalf("expected explicit no-provider prompt, got: %s", out)
	}
	if !strings.Contains(out, "command add canceled by user") {
		t.Fatalf("expected canceled-by-user message, got: %s", out)
	}
}

func TestCommandAdd_E2E_NoProvidersNonInteractiveFailsWithGuidance(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	commandsRoot := filepath.Join(repo, "commands-src")
	if err := os.MkdirAll(commandsRoot, 0o755); err != nil {
		t.Fatalf("mkdir commands source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandsRoot, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	seed := []byte("version: 1\nproviders: []\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: " + commandsRoot + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"command", "add", "pr-review"}, []string{"WEAVE_NON_INTERACTIVE=1"})
	if err == nil {
		t.Fatalf("expected non-interactive no-provider flow to fail")
	}
	if !strings.Contains(out, "No providers are currently enabled") || !strings.Contains(out, "--provider") {
		t.Fatalf("expected actionable non-interactive guidance, got: %s", out)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged on non-interactive no-provider failure")
	}
}

func TestCommandAdd_E2E_MultiProviderFailureRollsBackFilesystemAndConfig(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	commandsRoot := filepath.Join(repo, "commands-src")
	if err := os.MkdirAll(commandsRoot, 0o755); err != nil {
		t.Fatalf("mkdir commands source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandsRoot, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	seed := []byte("version: 1\nproviders:\n  - name: claude-code\n    enabled: true\n  - name: codex\n    enabled: true\n  - name: opencode\n    enabled: true\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: " + commandsRoot + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), seed, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	blockedProjectionRoot := filepath.Join(repo, ".codex", "commands", "__weave_commands__")
	if err := os.MkdirAll(blockedProjectionRoot, 0o755); err != nil {
		t.Fatalf("mkdir blocked projection root: %v", err)
	}
	if err := os.Chmod(blockedProjectionRoot, 0o555); err != nil {
		t.Fatalf("chmod blocked projection root: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(blockedProjectionRoot, 0o755) })

	out, err := runCLI(repo, root, []string{"command", "add", "pr-review"}, nil)
	if err == nil {
		t.Fatalf("expected command add to fail when provider projection apply fails")
	}
	if !strings.Contains(out, "rollback completed so no config or symlink changes were committed") {
		t.Fatalf("expected full rollback failure guidance, got: %s", out)
	}

	if _, err := os.Lstat(filepath.Join(repo, ".agents", "commands", "pr-review.md")); !os.IsNotExist(err) {
		t.Fatalf("expected shared command projection to be rolled back, got err: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(repo, ".codex", "commands", "__weave_commands__", "pr-review", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("expected codex command projection to be absent after rollback, got err: %v", err)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	if string(after) != string(seed) {
		t.Fatalf("expected weave.yaml unchanged after multi-provider rollback")
	}
}

func TestCommandAdd_E2E_MultiProviderSuccessPersistsTransactionalMetadata(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", ".."))

	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	commandsRoot := filepath.Join(repo, "commands-src")
	if err := os.MkdirAll(commandsRoot, 0o755); err != nil {
		t.Fatalf("mkdir commands source: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandsRoot, "pr-review.md"), []byte("# command"), 0o644); err != nil {
		t.Fatalf("write command source: %v", err)
	}

	cfg := []byte("version: 1\nproviders:\n  - name: claude-code\n    enabled: true\n  - name: codex\n    enabled: true\n  - name: opencode\n    enabled: true\nsources:\n  skills_dir: ~/.weave/skills\n  commands_dir: " + commandsRoot + "\nsync:\n  mode: symlink\nskills: []\ncommands: []\n")
	if err := os.WriteFile(filepath.Join(repo, "weave.yaml"), cfg, 0o644); err != nil {
		t.Fatalf("write weave.yaml: %v", err)
	}

	out, err := runCLI(repo, root, []string{"command", "add", "pr-review"}, nil)
	if err != nil {
		t.Fatalf("command add failed: %v\n%s", err, out)
	}

	after, err := os.ReadFile(filepath.Join(repo, "weave.yaml"))
	if err != nil {
		t.Fatalf("read weave.yaml: %v", err)
	}
	content := string(after)
	if !strings.Contains(content, "provider_compat:") || !strings.Contains(content, "shared_install: true") {
		t.Fatalf("expected transactional command metadata in config, got: %s", content)
	}
}

func runCLI(repo string, root string, args []string, extraEnv []string) (string, error) {
	cmd := exec.Command("go", "run", "./cmd/weave")
	if len(args) > 0 {
		cmd = exec.Command("go", append([]string{"run", "./cmd/weave"}, args...)...)
	}
	cmd.Dir = root
	baseEnv := append(os.Environ(), "WEAVE_WORKDIR="+repo)
	filteredEnv := make([]string, 0, len(extraEnv))
	stdinPayload := ""
	for _, kv := range extraEnv {
		if strings.HasPrefix(kv, "WEAVE_TEST_STDIN=") {
			stdinPayload = strings.TrimPrefix(kv, "WEAVE_TEST_STDIN=")
			continue
		}
		filteredEnv = append(filteredEnv, kv)
	}
	cmd.Env = append(baseEnv, filteredEnv...)
	if stdinPayload != "" {
		cmd.Stdin = bytes.NewBufferString(stdinPayload)
	}
	out, err := cmd.CombinedOutput()
	return string(out), err
}
