#!/usr/bin/env python3
import argparse
import json
import re
from pathlib import Path


def validate_sha(sha: str) -> None:
    if not re.fullmatch(r"[0-9a-f]{40}", sha):
        raise ValueError(f"invalid sha '{sha}', expected 40 lowercase hex chars")


def normalize_manifest_reference(path: Path) -> str:
    normalized = path.as_posix()
    marker = "openspec/evidence/"
    idx = normalized.find(marker)
    if idx >= 0:
        return normalized[idx:]
    return normalized


def build_unavailable(head_sha: str, manifest_path: str) -> str:
    return "\n".join(
        [
            "## OpenSpec Evidence References",
            "",
            f"- SHA: `{head_sha}`",
            "- Evidence state: `unavailable`",
            f"- Expected manifest: `{manifest_path}`",
            "- Reason: deterministic manifest not found yet for this SHA.",
            "",
        ]
    )


def build_from_manifest(head_sha: str, manifest_path: str, manifest: dict) -> str:
    workflow_lines = []
    for workflow in manifest.get("workflows", []):
        name = workflow.get("workflow_name", "unknown")
        status = workflow.get("status", "unknown")
        workflow_lines.append(f"- `{name}`: `{status}`")

    return "\n".join(
        [
            "## OpenSpec Evidence References",
            "",
            "Machine-generated references (no manual workflow URL copy required).",
            f"- SHA: `{head_sha}`",
            f"- Manifest: `{manifest_path}`",
            f"- PR comment snapshot: `openspec/evidence/pr-comment-{head_sha}.md`",
            "",
            "### Workflow status summary",
            *workflow_lines,
            "",
        ]
    )


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Generate OpenSpec evidence references from deterministic manifest"
    )
    parser.add_argument("--head-sha", required=True)
    parser.add_argument("--manifest-file", required=True)
    parser.add_argument("--output", required=True)
    args = parser.parse_args()

    validate_sha(args.head_sha)

    manifest_path = Path(args.manifest_file)
    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)

    manifest_reference = normalize_manifest_reference(manifest_path)

    if not manifest_path.exists():
        output_path.write_text(
            build_unavailable(args.head_sha, manifest_reference), encoding="utf-8"
        )
        print(f"Wrote deterministic unavailable evidence state to {output_path}")
        return 0

    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    if manifest.get("head_sha") != args.head_sha:
        raise ValueError(
            f"manifest head_sha mismatch: {manifest.get('head_sha')} != {args.head_sha}"
        )

    output_path.write_text(
        build_from_manifest(args.head_sha, manifest_reference, manifest),
        encoding="utf-8",
    )
    print(f"Wrote machine-generated evidence references to {output_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
