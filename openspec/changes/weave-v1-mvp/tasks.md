# Tasks — weave-v1-mvp

## Scope

Implement v1 CLI behavior from the PRD using TDD + SDD traceability.

## Status Legend

- [ ] Pending
- [-] In progress
- [x] Done
- [!] Blocked

## Priority Legend

- **P0**: Mandatory for v1 acceptance
- **P1**: Important hardening for v1
- **P2**: Nice-to-have follow-up for v1.x

## Execution Protocol (mandatory)

For every requirement (`R-*`):

1. Write at least one failing **success-path** test.
2. Write at least one failing **error/edge-case** test.
3. Implement minimal code to pass both.
4. Capture at least one **observable acceptance evidence** (exit code, filesystem state, config state, or JSON output).
5. Mark the requirement block done only when all checklist items are complete.

---

## Master Backlog (PRD-wide v1 requirements)

| Batch | Priority | Requirement | Title | Status |
|------:|:--------:|-------------|-------|:------:|
| B1 | P0 | R-CORE-01 | `forge` idempotent | [ ] |
| B1 | P0 | R-CORE-02 | `forge` no overwrite without consent | [ ] |
| B1 | P0 | R-CORE-03 | `forge` creates minimal defaults | [ ] |
| B1 | P0 | R-CONFIG-01 | versioned YAML schema | [ ] |
| B1 | P0 | R-CONFIG-02 | validate schema before mutate | [ ] |
| B1 | P0 | R-CONFIG-03 | deterministic regeneration from config | [ ] |
| B1 | P0 | R-CONFIG-04 | v1 sync.mode fixed to symlink | [ ] |
| B1 | P0 | R-CONFIG-05 | explicit inventory in `weave.yaml` | [ ] |
| B2 | P0 | R-CONFIG-06 | atomic update after successful fs ops | [ ] |
| B2 | P0 | R-CONFIG-07 | strict transactional persistence | [ ] |
| B2 | P0 | R-SKILL-04 | v1 skills sync mode = symlink | [ ] |
| B2 | P0 | R-CMD-03 | v1 commands sync mode = symlink | [ ] |
| B2 | P0 | R-SKILL-03 | configurable sources via config/flag/env | [ ] |
| B2 | P0 | R-CMD-01 | commands lifecycle mirrors skills | [ ] |
| B3 | P0 | R-PROV-01 | multi-provider enable in one project | [ ] |
| B3 | P0 | R-PROV-02 | provider adapter architecture | [ ] |
| B3 | P1 | R-PROV-03 | reversible provider operations (remove/repair) | [ ] |
| B3 | P0 | R-DEP-02 | validate provider binaries before success claim | [ ] |
| B3 | P1 | R-DEP-04 | actionable failure messages | [ ] |
| B4 | P0 | R-UX-02 | `doctor` explains status + repair | [ ] |
| B4 | P0 | R-ARCH-06 | inventory persistence transactional state | [ ] |
| B4 | P0 | R-ARCH-07 | symlink + config write as single logical transaction | [ ] |
| B4 | P0 | R-NFR-02 | deterministic exit codes | [ ] |
| B4 | P1 | R-UX-04 | script-friendly output (`--json`) | [ ] |
| B4 | P1 | R-POST-03 | errors include docs references | [ ] |
| B5 | P0 | R-DEP-01 | detect project root + safe scope guard | [ ] |
| B5 | P0 | R-DEP-03 | `--dry-run` for all mutating commands | [ ] |
| B5 | P0 | R-DEP-07 | shell-agnostic core CLI commands | [ ] |
| B5 | P1 | R-UX-01 | concise summary for mutating commands | [ ] |
| B5 | P1 | R-UX-03 | actionability-first UX language | [ ] |
| B5 | P1 | R-UX-05 | strict-mode rollback semantics in output | [ ] |
| B6 | P1 | R-SKILL-01 | conflict prompt (overwrite/skip/backup) | [ ] |
| B6 | P1 | R-SKILL-02 | non-interactive conflict flags | [ ] |
| B6 | P1 | R-CMD-02 | command metadata provider-compat markers (future-ready) | [ ] |
| B6 | P1 | R-ARCH-01 | UI-agnostic core business logic | [ ] |
| B6 | P1 | R-ARCH-02 | no provider leakage into command handlers | [ ] |
| B6 | P1 | R-ARCH-03 | dry-run/real-run share planning primitives | [ ] |
| B6 | P1 | R-ARCH-05 | `.agents` canonical, provider dirs are projections | [ ] |
| B6 | P1 | R-NFR-01 | backup-on-write for risky paths | [ ] |
| B6 | P1 | R-NFR-03 | unit coverage for critical paths | [ ] |
| B6 | P1 | R-POST-01 | help includes 60-second quickstart | [ ] |
| B6 | P1 | R-POST-02 | first-run contextual next-step guidance | [ ] |
| B6 | P1 | R-DIST-01 | signed checksums for release artifacts | [ ] |
| B6 | P1 | R-DIST-02 | install docs include quickstart + troubleshooting | [ ] |
| B6 | P1 | R-DIST-03 | semver binary naming/versioning | [ ] |
| B6 | P1 | R-UPD-01 | detect outdated schema + suggest/perform migration | [ ] |
| B6 | P1 | R-UPD-02 | breaking changes include migration guide | [ ] |
| B6 | P1 | R-UPD-03 | `doctor` flags stale provider integrations after upgrades | [ ] |
| B7 | P2 | R-ARCH-04 | TUI integration possible without domain refactor | [ ] |
| B7 | P2 | R-DEP-05 | no TUI runtime dependency in v1 binary | [ ] |
| B7 | P2 | R-DEP-06 | one-command install for Weave binary | [ ] |

