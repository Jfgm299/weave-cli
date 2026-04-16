package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Jfgm299/weave-cli/internal/config"
	"github.com/Jfgm299/weave-cli/internal/fsops"
)

type DoctorStatus string

const (
	DoctorStatusHealthy     DoctorStatus = "healthy"
	DoctorStatusIssuesFound DoctorStatus = "issues_found"
)

type DoctorIssue struct {
	Code          string `json:"code"`
	Summary       string `json:"summary"`
	RepairCommand string `json:"repair_command,omitempty"`
	DocsPath      string `json:"docs_path"`
	DocsURL       string `json:"docs_url"`
}

type DoctorResult struct {
	Status         DoctorStatus  `json:"status"`
	Issues         []DoctorIssue `json:"issues"`
	RepairCommands []string      `json:"repair_commands"`
}

type DoctorService struct {
	ProviderRegistry ProviderRegistry
}

func (d DoctorService) Run(_ context.Context, projectRoot string, cfg config.Config) (DoctorResult, error) {
	return d.run(projectRoot, cfg)
}

func (d DoctorService) run(projectRoot string, cfg config.Config) (DoctorResult, error) {
	issues := make([]DoctorIssue, 0)

	if err := (config.Validator{}).Validate(cfg); err != nil {
		issues = append(issues, DoctorIssue{
			Code:          "invalid_config",
			Summary:       fmt.Sprintf("weave.yaml is invalid: %v", err),
			RepairCommand: "Fix weave.yaml and re-run `weave doctor`",
			DocsPath:      DocsPathConfig,
			DocsURL:       DocsURL(DocsPathConfig),
		})
	}

	issues = append(issues, diagnoseAssets(projectRoot, cfg.Skills, true)...)
	issues = append(issues, diagnoseAssets(projectRoot, cfg.Commands, false)...)
	issues = append(issues, d.diagnoseProviderProjections(projectRoot, cfg)...)

	repairs := uniqueSortedRepairCommands(issues)
	if len(issues) == 0 {
		return DoctorResult{Status: DoctorStatusHealthy, Issues: []DoctorIssue{}, RepairCommands: []string{}}, nil
	}

	return DoctorResult{Status: DoctorStatusIssuesFound, Issues: issues, RepairCommands: repairs}, nil
}

func (d DoctorService) diagnoseProviderProjections(projectRoot string, cfg config.Config) []DoctorIssue {
	if d.ProviderRegistry == nil {
		return nil
	}

	enabled := ListEnabledProviders(cfg)
	issues := make([]DoctorIssue, 0)
	for _, provider := range enabled {
		adapter, ok := d.ProviderRegistry.Get(provider.Name)
		if !ok {
			issues = append(issues, DoctorIssue{
				Code:          "unknown_enabled_provider",
				Summary:       fmt.Sprintf("enabled provider %q is not supported by this weave binary", provider.Name),
				RepairCommand: fmt.Sprintf("Remove %q from weave.yaml or upgrade weave and run `weave provider repair %s`", provider.Name, provider.Name),
				DocsPath:      DocsPathProviders,
				DocsURL:       DocsURL(DocsPathProviders),
			})
			continue
		}

		opts, err := adapter.PlanSetup(projectRoot)
		if err != nil {
			issues = append(issues, DoctorIssue{
				Code:          "stale_provider_integration",
				Summary:       fmt.Sprintf("cannot evaluate provider %q projection health: %v", provider.Name, err),
				RepairCommand: fmt.Sprintf("weave provider repair %s", provider.Name),
				DocsPath:      DocsPathProviders,
				DocsURL:       DocsURL(DocsPathProviders),
			})
			continue
		}

		for _, op := range opts {
			if op.Type != fsops.OpCreateLink {
				continue
			}
			fi, err := os.Lstat(op.Path)
			if os.IsNotExist(err) {
				issues = append(issues, DoctorIssue{
					Code:          "stale_provider_integration",
					Summary:       fmt.Sprintf("provider %q projection is stale: missing %s", provider.Name, op.Path),
					RepairCommand: fmt.Sprintf("weave provider repair %s", provider.Name),
					DocsPath:      DocsPathProviders,
					DocsURL:       DocsURL(DocsPathProviders),
				})
				continue
			}
			if err != nil {
				issues = append(issues, DoctorIssue{
					Code:          "stale_provider_integration",
					Summary:       fmt.Sprintf("provider %q projection is unreadable at %s: %v", provider.Name, op.Path, err),
					RepairCommand: fmt.Sprintf("weave provider repair %s", provider.Name),
					DocsPath:      DocsPathProviders,
					DocsURL:       DocsURL(DocsPathProviders),
				})
				continue
			}
			if fi.Mode()&os.ModeSymlink == 0 {
				issues = append(issues, DoctorIssue{
					Code:          "stale_provider_integration",
					Summary:       fmt.Sprintf("provider %q projection at %s is not a symlink", provider.Name, op.Path),
					RepairCommand: fmt.Sprintf("weave provider repair %s", provider.Name),
					DocsPath:      DocsPathProviders,
					DocsURL:       DocsURL(DocsPathProviders),
				})
				continue
			}
			target, err := os.Readlink(op.Path)
			if err != nil {
				issues = append(issues, DoctorIssue{
					Code:          "stale_provider_integration",
					Summary:       fmt.Sprintf("provider %q projection link is unreadable at %s: %v", provider.Name, op.Path, err),
					RepairCommand: fmt.Sprintf("weave provider repair %s", provider.Name),
					DocsPath:      DocsPathProviders,
					DocsURL:       DocsURL(DocsPathProviders),
				})
				continue
			}
			if filepath.Clean(target) != filepath.Clean(op.Target) {
				issues = append(issues, DoctorIssue{
					Code:          "stale_provider_integration",
					Summary:       fmt.Sprintf("provider %q projection at %s points to %s but expected %s", provider.Name, op.Path, target, op.Target),
					RepairCommand: fmt.Sprintf("weave provider repair %s", provider.Name),
					DocsPath:      DocsPathProviders,
					DocsURL:       DocsURL(DocsPathProviders),
				})
			}
		}
	}

	return issues
}

