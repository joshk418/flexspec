---
id: T-005
name: Self-update package (CLI + skills)
parent_spec: ../README.md
status: done
satisfies: [FR-007, FR-008, FR-009, NF-002]
depends_on: []
verified_by: [TC-006, TC-007]
---

# T-005: Self-update package (CLI + skills)

> **Parent**: [Update command](../README.md) · **Status**: todo
> **Satisfies**: FR-007, FR-008, FR-009, NF-002 · **Depends on**: — · **Verified by**: TC-006, TC-007

## Objective

Create `internal/selfupdate` wrapping `go install …@latest` (CLI) and `npx skills …` (skills) behind an injectable runner, with plan + apply functions and toolchain detection.

## Context

Installed version is `cmd.version` (`cmd/root.go`, release-please managed). Tests must not spawn real `go`/`npx` → inject a `Runner`. Confirm exact `npx skills` args from the root `README` before wiring. See parent §2.1.

## Files In Scope

| File | Action |
| --- | --- |
| `internal/selfupdate/selfupdate.go` | create |
| `internal/selfupdate/selfupdate_test.go` | create |
| `README.md` | read (confirm `npx skills` invocation) |

## Implementation Steps

1. `type Action struct { Target, Command, Detail string }` (`Target` ∈ cli|skills).
2. `type Runner func(name string, args ...string) error`; default runner uses `exec.LookPath` + `exec.Command` (NF-002).
3. `PlanCLI(installed string) Action` → command `go install github.com/joshk418/flexspec@latest`, detail `installed <version>`.
4. `ApplyCLI(installed string, run Runner) (Action, error)` → LookPath `go` (error names `go` if missing, FR-009); run; wrap exit error.
5. `PlanSkills() Action` / `ApplySkills(run Runner) (Action, error)` → `npx skills <args>`; LookPath `npx` (FR-009).
6. Table-driven tests: TC-006 (PlanCLI text + version; ApplyCLI args; missing `go` errors); TC-007 (ApplySkills args via injected runner).

## Acceptance Criteria

- [ ] `PlanCLI/PlanSkills` build correct commands; no exec _(FR-007, FR-008)_
- [ ] `ApplyCLI/ApplySkills` exec via runner; errors surfaced _(FR-007, FR-008)_
- [ ] Missing toolchain named in error _(FR-009)_

## Testing

`go test ./internal/selfupdate/...`.

## Out of Scope

CLI flag parsing / orchestration (T-006); migrations.

## Open Questions

None — confirm `npx skills` args from README during step 5.

## References

- Parent §2.1, §2.2; `cmd/root.go` `version`
