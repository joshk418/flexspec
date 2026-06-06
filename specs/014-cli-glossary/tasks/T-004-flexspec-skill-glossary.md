---
depends_on:
    - T-003
id: T-004
name: FlexSpec skill glossary workflow
parent_spec: ../README.md
satisfies:
    - FR-005
    - FR-007
    - FR-008
    - FR-012
    - FR-013
    - NF-004
    - NF-005
status: done
verified_by:
    - TC-005
---

# T-004: FlexSpec skill glossary workflow

> **Parent spec**: [CLI glossary](../README.md) · **Status**: todo  
> **Satisfies**: FR-005, FR-007, FR-008, FR-012, FR-013, NF-004, NF-005 · **Depends on**: T-003 · **Verified by**: TC-005

## Objective

Update `/flexspec` so it uses the glossary during lifecycle work, automatically updates `.flexspec/charter.md` when specs imply charter changes, and ensure `/flexspec-charter` invokes glossary discovery during charter work.

## Context

The skill currently asks before charter updates during Phase 1. The user explicitly changed that policy on 2026-06-06: future `/flexspec` runs should update charter files without asking when the change is in scope.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `skills/flexspec/SKILL.md` | modify | Charter policy and glossary workflow |
| `skills/flexspec-charter/SKILL.md` | modify | Charter creation/update invokes glossary discovery |
| `.flexspec/charter.md` | modify | Document automatic charter update policy and glossary capability |

## Implementation Steps

1. Replace the charter delta question workflow with automatic charter update rules.
2. Add a glossary gate to Phase 1 and Phase 2: read known terms, watch for project-specific unknowns, and record clear/confirmed terms.
3. Require user interview only when a project-specific term is unclear.
4. Require `flexspec glossary add` for persisted terms; avoid manual YAML edits from the skill.
5. Update `/flexspec-charter` so charter creation/update invokes glossary discovery while keeping discovery runnable standalone.
6. Update charter sections for planned/current glossary capability and the new skill behavior.

## Acceptance Criteria

- [ ] `/flexspec` no longer asks whether to update charter for in-scope deltas. _(FR-012)_
- [ ] Skill workflow reads glossary and identifies project-specific unknowns. _(FR-005)_
- [ ] Unclear terms trigger a user question before persistence. _(FR-007, NF-005)_
- [ ] Persisted terms are written through the CLI. _(FR-008)_
- [ ] `/flexspec-charter` invokes glossary discovery during charter work. _(FR-013)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-005 | manual review | Skills document glossary use, unclear-term interview, CLI persistence, automatic charter updates, and charter-to-glossary discovery handoff | `skills/flexspec/SKILL.md`, `skills/flexspec-charter/SKILL.md` |

Run: `flexspec validate`

## Out of Scope

- Discovery skill creation.
- Go command implementation.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Depends on: [`T-003`](T-003-metadata-integration.md)
