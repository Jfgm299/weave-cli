# Weave CLI — Agent Config

## Stack
Go · Cobra · Viper · Bubble Tea (future)

## Critical Rules
- **Repository scope only** — commands and changes MUST stay inside the current repo
- **Never edit `.env` files** — use `.env.example` as reference only
- **Never commit without explicit user instruction** — stage only, wait for `/commit`
- **Never add Co-Authored-By** to commits
- **Never build after changes**
- **All documentation and command content MUST be in English**

## Branch & Commit Conventions
- Branch pattern: `(feat|fix|chore|docs|style|refactor|perf|test|build|ci|revert)/<name>`
- Conventional commits: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`, `perf:`, `style:`, `build:`, `ci:`, `revert:`
- Never commit directly to `main` or `develop`

## Config Structure

`.agents/` is the single source of truth for all AI agent config.
`.claude/` and `.opencode/` reference it via symlinks — never create config files directly in those directories.

### Current symlink map

| `.agents/` path | `.claude/` symlink | `.opencode/` symlink | Type |
|----------------|-------------------|---------------------|------|
| `AGENTS.md` | `CLAUDE.md` → `../.agents/AGENTS.md` | `AGENTS.md` → `../.agents/AGENTS.md` | file |
| `commands/` | `commands` → `../.agents/commands` | `commands` → `../.agents/commands` | **directory** |
| `docs/` | `docs` → `../.agents/docs` | `docs` → `../.agents/docs` | **directory** |
| `dev/` | — (not exposed, accessed by path) | — | none |
| `skills/` | — (not yet symlinked) | — | none |

### Rules for adding new content

**File inside an already-symlinked directory** (`commands/`, `docs/`) → **No extra work needed.**
The directory symlink covers all files inside it automatically.

```
# Example: adding a new command
touch .agents/commands/new-command.md
# → automatically available as .claude/commands/new-command.md and .opencode/commands/new-command.md
# → no symlinks required
```

**New directory under `.agents/`** → **Add directory symlinks from both tool dirs.**

```bash
# Example: adding a new top-level directory
mkdir .agents/new-dir/
ln -s ../.agents/new-dir .claude/new-dir
ln -s ../.agents/new-dir .opencode/new-dir
# Then register it in this table above
```

**New file at `.agents/` root** → **Add individual file symlinks.**

```bash
ln -s ../.agents/new-file.md .claude/new-file.md
ln -s ../.agents/new-file.md .opencode/new-file.md
```

### File locations

| Type | Location | Needs symlinks? |
|------|----------|----------------|
| Commands | `.agents/commands/<name>.md` | No — `commands/` dir already symlinked |
| Docs | `.agents/docs/<name>.md` | No — `docs/` dir already symlinked |
| Skills | `.agents/skills/<name>/SKILL.md` | No — skills accessed by path from AGENTS.md |
| PRDs | `.agents/dev/prd-<feature-name>.md` | No — accessed by path, not via tool dirs |

---

## Documentation

Load docs **on demand only** — do NOT pre-load all files at startup. Read a doc only when the task explicitly requires it.

| Doc | Path | Load when... |
|-----|------|--------------|
| Architecture | `@docs/architecture.md` | Startup sequence, schemas, docker setup |
| Module System | `@docs/module-system.md` | manifest.py, autodiscovery, automation contract spec |
| Patterns | `@docs/patterns.md` | Writing models, services, routers, exception handlers |
| Database | `@docs/database.md` | Alembic migrations, multi-schema, naming conventions |
| Testing | `@docs/testing.md` | conftest hierarchy, fixtures, how to run tests |
| Module Registry | `@docs/modules/README.md` | Overview of all modules and their status |
| gym_tracker | `@docs/modules/gym_tracker.md` | Working on gym_tracker module |
| expenses_tracker | `@docs/modules/expenses_tracker.md` | Working on expenses_tracker module |
| macro_tracker | `@docs/modules/macro_tracker.md` | Working on macro_tracker module |
| flights_tracker | `@docs/modules/flights_tracker.md` | Working on flights_tracker module |
| travels_tracker | `@docs/modules/travels_tracker.md` | Working on travels_tracker module |
| calendar_tracker | `@docs/modules/calendar_tracker.md` | Working on calendar_tracker module |
| automations_engine | `@docs/modules/automations_engine.md` | Working on the automation engine |

---

## Commands

| Command | File | What it does |
|---------|------|--------------|
| `/commit` | `@commands/commit.md` | Stage + commit with conventional commits |
| `/pr` | `@commands/pr.md` | Create PR following project template |
| `/test` | `@commands/test.md` | Run pytest at the right scope |
| `/deploy-check` | `@commands/deploy-check.md` | Pre-deploy 7-step checklist |
| `/new-module` | `@commands/new-module.md` | Scaffold a new backend module |
| `/update-docs` | `@commands/update-docs.md` | Update .agents/docs/ after code changes |
| `/prompt` | `@commands/prompt.md` | Optimize a prompt for AI agents |
| `/new-prd` | `@commands/new-prd.md` | Create a PRD — mandatory first step before any substantial feature |

---

## Skills

When working on this project, load the relevant skill(s) BEFORE writing any code.

### How to Use
1. Check the trigger column to find skills that match your current task
2. Load the skill by reading the SKILL.md file at the listed path
3. Follow ALL patterns and rules from the loaded skill
4. Multiple skills can apply simultaneously

| Skill | Trigger | Path |
|-------|---------|------|
| `version-batch-planner` | New version planning, batch definition, PRD-to-OpenSpec traceability, deferred/debt discipline | `.agents/skills/version-batch-planner/SKILL.md` |

---

## Quick Reference

```bash
# Weave CLI (target commands)
weave forge
weave provider add claude-code
weave provider add opencode
weave skill add <skill-name>
weave command add <command-name>
weave doctor
```
