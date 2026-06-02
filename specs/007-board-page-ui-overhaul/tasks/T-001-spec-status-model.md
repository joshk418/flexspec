---
id: T-001
name: Spec status model (Go + TS)
parent_spec: ../README.md
status: done
satisfies: [FR-002, FR-003, FR-005]
depends_on: []
verified_by: [TC-001, TC-002]
---

# T-001: Spec status model (Go + TS)

> **Parent**: [Board page UI overhaul](../README.md) · **Status**: todo
> **Satisfies**: FR-002, FR-003, FR-005 · **Depends on**: — · **Verified by**: TC-001, TC-002

## Objective

Add canonical spec lifecycle statuses in Go and mirror them in the UI so board columns and future validation share one ordered list.

## Context

Today `ui/src/api/client.ts` hard-codes six statuses including `refined` and `initial`. Parent spec: five statuses starting with `draft`. Normalize `refined` → `planned`, `initial` → `draft`.

## Files In Scope

| File | Action |
| --- | --- |
| `internal/spec/status.go` | create |
| `internal/spec/status_test.go` | create |
| `ui/src/api/status.ts` | create |
| `ui/src/api/client.ts` | modify — import columns from `status.ts` |

## Implementation Steps

1. Create `SpecStatuses() []string` returning `draft`, `planned`, `in_progress`, `in_review`, `complete` in order.
2. Add `NormalizeSpecStatus(s string) string` — trim/lowercase; map `refined` → `planned`, `initial` → `draft`; unknown → pass through for Unassigned.
3. Add `ColumnForSpecStatus(status string) string` — return status if in list, else `"unassigned"`.
4. Table-driven tests in `status_test.go` for TC-001/TC-002.
5. Move `SPEC_COLUMNS` and `columnForStatus` to `ui/src/api/status.ts`; re-export from `client.ts` if needed for imports.
6. Add comment in both files: "Keep in sync with internal/spec/status.go".

## Acceptance Criteria

- [ ] Five statuses in fixed order; no `refined` in canonical list _(FR-002)_
- [ ] Legacy `refined` normalizes per parent decision _(FR-002)_
- [ ] Unknown status → unassigned column _(FR-003)_
- [ ] Go tests pass _(TC-001, TC-002)_

## Testing

Run `go test ./internal/spec/... -run Status`.

## Out of Scope

Board CSS, templates, disk migration.

## Open Questions

None.

## References

- Parent §2.1, §2.5 FR-002/FR-003/FR-005
