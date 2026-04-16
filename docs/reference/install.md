# Installation

## One-command install

1. Install Go 1.22+.
2. Clone the repository.
3. Run:

```sh
./scripts/install.sh
```

The installer places `weave` in `~/.local/bin` (or `WEAVE_BIN_DIR` if set) and writes PATH configuration for:

- `~/.profile`
- `~/.bashrc`
- `~/.zshrc`
- `~/.config/fish/conf.d/weave_path.fish`

If your current shell session does not see `weave` immediately, open a new shell.

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

- **`project root not detected` when running `weave forge`**
  - Current behavior requires a Git repository root (`.git`) to be present.
  - In interactive mode, Weave prompts to run `git init` automatically.
  - In non-interactive mode (CI/scripts), Weave fails fast with actionable guidance to run `git init`.

- **Provider add fails with missing binary**
  - Install required provider binaries (`claude`, `codex`, `opencode`).
  - Re-run `weave provider repair <provider>`.
- **`weave` command is not found after install**
  - The installer writes PATH entries automatically for bash/zsh/fish/profile.
  - Open a new shell and run `weave --version` again.
  - If you use a custom shell config, add `~/.local/bin` (or your `WEAVE_BIN_DIR`) to PATH manually.
- **Config validation fails with outdated schema**
  - Run `weave migrate` (or `weave migrate --dry-run` first).
- **Doctor reports stale provider integration**
  - Run `weave provider repair <provider>`.
