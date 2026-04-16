#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="${WEAVE_BIN_DIR:-$HOME/.local/bin}"

append_path_block_if_missing() {
  local rc_file="$1"
  local target_dir="$2"

  mkdir -p "$(dirname "$rc_file")"
  touch "$rc_file"

  if grep -Fq "# >>> weave-cli path >>>" "$rc_file"; then
    return
  fi

  cat <<EOF >> "$rc_file"

# >>> weave-cli path >>>
if [ -d "$target_dir" ] && [[ ":\$PATH:" != *":$target_dir:"* ]]; then
  export PATH="$target_dir:\$PATH"
fi
# <<< weave-cli path <<<
EOF
}

configure_fish_path() {
  local target_dir="$1"
  local fish_conf_dir="$HOME/.config/fish/conf.d"
  local fish_conf_file="$fish_conf_dir/weave_path.fish"

  mkdir -p "$fish_conf_dir"

  cat <<EOF > "$fish_conf_file"
if test -d "$target_dir"
    if not contains "$target_dir" \$PATH
        set -gx PATH "$target_dir" \$PATH
    end
end
EOF
}

echo "Installing weave binary from source..."
cd "$ROOT_DIR"
mkdir -p "$BIN_DIR"
GOBIN="$BIN_DIR" go install ./cmd/weave

append_path_block_if_missing "$HOME/.profile" "$BIN_DIR"
append_path_block_if_missing "$HOME/.bashrc" "$BIN_DIR"
append_path_block_if_missing "$HOME/.zshrc" "$BIN_DIR"
configure_fish_path "$BIN_DIR"

if [[ ":$PATH:" != *":$BIN_DIR:"* ]]; then
  echo ""
  echo "Installed binary at: $BIN_DIR/weave"
  echo "PATH updates were written to: ~/.profile, ~/.bashrc, ~/.zshrc, and ~/.config/fish/conf.d/weave_path.fish"
  echo "Open a new shell (or source your shell config) before running: weave --version"
else
  echo "Installed binary at: $BIN_DIR/weave"
fi

echo "Install complete. Verify with: weave --version"
echo "Next steps: weave forge && weave doctor"
