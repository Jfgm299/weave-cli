# Test Strategy — weave-v3-community-readiness

## Purpose

Define auditable verification strategy for v3 planning scope with strict debt-first and deterministic-evidence discipline.

This strategy includes catalog-internet consistency verification with snapshot-bound determinism.

This strategy complements:

- `openspec/changes/weave-v3-community-readiness/spec.md`
- `openspec/changes/weave-v3-community-readiness/tasks.md`

## 1) Test levels

### Unit

- Evidence artifact schema validation and deterministic serializer behavior.
- Aggregator consolidation logic by head SHA.
- PR checklist/label/issue-link rule parser behavior.
- Version governance rule evaluation (semver bump classification + file/test sync checks).
- Catalog normalization schema validation and stable identity key derivation.
- Catalog dedup/merge conflict policy determinism and deterministic ranking tie-breakers.
- Trust policy validation (allowlist + pin-by-ref/hash) and malformed-data quarantine behavior.

### Integration

- Workflow evidence producers + aggregator manifest pipeline.
- PR quality gate execution against representative PR payloads.
- Version-governance gate execution against release-intended change payloads.
- Multi-source catalog ingestion (`github_curated_index` + `registry_json`) into one canonical snapshot.
- Snapshot commit/read path for offline query execution.
- Freshness/staleness metadata propagation and conflict resolution signaling in listing output.

### E2E / CI

- Full required workflow set emits evidence and is consolidated into one SHA manifest.
- PR receives deterministic manifest summary comment/reference.
- Coherence gate blocks inconsistent checklist/label/issue-link cases.
- Version governance gate blocks unsynchronized `version.go` and version tests.
- Same query + same snapshot hash returns identical listing output across repeated runs.
- Online/offline search parity holds for the same committed snapshot.
- Trust policy rejects non-allowlisted/unpinned sources deterministically.
- Malformed/partial source payloads are skipped/quarantined without corrupting active snapshot.

## 2) Mandatory requirement coverage rule

For each requirement in `spec.md`:

- at least one success-path test,
- at least one edge/error test,
- at least one observable acceptance evidence item (artifact content, manifest reference, PR comment, gate status, or exit code).

No requirement can be marked `[x]` in tasks without this rule.

This includes all catalog requirements `R-CAT-01..10`.

## 3) Batch-0 debt closure evidence requirements

Batch 0 MUST collect explicit closure evidence for:

- DEF-V2-001 (`__pycache__` hygiene closed and guarded).
- DEF-V2-002 (manual evidence capture replaced by Option A automation).
- DEF-V2-003 (PR coherence gate active and enforced).

Evidence format:

- workflow path or gate path,
- test name(s),
- deterministic manifest reference by SHA,
- closure statement linked to deferred register update.

## 4) Deterministic manifest assertions

For each tested SHA manifest:

- schema version present,
- required workflows represented,
- stable ordering verified,
- pass/fail/missing states explicit,
- PR-facing reference matches artifact payload.

## 5) Version governance assertions

For release-intended changes:

- semver bump rationale is present and valid,
- `internal/cli/version.go` changes are synchronized with version tests,
- governance gate outcome is deterministic across repeated runs.

## 6) Catalog internet consistency assertions

For each tested catalog snapshot:

- normalization schema required fields are present for all accepted entries,
- stable identity keys are reproducible across source-order permutations,
- dedup emits one canonical listing entry for equivalent multi-source records,
- deterministic ranking/sort is stable for repeated identical queries,
- freshness/staleness metadata and conflict state are present per listing record,
- snapshot metadata (`snapshot_version`, `snapshot_hash`) is present in sync/search outputs.

## 7) Snapshot repeatability and offline parity assertions

- Query replay contract: same query + same snapshot hash => byte-identical ordered listing payload.
- Offline parity: offline search against committed snapshot matches online output bound to same snapshot hash.
- Snapshot immutability: malformed-source ingestion attempts MUST NOT mutate the active committed snapshot.

## 8) Trust and malformed-data assertions

- Non-allowlisted sources are rejected before normalization.
- Unpinned sources (missing immutable ref/hash pin) are rejected deterministically.
- Malformed records are skipped or quarantined with actionable diagnostics.
- Partial records do not corrupt committed snapshots; valid records continue through deterministic pipeline.
