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
| B0 | P0 | R-DEBT-01 | Close DEF-001/002 CI automation debt | [-] |
| B0 | P0 | R-DEBT-02 | Close DEF-003 migration-note CI gate debt | [-] |
| B0 | P0 | R-DEBT-03 | Close DEF-004 guided git-init flow debt | [-] |
| B1 | P0 | R-PROV-04 | Add codex provider parity | [ ] |
| B1 | P0 | R-PROV-05 | Provider command projection strategy contract | [ ] |
| B2 | P0 | R-CMD-05 | `command add` default installs for all enabled providers | [ ] |
| B2 | P0 | R-CMD-06 | `--provider` exclusive install without shared `.agents` requirement | [ ] |
| B2 | P0 | R-CMD-07 | Codex namespace `__weave_commands__` | [ ] |
| B2 | P0 | R-CMD-08 | Interactive no-provider prompt `[y/N]` | [ ] |
| B2 | P0 | R-CMD-09 | Non-interactive no-provider auto-fail | [ ] |
| B3 | P0 | R-ARCH-08 | Keep provider logic out of CLI handlers | [ ] |
| B3 | P0 | R-ARCH-09 | Full rollback on multi-provider failure | [ ] |
| B3 | P0 | R-CONFIG-08 | Transactional provider-target persistence | [ ] |

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
  - `go test ./internal/cli -run 'TestProjectRootDetector_(FailsWhenNoGitAncestorExists|InteractivePromptDeclinedFails|InteractivePromptAcceptedRunsGitInit|InteractivePromptAcceptedGitInitFails)'`
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

- T1.4/T1.5 — `release-artifacts.yml` workflow_dispatch run (failure path before key-policy fix): https://github.com/Jfgm299/weave-cli/actions/runs/24534428611
- T1.4/T1.5 — `install-artifact-validation.yml` push run (success): https://github.com/Jfgm299/weave-cli/actions/runs/24534393193
- T2.4/T2.5 — `migration-note-gate.yml` pull_request run (success): https://github.com/Jfgm299/weave-cli/actions/runs/24534384197

### B0 CI Run URLs (fill to close Batch 0)

- T1.4/T1.5 — `release-artifacts.yml` run URL: _pending_
- T1.4/T1.5 — `install-artifact-validation.yml` run URL: _pending_
- T2.4/T2.5 — `migration-note-gate.yml` run URL: _pending_
- T3.5 — interactive git-init e2e evidence (command output/log link): _pending_

---

## Batch 1 — Codex provider integration baseline

### B1 Tasks

- [ ] Implement codex adapter registration and required binary checks.
- [ ] Implement provider add/remove/repair/list parity tests for codex.
- [ ] Update provider docs with codex support and constraints.

## Batch 2 — Provider-aware command add behavior

### B2 Tasks

- [ ] Implement default route for all enabled providers.
- [ ] Implement `--provider` exclusive route.
- [ ] Implement codex wrapper namespace `__weave_commands__`.
- [ ] Implement no-provider interactive prompt and non-interactive fail behavior.

## Batch 3 — Transaction + doctor consistency

### B3 Tasks

- [ ] Extend doctor checks for shared + exclusive command installs.
- [ ] Ensure full rollback for partial multi-provider failures.
- [ ] Align repair guidance with provider-aware install modes.

## Batch 4 — Verification + docs closure

### B4 Tasks

- [ ] Final requirement→test traceability audit.
- [ ] Unit/integration/e2e full pass for v2 delta.
- [ ] Docs update: `providers`, `doctor`, `transactions`, `config`, `install`.
