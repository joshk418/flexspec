---
id: T-004
name: Charter-check migration
parent_spec: ../README.md
status: done
satisfies: [FR-013]
depends_on: [T-001]
verified_by: [TC-005]
---

# T-004: Charter-check migration

> **Parent**: [Update command](../README.md) · **Status**: todo
> **Satisfies**: FR-013 · **Depends on**: T-001 · **Verified by**: TC-005

## Objective

Add a **report-only** migration `charter-check` that detects missing required charter sections or leftover template placeholders, without editing charter prose.

## Context

Charter at `.flexspec/charter.md` is user-authored; charter §8 boundaries forbid auto-editing docs. So this migration only reports. Required section headers come from the embedded `templates/charter.md` (`## N.` headings).

## Files In Scope

| File | Action |
| --- | --- |
| `internal/migrate/charter_check.go` | create |
| `internal/migrate/charter_check_test.go` | create |

## Implementation Steps

1. `charter-check`: read `.flexspec/charter.md`; derive expected `##` section headings from embedded `templates/charter.md` (injected `fs.FS`).
2. For each expected heading missing in the project charter → `Change{Kind:"report", Detail:"missing §N"}`.
3. Detect leftover placeholders (`{`...`}` or `<!--` guidance) → `Change{Kind:"report"}`.
4. `Apply` returns the same report changes but performs **no writes** (report-only).
5. Append to `Registry`.
6. Table-driven test: charter missing a section → detected; assert file unchanged after Apply (TC-005).

## Acceptance Criteria

- [ ] Missing sections/placeholders reported _(FR-013)_
- [ ] No writes to charter in Detect or Apply _(NF-001, TC-005)_

## Testing

`go test ./internal/migrate/ -run Charter`.

## Out of Scope

Editing charter content; status/template/config migrations.

## Open Questions

None.

## References

- Parent FR-013; charter §8
