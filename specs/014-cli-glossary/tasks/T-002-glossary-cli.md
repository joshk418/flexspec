---
depends_on:
    - T-001
id: T-002
name: Glossary CLI
parent_spec: ../README.md
satisfies:
    - FR-001
    - FR-002
    - FR-003
    - FR-004
    - NF-002
status: done
verified_by:
    - TC-002
    - TC-003
---

# T-002: Glossary CLI

> **Parent spec**: [CLI glossary](../README.md) · **Status**: todo  
> **Satisfies**: FR-001, FR-002, FR-003, FR-004, NF-002 · **Depends on**: T-001 · **Verified by**: TC-002, TC-003

## Objective

Add `flexspec glossary list`, `flexspec glossary query <text>`, and `flexspec glossary add <term>` using the glossary store.

## Context

Commands live one file per command group under `cmd/`. Human-readable output should use `internal/clioutput.WriteTable`; JSON output should use indented encoding like `list --json`.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `cmd/glossary.go` | create | Cobra command group and subcommands |
| `cmd/glossary_test.go` | create | CLI output and add/query flow tests |
| `internal/glossary/glossary.go` | read | Store API from T-001 |

## Implementation Steps

1. Register `glossaryCmd` under `rootCmd` with `list`, `query`, and `add` subcommands.
2. Add `--json` to `list` and `query`; render table output by default.
3. Add `add <term> --definition <text> [--alias <value>] [--category <value>] [--source <value>]`.
4. Validate required args/flags and preserve exit 0 for empty list/query results.
5. Write tests using temp project roots and command output buffers.

## Acceptance Criteria

- [ ] `flexspec glossary list` prints all terms or a no-terms message. _(FR-001, FR-004)_
- [ ] `flexspec glossary query "text"` returns ranked matches or a no-matches message. _(FR-002, FR-004)_
- [ ] `flexspec glossary add` persists an entry that can be queried. _(FR-003)_
- [ ] Table and JSON outputs are stable. _(NF-002)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-002 | unit | List/query output, JSON, empty state, missing args | `cmd/glossary_test.go` |
| TC-003 | integration | Add followed by query returns persisted entry | `cmd/glossary_test.go` |

Run: `go test ./cmd`

## Out of Scope

- `init`, `update`, or `validate` changes.
- Discovery heuristics.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Depends on: [`T-001`](T-001-glossary-store.md)
