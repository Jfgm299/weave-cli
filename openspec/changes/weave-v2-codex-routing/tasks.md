# Tasks — weave-v2-codex-routing

## Scope

Implement v2 delta with debt-first policy:

1. Close v1 deferred debt (`DEF-001..DEF-004`).
2. Add provider `codex`.
3. Make `command add` provider-aware while preserving UX.
4. Guarantee strict rollback and coherent doctor/repair behavior.

## Status Legend

- [ ] Pending
- [-] In progress
- [x] Done (100% implemented, tested, and operational)
- [!] Blocked

## Completion Policy

- Never mark `[x]` unless code + tests + acceptance evidence are complete.
- Deferred or baseline-only closure MUST stay `[-]` and be tracked in `deferred.md`.
- For this change, Batch 0 debt closure is mandatory before feature batches.

---

## Master Backlog

| Batch | Priority | Requirement | Title | Status |
|------:|:--------:|-------------|-------|:------:|
| B0 | P0 | R-DEBT-01 | Close DEF-001/002 CI automation debt | [x] |
| B0 | P0 | R-DEBT-02 | Close DEF-003 migration-note CI gate debt | [x] |
| B0 | P0 | R-DEBT-03 | Close DEF-004 guided git-init flow debt | [x] |
| B1 | P0 | R-PROV-04 | Add codex provider parity | [x] |
| B1 | P0 | R-PROV-05 | Provider command projection strategy contract | [x] |
| B2 | P0 | R-CMD-05 | `command add` default installs for all enabled providers | [x] |
| B2 | P0 | R-CMD-06 | `--provider` exclusive install without shared `.agents` requirement | [x] |
| B2 | P0 | R-CMD-07 | Codex namespace `__weave_commands__` | [x] |
| B2 | P0 | R-CMD-08 | Interactive no-provider prompt `[y/N]` | [x] |
| B2 | P0 | R-CMD-09 | Non-interactive no-provider auto-fail | [x] |
| B3 | P0 | R-ARCH-08 | Keep provider logic out of CLI handlers | [x] |
| B3 | P0 | R-ARCH-09 | Full rollback on multi-provider failure | [x] |
| B3 | P0 | R-CONFIG-08 | Transactional provider-target persistence | [x] |

---

## Batch Plan Overview

- **Batch 0 (B0) — Debt closure first (P0)**
- **Batch 1 (B1) — Codex provider integration baseline (P0)**
- **Batch 2 (B2) — Provider-aware command add behavior (P0)**
- **Batch 3 (B3) — Transaction/doctor consistency hardening (P0)**
- **Batch 4 (B4) — Final verification + docs + traceability closure (P1)**

---

## Batch 0 — Debt closure first (current)

### B0 Exit Criteria

- [x] DEF-001 and DEF-002 closure criteria are implemented and evidenced via local tests + live CI run URLs.
- [x] DEF-003 migration-note gate workflow is implemented and evidenced via local tests + live CI run URL.
- [x] DEF-004 guided git-init behavior is implemented with unit + targeted e2e evidence (interactive and non-interactive branches).
- [x] `openspec/changes/weave-v1-mvp/deferred.md` status is updated with closure evidence and explicit closure statement.

### B0 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E/CI | Acceptance Evidence | Status |
|-------------|------|-------------|--------|---------------------|--------|
| R-DEBT-01 (DEF-001/002) | B0-T1.1, B0-T1.2 | B0-T1.3 | B0-T1.4 | B0-T1.5 | [x] |
| R-DEBT-02 (DEF-003) | B0-T2.1, B0-T2.2 | B0-T2.3 | B0-T2.4 | B0-T2.5 | [x] |
| R-DEBT-03 (DEF-004) | B0-T3.1, B0-T3.2 | B0-T3.3 | B0-T3.4 | B0-T3.5 | [x] |

### B0 Tasks (live checklist)

#### R-DEBT-01 — close DEF-001 and DEF-002 (release signing/checksums + install validation)

