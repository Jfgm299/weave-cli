# Doctor Reference

`weave doctor` reports current project health and actionable repair steps.

## What it checks

- Config validity and schema compatibility
- Inventory drift between `weave.yaml` and project projections
- Provider integration freshness (stale/missing projections)
- Command projection integrity for both install modes:
  - shared mode (`weave command add <name>`) verifies `.agents/commands/<name>.md` plus provider projections
  - exclusive mode (`weave command add <name> --provider <name>`) verifies provider projection paths without requiring shared `.agents` command link

Doctor uses command metadata (`provider_compat`, `shared_install`) from `weave.yaml` to derive expected projection paths.

## Typical output

- `status: ok` when no issues are found
- `status: issues` plus issue list and suggested repair commands

## Next steps

- Run suggested commands from doctor output (for example `weave provider repair <name>`).
- For command projection drift, doctor suggests deterministic command repair paths:
  - shared install drift -> `weave command add <name>`
  - exclusive provider drift -> `weave command add <name> --provider <name>`
- Re-run `weave doctor` to confirm convergence.
