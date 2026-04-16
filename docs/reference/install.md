# Installation

## Quickstart

1. Install Go 1.24+.
2. Clone the repository.
3. Run:

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
- **Config validation fails with outdated schema**
  - Run `weave migrate` (or `weave migrate --dry-run` first).
- **Doctor reports stale provider integration**
  - Run `weave provider repair <provider>`.
