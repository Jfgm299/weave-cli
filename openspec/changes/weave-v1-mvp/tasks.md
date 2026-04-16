# Tasks — weave-v1-mvp

## Scope

Implement v1 CLI behavior from the PRD using TDD + SDD traceability.

## Status Legend

- [ ] Pending
- [-] In progress
- [x] Done (100% implemented, tested, and operational)
- [!] Blocked

## Completion Semantics Policy (mandatory)

To avoid ambiguity, this project enforces:

- **Never mark `[x]` unless the requirement is 100% implemented in code and validated by required evidence/tests.**
- If scope is intentionally cut for v1, use `[-]` and register the item in `deferred.md` with owner + closure criteria.
- Documentation/policy-only closure MUST NOT be marked `[x]` unless the requirement itself is explicitly docs-only.
- `E2E: N/A` is only allowed when the requirement is moved to deferred tracking.

### Completion type taxonomy

- `done-code`: fully implemented in code, tested, and operationally closed.
- `done-baseline`: policy/docs baseline accepted for v1, with deferred executable work tracked.
- `partial`: implemented in part but not yet fully closed.
- `deferred`: intentionally postponed to later batch/version.

In this file:
- `[x]` MUST only represent `done-code`.
- `done-baseline`, `partial`, and `deferred` MUST use `[-]` plus a linked deferred item when applicable.

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
| B1 | P0 | R-CORE-01 | `forge` idempotent | [x] |
| B1 | P0 | R-CORE-02 | `forge` no overwrite without consent | [x] |
| B1 | P0 | R-CORE-03 | `forge` creates minimal defaults | [x] |
| B1 | P0 | R-CONFIG-01 | versioned YAML schema | [x] |
| B1 | P0 | R-CONFIG-02 | validate schema before mutate | [x] |
| B1 | P0 | R-CONFIG-03 | deterministic regeneration from config | [x] |
| B1 | P0 | R-CONFIG-04 | v1 sync.mode fixed to symlink | [x] |
| B1 | P0 | R-CONFIG-05 | explicit inventory in `weave.yaml` | [x] |
| B2 | P0 | R-CONFIG-06 | atomic update after successful fs ops | [x] |
| B2 | P0 | R-CONFIG-07 | strict transactional persistence | [x] |
| B2 | P0 | R-SKILL-04 | v1 skills sync mode = symlink | [x] |
| B2 | P0 | R-CMD-03 | v1 commands sync mode = symlink | [x] |
| B2 | P1 | R-SKILL-05 | default shared skills directory (`~/.weave/skills`) | [x] |
| B2 | P1 | R-CMD-04 | default shared commands directory (`~/.weave/commands`) | [x] |
| B2 | P0 | R-SKILL-03 | configurable sources via config/flag/env | [x] |
| B2 | P0 | R-CMD-01 | commands lifecycle mirrors skills | [x] |
| B3 | P0 | R-PROV-01 | multi-provider enable in one project | [x] |
| B3 | P0 | R-PROV-02 | provider adapter architecture | [x] |
| B3 | P1 | R-PROV-03 | reversible provider operations (remove/repair) | [x] |
| B3 | P0 | R-DEP-02 | validate provider binaries before success claim | [x] |
| B3 | P1 | R-DEP-04 | actionable failure messages | [x] |
| B4 | P0 | R-UX-02 | `doctor` explains status + repair | [x] |
| B4 | P0 | R-ARCH-06 | inventory persistence transactional state | [x] |
| B4 | P0 | R-ARCH-07 | symlink + config write as single logical transaction | [x] |
| B4 | P0 | R-NFR-02 | deterministic exit codes | [x] |
| B4 | P1 | R-UX-04 | script-friendly output (`--json`) | [x] |
| B4 | P1 | R-POST-03 | errors include docs references | [x] |
| B5 | P0 | R-DEP-01 | detect project root + safe scope guard | [x] |
| B5 | P0 | R-DEP-03 | `--dry-run` for all mutating commands | [x] |
| B5 | P0 | R-DEP-07 | shell-agnostic core CLI commands | [x] |
| B5 | P1 | R-UX-01 | concise summary for mutating commands | [x] |
| B5 | P1 | R-UX-03 | actionability-first UX language | [x] |
| B5 | P1 | R-UX-05 | strict-mode rollback semantics in output | [x] |
| B6 | P1 | R-SKILL-01 | conflict prompt (overwrite/skip/backup) | [x] |
| B6 | P1 | R-SKILL-02 | non-interactive conflict flags | [x] |
| B6 | P1 | R-CMD-02 | command metadata provider-compat markers (future-ready) | [x] |
| B6 | P1 | R-ARCH-01 | UI-agnostic core business logic | [x] |
| B6 | P1 | R-ARCH-02 | no provider leakage into command handlers | [x] |
| B6 | P1 | R-ARCH-03 | dry-run/real-run share planning primitives | [x] |
| B6 | P1 | R-ARCH-05 | `.agents` canonical, provider dirs are projections | [x] |
| B6 | P1 | R-NFR-01 | backup-on-write for risky paths | [x] |
| B6 | P1 | R-NFR-03 | unit coverage for critical paths | [x] |
| B6 | P1 | R-POST-01 | help includes 60-second quickstart | [x] |
| B6 | P1 | R-POST-02 | first-run contextual next-step guidance | [x] |
| B6 | P1 | R-DIST-01 | signed checksums for release artifacts | [-] |
| B6 | P1 | R-DIST-02 | install docs include quickstart + troubleshooting | [-] |
| B6 | P1 | R-DIST-03 | semver binary naming/versioning | [x] |
| B6 | P1 | R-UPD-01 | detect outdated schema + suggest/perform migration | [x] |
| B6 | P1 | R-UPD-02 | breaking changes include migration guide | [-] |
| B6 | P1 | R-UPD-03 | `doctor` flags stale provider integrations after upgrades | [x] |
| B7 | P2 | R-ARCH-04 | TUI integration possible without domain refactor | [x] |
| B7 | P2 | R-DEP-05 | no TUI runtime dependency in v1 binary | [x] |
| B7 | P2 | R-DEP-06 | one-command install for Weave binary | [x] |

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

- [x] `forge` is idempotent and safe on existing projects.
- [x] Minimal project defaults are generated once and preserved.
- [x] `weave.yaml` is schema-validated and deterministic.
- [x] v1 config enforces symlink mode and explicit inventory model.

### B1 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-CORE-01 (`forge` idempotent) | B1-T1.1, B1-T1.2 | B1-T1.3 | B1-T1.4 | B1-T1.5 | [x] |
| R-CORE-02 (no overwrite without consent) | B1-T2.1, B1-T2.2 | B1-T2.3 | B1-T2.4 | B1-T2.5 | [x] |
| R-CORE-03 (minimal defaults) | B1-T3.1, B1-T3.2 | B1-T3.3 | B1-T3.4 | B1-T3.5 | [x] |
| R-CONFIG-01 (schema versioned) | B1-T4.1, B1-T4.2 | B1-T4.3 | B1-T4.4 | B1-T4.5 | [x] |
| R-CONFIG-02 (validate before mutate) | B1-T5.1, B1-T5.2 | B1-T5.3 | B1-T5.4 | B1-T5.5 | [x] |
| R-CONFIG-03 (deterministic regeneration) | B1-T6.1, B1-T6.2 | B1-T6.3 | B1-T6.4 | B1-T6.5 | [x] |
| R-CONFIG-04 (symlink mode fixed) | B1-T7.1, B1-T7.2 | B1-T7.3 | B1-T7.4 | B1-T7.5 | [x] |
| R-CONFIG-05 (explicit inventory) | B1-T8.1, B1-T8.2 | B1-T8.3 | B1-T8.4 | B1-T8.5 | [x] |

### B1 Tasks (live checklist)

#### R-CORE-01 — `forge` MUST be idempotent

