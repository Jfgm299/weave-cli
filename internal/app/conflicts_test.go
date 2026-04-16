package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestConflictPlanner_BackupPolicyProducesBackupAndCreateLinkOps(t *testing.T) {
	t.Parallel()

	planner := conflictPlanner{now: func() time.Time { return time.Date(2026, 4, 16, 12, 34, 56, 0, time.UTC) }}
	ops, skipped, err := planner.plan("/repo/.agents/skills/a/SKILL.md", "/src/a/SKILL.md", ConflictPolicyBackup)
	if err != nil {
		t.Fatalf("unexpected planner error: %v", err)
	}
	if skipped {
		t.Fatalf("backup policy must not skip")
	}
	if len(ops) != 2 || ops[0].Type != fsops.OpBackupPath || ops[1].Type != fsops.OpCreateLink {
		t.Fatalf("unexpected ops: %+v", ops)
	}
	if ops[0].Target == "" {
		t.Fatalf("expected backup target path")
	}
}

func TestResolveConflictPolicy_PromptWithoutPrompterFails(t *testing.T) {
	t.Parallel()

	_, err := resolveConflictPolicy(context.Background(), ConflictPolicyPrompt, nil, ConflictPromptInput{Kind: "skill", Name: "a", Path: "/repo/.agents/skills/a/SKILL.md"})
	if err == nil {
		t.Fatalf("expected prompt-required error")
	}
	if !errors.Is(err, ErrConflictPromptRequired) {
		t.Fatalf("expected ErrConflictPromptRequired, got %v", err)
	}
}

func TestResolveConflictPolicy_PrompterChoiceIsRespected(t *testing.T) {
	t.Parallel()

	resolved, err := resolveConflictPolicy(context.Background(), ConflictPolicyPrompt, promptStub{policy: ConflictPolicySkip}, ConflictPromptInput{Kind: "skill", Name: "a", Path: "/repo/.agents/skills/a/SKILL.md"})
	if err != nil {
		t.Fatalf("unexpected resolve error: %v", err)
	}
	if resolved != ConflictPolicySkip {
		t.Fatalf("expected skip policy from prompter, got %q", resolved)
	}
}

type promptStub struct{ policy ConflictPolicy }

func (s promptStub) ResolveConflict(_ context.Context, _ ConflictPromptInput) (ConflictPolicy, error) {
	return s.policy, nil
}
