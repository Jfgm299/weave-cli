# Providers Reference

Providers are integrated via adapter projections from canonical `.agents/` content.

## Supported providers

- `claude-code`
- `codex`
- `opencode`

All provider projections remain **symlink-based** from canonical `.agents/` content:

- `.claude/*` -> `../.agents/*`
- `.codex/*` -> `../.agents/*`
- `.opencode/*` -> `../.agents/*`

## Commands

```sh
weave provider add <name>
weave provider remove <name>
weave provider repair <name>
weave provider list
```

Provider-aware command install behavior:

```sh
weave command add <name>
weave command add <name> --provider <name>
```

- Default (`weave command add <name>`):
  - installs canonical shared command at `.agents/commands/<name>.md`
  - projects to all enabled providers.
- Exclusive (`--provider <name>`):
  - installs only the target provider projection
  - does not require shared `.agents/commands` install.

Codex command projection uses wrapper namespace:

- `.codex/commands/__weave_commands__/<name>/SKILL.md`

No-provider behavior:

- Interactive sessions prompt: `No providers are currently enabled. Continue anyway? [y/N]:`
- Non-interactive sessions fail automatically with actionable guidance.

## Dependency checks

Provider binaries are validated before setup or repair is considered successful.

Required binaries:

- `claude-code` -> `claude`
- `codex` -> `codex`
- `opencode` -> `opencode`

If a binary is missing:

1. install the missing binary,
2. run `weave provider repair <name>`,
3. verify with `weave doctor`.
