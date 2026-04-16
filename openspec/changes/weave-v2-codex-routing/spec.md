# Spec — weave-v2-codex-routing

## Scope

This spec defines the v2 delta for:

1. Codex provider support.
2. Provider-aware command installation with unchanged user-facing CLI (`weave command add <name>`).
3. Explicit install-mode behavior for default multi-provider and `--provider` exclusive mode.
4. Strict all-or-nothing rollback semantics for multi-provider command operations.
5. Batch-0 closure of v1 deferred debt (`DEF-001..DEF-004`).

## Requirements

### Provider support

#### R-PROV-04
CLI MUST support provider `codex` with add/remove/repair/list parity relative to existing providers.

#### R-PROV-05
Provider adapters MUST declare command projection targets so command installation logic remains provider-agnostic in CLI handlers.

### Command routing and install modes

#### R-CMD-05
`weave command add <name>` without `--provider` MUST install for all enabled providers.

#### R-CMD-06
`weave command add <name> --provider <p>` MUST perform exclusive installation for provider `<p>` and MUST NOT require shared `.agents` command installation.

#### R-CMD-07
Codex command wrappers MUST be projected under namespace `__weave_commands__`.

#### R-CMD-08
When no providers are enabled, interactive sessions MUST prompt in English with default `[y/N]` before continuing with shared-only install behavior.

#### R-CMD-09
When no providers are enabled and session is non-interactive, `weave command add <name>` MUST fail automatically with actionable guidance and no mutations.

### Architecture and transactions

#### R-ARCH-08
Provider-specific command projection logic MUST remain outside CLI parsing and be encapsulated in app/provider strategy layer.

#### R-ARCH-09
Any failure during multi-provider command installation MUST trigger full rollback of filesystem + config effects.

#### R-CONFIG-08
Config persistence for command installs MUST remain transactional across provider targets (no partial persisted state).

### Debt closure

#### R-DEBT-01
`DEF-001` and `DEF-002` MUST be closed with automated CI evidence for release signing/checksum generation and install artifact validation.

#### R-DEBT-02
`DEF-003` MUST be closed with a CI gate that enforces migration notes when breaking changes are present.

#### R-DEBT-03
`DEF-004` MUST be closed with guided `git init` prompt behavior for interactive mode and deterministic non-interactive handling.

## Scenarios

### S1: Default install (all enabled providers)
Given enabled providers: `claude-code`, `opencode`, `codex`  
When `weave command add commit` runs  
Then command is projected to each enabled provider strategy  
And operation succeeds atomically or rolls back fully.

### S2: Exclusive provider install
Given enabled providers exist  
When `weave command add commit --provider codex`  
Then only Codex target is written  
And shared `.agents/commands/commit.md` is not required by that operation.

### S3: No providers enabled, interactive
Given no enabled providers and interactive TTY  
When `weave command add commit` runs  
Then prompt `No providers are currently enabled. Continue anyway? [y/N]` is shown  
And default answer is `No`.

### S4: No providers enabled, non-interactive
Given no enabled providers and non-interactive context  
When `weave command add commit` runs  
Then command fails with actionable error and no mutations.

### S5: Multi-provider failure rollback
Given multi-provider default install plan and one provider projection fails  
When apply is executed  
Then all previously applied projections are rolled back  
And `weave.yaml` remains unchanged.

### S6: Batch-0 debt closure evidence
Given CI workflows for release and install validation  
When release/debt checks run  
Then all DEF-001..DEF-003 closure criteria are enforced automatically  
And DEF-004 behavior is covered by tests.
