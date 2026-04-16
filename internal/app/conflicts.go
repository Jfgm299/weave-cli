package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type ConflictPolicy string

const (
	ConflictPolicyPrompt    ConflictPolicy = "prompt"
	ConflictPolicyOverwrite ConflictPolicy = "overwrite"
	ConflictPolicySkip      ConflictPolicy = "skip"
	ConflictPolicyBackup    ConflictPolicy = "backup"
)

var (
	ErrConflictPromptRequired = errors.New("conflict prompt is required")
)

type ConflictPromptInput struct {
	Kind string
	Name string
	Path string
}

type ConflictPrompter interface {
	ResolveConflict(ctx context.Context, input ConflictPromptInput) (ConflictPolicy, error)
}

type conflictPlanner struct {
	now func() time.Time
}

func (p conflictPlanner) plan(path string, source string, policy ConflictPolicy) ([]fsops.Operation, bool, error) {
	resolved, err := normalizeConflictPolicy(policy)
	if err != nil {
		return nil, false, err
	}

	switch resolved {
	case ConflictPolicySkip:
		return []fsops.Operation{}, true, nil
	case ConflictPolicyOverwrite:
		return []fsops.Operation{{Type: fsops.OpCreateLink, Path: path, Target: source}}, false, nil
	case ConflictPolicyBackup:
		ts := p.nowFn().UTC().Format("20060102T150405Z")
		backupPath := fmt.Sprintf("%s.bak.%s", path, ts)
		return []fsops.Operation{
			{Type: fsops.OpBackupPath, Path: path, Target: backupPath},
			{Type: fsops.OpCreateLink, Path: path, Target: source},
		}, false, nil
	default:
		return nil, false, fmt.Errorf("unsupported conflict policy: %s", resolved)
	}
}

func (p conflictPlanner) nowFn() time.Time {
	if p.now != nil {
		return p.now()
	}
	return time.Now()
}

func normalizeConflictPolicy(policy ConflictPolicy) (ConflictPolicy, error) {
	v := ConflictPolicy(strings.TrimSpace(string(policy)))
	if v == "" {
		return ConflictPolicyPrompt, nil
	}
	switch v {
	case ConflictPolicyPrompt, ConflictPolicyOverwrite, ConflictPolicySkip, ConflictPolicyBackup:
		return v, nil
	default:
		return "", fmt.Errorf("unsupported conflict policy: %s", policy)
	}
}

func classifyAssetKind(kind AssetKind) string {
	if kind == AssetKindCommand {
		return "command"
	}
	return "skill"
}

func hasConflictAtPath(path string) (bool, error) {
	if path == "" {
		return false, nil
	}
	if _, err := filepath.Abs(path); err != nil {
		return false, err
	}
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type defaultPathChecker struct{}

func (defaultPathChecker) Exists(path string) (bool, error) {
	return hasConflictAtPath(path)
}

func resolveConflictPolicy(ctx context.Context, policy ConflictPolicy, prompter ConflictPrompter, input ConflictPromptInput) (ConflictPolicy, error) {
	resolved, err := normalizeConflictPolicy(policy)
	if err != nil {
		return "", err
	}
	if resolved != ConflictPolicyPrompt {
		return resolved, nil
	}
	if prompter == nil {
		return "", fmt.Errorf("%w for %s %q at %s. Re-run with one of --overwrite, --skip, --backup", ErrConflictPromptRequired, input.Kind, input.Name, input.Path)
	}
	choice, err := prompter.ResolveConflict(ctx, input)
	if err != nil {
		return "", err
	}
	return normalizeConflictPolicy(choice)
}
