#!/usr/bin/env python3
import json
import subprocess
import sys
import tempfile
import unittest
from pathlib import Path


SCRIPT_PATH = Path(__file__).with_name("check_pr_metadata_coherence.py")


def run_gate(payload: dict) -> subprocess.CompletedProcess[str]:
    with tempfile.NamedTemporaryFile("w", delete=False, suffix=".json") as tmp:
        json.dump(payload, tmp)
        event_path = tmp.name

    try:
        return subprocess.run(
            [sys.executable, str(SCRIPT_PATH), "--from-file", event_path],
            text=True,
            capture_output=True,
            check=False,
        )
    finally:
        Path(event_path).unlink(missing_ok=True)


def base_body(
    primary: str = "feature",
    breaking_checked: bool = False,
    closes: str = "#123",
    extra_type_rows: list[str] | None = None,
) -> str:
    primary_rows = [
        f"- [{'x' if primary == 'bug' else ' '}] `type:bug` — Bug fix",
        f"- [{'x' if primary == 'feature' else ' '}] `type:feature` — New feature",
        f"- [{'x' if primary == 'docs' else ' '}] `type:docs` — Documentation only",
        f"- [{'x' if primary == 'refactor' else ' '}] `type:refactor` — Refactor",
        f"- [{'x' if primary == 'chore' else ' '}] `type:chore` — Chore",
    ]
    breaking_row = f"- [{'x' if breaking_checked else ' '}] `type:breaking-change` — Breaking change"
    return (
        "## 🔗 Linked Issue\n\n"
        f"Closes {closes}\n\n"
        "## 🏷️ PR Type\n\n"
        + "\n".join(primary_rows)
        + "\n"
        + breaking_row
        + "\n"
        + ("\n" + "\n".join(extra_type_rows) if extra_type_rows else "")
    )


class PRMetadataCoherenceTests(unittest.TestCase):
    def test_passes_with_matching_primary_type_and_issue_reference(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}],
                "body": base_body(primary="feature", closes="#42"),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 0)
        self.assertIn("passed", result.stdout)

    def test_passes_when_issue_is_na(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:docs"}],
                "body": base_body(primary="docs", closes="N/A"),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 0)

    def test_passes_when_breaking_change_checkbox_and_label_match(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}, {"name": "type:breaking-change"}],
                "body": base_body(
                    primary="feature", breaking_checked=True, closes="#77"
                ),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 0)

    def test_fails_when_primary_checklist_does_not_match_label(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}],
                "body": base_body(primary="bug", closes="#123"),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 1)
        self.assertIn("does not match", result.stderr)

    def test_fails_when_breaking_change_checkbox_and_label_do_not_match(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}],
                "body": base_body(
                    primary="feature", breaking_checked=True, closes="#123"
                ),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 1)
        self.assertIn(
            "type:breaking-change checklist selection must match", result.stderr
        )

    def test_fails_when_multiple_closes_declarations_exist(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}],
                "body": base_body(primary="feature", closes="#123") + "\nCloses #456\n",
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 1)
        self.assertIn("exactly one 'Closes ...' declaration", result.stderr)

    def test_fails_when_closes_line_is_not_exact_issue_reference(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}],
                "body": base_body(primary="feature", closes="owner/repo#321"),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 1)
        self.assertIn("must be exactly 'Closes #<id>' or 'Closes N/A'", result.stderr)

    def test_fails_when_unsupported_checklist_type_is_selected(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}],
                "body": base_body(
                    primary="feature",
                    closes="#123",
                    extra_type_rows=["- [x] `type:style` — Style changes"],
                ),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 1)
        self.assertIn("unsupported type entries", result.stderr)

    def test_fails_when_issue_line_is_missing_or_malformed(self) -> None:
        payload = {
            "pull_request": {
                "labels": [{"name": "type:feature"}],
                "body": base_body(primary="feature", closes="pending"),
            }
        }
        result = run_gate(payload)
        self.assertEqual(result.returncode, 1)
        self.assertIn("Linked Issue", result.stderr)


if __name__ == "__main__":
    unittest.main()
