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
