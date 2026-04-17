#!/usr/bin/env python3
import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


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


if __name__ == "__main__":
    unittest.main()
