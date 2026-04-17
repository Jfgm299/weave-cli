# Tasks — weave-v3-community-readiness

## Scope

Implement v3 planning delta with strict v1-style batch discipline:

1. Close v2 deferred debt (`DEF-V2-001..003`) in Batch 0.
2. Implement CI evidence automation Option A (deterministic SHA manifest + PR reference).
3. Enforce release/versioning governance and semver policy.
4. Enforce PR checklist-label-issue coherence quality gates.
5. Add catalog-internet MVP and consistency hardening with deterministic snapshot semantics.

## Status Legend

- [ ] Pending
- [-] In progress
- [x] Done (100% implemented, tested, and operational)
- [!] Blocked

## Completion Policy (mandatory)

- Never mark `[x]` unless code + tests + acceptance evidence are complete.
- Deferred or baseline-only closure MUST stay `[-]` and be tracked in `deferred.md`.
- For this change, Batch 0 debt closure is mandatory before all other batches.

---

## Master Backlog

| Batch | Priority | Requirement | Title | Status |
|------:|:--------:|-------------|-------|:------:|
| B0 | P0 | R-DEBT-V3-01 | Close DEF-V2-001 (`__pycache__` tracked artifacts) | [-] |
| B0 | P0 | R-DEBT-V3-02 | Close DEF-V2-002 (manual CI evidence capture) | [-] |
| B0 | P0 | R-DEBT-V3-03 | Close DEF-V2-003 (PR checklist-label coherence drift) | [-] |
| B1 | P0 | R-CI-01 | Per-workflow JSON evidence emission contract | [ ] |
| B1 | P0 | R-CI-02 | `workflow_run` SHA aggregator for deterministic consolidation | [ ] |
| B1 | P0 | R-CI-03 | Deterministic manifest reference + PR visibility | [ ] |
| B1 | P0 | R-CI-04 | OpenSpec evidence flow defaults to machine-generated references | [ ] |
| B2 | P0 | R-VER-01 | Semver bump governance policy for weave binary | [ ] |
| B2 | P0 | R-VER-02 | Enforce synchronized updates (`internal/cli/version.go` + version tests) | [ ] |
| B2 | P0 | R-VER-03 | CI/release gate for version governance compliance | [ ] |
| B3 | P0 | R-PRQ-01 | PR checklist-label coherence quality gate | [ ] |
| B3 | P0 | R-PRQ-02 | PR issue-link semantics quality gate | [ ] |
| B4 | P0 | R-CAT-01 | Cross-source consistent listing contract (same query + same snapshot) | [ ] |
| B4 | P0 | R-CAT-02 | Multi-source provider model (`github_curated_index`, `registry_json`) | [ ] |
| B4 | P0 | R-CAT-03 | Deterministic normalization schema for catalog entries | [ ] |
| B4 | P0 | R-CAT-04 | Stable identity + deterministic dedup/merge rules | [ ] |
| B4 | P0 | R-CAT-05 | Deterministic ranking/sort for repeatable listing order | [ ] |
| B4 | P0 | R-CAT-06 | Freshness/staleness metadata + conflict-resolution policy | [ ] |
| B4 | P0 | R-CAT-07 | Snapshot-based sync with version/hash for repeatability | [ ] |
| B4 | P0 | R-CAT-08 | Offline search against committed synced snapshot | [ ] |
| B4 | P0 | R-CAT-09 | Trust policy enforcement (allowlist + pin by ref/hash) | [ ] |
| B4 | P0 | R-CAT-10 | Malformed/partial source handling (quarantine/skip + diagnostics) | [ ] |

---

## Batch Plan Overview

- **Batch 0 (B0) — Debt closure first (P0, mandatory gate)**
- **Batch 1 (B1) — CI evidence automation Option A (P0)**
- **Batch 2 (B2) — Release/versioning governance enforcement (P0)**
- **Batch 3 (B3) — PR quality coherence gates (P0)**
- **Batch 4 (B4) — Catalog internet MVP (P0)**
- **Batch 5 (B5) — Catalog consistency hardening + final verification and traceability closure (P0/P1)**

