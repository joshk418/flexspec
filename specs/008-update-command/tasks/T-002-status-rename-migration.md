---
id: T-002
name: Status-rename migration
parent_spec: ../README.md
status: done
satisfies: [FR-010]
depends_on: [T-001]
verified_by: [TC-002]
---

# T-002: Status-rename migration

> **Parent**: [Update command](../README.md) · **Status**: todo
> **Satisfies**: FR-010 · **Depends on**: T-001 · **Verified by**: TC-002

## Objective

Add a migration `status-rename` that detects legacy spec/task statuses and rewrites them on apply, reusing spec 007's `NormalizeSpecStatus`.

## Context

Spec 007 adds `internal/spec/status.go::NormalizeSpecStatus` (`refined`→`planned`, `initial`→`draft`). Frontmatter writes use existing `spec.SetFileStatus` (preserves body + other keys). Scan specs via `spec.List` or directory walk.

## Files In Scope

| File | Action |
| --- | --- |
| `internal/migrate/status_rename.go` | create |
| `internal/migrate/status_rename_test.go` | create |
| `internal/spec/status.go` | read (from 007) |
| `internal/spec/frontmatter.go` | read |

## Implementation Steps

1. New type implementing `Migration`; `ID()="status-rename"`.
2. `Detect`: walk spec READMEs + task files; for each, if `NormalizeSpecStatus(raw) != raw` (and raw non-empty), emit a `Change{Kind:"rewrite", Path, Detail:"refined→planned"}`.
3. `Apply`: for each detected file, call `spec.SetFileStatus(path, NormalizeSpecStatus(raw))`; return applied changes.
4. Append to `Registry`.
5. Table-driven test: temp project with `status: refined` README → Detect finds it; Apply rewrites to `planned`; assert body unchanged (TC-002).

## Acceptance Criteria

- [ ] Legacy statuses detected, not written during Detect _(FR-010, NF-001)_
- [ ] Apply rewrites frontmatter only, body intact _(TC-002)_

## Testing

`go test ./internal/migrate/ -run StatusRename`.

## Out of Scope

Template/config/charter migrations; CLI.

## Open Questions

None.

## References

- Parent FR-010; spec 007 `NormalizeSpecStatus`
