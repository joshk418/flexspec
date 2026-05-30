---
id: T-011
name: UI and CLI integration tests
parent_spec: '../README.md'
status: done
satisfies: [NF-003]
depends_on: [T-004, T-009, T-010, T-012]
verified_by: [TC-008, TC-010]
---

# T-011: UI and CLI integration tests

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: NF-003 · **Depends on**: T-004, T-009, T-010, T-012 · **Verified by**: TC-008, TC-010

## Objective

Complete table-driven test coverage for UI handlers, frontmatter, CLI flags; ensure `go test -race ./...` passes.

## Context

Follow charter §7: one `_test.go` per source file. Use temp dirs with minimal `.flexspec/` + spec fixtures (copy pattern from `internal/validate/testutil_test.go`).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/ui/*_test.go` | modify/create | fill gaps from T-001–T-004 |
| `cmd/ui_test.go` | create | |
| `cmd/status_test.go` | create | |
| `cmd/list_test.go` | modify | JSON output |

## Implementation Steps

1. Shared test helper: `writeMinimalProject(root)` with config + one simple spec.
2. Cover TC-001 through TC-007 per parent spec table.
3. Table-driven cases for invalid config PUT, unknown spec PATCH, malformed frontmatter.
4. Run full suite with race detector in CI (existing workflow).
5. Run TC-010 grep checklist from T-012 after docs land.

## Acceptance Criteria

- [ ] All TC-* mapped to passing tests _(NF-003)_
- [ ] No flaky SSE tests (use generous timeout or inject hub directly) _(NF-003)_

## Testing

Run: `go test -race ./...`

## Out of Scope

- Frontend unit tests (optional later).

## Open Questions

- None.

## References

- `specs/001-cli-validate` test patterns