---

## Batch 0 — Debt closure first (mandatory)

### B0 Exit Criteria

- [-] DEF-V2-001 closure criteria satisfied and guarded against regression.
- [-] DEF-V2-002 closure criteria satisfied with Option A baseline active.
- [-] DEF-V2-003 closure criteria satisfied with enforced quality gate.
- [-] v2 deferred register updated with closure evidence references.

### B0 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E/CI | Acceptance Evidence | Status |
|-------------|------|-------------|--------|---------------------|--------|
| R-DEBT-V3-01 (DEF-V2-001) | B0-T1.1, B0-T1.2 | B0-T1.3 | B0-T1.4 | B0-T1.5 | [-] |
| R-DEBT-V3-02 (DEF-V2-002) | B0-T2.1, B0-T2.2 | B0-T2.3 | B0-T2.4 | B0-T2.5 | [-] |
| R-DEBT-V3-03 (DEF-V2-003) | B0-T3.1, B0-T3.2 | B0-T3.3 | B0-T3.4 | B0-T3.5 | [-] |

### B0 Tasks (live checklist)

#### R-DEBT-V3-01 — close DEF-V2-001 (`__pycache__` tracked artifacts)

- [x] **B0-T1.1 Unit (success):** repository hygiene rule detects Python cache artifacts deterministically.
- [x] **B0-T1.2 Unit (edge):** hygiene rule fails when tracked `__pycache__` artifacts are reintroduced.
- [x] **B0-T1.3 Integration:** CI hygiene check blocks PR when tracked cache artifacts exist.
- [-] **B0-T1.4 E2E/CI:** representative run shows pass and fail behavior with actionable diagnostics. (local pass/fail complete; live GitHub Actions run evidence pending)
- [-] **B0-T1.5 Evidence:** workflow/gate definition + deterministic evidence manifest reference + deferred closure note. (live run URL pending)

#### R-DEBT-V3-02 — close DEF-V2-002 (manual CI evidence capture)

- [x] **B0-T2.1 Unit (success):** evidence producer emits valid per-run JSON payload.
- [x] **B0-T2.2 Unit (edge):** malformed/missing evidence payload is rejected by aggregator validation.
- [x] **B0-T2.3 Integration:** aggregator consolidates required workflow artifacts by SHA.
- [-] **B0-T2.4 E2E/CI:** PR flow exposes deterministic manifest reference without manual URL copy. (baseline generated locally; live PR automation run pending)
- [-] **B0-T2.5 Evidence:** `openspec/evidence/<sha>.json` reference + PR comment reference + closure update in deferred register. (live PR comment evidence pending)

#### R-DEBT-V3-03 — close DEF-V2-003 (PR checklist-label coherence)

- [x] **B0-T3.1 Unit (success):** coherence parser validates compliant checklist/label mappings.
- [x] **B0-T3.2 Unit (edge):** parser rejects incoherent checklist/label combinations and invalid issue-link semantics.
- [x] **B0-T3.3 Integration:** CI quality gate blocks incoherent PR metadata.
- [-] **B0-T3.4 E2E/CI:** representative PR validation run demonstrates deterministic pass/fail results. (script-level pass/fail covered; live PR-triggered run pending)
- [-] **B0-T3.5 Evidence:** gate definition + test outputs + manifest reference + deferred closure update. (live run URL pending)

### B0 Evidence Log (current)

- Local unit tests (R-DEBT-V3-01/02/03):
  - `python3 -m unittest scripts/ci/check_repo_hygiene_test.py scripts/ci/collect_workflow_evidence_test.py scripts/ci/evidence_manifest_test.py scripts/ci/check_pr_metadata_coherence_test.py`
  - Outcome: `Ran 9 tests ... OK`.