- [x] **B1-T1.1 Unit (success):** planner returns no-op when project state is already converged.
- [x] **B1-T1.2 Unit (edge):** partial existing structure does not trigger destructive plan.
- [x] **B1-T1.3 Integration:** repeated apply yields the same filesystem/config state.
- [x] **B1-T1.4 E2E:** run `weave forge` twice in the same repo.
- [x] **B1-T1.5 Evidence:** exit code `0` + no destructive mutations on second run.

#### R-CORE-02 — `forge` MUST avoid overwriting existing files without explicit consent

- [x] **B1-T2.1 Unit (success):** planner marks existing protected files as keep/no-op by default.
- [x] **B1-T2.2 Unit (error/edge):** overwrite attempt without consent returns consent-required error.
- [x] **B1-T2.3 Integration:** apply respects existing files and preserves content unless explicit overwrite option is present.
- [x] **B1-T2.4 E2E:** pre-existing `weave.yaml`/`.agents/AGENTS.md` remain unchanged on default `forge`.
- [x] **B1-T2.5 Evidence:** checksum/content assertions prove no overwrite occurred.

#### R-CORE-03 — `forge` MUST create minimal defaults

- [x] **B1-T3.1 Unit (success):** default scaffold plan includes required baseline paths/files.
- [x] **B1-T3.2 Unit (edge):** missing optional dirs do not block minimal scaffold creation.
- [x] **B1-T3.3 Integration:** fresh project receives expected baseline structure.
- [x] **B1-T3.4 E2E:** `weave forge` in empty repo creates minimal defaults.
- [x] **B1-T3.5 Evidence:** assert expected tree and key files exist.

#### R-CONFIG-01 — YAML schema MUST be versioned

- [x] **B1-T4.1 Unit (success):** schema validator accepts supported `version` values.
- [x] **B1-T4.2 Unit (error/edge):** missing/invalid `version` fails validation.
- [x] **B1-T4.3 Integration:** config load pipeline rejects invalid schema before planning.
- [x] **B1-T4.4 E2E:** invalid `weave.yaml` version returns non-zero and actionable message.
- [x] **B1-T4.5 Evidence:** deterministic validation error includes version guidance.

#### R-CONFIG-02 — CLI MUST validate schema before mutating filesystem

- [x] **B1-T5.1 Unit (success):** valid config enables plan/apply path.
- [x] **B1-T5.2 Unit (error):** invalid config short-circuits before fs mutation calls.
- [x] **B1-T5.3 Integration:** mutation executor is never invoked when validation fails.
- [x] **B1-T5.4 E2E:** run mutating command with invalid config and verify no fs changes.
- [x] **B1-T5.5 Evidence:** pre/post filesystem snapshot unchanged.

#### R-CONFIG-03 — CLI MUST support deterministic regeneration from config

- [x] **B1-T6.1 Unit (success):** same config input produces stable operation plan ordering.
- [x] **B1-T6.2 Unit (edge):** unordered YAML keys normalize to deterministic internal model.
- [x] **B1-T6.3 Integration:** regenerate twice from same config yields identical outputs.
- [x] **B1-T6.4 E2E:** same command + same config => same resulting state and output signature.
- [x] **B1-T6.5 Evidence:** deterministic hash/snapshot for plan or resulting config.

#### R-CONFIG-04 — `sync.mode` MUST be fixed to `symlink` in v1

- [x] **B1-T7.1 Unit (success):** config parser accepts `symlink` mode.
- [x] **B1-T7.2 Unit (error):** parser rejects any non-`symlink` mode in v1.
- [x] **B1-T7.3 Integration:** planner refuses non-symlink mode before operation planning.
- [x] **B1-T7.4 E2E:** config with `sync.mode: copy` fails with explicit v1 constraint message.
- [x] **B1-T7.5 Evidence:** no mutation + clear policy error.

#### R-CONFIG-05 — `weave.yaml` MUST store explicit desired inventory

- [x] **B1-T8.1 Unit (success):** config model persists explicit skills/commands inventory entries.
- [x] **B1-T8.2 Unit (edge):** duplicate inventory entries are normalized/rejected consistently.
- [x] **B1-T8.3 Integration:** write/read roundtrip preserves explicit inventory exactly.
- [x] **B1-T8.4 E2E:** after baseline ops, `weave.yaml` shows explicit desired inventory structure.
- [x] **B1-T8.5 Evidence:** schema-compliant config snapshot with explicit inventory fields.

---

## Batch 2 — Transactionality + Symlink-only lifecycle

### B2 Exit Criteria

- [x] `skill add` and `command add` create symlink operations only in v1.
- [x] source resolution precedence is deterministic: flag > env > config > defaults.
- [x] `weave.yaml` persists only after successful fs ops.
- [x] strict mode guarantees no config mutation when symlink creation fails.
- [x] commands lifecycle mirrors skills lifecycle for add behavior.

### B2 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-CONFIG-06 (atomic update after successful fs ops) | B2-T1.1, B2-T1.2 | B2-T1.3 | B2-T1.4 | B2-T1.5 | [x] |
| R-CONFIG-07 (strict transactional persistence) | B2-T2.1, B2-T2.2 | B2-T2.3 | B2-T2.4 | B2-T2.5 | [x] |
| R-SKILL-04 (skills sync mode = symlink) | B2-T3.1, B2-T3.2 | B2-T3.3 | B2-T3.4 | B2-T3.5 | [x] |
| R-CMD-03 (commands sync mode = symlink) | B2-T4.1, B2-T4.2 | B2-T4.3 | B2-T4.4 | B2-T4.5 | [x] |
| R-SKILL-05 (default shared skills directory path) | B2-T7.1, B2-T7.2 | B2-T7.3 | B2-T7.4 | B2-T7.5 | [x] |
| R-CMD-04 (default shared commands directory path) | B2-T8.1, B2-T8.2 | B2-T8.3 | B2-T8.4 | B2-T8.5 | [x] |
| R-SKILL-03 (configurable skill sources + override precedence) | B2-T5.1, B2-T5.2 | B2-T5.3 | B2-T5.4 | B2-T5.5 | [x] |
| R-CMD-01 (commands lifecycle mirrors skills lifecycle) | B2-T6.1, B2-T6.2 | B2-T6.3 | B2-T6.4 | B2-T6.5 | [x] |

### B2 Tasks (live checklist)

#### R-CONFIG-06 — `skill add` / `command add` MUST update config atomically only after fs success

- [x] **B2-T1.1 Unit (success):** atomic writer updates `weave.yaml` via temp file + rename semantics.
- [x] **B2-T1.2 Unit (edge):** writer failure preserves previous `weave.yaml` content.
- [x] **B2-T1.3 Integration:** add lifecycle executes symlink op before config persistence path.
- [x] **B2-T1.4 E2E:** successful `skill add` and `command add` update links and config in one logical flow.
- [x] **B2-T1.5 Evidence:** post-run config contains new inventory entry only after symlink exists.

Implementation evidence:

- [x] `internal/config/atomic_writer_test.go::TestAtomicFileWriter_Write_UpdatesExistingConfigAtomically`
- [x] `internal/config/atomic_writer_test.go::TestAtomicFileWriter_Write_RenameFailureLeavesOriginalUnchanged`
- [x] `test/integration/asset_add_integration_test.go::TestAddAsset_Integration_SuccessCreatesSymlinkAndConfigEntry`

#### R-CONFIG-07 — strict mode MUST avoid config mutation on symlink failure

- [x] **B2-T2.1 Unit (success):** service persists config when executor succeeds.
- [x] **B2-T2.2 Unit (error):** service never calls config writer when symlink op fails.
- [x] **B2-T2.3 Integration:** failing symlink destination returns error and leaves `weave.yaml` unchanged.
- [x] **B2-T2.4 E2E:** `weave skill add` failure keeps config byte-identical.
- [x] **B2-T2.5 Evidence:** non-zero exit + unchanged config snapshot + actionable error.

Implementation evidence:

- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_StrictModeDoesNotPersistOnSymlinkFailure`
- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_StrictModeDoesNotPersistOnSymlinkFailureForCommand`
- [x] `test/e2e/forge_e2e_test.go::TestSkillAdd_E2E_StrictFailureKeepsConfigUnchanged`
- [x] `test/e2e/forge_e2e_test.go::TestCommandAdd_E2E_StrictFailureKeepsConfigUnchanged`

#### R-SKILL-04 — skills sync mode in v1 MUST be symlink-only

- [x] **B2-T3.1 Unit (success):** skill add planner emits `create_link` op for destination in `.agents/skills`.
- [x] **B2-T3.2 Unit (edge):** no copy fallback operation is emitted in v1 code path.
- [x] **B2-T3.3 Integration:** service + executor produce an actual symlink for added skill.
- [x] **B2-T3.4 E2E:** `weave skill add <name>` creates symlink, not copied payload.
- [x] **B2-T3.5 Evidence:** `lstat` reports `ModeSymlink` for installed skill path.

Implementation evidence:

- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_UsesCreateLinkForSkillAndCommand/skill`
- [x] `internal/cli/assets_handler_test.go::TestAssetAddService_Add_CreatesSymlinkAndPersistsConfig`
- [x] `test/e2e/forge_e2e_test.go::TestAssetAdd_E2E_UsesSourcePrecedenceAndSymlinkOnly`

#### R-CMD-03 — commands sync mode in v1 MUST be symlink-only

- [x] **B2-T4.1 Unit (success):** command add planner emits `create_link` op for `.agents/commands`.
- [x] **B2-T4.2 Unit (edge):** command add path has no copy fallback in v1.
- [x] **B2-T4.3 Integration:** command add creates symlink in project commands directory.
- [x] **B2-T4.4 E2E:** `weave command add <name>` results in symlink install behavior.
- [x] **B2-T4.5 Evidence:** `lstat` confirms command installation path is a symlink.

Implementation evidence:

- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_UsesCreateLinkForSkillAndCommand/command`
- [x] `test/integration/asset_add_integration_test.go::TestAddAsset_Integration_CommandSuccessCreatesSymlinkAndConfigEntry`
- [x] `test/e2e/forge_e2e_test.go::TestAssetAdd_E2E_UsesSourcePrecedenceAndSymlinkOnly`

#### R-SKILL-03 — skill source directory MUST be configurable with precedence

- [x] **B2-T5.1 Unit (success):** source resolver uses `--from` over env/config/default.
- [x] **B2-T5.2 Unit (edge):** source resolver falls back env > config > default when higher precedence is absent.
- [x] **B2-T5.3 Integration:** configured source directory is used when no flag/env override is provided.
- [x] **B2-T5.4 E2E:** verify precedence by setting config + env + flag and asserting selected source.
- [x] **B2-T5.5 Evidence:** persisted asset source path matches precedence-selected directory.

Implementation evidence:

- [x] `internal/cli/source_resolver_test.go::TestSourceResolver_Resolve_FlagOverridesEnvConfigAndDefault`
- [x] `internal/cli/source_resolver_test.go::TestSourceResolver_Resolve_EnvOverridesConfigAndDefault`
- [x] `internal/cli/source_resolver_test.go::TestSourceResolver_Resolve_CommandEnvOverridesConfigAndDefault`
- [x] `test/e2e/forge_e2e_test.go::TestAssetAdd_E2E_UsesSourcePrecedenceAndSymlinkOnly`

#### R-CMD-01 — command add lifecycle MUST mirror skill add lifecycle

- [x] **B2-T6.1 Unit (success):** command add shares the same add service lifecycle stages as skill add.
- [x] **B2-T6.2 Unit (edge):** command add preserves strict persistence behavior on symlink failure.
- [x] **B2-T6.3 Integration:** command add follows detect/validate/symlink/persist flow equivalent to skills.
- [x] **B2-T6.4 E2E:** command add behavior aligns with skill add for success and failure semantics.
- [x] **B2-T6.5 Evidence:** command inventory update + symlink creation mirror skill lifecycle guarantees.

Implementation evidence:

- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_PersistsConfigAfterSymlinkSuccess`
- [x] `test/integration/asset_add_integration_test.go::TestAddAsset_Integration_CommandSuccessCreatesSymlinkAndConfigEntry`
- [x] `test/e2e/forge_e2e_test.go::TestCommandAdd_E2E_StrictFailureKeepsConfigUnchanged`

#### R-SKILL-05 — default shared skills directory path MUST be `~/.weave/skills`

- [x] **B2-T7.1 Unit (success):** default config sets `sources.skills_dir` to `~/.weave/skills`.
- [x] **B2-T7.2 Unit (edge):** resolver falls back to `~/.weave/skills` when flag/env/config are absent.
- [x] **B2-T7.3 Integration:** default project config roundtrip preserves `~/.weave/skills`.
- [x] **B2-T7.4 E2E:** default source resolution selects the `~/.weave/skills` baseline.
- [x] **B2-T7.5 Evidence:** `internal/config/defaults.go`, `internal/cli/source_resolver.go`, and tests assert the default path.

#### R-CMD-04 — default shared commands directory path MUST be `~/.weave/commands`

- [x] **B2-T8.1 Unit (success):** default config sets `sources.commands_dir` to `~/.weave/commands`.
- [x] **B2-T8.2 Unit (edge):** resolver falls back to `~/.weave/commands` when flag/env/config are absent.
- [x] **B2-T8.3 Integration:** default project config roundtrip preserves `~/.weave/commands`.
- [x] **B2-T8.4 E2E:** default source resolution selects the `~/.weave/commands` baseline.
- [x] **B2-T8.5 Evidence:** `internal/config/defaults.go`, `internal/cli/source_resolver.go`, and tests assert the default path.

---

## Batch 3 — Provider integration baseline

### B3 Exit Criteria

- [x] `provider add` supports enabling multiple providers in one project state.
- [x] provider behavior is isolated behind adapter interfaces and registry boundaries.
- [x] provider binary prerequisites are validated before setup success is reported.
- [x] baseline reversible provider operations (`remove` + `repair`) are available.
- [x] provider failures include actionable next steps.

### B3 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-PROV-01 (multi-provider enable in one project) | B3-T1.1, B3-T1.2 | B3-T1.3 | B3-T1.4 | B3-T1.5 | [x] |
| R-PROV-02 (provider adapter architecture) | B3-T2.1, B3-T2.2 | B3-T2.3 | B3-T2.4 | B3-T2.5 | [x] |
| R-PROV-03 (reversible provider operations baseline remove/repair) | B3-T3.1, B3-T3.2 | B3-T3.3 | B3-T3.4 | B3-T3.5 | [x] |
| R-DEP-02 (validate required provider binaries before success claim) | B3-T4.1, B3-T4.2 | B3-T4.3 | B3-T4.4 | B3-T4.5 | [x] |
| R-DEP-04 (actionable failure messages) | B3-T5.1, B3-T5.2 | B3-T5.3 | B3-T5.4 | B3-T5.5 | [x] |

### B3 Tasks (live checklist)

#### R-PROV-01 — CLI MUST allow enabling multiple providers in one project

- [x] **B3-T1.1 Unit (success):** provider service appends enabled provider without removing existing enabled providers.
- [x] **B3-T1.2 Unit (edge):** provider list returns only enabled providers in deterministic order.
- [x] **B3-T1.3 Integration:** adding `claude-code` then `opencode` persists both providers in `weave.yaml`.
- [x] **B3-T1.4 E2E:** `weave provider add claude-code` + `weave provider add opencode` + `weave provider list` in same repo.
- [x] **B3-T1.5 Evidence:** `weave.yaml` contains both provider entries and list output includes both names.

Implementation evidence:

- [x] `internal/app/provider_service_test.go::TestProviderService_AddProvider_EnablesMultipleProviders`
- [x] `internal/cli/provider_handler_test.go::TestProviderRunAction_List_ReturnsSortedEnabledProviders`
- [x] `test/integration/provider_integration_test.go::TestProviderAdd_Integration_EnablesMultipleProvidersAndCreatesProjectionLinks`
- [x] `test/e2e/provider_e2e_test.go::TestProviderAddAndList_E2E_MultiProviderFlowAndActionableErrors`

#### R-PROV-02 — Provider integration MUST be adapter-based (interface boundaries)

- [x] **B3-T2.1 Unit (success):** registry resolves known providers via adapter interface lookup.
- [x] **B3-T2.2 Unit (edge):** unsupported provider returns deterministic supported-provider set.
- [x] **B3-T2.3 Integration:** provider add flow invokes adapter-planned operations only (no provider leakage in command parser).
- [x] **B3-T2.4 E2E:** provider add for both adapters creates expected provider projection links.
- [x] **B3-T2.5 Evidence:** `internal/providers` registry + adapter implementations are the only provider-specific operation planners.

Implementation evidence:

- [x] `internal/providers/registry.go`
- [x] `internal/providers/registry_test.go::TestDefaultRegistry_Get_ReturnsKnownProviders`
- [x] `internal/providers/registry_test.go::TestClaudeAdapter_PlanSetup_CreatesProviderProjectionLinks`
- [x] `internal/cli/root.go::runProvider`

#### R-PROV-03 — Provider operations MUST be reversible (remove/repair) [baseline]

- [x] **B3-T3.1 Unit (success):** remove operation drops provider inventory entry and applies adapter remove plan.
- [x] **B3-T3.2 Unit (edge):** repair operation re-applies adapter setup plan and ensures provider stays enabled.
- [x] **B3-T3.3 Integration:** remove deletes projection directory and updates config, repair restores projection links.
- [x] **B3-T3.4 E2E:** `provider remove` then `provider repair` re-establishes provider links.
- [x] **B3-T3.5 Evidence:** `.claude` projection disappears after remove and reappears after repair; config remains consistent.

Implementation evidence:

- [x] `internal/app/provider_service.go::RemoveProvider`
- [x] `internal/app/provider_service.go::RepairProvider`
- [x] `internal/cli/provider_handler_test.go::TestProviderRunAction_Remove_DeletesProjectionAndConfigEntry`
- [x] `test/e2e/provider_e2e_test.go::TestProviderRemoveAndRepair_E2E_ReversibleOperations`

#### R-DEP-02 — CLI MUST validate required provider binaries before claiming provider setup success

- [x] **B3-T4.1 Unit (success):** provider add proceeds when all adapter-required binaries resolve in PATH.
- [x] **B3-T4.2 Unit (error):** provider add fails before fs/config mutation when any required binary is missing.
- [x] **B3-T4.3 Integration:** missing binary keeps `weave.yaml` byte-identical and skips projection writes.
- [x] **B3-T4.4 E2E:** provider add fails with missing binary in constrained PATH.
- [x] **B3-T4.5 Evidence:** non-zero exit + actionable dependency message + unchanged config snapshot.

Implementation evidence:

- [x] `internal/app/provider_service_test.go::TestProviderService_AddProvider_MissingBinaryReturnsActionableError`
- [x] `test/integration/provider_integration_test.go::TestProviderAdd_Integration_MissingBinaryKeepsConfigUnchanged`
- [x] `test/e2e/provider_e2e_test.go::TestProviderAdd_E2E_MissingBinaryActionableFailureAndNoConfigMutation`

#### R-DEP-04 — CLI MUST provide readable failure messages with actionable next steps

- [x] **B3-T5.1 Unit (error):** unsupported provider errors include supported provider list.
- [x] **B3-T5.2 Unit (error):** missing binary errors include installation + repair instruction.
- [x] **B3-T5.3 Integration:** dependency/setup errors flow through service boundary with preserved actionable context.
- [x] **B3-T5.4 E2E:** provider add failure output includes clear remediation command.
- [x] **B3-T5.5 Evidence:** stderr contains explicit “Install the missing binaries…” and `weave provider repair <provider>` guidance.

Implementation evidence:

- [x] `internal/app/provider_service.go` error wrapping messages
- [x] `internal/app/provider_service_test.go::TestProviderService_AddProvider_UnsupportedProviderReturnsSupportedList`
- [x] `internal/cli/provider_handler_test.go::TestProviderRunAction_Add_MissingBinaryReturnsActionableMessage`

---

## Batch 4 — Doctor + transactional observability + exit semantics

### B4 Exit Criteria

- [x] `doctor` reports project status with actionable repair path.
- [x] Inventory persistence is treated as transactional state (`weave.yaml` remains the source of truth and drift is diagnosable).
- [x] In v1 strict mode, symlink operation + config persistence behave as one logical transaction (rollback link on config persistence failure).
- [x] Exit-code mapping is deterministic for doctor and service error classes.
- [x] `doctor --json` emits machine-readable diagnostics.
- [x] Error messages include docs references (path + URL).

### B4 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-UX-02 (`doctor` explains status + repair) | B4-T1.1, B4-T1.2 | B4-T1.3 | B4-T1.4 | B4-T1.5 | [x] |
| R-ARCH-06 (inventory persistence transactional state) | B4-T2.1, B4-T2.2 | B4-T2.3 | B4-T2.4 | B4-T2.5 | [x] |
| R-ARCH-07 (symlink + config as single logical transaction) | B4-T3.1, B4-T3.2 | B4-T3.3 | B4-T3.4 | B4-T3.5 | [x] |
| R-NFR-02 (deterministic exit codes) | B4-T4.1, B4-T4.2 | B4-T4.3 | B4-T4.4 | B4-T4.5 | [x] |
| R-UX-04 (script-friendly `--json`) | B4-T5.1, B4-T5.2 | B4-T5.3 | B4-T5.4 | B4-T5.5 | [x] |
| R-POST-03 (errors include docs references) | B4-T6.1, B4-T6.2 | B4-T6.3 | B4-T6.4 | B4-T6.5 | [x] |

### B4 Tasks (live checklist)

#### R-UX-02 — `doctor` MUST explain current status and repair path

- [x] **B4-T1.1 Unit (success):** doctor returns `healthy` with empty issues/repair list for converged project state.
- [x] **B4-T1.2 Unit (edge):** doctor returns issue diagnostics and repair commands when expected symlink is missing.
- [x] **B4-T1.3 Integration:** CLI `doctor` path maps service result into deterministic status text.
- [x] **B4-T1.4 E2E:** `weave doctor` on inconsistent state prints issues and suggested repair commands.
- [x] **B4-T1.5 Evidence:** output contains status (`healthy` or `issues_found`) and explicit repair path (`weave skill add ...` / `weave command add ...`).

Implementation evidence:

- [x] `internal/app/doctor_test.go::TestDoctorService_Run_HealthyProjectReturnsNoIssues`
- [x] `internal/app/doctor_test.go::TestDoctorService_Run_MissingSkillSymlinkReturnsRepairGuidance`
- [x] `internal/cli/doctor_test.go::TestRunDoctor_MissingAssetReturnsExitDoctorIssuesAndRepairPath`
- [x] `test/e2e/doctor_e2e_test.go::TestDoctor_E2E_ReportsIssuesAndDeterministicExitCode`

#### R-ARCH-06 — inventory persistence in `weave.yaml` MUST be transactional state

- [x] **B4-T2.1 Unit (success):** config persistence still succeeds after symlink apply on normal path.
- [x] **B4-T2.2 Unit (error):** config persistence failure triggers rollback operation for installed link.
- [x] **B4-T2.3 Integration:** failed persistence leaves no newly installed link and no new `weave.yaml` state.
- [x] **B4-T2.4 E2E:** doctor can detect inventory/install drift as recoverable state.
- [x] **B4-T2.5 Evidence:** rollback remove op is executed and missing-link issue becomes observable via diagnostics.

Implementation evidence:

- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_ConfigWriteFailureRollsBackSymlink`
- [x] `internal/app/asset_add.go` rollback-on-writer-failure flow
- [x] `test/integration/doctor_integration_test.go::TestDoctor_Integration_ConfigWriteFailureRollsBackSymlinkInStrictMode`