- [x] **B0-T1.1 Unit (success):** release artifact metadata generation includes expected checksums/signature references.
- [x] **B0-T1.2 Unit (edge):** release validation fails on missing checksum/signature artifacts.
- [x] **B0-T1.3 Integration:** CI workflow job validates install flow against produced artifacts.
- [x] **B0-T1.4 CI/E2E:** workflow run demonstrates pass path and fail path behavior.
- [x] **B0-T1.5 Evidence:** workflow definitions + local script test logs + live CI run URLs captured.

Implementation evidence:

- [x] `.github/workflows/release-artifacts.yml`
- [x] `.github/workflows/install-artifact-validation.yml`
- [x] `scripts/release/generate_artifacts.sh`
- [x] `scripts/release/verify_release_artifacts.sh`
- [x] `scripts/release/release_artifacts_test.py`

#### R-DEBT-02 — close DEF-003 (migration-note CI gate)

- [x] **B0-T2.1 Unit (success):** gate rule recognizes required migration section when breaking change marker exists.
- [x] **B0-T2.2 Unit (edge):** gate fails when marker exists and migration section is missing.
- [x] **B0-T2.3 Integration:** PR/release validation workflow blocks merge/release on missing migration section.
- [x] **B0-T2.4 CI/E2E:** workflow demonstrates expected failure message and remediation guidance.
- [x] **B0-T2.5 Evidence:** workflow config + local gate test evidence + live CI run URL captured.

Implementation evidence:

- [x] `.github/workflows/migration-note-gate.yml`
- [x] `scripts/release/check_migration_gate.py`
- [x] `scripts/release/check_migration_gate_test.py`
- [x] `docs/reference/releases.md` (automated enforcement section)

#### R-DEBT-03 — close DEF-004 (guided git-init flow)

- [x] **B0-T3.1 Unit (interactive success):** missing git root in interactive mode prompts to run `git init`.
- [x] **B0-T3.2 Unit (non-interactive edge):** missing git root in non-interactive mode fails deterministically with actionable message.
- [x] **B0-T3.3 Integration:** command flow respects user choice in interactive prompt.
- [x] **B0-T3.4 E2E:** non-interactive CLI behavior validated with actionable output expectations.
- [x] **B0-T3.5 Evidence:** unit + non-interactive and interactive targeted e2e evidence captured.

Implementation evidence:

- [x] `internal/cli/forge_handler.go`
- [x] `internal/cli/forge_handler_test.go`
- [x] `test/e2e/git_init_guidance_e2e_test.go`
- [x] `test/e2e/git_init_guidance_interactive_e2e_test.go`
- [x] `docs/reference/install.md`

### B0 Evidence Log (current)

- Local unit test evidence (release + migration scripts):
  - `python3 -m unittest scripts/release/release_artifacts_test.py scripts/release/check_migration_gate_test.py`
  - Result: `Ran 7 tests ... OK`
- Local unit test evidence (guided git-init flow):
  - `go test ./internal/cli -run 'TestProjectRootDetector_(FailsWhenNoGitAncestorExists|InteractivePromptDeclinedFails|InteractivePromptAcceptedRunsGitInit|InteractivePromptAcceptedGitInitFails)|TestDefaultIsInteractiveSession_(RespectsNonInteractiveEnvOverride|ForceInteractiveOverrideWins)'`
  - Result: `ok github.com/Jfgm299/weave-cli/internal/cli`
- Local e2e test definition added for non-interactive root-detection behavior:
  - `test/e2e/git_init_guidance_e2e_test.go::TestForge_E2E_NoGitRoot_NonInteractiveShowsActionableGuidance`
- Local e2e execution evidence (non-interactive root detection):
  - `go test -tags=e2e ./test/e2e -run TestForge_E2E_NoGitRoot_NonInteractiveShowsActionableGuidance`
  - Result: `ok github.com/Jfgm299/weave-cli/test/e2e`
