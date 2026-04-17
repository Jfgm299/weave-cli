#!/usr/bin/env python3
import argparse
import json
import re
import subprocess
import sys
from pathlib import Path


SEMVER_RE = re.compile(r"^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$")
SEMVER_LABEL_RE = re.compile(r"^semver:(patch|minor|major)$")
RATIONALE_LINE_RE = re.compile(r"^\s*Semver rationale\s*:\s*(.+)\s*$", re.IGNORECASE)
GO_VERSION_RE = re.compile(r'const\s+version\s*=\s*"([0-9]+\.[0-9]+\.[0-9]+)"')
TEST_VERSION_RE = re.compile(r'"([0-9]+\.[0-9]+\.[0-9]+)"')


def parse_semver(value: str) -> tuple[int, int, int]:
    match = SEMVER_RE.fullmatch(value)
    if not match:
        raise ValueError(f"invalid semver '{value}'")
    return tuple(int(part) for part in match.groups())


def classify_bump(previous: str, current: str) -> str:
    prev = parse_semver(previous)
    curr = parse_semver(current)

    if curr == prev:
        raise ValueError("version did not change")

    if curr[0] == prev[0] and curr[1] == prev[1] and curr[2] == prev[2] + 1:
        return "patch"
    if curr[0] == prev[0] and curr[1] == prev[1] + 1 and curr[2] == 0:
        return "minor"
    if curr[0] == prev[0] + 1 and curr[1] == 0 and curr[2] == 0:
        return "major"

    raise ValueError(
        f"non-compliant semver transition {previous} -> {current} (expected patch/minor/major canonical increment)"
    )


def parse_version_from_go(path: Path) -> str:
    content = path.read_text(encoding="utf-8")
    match = GO_VERSION_RE.search(content)
    if not match:
        raise ValueError(f"unable to extract const version from {path}")
    return match.group(1)


def parse_expected_version_from_test(path: Path) -> str:
    content = path.read_text(encoding="utf-8")
    matches = TEST_VERSION_RE.findall(content)
    if not matches:
        raise ValueError(f"unable to extract semver assertion from {path}")
    return matches[0]


