package config

func Default() Config {
	return Config{
		Version: 1,
		Sources: Sources{
			SkillsDir:   "~/.weave/skills",
			CommandsDir: "~/.weave/commands",
		},
		Sync: Sync{
			Mode:           "symlink",
			ConflictPolicy: "prompt",
		},
		Providers: []Provider{},
		Skills:    []Asset{},
		Commands:  []Asset{},
	}
}
