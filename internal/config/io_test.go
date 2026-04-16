package config

import (
	"bytes"
	"strings"
	"testing"
)

func TestMarshalDeterministic_SortsInventoryAndStabilizesOutput(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Version: 1,
		Sync:    Sync{Mode: "symlink"},
		Skills: []Asset{
			{Name: "zeta", Source: "z"},
			{Name: "alpha", Source: "a"},
		},
		Commands: []Asset{
			{Name: "run", Source: "r"},
			{Name: "build", Source: "b"},
		},
	}

	b1, err := MarshalDeterministic(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b2, err := MarshalDeterministic(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(b1, b2) {
		t.Fatalf("expected deterministic bytes, got diff")
	}
}

func TestMarshalDeterministic_MissingVersionFails(t *testing.T) {
	t.Parallel()

	_, err := MarshalDeterministic(Config{Sync: Sync{Mode: "symlink"}})
	if err == nil {
		t.Fatalf("expected error for missing version")
	}
}

func TestMarshalDeterministic_AlwaysEmitsExplicitInventory(t *testing.T) {
	t.Parallel()

	b, err := MarshalDeterministic(Config{Version: 1, Sync: Sync{Mode: "symlink"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(b)
	if !strings.Contains(out, "skills: []") {
		t.Fatalf("expected explicit skills inventory, got: %s", out)
	}

	if !strings.Contains(out, "commands: []") {
		t.Fatalf("expected explicit commands inventory, got: %s", out)
	}

	if !strings.Contains(out, "providers: []") {
		t.Fatalf("expected explicit providers inventory, got: %s", out)
	}

	if !strings.Contains(string(DefaultMustYAML(t)), "conflict_policy: prompt") {
		t.Fatalf("expected explicit default conflict_policy in defaults serialization")
	}

	if !strings.Contains(out, "mode: symlink") {
		t.Fatalf("expected explicit default conflict_policy, got: %s", out)
	}

	if !strings.Contains(out, "skills_dir: ~/.weave/skills") {
		t.Fatalf("expected explicit skills_dir source, got: %s", out)
	}

	if !strings.Contains(out, "commands_dir: ~/.weave/commands") {
		t.Fatalf("expected explicit commands_dir source, got: %s", out)
	}
}

func TestMarshalDeterministic_DuplicateProvidersFail(t *testing.T) {
	t.Parallel()

	_, err := MarshalDeterministic(Config{
		Version: 1,
		Sync:    Sync{Mode: "symlink"},
		Providers: []Provider{
			{Name: "claude-code", Enabled: true},
			{Name: "claude-code", Enabled: true},
		},
	})
	if err == nil {
		t.Fatalf("expected duplicate providers inventory to fail")
	}
}

func TestMarshalDeterministic_DuplicateInventoryFails(t *testing.T) {
	t.Parallel()

	_, err := MarshalDeterministic(Config{
		Version: 1,
		Sync:    Sync{Mode: "symlink"},
		Skills: []Asset{
			{Name: "dup", Source: "a"},
			{Name: "dup", Source: "b"},
		},
	})
	if err == nil {
		t.Fatalf("expected duplicate inventory to fail")
	}
}

func TestMarshalDeterministic_CommandMetadataProviderCompatIsNormalized(t *testing.T) {
	t.Parallel()

	b, err := MarshalDeterministic(Config{
		Version: 1,
		Sync:    Sync{Mode: "symlink"},
		Commands: []Asset{{
			Name:   "pr-review",
			Source: "/tmp/pr-review.md",
			Meta:   &CommandMetaV1{ProviderCompat: []string{"opencode", "", "claude-code", "opencode"}},
		}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(b)
	if !strings.Contains(out, "provider_compat") || !strings.Contains(out, "- claude-code") || !strings.Contains(out, "- opencode") {
		t.Fatalf("expected normalized provider_compat metadata, got: %s", out)
	}
}

func DefaultMustYAML(t *testing.T) string {
	t.Helper()
	b, err := MarshalDeterministic(Default())
	if err != nil {
		t.Fatalf("default marshal failed: %v", err)
	}
	return string(b)
}
