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

The installer places `weave` in `~/.local/bin` (or `WEAVE_BIN_DIR` if set) and updates PATH configuration for bash, zsh, fish, and profile shells.
If your current session does not pick it up immediately, open a new shell.

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
go run ./cmd/weave provider add codex
go run ./cmd/weave command add pr-review
go run ./cmd/weave command add pr-review --provider codex
go run ./cmd/weave doctor
```

## CLI commands (current)

```sh
weave --help
weave --version

weave forge [--dry-run]

weave skill add <name> [--from <dir>] [--overwrite|--skip|--backup] [--dry-run]
weave command add <name> [--provider <name>] [--from <dir>] [--overwrite|--skip|--backup] [--dry-run]

weave provider add <name> [--dry-run]
weave provider remove <name> [--dry-run]
weave provider repair <name> [--dry-run]
weave provider list

weave doctor [--json]
weave migrate [--dry-run]
```

## v2 highlights

- **Provider support**: `claude-code`, `codex`, `opencode`.
- **Provider-aware command install**:
  - default: `weave command add <name>` installs shared command + provider projections for all enabled providers
  - exclusive: `weave command add <name> --provider <name>` installs only the target provider projection
- **Codex command wrapper namespace**:
  - `.codex/commands/__weave_commands__/<name>/SKILL.md`
- **No-provider behavior**:
  - interactive: prompt `No providers are currently enabled. Continue anyway? [y/N]:`
  - non-interactive: fail fast with actionable guidance
- **Transactional guarantees**:
  - if multi-provider apply fails, rollback is attempted for all planned shared/provider projection paths
  - `weave.yaml` persists only after successful filesystem apply

## Asset source and project layout

### Project root requirement

Mutating commands require a Git project root (`.git`) in the current directory or one of its parents.

### Canonical project directories

Within the target project, Weave manages canonical assets under:

- `.agents/skills`
- `.agents/commands`
- `.agents/docs`

Provider directories are treated as projections/symlinks from this canonical structure.

Current provider projections:

- `.claude/*` -> `../.agents/*`
- `.codex/*` -> `../.agents/*`
- `.opencode/*` -> `../.agents/*`

### Default shared source directories

By default, named asset resolution uses:

- Skills: `~/.weave/skills`
- Commands: `~/.weave/commands`

### `--from` option

Use `--from` to override source resolution for a single command invocation.

```sh
weave skill add sdd-orchestrator --from "$HOME/custom-skills"
weave command add pr-review --from "$HOME/custom-commands"
```

Source resolution precedence:

1. `--from`
2. Environment variables
3. `weave.yaml` configured sources
4. v1 defaults (`~/.weave/skills`, `~/.weave/commands`)

### Configure project-level default source directories

You can override default source paths in `weave.yaml`:

```yaml
sources:
  skills_dir: "/Users/you/custom-skills"
  commands_dir: "/Users/you/custom-commands"
```

After this, `weave skill add <name>` and `weave command add <name>` will resolve from those directories unless a higher-precedence source is provided (for example `--from`).

## Docs

- Install + troubleshooting: `docs/reference/install.md`
- Distribution baseline: `docs/reference/distribution.md`
- Migration guide: `docs/reference/migration.md`
- Release notes policy: `docs/reference/releases.md`
- Providers + provider-aware command behavior: `docs/reference/providers.md`
- Doctor checks + repair guidance: `docs/reference/doctor.md`
- Transaction semantics and rollback guarantees: `docs/reference/transactions.md`
- Config metadata (`provider_compat`, `shared_install`): `docs/reference/config.md`