- Local e2e execution evidence (interactive decline branch):
  - `go test -tags=e2e ./test/e2e -run TestForge_E2E_NoGitRoot_ForcedInteractiveDeclineShowsDeclinedGuidance`
  - Result: `ok github.com/Jfgm299/weave-cli/test/e2e`

### B0 CI Run URLs

- T1.4/T1.5 — `release-artifacts.yml` workflow_dispatch run (initial failure before key-policy fix): https://github.com/Jfgm299/weave-cli/actions/runs/24534428611
- T1.4/T1.5 — `release-artifacts.yml` workflow_dispatch run (success after policy fix): https://github.com/Jfgm299/weave-cli/actions/runs/24534779007
- T1.4/T1.5 — `install-artifact-validation.yml` push run (success): https://github.com/Jfgm299/weave-cli/actions/runs/24534774609
- T2.4/T2.5 — `migration-note-gate.yml` pull_request run (success): https://github.com/Jfgm299/weave-cli/actions/runs/24534770308
- T3.5 — interactive git-init e2e evidence (local command + output):
  - `go test -tags=e2e ./test/e2e -run TestForge_E2E_NoGitRoot_ForcedInteractiveDeclineShowsDeclinedGuidance`
  - Result: `ok github.com/Jfgm299/weave-cli/test/e2e`

---

## Batch 1 — Codex provider integration baseline

### B1 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-PROV-04 (codex provider parity) | B1-T1.1, B1-T1.2 | B1-T1.3 | B1-T1.4 | B1-T1.5 | [x] |
| R-PROV-05 (provider projection strategy contract) | B1-T2.1, B1-T2.2 | B1-T2.3 | B1-T2.4 | B1-T2.5 | [x] |

### B1 Tasks

#### R-PROV-04 — codex provider MUST have add/remove/repair/list parity

- [x] **B1-T1.1 Unit (success):** codex adapter is registered and exposes required binaries (`codex`).
- [x] **B1-T1.2 Unit (edge):** missing `codex` binary yields actionable error and blocks mutation.
- [x] **B1-T1.3 Integration:** add/remove/repair/list codex behaves consistently with existing providers.
- [x] **B1-T1.4 E2E:** codex provider flow validated from CLI in sandbox project.
- [x] **B1-T1.5 Evidence:** config and filesystem projections are deterministic and reversible.

#### R-PROV-05 — provider projection strategy MUST remain adapter-driven

- [x] **B1-T2.1 Unit (success):** provider projection ops are produced by registry adapters, not ad-hoc CLI logic.
- [x] **B1-T2.2 Unit (edge):** unsupported provider handling is adapter/registry-gated.
- [x] **B1-T2.3 Integration:** provider projection ops are executed transactionally via service layer.
- [x] **B1-T2.4 E2E:** projection links for codex are created/removed/repaired correctly.
- [x] **B1-T2.5 Evidence:** docs and tests confirm adapter contract and projection paths.

### B1 Evidence Log (current)

- Local unit test evidence (provider registry + service + CLI):
  - `go test ./internal/providers -run 'Test(DefaultRegistry_Get_ReturnsKnownProviders|CodexAdapter_RequiredBinaries_DeclaresCodexBinary|CodexAdapter_PlanSetup_CreatesProviderProjectionLinks)'`
  - `go test ./internal/app -run 'TestProviderService_'`
  - `go test ./internal/cli -run 'TestProviderRunAction_'`
  - Result: all pass (`ok`).
- Local integration test evidence (provider parity + binary failure no-mutation):
  - `go test ./test/integration -run 'TestProviderAdd_Integration_'`
  - Result: pass (`ok`).
