---
id: T-003
name: Templates + config migrations
parent_spec: ../README.md
status: done
satisfies: [FR-011, FR-012]
depends_on: [T-001]
verified_by: [TC-003, TC-004]
---

# T-003: Templates + config migrations

> **Parent**: [Update command](../README.md) · **Status**: todo
> **Satisfies**: FR-011, FR-012 · **Depends on**: T-001 · **Verified by**: TC-003, TC-004

## Objective

Add `templates-resync` (restore/compare embedded templates) and `config-keys` (reconcile config.yaml keys) migrations.

## Context

Embedded templates arrive as injected `fs.FS` (see `cmd.TemplatesFS` and `cmd/init.go::copyTemplates`). Config via `internal/config`. `--force` (a field on the migration, set by CLI) gates overwrite of differing files.

## Files In Scope

| File | Action |
| --- | --- |
| `internal/migrate/templates_resync.go` | create |
| `internal/migrate/templates_resync_test.go` | create |
| `internal/migrate/config_keys.go` | create |
| `internal/migrate/config_keys_test.go` | create |
| `cmd/init.go` | read (template walk pattern) |
| `internal/config/config.go` | read |

## Implementation Steps

1. `templates-resync`: walk injected `fs.FS`; for each template path compare to `.flexspec/templates/<rel>`.
   - missing → `Change{Kind:"create"}`; Apply writes it.
   - present but differing bytes → `Change{Kind:"report"}`; Apply overwrites **only if** force flag set.
2. `config-keys`: load config; known keys = `specs_dir`, `always_one_shot`, `spec_template`. Missing key → `Change{Kind:"rewrite", Detail:"add spec_template"}`; Apply writes key with documented default. Unknown keys → `Change{Kind:"report"}`.
3. Append both to `Registry`; pass `force` from CLI into `templates-resync`.
4. Tests: missing-template restore + differing-not-overwritten-without-force (TC-003); config missing `spec_template` added (TC-004).

## Acceptance Criteria

- [ ] Missing templates restored; differing reported unless `--force` _(FR-011, TC-003)_
- [ ] Missing config keys added with defaults _(FR-012, TC-004)_

## Testing

`go test ./internal/migrate/ -run 'Templates|Config'`.

## Out of Scope

Status/charter migrations; CLI flag parsing (consumes force only).

## Open Questions

None.

## References

- Parent FR-011, FR-012; `cmd/init.go`
