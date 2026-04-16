# Transaction Semantics

Mutating operations follow strict v1 transaction semantics.

## Rules

1. Validate config and prerequisites first.
2. Apply filesystem operations.
3. Persist `weave.yaml` only after successful fs operations.

## Strict rollback behavior

If config persistence fails after symlink/fs apply:

- rollback is attempted for already-applied risky operations,
- errors must clearly state whether rollback succeeded or partial state may remain.

If filesystem apply fails during a multi-provider command transaction:

- rollback is attempted for every planned shared + provider projection path,
- command fails non-zero with explicit rollback outcome,
- `weave.yaml` remains unchanged because config persistence occurs only after successful apply.

## Recovery

When rollback fails or state is uncertain:

```sh
weave doctor
```

Then run the suggested repair command.
