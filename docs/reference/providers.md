# Providers Reference

Providers are integrated via adapter projections from canonical `.agents/` content.

## Supported v1 providers

- `claude-code`
- `opencode`

## Commands

```sh
weave provider add <name>
weave provider remove <name>
weave provider repair <name>
weave provider list
```

## Dependency checks

Provider binaries are validated before setup is considered successful.

If a binary is missing:

1. install the missing binary,
2. run `weave provider repair <name>`,
3. verify with `weave doctor`.
