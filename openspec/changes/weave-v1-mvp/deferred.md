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

- DEF-001: workflow/script implementation complete with local unit evidence (`scripts/release/release_artifacts_test.py`); pending live CI workflow run evidence before closure.
- DEF-002: install-validation workflow wired and locally exercised through artifact tests; pending live CI workflow run evidence before closure.
- DEF-003: migration-note gate implementation complete with local unit evidence (`scripts/release/check_migration_gate_test.py`); pending live CI workflow run evidence before closure.
- DEF-004: guided `git init` flow implementation complete with targeted unit tests (`internal/cli/forge_handler_test.go`) and executed non-interactive e2e evidence (`test/e2e/git_init_guidance_e2e_test.go`); pending interactive e2e branch evidence before closure.
