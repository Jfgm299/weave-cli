# Design — weave-v2-codex-routing

## Architecture
Introduce provider-aware command projection strategy:

- Shared mode (default command add):
  - source from configured command catalog
  - project to all enabled providers using per-provider strategy
- Exclusive mode (`--provider`):
  - bypass shared `.agents` install requirement
  - install only target provider projection

## Provider command projection strategies
- `claude-code`: markdown command projection
- `opencode`: markdown command projection
- `codex`: skill-wrapper projection under `__weave_commands__/<command>/SKILL.md`

## Transaction model
One logical transaction per `command add`:
1. plan all operations
2. apply operations
3. persist config
4. on any failure, rollback all applied operations and skip config mutation

## Prompt behavior
If no enabled providers:
- interactive: `[y/N]` prompt
- non-interactive: immediate fail

## Doctor/repair model
Doctor must validate:
- shared projection assets
- exclusive provider installs
- codex wrapper namespace integrity
Repair commands must stay actionable and deterministic.