#### R-ARCH-07 — symlink operation + config write MUST behave as one logical transaction in v1 strict mode

- [x] **B4-T3.1 Unit (success):** add flow preserves existing strict ordering (`create_link` then persist config).
- [x] **B4-T3.2 Unit (error):** add flow compensates with rollback when persistence fails after fs apply.
- [x] **B4-T3.3 Integration:** strict-mode add failure does not leave partially committed state.
- [x] **B4-T3.4 E2E:** issue state is deterministic and recoverable through doctor/repair commands.
- [x] **B4-T3.5 Evidence:** atomic semantic represented by either (a) link+config committed, or (b) neither committed.

Implementation evidence:

- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_PersistsConfigAfterSymlinkSuccess`
- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_ConfigWriteFailureRollsBackSymlink`
- [x] `test/integration/doctor_integration_test.go::TestDoctor_Integration_ConfigWriteFailureRollsBackSymlinkInStrictMode`

#### R-NFR-02 — exit codes MUST be deterministic

- [x] **B4-T4.1 Unit (success):** static exit constants include doctor issue class and remain stable.
- [x] **B4-T4.2 Unit (edge):** error-class mapping function deterministically maps invalid config, missing dependencies, and runtime failures.
- [x] **B4-T4.3 Integration:** CLI doctor returns stable issue exit code (`5`) without treating issue-reporting as runtime error.
- [x] **B4-T4.4 E2E:** shell-visible process exit code for doctor issue state is deterministic.
- [x] **B4-T4.5 Evidence:** repeated failing doctor runs return the same process exit code (`5`).

Implementation evidence:

- [x] `internal/cli/exitcodes_test.go::TestExitCodes_AreStable`
- [x] `internal/cli/exitcodes_test.go::TestExitCodes_ErrorMappingIsDeterministic`
- [x] `test/integration/doctor_integration_test.go::TestDoctor_Integration_ExitCodeMappingForIssueState`
- [x] `test/e2e/doctor_e2e_test.go::TestDoctor_E2E_ReportsIssuesAndDeterministicExitCode`

#### R-UX-04 — relevant commands MUST support script-friendly output (`--json`, at least doctor)

- [x] **B4-T5.1 Unit (success):** `doctor --json` produces valid JSON payload with status/diagnostics schema.
- [x] **B4-T5.2 Unit (edge):** unknown flags are rejected deterministically.
- [x] **B4-T5.3 Integration:** doctor CLI printer switches between human text and JSON output deterministically.
- [x] **B4-T5.4 E2E:** `weave doctor --json` emits parseable JSON for healthy state.
- [x] **B4-T5.5 Evidence:** JSON includes `status`, `issues`, and `repair_commands` keys.

Implementation evidence:

- [x] `internal/cli/doctor_test.go::TestRunDoctor_JSONOutput`
- [x] `internal/cli/root.go::runDoctor`
- [x] `test/e2e/doctor_e2e_test.go::TestDoctor_E2E_JSONOutputIsScriptFriendly`

#### R-POST-03 — errors MUST include docs reference paths/URLs

- [x] **B4-T6.1 Unit (success):** invalid-config wrapper includes docs path + URL in emitted error.
- [x] **B4-T6.2 Unit (edge):** provider/setup transactional errors include docs references for troubleshooting.
- [x] **B4-T6.3 Integration:** doctor issue diagnostics include both docs path and canonical docs URL.
- [x] **B4-T6.4 E2E:** doctor issue output includes docs references in human-readable mode.
- [x] **B4-T6.5 Evidence:** user-facing errors/diagnostics include `docs/reference/...` and `https://.../docs/...`.

Implementation evidence:

- [x] `internal/app/docs_refs.go`
- [x] `internal/app/errors.go::WrapInvalidConfig`
- [x] `internal/app/provider_service.go` error paths
- [x] `internal/app/doctor_test.go::TestDoctorService_Run_MissingSkillSymlinkReturnsRepairGuidance`
- [x] `test/e2e/doctor_e2e_test.go::TestDoctor_E2E_ReportsIssuesAndDeterministicExitCode`

---

## Progress Log

- [x] Batch 1 started
- [x] Batch 1 completed
- [x] Batch 2 started
- [x] Batch 2 completed
- [x] Batch 3 started
- [x] Batch 3 completed
- [x] Batch 4 started
- [x] Batch 4 completed
- [x] Batch 5 started
- [x] Batch 5 completed
- [x] Batch 6 started
- [x] Batch 6 completed
- [x] Batch 7 started
- [x] Batch 7 completed

---

## Batch 5 — Safety guards + dry-run + shell-agnostic behavior

### B5 Exit Criteria

- [x] CLI detects project root from nested working directories and rejects mutating runs outside repository scope.
- [x] All mutating commands support deterministic `--dry-run` behavior without filesystem/config mutation.
- [x] Mutating command summaries are concise, actionable, and consistent across shell invocation styles.
- [x] Strict-mode rollback output explicitly states whether rollback completed or partial state remains.

### B5 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-DEP-01 (detect project root + safe scope guard) | B5-T1.1, B5-T1.2 | B5-T1.3 | B5-T1.4 | B5-T1.5 | [x] |
| R-DEP-03 (`--dry-run` for mutating commands) | B5-T2.1, B5-T2.2 | B5-T2.3 | B5-T2.4 | B5-T2.5 | [x] |
| R-DEP-07 (shell-agnostic core CLI commands) | B5-T3.1, B5-T3.2 | B5-T3.3 | B5-T3.4 | B5-T3.5 | [x] |
| R-UX-01 (concise summary for mutating commands) | B5-T4.1, B5-T4.2 | B5-T4.3 | B5-T4.4 | B5-T4.5 | [x] |
| R-UX-03 (actionability-first UX language) | B5-T5.1, B5-T5.2 | B5-T5.3 | B5-T5.4 | B5-T5.5 | [x] |
| R-UX-05 (strict-mode rollback semantics in output) | B5-T6.1, B5-T6.2 | B5-T6.3 | B5-T6.4 | B5-T6.5 | [x] |

### B5 Tasks (live checklist)

#### R-DEP-01 — detect project root + safe scope guard

- [x] **B5-T1.1 Unit (success):** root detector resolves repository root from nested workdir.
- [x] **B5-T1.2 Unit (edge):** mutation guard rejects fs operations outside detected project root.
- [x] **B5-T1.3 Integration:** service-layer mutating flows enforce root-scoped operation guards before apply.
- [x] **B5-T1.4 E2E:** mutating commands run from nested contexts still operate on repository root.
- [x] **B5-T1.5 Evidence:** deterministic `ErrUnsafeMutationPath` on out-of-root operation plans.

Implementation evidence:

- [x] `internal/cli/forge_handler.go::detectProjectRootFrom`
- [x] `internal/app/mutation_guard.go`
- [x] `internal/app/mutation_guard_test.go::TestEnsureOpsWithinRoot_RejectsPathOutsideRoot`
- [x] `internal/app/forge_test.go::TestForgeService_Run_RejectsOpsOutsideDetectedRoot`

#### R-DEP-03 — `--dry-run` for all mutating commands

- [x] **B5-T2.1 Unit (success):** forge/asset/provider services return planned operation counts in dry-run without apply/write.
- [x] **B5-T2.2 Unit (edge):** dry-run parsing rejects unsupported combinations (e.g. `provider list --dry-run`).
- [x] **B5-T2.3 Integration:** dry-run mode preserves existing `weave.yaml` bytes and does not create symlinks/projection dirs.
- [x] **B5-T2.4 E2E:** `forge`, `skill add`, and `provider add` dry-runs emit success exit code with no mutations.
- [x] **B5-T2.5 Evidence:** pre/post filesystem + config snapshots remain unchanged under dry-run.

