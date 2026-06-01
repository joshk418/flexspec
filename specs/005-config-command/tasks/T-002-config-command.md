---
id: T-002
name: flexspec config command
parent_spec: '../README.md'
status: done
satisfies: [FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, FR-007, NF-001, NF-002]
depends_on: [T-001]
verified_by: [TC-001, TC-002, TC-003, TC-004, TC-005]
---

# T-002: flexspec config command

> **Parent spec**: [Config command](../README.md) · **Status**: todo
> **Satisfies**: FR-001–FR-007, NF-001–NF-002 · **Depends on**: T-001 · **Verified by**: TC-001–TC-005

## Objective

Implement `flexspec config` with human table and `--json` output.

## Context

Follow `cmd/list.go`: `os.Getwd`, `config.Load`, `tabwriter`, flag `--json`. Register in `init()` with `rootCmd.AddCommand`. Use T-001 helpers for rows/JSON.

## Files In Scope

- `cmd/config.go`
- `cmd/config_test.go`

## Implementation Steps

1. Create `configCmd` with `Use: "config"`, appropriate Short/Long.
2. `RunE`: load config; on error return as-is.
3. If `--json`: encode `config.JSONMap(cfg)` with indent; return.
4. Else: print `KEY\tVALUE` header, flush rows from `DisplayEntries`, then `fmt.Printf("config: %s\n", config.ConfigPath(root))`.
5. Add `configJSON` bool flag `--json`.
6. `config_test.go`: temp dir fixtures via pattern from `list_test.go` / `config_test.go`; capture stdout; TC-001–TC-005.

## Acceptance Criteria

- [ ] `flexspec config` and `flexspec config --json` work from project root
- [ ] Help lists subcommand
- [ ] All TC-001–TC-005 covered

## Testing

```bash
go test ./cmd/... -run Config
go test -race ./...
```

## Out of Scope

- Docs/skill (T-003)
- validate rule changes (existing `config.Load` sufficient)

## Open Questions

(none)

## References

- `cmd/list.go`
- T-001 helpers
