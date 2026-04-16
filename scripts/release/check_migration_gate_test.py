#!/usr/bin/env python3
import json
import os
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SCRIPT_PATH = Path(__file__).with_name("check_migration_gate.py")


def run_gate(event_payload: dict) -> subprocess.CompletedProcess[str]:
    with tempfile.NamedTemporaryFile("w", delete=False, suffix=".json") as tmp:
        json.dump(event_payload, tmp)
        event_path = tmp.name

    env = os.environ.copy()
    env["GITHUB_EVENT_PATH"] = event_path

    try:
        return subprocess.run(
            [sys.executable, str(SCRIPT_PATH)],
            env=env,
            text=True,
            capture_output=True,
            check=False,
        )
    finally:
        os.unlink(event_path)


class MigrationGateTests(unittest.TestCase):
    def test_passes_without_breaking_change_label(self) -> None:
        result = run_gate(
            {
                "pull_request": {
                    "labels": [{"name": "enhancement"}],
                    "body": "## Migration\nNot required",
                }
            }
        )

        self.assertEqual(result.returncode, 0)
        self.assertIn("No breaking-change label present", result.stdout)

    def test_passes_with_breaking_change_and_migration_heading(self) -> None:
        result = run_gate(
            {
                "pull_request": {
                    "labels": [{"name": "breaking-change"}],
                    "body": "## Migration\nRun weave migrate",
                }
            }
        )

        self.assertEqual(result.returncode, 0)
        self.assertIn("Migration section found", result.stdout)

    def test_passes_with_deeper_migration_heading(self) -> None:
        result = run_gate(
            {
                "pull_request": {
                    "labels": [{"name": "breaking-change"}],
                    "body": "### Migration\nRun weave migrate",
                }
            }
        )

        self.assertEqual(result.returncode, 0)
        self.assertIn("Migration section found", result.stdout)

    def test_fails_when_breaking_change_missing_migration_heading(self) -> None:
        result = run_gate(
            {
                "pull_request": {
                    "labels": [{"name": "breaking-change"}],
                    "body": "## Notes\nNo migration details",
                }
            }
        )

        self.assertEqual(result.returncode, 1)
        self.assertIn("Migration gate failed", result.stderr)


if __name__ == "__main__":
    unittest.main()
