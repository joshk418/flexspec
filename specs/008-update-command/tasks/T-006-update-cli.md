---
id: T-006
name: flexspec update orchestrator
parent_spec: ../README.md
status: done
satisfies: [FR-001, FR-002, FR-003, FR-004, FR-006]
depends_on: [T-002, T-003, T-004, T-005]
verified_by: [TC-008, TC-009, TC-010]
---

# T-006: flexspec update orchestrator

> **Parent**: [Update command](../README.md) · **Status**: todo
> **Satisfies**: FR-001–FR-004, FR-006 · **Depends on**: T-002, T-003, T-004, T-005 · **Verified by**: TC-008, TC-009, TC-010

## Objective

Add `cmd/update.go` that runs CLI install + skills install + migrations. Default applies all; step flags select; `--dry-run` previews; `--check` is a migration CI gate.

## Context

One command per file (charter §7). Output style mirrors `cmd/validate.go`. `main.go` sets `cmd.TemplatesFS`; pass its subtree to `migrate.Registry`. Reuses `migrate` (T-001–T-004) and `selfupdate` (T-005).

## Files In Scope

| File | Action |
| --- | --- |
| `cmd/update.go` | create |
| `cmd/update_test.go` | create |
| `cmd/root.go` | modify (register) |
| `main.go` | modify if subtree needed |

## Implementation Steps

1. `updateCmd` flags: `--cli`, `--skills`, `--migrate` (bool), `--dry-run`, `--check`, `--force` (bool), `--only` (stringSlice).
2. Resolve step set: if none of `--cli/--skills/--migrate` set → all three (FR-001); else only those set (FR-002).
3. Resolve root + `config.Load`. Build migrations `migrate.Registry(templatesSubFS)`; apply `migrate.Select` if `--only`.
4. **Apply order: migrations → skills → CLI** (CLI last; new binary applies next run).
5. `--check`: detection only; print plan; exit 1 if any migration pending, else 0 (FR-006). No installs.
6. `--dry-run`: print migration plan + `selfupdate.PlanCLI/PlanSkills` for selected steps; no writes/exec (FR-003, NF-001).
7. Default (apply): run `migrate.Apply`, then `selfupdate.ApplySkills`, then `ApplyCLI` for selected steps; pass `--force` into templates migration.
8. Print grouped, tab-separated report + combined summary + exit code (FR-004). Advise re-run after CLI upgrade.
9. Register in `root.go`.
10. Tests: TC-008 (no flags ⇒ all; single flag restricts), TC-009 (`--dry-run` no writes/exec), TC-010 (`--check` exit codes).

## Acceptance Criteria

- [ ] No flags ⇒ all three applied; step flags restrict _(FR-001, FR-002, TC-008)_
- [ ] `--dry-run` writes/execs nothing _(FR-003, TC-009)_
- [ ] `--check` exit 1 when pending _(FR-006, TC-010)_
- [ ] Grouped report + single exit code _(FR-004)_

## Testing

`go test ./cmd/ -run Update`.

## Out of Scope

Migration/selfupdate internals; docs (T-007).

## Open Questions

None.

## References

- Parent §2.2, §2.4; `cmd/validate.go`, `cmd/init.go`
