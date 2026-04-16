package cli

import (
	"testing"

	"github.com/Jfgm299/weave-cli/internal/config"
)

func TestSourceResolver_Resolve_FlagOverridesEnvConfigAndDefault(t *testing.T) {
	t.Parallel()

	r := sourceResolver{
		lookupEnv: func(_ string) (string, bool) { return "/env/skills", true },
		homeDir:   func() (string, error) { return "/home/dev", nil },
	}

	cfg := config.Config{Sources: config.Sources{SkillsDir: "/cfg/skills"}}
	g, err := r.Resolve(assetKindSkill, "/flag/skills", cfg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if g != "/flag/skills" {
		t.Fatalf("expected flag source, got %q", g)
	}
}

func TestSourceResolver_Resolve_EnvOverridesConfigAndDefault(t *testing.T) {
	t.Parallel()

	r := sourceResolver{
		lookupEnv: func(key string) (string, bool) {
			if key == "WEAVE_SKILLS_DIR" {
				return "/env/skills", true
			}
			return "", false
		},
		homeDir: func() (string, error) { return "/home/dev", nil },
	}

	cfg := config.Config{Sources: config.Sources{SkillsDir: "/cfg/skills"}}
	g, err := r.Resolve(assetKindSkill, "", cfg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if g != "/env/skills" {
		t.Fatalf("expected env source, got %q", g)
	}
}

func TestSourceResolver_Resolve_ConfigOverridesDefault(t *testing.T) {
	t.Parallel()

	r := sourceResolver{
		lookupEnv: func(_ string) (string, bool) { return "", false },
		homeDir:   func() (string, error) { return "/home/dev", nil },
	}

	cfg := config.Config{Sources: config.Sources{CommandsDir: "/cfg/commands"}}
	g, err := r.Resolve(assetKindCommand, "", cfg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if g != "/cfg/commands" {
		t.Fatalf("expected config source, got %q", g)
	}
}

func TestSourceResolver_Resolve_CommandEnvOverridesConfigAndDefault(t *testing.T) {
	t.Parallel()

	r := sourceResolver{
		lookupEnv: func(key string) (string, bool) {
			if key == "WEAVE_COMMANDS_DIR" {
				return "/env/commands", true
			}
			return "", false
		},
		homeDir: func() (string, error) { return "/home/dev", nil },
	}

	cfg := config.Config{Sources: config.Sources{CommandsDir: "/cfg/commands"}}
	g, err := r.Resolve(assetKindCommand, "", cfg)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if g != "/env/commands" {
		t.Fatalf("expected env source, got %q", g)
	}
}

func TestSourceResolver_Resolve_DefaultWhenNoOverrides(t *testing.T) {
	t.Parallel()

	r := sourceResolver{
		lookupEnv: func(_ string) (string, bool) { return "", false },
		homeDir:   func() (string, error) { return "/home/dev", nil },
	}

	g, err := r.Resolve(assetKindCommand, "", config.Default())
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if g != "/home/dev/.weave/commands" {
		t.Fatalf("expected default command source, got %q", g)
	}
}
