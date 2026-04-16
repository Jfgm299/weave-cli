package config

type Config struct {
	Version   int        `yaml:"version"`
	Providers []Provider `yaml:"providers"`
	Sources   Sources    `yaml:"sources"`
	Sync      Sync       `yaml:"sync"`
	Skills    []Asset    `yaml:"skills"`
	Commands  []Asset    `yaml:"commands"`
}

type Provider struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
}

type Sources struct {
	SkillsDir   string `yaml:"skills_dir"`
	CommandsDir string `yaml:"commands_dir"`
}

type Sync struct {
	Mode           string `yaml:"mode"`
	ConflictPolicy string `yaml:"conflict_policy,omitempty"`
}

type Asset struct {
	Name   string         `yaml:"name"`
	Source string         `yaml:"source"`
	Meta   *CommandMetaV1 `yaml:"metadata,omitempty"`
}

type CommandMetaV1 struct {
	ProviderCompat []string `yaml:"provider_compat,omitempty"`
}
