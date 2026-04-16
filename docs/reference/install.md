# Installation

## One-command install

1. Install Go 1.22+.
2. Clone the repository.
3. Run:

```sh
./scripts/install.sh
```

## Verify installation

```sh
weave --version
weave forge
weave doctor
```

## Contributor quickstart (run without installing)

```sh
go run ./cmd/weave --help
go run ./cmd/weave forge
go run ./cmd/weave provider add claude-code
go run ./cmd/weave doctor
```

## Troubleshooting

- **Provider add fails with missing binary**
  - Install required provider binaries (`claude`, `opencode`).
  - Re-run `weave provider repair <provider>`.
- **`weave` command is not found after install**
  - Ensure `$(go env GOPATH)/bin` (or `~/go/bin`) is in your `PATH`.
  - Open a new shell and run `weave --version` again.
- **Config validation fails with outdated schema**
  - Run `weave migrate` (or `weave migrate --dry-run` first).
- **Doctor reports stale provider integration**
  - Run `weave provider repair <provider>`.
