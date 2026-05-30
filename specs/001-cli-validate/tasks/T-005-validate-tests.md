---
id: T-005
name: Validation test coverage and CI
parent_spec: ../README.md
status: done
satisfies: [NF-002]
depends_on: [T-004]
verified_by: [TC-008]
---

# T-005: Validation test coverage and CI

> **Parent spec**: [CLI validate command](../README.md) · **Status**: todo
> **Satisfies**: NF-002 · **Depends on**: T-004 · **Verified by**: TC-008

## Objective

Ensure all validate rules have table-driven tests; `go test -race ./...` passes.

## Context

Charter §7: one test file per source file, table-driven per function. Fill gaps in T-001–T-004 tests; add edge cases (empty specs dir, expanded with zero tasks, BOM in frontmatter if `spec` already handles).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/validate/*_test.go` | modify | Complete coverage |
| `.github/workflows/ci.yml` | reference | No change required if tests already run |

## Implementation Steps

1. Audit each FR rule code has at least one test case.
2. Add integration-style test for `validate` command via `cmd` test helper.
3. Run `go test -race ./...` and fix failures.
4. Run `gofmt`, `go vet`, `golangci-lint` locally.

## Acceptance Criteria

- [ ] TC-008: full test suite green with `-race`.
- [ ] Every `Rule` constant used in production has test coverage.

## Testing

- `go test -race ./...`

## Out of Scope

- New CI workflow jobs (existing pipeline sufficient).

## Open Questions

- (none)

## References

- Parent spec §4 Testing Criteria
