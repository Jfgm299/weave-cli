#!/usr/bin/env python3
import argparse
import io
import json
import urllib.parse
import urllib.request
import zipfile
from pathlib import Path

from collect_workflow_evidence import REQUIRED_KEYS
from evidence_manifest import REQUIRED_WORKFLOWS, SCHEMA_VERSION, validate_sha


GITHUB_API = "https://api.github.com"


def build_failure_payload(workflow_name: str, head_sha: str, reason: str) -> dict:
    return {
        "schema_version": 1,
        "workflow_name": workflow_name,
        "run_id": "aggregator",
        "run_url": "",
        "head_sha": head_sha,
        "status": "failure",
        "timestamp": "1970-01-01T00:00:00+00:00",
        "checks": {
            "summary": f"aggregator: {reason}",
        },
    }


def validate_downloaded_payload(
    payload: dict, workflow_name: str, head_sha: str
) -> None:
    if not isinstance(payload, dict):
        raise ValueError("payload must be object")
    if payload.get("workflow_name") != workflow_name:
        raise ValueError(
            f"workflow mismatch payload={payload.get('workflow_name')} expected={workflow_name}"
        )
    if payload.get("head_sha") != head_sha:
        raise ValueError(
            f"head_sha mismatch payload={payload.get('head_sha')} expected={head_sha}"
        )
    if payload.get("schema_version") != SCHEMA_VERSION:
        raise ValueError(f"unsupported schema version {payload.get('schema_version')}")
    for key in REQUIRED_KEYS:
        if key not in payload:
            raise ValueError(f"missing required key {key}")


def api_get_json(path: str, token: str) -> dict:
    req = urllib.request.Request(
        f"{GITHUB_API}{path}",
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
        },
    )
    with urllib.request.urlopen(req) as response:
        return json.loads(response.read().decode("utf-8"))


def download_artifact_payload(archive_url: str, token: str) -> dict:
    req = urllib.request.Request(
        archive_url,
        headers={
            "Authorization": f"Bearer {token}",
            "Accept": "application/vnd.github+json",
            "X-GitHub-Api-Version": "2022-11-28",
        },
    )
    with urllib.request.urlopen(req) as response:
        content = response.read()

    with zipfile.ZipFile(io.BytesIO(content)) as zf:
        json_entries = sorted(
            [name for name in zf.namelist() if name.endswith(".json")]
        )
        if not json_entries:
            raise ValueError("artifact archive has no json entries")
        with zf.open(json_entries[0]) as raw:
            return json.loads(raw.read().decode("utf-8"))


def select_latest_completed_runs(runs: list[dict], head_sha: str) -> dict[str, dict]:
    selected: dict[str, dict] = {}
    for run in runs:
        workflow_name = run.get("name")
        if workflow_name not in REQUIRED_WORKFLOWS:
            continue
        if run.get("head_sha") != head_sha:
            continue
        if run.get("status") != "completed":
            continue

        existing = selected.get(workflow_name)
        if existing is None:
            selected[workflow_name] = run
            continue

        left = (
            int(existing.get("run_attempt") or 0),
            int(existing.get("id") or 0),
        )
        right = (
            int(run.get("run_attempt") or 0),
            int(run.get("id") or 0),
        )
        if right > left:
            selected[workflow_name] = run
    return selected


def collect_payloads(repo: str, head_sha: str, token: str) -> dict[str, dict]:
    query = urllib.parse.urlencode({"head_sha": head_sha, "per_page": 100})
    runs_payload = api_get_json(f"/repos/{repo}/actions/runs?{query}", token)
    runs = runs_payload.get("workflow_runs", [])
    selected = select_latest_completed_runs(runs, head_sha)

    collected: dict[str, dict] = {}
    for workflow_name in sorted(REQUIRED_WORKFLOWS):
        run = selected.get(workflow_name)
        if run is None:
            collected[workflow_name] = build_failure_payload(
                workflow_name, head_sha, "missing workflow run for head sha"
            )
            continue

        run_id = run["id"]
        artifacts_payload = api_get_json(
            f"/repos/{repo}/actions/runs/{run_id}/artifacts?per_page=100", token
        )
        artifacts = artifacts_payload.get("artifacts", [])
        expected_artifact = f"ci-evidence-{workflow_name}"
        match = next((a for a in artifacts if a.get("name") == expected_artifact), None)

        if match is None:
            collected[workflow_name] = build_failure_payload(
                workflow_name,
                head_sha,
                f"missing artifact '{expected_artifact}' for run {run_id}",
            )
            continue

        try:
            payload = download_artifact_payload(match["archive_download_url"], token)
            validate_downloaded_payload(payload, workflow_name, head_sha)
            collected[workflow_name] = payload
        except Exception as exc:
            collected[workflow_name] = build_failure_payload(
                workflow_name, head_sha, f"malformed artifact payload: {exc}"
            )

    return collected


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Collect required workflow evidence artifacts by head SHA"
    )
    parser.add_argument("--repo", required=True, help="owner/repo")
    parser.add_argument("--head-sha", required=True)
    parser.add_argument("--token", required=True)
    parser.add_argument("--output-dir", required=True)
    args = parser.parse_args()

    validate_sha(args.head_sha)

    output_dir = Path(args.output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)

    payloads = collect_payloads(args.repo, args.head_sha, args.token)
    for workflow_name in sorted(REQUIRED_WORKFLOWS):
        payload = payloads[workflow_name]
        output_file = output_dir / f"{workflow_name}.json"
        output_file.write_text(
            json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8"
        )

    print(f"Collected workflow evidence payloads for {args.head_sha} into {output_dir}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
