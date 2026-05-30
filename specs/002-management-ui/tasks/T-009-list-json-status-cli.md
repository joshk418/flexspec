---
id: T-009
name: list --json and status set CLI
parent_spec: '../README.md'
status: done
satisfies: [FR-011, FR-012]
depends_on: [T-004]
verified_by: [TC-007, TC-003]
---

# T-009: list --json and status set CLI

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-011, FR-012 · **Depends on**: T-004 · **Verified by**: TC-007, TC-003

## Objective

Add `flexspec list --json` and `flexspec status set` sharing frontmatter logic with UI server.

## Context

Extract shared `SetStatus` to `internal/spec` or `internal/ui` imported by both cmd and handlers to avoid duplication.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `cmd/list.go` | modify | `--json` flag, `encoding/json` output |
| `cmd/status.go` | create | `status set` subcommand |
| `cmd/root.go` | modify | register status |
| `internal/ui/frontmatter.go` | modify | move to `internal/spec/status.go` if cleaner |

## Implementation Steps

1. `list --json`: encode same shape as `GET /api/specs` to stdout.
2. `status set <target> --status <value>` where target is spec dir name (`002-management-ui`) or numeric id (`002`).
3. `--task T-001-slug.md` optional for task status update.
4. Validate status against allowed enums (warn or error on unknown).
5. Reuse `SetStatus` from T-004.

## Acceptance Criteria

- [ ] `list --json` parses with `jq` _(FR-011)_
- [ ] `status set` updates frontmatter only _(FR-012)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-003 | integration | JSON matches API | `cmd/list_test.go` |
| TC-007 | unit | CLI status set | `cmd/status_test.go` |

Run: `go test -race ./cmd/...`

## Out of Scope

- `flexspec status get` (list covers read).

## Open Questions

- None.

## References

- LeanSpec `update --status` naming; we use `status set` for clarity
