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

## One-command install (Weave binary)

Run:

```sh
./scripts/install.sh
```

Then verify installation:

```sh
weave --version
weave forge
weave doctor
```

## Manual test command

```sh
go test ./...
```

## Quickstart

```sh
go run ./cmd/weave --help
go run ./cmd/weave forge
go run ./cmd/weave provider add claude-code
go run ./cmd/weave doctor
```

## Docs

- Install + troubleshooting: `docs/reference/install.md`
- Distribution baseline: `docs/reference/distribution.md`
- Migration guide: `docs/reference/migration.md`
- Release notes policy: `docs/reference/releases.md`
