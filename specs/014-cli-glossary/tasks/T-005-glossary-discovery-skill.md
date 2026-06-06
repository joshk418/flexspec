---
depends_on:
    - T-004
id: T-005
name: Glossary discovery skill
parent_spec: ../README.md
satisfies:
    - FR-006
    - FR-007
    - FR-008
    - FR-009
    - FR-013
    - NF-004
    - NF-005
status: done
verified_by:
    - TC-006
---

# T-005: Glossary discovery skill

> **Parent spec**: [CLI glossary](../README.md) · **Status**: todo  
> **Satisfies**: FR-006, FR-007, FR-008, FR-009, FR-013, NF-004, NF-005 · **Depends on**: T-004 · **Verified by**: TC-006

## Objective

Create `flexspec-glossary-discovery`, a skill that scans project language, asks about unclear terms, writes confirmed definitions through the CLI, and remains runnable standalone when `/flexspec-charter` also invokes it.

## Context

The skill should be token-conscious: scan with `rg`, process candidate lists in code when large, compare against `flexspec glossary list --json`, and avoid web lookups.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `skills/flexspec-glossary-discovery/SKILL.md` | create | Discovery workflow, triggers, scan policy, interview rules |
| `.flexspec/charter.md` | modify | Add skill to planned/current skills list and glossary |

## Implementation Steps

1. Add valid skill frontmatter with triggers for glossary discovery requests.
2. Define scan sources: specs, charter, likely docs, code identifiers, and config names.
3. Define filtering rules: exclude common language, existing glossary terms, and vendor/library names unless project-specific.
4. Define interview workflow for unclear terms and CLI persistence for confirmed terms.
5. Define final report shape with added, skipped, and ambiguous terms.
6. Include standalone trigger aliases, including `flexspec-glossary`, for manual glossary updates.

## Acceptance Criteria

- [ ] Skill finds candidate project terms from repository text and identifiers. _(FR-006)_
- [ ] Skill asks for exact meaning only when meaning is unclear. _(FR-007, NF-005)_
- [ ] Skill writes confirmed terms through `flexspec glossary add`. _(FR-008)_
- [ ] Skill reports added/skipped/ambiguous terms. _(FR-009)_
- [ ] Skill remains runnable standalone after `/flexspec-charter` integration. _(FR-013)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-006 | manual review | Discovery workflow covers scan, filtering, interview, persistence, and report | `skills/flexspec-glossary-discovery/SKILL.md` |

Run: `flexspec validate`

## Out of Scope

- Embeddings or network searches.
- Writing project documentation outside FlexSpec metadata.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Depends on: [`T-004`](T-004-flexspec-skill-glossary.md)