Implementation evidence:

- [x] `internal/app/forge.go::RunWithOptions`
- [x] `internal/app/asset_add.go::AddAssetWithOptions`
- [x] `internal/app/provider_service.go::*WithOptions`
- [x] `internal/cli/root.go::parseDryRunOnly|parseAddFlags|parseProviderAction`
- [x] `test/e2e/forge_e2e_test.go::TestForge_E2E_DryRunDoesNotMutateAndPrintsActionableSummary`
- [x] `test/e2e/provider_e2e_test.go::TestProviderAdd_E2E_DryRunDoesNotMutateAndPrintsActionableSummary`

#### R-DEP-07 — shell-agnostic core CLI commands

- [x] **B5-T3.1 Unit (success):** argument parsing accepts explicit flag/value and equals forms deterministically.
- [x] **B5-T3.2 Unit (edge):** unsupported flags/extra args fail with deterministic parsing errors.
- [x] **B5-T3.3 Integration:** command dispatch remains independent of shell-specific wrappers/aliases.
- [x] **B5-T3.4 E2E:** command behavior is stable across `go run ./cmd/weave <args>` entrypoints.
- [x] **B5-T3.5 Evidence:** repeated invocations with equivalent args produce identical summaries and state transitions.

Implementation evidence:

- [x] `internal/cli/root_test.go::TestParseProviderAction_AddParsesDryRun`
- [x] `internal/cli/root_test.go::TestParseDryRunOnly_RejectsUnknownFlag`
- [x] `test/e2e/forge_e2e_test.go` dry-run + precedence flows

#### R-UX-01 — concise summary for mutating commands

- [x] **B5-T4.1 Unit (success):** forge summary reports applied/planned counts in one concise line.
- [x] **B5-T4.2 Unit (edge):** no-op and dry-run variants produce concise summary text.
- [x] **B5-T4.3 Integration:** provider and asset commands print deterministic concise summaries after successful execution.
- [x] **B5-T4.4 E2E:** mutating command output includes single-line concise status summaries.
- [x] **B5-T4.5 Evidence:** output assertions match concise summary prefixes (`forge:`, `skill add`, `provider added`).

Implementation evidence:

- [x] `internal/cli/root.go::printForgeSummary`
- [x] `internal/cli/root.go::printAssetAddSummary`
- [x] `internal/cli/root.go::printProviderSummary`

#### R-UX-03 — actionability-first UX language

- [x] **B5-T5.1 Unit (success):** dry-run summaries include explicit next step (`rerun without --dry-run`).
- [x] **B5-T5.2 Unit (edge):** argument/flag errors identify exactly what to change.
- [x] **B5-T5.3 Integration:** transactional/provider failure messages preserve repair command guidance.
- [x] **B5-T5.4 E2E:** user-facing dry-run and failure outputs include actionable next commands.
- [x] **B5-T5.5 Evidence:** assertions verify guidance text in CLI output.

Implementation evidence:

- [x] `internal/cli/root.go` summary strings
- [x] `internal/cli/root_test.go` parsing failure tests
- [x] `test/e2e/forge_e2e_test.go::TestSkillAdd_E2E_DryRunDoesNotMutateAndPrintsActionableSummary`

#### R-UX-05 — strict-mode rollback semantics in output

- [x] **B5-T6.1 Unit (success):** rollback-complete errors explicitly state no config/symlink changes were committed.
- [x] **B5-T6.2 Unit (edge):** rollback-failure errors explicitly warn about possible partial state and immediate repair path.
- [x] **B5-T6.3 Integration:** rollback semantics surface from service boundary to CLI stderr unchanged.
- [x] **B5-T6.4 E2E:** strict rollback failure output remains actionable and deterministic.
- [x] **B5-T6.5 Evidence:** output includes semantic markers (`rollback completed` or `project may be partially modified`).

Implementation evidence:

