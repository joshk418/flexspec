---
id: T-001
name: Validate core types and orchestration
parent_spec: ../README.md
status: done
satisfies: [FR-010, FR-011, NF-001]
depends_on: []
verified_by: [TC-007]
---

# T-001: Validate core types and orchestration

> **Parent spec**: [CLI validate command](../README.md) · **Status**: todo
> **Satisfies**: FR-010, FR-011, NF-001 · **Depends on**: — · **Verified by**: TC-007

## Objective

Introduce `internal/validate` with `Finding`, `Severity`, `Options`, and `Run(root, cfg, opts)` that aggregates findings from registered check functions.

## Context

Follow `internal/config` and `internal/spec` package style: small files, wrapped errors at boundaries, exported types only where needed. Checks append to a slice; no I/O in `validate.go` itself.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/validate/validate.go` | create | Types + `Run` |
| `internal/validate/validate_test.go` | create | Minimal orchestration test |

## Implementation Steps

1. Add `Severity` type (`Error`, `Warning`) and `Finding` struct: `Severity`, `Path`, `Rule`, `Message`.
2. Add `Options` with `Strict bool` (used later by T-002/T-003).
3. Define `type Check func(root string, cfg config.Config, opts Options) []Finding`.
4. Implement `Run(root string, cfg config.Config, opts Options, checks ...Check) []Finding` that concatenates results.
5. Add helper `HasErrors(findings []Finding) bool` for cmd exit logic.

## Acceptance Criteria

- [ ] Package compiles; `Run` invokes checks in order and returns combined slice.
- [ ] `HasErrors` true iff any `SeverityError` finding exists.
- [ ] No new dependencies beyond stdlib and `internal/config`.

## Testing

- TC-007 (partial): unit test `Run` with stub check returning one error and one warning.

## Out of Scope

- Individual validation rules (T-002, T-003).
- Cobra command (T-004).

## Open Questions

- (none — resolved at spec level)

## References

- Parent spec §2.1, FR-010, FR-011