- Local integration/e2e representative (R-DEBT-V3-01):
  - `python3 scripts/ci/check_repo_hygiene.py` (before untracking pycache)
  - Outcome: failed with actionable diagnostics listing tracked `scripts/release/__pycache__/*.pyc`.
  - `git rm --cached scripts/release/__pycache__/check_migration_gate_test.cpython-314.pyc scripts/release/__pycache__/release_artifacts_test.cpython-314.pyc`
  - `python3 scripts/ci/check_repo_hygiene.py` (after untracking)
  - Outcome: `Repository hygiene passed: no tracked Python cache artifacts found`.
- Local integration/e2e representative (R-DEBT-V3-02 Option A baseline):
  - `python3 scripts/ci/run_evidence_baseline.py`
  - Outcome: deterministic manifest generated at `openspec/evidence/ef37b1acf4fdc5e8594f122a67a502ba3e9c5a6c.json`.
  - PR-facing reference generated at `openspec/evidence/pr-comment-ef37b1acf4fdc5e8594f122a67a502ba3e9c5a6c.md`.
- CI gate definitions added for Batch 0 closure:
  - `.github/workflows/repo-hygiene-gate.yml`
  - `.github/workflows/pr-metadata-coherence.yml`
  - Existing workflows now emit per-run evidence payloads using `scripts/ci/collect_workflow_evidence.py`.
- Deferred for full closure confirmation:
  - `openspec/changes/weave-v3-community-readiness/deferred.md#DEF-V3-B0-001`
  - `openspec/changes/weave-v3-community-readiness/deferred.md#DEF-V3-B0-002`

### B0 Blockers (current)

- None.

---

## Batch 1 — CI evidence automation Option A

### B1 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E/CI | Acceptance Evidence | Status |
|-------------|------|-------------|--------|---------------------|--------|
| R-CI-01 (per-workflow JSON evidence) | B1-T1.1, B1-T1.2 | B1-T1.3 | B1-T1.4 | B1-T1.5 | [ ] |
| R-CI-02 (`workflow_run` SHA aggregator) | B1-T2.1, B1-T2.2 | B1-T2.3 | B1-T2.4 | B1-T2.5 | [ ] |
| R-CI-03 (deterministic manifest + PR visibility) | B1-T3.1, B1-T3.2 | B1-T3.3 | B1-T3.4 | B1-T3.5 | [ ] |
| R-CI-04 (OpenSpec evidence default automation) | B1-T4.1, B1-T4.2 | B1-T4.3 | B1-T4.4 | B1-T4.5 | [ ] |

### B1 Tasks (live checklist)

#### R-CI-01 — per-workflow JSON evidence emission

- [ ] **B1-T1.1 Unit (success):** evidence schema validator accepts required fields and deterministic keys.
- [ ] **B1-T1.2 Unit (edge):** missing or inconsistent fields fail schema validation.
- [ ] **B1-T1.3 Integration:** required workflows emit artifacts conforming to schema.
- [ ] **B1-T1.4 E2E/CI:** complete run set produces evidence artifacts for all required workflows.
- [ ] **B1-T1.5 Evidence:** artifact list and schema validation output captured in manifest.

#### R-CI-02 — `workflow_run` aggregator consolidation by head SHA

- [ ] **B1-T2.1 Unit (success):** aggregator groups evidence by exact head SHA deterministically.
- [ ] **B1-T2.2 Unit (edge):** partial/missing workflow evidence is represented explicitly as missing/failure state.
- [ ] **B1-T2.3 Integration:** aggregator pipeline ingests per-run artifacts and emits one manifest per SHA.
- [ ] **B1-T2.4 E2E/CI:** reruns for same SHA preserve deterministic consolidated output contract.
- [ ] **B1-T2.5 Evidence:** consolidated manifest payload contains stable ordered workflow entries.

#### R-CI-03 — deterministic manifest reference + PR visibility

