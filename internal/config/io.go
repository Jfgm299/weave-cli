package config

import (
	"bytes"
	"fmt"
	"sort"

	"gopkg.in/yaml.v3"
)

func MarshalDeterministic(cfg Config) ([]byte, error) {
	copyCfg := cfg

	sort.Slice(copyCfg.Skills, func(i, j int) bool {
		return copyCfg.Skills[i].Name < copyCfg.Skills[j].Name
	})

	sort.Slice(copyCfg.Commands, func(i, j int) bool {
		return copyCfg.Commands[i].Name < copyCfg.Commands[j].Name
	})

	sort.Slice(copyCfg.Providers, func(i, j int) bool {
		return copyCfg.Providers[i].Name < copyCfg.Providers[j].Name
	})

	if copyCfg.Version == 0 {
		return nil, fmt.Errorf("version is required")
	}

	if copyCfg.Sync.Mode == "" {
		copyCfg.Sync.Mode = "symlink"
	}

	if copyCfg.Sources.SkillsDir == "" {
		copyCfg.Sources.SkillsDir = "~/.weave/skills"
	}

	if copyCfg.Sources.CommandsDir == "" {
		copyCfg.Sources.CommandsDir = "~/.weave/commands"
	}

	if copyCfg.Skills == nil {
		copyCfg.Skills = []Asset{}
	}

	if copyCfg.Commands == nil {
		copyCfg.Commands = []Asset{}
	}

	if copyCfg.Providers == nil {
		copyCfg.Providers = []Provider{}
	}

	if err := validateUniqueNames(copyCfg.Skills, "skills"); err != nil {
		return nil, err
	}

	if err := validateUniqueNames(copyCfg.Commands, "commands"); err != nil {
		return nil, err
	}

	normalizeAssetMetadata(copyCfg.Commands)
	normalizeAssetMetadata(copyCfg.Skills)

	if err := validateUniqueProviders(copyCfg.Providers); err != nil {
		return nil, err
	}

	b, err := yaml.Marshal(copyCfg)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(b), nil
}

func validateUniqueProviders(providers []Provider) error {
	seen := make(map[string]struct{}, len(providers))
	for _, provider := range providers {
		if _, ok := seen[provider.Name]; ok {
			return fmt.Errorf("duplicate providers inventory entry: %s", provider.Name)
		}
		seen[provider.Name] = struct{}{}
	}
	return nil
}

func validateUniqueNames(assets []Asset, field string) error {
	seen := make(map[string]struct{}, len(assets))
	for _, asset := range assets {
		if _, ok := seen[asset.Name]; ok {
			return fmt.Errorf("duplicate %s inventory entry: %s", field, asset.Name)
		}
		seen[asset.Name] = struct{}{}
	}
	return nil
}

func normalizeAssetMetadata(assets []Asset) {
	for i := range assets {
		if assets[i].Meta == nil {
			continue
		}
		if len(assets[i].Meta.ProviderCompat) == 0 {
			assets[i].Meta = nil
			continue
		}
		sorted := append([]string(nil), assets[i].Meta.ProviderCompat...)
		sort.Strings(sorted)
		assets[i].Meta.ProviderCompat = uniqueNonEmpty(sorted)
		if len(assets[i].Meta.ProviderCompat) == 0 {
			assets[i].Meta = nil
		}
	}
}

func uniqueNonEmpty(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