- Local e2e evidence (add/remove/repair/list parity + actionable missing binary):
  - `go test -tags=e2e ./test/e2e -run 'TestProvider(AddAndList_E2E_MultiProviderFlowAndActionableErrors|Add_E2E_MissingBinaryActionableFailureAndNoConfigMutation|Add_E2E_CodexMissingBinaryActionableFailureAndNoConfigMutation|RemoveAndRepair_E2E_ReversibleOperations|Add_E2E_DryRunDoesNotMutateAndPrintsActionableSummary|RemoveAndRepair_E2E_DryRunNoMutation)'`
  - Result: pass (`ok`).

Blockers:

- None.

## Batch 2 — Provider-aware command add behavior

### B2 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-CMD-05 (default all enabled providers) | B2-T1.1, B2-T1.2 | B2-T1.3 | B2-T1.4 | B2-T1.5 | [x] |
| R-CMD-06 (`--provider` exclusive mode) | B2-T2.1, B2-T2.2 | B2-T2.3 | B2-T2.4 | B2-T2.5 | [x] |
| R-CMD-07 (codex namespace `__weave_commands__`) | B2-T3.1, B2-T3.2 | B2-T3.3 | B2-T3.4 | B2-T3.5 | [x] |
| R-CMD-08/09 (no-provider interactive/non-interactive behavior) | B2-T4.1, B2-T4.2 | B2-T4.3 | B2-T4.4 | B2-T4.5 | [x] |

### B2 Tasks

#### R-CMD-05 — default command install MUST target all enabled providers

- [x] **B2-T1.1 Unit (success):** planner resolves enabled provider targets deterministically.
- [x] **B2-T1.2 Unit (edge):** unsupported/invalid provider metadata fails with actionable error.
- [x] **B2-T1.3 Integration:** shared install writes canonical command plus provider projections.
- [x] **B2-T1.4 E2E:** `weave command add <name>` projects for all enabled providers.
- [x] **B2-T1.5 Evidence:** filesystem + config assertions confirm all-enabled projection behavior.

#### R-CMD-06 — `--provider` MUST perform exclusive provider install

- [x] **B2-T2.1 Unit (success):** parser and planner resolve exclusive provider mode.
- [x] **B2-T2.2 Unit (edge):** invalid/duplicate `--provider` flag usage fails deterministically.
- [x] **B2-T2.3 Integration:** exclusive mode skips shared `.agents` requirement.
- [x] **B2-T2.4 E2E:** `weave command add <name> --provider codex` succeeds without shared install.
- [x] **B2-T2.5 Evidence:** shared paths untouched in exclusive-mode assertions.

#### R-CMD-07 — codex command wrappers MUST use `__weave_commands__`

- [x] **B2-T3.1 Unit (success):** codex projection path builder emits `__weave_commands__/<name>/SKILL.md`.
- [x] **B2-T3.2 Unit (edge):** projection strategy rejects unsupported provider mapping.
- [x] **B2-T3.3 Integration:** codex projection writes namespace wrapper symlink correctly.
- [x] **B2-T3.4 E2E:** codex wrapper path exists and resolves to expected command source.
- [x] **B2-T3.5 Evidence:** path assertions and provider docs confirm namespace contract.

#### R-CMD-08/09 — no-provider behavior MUST be interactive-safe and non-interactive-safe

- [x] **B2-T4.1 Unit (success):** interactive path prompts in English with default `[y/N]`.
- [x] **B2-T4.2 Unit (edge):** non-interactive path fails with actionable guidance.
- [x] **B2-T4.3 Integration:** command add control flow honors confirmation callback and cancellation.
- [x] **B2-T4.4 E2E:** interactive default-no and non-interactive fail behaviors are observable.
- [x] **B2-T4.5 Evidence:** output, exit code, and config immutability assertions captured.

### B2 Evidence Log (current)

- Local unit test evidence (CLI parser/planning + handler projections + transactional rollback):
  - `go test ./internal/cli`
  - `go test ./internal/app -run 'TestForgeService_AddAsset_(WithAdditionalOperations_AppliesAllAndPersistsConfig|ConfigWriteFailureRollsBackAdditionalOperations)'`
  - Result: pass (`ok`).