- [ ] **B1-T3.1 Unit (success):** manifest reference path generation is deterministic for SHA input.
- [ ] **B1-T3.2 Unit (edge):** invalid SHA input fails with actionable diagnostics.
- [ ] **B1-T3.3 Integration:** PR comment payload references generated manifest deterministically.
- [ ] **B1-T3.4 E2E/CI:** PR flow displays manifest reference after required workflows complete.
- [ ] **B1-T3.5 Evidence:** PR comment snapshot and artifact reference align with same SHA manifest.

#### R-CI-04 — OpenSpec evidence defaults to machine-generated references

- [ ] **B1-T4.1 Unit (success):** evidence extraction utility resolves manifest references for tasks without manual URL input.
- [ ] **B1-T4.2 Unit (edge):** missing manifest reference triggers deterministic “evidence unavailable” state.
- [ ] **B1-T4.3 Integration:** tasks evidence logging consumes manifest references directly.
- [ ] **B1-T4.4 E2E/CI:** end-to-end batch evidence closure uses automated references only.
- [ ] **B1-T4.5 Evidence:** sample task closure demonstrates no manual run URL copy requirement.

### B1 Evidence Log (current)

- Pending — no execution evidence captured yet.

### B1 Blockers (current)

- None.

---

## Batch 2 — Release/versioning governance enforcement

### B2 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E/CI | Acceptance Evidence | Status |
|-------------|------|-------------|--------|---------------------|--------|
| R-VER-01 (semver bump governance) | B2-T1.1, B2-T1.2 | B2-T1.3 | B2-T1.4 | B2-T1.5 | [ ] |
| R-VER-02 (sync `version.go` + version tests) | B2-T2.1, B2-T2.2 | B2-T2.3 | B2-T2.4 | B2-T2.5 | [ ] |
| R-VER-03 (CI/release governance gate) | B2-T3.1, B2-T3.2 | B2-T3.3 | B2-T3.4 | B2-T3.5 | [ ] |

### B2 Tasks (live checklist)

#### R-VER-01 — semver bump governance policy

- [ ] **B2-T1.1 Unit (success):** policy evaluator classifies patch/minor/major scenarios deterministically.
- [ ] **B2-T1.2 Unit (edge):** ambiguous or missing bump rationale fails validation.
- [ ] **B2-T1.3 Integration:** release metadata ingestion enforces bump policy mapping.
- [ ] **B2-T1.4 E2E/CI:** representative release candidate checks pass/fail according to semver rules.
- [ ] **B2-T1.5 Evidence:** policy matrix output is present in deterministic manifest.

#### R-VER-02 — synchronized `internal/cli/version.go` and version tests

- [ ] **B2-T2.1 Unit (success):** validator passes when version source and version tests are updated coherently.
- [ ] **B2-T2.2 Unit (edge):** validator fails when one of source/tests changes without counterpart update.
- [ ] **B2-T2.3 Integration:** CI check inspects diff and validates synchronized version updates.
- [ ] **B2-T2.4 E2E/CI:** unsynchronized version update PR is blocked deterministically.
- [ ] **B2-T2.5 Evidence:** CI output and manifest record include explicit version-sync result.

#### R-VER-03 — CI/release governance gate

- [ ] **B2-T3.1 Unit (success):** governance gate accepts fully compliant release candidate inputs.
- [ ] **B2-T3.2 Unit (edge):** governance gate rejects inconsistent semver/version/test metadata combinations.
- [ ] **B2-T3.3 Integration:** governance gate is wired into release validation workflow.
- [ ] **B2-T3.4 E2E/CI:** release workflow demonstrates deterministic governance gate outcomes.
- [ ] **B2-T3.5 Evidence:** gate result is included in SHA manifest and PR/release summary references.

### B2 Evidence Log (current)

- Pending — no execution evidence captured yet.

### B2 Blockers (current)

- None.

---

## Batch 3 — PR quality coherence gates

