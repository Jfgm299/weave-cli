#!/usr/bin/env python3
import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent))

from collect_workflow_evidence import validate_payload


SCRIPT_PATH = Path(__file__).with_name("collect_workflow_evidence.py")


class CollectWorkflowEvidenceTests(unittest.TestCase):
    def test_emits_valid_payload_with_required_keys(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            output = Path(temp_dir) / "release-artifacts.json"
            env = os.environ.copy()
            env.update(
                {
                    "GITHUB_SHA": "c" * 40,
                    "GITHUB_RUN_ID": "999",
                    "GITHUB_REPOSITORY": "Jfgm299/weave-cli",
                }
            )
            result = subprocess.run(
                [
                    sys.executable,
                    str(SCRIPT_PATH),
                    "--workflow-name",
                    "release-artifacts",
                    "--status",
                    "success",
                    "--output",
                    str(output),
                ],
                env=env,
                text=True,
                capture_output=True,
                check=False,
            )
            self.assertEqual(result.returncode, 0, msg=result.stderr)

            payload = json.loads(output.read_text(encoding="utf-8"))
            for key in [
                "schema_version",
                "workflow_name",
                "run_id",
                "run_url",
                "head_sha",
                "status",
                "timestamp",
                "checks",
            ]:
                self.assertIn(key, payload)

            self.assertEqual(payload["head_sha"], "c" * 40)
            self.assertEqual(payload["workflow_name"], "release-artifacts")

            raw = output.read_text(encoding="utf-8")
            self.assertIn('"checks"', raw)
            self.assertIn('"head_sha"', raw)
            self.assertLess(raw.index('"checks"'), raw.index('"head_sha"'))

    def test_validate_payload_rejects_missing_required_fields(self) -> None:
        invalid_payload = {
            "schema_version": 1,
            "workflow_name": "release-artifacts",
            "run_id": "123",
            "run_url": "https://github.com/Jfgm299/weave-cli/actions/runs/123",
            "head_sha": "a" * 40,
            "status": "success",
            "timestamp": "2026-04-17T12:00:00+00:00",
        }

        with self.assertRaisesRegex(ValueError, "missing required keys"):
            validate_payload(invalid_payload)


if __name__ == "__main__":
    unittest.main()
