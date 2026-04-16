# Test Strategy — weave-v2-codex-routing

## Purpose

Define auditable TDD coverage for v2 delta with debt-first execution.

This strategy complements:

- `openspec/changes/weave-v2-codex-routing/spec.md`
- `openspec/changes/weave-v2-codex-routing/tasks.md`

## 1) Test levels

### Unit

- Provider registry includes `codex` and exposes deterministic strategy metadata.
- Command projection strategy resolves targets by provider and install mode.
- No-provider prompt policy is deterministic (interactive prompt vs non-interactive fail).
- Rollback planner composes complete compensating operations.

### Integration

- Multi-provider command add applies all projections under one transaction envelope.
- Any provider-path failure triggers full rollback and zero config mutation.
- Exclusive `--provider codex` does not require shared `.agents` command path.
- Guided git-init flow matches interactive/non-interactive policies.

### E2E / CI

- `weave provider add codex` parity behavior.
- `weave command add <name>` with mixed providers enabled.
- no-provider interactive prompt behavior.
- no-provider non-interactive auto-fail behavior.
- release/signing/install-validation debt closure checks (DEF-001/002/003).

## 2) Mandatory requirement coverage rule

For each requirement in `spec.md`:

- at least one success-path test
- at least one edge/error test
- at least one observable acceptance evidence (filesystem, config, output, exit code, CI check)

No requirement can be marked `[x]` in tasks without this rule.

## 3) Batch-0 debt-specific evidence

Batch 0 MUST collect explicit closure evidence for:

- DEF-001: release checksum/signature automation in workflow.
- DEF-002: install validation from release artifacts.
- DEF-003: migration-note CI gate for breaking changes.
- DEF-004: guided git-init behavior with tests.

Evidence format:

- workflow path
- test name(s)
- run/log reference
- acceptance statement
