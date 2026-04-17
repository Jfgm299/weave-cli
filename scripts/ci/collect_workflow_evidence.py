#!/usr/bin/env python3
import argparse
import json
import os
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path


REQUIRED_KEYS = [
    "schema_version",
    "workflow_name",
    "run_id",
    "run_url",
    "head_sha",
    "status",
    "timestamp",
    "checks",
]


def current_sha() -> str:
    sha = os.environ.get("GITHUB_SHA")
    if sha:
        return sha

    result = subprocess.run(
        ["git", "rev-parse", "HEAD"], text=True, capture_output=True, check=False
    )
    if result.returncode != 0:
        raise RuntimeError(f"unable to resolve git sha: {result.stderr.strip()}")
    return result.stdout.strip()


def build_run_url(repo: str, run_id: str) -> str:
    if repo and run_id:
        return f"https://github.com/{repo}/actions/runs/{run_id}"
    return ""


def create_payload(workflow_name: str, status: str) -> dict:
    run_id = os.environ.get("GITHUB_RUN_ID", "local")
    repo = os.environ.get("GITHUB_REPOSITORY", "")
    payload = {
        "schema_version": 1,
        "workflow_name": workflow_name,
        "run_id": run_id,
        "run_url": build_run_url(repo, run_id),
        "head_sha": current_sha(),
        "status": status,
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "checks": {
            "summary": f"{workflow_name}:{status}",
        },
    }
    return payload


def validate_payload(payload: dict) -> None:
    missing = [key for key in REQUIRED_KEYS if key not in payload]
    if missing:
        raise ValueError(f"missing required keys: {', '.join(missing)}")


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Emit per-workflow evidence JSON payload"
    )
    parser.add_argument("--workflow-name", required=True)
    parser.add_argument(
        "--status",
        default="success",
        choices=["success", "failure", "cancelled", "neutral", "skipped"],
    )
    parser.add_argument("--output", required=True)
    args = parser.parse_args()

    try:
        payload = create_payload(args.workflow_name, args.status)
        validate_payload(payload)
    except Exception as exc:
        print(f"Failed to produce evidence payload: {exc}", file=sys.stderr)
        return 1

    output_path = Path(args.output)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(
        json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )
    print(f"Wrote evidence payload to {output_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
