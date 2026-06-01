---
id: "T-001"
name: "CLI config set"
parent_spec: "../README.md"
status: done
satisfies: [FR-001, FR-002, FR-003, FR-004, FR-005, FR-006, NF-001, NF-002, NF-004]
depends_on: []
verified_by: [TC-001, TC-002, TC-003]
---

# T-001: CLI config set

> **Parent spec**: [Config update command and UI](../README.md) · **Status**: todo  
> **Satisfies**: FR-001-FR-006, NF-001, NF-002, NF-004 · **Depends on**: none · **Verified by**: TC-001, TC-002, TC-003

## Objective

Add `flexspec config set <key> <value>` so users can update one known config key from CLI with existing validation and table output.

## Context

`cmd/config.go` currently prints config as table/JSON. `internal/config.Config` has `SpecsDir`, `AlwaysOneShot`, `SpecTemplate`, plus `Load`, `Save`, `DisplayEntries`, and validation. Keep read behavior unchanged.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `cmd/config.go` | modify | Add `set` child command and shared table render helper |
| `cmd/config_test.go` | modify | Add set command coverage; preserve read tests |
| `internal/config/config.go` | modify | Add key update/parser helper |
| `internal/config/config_test.go` | modify | Table-driven parser/update tests |
| `internal/config/config_save_test.go` | read/modify | Reuse save expectations if needed |

## Implementation Steps

1. Add config helper, e.g. `ApplyUpdate(cfg Config, key, value string) (Config, error)`, accepting `specs_dir`, `always_one_shot`, `spec_template`.
2. Parse `always_one_shot` with `strconv.ParseBool`; require non-empty `specs_dir`; allow `spec_template` values `simple`, `expanded`, or empty string.
3. In `cmd/config.go`, add `config set <key> <value>` as child command under `configCmd`; keep `configJSON` only for read command.
4. Move table printing into a local helper reused by read and set.
5. `set` flow: cwd -> `config.Load` -> `config.ApplyUpdate` -> `config.Save` -> print updated table.
6. Ensure unknown keys and invalid values return errors without writing.

## Acceptance Criteria

- [ ] `flexspec config` and `flexspec config --json` output stay unchanged. _(NF-004)_
- [ ] `flexspec config set specs_dir custom-specs` persists and reloads. _(FR-003, FR-005)_
- [ ] `always_one_shot` accepts valid bool strings and rejects invalid strings. _(FR-004)_
- [ ] `spec_template` accepts `simple`, `expanded`, or empty string. _(FR-004)_
- [ ] Successful set prints updated table. _(FR-006)_
- [ ] No new Go dependencies. _(NF-001)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-001 | unit | arg validation and missing config error | `cmd/config_test.go` |
| TC-002 | unit | each key updates, persists, and prints table | `cmd/config_test.go`, `internal/config/config_test.go` |
| TC-003 | unit | invalid key/value does not write | `cmd/config_test.go`, `internal/config/config_test.go` |

Run: `go test -race ./...`

## Out of Scope

- Adding config keys.
- Preserving YAML comments.
- Batch update syntax.
- Changing JSON read output.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Prior config read spec: [`../../005-config-command/README.md`](../../005-config-command/README.md)
