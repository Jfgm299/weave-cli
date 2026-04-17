# Spec — weave-v3-community-readiness

## Scope

This spec defines the v3 delta for:

1. CI evidence automation Option A with deterministic manifesting by head SHA.
2. Debt-first closure of v2 deferred items (`DEF-V2-001..003`) in Batch 0.
3. Release/versioning governance for semver and source-of-truth version files/tests.
4. PR quality gate coherence checks between checklist state, labels, and issue-link semantics.
5. Catalog-internet MVP with strong consistency contract across multiple initial sources.

## Requirements

### Debt closure (mandatory Batch 0)

#### R-DEBT-V3-01
`DEF-V2-001` MUST be closed by removing tracked `__pycache__` artifacts and enforcing ignore coverage that prevents reintroduction.

#### R-DEBT-V3-02
`DEF-V2-002` MUST be closed via CI evidence automation Option A, eliminating manual URL copy as the default task evidence path.

#### R-DEBT-V3-03
`DEF-V2-003` MUST be closed by a CI quality gate that validates PR checklist-label coherence and linked-issue semantics.

### CI evidence automation Option A

#### R-CI-01
Each required CI workflow MUST emit a per-run JSON evidence artifact with deterministic schema fields (workflow name, run ID, run URL, SHA, status, timestamp, checks metadata).

#### R-CI-02
A `workflow_run` aggregator MUST consolidate workflow evidence by head SHA and produce a deterministic manifest for that SHA.

#### R-CI-03
The consolidated manifest MUST be referenceable as `openspec/evidence/<sha>.json` and MUST be surfaced in PR context (for example via automated PR comment).

#### R-CI-04
OpenSpec task evidence workflow MUST default to machine-generated references from the aggregator output and MUST NOT require manual URL copy.

### Release/versioning governance

#### R-VER-01
Release policy MUST define semver bump rules for `weave` binary changes (patch/minor/major) with explicit decision criteria.

#### R-VER-02
Any version bump MUST include synchronized update of `internal/cli/version.go` and corresponding version tests in the same change set.

#### R-VER-03
CI/release validation MUST fail when version-governance rules are violated (missing bump rationale, unsynchronized version file/tests, or inconsistent release metadata).

### PR quality gates

#### R-PRQ-01
PR quality gate MUST validate checklist-label coherence (for required labels/checklist items) before merge eligibility.

#### R-PRQ-02
PR quality gate MUST validate linked-issue semantics (`Closes #<id>` or explicit `N/A`) according to repository policy.

### Catalog internet consistency contract

#### R-CAT-01
For the same user query and the same synced catalog snapshot, listing output MUST be consistent across all initial source providers after normalization and deduplication.

#### R-CAT-02
Catalog ingestion MUST implement a multi-source provider model with at least `github_curated_index` and `registry_json` providers in v3 planning scope.

#### R-CAT-03
Ingestion MUST normalize all provider entries into a deterministic canonical schema with required fields: `entity_type`, `canonical_name`, `provider_source`, `provider_ref`, `version_ref`, `description`, `tags`, and `updated_at`.

#### R-CAT-04
Catalog identity and deduplication MUST use stable identity keys and deterministic merge rules so equivalent skills/commands from different sources yield one consistent listing record.

#### R-CAT-05
Sorting and ranking MUST be deterministic; same query + same snapshot hash MUST always produce the same ordered results.

#### R-CAT-06
Each listing record MUST include source freshness/staleness metadata and apply deterministic conflict-resolution policy when equivalent records disagree across sources.

#### R-CAT-07
Sync MUST produce a catalog snapshot with explicit snapshot version and snapshot hash; query execution MUST be bound to one snapshot to guarantee repeatability.

#### R-CAT-08
Catalog search MUST support offline execution against the latest successful synced snapshot and report snapshot metadata used for the result set.

#### R-CAT-09
Catalog internet ingestion MUST enforce trust policy (allowlisted sources + pin by ref/hash) before source data is accepted into normalization.

#### R-CAT-10
Malformed or partial source data MUST be quarantined or skipped with diagnostics and MUST NOT invalidate or corrupt committed snapshots.

## Scenarios

### S1: Deterministic evidence manifest per commit SHA
Given multiple CI workflows run for the same head SHA  
When all required workflows complete  
Then the aggregator produces a single deterministic manifest for that SHA  
And publishes reference information for OpenSpec and PR review.

### S2: No manual URL copy in planning tasks
Given OpenSpec task evidence logging for a completed batch  
When CI automation is enabled  
Then evidence references come from the generated manifest or PR comment  
And manual copy/paste of run URLs is not required for standard closure.

### S3: Deferred debt closure enforcement
Given v3 Batch 0 execution  
When debt closure checks run  
Then `DEF-V2-001`, `DEF-V2-002`, and `DEF-V2-003` closure criteria are fully satisfied  
And deferred register entries are updated with closure evidence.

### S4: Semver governance for releases
Given a release-intended change  
When version bump policy evaluation runs  
Then bump type rationale is explicit and compliant  
And `internal/cli/version.go` + version tests are synchronized.

### S5: PR quality coherence guard
Given a PR with checklist, labels, and issue-link metadata  
When quality gate validation runs  
Then coherence mismatches fail the gate with actionable errors  
And compliant PRs pass deterministically.

### S6: Multi-source consistent listing
Given the same logical skill/command exists in both `github_curated_index` and `registry_json`  
When the same query is executed on the same snapshot hash  
Then output contains one deduplicated canonical entry  
And ordering remains deterministic across repeated runs.

### S7: Snapshot-bound repeatable search
Given a synced snapshot identified by version/hash  
When search is executed online and offline against that same snapshot  
Then listing set and order are identical  
And result metadata includes the snapshot version/hash.

### S8: Trust and malformed data handling
Given a source payload that is non-allowlisted, unpinned, malformed, or partial  
When ingestion validation executes  
Then entries are rejected or quarantined with diagnostics  
And the previously committed snapshot remains valid and searchable.
