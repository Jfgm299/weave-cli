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
| B2 | P0 | R-SKILL-03 | configurable sources via config/flag/env | [x] |
| B2 | P0 | R-CMD-01 | commands lifecycle mirrors skills | [x] |
| B3 | P0 | R-PROV-01 | multi-provider enable in one project | [x] |
| B3 | P0 | R-PROV-02 | provider adapter architecture | [x] |
| B3 | P1 | R-PROV-03 | reversible provider operations (remove/repair) | [x] |
| B3 | P0 | R-DEP-02 | validate provider binaries before success claim | [x] |
| B3 | P1 | R-DEP-04 | actionable failure messages | [x] |
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

## Progress Log

- [x] Batch 1 started
- [x] Batch 1 completed
- [x] Batch 2 started
- [x] Batch 2 completed
- [x] Batch 3 started
- [x] Batch 3 completed
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