func diagnoseAssets(projectRoot string, assets []config.Asset, skills bool) []DoctorIssue {
	out := make([]DoctorIssue, 0)
	for _, asset := range assets {
		installedPath := installedAssetPath(projectRoot, asset.Name, skills)
		fi, err := os.Lstat(installedPath)
		if os.IsNotExist(err) {
			out = append(out, DoctorIssue{
				Code:          "missing_symlink",
				Summary:       fmt.Sprintf("missing installed %s %q at %s", assetKindName(skills), asset.Name, installedPath),
				RepairCommand: doctorRepairCommand(asset.Name, skills),
				DocsPath:      DocsPathDoctor,
				DocsURL:       DocsURL(DocsPathDoctor),
			})
			continue
		}
		if err != nil {
			out = append(out, DoctorIssue{
				Code:          "unreadable_asset_path",
				Summary:       fmt.Sprintf("cannot inspect %s %q at %s: %v", assetKindName(skills), asset.Name, installedPath, err),
				RepairCommand: doctorRepairCommand(asset.Name, skills),
				DocsPath:      DocsPathDoctor,
				DocsURL:       DocsURL(DocsPathDoctor),
			})
			continue
		}

		if fi.Mode()&os.ModeSymlink == 0 {
			out = append(out, DoctorIssue{
				Code:          "not_a_symlink",
				Summary:       fmt.Sprintf("installed %s %q is not a symlink: %s", assetKindName(skills), asset.Name, installedPath),
				RepairCommand: doctorRepairCommand(asset.Name, skills),
				DocsPath:      DocsPathTransactions,
				DocsURL:       DocsURL(DocsPathTransactions),
			})
			continue
		}

		target, err := os.Readlink(installedPath)
		if err != nil {
			out = append(out, DoctorIssue{
				Code:          "broken_symlink",
				Summary:       fmt.Sprintf("cannot read symlink target for %s %q: %v", assetKindName(skills), asset.Name, err),
				RepairCommand: doctorRepairCommand(asset.Name, skills),
				DocsPath:      DocsPathDoctor,
				DocsURL:       DocsURL(DocsPathDoctor),
			})
			continue
		}

		actual := target
		if !filepath.IsAbs(actual) {
			actual = filepath.Clean(filepath.Join(filepath.Dir(installedPath), actual))
		}

		expected := filepath.Clean(asset.Source)
		if filepath.Clean(actual) != expected {
			out = append(out, DoctorIssue{
				Code:          "symlink_target_mismatch",
				Summary:       fmt.Sprintf("installed %s %q points to %s but weave.yaml expects %s", assetKindName(skills), asset.Name, actual, expected),
				RepairCommand: doctorRepairCommand(asset.Name, skills),
				DocsPath:      DocsPathDoctor,
				DocsURL:       DocsURL(DocsPathDoctor),
			})
		}
	}

	return out
}

func installedAssetPath(projectRoot, name string, skills bool) string {
	if skills {
		return filepath.Join(projectRoot, ".agents", "skills", name, "SKILL.md")
	}
	return filepath.Join(projectRoot, ".agents", "commands", fmt.Sprintf("%s.md", name))
}

func doctorRepairCommand(name string, skills bool) string {
	if skills {
		return fmt.Sprintf("weave skill add %s", name)
	}
	return fmt.Sprintf("weave command add %s", name)
}

func assetKindName(skills bool) string {
	if skills {
		return "skill"
	}
	return "command"
}

func uniqueSortedRepairCommands(issues []DoctorIssue) []string {
	seen := map[string]struct{}{}
	for _, issue := range issues {
		if issue.RepairCommand == "" {
			continue
		}
		seen[issue.RepairCommand] = struct{}{}
	}

	out := make([]string, 0, len(seen))
	for cmd := range seen {
		out = append(out, cmd)
	}
	sort.Strings(out)
	return out
}
