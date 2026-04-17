#!/usr/bin/env python3
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SCRIPT_PATH = Path(__file__).with_name("check_version_governance.py")


def run_gate(args: list[str]) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, str(SCRIPT_PATH), *args],
        text=True,
        capture_output=True,
        check=False,
    )


def write_event(path: Path, *, semver_label: str | None, rationale: str) -> None:
    labels = []
    if semver_label:
        labels.append({"name": semver_label})

    payload = {
        "pull_request": {
            "base": {"sha": "a" * 40},
            "head": {"sha": "b" * 40},
            "labels": labels,
            "body": f"## Release\n\nSemver rationale: {rationale}\n",
        }
    }
    path.write_text(
        json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8"
    )


def write_version_files(root: Path, *, version: str) -> tuple[Path, Path]:
    version_go = root / "version.go"
    version_test = root / "version_test.go"

    version_go.write_text(
        "package cli\n\n"
        + f'const version = "{version}"\n\n'
        + "func Version() string {\n\treturn version\n}\n",
        encoding="utf-8",
    )
    version_test.write_text(
        'package cli\n\nimport "testing"\n\n'
        + "func TestVersion_IsSemver(t *testing.T) {\n"
        + f'\tif Version() != "{version}" {{\n'
        + '\t\tt.Fatalf("unexpected version")\n'
        + "\t}\n"
        + "}\n",
        encoding="utf-8",
    )
    return version_go, version_test


class VersionGovernanceTests(unittest.TestCase):
    def test_passes_for_patch_bump_with_synced_files_and_metadata(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            changed = root / "changed.txt"
            event = root / "event.json"
            version_go, version_test = write_version_files(root, version="0.1.1")

            changed.write_text(
                f"{version_go}\n{version_test}\n",
                encoding="utf-8",
            )
            write_event(
                event, semver_label="semver:patch", rationale="Fix-only release"
            )

            result = run_gate(
                [
                    "--event-file",
                    str(event),
                    "--changed-files-file",
                    str(changed),
                    "--base-version",
                    "0.1.0",
                    "--current-version",
                    "0.1.1",
                    "--version-file",
                    str(version_go),
                    "--version-test-file",
                    str(version_test),
                ]
            )

            self.assertEqual(result.returncode, 0, msg=result.stderr)
            self.assertIn("passed", result.stdout)

    def test_fails_when_version_source_and_test_are_not_synchronized(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            changed = root / "changed.txt"
            event = root / "event.json"
            version_go, version_test = write_version_files(root, version="0.1.1")

            changed.write_text(f"{version_go}\n", encoding="utf-8")
            write_event(
                event, semver_label="semver:patch", rationale="Fix-only release"
            )

            result = run_gate(
                [
                    "--event-file",
                    str(event),
                    "--changed-files-file",
                    str(changed),
                    "--base-version",
                    "0.1.0",
                    "--current-version",
                    "0.1.1",
                    "--version-file",
                    str(version_go),
                    "--version-test-file",
                    str(version_test),
                ]
            )

            self.assertEqual(result.returncode, 1)
            self.assertIn("must be updated together", result.stderr)

    def test_fails_when_declared_bump_does_not_match_computed_transition(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            changed = root / "changed.txt"
            event = root / "event.json"
            version_go, version_test = write_version_files(root, version="0.1.1")

            changed.write_text(
                f"{version_go}\n{version_test}\n",
                encoding="utf-8",
            )
            write_event(
                event, semver_label="semver:minor", rationale="Fix-only release"
            )

            result = run_gate(
                [
                    "--event-file",
                    str(event),
                    "--changed-files-file",
                    str(changed),
                    "--base-version",
                    "0.1.0",
                    "--current-version",
                    "0.1.1",
                    "--version-file",
                    str(version_go),
                    "--version-test-file",
                    str(version_test),
                ]
            )

            self.assertEqual(result.returncode, 1)
            self.assertIn("does not match computed bump", result.stderr)

    def test_fails_when_rationale_is_missing(self) -> None:
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            changed = root / "changed.txt"
            event = root / "event.json"
            version_go, version_test = write_version_files(root, version="0.1.1")

            changed.write_text(
                f"{version_go}\n{version_test}\n",
                encoding="utf-8",
            )
            payload = {
                "pull_request": {
                    "base": {"sha": "a" * 40},
                    "head": {"sha": "b" * 40},
                    "labels": [{"name": "semver:patch"}],
                    "body": "## Release\n\nNo rationale field here\n",
                }
            }
            event.write_text(
                json.dumps(payload, indent=2, sort_keys=True) + "\n", encoding="utf-8"
            )

            result = run_gate(
                [
                    "--event-file",
                    str(event),
                    "--changed-files-file",
                    str(changed),
                    "--base-version",
                    "0.1.0",
                    "--current-version",
                    "0.1.1",
                    "--version-file",
                    str(version_go),
                    "--version-test-file",
                    str(version_test),
                ]
            )

            self.assertEqual(result.returncode, 1)
            self.assertIn("missing semver rationale", result.stderr)


if __name__ == "__main__":
    unittest.main()
