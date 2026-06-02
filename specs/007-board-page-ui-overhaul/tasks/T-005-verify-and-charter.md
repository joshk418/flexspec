---
id: T-005
name: Verify, build UI, charter update
parent_spec: ../README.md
status: done
satisfies: [FR-004, FR-005, FR-006, NF-001, NF-003]
depends_on: [T-002, T-003]
verified_by: [TC-006, TC-007]
---

# T-005: Verify, build UI, charter update

> **Parent**: [Board page UI overhaul](../README.md) · **Status**: todo
> **Satisfies**: FR-004–FR-006, NF-001, NF-003 · **Depends on**: T-002, T-003 · **Verified by**: TC-006, TC-007

## Objective

Run full test/build pipeline and update `.flexspec/charter.md` if user confirmed §4/§9 deltas.

## Context

User confirmed charter update. Append §11 row citing `007-board-page-ui-overhaul`.

## Files In Scope

| File | Action |
| --- | --- |
| `.flexspec/charter.md` | modify if approved |
| `internal/ui/server_test.go` | modify if column expectations |
| `cmd/list_test.go` | modify if needed |

## Implementation Steps

1. `go test -race ./...` _(TC-006)_.
2. `make build-ui` _(TC-007)_.
3. Update tests for five columns and `draft` status.
4. Charter §4 board UX + column picker; §9 glossary (`draft`, no `refined`/`initial`); §11 row.

## Acceptance Criteria

- [ ] CI-equivalent tests pass _(NF-001)_
- [ ] Embedded UI build succeeds _(NF-003)_
- [ ] Charter updated if user said yes _(charter process)_

## Testing

`go test -race ./...` and `make build-ui`.

## Out of Scope

New features beyond parent spec.

## Open Questions

None.

## References

- Parent §5 charter freshness
