#!/usr/bin/env python3
import argparse
import os
import subprocess
import sys
import tempfile
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
COLLECT = ROOT / "scripts" / "ci" / "collect_workflow_evidence.py"
MANIFEST = ROOT / "scripts" / "ci" / "evidence_manifest.py"
WORKFLOWS = [
    "install-artifact-validation",
    "migration-note-gate",
    "release-artifacts",
    "pr-metadata-coherence",
    "repo-hygiene-gate",
]


def run(cmd: list[str], env: dict[str, str]) -> None:
    result = subprocess.run(
        cmd, cwd=str(ROOT), env=env, text=True, capture_output=True, check=False
    )
    if result.returncode != 0:
        raise RuntimeError(
            f"command failed: {' '.join(cmd)}\nstdout:\n{result.stdout}\nstderr:\n{result.stderr}"
        )


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Generate local Option-A evidence baseline for current HEAD"
    )
    parser.add_argument("--sha", help="Explicit SHA (defaults to current HEAD)")
    args = parser.parse_args()

    sha = args.sha
    if not sha:
        result = subprocess.run(
            ["git", "rev-parse", "HEAD"],
            cwd=str(ROOT),
            text=True,
            capture_output=True,
            check=False,
        )
        if result.returncode != 0:
            print(f"Unable to resolve HEAD sha: {result.stderr}", file=sys.stderr)
            return 1
        sha = result.stdout.strip()

    manifest_output = ROOT / "openspec" / "evidence" / f"{sha}.json"
    comment_output = ROOT / "openspec" / "evidence" / f"pr-comment-{sha}.md"

    with tempfile.TemporaryDirectory() as temp_dir:
        temp_path = Path(temp_dir)
        env = os.environ.copy()
        env.update(
            {
                "GITHUB_SHA": sha,
                "GITHUB_RUN_ID": "local-baseline",
                "GITHUB_REPOSITORY": "Jfgm299/weave-cli",
            }
        )

        for workflow in WORKFLOWS:
            output = temp_path / f"{workflow}.json"
            run(
                [
                    sys.executable,
                    str(COLLECT),
                    "--workflow-name",
                    workflow,
                    "--status",
                    "success",
                    "--output",
                    str(output),
                ],
                env,
            )

        run(
            [
                sys.executable,
                str(MANIFEST),
                "--head-sha",
                sha,
                "--input-dir",
                str(temp_path),
                "--manifest-output",
                str(manifest_output),
                "--pr-comment-output",
                str(comment_output),
            ],
            env,
        )

    print(f"Baseline manifest generated: {manifest_output}")
    print(f"PR comment reference generated: {comment_output}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
