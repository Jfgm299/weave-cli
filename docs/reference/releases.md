# Release Notes Policy

## Migration Guide requirement for breaking changes

Any release that introduces breaking behavior MUST include:

- A `Migration` section with before/after examples.
- Explicit commands to migrate (`weave migrate` when applicable).
- Rollback guidance.

## Minimum release checklist

- [ ] Semver tag
- [ ] Versioned artifact names
- [ ] checksums.txt + checksums.txt.sig
- [ ] Migration section when breaking changes exist

## Automated enforcement

- Release artifact/signature generation and verification: `.github/workflows/release-artifacts.yml`
- Install validation from generated artifacts: `.github/workflows/install-artifact-validation.yml`
- Migration-note gate for PRs labeled `breaking-change`: `.github/workflows/migration-note-gate.yml`
- Repository hygiene gate for tracked Python cache artifacts: `.github/workflows/repo-hygiene-gate.yml`
- PR checklist-label-issue coherence gate: `.github/workflows/pr-metadata-coherence.yml`

## CI evidence automation (Option A full)

- Per-workflow evidence payload producer: `scripts/ci/collect_workflow_evidence.py`
- `workflow_run` evidence collector by head SHA: `scripts/ci/collect_workflow_run_evidence.py`
- Deterministic SHA manifest generator + digest metadata: `scripts/ci/evidence_manifest.py`
- PR comment upsert utility (manifest reference visibility): `scripts/ci/post_manifest_pr_comment.py`
- OpenSpec machine-reference generator (no manual URL copy): `scripts/ci/generate_openspec_evidence_refs.py`
- CI orchestration workflow: `.github/workflows/ci-evidence-aggregator.yml`
- Local baseline helper (writes `openspec/evidence/<sha>.json`, `pr-comment-<sha>.md`, `refs-<sha>.md`): `scripts/ci/run_evidence_baseline.py`
