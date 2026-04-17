#!/usr/bin/env python3
import argparse
import hashlib
import json
import re
import sys
from pathlib import Path


SCHEMA_VERSION = 1
REQUIRED_WORKFLOWS = [
    "install-artifact-validation",
    "migration-note-gate",
    "release-artifacts",
    "pr-metadata-coherence",
    "repo-hygiene-gate",
    "version-governance-gate",
]
REQUIRED_KEYS = {
    "schema_version",
    "workflow_name",
    "run_id",
    "run_url",
    "head_sha",
    "status",
    "timestamp",
    "checks",
}


def load_payload(path: Path) -> dict:
    data = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise ValueError(f"payload must be object: {path}")
    missing = sorted(REQUIRED_KEYS - set(data.keys()))
    if missing:
        raise ValueError(f"{path}: missing required keys: {', '.join(missing)}")
    if data.get("schema_version") != SCHEMA_VERSION:
        raise ValueError(
            f"{path}: unsupported schema_version={data.get('schema_version')}"
        )
    return data


def validate_sha(sha: str) -> None:
    if not re.fullmatch(r"[0-9a-f]{40}", sha):
        raise ValueError(f"invalid sha '{sha}', expected 40 lowercase hex chars")


def build_manifest(head_sha: str, payloads: list[dict]) -> dict:
    grouped = {payload["workflow_name"]: payload for payload in payloads}

    workflows = []
    for name in sorted(REQUIRED_WORKFLOWS):
        payload = grouped.get(name)
        if payload is None:
            workflows.append(
                {
                    "workflow_name": name,
                    "status": "missing",
                    "head_sha": head_sha,
                }
            )
        else:
            if payload["head_sha"] != head_sha:
                raise ValueError(
                    f"payload for workflow '{name}' has head_sha={payload['head_sha']} but expected {head_sha}"
                )
            workflows.append(payload)

    manifest_core = {
        "schema_version": SCHEMA_VERSION,
        "head_sha": head_sha,
        "workflows": workflows,
    }

    canonical = json.dumps(manifest_core, sort_keys=True, separators=(",", ":"))
    manifest_digest = hashlib.sha256(canonical.encode("utf-8")).hexdigest()

    return {
        **manifest_core,
        "snapshot": {
            "workflow_count": len(workflows),
            "manifest_digest_sha256": manifest_digest,
        },
    }


def write_pr_comment(
    output_file: Path, manifest_path: str, head_sha: str, workflows: list[dict]
) -> None:
    lines = [
        "## CI Evidence Manifest",
        "",
        f"- SHA: `{head_sha}`",
        f"- Manifest: `{manifest_path}`",
        "",
        "| Workflow | Status |",
        "|---|---|",
    ]

    for workflow in workflows:
        lines.append(
            f"| `{workflow['workflow_name']}` | `{workflow.get('status', 'unknown')}` |"
        )

    output_file.parent.mkdir(parents=True, exist_ok=True)
    output_file.write_text("\n".join(lines) + "\n", encoding="utf-8")


def manifest_reference_path(path: Path) -> str:
    normalized = path.as_posix()
    marker = "openspec/evidence/"
    idx = normalized.find(marker)
    if idx >= 0:
        return normalized[idx:]
    return normalized


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate workflow evidence and build deterministic SHA manifest"
    )
    parser.add_argument("--head-sha", required=True)
    parser.add_argument("--input-dir", required=True)
    parser.add_argument("--manifest-output", required=True)
    parser.add_argument("--pr-comment-output")
    args = parser.parse_args()

    try:
        validate_sha(args.head_sha)
        input_dir = Path(args.input_dir)
        payload_paths = sorted(input_dir.glob("*.json"), key=lambda p: p.name)
        payloads = [load_payload(path) for path in payload_paths]
        manifest = build_manifest(args.head_sha, payloads)
    except Exception as exc:
        print(f"Evidence manifest generation failed: {exc}", file=sys.stderr)
        return 1

    manifest_output = Path(args.manifest_output)
    manifest_output.parent.mkdir(parents=True, exist_ok=True)
    manifest_output.write_text(
        json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    print(f"Wrote deterministic manifest to {manifest_output}")

    if args.pr_comment_output:
        manifest_reference = manifest_reference_path(manifest_output)
        write_pr_comment(
            Path(args.pr_comment_output),
            manifest_reference,
            args.head_sha,
            manifest["workflows"],
        )
        print(f"Wrote PR comment summary to {args.pr_comment_output}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
