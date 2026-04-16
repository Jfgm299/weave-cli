#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

install_go_macos() {
  if command -v go >/dev/null 2>&1; then
    echo "Go already installed: $(go version)"
    return
  fi

  if ! command -v brew >/dev/null 2>&1; then
    echo "Homebrew is required on macOS to install Go."
    echo "Install Homebrew first: https://brew.sh"
    exit 1
  fi

  echo "Installing Go with Homebrew..."
  brew install go
}

install_go_linux() {
  if command -v go >/dev/null 2>&1; then
    echo "Go already installed: $(go version)"
    return
  fi

  if command -v apt-get >/dev/null 2>&1; then
    echo "Installing Go with apt..."
    sudo apt-get update
    sudo apt-get install -y golang-go
    return
  fi

  if command -v dnf >/dev/null 2>&1; then
    echo "Installing Go with dnf..."
    sudo dnf install -y golang
    return
  fi

  if command -v pacman >/dev/null 2>&1; then
    echo "Installing Go with pacman..."
    sudo pacman -Sy --noconfirm go
    return
  fi

  echo "Unsupported Linux package manager. Install Go manually: https://go.dev/doc/install"
  exit 1
}

OS="$(uname -s)"
case "$OS" in
  Darwin)
    install_go_macos
    ;;
  Linux)
    install_go_linux
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

if ! command -v go >/dev/null 2>&1; then
  echo "Go installation failed or Go is not in PATH."
  exit 1
fi

echo "Using $(go version)"

cd "$ROOT_DIR"

echo "Syncing Go modules..."
go mod tidy
go mod download

echo "Running test suite..."
go test ./...

echo "Bootstrap complete."
