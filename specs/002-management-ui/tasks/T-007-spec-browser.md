---
id: T-007
name: Spec browser and detail
parent_spec: '../README.md'
status: done
satisfies: [FR-006, FR-007, FR-010]
depends_on: [T-005, T-002]
verified_by: [TC-004]
---

# T-007: Spec browser and detail

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-006, FR-007, FR-010 · **Depends on**: T-005, T-002

## Objective

Spec list sidebar + detail pane: render spec README as GFM HTML; expanded specs show collapsible task list with task markdown bodies.

## Context

Use `react-markdown` + `remark-gfm`; `react-syntax-highlighter` or shiki for code fences. Task accordion: header shows `id`, `name`, `status`; expand loads body from detail API (already included in `GET /api/specs/{dir}`).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `web/src/pages/SpecsPage.tsx` | modify | list + detail layout |
| `web/src/components/SpecMarkdown.tsx` | create | markdown renderer |
| `web/src/components/TaskAccordion.tsx` | create | collapsible tasks |
| `web/package.json` | modify | markdown deps |

## Implementation Steps

1. Left panel: spec list from `/api/specs`; highlight selected `dir`.
2. Route `/specs/:dir` loads `GET /api/specs/{dir}`.
3. Render `markdown` field with `SpecMarkdown`.
4. If `tasks.length > 0`, render `TaskAccordion` per task (`markdown` body when expanded).
5. SSE: refetch active spec when `specs-changed` and `dir` matches.
6. Link from board cards to this route.

## Acceptance Criteria

- [ ] README renders headings, lists, code blocks _(FR-006)_
- [ ] Expanded spec shows collapsible tasks with bodies _(FR-007)_
- [ ] Live update when task file changes _(FR-010)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-004 | integration | detail API includes bodies | Go test in T-001; manual UI |

## Out of Scope

- In-browser editing of spec markdown.

## Open Questions

- None.

## References

- Parent spec: [§1 Summary views](../README.md#1-summary)
