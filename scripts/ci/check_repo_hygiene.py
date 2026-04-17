#!/usr/bin/env python3
import argparse
import subprocess
import sys
from pathlib import Path


def is_forbidden_path(path: str) -> bool:
    normalized = path.replace("\\", "/")
    return "/__pycache__/" in f"/{normalized}" or normalized.endswith(
        (".pyc", ".pyo", ".pyd")
    )


def collect_tracked_paths(from_file: str | None = None) -> list[str]:
    if from_file:
        return [
            line.strip()
            for line in Path(from_file).read_text(encoding="utf-8").splitlines()
            if line.strip()
        ]

    result = subprocess.run(
        ["git", "ls-files"],
        text=True,
        capture_output=True,
        check=False,
    )
    if result.returncode != 0:
        raise RuntimeError(f"git ls-files failed: {result.stderr.strip()}")

    return [line.strip() for line in result.stdout.splitlines() if line.strip()]


def find_forbidden(paths: list[str]) -> list[str]:
    return sorted([path for path in paths if is_forbidden_path(path)])


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Validate repository hygiene for Python cache artifacts"
    )
    parser.add_argument(
        "--from-file", help="Read tracked file list from file (one path per line)"
    )
    args = parser.parse_args()

    try:
        tracked = collect_tracked_paths(args.from_file)
    except Exception as exc:  # pragma: no cover
        print(f"Hygiene check failed to collect tracked files: {exc}", file=sys.stderr)
        return 1

    forbidden = find_forbidden(tracked)
    if not forbidden:
        print("Repository hygiene passed: no tracked Python cache artifacts found")
        return 0

    print(
        "Repository hygiene failed: tracked Python cache artifacts detected",
        file=sys.stderr,
    )
    for path in forbidden:
        print(f" - {path}", file=sys.stderr)
    return 1


if __name__ == "__main__":
    raise SystemExit(main())
