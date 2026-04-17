# Deferred Work Register — weave-v3-community-readiness

## Deferred Items
| ID | Requirement | Deferred From | Why Deferred | Closure Criteria | Target |
|----|-------------|---------------|--------------|------------------|--------|
| DEF-V3-B1-001 | R-CI-01 / R-CI-02 / R-CI-03 / R-CI-04 | Batch 1 (Option A full) | Live GitHub execution proof is not yet captured in this branch run: `workflow_run` aggregator needs real upstream workflow completions and PR event linkage to confirm end-to-end automation in hosted CI. | Capture real Actions evidence for one PR SHA showing: (1) all required workflow artifacts uploaded, (2) `ci-evidence-aggregator` run completed, (3) manifest artifact `openspec/evidence/<sha>.json` produced, (4) PR comment upserted with same manifest reference, (5) OpenSpec refs snapshot generated from manifest without manual URL copy. | Next CI pass on active PR |
