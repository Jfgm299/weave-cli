package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestDoctorService_Run_HealthyProjectReturnsNoIssues(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	source := filepath.Join(repo, "shared", "skills", "sdd-orchestrator", "SKILL.md")
	installed := filepath.Join(repo, ".agents", "skills", "sdd-orchestrator", "SKILL.md")

	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.WriteFile(source, []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(installed), 0o755); err != nil {
		t.Fatalf("mkdir installed: %v", err)
	}
	if err := os.Symlink(source, installed); err != nil {
		t.Fatalf("seed symlink: %v", err)
	}

	svc := DoctorService{}
	result, err := svc.Run(context.Background(), repo, config.Config{
		Version: 1,
		Sync:    config.Sync{Mode: "symlink"},
		Skills:  []config.Asset{{Name: "sdd-orchestrator", Source: source}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != DoctorStatusHealthy {
		t.Fatalf("expected healthy status, got %q", result.Status)
	}

	if len(result.Issues) != 0 {
		t.Fatalf("expected no issues, got %+v", result.Issues)
	}

	if len(result.RepairCommands) != 0 {
		t.Fatalf("expected no repair commands, got %+v", result.RepairCommands)
	}
}

func TestDoctorService_Run_MissingSkillSymlinkReturnsRepairGuidance(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	source := filepath.Join(repo, "shared", "skills", "sdd-orchestrator", "SKILL.md")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("mkdir source: %v", err)
	}
	if err := os.WriteFile(source, []byte("# skill"), 0o644); err != nil {
		t.Fatalf("write source: %v", err)
	}

	svc := DoctorService{}
	result, err := svc.Run(context.Background(), repo, config.Config{
		Version: 1,
		Sync:    config.Sync{Mode: "symlink"},
		Skills:  []config.Asset{{Name: "sdd-orchestrator", Source: source}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != DoctorStatusIssuesFound {
		t.Fatalf("expected issues status, got %q", result.Status)
	}

	if len(result.Issues) == 0 {
		t.Fatalf("expected at least one issue")
	}

	if result.Issues[0].DocsPath == "" || result.Issues[0].DocsURL == "" {
		t.Fatalf("expected docs references in issue, got %+v", result.Issues[0])
	}

	if len(result.RepairCommands) == 0 || !strings.Contains(result.RepairCommands[0], "weave skill add sdd-orchestrator") {
		t.Fatalf("expected repair guidance for missing skill symlink, got %+v", result.RepairCommands)
	}
}

func TestDoctorService_Run_UnknownEnabledProviderFlagsStaleIntegration(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()

	svc := DoctorService{ProviderRegistry: doctorRegistryStub{}}
	result, err := svc.Run(context.Background(), repo, config.Config{
		Version:   1,
		Sync:      config.Sync{Mode: "symlink"},
		Providers: []config.Provider{{Name: "future-provider", Enabled: true}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != DoctorStatusIssuesFound {
		t.Fatalf("expected issues status, got %q", result.Status)
	}

	found := false
	for _, issue := range result.Issues {
		if issue.Code == "unknown_enabled_provider" {
			found = true
			if !strings.Contains(issue.RepairCommand, "upgrade weave") {
				t.Fatalf("expected upgrade/repair guidance, got %q", issue.RepairCommand)
			}
			break
		}
	}
	if !found {
		t.Fatalf("expected unknown_enabled_provider issue, got %+v", result.Issues)
	}
}

type doctorRegistryStub struct{}

func (doctorRegistryStub) Get(name string) (ProviderAdapter, bool) {
	if name == "claude-code" {
		return doctorAdapterStub{name: "claude-code"}, true
	}
	return nil, false
}

func (doctorRegistryStub) SupportedNames() []string {
	return []string{"claude-code"}
}

type doctorAdapterStub struct{ name string }

func (a doctorAdapterStub) Name() string { return a.name }

func (doctorAdapterStub) RequiredBinaries() []string { return nil }

func (doctorAdapterStub) PlanSetup(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{{Type: fsops.OpCreateLink, Path: filepath.Join(projectRoot, ".claude", "CLAUDE.md"), Target: filepath.Join("..", ".agents", "AGENTS.md")}}, nil
}

func (doctorAdapterStub) PlanRepair(projectRoot string) ([]fsops.Operation, error) {
	return doctorAdapterStub{}.PlanSetup(projectRoot)
}

func (doctorAdapterStub) PlanRemove(projectRoot string) ([]fsops.Operation, error) {
	return []fsops.Operation{{Type: fsops.OpRemovePath, Path: filepath.Join(projectRoot, ".claude")}}, nil
}