- Local integration test evidence (shared + exclusive provider-aware command install):
  - `go test ./test/integration -run 'TestAddAsset_Integration_(CommandWithProviderProjections_WritesSharedAndProviderTargets|CommandExclusiveProvider_DoesNotRequireSharedPath)'`
  - Result: pass (`ok`).
- Local e2e evidence (default all-enabled, exclusive `--provider`, interactive/non-interactive no-provider behavior):
  - `go test -tags=e2e ./test/e2e -run 'TestCommandAdd_E2E_(DefaultAllEnabledProviders_ProjectsForEachProvider|ProviderExclusiveCodex_DoesNotRequireSharedAgentsPath|NoProvidersInteractivePrompt_DefaultNo|NoProvidersInteractivePrompt_EmptyInputDefaultsToNo|NoProvidersNonInteractiveFailsWithGuidance)'`
  - Result: pass (`ok`).

Blockers:

- None.

## Batch 3 — Transaction + doctor consistency

### B3 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| R-ARCH-08 (provider logic outside CLI parsing) | B3-T1.1, B3-T1.2 | B3-T1.3 | B3-T1.4 | B3-T1.5 | [x] |
| R-ARCH-09 (full rollback on multi-provider failure) | B3-T2.1, B3-T2.2 | B3-T2.3 | B3-T2.4 | B3-T2.5 | [x] |
| R-CONFIG-08 (transactional provider-target persistence) | B3-T3.1, B3-T3.2 | B3-T3.3 | B3-T3.4 | B3-T3.5 | [x] |

### B3 Tasks (live checklist)

#### R-ARCH-08 — provider projection logic MUST remain outside CLI parsing

- [x] **B3-T1.1 Unit (success):** command install planning composes provider operations through strategy/handler layer, not inline CLI branching.
- [x] **B3-T1.2 Unit (edge):** unsupported provider route fails from planning layer with actionable error, not from parser-level hardcoding.
- [x] **B3-T1.3 Integration:** end-to-end command planning uses provider strategy hooks for shared + exclusive modes.
- [x] **B3-T1.4 E2E:** provider-aware command add behavior remains consistent after parser changes/refactors.
- [x] **B3-T1.5 Evidence:** code-path assertions + command outputs show CLI only parses flags and delegates strategy.

### B3 Tasks

#### R-ARCH-09 — full rollback on multi-provider failures

- [x] **B3-T2.1 Unit (success):** rollback planner generates compensating operations for every additional provider projection op.
- [x] **B3-T2.2 Unit (edge):** failure at Nth provider projection rolls back prior applied ops and keeps system converged.
- [x] **B3-T2.3 Integration:** simulated provider projection failure leaves no partial provider directories/symlinks.
- [x] **B3-T2.4 E2E:** injected projection failure during default all-enabled `command add` returns non-zero and clean filesystem/config.
- [x] **B3-T2.5 Evidence:** `weave.yaml` unchanged + no orphan symlinks after failure scenario.

#### R-CONFIG-08 — config persistence MUST be transactional across provider-target installs

- [x] **B3-T3.1 Unit (success):** config writer executes only after successful application of shared + provider-target ops.
- [x] **B3-T3.2 Unit (edge):** config write failure triggers rollback of applied provider-target ops.
- [x] **B3-T3.3 Integration:** command add with multi-provider target persists config atomically when all ops succeed.
- [x] **B3-T3.4 E2E:** failure path confirms config remains unchanged while success path records deterministic inventory.
- [x] **B3-T3.5 Evidence:** before/after config snapshots + exit code assertions captured.

### B3 Evidence Log (current)

- Local unit test evidence (architecture delegation + projection strategy + rollback semantics):
  - `go test ./internal/app -run 'Test(CommandInstallPlanner_|BuildCommandProjectionOps_|ForgeService_AddAsset_ApplyFailureRollsBackAllPlannedOperations|ForgeService_AddAsset_ConfigWriteFailureRollsBackAdditionalOperations)'`
  - `go test ./internal/cli -run 'TestResolveCommandInstallPlan_|TestResolveCommandInstallPlanWithDeps_|TestBuildAddAssetInput_Command_'`
  - Result: pass (`ok`).
