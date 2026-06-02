---
id: T-003
name: Board kanban layout overhaul
parent_spec: ../README.md
status: done
satisfies: [FR-001, FR-007, FR-008, FR-009, NF-002]
depends_on: [T-001]
verified_by: [TC-003]
---

# T-003: Board kanban layout overhaul

> **Parent**: [Board page UI overhaul](../README.md) · **Status**: todo
> **Satisfies**: FR-001, FR-007, FR-008, FR-009, NF-002 · **Depends on**: T-001 · **Verified by**: TC-003, TC-008

## Objective

Refactor kanban so all lifecycle columns fit in the main content area at desktop widths without document-level horizontal scroll.

## Context

Current `BoardPage.tsx` uses `display:flex`, `minWidth:220`, `overflowX:auto` on the row — causes browser-level scrollbar. `.app-main` max-width 1400px. After T-001, six column keys: five statuses + `unassigned`.

## Files In Scope

| File | Action |
| --- | --- |
| `ui/src/pages/BoardPage.tsx` | modify |
| `ui/src/components/Layout.tsx` | modify |
| `ui/src/index.css` | modify |

## Implementation Steps

1. Import `SPEC_COLUMNS` from `api/status.ts` (post T-001).
2. Wrap kanban in `.board-kanban`; columns `.board-column`.
3. Column picker (FR-008): multi-select or dropdown; default all `SPEC_COLUMNS`; persist `flexspec.boardColumns` in localStorage; grid uses visible subset only; always include `unassigned` if any unassigned specs exist.
4. CSS grid: `repeat(N, minmax(0, 1fr))` for visible N; no body overflow at ≥1280px.
5. Remove inline `minWidth:220`; board-only `overflow-x:auto` for NF-002 narrow viewports.
6. Tighten cards; style in `index.css` _(FR-007)_.
7. Light shell polish: `.app-nav`, `.app-main` spacing/overflow _(FR-009)_.
8. Table view unchanged.

## Acceptance Criteria

- [ ] At 1280px+ width, no horizontal scroll on `body` with kanban visible _(FR-001)_
- [ ] Narrow viewport scroll contained in `.board-kanban` _(NF-002)_
- [ ] Visual polish uses theme tokens _(FR-007)_

## Testing

Manual TC-003: run `flexspec ui`, open `/board`, resize window, inspect overflow on `document.documentElement`.

## Out of Scope

Status migration, API changes.

## Open Questions

None.

## References

- Parent FR-001, NF-002, NF-003
