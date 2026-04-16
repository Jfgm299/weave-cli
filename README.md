# Weave CLI

Declarative project control plane for reproducible AI-agent setup.

## One-command bootstrap (from zero)

Run:

```sh
./scripts/bootstrap.sh
```

What it does:

1. Installs Go if missing (macOS via Homebrew, Linux via apt/dnf/pacman).
2. Runs `go mod tidy` and `go mod download`.
3. Runs `go test ./...` to validate the environment.

## Manual test command

```sh
go test ./...
```