---

## Batch Plan Overview

- **Batch 1 (B1) — Foundation + Forge + Config Baseline (P0)**
- **Batch 2 (B2) — Transactionality + Symlink-only lifecycle (P0)**
- **Batch 3 (B3) — Provider integration baseline (P0/P1)**
- **Batch 4 (B4) — Doctor + transactional observability + exit semantics (P0/P1)**
- **Batch 5 (B5) — Safety guards + dry-run + shell-agnostic behavior (P0/P1)**
- **Batch 6 (B6) — UX hardening + architecture guardrails + release/update requirements (P1)**
- **Batch 7 (B7) — Deferred v1.x hardening (P2)**

---

## Batch 1 — Foundation + Forge + Config Baseline (current)

### B1 Exit Criteria

- [ ] `forge` is idempotent and safe on existing projects.
- [ ] Minimal project defaults are generated once and preserved.
- [ ] `weave.yaml` is schema-validated and deterministic.
- [ ] v1 config enforces symlink mode and explicit inventory model.

### B1 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-CORE-01 (`forge` idempotent) | B1-T1.1, B1-T1.2 | B1-T1.3 | B1-T1.4 | B1-T1.5 | [ ] |
| R-CORE-02 (no overwrite without consent) | B1-T2.1, B1-T2.2 | B1-T2.3 | B1-T2.4 | B1-T2.5 | [ ] |
| R-CORE-03 (minimal defaults) | B1-T3.1, B1-T3.2 | B1-T3.3 | B1-T3.4 | B1-T3.5 | [ ] |
| R-CONFIG-01 (schema versioned) | B1-T4.1, B1-T4.2 | B1-T4.3 | B1-T4.4 | B1-T4.5 | [ ] |
| R-CONFIG-02 (validate before mutate) | B1-T5.1, B1-T5.2 | B1-T5.3 | B1-T5.4 | B1-T5.5 | [ ] |
| R-CONFIG-03 (deterministic regeneration) | B1-T6.1, B1-T6.2 | B1-T6.3 | B1-T6.4 | B1-T6.5 | [ ] |
| R-CONFIG-04 (symlink mode fixed) | B1-T7.1, B1-T7.2 | B1-T7.3 | B1-T7.4 | B1-T7.5 | [ ] |
| R-CONFIG-05 (explicit inventory) | B1-T8.1, B1-T8.2 | B1-T8.3 | B1-T8.4 | B1-T8.5 | [ ] |

### B1 Tasks (live checklist)

#### R-CORE-01 — `forge` MUST be idempotent

- [ ] **B1-T1.1 Unit (success):** planner returns no-op when project state is already converged.
- [ ] **B1-T1.2 Unit (edge):** partial existing structure does not trigger destructive plan.
- [ ] **B1-T1.3 Integration:** repeated apply yields the same filesystem/config state.
- [ ] **B1-T1.4 E2E:** run `weave forge` twice in the same repo.
- [ ] **B1-T1.5 Evidence:** exit code `0` + no destructive mutations on second run.

#### R-CORE-02 — `forge` MUST avoid overwriting existing files without explicit consent

