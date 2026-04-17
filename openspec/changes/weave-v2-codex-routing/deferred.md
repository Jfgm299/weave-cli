# Deferred Work Register — weave-v2-codex-routing

## Deferred Items
| ID | Requirement | Deferred From | Why Deferred | Closure Criteria | Target |
|----|-------------|---------------|--------------|------------------|--------|
| DEF-V2-001 | R-DEBT-01 / Release hygiene | v2 post-merge hardening | Python bytecode cache files under `scripts/release/__pycache__/*.pyc` were tracked in git, generating avoidable repo noise and potential cross-version artifacts. | Closed in v3 Batch 0: artifacts untracked, ignore/hygiene checks in place, and live CI evidence captured (`repo-hygiene-gate` run: https://github.com/Jfgm299/weave-cli/actions/runs/24571204951). | closed (v3 batch 0) |
| DEF-V2-002 | R-DEBT-01 / CI evidence automation | v2 post-merge hardening | Release evidence depended on manual workflow_dispatch/tag handling and manual URL copy into OpenSpec. | Closed in v3 Batch 0 baseline: deterministic evidence scripts/manifests added and live CI references captured (`install-artifact-validation`: https://github.com/Jfgm299/weave-cli/actions/runs/24571204913, `migration-note-gate`: https://github.com/Jfgm299/weave-cli/actions/runs/24571409009). | closed (v3 batch 0) |
| DEF-V2-003 | R-UX-06 / Workflow consistency | v2 process hardening | PR template/checklist could drift from actual label/issue state, causing inconsistent reviewer signals. | Closed in v3 Batch 0: metadata coherence gate added with live CI evidence (`pr-metadata-coherence`: https://github.com/Jfgm299/weave-cli/actions/runs/24571408997). | closed (v3 batch 0) |
