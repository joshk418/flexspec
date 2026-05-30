---
id: T-002
name: Config and .flexspec layout validation
parent_spec: ../README.md
status: done
satisfies: [FR-001, FR-002, FR-003, FR-004, FR-005]
depends_on: [T-001]
verified_by: [TC-001, TC-002, TC-003]
---

# T-002: Config and .flexspec layout validation

> **Parent spec**: [CLI validate command](../README.md) · **Status**: todo
> **Satisfies**: FR-001–FR-005 · **Depends on**: T-001 · **Verified by**: TC-001–TC-003

## Objective

Implement checks for `config.yaml`, charter frontmatter, and required template files under `.flexspec/`.

## Context

Mirror `config.Load` errors as findings with rule codes `config.missing`, `config.parse`, `config.specs_dir`, `config.spec_template`. Template list matches `cmd/init.go` scaffold paths. Use `spec.splitFrontmatter` behavior indirectly via reading charter — either export a small helper in `spec` or duplicate minimal `---` check in validate (prefer calling existing parse if charter uses same frontmatter shape).

Charter is markdown with YAML frontmatter like specs.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/validate/config.go` | create | `CheckConfig` |
| `internal/validate/flexspec.go` | create | `CheckFlexspecDir` |
| `internal/validate/config_test.go` | create | TC-001, TC-002 |
| `internal/validate/flexspec_test.go` | create | TC-003 |

## Implementation Steps

1. `CheckConfig`: if config missing → single error, return (skip dependent checks if spec says so).
2. Else `config.Load`; on error → map to finding with path `.flexspec/config.yaml`.
3. If `spec_template` set and not `simple`/`expanded` → error `config.spec_template`.
4. `CheckFlexspecDir`: stat `.flexspec/charter.md`; read and verify frontmatter delimiters (reuse `spec.ParseSpecMeta` pattern or add `ParseCharterMeta` only if needed).
5. For each required template path under `.flexspec/templates/`, stat file; missing → error `templates.missing`.
6. If `opts.Strict`: charter placeholder checks per spec §5 answer (implement only if user confirms).

## Acceptance Criteria

- [ ] TC-001: no config file → error finding, rule `config.missing`.
- [ ] TC-002: bad YAML and invalid `spec_template` produce errors.
- [ ] TC-003: missing charter or template file produces errors.

## Testing

- Table-driven tests with `t.TempDir()` and minimal file trees.

## Out of Scope

- Specs directory (T-003).
- Fixing or writing files.

## Open Questions

- (none)

## References

- `cmd/init.go` template paths
- `internal/config/config.go`