### B3 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E/CI | Acceptance Evidence | Status |
|-------------|------|-------------|--------|---------------------|--------|
| R-PRQ-01 (checklist-label coherence) | B3-T1.1, B3-T1.2 | B3-T1.3 | B3-T1.4 | B3-T1.5 | [ ] |
| R-PRQ-02 (issue-link semantics) | B3-T2.1, B3-T2.2 | B3-T2.3 | B3-T2.4 | B3-T2.5 | [ ] |

### B3 Tasks (live checklist)

#### R-PRQ-01 — checklist-label coherence gate

- [ ] **B3-T1.1 Unit (success):** parser validates expected checklist/label coherence mapping.
- [ ] **B3-T1.2 Unit (edge):** parser rejects mismatch cases with actionable diagnostics.
- [ ] **B3-T1.3 Integration:** PR validation workflow blocks incoherent metadata combinations.
- [ ] **B3-T1.4 E2E/CI:** representative compliant/non-compliant PR payloads produce deterministic pass/fail.
- [ ] **B3-T1.5 Evidence:** gate results and diagnostics captured in SHA manifest.

#### R-PRQ-02 — issue-link semantics gate (`Closes #<id>` vs `N/A`)

- [ ] **B3-T2.1 Unit (success):** validator accepts policy-compliant issue-link declarations.
- [ ] **B3-T2.2 Unit (edge):** validator rejects malformed/missing issue-link declarations.
- [ ] **B3-T2.3 Integration:** gate is required in PR checks and blocks merge on failure.
- [ ] **B3-T2.4 E2E/CI:** deterministic results demonstrated for compliant and violating payloads.
- [ ] **B3-T2.5 Evidence:** manifest includes issue-link semantics gate outcome and error details.

### B3 Evidence Log (current)

- Pending — no execution evidence captured yet.

### B3 Blockers (current)

- None.

---

## Batch 4 — Catalog internet MVP

### B4 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E/CI | Acceptance Evidence | Status |
|-------------|------|-------------|--------|---------------------|--------|
| R-CAT-02 (multi-source providers) | B4-T1.1, B4-T1.2 | B4-T1.3 | B4-T1.4 | B4-T1.5 | [ ] |
| R-CAT-03 (deterministic normalization schema) | B4-T2.1, B4-T2.2 | B4-T2.3 | B4-T2.4 | B4-T2.5 | [ ] |
| R-CAT-04 (stable identity + dedup) | B4-T3.1, B4-T3.2 | B4-T3.3 | B4-T3.4 | B4-T3.5 | [ ] |
| R-CAT-07 (snapshot version/hash repeatability) | B4-T4.1, B4-T4.2 | B4-T4.3 | B4-T4.4 | B4-T4.5 | [ ] |
| R-CAT-08 (offline snapshot search) | B4-T5.1, B4-T5.2 | B4-T5.3 | B4-T5.4 | B4-T5.5 | [ ] |
| R-CAT-09 (trust policy allowlist/pin) | B4-T6.1, B4-T6.2 | B4-T6.3 | B4-T6.4 | B4-T6.5 | [ ] |
| R-CAT-10 (malformed/partial handling) | B4-T7.1, B4-T7.2 | B4-T7.3 | B4-T7.4 | B4-T7.5 | [ ] |

### B4 Tasks (live checklist)

#### R-CAT-02 — multi-source provider model (`github_curated_index`, `registry_json`)

- [ ] **B4-T1.1 Unit (success):** provider registry resolves both initial provider types deterministically.
- [ ] **B4-T1.2 Unit (edge):** unsupported provider type fails validation with actionable diagnostics.
- [ ] **B4-T1.3 Integration:** ingestion pipeline consumes both provider payload contracts in one sync transaction.
- [ ] **B4-T1.4 E2E/CI:** representative sync run loads both providers and reports unified summary.
- [ ] **B4-T1.5 Evidence:** provider-resolution outputs included in deterministic snapshot evidence.

