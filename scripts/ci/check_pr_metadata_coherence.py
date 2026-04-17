#!/usr/bin/env python3
import argparse
import json
import re
import sys
from pathlib import Path


PR_TYPE_RE = re.compile(r"- \[([ xX])\] `type:([a-z0-9-]+)`")
ISSUE_LINE_RE = re.compile(r"^\s*Closes\s+(.+)\s*$", re.IGNORECASE | re.MULTILINE)
ISSUE_REF_RE = re.compile(r"#\d+")


def parse_selected_types(body: str) -> list[str]:
    selected = []
    for mark, label_type in PR_TYPE_RE.findall(body):
        if mark.lower() == "x":
            selected.append(label_type)
    return selected


def extract_issue_line(body: str) -> str | None:
    match = ISSUE_LINE_RE.search(body)
    if not match:
        return None
    return match.group(1).strip()


def normalize_labels(event: dict) -> set[str]:
    pr = event.get("pull_request") or {}
    labels = {label.get("name", "") for label in pr.get("labels", [])}
    return {label for label in labels if label}


def validate(event: dict) -> list[str]:
    pr = event.get("pull_request") or {}
    body = pr.get("body") or ""
    labels = normalize_labels(event)

    errors: list[str] = []

    selected_types = parse_selected_types(body)
    primary_types = [value for value in selected_types if value != "breaking-change"]
    selected_breaking = "breaking-change" in selected_types

    primary_labels = sorted(
        [
            label
            for label in labels
            if label.startswith("type:") and label != "type:breaking-change"
        ]
    )

    if len(primary_types) != 1:
        errors.append(
            "PR checklist must select exactly one primary type:* item (excluding type:breaking-change)"
        )
    if len(primary_labels) != 1:
        errors.append(
            "PR labels must include exactly one primary type:* label (excluding type:breaking-change)"
        )

    if len(primary_types) == 1 and len(primary_labels) == 1:
        expected = f"type:{primary_types[0]}"
        if primary_labels[0] != expected:
            errors.append(
                f"Primary checklist type '{expected}' does not match label '{primary_labels[0]}'"
            )

    has_breaking_label = "type:breaking-change" in labels
    if selected_breaking != has_breaking_label:
        errors.append(
            "type:breaking-change checklist selection must match presence of 'type:breaking-change' label"
        )

    issue_line = extract_issue_line(body)
    if issue_line is None:
        errors.append(
            "PR body must include a 'Closes ...' line in Linked Issue section"
        )
    else:
        normalized = issue_line.strip()
        if normalized.upper() == "N/A":
            pass
        elif not ISSUE_REF_RE.search(normalized):
            errors.append(
                "Linked Issue must be 'N/A' or include at least one issue reference like #123"
            )

    return errors


def load_event(from_file: str | None) -> dict:
    if not from_file:
        raise ValueError("--from-file is required")

    path = Path(from_file)
    return json.loads(path.read_text(encoding="utf-8"))


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate PR checklist-label-issue coherence"
    )
    parser.add_argument("--from-file", required=True)
    args = parser.parse_args()

    try:
        event = load_event(args.from_file)
        errors = validate(event)
    except Exception as exc:
        print(f"PR metadata coherence check failed to execute: {exc}", file=sys.stderr)
        return 1

    if errors:
        print("PR metadata coherence gate failed:", file=sys.stderr)
        for error in errors:
            print(f" - {error}", file=sys.stderr)
        return 1

    print("PR metadata coherence gate passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