- [ ] **B1-T2.1 Unit (success):** planner marks existing protected files as keep/no-op by default.
- [ ] **B1-T2.2 Unit (error/edge):** overwrite attempt without consent returns consent-required error.
- [ ] **B1-T2.3 Integration:** apply respects existing files and preserves content unless explicit overwrite option is present.
- [ ] **B1-T2.4 E2E:** pre-existing `weave.yaml`/`.agents/AGENTS.md` remain unchanged on default `forge`.
- [ ] **B1-T2.5 Evidence:** checksum/content assertions prove no overwrite occurred.

#### R-CORE-03 — `forge` MUST create minimal defaults

- [ ] **B1-T3.1 Unit (success):** default scaffold plan includes required baseline paths/files.
- [ ] **B1-T3.2 Unit (edge):** missing optional dirs do not block minimal scaffold creation.
- [ ] **B1-T3.3 Integration:** fresh project receives expected baseline structure.
- [ ] **B1-T3.4 E2E:** `weave forge` in empty repo creates minimal defaults.
- [ ] **B1-T3.5 Evidence:** assert expected tree and key files exist.

#### R-CONFIG-01 — YAML schema MUST be versioned

- [ ] **B1-T4.1 Unit (success):** schema validator accepts supported `version` values.
- [ ] **B1-T4.2 Unit (error/edge):** missing/invalid `version` fails validation.
- [ ] **B1-T4.3 Integration:** config load pipeline rejects invalid schema before planning.
- [ ] **B1-T4.4 E2E:** invalid `weave.yaml` version returns non-zero and actionable message.
- [ ] **B1-T4.5 Evidence:** deterministic validation error includes version guidance.

#### R-CONFIG-02 — CLI MUST validate schema before mutating filesystem

- [ ] **B1-T5.1 Unit (success):** valid config enables plan/apply path.
- [ ] **B1-T5.2 Unit (error):** invalid config short-circuits before fs mutation calls.
- [ ] **B1-T5.3 Integration:** mutation executor is never invoked when validation fails.
- [ ] **B1-T5.4 E2E:** run mutating command with invalid config and verify no fs changes.
- [ ] **B1-T5.5 Evidence:** pre/post filesystem snapshot unchanged.

#### R-CONFIG-03 — CLI MUST support deterministic regeneration from config

- [ ] **B1-T6.1 Unit (success):** same config input produces stable operation plan ordering.
- [ ] **B1-T6.2 Unit (edge):** unordered YAML keys normalize to deterministic internal model.
- [ ] **B1-T6.3 Integration:** regenerate twice from same config yields identical outputs.
- [ ] **B1-T6.4 E2E:** same command + same config => same resulting state and output signature.
- [ ] **B1-T6.5 Evidence:** deterministic hash/snapshot for plan or resulting config.

#### R-CONFIG-04 — `sync.mode` MUST be fixed to `symlink` in v1

- [ ] **B1-T7.1 Unit (success):** config parser accepts `symlink` mode.
- [ ] **B1-T7.2 Unit (error):** parser rejects any non-`symlink` mode in v1.
- [ ] **B1-T7.3 Integration:** planner refuses non-symlink mode before operation planning.
- [ ] **B1-T7.4 E2E:** config with `sync.mode: copy` fails with explicit v1 constraint message.
- [ ] **B1-T7.5 Evidence:** no mutation + clear policy error.

#### R-CONFIG-05 — `weave.yaml` MUST store explicit desired inventory

- [ ] **B1-T8.1 Unit (success):** config model persists explicit skills/commands inventory entries.
- [ ] **B1-T8.2 Unit (edge):** duplicate inventory entries are normalized/rejected consistently.
- [ ] **B1-T8.3 Integration:** write/read roundtrip preserves explicit inventory exactly.
- [ ] **B1-T8.4 E2E:** after baseline ops, `weave.yaml` shows explicit desired inventory structure.
- [ ] **B1-T8.5 Evidence:** schema-compliant config snapshot with explicit inventory fields.

---

## Progress Log

- [ ] Batch 1 started
- [ ] Batch 1 completed
- [ ] Batch 2 started
- [ ] Batch 2 completed
- [ ] Batch 3 started
- [ ] Batch 3 completed
- [ ] Batch 4 started
- [ ] Batch 4 completed
- [ ] Batch 5 started
- [ ] Batch 5 completed
- [ ] Batch 6 started
- [ ] Batch 6 completed
- [ ] Batch 7 started
- [ ] Batch 7 completed

---

## Implementation Notes

- Keep this file updated during implementation (live checklist behavior).
- Only mark a requirement block complete after all test levels + evidence are complete.
- Add regression tasks under the relevant requirement when bugs are found.
