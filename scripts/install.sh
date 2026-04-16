#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "Installing weave binary from source..."
cd "$ROOT_DIR"
go install ./cmd/weave

echo "Install complete. Verify with: weave --version"
echo "Next steps: weave forge && weave doctor"
