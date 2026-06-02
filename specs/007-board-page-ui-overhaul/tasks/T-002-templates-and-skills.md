---
id: T-002
name: Templates and flexspec skill
parent_spec: ../README.md
status: done
satisfies: [FR-002, FR-006]
depends_on: [T-001]
verified_by: [TC-005]
---

# T-002: Templates and flexspec skill

> **Parent**: [Board page UI overhaul](../README.md) · **Status**: todo
> **Satisfies**: FR-002, FR-006 · **Depends on**: T-001 · **Verified by**: TC-005

## Objective

Remove `refined`; rename `initial` → `draft` in templates and `/flexspec` skill phase routing.

## Context

Templates and skill currently use `initial` / `refined`. Target: Author `draft` → `planned` only.

## Files In Scope

| File | Action |
| --- | --- |
| `.flexspec/templates/flexspec-simple.md` | modify |
| `.flexspec/templates/expanded/flexspec-expanded.md` | modify |
| `.flexspec/templates/README.md` | modify |
| `templates/flexspec-simple.md` | modify |
| `templates/expanded/flexspec-expanded.md` | modify |
| `templates/README.md` | modify |
| `skills/flexspec/SKILL.md` | modify |
| `README.md` | modify status docs |

## Implementation Steps

1. Update frontmatter `status:` hints to five values only.
2. Replace "must not advance past `refined`" prose with "must not set `planned` while blocking open questions remain" (Section 5).
3. Update status tables in both `templates/README.md` files.
4. In `skills/flexspec/SKILL.md`: Phase 1 → `none` / `draft` → `planned`; remove `refined` and `initial`.
5. Grep repo for `refined` and spec-status `initial` in templates/skills; fix root README.

## Acceptance Criteria

- [ ] No `refined` in template status enums _(FR-002)_
- [ ] SKILL phase table matches parent §2.1 _(FR-006, TC-005)_

## Testing

Manual: `rg 'refined' skills/flexspec .flexspec/templates templates` — only historical/changelog context allowed.

## Out of Scope

Board UI, spec dir migration.

## Open Questions

—

## References

- Parent FR-006; spec 004-interviews still valid (gate before `planned`)