def read_changed_files(
    from_file: str | None, base_sha: str | None, head_sha: str | None
) -> list[str]:
    if from_file:
        return [
            line.strip()
            for line in Path(from_file).read_text(encoding="utf-8").splitlines()
            if line.strip()
        ]

    if not base_sha or not head_sha:
        return []

    result = subprocess.run(
        ["git", "diff", "--name-only", base_sha, head_sha],
        text=True,
        capture_output=True,
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError(f"git diff failed: {result.stderr.strip()}")
    return [line.strip() for line in result.stdout.splitlines() if line.strip()]


def resolve_declared_bump(labels: list[str], explicit: str | None) -> str | None:
    if explicit:
        return explicit

    selected = []
    for label in labels:
        match = SEMVER_LABEL_RE.fullmatch(label)
        if match:
            selected.append(match.group(1))

    if not selected:
        return None
    if len(selected) > 1:
        raise ValueError(
            "release metadata must include exactly one semver:* label (semver:patch|semver:minor|semver:major)"
        )
    return selected[0]


def resolve_rationale(body: str, explicit: str | None) -> str | None:
    if explicit:
        return explicit.strip()

    for line in body.splitlines():
        match = RATIONALE_LINE_RE.match(line)
        if match:
            return match.group(1).strip()
    return None


def load_event(path: str | None) -> dict:
    if not path:
        return {}
    return json.loads(Path(path).read_text(encoding="utf-8"))


def read_version_from_git(base_sha: str, version_file: str) -> str:
    result = subprocess.run(
        ["git", "show", f"{base_sha}:{version_file}"],
        text=True,
        capture_output=True,
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError(
            f"unable to read previous version from {base_sha}:{version_file}: {result.stderr.strip()}"
        )
    match = GO_VERSION_RE.search(result.stdout)
    if not match:
        raise ValueError(
            f"unable to extract const version from {base_sha}:{version_file}"
        )
    return match.group(1)


def validate(args: argparse.Namespace) -> list[str]:
    event = load_event(args.event_file)
    pr = event.get("pull_request") or {}

    base_sha = args.base_sha or (pr.get("base") or {}).get("sha")
    head_sha = args.head_sha or (pr.get("head") or {}).get("sha")

    labels = [
        label.get("name", "") for label in pr.get("labels", []) if label.get("name")
    ]
    body = pr.get("body") or ""

    changed = read_changed_files(args.changed_files_file, base_sha, head_sha)
    version_file = args.version_file
    version_test_file = args.version_test_file

    version_changed = version_file in changed
    test_changed = version_test_file in changed

    errors: list[str] = []

    if version_changed != test_changed:
        errors.append(
            f"version governance violation: '{version_file}' and '{version_test_file}' must be updated together"
        )

    if not version_changed and not test_changed:
        if args.declared_bump or args.bump_rationale:
            errors.append(
                "release metadata declares semver bump, but neither version source nor version test changed"
            )
        return errors

    try:
        current_version = args.current_version or parse_version_from_go(
            Path(version_file)
        )
    except Exception as exc:
        errors.append(str(exc))
        return errors

    try:
        asserted_test_version = parse_expected_version_from_test(
            Path(version_test_file)
        )
        if asserted_test_version != current_version:
            errors.append(
                f"version test mismatch: {version_test_file} asserts {asserted_test_version} but {version_file} is {current_version}"
            )
    except Exception as exc:
        errors.append(str(exc))

    previous_version = args.base_version
    if not previous_version:
        if not base_sha:
            errors.append(
                "unable to determine base version (missing base SHA); provide --base-version or --base-sha"
            )
            return errors
        try:
            previous_version = read_version_from_git(base_sha, version_file)
        except Exception as exc:
            errors.append(str(exc))
            return errors

    try:
        computed_bump = classify_bump(previous_version, current_version)
    except Exception as exc:
        errors.append(str(exc))
        return errors

    try:
        declared_bump = resolve_declared_bump(labels, args.declared_bump)
    except Exception as exc:
        errors.append(str(exc))
        return errors

    rationale = resolve_rationale(body, args.bump_rationale)

    if not declared_bump:
        errors.append(
            "missing semver declaration: add exactly one label semver:patch|semver:minor|semver:major"
        )
    elif declared_bump != computed_bump:
        errors.append(
            f"declared bump semver:{declared_bump} does not match computed bump '{computed_bump}' from {previous_version} -> {current_version}"
        )

    if not rationale:
        errors.append(
            "missing semver rationale: include 'Semver rationale: <why>' in PR body or pass --bump-rationale"
        )
    elif rationale.strip().upper() == "N/A":
        errors.append(
            "semver rationale cannot be N/A for release-intended version bump"
        )

    return errors


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate semver bump governance and version source/test synchronization"
    )
    parser.add_argument("--event-file", help="Path to GitHub event payload JSON")
    parser.add_argument("--base-sha")
    parser.add_argument("--head-sha")
    parser.add_argument("--changed-files-file", help="Optional list of changed files")
    parser.add_argument("--version-file", default="internal/cli/version.go")
    parser.add_argument("--version-test-file", default="internal/cli/version_test.go")
    parser.add_argument("--declared-bump", choices=["patch", "minor", "major"])
    parser.add_argument("--bump-rationale")
    parser.add_argument("--base-version")
    parser.add_argument("--current-version")
    args = parser.parse_args()

    try:
        errors = validate(args)
    except Exception as exc:
        print(f"Version governance check failed to execute: {exc}", file=sys.stderr)
        return 1

    if errors:
        print("Version governance gate failed:", file=sys.stderr)
        for error in errors:
            print(f" - {error}", file=sys.stderr)
        return 1

    print("Version governance gate passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
