package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jfgm299/weave-cli/internal/config"
)

type assetKind string

const (
	assetKindSkill   assetKind = "skill"
	assetKindCommand assetKind = "command"
)

type sourceResolver struct {
	lookupEnv func(string) (string, bool)
	homeDir   func() (string, error)
}

func newSourceResolver() sourceResolver {
	return sourceResolver{
		lookupEnv: os.LookupEnv,
		homeDir:   os.UserHomeDir,
	}
}

func (r sourceResolver) Resolve(kind assetKind, fromFlag string, cfg config.Config) (string, error) {
	if strings.TrimSpace(fromFlag) != "" {
		return expandHome(fromFlag, r.homeDir)
	}

	envKey := ""
	cfgValue := ""
	defaultValue := ""

	switch kind {
	case assetKindSkill:
		envKey = "WEAVE_SKILLS_DIR"
		cfgValue = cfg.Sources.SkillsDir
		defaultValue = "~/.weave/skills"
	case assetKindCommand:
		envKey = "WEAVE_COMMANDS_DIR"
		cfgValue = cfg.Sources.CommandsDir
		defaultValue = "~/.weave/commands"
	default:
		return "", fmt.Errorf("unsupported asset kind: %s", kind)
	}

	if raw, ok := r.lookupEnv(envKey); ok && strings.TrimSpace(raw) != "" {
		return expandHome(raw, r.homeDir)
	}

	if strings.TrimSpace(cfgValue) != "" {
		return expandHome(cfgValue, r.homeDir)
	}

	return expandHome(defaultValue, r.homeDir)
}

func expandHome(raw string, homeDirFn func() (string, error)) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil
	}

	if raw == "~" {
		home, err := homeDirFn()
		if err != nil {
			return "", err
		}
		return home, nil
	}

	if strings.HasPrefix(raw, "~/") {
		home, err := homeDirFn()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, strings.TrimPrefix(raw, "~/")), nil
	}

	return raw, nil
}
