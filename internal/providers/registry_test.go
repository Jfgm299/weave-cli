package providers

import (
	"path/filepath"
	"testing"

	"github.com/Jfgm299/weave-cli/internal/fsops"
)

func TestDefaultRegistry_Get_ReturnsKnownProviders(t *testing.T) {
	t.Parallel()

	r := NewDefaultRegistry()

	if _, ok := r.Get("claude-code"); !ok {
		t.Fatalf("expected claude-code adapter to be registered")
	}

	if _, ok := r.Get("opencode"); !ok {
		t.Fatalf("expected opencode adapter to be registered")
	}

	if _, ok := r.Get("unknown"); ok {
		t.Fatalf("expected unknown provider not to be registered")
	}
}

func TestClaudeAdapter_PlanSetup_CreatesProviderProjectionLinks(t *testing.T) {
	t.Parallel()

	a := ClaudeCodeAdapter{}
	root := "/tmp/repo"

	ops, err := a.PlanSetup(root)
	if err != nil {
		t.Fatalf("unexpected setup plan error: %v", err)
	}

	expected := []fsops.Operation{
		{Type: fsops.OpCreateLink, Path: filepath.Join(root, ".claude", "CLAUDE.md"), Target: filepath.Join("..", ".agents", "AGENTS.md")},
		{Type: fsops.OpCreateLink, Path: filepath.Join(root, ".claude", "commands"), Target: filepath.Join("..", ".agents", "commands")},
		{Type: fsops.OpCreateLink, Path: filepath.Join(root, ".claude", "docs"), Target: filepath.Join("..", ".agents", "docs")},
	}

	if len(ops) != len(expected) {
		t.Fatalf("expected %d operations, got %d", len(expected), len(ops))
	}

	for i := range expected {
		if ops[i].Type != expected[i].Type || ops[i].Path != expected[i].Path || ops[i].Target != expected[i].Target {
			t.Fatalf("unexpected op[%d]: got %+v want %+v", i, ops[i], expected[i])
		}
	}
}

func TestOpenCodeAdapter_RequiredBinaries_DeclaresOpencodeBinary(t *testing.T) {
	t.Parallel()

	a := OpenCodeAdapter{}
	bins := a.RequiredBinaries()
	if len(bins) != 1 || bins[0] != "opencode" {
		t.Fatalf("expected required binaries [opencode], got %+v", bins)
	}
}
