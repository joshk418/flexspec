---
id: T-004
name: Cobra validate command
parent_spec: ../README.md
status: done
satisfies: [FR-010, FR-011, NF-003]
depends_on: [T-002, T-003]
verified_by: [TC-007]
---

# T-004: Cobra validate command

> **Parent spec**: [CLI validate command](../README.md) · **Status**: todo
> **Satisfies**: FR-010, FR-011, NF-003 · **Depends on**: T-002, T-003 · **Verified by**: TC-007

## Objective

Add `flexspec validate` command: run all checks, print findings, exit 0/1 per FR-011.

## Context

Match `cmd/list.go`: `os.Getwd()`, `config.Load`, delegate to `validate.Run`. One finding per line: `error  path  rule  message` (exact format stable for CI). Summary: `N error(s), M warning(s)`.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `cmd/validate.go` | create | Cobra + flags |
| `cmd/root.go` | modify | `rootCmd.AddCommand(validateCmd)` |

## Implementation Steps

1. Create `validateCmd` with `Use: "validate"`, short/long help.
2. Flags: `--strict` (bool), `--json` (bool) — wire JSON only if spec §5 confirms; else omit flag.
3. `RunE`: load config (may fail); pass cfg zero value if missing but still run checks per §5 Q4 answer.
4. Call `validate.Run` with `CheckConfig`, `CheckFlexspecDir`, `CheckSpecs` as appropriate.
5. Print findings; if `validate.HasErrors` → return error or `os.Exit(1)` via Cobra pattern used elsewhere.
6. Update root help text if project documents commands in README (out of scope per charter §8).

## Acceptance Criteria

- [ ] `flexspec validate` registered and appears in `flexspec --help`.
- [ ] Exit 0 on clean temp project; exit 1 when config missing (per FR-011).
- [ ] TC-007 golden output test passes.

## Testing

- `cmd` package test executing validate subcommand or test `RunE` directly.

## Out of Scope

- README documentation updates.

## Open Questions

- (none)

## References

- `cmd/list.go`