- Local unit test evidence (doctor shared+exclusive repair guidance):
  - `go test ./internal/app -run 'Test(DoctorService_Run_CommandSharedInstallMissingSymlinkHasCommandRepairGuidance|DoctorService_Run_CommandExclusiveProjectionMissingUsesProviderSpecificRepair|DoctorService_Run_UnknownEnabledProviderFlagsStaleIntegration)'`
  - `go test ./internal/cli -run 'TestRunDoctor_(CommandExclusiveProjectionMissingShowsProviderRepairGuidance|StaleProviderIntegrationIsReported)'`
  - Result: pass (`ok`).
- Local integration test evidence (transaction rollback + projection/doctor consistency):
  - `go test ./test/integration -run 'Test(AddAsset_Integration_ApplyFailureRollsBackSharedAndProviderTargetsAndKeepsConfigUnchanged|AddAsset_Integration_CommandWithProviderProjections_WritesSharedAndProviderTargets|Doctor_Integration_CommandExclusiveProjectionDriftReportsProviderScopedRepair)'`
  - Result: pass (`ok`).
- Local e2e evidence (all-enabled + exclusive command routing, rollback, transactional config, doctor guidance):
  - `go test -tags=e2e ./test/e2e -run 'Test(CommandAdd_E2E_MultiProviderFailureRollsBackFilesystemAndConfig|CommandAdd_E2E_MultiProviderSuccessPersistsTransactionalMetadata|Doctor_E2E_ExclusiveCommandProjectionMissingShowsProviderScopedRepair|CommandAdd_E2E_DefaultAllEnabledProviders_ProjectsForEachProvider|CommandAdd_E2E_ProviderExclusiveCodex_DoesNotRequireSharedAgentsPath)'`
  - Result: pass (`ok`).

Blockers:

- None.

## Batch 4 — Verification + docs closure

### B4 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E | Acceptance Evidence | Status |
|-------------|------|-------------|-----|---------------------|--------|
| v2 requirements traceability audit (R-PROV-04/05, R-CMD-05/06/07/08/09, R-ARCH-08/09, R-CONFIG-08) | B4-T1.1 | B4-T1.2 | B4-T1.3 | B4-T1.4 | [x] |
| Docs closure for v2 command/provider semantics | B4-T2.1 | B4-T2.2 | B4-T2.3 | B4-T2.4 | [x] |

### B4 Tasks

- [x] **B4-T1.1 Unit:** confirm all v2 requirements map to at least one success and one edge test.
- [x] **B4-T1.2 Integration:** run and verify integration matrix for B1–B3 requirements.
- [x] **B4-T1.3 E2E:** run and verify e2e matrix for B1–B3 requirements.
- [x] **B4-T1.4 Evidence:** publish final requirement→test traceability closure report.
- [x] **B4-T2.1 Unit/docs lint:** ensure docs examples align with implemented command flags/behavior.
- [x] **B4-T2.2 Integration/docs:** validate repair guidance/docs links from runtime errors and doctor output.
- [x] **B4-T2.3 E2E/docs:** verify documented provider/command flows in sandbox project.
- [x] **B4-T2.4 Evidence:** update `providers`, `doctor`, `transactions`, `config`, `install` docs with final confirmed behavior.

### B4 Evidence Log (current)

