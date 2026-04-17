#!/usr/bin/env python3
import unittest
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).parent))

from collect_workflow_run_evidence import (
    build_failure_payload,
    download_artifact_payload,
    select_latest_completed_runs,
    validate_downloaded_payload,
)


class CollectWorkflowRunEvidenceTests(unittest.TestCase):
    def test_select_latest_completed_runs_by_attempt_then_id(self) -> None:
        sha = "d" * 40
        runs = [
            {
                "name": "release-artifacts",
                "head_sha": sha,
                "status": "completed",
                "run_attempt": 1,
                "id": 10,
            },
            {
                "name": "release-artifacts",
                "head_sha": sha,
                "status": "completed",
                "run_attempt": 2,
                "id": 9,
            },
            {
                "name": "release-artifacts",
                "head_sha": sha,
                "status": "in_progress",
                "run_attempt": 99,
                "id": 99,
            },
            {
                "name": "migration-note-gate",
                "head_sha": sha,
                "status": "completed",
                "run_attempt": 1,
                "id": 3,
            },
            {
                "name": "migration-note-gate",
                "head_sha": "e" * 40,
                "status": "completed",
                "run_attempt": 99,
                "id": 99,
            },
        ]

        selected = select_latest_completed_runs(runs, sha)

        self.assertEqual(selected["release-artifacts"]["run_attempt"], 2)
        self.assertEqual(selected["release-artifacts"]["id"], 9)
        self.assertEqual(selected["migration-note-gate"]["head_sha"], sha)

    def test_validate_downloaded_payload_rejects_malformed(self) -> None:
        sha = "f" * 40
        malformed = {
            "schema_version": 1,
            "workflow_name": "release-artifacts",
            "head_sha": sha,
        }

        with self.assertRaisesRegex(ValueError, "missing required key"):
            validate_downloaded_payload(malformed, "release-artifacts", sha)

    def test_failure_payload_contains_deterministic_failure_shape(self) -> None:
        sha = "a" * 40
        payload = build_failure_payload("repo-hygiene-gate", sha, "missing artifact")

        self.assertEqual(payload["schema_version"], 1)
        self.assertEqual(payload["workflow_name"], "repo-hygiene-gate")
        self.assertEqual(payload["head_sha"], sha)
        self.assertEqual(payload["status"], "failure")
        self.assertIn("aggregator:", payload["checks"]["summary"])

    def test_download_artifact_payload_uses_no_authorization_header(self) -> None:
        captured = {}

        class FakeResponse:
            def __enter__(self):
                return self

            def __exit__(self, exc_type, exc, tb):
                return False

            def read(self):
                import io
                import json
                import zipfile

                buf = io.BytesIO()
                with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
                    zf.writestr(
                        "artifact.json",
                        json.dumps(
                            {
                                "schema_version": 1,
                                "workflow_name": "release-artifacts",
                                "run_id": "1",
                                "run_url": "https://example.test",
                                "head_sha": "a" * 40,
                                "status": "success",
                                "timestamp": "1970-01-01T00:00:00+00:00",
                                "checks": {"summary": "ok"},
                            }
                        ),
                    )
                return buf.getvalue()

        def fake_urlopen(req):
            captured["authorization"] = req.headers.get("Authorization")
            return FakeResponse()

        import urllib.request

        original = urllib.request.urlopen
        urllib.request.urlopen = fake_urlopen
        try:
            _ = download_artifact_payload("https://example.test/archive.zip", "token")
        finally:
            urllib.request.urlopen = original

        self.assertIsNone(captured["authorization"])


if __name__ == "__main__":
    unittest.main()