- [x] `internal/app/asset_add.go` rollback error messaging
- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_ConfigWriteFailureRollsBackSymlink`

---

## Batch 6 — UX hardening + architecture guardrails + release/update requirements

### B6 Exit Criteria

- [x] Skill/command add supports conflict resolution prompt path and deterministic non-interactive policies.
- [x] Conflict handling, backup-on-write, and dry-run reuse one planning primitive.
- [x] Command inventory supports future-ready metadata markers for provider compatibility.
- [x] Help and first-run output include contextual quickstart guidance.
- [x] Distribution/update baseline artifacts exist: signed-checksum policy, install/troubleshooting docs, semver naming, migration guide.
- [x] Doctor reports stale provider integrations and unknown enabled providers post-upgrade.

### B6 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-SKILL-01 (conflict prompt overwrite/skip/backup) | B6-T1.1, B6-T1.2 | B6-T1.3 | B6-T1.4 | B6-T1.5 | [x] |
| R-SKILL-02 (non-interactive conflict flags) | B6-T2.1, B6-T2.2 | B6-T2.3 | B6-T2.4 | B6-T2.5 | [x] |
| R-CMD-02 (command metadata provider-compat markers) | B6-T3.1, B6-T3.2 | B6-T3.3 | B6-T3.4 | B6-T3.5 | [x] |
| R-ARCH-01 (UI-agnostic core business logic) | B6-T4.1, B6-T4.2 | B6-T4.3 | B6-T4.4 | B6-T4.5 | [x] |
| R-ARCH-02 (no provider leakage into command handlers) | B6-T5.1, B6-T5.2 | B6-T5.3 | B6-T5.4 | B6-T5.5 | [x] |
| R-ARCH-03 (dry-run/real-run share planning primitives) | B6-T6.1, B6-T6.2 | B6-T6.3 | B6-T6.4 | B6-T6.5 | [x] |
| R-ARCH-05 (`.agents` canonical, provider dirs projections) | B6-T7.1, B6-T7.2 | B6-T7.3 | B6-T7.4 | B6-T7.5 | [x] |
| R-NFR-01 (backup-on-write for risky paths) | B6-T8.1, B6-T8.2 | B6-T8.3 | B6-T8.4 | B6-T8.5 | [x] |
| R-NFR-03 (critical path unit coverage) | B6-T9.1, B6-T9.2 | B6-T9.3 | B6-T9.4 | B6-T9.5 | [x] |
| R-POST-01 (help includes 60-second quickstart) | B6-T10.1, B6-T10.2 | B6-T10.3 | B6-T10.4 | B6-T10.5 | [x] |
| R-POST-02 (first-run contextual next-step guidance) | B6-T11.1, B6-T11.2 | B6-T11.3 | B6-T11.4 | B6-T11.5 | [x] |
| R-DIST-01 (signed checksums baseline) | B6-T12.1, B6-T12.2 | B6-T12.3 | B6-T12.4 | B6-T12.5 | [-] |
| R-DIST-02 (install docs quickstart + troubleshooting) | B6-T13.1, B6-T13.2 | B6-T13.3 | B6-T13.4 | B6-T13.5 | [-] |
| R-DIST-03 (semver binary naming/versioning) | B6-T14.1, B6-T14.2 | B6-T14.3 | B6-T14.4 | B6-T14.5 | [x] |
| R-UPD-01 (detect outdated schema + suggest/perform migration) | B6-T15.1, B6-T15.2 | B6-T15.3 | B6-T15.4 | B6-T15.5 | [x] |
| R-UPD-02 (breaking changes include migration guide) | B6-T16.1, B6-T16.2 | B6-T16.3 | B6-T16.4 | B6-T16.5 | [-] |
| R-UPD-03 (doctor flags stale provider integrations) | B6-T17.1, B6-T17.2 | B6-T17.3 | B6-T17.4 | B6-T17.5 | [x] |

### B6 Tasks (live checklist)

#### R-SKILL-01 — conflict prompt (overwrite/skip/backup)

- [x] **B6-T1.1 Unit (success):** conflict planner returns overwrite/skip/backup plans deterministically.
- [x] **B6-T1.2 Unit (error/edge):** prompt-required error is raised when conflict policy is prompt and no prompter is available.
- [x] **B6-T1.3 Integration:** skill add with pre-existing destination and `--skip` performs no write/mutation.
- [x] **B6-T1.4 E2E:** existing destination with conflict policy path returns deterministic summary.
- [x] **B6-T1.5 Evidence:** output includes “conflict detected … skipped” or backup path marker.

Implementation evidence:

- [x] `internal/app/conflicts_test.go::TestConflictPlanner_BackupPolicyProducesBackupAndCreateLinkOps`
- [x] `internal/app/conflicts_test.go::TestResolveConflictPolicy_PromptWithoutPrompterFails`
- [x] `internal/cli/assets_handler_test.go::TestAssetAddService_Add_ConflictSkipPolicyProducesNoMutation`

#### R-SKILL-02 — non-interactive conflict flags

- [x] **B6-T2.1 Unit (success):** parser accepts `--overwrite`, `--skip`, `--backup`.
- [x] **B6-T2.2 Unit (edge):** parser rejects multiple conflict flags in one invocation.
- [x] **B6-T2.3 Integration:** conflict flag is propagated from CLI parser to AddAsset service input.
- [x] **B6-T2.4 E2E:** non-interactive conflict policy execution is deterministic.
- [x] **B6-T2.5 Evidence:** parse error and execution summaries match selected policy.

Implementation evidence:

- [x] `internal/cli/root_test.go::TestParseAddFlags_ParsesConflictPolicyFlags`
- [x] `internal/cli/root_test.go::TestParseAddFlags_RejectsMultipleConflictPolicyFlags`
- [x] `internal/cli/root.go::parseAddFlags`

#### R-CMD-02 — command metadata provider-compat markers (future-ready baseline)

- [x] **B6-T3.1 Unit (success):** command config model supports metadata/provider_compat field.
- [x] **B6-T3.2 Unit (edge):** metadata normalization removes empty/duplicate entries and sorts values.
- [x] **B6-T3.3 Integration:** command add path persists baseline metadata struct without provider leakage.
- [x] **B6-T3.4 E2E:** command inventory remains schema-compatible with metadata envelope.
- [x] **B6-T3.5 Evidence:** deterministic YAML includes normalized `metadata.provider_compat` when provided.

Implementation evidence:

- [x] `internal/config/config.go::Asset.Meta`
- [x] `internal/config/io.go::normalizeAssetMetadata`
- [x] `internal/config/io_test.go::TestMarshalDeterministic_CommandMetadataProviderCompatIsNormalized`

#### R-ARCH-01 — UI-agnostic core business logic

- [x] **B6-T4.1 Unit (success):** conflict planning and migration logic implemented in `internal/app` services.
- [x] **B6-T4.2 Unit (edge):** CLI-specific prompt fallback remains adapter-side, not domain-side.
- [x] **B6-T4.3 Integration:** CLI handlers invoke app services without embedding business rules.
- [x] **B6-T4.4 E2E:** CLI observable behavior derived from service outputs.
- [x] **B6-T4.5 Evidence:** `internal/app/conflicts.go` + `internal/app/migrate.go` are reusable by future UI.

#### R-ARCH-02 — no provider leakage into command handlers

- [x] **B6-T5.1 Unit (success):** provider diagnostics in doctor use provider registry interface only.
- [x] **B6-T5.2 Unit (edge):** unknown provider path emits actionable issue without provider-specific branching in CLI.
- [x] **B6-T5.3 Integration:** CLI `doctor` delegates provider checks to `app.DoctorService`.
- [x] **B6-T5.4 E2E:** stale provider diagnosis is emitted through generic doctor issue channel.
- [x] **B6-T5.5 Evidence:** provider-specific path remains in adapters/registry.

#### R-ARCH-03 — dry-run/real-run share planning primitives

- [x] **B6-T6.1 Unit (success):** conflict planning returns operation list reused by dry-run and real-run.
- [x] **B6-T6.2 Unit (edge):** skip policy returns empty plan and short-circuits both modes.
- [x] **B6-T6.3 Integration:** AddAsset applies same plan in dry-run and execution path.
- [x] **B6-T6.4 E2E:** dry-run output reflects planned operation counts from shared primitive.
- [x] **B6-T6.5 Evidence:** `app.AddAssetWithOptions` computes `ops` once before mode branch.

#### R-ARCH-05 — `.agents` canonical, provider dirs are projections

- [x] **B6-T7.1 Unit (success):** doctor validates projection links expected from provider adapter plan.
- [x] **B6-T7.2 Unit (edge):** missing/non-symlink provider projection emits stale integration issue.
- [x] **B6-T7.3 Integration:** doctor + provider registry checks `.claude` / `.opencode` as projections.
- [x] **B6-T7.4 E2E:** stale provider projection guidance uses `weave provider repair <provider>`.
- [x] **B6-T7.5 Evidence:** issue code `stale_provider_integration` with repair command.

Implementation evidence:

- [x] `internal/app/doctor.go::diagnoseProviderProjections`
- [x] `internal/cli/doctor_test.go::TestRunDoctor_StaleProviderIntegrationIsReported`

#### R-NFR-01 — backup-on-write for risky paths

- [x] **B6-T8.1 Unit (success):** fs engine supports `backup_path` operation.
- [x] **B6-T8.2 Unit (edge):** backup op preserves original payload in backup destination.
- [x] **B6-T8.3 Integration:** conflict backup policy plans backup + create_link sequence.
- [x] **B6-T8.4 E2E:** conflict backup path included in command summary/evidence.
- [x] **B6-T8.5 Evidence:** backup file naming with timestamp suffix (`.bak.<UTC>`).

Implementation evidence:

- [x] `internal/fsops/ops.go::OpBackupPath`
- [x] `internal/fsops/engine_test.go::TestEngine_Apply_BackupPathRenamesExistingPath`
- [x] `internal/app/forge_test.go::TestForgeService_AddAsset_ConflictBackupPolicyPlansBackupAndCreateLink`

#### R-NFR-03 — unit coverage for critical paths

- [x] **B6-T9.1 Unit (success):** conflict planner coverage (backup/prompt/policy resolution).
- [x] **B6-T9.2 Unit (edge):** migration coverage (dry-run/no-op/upgrade paths).
- [x] **B6-T9.3 Integration:** doctor integration retains deterministic issue mapping.
- [x] **B6-T9.4 E2E:** full command paths continue to pass end-to-end tests.
- [x] **B6-T9.5 Evidence:** `go test ./...` and `go test -tags=e2e ./test/e2e/...` pass.

#### R-POST-01 — help includes 60-second quickstart

- [x] **B6-T10.1 Unit (success):** `--help` output includes quickstart block.
- [x] **B6-T10.2 Unit (edge):** help path exits `0` deterministically.
- [x] **B6-T10.3 Integration:** root command dispatch supports explicit help aliases.
- [x] **B6-T10.4 E2E:** quickstart commands are valid in current command surface.
- [x] **B6-T10.5 Evidence:** help text contains “60-second quickstart”.

Implementation evidence:

- [x] `internal/cli/root.go::printHelp`
- [x] `internal/cli/root_test.go::TestRun_HelpPrintsQuickstart`

#### R-POST-02 — first-run contextual next-step guidance

- [x] **B6-T11.1 Unit (success):** forge flow detects first-run when `weave.yaml` absent pre-run.
- [x] **B6-T11.2 Unit (edge):** guidance omitted for non-first runs and dry-run.
- [x] **B6-T11.3 Integration:** contextual next-step guidance printed after first successful forge.
- [x] **B6-T11.4 E2E:** guidance text maps to actionable provider/skill/doctor commands.
- [x] **B6-T11.5 Evidence:** output includes “First run complete. Suggested next steps”.

#### R-DIST-01 — signed checksums for release artifacts

- [x] **B6-T12.1 Unit (success):** distribution policy doc defines checksums + signature artifacts.
- [x] **B6-T12.2 Unit (edge):** checklist requires signed checksum publication.
- [x] **B6-T12.3 Integration:** README links to distribution baseline docs.
- [-] **B6-T12.4 E2E:** release-signing automation not implemented in v1 baseline (tracked in `deferred.md`).
- [x] **B6-T12.5 Evidence:** `docs/reference/distribution.md` signed-checksum section.

#### R-DIST-02 — install docs include quickstart + troubleshooting

- [x] **B6-T13.1 Unit (success):** install doc includes quickstart command set.
- [x] **B6-T13.2 Unit (edge):** troubleshooting includes missing binaries/schema migration/stale projection flows.
- [x] **B6-T13.3 Integration:** README docs index references install guide.
- [-] **B6-T13.4 E2E:** installation artifact pipeline validation deferred (tracked in `deferred.md`).
- [x] **B6-T13.5 Evidence:** `docs/reference/install.md`.

#### R-DIST-03 — semver binary naming/versioning

- [x] **B6-T14.1 Unit (success):** CLI version baseline exposed as semver string.
- [x] **B6-T14.2 Unit (edge):** distribution doc codifies artifact naming contract.
- [x] **B6-T14.3 Integration:** root command supports `--version` output.
- [x] **B6-T14.4 E2E:** version can be invoked in CLI-first flow.
- [x] **B6-T14.5 Evidence:** `internal/cli/version.go` + `docs/reference/distribution.md` naming examples.

#### R-UPD-01 — detect outdated schema + suggest/perform migration

- [x] **B6-T15.1 Unit (success):** validator classifies outdated schema with migration guidance.
- [x] **B6-T15.2 Unit (edge):** validator classifies unsupported future schema.
- [x] **B6-T15.3 Integration:** migrate command upgrades schema and supports dry-run.
- [x] **B6-T15.4 E2E:** migration command path available in CLI command set.
- [x] **B6-T15.5 Evidence:** `weave migrate` output and wrapped invalid-config references migration docs.

Implementation evidence:

- [x] `internal/config/validator.go::ErrOutdatedSchema|ErrUnsupportedSchema`
- [x] `internal/app/migrate.go`
- [x] `internal/cli/root.go::runMigrate`

#### R-UPD-02 — breaking changes include migration guide

- [x] **B6-T16.1 Unit (success):** migration guide doc exists in docs baseline.
- [x] **B6-T16.2 Unit (edge):** release notes policy includes mandatory migration section for breaking releases.
- [x] **B6-T16.3 Integration:** docs references available from README/docs index.
- [-] **B6-T16.4 E2E:** release-note enforcement automation deferred (tracked in `deferred.md`).
- [x] **B6-T16.5 Evidence:** `docs/reference/migration.md`, `docs/reference/releases.md`.

#### R-UPD-03 — doctor flags stale provider integrations after upgrades

- [x] **B6-T17.1 Unit (success):** doctor reports unknown enabled provider when registry lacks adapter.
- [x] **B6-T17.2 Unit (edge):** doctor reports stale provider projections (missing/not symlink/wrong target).
- [x] **B6-T17.3 Integration:** CLI doctor exits with `ExitDoctorIssues` for stale provider state.
- [x] **B6-T17.4 E2E:** doctor output includes actionable `weave provider repair <provider>` guidance.
- [x] **B6-T17.5 Evidence:** issue code `stale_provider_integration` in doctor output.

Implementation evidence:

- [x] `internal/app/doctor_test.go::TestDoctorService_Run_UnknownEnabledProviderFlagsStaleIntegration`
- [x] `internal/cli/doctor_test.go::TestRunDoctor_StaleProviderIntegrationIsReported`

---

## Batch 7 — Deferred v1.x hardening

### B7 Exit Criteria

- [x] Core app/domain contracts remain reusable by future UI adapters (including a TUI) without refactoring service boundaries.
- [x] v1 runtime and module graph remain free of TUI runtime dependencies.
- [x] One-command install flow exists and is documented consistently across script + docs.

### B7 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-ARCH-04 (TUI integration possible without domain refactor) | B7-T1.1, B7-T1.2 | B7-T1.3 | B7-T1.4 | B7-T1.5 | [x] |
| R-DEP-05 (no TUI runtime dependency in v1 binary) | B7-T2.1, B7-T2.2 | B7-T2.3 | B7-T2.4 | B7-T2.5 | [x] |
| R-DEP-06 (one-command install for Weave binary) | B7-T3.1, B7-T3.2 | B7-T3.3 | B7-T3.4 | B7-T3.5 | [x] |

### B7 Tasks (live checklist)

#### R-ARCH-04 — TUI integration MUST be possible without refactoring domain contracts

- [x] **B7-T1.1 Unit (success):** boundary test verifies `internal/app` has no `internal/cli` imports.
- [x] **B7-T1.2 Unit (edge):** boundary test verifies `internal/app` has no direct imports of common TUI modules.
- [x] **B7-T1.3 Integration:** repository-level integration guard asserts app layer remains reusable outside CLI adapter boundaries.
- [x] **B7-T1.4 E2E:** existing command-surface E2E suite remains green after adding boundary guards (no command behavior regression).
- [x] **B7-T1.5 Evidence:** `internal/cli/batch7_boundaries_test.go` and `test/integration/batch7_boundary_integration_test.go` enforce adapter boundaries.

Implementation evidence:

- [x] `internal/cli/batch7_boundaries_test.go::TestArchitectureBoundary_AppLayerHasNoCLINorTUIImports`
- [x] `test/integration/batch7_boundary_integration_test.go::TestBatch7_Integration_AppLayerIsReusableOutsideCLI`

#### R-DEP-05 — v1 binary MUST have no TUI runtime dependency

- [x] **B7-T2.1 Unit (success):** dependency guard test verifies `go.mod` does not include known TUI runtime modules.
- [x] **B7-T2.2 Unit (edge):** runtime-source guard test verifies runtime code does not import known TUI modules.
- [x] **B7-T2.3 Integration:** integration guard verifies module-level dependency contract remains TUI-free.
- [x] **B7-T2.4 E2E:** existing e2e command suite executes with current non-TUI runtime graph.
- [x] **B7-T2.5 Evidence:** `go.mod` remains minimal (`gopkg.in/yaml.v3` only) and tests fail-fast on any future TUI import/dependency leak.

Implementation evidence:

- [x] `internal/cli/batch7_boundaries_test.go::TestRuntimeBoundary_V1HasNoTUIRuntimeDependency`
- [x] `test/integration/batch7_boundary_integration_test.go::TestBatch7_Integration_GoModuleStaysFreeOfTUIRuntimeDeps`

#### R-DEP-06 — distribution MUST support one-command install for Weave binary

- [x] **B7-T3.1 Unit (success):** install-flow test verifies `scripts/install.sh` exists, is executable, and installs with `go install ./cmd/weave`.
- [x] **B7-T3.2 Unit (edge):** install-flow test verifies docs alignment between README and install reference for command path + verification commands.
- [x] **B7-T3.3 Integration:** docs/script alignment guard ensures install script and docs stay synchronized.
- [x] **B7-T3.4 E2E:** `go test -tags=e2e ./test/e2e/...` remains green with install-flow changes.
- [x] **B7-T3.5 Evidence:** `scripts/install.sh`, `README.md`, and `docs/reference/install.md` provide one-command install + verification path.

Implementation evidence:

- [x] `internal/cli/batch7_boundaries_test.go::TestInstallFlow_OneCommandScriptAndDocsAreAligned`
- [x] `scripts/install.sh`
- [x] `README.md`
- [x] `docs/reference/install.md`

---

## Implementation Notes

- Keep this file updated during implementation (live checklist behavior).
- Only mark a requirement block complete after all test levels + evidence are complete.
- Add regression tasks under the relevant requirement when bugs are found.
- Register all intentionally deferred executable work in `deferred.md`.
