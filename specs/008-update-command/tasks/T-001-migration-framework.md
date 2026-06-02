---
id: T-001
name: Migration framework
parent_spec: ../README.md
status: done
satisfies: [FR-005, NF-001]
depends_on: []
verified_by: [TC-001]
---

# T-001: Migration framework

> **Parent**: [Update command](../README.md) · **Status**: todo
> **Satisfies**: FR-005, NF-001 · **Depends on**: — · **Verified by**: TC-001

## Objective

Create the `internal/migrate` package: `Change` type, `Migration` interface, ordered `Registry`, and `Plan`/`Apply` orchestrators. No migrations yet (added in T-002–T-004).

## Context

Mirrors `internal/validate` output model. Migrations inject embedded templates via `fs.FS` to avoid `embed`/import cycles. See parent §2.1.

## Files In Scope

| File | Action |
| --- | --- |
| `internal/migrate/migrate.go` | create |
| `internal/migrate/migrate_test.go` | create |

## Implementation Steps

1. Define `type Change struct { Migration, Path, Kind, Detail string }` (`Kind` ∈ rewrite|create|delete|report).
2. Define `Migration` interface: `ID()`, `Description()`, `Detect(root string, cfg config.Config) ([]Change, error)`, `Apply(root string, cfg config.Config) ([]Change, error)`.
3. Add unexported registry slice + `Registry(tmpl fs.FS) []Migration` returning migrations in deterministic order (empty list for now; T-002–T-004 append).
4. `Plan(root string, cfg config.Config, migs []Migration) ([]Change, error)` — call each `Detect`, aggregate; never write.
5. `Apply(root string, cfg config.Config, migs []Migration) ([]Change, error)` — call each `Apply`, aggregate.
6. Helper `Select(migs []Migration, ids []string) ([]Migration, error)` — filter by id; error on unknown.
7. Table-driven tests with a fake migration verifying Plan aggregates and writes nothing (TC-001, NF-001).

## Acceptance Criteria

- [ ] `Migration` iface + `Change` exported _(FR-005)_
- [ ] `Plan` performs zero writes _(NF-001, TC-001)_
- [ ] `Select` errors on unknown id

## Testing

`go test ./internal/migrate/...`.

## Out of Scope

Concrete migrations; CLI.

## Open Questions

None.

## References

- Parent §2.1, §2.2
