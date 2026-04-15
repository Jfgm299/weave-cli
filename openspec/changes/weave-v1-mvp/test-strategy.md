# Test Strategy — weave-v1-mvp

## Purpose

Define a consistent, auditable test structure for v1 implementation using TDD + SDD.

This document complements:

- `PRD.md` section 13 (Verification Strategy)
- `openspec/changes/weave-v1-mvp/tasks.md` (requirement checklist)

---

## 1) Test Levels and Intent

### Unit tests

- Scope: single function/service/module.
- Goal: validate business rules and edge-case handling quickly.
- Examples: config validation, plan generation, path normalization, error mapping.

### Integration tests

- Scope: multiple modules + filesystem boundaries.
- Goal: validate orchestration and transactional behavior.
- Examples: symlink creation + config persistence ordering, provider adapter wiring, doctor diagnostics.

### E2E CLI tests

- Scope: command invocation as a user would run it.
- Goal: validate observable behavior (exit codes, output, resulting project state).
- Examples: `forge`, `provider add`, `skill/command add/list/remove`, `doctor`.

---

## 2) Repository Test Layout (v1)

```text
internal/
  app/
    ...
    *_test.go                  # unit tests colocated with implementation
  config/
    ...
    *_test.go                  # unit tests
  fsops/
    ...
    *_test.go                  # unit tests
  providers/
    ...
    *_test.go                  # unit tests

test/
  integration/
    *.go                       # integration tests (filesystem + service orchestration)
  e2e/
    *.go                       # CLI E2E tests
  fixtures/
    ...                        # reusable fixture templates
  testdata/
    ...                        # golden files, static inputs
```

Notes:

- Unit tests stay colocated (`*_test.go`) for proximity and maintainability.
- Integration/E2E tests are centralized in `test/` for clearer execution boundaries.

---

## 3) Naming Conventions

### Test files

- Use Go standard naming: `*_test.go`.

### Test functions

- Pattern: `Test<Subject>_<Scenario>_<ExpectedResult>`
- Examples:
  - `TestForgePlanner_ConvergedProject_ReturnsNoOp`
  - `TestSkillAdd_InvalidSource_LeavesConfigUnchanged`
  - `TestExitCodes_RepeatedErrorScenario_IsDeterministic`

### Subtests

- Use `t.Run("<scenario>", ...)` for matrix-style edge cases.

---

## 4) Execution Commands (v1)

### All tests

```sh
go test ./...
```

### Unit tests only

```sh
go test ./internal/...
```

### Integration tests

Use build tags:

```sh
go test -tags=integration ./test/integration/...
```

### E2E CLI tests

Use build tags:

```sh
go test -tags=e2e ./test/e2e/...
```

### CI recommendation

```sh
go test ./... && go test -tags=integration ./test/integration/... && go test -tags=e2e ./test/e2e/...
```

---

## 5) Isolation and Determinism Rules

1. Every integration/E2E test MUST run in an isolated temp workspace.
2. Tests MUST NOT mutate user home directories (`~/.weave`, `~/.agents`) directly.
3. External tools MUST be mocked/stubbed unless the test explicitly validates tool detection behavior.
4. Tests MUST assert deterministic outputs where relevant (stable ordering, stable exit codes).
5. Time/random/network side effects SHOULD be abstracted and controlled in tests.

---

## 6) Requirement Coverage Rule (mandatory)

For each requirement (`R-*`) included in a batch:

- At least 1 success-path test
- At least 1 error/edge-case test
- At least 1 observable acceptance evidence

No requirement is marked done in `tasks.md` without satisfying this rule.

---

## 7) Traceability Workflow (SDD)

When implementing a task in `tasks.md`:

1. Add/update relevant tests first (TDD red).
2. Implement minimal code (green).
3. Record observable evidence in task completion notes.
4. Mark checklist item as done.

Each completed requirement should be traceable from:

- Requirement ID (`R-*`) → Task ID (`B*-T*`) → Test name(s) → Observable evidence.

---

## 8) Batch 1 Test Scope (immediate)

Batch 1 MUST include tests for:

- `R-CORE-01`, `R-CORE-02`, `R-CORE-03`
- `R-CONFIG-01`, `R-CONFIG-02`, `R-CONFIG-03`, `R-CONFIG-04`, `R-CONFIG-05`

At least one unit + one integration + one E2E assertion path must exist per requirement block as defined in `tasks.md`.
