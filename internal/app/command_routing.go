package app

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type CommandInstallPlan struct {
	ProviderTargets []string
	SharedInstall   bool
}

type CommandInstallPlanner struct {
	Registry ProviderRegistry
}

func (p CommandInstallPlanner) Resolve(providerFlag string, cfg config.Config, interactive bool, confirmNoProviders func() (bool, error)) (CommandInstallPlan, error) {
	if strings.TrimSpace(providerFlag) != "" {
		providerName := strings.TrimSpace(providerFlag)
		if p.Registry == nil {
			return CommandInstallPlan{}, fmt.Errorf("provider registry is not configured")
		}
		if _, ok := p.Registry.Get(providerName); !ok {
			supported := p.Registry.SupportedNames()
			sort.Strings(supported)
			return CommandInstallPlan{}, fmt.Errorf("unsupported provider: %s. Supported providers: %s. See %s (%s)", providerName, strings.Join(supported, ", "), DocsPathProviders, DocsURL(DocsPathProviders))
		}
		return CommandInstallPlan{ProviderTargets: []string{providerName}, SharedInstall: false}, nil
	}

	enabled := ListEnabledProviders(cfg)
	providerTargets := make([]string, 0, len(enabled))
	for _, provider := range enabled {
		providerTargets = append(providerTargets, provider.Name)
	}
	sort.Strings(providerTargets)

	if len(providerTargets) > 0 {
		return CommandInstallPlan{ProviderTargets: providerTargets, SharedInstall: true}, nil
	}

	if !interactive {
		return CommandInstallPlan{}, fmt.Errorf("No providers are currently enabled. Run `weave provider add <name>` first, or re-run with `weave command add <name> --provider <name>` for exclusive install")
	}
	if confirmNoProviders == nil {
		return CommandInstallPlan{}, fmt.Errorf("missing confirmation callback for no-provider interactive flow")
	}

	confirmed, err := confirmNoProviders()
	if err != nil {
		return CommandInstallPlan{}, err
	}
	if !confirmed {
		return CommandInstallPlan{}, fmt.Errorf("command add canceled by user because no providers are enabled")
	}

	return CommandInstallPlan{ProviderTargets: nil, SharedInstall: true}, nil
}

func BuildCommandProjectionOps(projectRoot string, command string, sourcePath string, plan CommandInstallPlan) ([]fsops.Operation, error) {
	ops := make([]fsops.Operation, 0, len(plan.ProviderTargets))
	for _, providerName := range plan.ProviderTargets {
		switch providerName {
		case "codex":
			ops = append(ops, fsops.Operation{
				Type:   fsops.OpCreateLink,
				Path:   filepath.Join(projectRoot, ".codex", "commands", "__weave_commands__", command, "SKILL.md"),
				Target: sourcePath,
			})
		case "claude-code":
			if plan.SharedInstall {
				continue
			}
			ops = append(ops, fsops.Operation{
				Type:   fsops.OpCreateLink,
				Path:   filepath.Join(projectRoot, ".claude", "commands", command+".md"),
				Target: sourcePath,
			})
		case "opencode":
			if plan.SharedInstall {
				continue
			}
			ops = append(ops, fsops.Operation{
				Type:   fsops.OpCreateLink,
				Path:   filepath.Join(projectRoot, ".opencode", "commands", command+".md"),
				Target: sourcePath,
			})
		default:
			return nil, fmt.Errorf("unsupported provider projection strategy: %s. See %s (%s)", providerName, DocsPathProviders, DocsURL(DocsPathProviders))
		}
	}

	return ops, nil
}
