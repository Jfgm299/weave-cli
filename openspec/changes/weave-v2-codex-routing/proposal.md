# Proposal — weave-v2-codex-routing

## Context
Weave v1 delivered canonical `.agents` workflows and provider projections for `claude-code` and `opencode`, but left deferred debt and lacks Codex support. Command installation is not provider-aware.

## Goals
- Add `codex` provider with symlink-based projection.
- Make `weave command add <name>` provider-aware:
  - default: install for all enabled providers
  - `--provider <name>`: exclusive install for that provider (without `.agents`)
- Use Codex command wrapper namespace: `__weave_commands__`.
- Enforce all-or-nothing rollback across multi-provider command installs.
- Close v1 deferred debt (`DEF-001..DEF-004`) before implementation batches.

## Non-Goals
- Remote registries/marketplace.
- New copy-based sync mode (v1 stays symlink-first).
- UI/TUI changes.

## Decisions
- No providers enabled:
  - interactive: ask `No providers are currently enabled. Continue anyway? [y/N]`
  - non-interactive: fail automatically with actionable message.
- Keep user-facing command unchanged (`weave command add ...`), route internally by provider strategy.

## Risks
- Current architecture assumes `.agents` as canonical for command installs.
- Doctor/repair must understand shared and exclusive command installs.
- Config schema may need provider-target metadata for deterministic repair.
