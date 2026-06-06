---
id: "T-001"
name: "Glossary store"
parent_spec: "../README.md"
status: todo
satisfies: [FR-002, FR-003, NF-001, NF-003]
depends_on: []
verified_by: [TC-001]
---

# T-001: Glossary store

> **Parent spec**: [CLI glossary](../README.md) · **Status**: todo  
> **Satisfies**: FR-002, FR-003, NF-001, NF-003 · **Depends on**: none · **Verified by**: TC-001

## Objective

Create the structured glossary package that loads, saves, searches, and upserts `.flexspec/glossary.yaml`.

## Context

Follow existing Go patterns: small internal package, exported functions documented, errors wrapped with `%w`, table-driven tests. This owns CLI code-map steps 2-5.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/glossary/glossary.go` | create | YAML document model and helpers |
| `internal/glossary/glossary_test.go` | create | Store/search/upsert tests |

## Implementation Steps

1. Define `Document` and `Entry` structs for `version`, `updated`, and sorted `terms`.
2. Implement `Load(root string) (Document, error)` that returns an empty document when the glossary file is missing.
3. Implement `Save(root string, doc Document) error` with deterministic term sorting and `.flexspec/` path creation.
4. Implement `Query(doc Document, text string) []Entry` with exact term/alias matches before substring matches.
5. Implement `Upsert(doc Document, entry Entry) (Document, error)` requiring non-empty term and definition.

## Acceptance Criteria

- [ ] Missing glossary returns an empty document without error. _(NF-001)_
- [ ] Malformed YAML and invalid entries report path-aware errors. _(NF-003)_
- [ ] Save/upsert output is deterministic by term. _(FR-003)_
- [ ] Query ranks exact term and alias matches before substring matches. _(FR-002)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-001 | unit | Missing file, malformed YAML, sort order, upsert preservation, and query ranking | `internal/glossary/glossary_test.go` |

Run: `go test ./internal/glossary`

## Out of Scope

- Cobra command wiring.
- Init/update/validate integration.
- Skill files.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
