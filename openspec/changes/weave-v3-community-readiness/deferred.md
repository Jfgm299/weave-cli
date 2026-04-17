# Deferred Work Register — weave-v3-community-readiness

## Deferred Items
| ID | Requirement | Deferred From | Why Deferred | Closure Criteria | Target |
|----|-------------|---------------|--------------|------------------|--------|
| DEF-V3-B0-001 | R-DEBT-V3-02 / R-DEBT-V3-03 evidence finalization | v3 batch 0 implementation | Batch 0 implemented Option A baseline and PR coherence gate, but live GitHub Actions/PR-trigger execution evidence is not yet captured in this workspace session. | Record live run URL(s) for `pr-metadata-coherence.yml` and Option-A evidence flow, confirm deterministic manifest + PR reference from CI context, then flip B0-T2.4/T2.5 and B0-T3.4/T3.5 to `[x]`. | v3 batch 0 follow-up |
| DEF-V3-B0-002 | R-DEBT-V3-01 evidence finalization | v3 batch 0 implementation | Repo hygiene gate and local pass/fail evidence are complete, but missing live GitHub Actions run URL proving gate behavior in CI context. | Capture at least one fail/pass run URL for `repo-hygiene-gate.yml` and link in tasks evidence log; then flip B0-T1.4/T1.5 to `[x]`. | v3 batch 0 follow-up |
