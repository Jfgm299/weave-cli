#!/usr/bin/env python3
import unittest
from unittest.mock import patch
from pathlib import Path
import sys

sys.path.insert(0, str(Path(__file__).parent))

from post_manifest_pr_comment import MARKER, build_comment_body, upsert_comment


class PostManifestPRCommentTests(unittest.TestCase):
    def test_build_comment_body_contains_marker_sha_and_manifest(self) -> None:
        body = build_comment_body(
            "openspec/evidence/" + "a" * 40 + ".json",
            "a" * 40,
            "| Workflow | Status |\n|---|---|\n| `x` | `success` |",
        )

        self.assertIn(MARKER, body)
        self.assertIn("`" + "a" * 40 + "`", body)
        self.assertIn("openspec/evidence/" + "a" * 40 + ".json", body)

    @patch("post_manifest_pr_comment.api_request")
    def test_upsert_comment_updates_existing_marker_comment(self, api_request) -> None:
        api_request.side_effect = [
            [{"id": 101, "body": f"{MARKER}\nold"}],
            {"html_url": "https://example.com/updated"},
        ]

        url = upsert_comment("Jfgm299/weave-cli", 5, "token", "new body")

        self.assertEqual(url, "https://example.com/updated")
        self.assertEqual(api_request.call_args_list[1].args[0], "PATCH")


if __name__ == "__main__":
    unittest.main()
