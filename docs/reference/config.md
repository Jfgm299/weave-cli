# Configuration Reference

`weave.yaml` is the desired-state contract for project setup.

## Core fields

- `version`
- `providers`
- `skills`
- `commands`
- `sources.skills_dir`
- `sources.commands_dir`
- `sync.mode` (v1 requires `symlink`)

### Command metadata (v2 routing)

Each command entry can include deterministic routing metadata:

- `metadata.provider_compat`: provider targets included in the install transaction
- `metadata.shared_install`: whether shared `.agents/commands/<name>.md` is part of the transaction (`true` for default installs, `false` for exclusive `--provider` installs)

This metadata is used by `weave doctor` to validate shared/exclusive command projection integrity.

## Validation behavior

- Invalid or outdated schema returns actionable errors.
- `sync.mode` values other than `symlink` are rejected in v1.

## Migration

If schema is outdated, run:

```sh
weave migrate
```

Use dry-run first when needed:

```sh
weave migrate --dry-run
```