#### R-CAT-03 — deterministic normalization schema

- [ ] **B4-T2.1 Unit (success):** normalization maps provider records to required canonical fields.
- [ ] **B4-T2.2 Unit (edge):** missing required canonical fields are flagged deterministically.
- [ ] **B4-T2.3 Integration:** mixed-provider payloads produce identical canonical shape under repeated runs.
- [ ] **B4-T2.4 E2E/CI:** sync output demonstrates schema stability and canonical field completeness.
- [ ] **B4-T2.5 Evidence:** canonical schema validation report captured in snapshot evidence bundle.

#### R-CAT-04 — stable identity + deterministic dedup

- [ ] **B4-T3.1 Unit (success):** identity-key derivation yields stable key for equivalent records.
- [ ] **B4-T3.2 Unit (edge):** conflicting records follow deterministic merge tie-break chain.
- [ ] **B4-T3.3 Integration:** duplicate entries across providers collapse into one canonical listing row.
- [ ] **B4-T3.4 E2E/CI:** same logical skill/command appears once in catalog listing after sync.
- [ ] **B4-T3.5 Evidence:** dedup/merge decision log persisted with deterministic ordering.

#### R-CAT-07 — snapshot-based sync (version/hash)

- [ ] **B4-T4.1 Unit (success):** snapshot hash generation is deterministic for equivalent canonical datasets.
- [ ] **B4-T4.2 Unit (edge):** tampered/incomplete snapshot metadata is rejected.
- [ ] **B4-T4.3 Integration:** sync commits atomic snapshot with explicit version/hash metadata.
- [ ] **B4-T4.4 E2E/CI:** repeated sync with unchanged sources preserves snapshot hash and output contract.
- [ ] **B4-T4.5 Evidence:** snapshot metadata report includes version/hash and canonical entry count.

#### R-CAT-08 — offline search against synced snapshot

- [ ] **B4-T5.1 Unit (success):** search executor reads active snapshot without network dependency.
- [ ] **B4-T5.2 Unit (edge):** missing active snapshot returns deterministic “snapshot unavailable” diagnostics.
- [ ] **B4-T5.3 Integration:** offline query execution returns canonical entries and snapshot metadata.
- [ ] **B4-T5.4 E2E/CI:** online/offline search parity verified for same snapshot hash.
- [ ] **B4-T5.5 Evidence:** parity assertion output captured with query + snapshot hash references.

#### R-CAT-09 — trust policy allowlist + pin by ref/hash

- [ ] **B4-T6.1 Unit (success):** allowlisted + pinned sources pass trust validation.
- [ ] **B4-T6.2 Unit (edge):** non-allowlisted or unpinned sources fail trust checks deterministically.
- [ ] **B4-T6.3 Integration:** trust gate executes before normalization and blocks unsafe ingestion.
- [ ] **B4-T6.4 E2E/CI:** sync run demonstrates deterministic failure on trust policy violations.
- [ ] **B4-T6.5 Evidence:** trust validation outcomes and offending source diagnostics captured in evidence.

#### R-CAT-10 — malformed/partial source handling

- [ ] **B4-T7.1 Unit (success):** valid records proceed while malformed records are skipped/quarantined.
- [ ] **B4-T7.2 Unit (edge):** malformed payload does not abort snapshot commit when valid set remains.
- [ ] **B4-T7.3 Integration:** diagnostics channel records quarantine/skip reasons per source record.
- [ ] **B4-T7.4 E2E/CI:** snapshot remains searchable and uncorrupted after malformed-source ingestion attempt.
- [ ] **B4-T7.5 Evidence:** quarantine/skip diagnostics report linked to snapshot evidence.

### B4 Evidence Log (current)

- Pending — no execution evidence captured yet.

### B4 Blockers (current)

- None.

---

## Batch 5 — Catalog consistency hardening + final verification and traceability closure

### B5 Requirement Traceability Matrix

