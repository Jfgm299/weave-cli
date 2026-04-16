# Deferred Work Register — weave-v1-mvp

This register tracks intentionally deferred items.

Policy:

- A requirement cannot be marked `[x]` in `tasks.md` unless it is 100% implemented and operational.
- If a requirement is partially closed as docs/policy baseline, the remaining executable work MUST be tracked here.

## Deferred Items

| ID | Requirement | Deferred From | Why Deferred | Closure Criteria | Target |
|----|-------------|---------------|--------------|------------------|--------|
| DEF-001 | R-DIST-01 | B6-T12.4 | Release signing automation is not implemented in v1 baseline | Add automated checksum+signature generation in release workflow and validate in CI | v1.x |
| DEF-002 | R-DIST-02 | B6-T13.4 | Installation artifact pipeline validation is not automated | Add automated install validation against release artifacts in CI | v1.x |
| DEF-003 | R-UPD-02 | B6-T16.4 | Release-note migration enforcement is policy-only | Add CI gate that fails release PRs missing migration section when breaking changes are present | v1.x |
| DEF-004 | R-DEP-01 | Post-v1 | Missing Git root currently returns an error with no guided initialization prompt | Add interactive/non-interactive prompt flow to offer `git init` when no `.git` root is detected | v1.x |

## Carry-over to v2

Deferred implementation closure for `DEF-001..DEF-004` is scheduled as Batch 0 in:

- `openspec/changes/weave-v2-codex-routing/tasks.md`

Current execution status (v2 Batch 0):

- DEF-001: closure criteria implemented and evidenced via local script tests + CI run URL in `openspec/changes/weave-v2-codex-routing/tasks.md` (B0 CI Run URLs section).
- DEF-002: closure criteria implemented and evidenced via local install validation + CI run URL in `openspec/changes/weave-v2-codex-routing/tasks.md` (B0 CI Run URLs section).
- DEF-003: closure criteria implemented and evidenced via local migration gate tests + CI run URL in `openspec/changes/weave-v2-codex-routing/tasks.md` (B0 CI Run URLs section).
- DEF-004: guided `git init` behavior implemented and evidenced by unit tests plus targeted e2e for both non-interactive and interactive-decline branches.

Batch 0 status: closed in `openspec/changes/weave-v2-codex-routing/tasks.md`.
