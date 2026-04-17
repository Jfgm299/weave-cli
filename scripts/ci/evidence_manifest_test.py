#!/usr/bin/env python3
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SCRIPT_PATH = Path(__file__).with_name("evidence_manifest.py")


def write_payload(path: Path, workflow: str, sha: str) -> None:
    payload = {
        "schema_version": 1,
        "workflow_name": workflow,
        "run_id": "12345",
        "run_url": "https://github.com/Jfgm299/weave-cli/actions/runs/12345",
        "head_sha": sha,
        "status": "success",
        "timestamp": "2026-04-17T12:00:00+00:00",
        "checks": {"summary": "ok"},
    }
    path.write_text(
        json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )


class EvidenceManifestTests(unittest.TestCase):
    def test_builds_manifest_and_marks_missing_workflows(self) -> None:
        sha = "a" * 40
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            inputs = root / "inputs"
            inputs.mkdir(parents=True)
            write_payload(
                inputs / "install-artifact-validation.json",
                "install-artifact-validation",
                sha,
            )
            write_payload(inputs / "release-artifacts.json", "release-artifacts", sha)

            manifest_output = root / "openspec" / "evidence" / f"{sha}.json"
            pr_output = root / "pr-comment.md"

            result = subprocess.run(
                [
                    sys.executable,
                    str(SCRIPT_PATH),
                    "--head-sha",
                    sha,
                    "--input-dir",
                    str(inputs),
                    "--manifest-output",
                    str(manifest_output),
                    "--pr-comment-output",
                    str(pr_output),
                ],
                text=True,
                capture_output=True,
                check=False,
            )

            self.assertEqual(result.returncode, 0, msg=result.stderr)
            manifest = json.loads(manifest_output.read_text(encoding="utf-8"))
            self.assertEqual(manifest["head_sha"], sha)
            names = [w["workflow_name"] for w in manifest["workflows"]]
            self.assertEqual(names, sorted(names))
            statuses = {w["workflow_name"]: w["status"] for w in manifest["workflows"]}
            self.assertEqual(statuses["migration-note-gate"], "missing")
            self.assertEqual(statuses["release-artifacts"], "success")

            comment = pr_output.read_text(encoding="utf-8")
            self.assertIn(f"openspec/evidence/{sha}.json", comment)

    def test_rejects_payload_missing_required_keys(self) -> None:
        sha = "b" * 40
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            inputs = root / "inputs"
            inputs.mkdir(parents=True)
            invalid_payload = {
                "schema_version": 1,
                "workflow_name": "release-artifacts",
                "head_sha": sha,
            }
            (inputs / "release-artifacts.json").write_text(
                json.dumps(invalid_payload, indent=2, sort_keys=True) + "\n",
                encoding="utf-8",
            )
            manifest_output = root / "manifest.json"

            result = subprocess.run(
                [
                    sys.executable,
                    str(SCRIPT_PATH),
                    "--head-sha",
                    sha,
                    "--input-dir",
                    str(inputs),
                    "--manifest-output",
                    str(manifest_output),
                ],
                text=True,
                capture_output=True,
                check=False,
            )

            self.assertEqual(result.returncode, 1)
            self.assertIn("missing required keys", result.stderr)


if __name__ == "__main__":
    unittest.main()
