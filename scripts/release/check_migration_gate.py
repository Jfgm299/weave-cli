#!/usr/bin/env python3
import json
import os
import re
import sys


MIGRATION_HEADING_RE = re.compile(
    r"^\s{0,3}#{2,6}\s+migration\b", re.IGNORECASE | re.MULTILINE
)


def has_migration_section(body: str) -> bool:
    return bool(MIGRATION_HEADING_RE.search(body))


def main() -> int:
    event_path = os.environ.get("GITHUB_EVENT_PATH")
    if not event_path:
        print("GITHUB_EVENT_PATH is not set; skipping migration gate")
        return 0

    with open(event_path, "r", encoding="utf-8") as f:
        event = json.load(f)

    pr = event.get("pull_request") or {}
    labels = [l.get("name", "") for l in pr.get("labels", [])]
    has_breaking = "breaking-change" in labels

    if not has_breaking:
        print("No breaking-change label present; migration gate passed")
        return 0

    body = pr.get("body") or ""
    if has_migration_section(body):
        print("Migration section found in PR body; migration gate passed")
        return 0

    print(
        "Migration gate failed: PR has label 'breaking-change' but body is missing a '## Migration' section.",
        file=sys.stderr,
    )
    return 1


if __name__ == "__main__":
    raise SystemExit(main())
