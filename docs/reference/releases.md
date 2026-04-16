# Release Notes Policy

## Migration Guide requirement for breaking changes

Any release that introduces breaking behavior MUST include:

- A `Migration` section with before/after examples.
- Explicit commands to migrate (`weave migrate` when applicable).
- Rollback guidance.

## Minimum release checklist

- [ ] Semver tag
- [ ] Versioned artifact names
- [ ] checksums.txt + checksums.txt.sig
- [ ] Migration section when breaking changes exist
