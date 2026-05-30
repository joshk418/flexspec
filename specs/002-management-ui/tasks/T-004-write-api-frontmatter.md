---
id: T-004
name: Write API and frontmatter helpers
parent_spec: '../README.md'
status: done
satisfies: [FR-009, FR-012, NF-003]
depends_on: [T-001]
verified_by: [TC-005, TC-007]
---

# T-004: Write API and frontmatter helpers

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-009, FR-012, NF-003 · **Depends on**: T-001 · **Verified by**: TC-005, TC-007

## Objective

Implement safe frontmatter `status` updates and config save; expose `PUT /api/config`, `PATCH` status routes.

## Context

Reuse `splitFrontmatter` logic from `internal/spec` (export or duplicate minimally in `internal/ui/frontmatter.go`). Add `config.Save(root, Config)` validating same rules as `Load`.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/ui/frontmatter.go` | create | `SetStatus(path, status string) error` |
| `internal/ui/handlers.go` | modify | PUT/PATCH handlers |
| `internal/config/config.go` | modify | `Save`, validation |
| `internal/ui/frontmatter_test.go` | create | round-trip tests |

## Implementation Steps

1. `SetStatus`: parse frontmatter YAML map, set `status` key, re-serialize preserving body unchanged.
2. `PUT /api/config`: decode JSON → `Config` struct → validate → write YAML to `.flexspec/config.yaml`.
3. `PATCH /api/specs/{dir}/status`: body `{status}` → update spec README.
4. `PATCH /api/specs/{dir}/tasks/{filename}/status`: update task file.
5. Return 400 with message on validation failure; 404 on missing paths.

## Acceptance Criteria

- [ ] Invalid `spec_template` on save returns 400 _(FR-009)_
- [ ] Status patch changes only `status` in frontmatter _(FR-012)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-005 | unit | config PUT validation | `internal/ui/handlers_test.go` |
| TC-007 | unit | frontmatter round-trip | `internal/ui/frontmatter_test.go` |

Run: `go test -race ./internal/...`

## Out of Scope

- Full markdown body editing.

## Open Questions

- None.

## References

- `internal/spec/spec.go` — `splitFrontmatter`
