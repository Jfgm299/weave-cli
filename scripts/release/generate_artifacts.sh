#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
DIST_DIR="${DIST_DIR:-$ROOT_DIR/dist}"
VERSION="${WEAVE_VERSION:-dev}"
TARGETS="${WEAVE_TARGETS:-darwin/arm64 linux/amd64}"
REQUIRE_SIGNING_KEY="${REQUIRE_SIGNING_KEY:-0}"

mkdir -p "$DIST_DIR"
rm -f "$DIST_DIR"/checksums.txt "$DIST_DIR"/checksums.txt.sig

for target in $TARGETS; do
  GOOS="${target%/*}"
  GOARCH="${target#*/}"
  ARTIFACT="weave_${VERSION}_${GOOS}_${GOARCH}.tar.gz"
  BUILD_DIR="$(mktemp -d)"

  GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -o "$BUILD_DIR/weave" ./cmd/weave
  tar -C "$BUILD_DIR" -czf "$DIST_DIR/$ARTIFACT" weave
  rm -rf "$BUILD_DIR"
done

(
  cd "$DIST_DIR"
  shasum -a 256 weave_*.tar.gz > checksums.txt
)

if [[ -n "${WEAVE_SIGNING_PRIVATE_KEY:-}" ]]; then
  GNUPGHOME="$(mktemp -d)"
  export GNUPGHOME
  printf '%s' "$WEAVE_SIGNING_PRIVATE_KEY" | gpg --batch --import
  gpg --batch --yes --armor --output "$DIST_DIR/checksums.txt.sig" --detach-sign "$DIST_DIR/checksums.txt"
  rm -rf "$GNUPGHOME"
elif [[ "$REQUIRE_SIGNING_KEY" == "1" ]]; then
  echo "WEAVE_SIGNING_PRIVATE_KEY is required when REQUIRE_SIGNING_KEY=1" >&2
  exit 1
else
  GNUPGHOME="$(mktemp -d)"
  export GNUPGHOME
  cat > "$GNUPGHOME/genkey" <<'EOF'
%no-protection
Key-Type: RSA
Key-Length: 3072
Subkey-Type: RSA
Subkey-Length: 3072
Name-Real: Weave CI Signing
Name-Email: ci@weave.local
Expire-Date: 0
EOF
  gpg --batch --generate-key "$GNUPGHOME/genkey"
  gpg --batch --yes --armor --output "$DIST_DIR/checksums.txt.sig" --detach-sign "$DIST_DIR/checksums.txt"
  gpg --batch --armor --export > "$DIST_DIR/checksums.pub"
  rm -rf "$GNUPGHOME"
fi

echo "Artifacts generated in $DIST_DIR"
