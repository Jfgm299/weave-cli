# Doctor Reference

`weave doctor` reports current project health and actionable repair steps.

## What it checks

- Config validity and schema compatibility
- Inventory drift between `weave.yaml` and project projections
- Provider integration freshness (stale/missing projections)

## Typical output

- `status: ok` when no issues are found
- `status: issues` plus issue list and suggested repair commands

## Next steps

- Run suggested commands from doctor output (for example `weave provider repair <name>`).
- Re-run `weave doctor` to confirm convergence.
