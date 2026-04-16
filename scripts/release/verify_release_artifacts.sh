#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DIST_DIR="${DIST_DIR:-$ROOT_DIR/dist}"

if [[ ! -f "$DIST_DIR/checksums.txt" ]]; then
  echo "missing checksums.txt in $DIST_DIR" >&2
  exit 1
fi

if [[ ! -f "$DIST_DIR/checksums.txt.sig" ]]; then
  echo "missing checksums.txt.sig in $DIST_DIR" >&2
  exit 1
fi

(
  cd "$DIST_DIR"
  shasum -a 256 -c checksums.txt
)

GNUPGHOME="$(mktemp -d)"
export GNUPGHOME

if [[ -n "${WEAVE_SIGNING_PUBLIC_KEY:-}" ]]; then
  printf '%s' "$WEAVE_SIGNING_PUBLIC_KEY" | gpg --batch --import
elif [[ -f "$DIST_DIR/checksums.pub" ]]; then
  gpg --batch --import "$DIST_DIR/checksums.pub"
else
  echo "No signing public key available (set WEAVE_SIGNING_PUBLIC_KEY or provide checksums.pub)" >&2
  rm -rf "$GNUPGHOME"
  exit 1
fi

gpg --batch --verify "$DIST_DIR/checksums.txt.sig" "$DIST_DIR/checksums.txt"
rm -rf "$GNUPGHOME"

echo "Release artifact verification passed"