- Local unit traceability audit (success + edge paths for v2 requirements):
  - `go test ./internal/app -run 'TestProviderService_AddProvider_EnablesMultipleProviders|TestProviderService_AddProvider_CodexMissingBinaryReturnsActionableError|TestCommandInstallPlanner_Resolve_DefaultUsesEnabledProvidersFromStrategyLayer|TestCommandInstallPlanner_Resolve_UnsupportedProviderFailsInPlanningLayer|TestBuildCommandProjectionOps_DefaultSharedModeProjectsCodexOnly|TestBuildCommandProjectionOps_UnsupportedProviderFails|TestDoctorService_Run_CommandSharedInstallMissingSymlinkHasCommandRepairGuidance|TestDoctorService_Run_CommandExclusiveProjectionMissingUsesProviderSpecificRepair|TestForgeService_AddAsset_ApplyFailureRollsBackAllPlannedOperations|TestForgeService_AddAsset_ConfigWriteFailureRollsBackAdditionalOperations|TestDocsReferencePaths_ExistInRepository'`
  - `go test ./internal/cli -run 'TestBuildAddAssetInput_Command_DefaultAllEnabledProvidersIncludesSharedAndProjections|TestBuildAddAssetInput_Command_ExclusiveProviderSkipsSharedInstall|TestResolveCommandInstallPlan_NoProvidersNonInteractiveFails|TestRunDoctor_CommandExclusiveProjectionMissingShowsProviderRepairGuidance'`
  - `go test ./internal/providers -run 'TestDefaultRegistry_Get_ReturnsKnownProviders|TestCodexAdapter_RequiredBinaries_DeclaresCodexBinary|TestCodexAdapter_PlanSetup_CreatesProviderProjectionLinks'`
  - Result: pass (`ok`).
- Local integration matrix (B1–B3 coverage):
  - `go test ./test/integration -run 'TestProviderAdd_Integration_EnablesMultipleProvidersAndCreatesProjectionLinks|TestProviderAdd_Integration_CodexMissingBinaryKeepsConfigUnchanged|TestAddAsset_Integration_CommandWithProviderProjections_WritesSharedAndProviderTargets|TestAddAsset_Integration_CommandExclusiveProvider_DoesNotRequireSharedPath|TestAddAsset_Integration_ApplyFailureRollsBackSharedAndProviderTargetsAndKeepsConfigUnchanged|TestDoctor_Integration_CommandExclusiveProjectionDriftReportsProviderScopedRepair'`
  - Result: pass (`ok`).
- Local e2e matrix (provider-aware command flows, rollback, doctor repair guidance):
  - `go test -tags=e2e ./test/e2e -run 'TestProviderAddAndList_E2E_MultiProviderFlowAndActionableErrors|TestProviderAdd_E2E_CodexMissingBinaryActionableFailureAndNoConfigMutation|TestProviderRemoveAndRepair_E2E_ReversibleOperations|TestCommandAdd_E2E_DefaultAllEnabledProviders_ProjectsForEachProvider|TestCommandAdd_E2E_ProviderExclusiveCodex_DoesNotRequireSharedAgentsPath|TestCommandAdd_E2E_NoProvidersInteractivePrompt_DefaultNo|TestCommandAdd_E2E_NoProvidersInteractivePrompt_EmptyInputDefaultsToNo|TestCommandAdd_E2E_NoProvidersNonInteractiveFailsWithGuidance|TestCommandAdd_E2E_MultiProviderFailureRollsBackFilesystemAndConfig|TestCommandAdd_E2E_MultiProviderSuccessPersistsTransactionalMetadata|TestDoctor_E2E_ExclusiveCommandProjectionMissingShowsProviderScopedRepair'`
  - Result: pass (`ok`).
- Docs closure validation:
  - Runtime docs links/paths validated by tests: `TestDocsReferencePaths_ExistInRepository`, `TestRunDoctor_MissingAssetReturnsExitDoctorIssuesAndRepairPath`, `TestRunDoctor_CommandExclusiveProjectionMissingShowsProviderRepairGuidance`, `TestDoctor_E2E_ExclusiveCommandProjectionMissingShowsProviderScopedRepair`.
  - `docs/reference/install.md` updated to include `codex` in provider-binary troubleshooting list.

### B4 Final v2 Requirement → Test Traceability (closure)

