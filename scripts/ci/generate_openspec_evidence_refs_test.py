#!/usr/bin/env python3
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SCRIPT_PATH = Path(__file__).with_name("generate_openspec_evidence_refs.py")


class GenerateOpenSpecEvidenceRefsTests(unittest.TestCase):
    def test_generates_manifest_references_when_manifest_exists(self) -> None:
        sha = "b" * 40
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            manifest = root / "openspec" / "evidence" / f"{sha}.json"
            output = root / "refs.md"
            manifest.parent.mkdir(parents=True)

            manifest.write_text(
                json.dumps(
                    {
                        "schema_version": 1,
                        "head_sha": sha,
                        "workflows": [
                            {
                                "workflow_name": "release-artifacts",
                                "status": "success",
                            }
                        ],
                    },
                    indent=2,
                    sort_keys=True,
                )
                + "\n",
                encoding="utf-8",
            )

            result = subprocess.run(
                [
                    sys.executable,
                    str(SCRIPT_PATH),
                    "--head-sha",
                    sha,
                    "--manifest-file",
                    str(manifest),
                    "--output",
                    str(output),
                ],
                text=True,
                capture_output=True,
                check=False,
            )

            self.assertEqual(result.returncode, 0, msg=result.stderr)
            content = output.read_text(encoding="utf-8")
            self.assertIn("Machine-generated references", content)
            self.assertIn("release-artifacts", content)
            self.assertIn(f"openspec/evidence/{sha}.json", content)

    def test_generates_unavailable_state_when_manifest_missing(self) -> None:
        sha = "c" * 40
        with tempfile.TemporaryDirectory() as temp_dir:
            root = Path(temp_dir)
            output = root / "refs.md"

            result = subprocess.run(
                [
                    sys.executable,
                    str(SCRIPT_PATH),
                    "--head-sha",
                    sha,
                    "--manifest-file",
                    str(root / "missing.json"),
                    "--output",
                    str(output),
                ],
                text=True,
                capture_output=True,
                check=False,
            )

            self.assertEqual(result.returncode, 0, msg=result.stderr)
            content = output.read_text(encoding="utf-8")
            self.assertIn("Evidence state: `unavailable`", content)


if __name__ == "__main__":
    unittest.main()
