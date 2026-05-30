---
id: T-003
name: Specs directory and task validation
parent_spec: ../README.md
status: done
satisfies: [FR-006, FR-007, FR-008, FR-009]
depends_on: [T-001]
verified_by: [TC-004, TC-005, TC-006]
---

# T-003: Specs directory and task validation

> **Parent spec**: [CLI validate command](../README.md) · **Status**: todo
> **Satisfies**: FR-006–FR-009 · **Depends on**: T-001 · **Verified by**: TC-004–TC-006

## Objective

Validate every spec under `cfg.SpecsDir`: README frontmatter, expanded tasks, orphans, duplicate sequence numbers.

## Context

Reuse `spec.ParseSpecMeta` and `spec.ParseTaskMeta` — same code path as `list`. Spec folder pattern: `^\d{3}-` prefix (align with `spec.specID`). Non-matching child dirs → warning `specs.orphan_dir`. Duplicate numeric prefix → warning `specs.duplicate_sequence`.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/validate/specs.go` | create | `CheckSpecs` |
| `internal/validate/specs_test.go` | create | TC-004–TC-006 |

## Implementation Steps

1. If `specs_dir` does not exist → warning only (not error; `new` creates on demand).
2. Read subdirs; track seen sequence ints for FR-009.
3. For dirs matching `^\d{3}-`: require `README.md`; call `ParseSpecMeta`; errors → `specs.frontmatter`.
4. If `meta.SpecType` expanded: require `tasks/` exists; foreach `T-*.md` call `ParseTaskMeta`.
5. Dirs not matching pattern or missing README → warning `specs.orphan_dir`.
6. Register `CheckSpecs` in `Run` only when config loaded successfully (cfg.SpecsDir set).

## Acceptance Criteria

- [ ] TC-004: broken spec frontmatter → error.
- [ ] TC-005: expanded spec bad task file → error.
- [ ] TC-006: orphan dir warning; duplicate `001-*` warning.

## Testing

- Temp dirs with minimal markdown frontmatter fixtures.

## Out of Scope

- Cross-referencing FR/T/TC IDs in markdown bodies.
- `strict` status enum validation unless spec §5 confirms.

## Open Questions

- (none)

## References

- `internal/spec/spec.go` (`List`, `ParseSpecMeta`, `loadTasks`)
