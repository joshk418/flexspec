---
id: T-001
name: Config display helpers
parent_spec: '../README.md'
status: done
satisfies: [FR-003, FR-004, FR-005, FR-006]
depends_on: []
verified_by: [TC-002, TC-003]
---

# T-001: Config display helpers

> **Parent spec**: [Config command](../README.md) · **Status**: todo
> **Satisfies**: FR-003–FR-006 · **Depends on**: — · **Verified by**: TC-002, TC-003

## Objective

Add small helpers in `internal/config` that expose config as ordered key/value rows and as a JSON-friendly map for `flexspec config`.

## Context

`Config` already has `Load`, `Save`, `validate`. `cmd/list.go` uses `displayOrDash` for empty strings — config command can reuse that pattern from `cmd` or duplicate a one-liner in config package (prefer keeping display logic callable from `cmd/config.go` without duplicating field names in two places).

Fixed key order: `specs_dir`, `always_one_shot`, `spec_template`.

## Files In Scope

- `internal/config/config.go`
- `internal/config/config_test.go`

## Implementation Steps

1. Add type `Entry struct { Key, Value string }` (or `DisplayEntries(cfg Config) []Entry`).
2. Implement `DisplayEntries(cfg Config) []Entry` mapping struct fields to string values (`strconv.FormatBool` for bool).
3. Implement `JSONMap(cfg Config) map[string]any` with same three keys; `spec_template` as empty string when unset.
4. Add unit tests for entries order and JSON map contents.

## Acceptance Criteria

- [ ] Three entries in fixed order for a sample `Config`
- [ ] Empty `SpecTemplate` → value `-` in display entries, `""` in JSON map
- [ ] Tests in `config_test.go` pass

## Testing

- `go test ./internal/config/...`

## Out of Scope

- CLI command wiring (T-002)
- New config fields

## Open Questions

(none)

## References

- Parent spec §2.4 FR-003–FR-006
- `cmd/list.go` `displayOrDash`
