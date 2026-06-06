---
id: "T-003"
name: "Metadata integration"
parent_spec: "../README.md"
status: todo
satisfies: [FR-010, FR-011, NF-001, NF-003]
depends_on: [T-002]
verified_by: [TC-004]
---

# T-003: Metadata integration

> **Parent spec**: [CLI glossary](../README.md) · **Status**: todo  
> **Satisfies**: FR-010, FR-011, NF-001, NF-003 · **Depends on**: T-002 · **Verified by**: TC-004

## Objective

Make `.flexspec/glossary.yaml` a normal FlexSpec metadata file created by init/update and checked by validate.

## Context

`init` already scaffolds `.flexspec/` files from embedded templates. `validate` has focused checks under `internal/validate`. `update --migrate` has registered migrations for metadata backfills.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `templates/glossary.yaml` | create | Empty seed document |
| `cmd/init.go` | modify | Copy glossary seed without clobbering |
| `internal/validate/flexspec.go` | modify | Check glossary shape |
| `internal/validate/flexspec_test.go` | modify | Missing/malformed glossary cases |
| `internal/update/migrations/*` | modify/create | Backfill glossary file for existing projects |
| `internal/update/migrations/*_test.go` | modify/create | Migration test coverage |

## Implementation Steps

1. Add an empty glossary seed template with `version`, `updated`, and empty `terms`.
2. Extend init scaffolding to create `.flexspec/glossary.yaml` without overwriting.
3. Add validate findings for malformed glossary YAML or invalid term entries.
4. Register an update migration that creates the glossary only when missing.
5. Add focused tests for init, validate, and migration behavior.

## Acceptance Criteria

- [ ] New projects receive `.flexspec/glossary.yaml`. _(FR-010)_
- [ ] Existing glossary content is never clobbered. _(FR-010)_
- [ ] `flexspec update --migrate` backfills missing glossary files. _(FR-010)_
- [ ] `flexspec validate` reports schema errors. _(FR-011, NF-003)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-004 | unit/integration | Init, migration, and validate behavior for glossary metadata | `cmd/init_test.go`, `internal/validate/*`, `internal/update/migrations/*` |

Run: `go test ./cmd ./internal/validate ./internal/update/...`

## Out of Scope

- CLI list/query/add command behavior.
- Skill changes.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Depends on: [`T-002`](T-002-glossary-cli.md)
