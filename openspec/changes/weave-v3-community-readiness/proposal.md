# Proposal — weave-v3-community-readiness

## Context

v2 delivered codex routing and transaction hardening, but community-readiness is still limited by process debt and non-deterministic evidence collection.

Deferred items from v2 remain open:

- DEF-V2-001: tracked `__pycache__` artifacts.
- DEF-V2-002: manual CI evidence capture.
- DEF-V2-003: PR checklist/label coherence drift.

## Goals

1. Adopt CI evidence automation **Option A** end-to-end:
   - each workflow emits per-run JSON evidence artifact,
   - `workflow_run` aggregator consolidates by head SHA,
   - deterministic manifest is produced and referenced,
   - no manual URL copy required in OpenSpec tasks.
2. Close v2 deferred debt in Batch 0 before feature hardening.
3. Define and enforce release/versioning governance for the `weave` binary:
   - semver bump rules,
   - explicit version update process for `internal/cli/version.go` + version tests.
4. Add PR quality gates that enforce checklist-label coherence.

## Non-Goals

- New runtime product features unrelated to CI governance and release quality.
- Replacing current provider architecture.
- Introducing remote registries or marketplace behavior.

## Key Decisions

- Use SHA-addressed deterministic manifests as the single source of CI audit truth.
- Keep OpenSpec evidence references machine-generated (artifact + PR comment), not manually curated.
- Treat version governance as a release contract, not best-effort documentation.

## Risks

- Aggregator race conditions if multiple workflows for the same SHA complete out-of-order.
- False positives/false negatives in checklist-label coherence parsing if PR template variants drift.
- Version governance can regress if tests are not mandatory in release CI gates.
