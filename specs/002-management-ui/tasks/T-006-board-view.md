---
id: T-006
name: Board view (kanban and table)
parent_spec: '../README.md'
status: done
satisfies: [FR-004, FR-005, FR-010, NF-005]
depends_on: [T-005, T-002]
verified_by: [TC-003]
---

# T-006: Board view (kanban and table)

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-004, FR-005, FR-010, NF-005 · **Depends on**: T-005, T-002

## Objective

Implement board page with kanban columns per spec status and a table toggle persisted in `localStorage`; refresh on SSE.

## Context

Column order: `initial`, `refined`, `planned`, `in_progress`, `in_review`, `complete`, plus `Unassigned` for empty/unknown status. Cards link to `/specs/:dir`.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `web/src/pages/BoardPage.tsx` | modify | full implementation |
| `web/src/components/BoardKanban.tsx` | create | columns + cards |
| `web/src/components/BoardTable.tsx` | create | sortable table |
| `web/src/hooks/useSpecs.ts` | create | fetch + SSE refetch |

## Implementation Steps

1. Load specs via `GET /api/specs`; group by `status` into columns.
2. Toggle control: Kanban | Table; store `boardView` in `localStorage`.
3. Kanban: horizontal scroll columns, card shows name, id, truncated description, type badge.
4. Table: columns ID, Name, Status, Type, Description; click row → spec detail.
5. Subscribe to SSE; on `specs-changed`, refetch list.
6. Loading and empty states.

## Acceptance Criteria

- [ ] All lifecycle statuses have columns; unknown status → Unassigned _(FR-004)_
- [ ] Toggle persists across reload _(FR-005)_
- [ ] External file edit updates board without manual refresh _(FR-010)_

## Testing

Manual + TC-003 API parity. Optional: Playwright smoke (out of scope unless time).

## Out of Scope

- Drag-and-drop status change (v1 stretch; API exists in T-004).

## Open Questions

- None.

## References

- `templates/README.md` — status lifecycle
