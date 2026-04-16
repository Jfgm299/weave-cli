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
