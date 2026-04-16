# Migration Guide

## Schema migration baseline (v1)

When `weave.yaml` has an outdated schema version, run:

```sh
weave migrate
```

To preview without writing:

```sh
weave migrate --dry-run
```

## Breaking changes policy

Any future breaking schema change MUST include:

1. Updated migration guide section in this document.
2. Explicit release-notes migration steps.
3. `doctor` diagnostics for stale provider or schema integrations.
