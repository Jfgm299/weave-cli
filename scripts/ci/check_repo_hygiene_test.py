#!/usr/bin/env python3
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SCRIPT_PATH = Path(__file__).with_name("check_repo_hygiene.py")


def run_hygiene(paths: list[str]) -> subprocess.CompletedProcess[str]:
    with tempfile.NamedTemporaryFile("w", delete=False, suffix=".txt") as tmp:
        tmp.write("\n".join(paths) + "\n")
        sample_path = tmp.name

    try:
        return subprocess.run(
            [sys.executable, str(SCRIPT_PATH), "--from-file", sample_path],
            text=True,
            capture_output=True,
            check=False,
        )
    finally:
        Path(sample_path).unlink(missing_ok=True)


class RepoHygieneTests(unittest.TestCase):
    def test_passes_when_no_python_cache_artifacts_are_tracked(self) -> None:
        result = run_hygiene(
            [
                "scripts/release/check_migration_gate.py",
                "scripts/release/release_artifacts_test.py",
                "README.md",
            ]
        )

        self.assertEqual(result.returncode, 0)
        self.assertIn("Repository hygiene passed", result.stdout)

    def test_fails_when_tracked_pycache_artifacts_are_present(self) -> None:
        result = run_hygiene(
            [
                "scripts/release/__pycache__/check_migration_gate_test.cpython-314.pyc",
                "scripts/release/check_migration_gate.py",
            ]
        )

        self.assertEqual(result.returncode, 1)
        self.assertIn("Repository hygiene failed", result.stderr)
        self.assertIn("__pycache__", result.stderr)


if __name__ == "__main__":
    unittest.main()