| Requirement | Unit | Integration | E2E/CI | Acceptance Evidence | Status |
|-------------|------|-------------|--------|---------------------|--------|
| R-CAT-01 (cross-source consistency guarantee) | B5-T1.1, B5-T1.2 | B5-T1.3 | B5-T1.4 | B5-T1.5 | [ ] |
| R-CAT-05 (deterministic ranking/sort) | B5-T2.1, B5-T2.2 | B5-T2.3 | B5-T2.4 | B5-T2.5 | [ ] |
| R-CAT-06 (freshness/staleness + conflict resolution) | B5-T3.1, B5-T3.2 | B5-T3.3 | B5-T3.4 | B5-T3.5 | [ ] |
| v3 traceability closure audit (R-DEBT-V3-01/02/03, R-CI-01/02/03/04, R-VER-01/02/03, R-PRQ-01/02, R-CAT-01..10) | B5-T4.1 | B5-T4.2 | B5-T4.3 | B5-T4.4 | [ ] |

### B5 Tasks (live checklist)

#### R-CAT-01 — cross-source consistent listing contract

- [ ] **B5-T1.1 Unit (success):** same-query evaluator over fixed snapshot yields stable canonical set from mixed-source corpus.
- [ ] **B5-T1.2 Unit (edge):** source-order permutation does not change listing output for same snapshot hash.
- [ ] **B5-T1.3 Integration:** end-to-end query path validates consistent listing over all initial sources.
- [ ] **B5-T1.4 E2E/CI:** replay suite confirms same query + same snapshot => byte-identical listing payload.
- [ ] **B5-T1.5 Evidence:** consistency proof artifacts include query input, snapshot hash, and deterministic output digest.

#### R-CAT-05 — deterministic ranking/sorting

- [ ] **B5-T2.1 Unit (success):** ranking function is deterministic for identical input/query/policy.
- [ ] **B5-T2.2 Unit (edge):** score ties resolve deterministically by canonical tie-break keys.
- [ ] **B5-T2.3 Integration:** ranking outputs remain stable across repeated catalog query executions.
- [ ] **B5-T2.4 E2E/CI:** deterministic ordering validated in CI replay job for representative queries.
- [ ] **B5-T2.5 Evidence:** sorted output checksum comparisons recorded in evidence manifest.

#### R-CAT-06 — freshness/staleness metadata + conflict policy

- [ ] **B5-T3.1 Unit (success):** freshness and staleness classification computed deterministically.
- [ ] **B5-T3.2 Unit (edge):** conflicting multi-source records resolve according to policy order or quarantine deterministically.
- [ ] **B5-T3.3 Integration:** listing output includes freshness/conflict metadata per canonical entry.
- [ ] **B5-T3.4 E2E/CI:** conflict scenarios produce stable winner/quarantine outcomes across reruns.
- [ ] **B5-T3.5 Evidence:** conflict-resolution decision table and staleness metadata snapshots captured.

#### Final traceability closure (all v3 requirements)

- [ ] **B5-T4.1 Unit:** verify each v3 requirement has success + edge test linkage.
- [ ] **B5-T4.2 Integration:** verify integration matrix completeness across B0–B5 requirements.
- [ ] **B5-T4.3 E2E/CI:** verify deterministic evidence references exist for all required checks.
- [ ] **B5-T4.4 Evidence:** publish final v3 requirement→test→evidence closure report.

### B5 Evidence Log (current)

- Pending — no execution evidence captured yet.

### B5 Blockers (current)

- None.

---

## Deferred / Debt Policy (explicit)

- If any R-CAT requirement cannot be guaranteed operationally in v3, it MUST remain `[-]` and be registered in `deferred.md` with:
  - deferred ID,
  - rationale,
  - owner,
  - concrete closure criteria,
  - target milestone.
- No catalog consistency requirement may be marked `[x]` based on policy/docs only; executable evidence is mandatory.