| Requirement | Success Coverage | Edge/Failure Coverage |
|-------------|------------------|-----------------------|
| R-PROV-04 | `TestProviderService_AddProvider_EnablesMultipleProviders`, `TestProviderAdd_Integration_EnablesMultipleProvidersAndCreatesProjectionLinks`, `TestProviderAddAndList_E2E_MultiProviderFlowAndActionableErrors` | `TestProviderService_AddProvider_CodexMissingBinaryReturnsActionableError`, `TestProviderAdd_Integration_CodexMissingBinaryKeepsConfigUnchanged`, `TestProviderAdd_E2E_CodexMissingBinaryActionableFailureAndNoConfigMutation` |
| R-PROV-05 | `TestDefaultRegistry_Get_ReturnsKnownProviders`, `TestCodexAdapter_PlanSetup_CreatesProviderProjectionLinks`, `TestAddAsset_Integration_CommandWithProviderProjections_WritesSharedAndProviderTargets` | `TestCommandInstallPlanner_Resolve_UnsupportedProviderFailsInPlanningLayer`, `TestBuildCommandProjectionOps_UnsupportedProviderFails` |
| R-CMD-05 | `TestBuildAddAssetInput_Command_DefaultAllEnabledProvidersIncludesSharedAndProjections`, `TestCommandAdd_E2E_DefaultAllEnabledProviders_ProjectsForEachProvider` | `TestCommandAdd_E2E_MultiProviderFailureRollsBackFilesystemAndConfig` |
| R-CMD-06 | `TestBuildAddAssetInput_Command_ExclusiveProviderSkipsSharedInstall`, `TestAddAsset_Integration_CommandExclusiveProvider_DoesNotRequireSharedPath`, `TestCommandAdd_E2E_ProviderExclusiveCodex_DoesNotRequireSharedAgentsPath` | `TestResolveCommandInstallPlan_NoProvidersNonInteractiveFails` |
| R-CMD-07 | `TestBuildCommandProjectionOps_DefaultSharedModeProjectsCodexOnly`, `TestCommandAdd_E2E_DefaultAllEnabledProviders_ProjectsForEachProvider` | `TestBuildCommandProjectionOps_UnsupportedProviderFails` |
| R-CMD-08 | `TestCommandAdd_E2E_NoProvidersInteractivePrompt_DefaultNo` (prompt contract visible) | `TestCommandAdd_E2E_NoProvidersInteractivePrompt_EmptyInputDefaultsToNo` (default-no on empty input) |
| R-CMD-09 | `TestRunDoctor_CommandExclusiveProjectionMissingShowsProviderRepairGuidance` (actionable guidance pattern) | `TestResolveCommandInstallPlan_NoProvidersNonInteractiveFails`, `TestCommandAdd_E2E_NoProvidersNonInteractiveFailsWithGuidance` (non-interactive fail + no mutation) |
| R-ARCH-08 | `TestCommandInstallPlanner_Resolve_DefaultUsesEnabledProvidersFromStrategyLayer` | `TestCommandInstallPlanner_Resolve_UnsupportedProviderFailsInPlanningLayer` |
| R-ARCH-09 | `TestForgeService_AddAsset_ApplyFailureRollsBackAllPlannedOperations`, `TestAddAsset_Integration_ApplyFailureRollsBackSharedAndProviderTargetsAndKeepsConfigUnchanged` | `TestCommandAdd_E2E_MultiProviderFailureRollsBackFilesystemAndConfig` |
| R-CONFIG-08 | `TestForgeService_AddAsset_ConfigWriteFailureRollsBackAdditionalOperations`, `TestCommandAdd_E2E_MultiProviderSuccessPersistsTransactionalMetadata` | `TestAddAsset_Integration_ApplyFailureRollsBackSharedAndProviderTargetsAndKeepsConfigUnchanged`, `TestCommandAdd_E2E_MultiProviderFailureRollsBackFilesystemAndConfig` |

Blockers:

- None.
