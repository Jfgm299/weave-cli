---
name: version-batch-planner
description: >
  Plan and execute new Weave versions using v1-style SDD/TDD discipline.
  Trigger: when the user asks to start a new version, define batches, update PRD requirements, or manage OpenSpec tasks/deferred progressively.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- User asks to create a new version roadmap (v3, v4, etc.).
- User asks to define or refine batch backlog.
- User asks to align work with PRD requirements and OpenSpec traceability.
- User asks to continue implementation batch-by-batch with explicit test evidence.

## Critical Patterns

1. **PRD-first, then execution**
   - Extract/define requirements in PRD before implementation.
   - Every batch task must map to requirement IDs.

2. **v1-style planning artifacts**
   - Default artifacts for substantial versions:
     - `proposal.md` (why/scope)
     - `spec.md` (requirements + scenarios)
     - `design.md` (architecture decisions)
     - `tasks.md` (master backlog + per-batch traceability)
     - `test-strategy.md` (unit/integration/e2e evidence model)
     - `deferred.md` (explicit debt ledger)
   - If an artifact does not add value, document why it is intentionally omitted.

3. **Master backlog must be explicit**
   - Maintain a single master backlog table with batch, requirement, status.
   - Each batch must include:
     - requirement traceability matrix
     - test checklist (unit/integration/e2e/evidence)
     - evidence log
     - blockers section

4. **TDD/SDD flow per batch**
   - Before coding a batch:
     - write/update test expectations in `tasks.md`
     - define acceptance evidence for each requirement
   - Then implement incrementally and update status honestly.

5. **Deferred discipline is mandatory**
   - If any scope is incomplete/unconfirmed, DO NOT mark `[x]`.
   - Add entry to version `deferred.md` with:
     - ID
     - requirement
     - deferred from (batch/task)
     - reason
     - closure criteria
     - target version/batch

6. **No fake closure**
   - `[x]` only when implementation + verification evidence exist.
   - Keep `[ - ]`/`[ ]` if evidence is missing.

## Execution Workflow

1. **Version kickoff**
   - Audit prior deferred debt.
   - Add carry-over items to new version backlog (typically Batch 0).
   - Propose initial batch plan with tradeoffs.

2. **Batch setup**
   - Expand `tasks.md` for target batch with requirement-level checklist.
   - Define exact tests and expected evidence before implementation.

3. **Implementation + validation**
   - Implement only current batch scope.
   - Run targeted tests and collect concrete outputs/URLs.
   - Update evidence log.

4. **Batch closure gate**
   - If all requirements verified -> mark `[x]`.
   - Otherwise add deferred entry and keep status in progress.

5. **Move to next batch**
   - Re-check master backlog consistency before starting next batch.

## Tasks.md Structure Contract (v1-style)

Minimum sections per change:

- `Master Backlog`
- `Batch Plan Overview`
- For each batch:
  - `Requirement Traceability Matrix`
  - `Tasks (live checklist)` grouped by requirement
  - `Evidence Log (current)`
  - `Blockers`

## Deferred.md Structure Contract

Use table columns:

| ID | Requirement | Deferred From | Why Deferred | Closure Criteria | Target |

No implicit debt allowed.

## Commands

```bash
# Audit current change artifacts
ls openspec/changes/<change-name>

# Track repo state before batch execution
git status --short --branch

# Run targeted tests for current batch
go test ./...                    # only when full-suite confirmation is required
go test ./internal/... -run ...  # targeted unit/integration
go test -tags=e2e ./test/e2e -run ...
```

## Resources

- **Primary baseline**: `openspec/changes/weave-v1-mvp/tasks.md`
- **Current version pattern**: `openspec/changes/weave-v2-codex-routing/tasks.md`
- **Deferred ledger examples**:
  - `openspec/changes/weave-v1-mvp/deferred.md`
  - `openspec/changes/weave-v2-codex-routing/deferred.md`
