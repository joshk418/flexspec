---
depends_on:
    - T-005
id: T-006
name: Verification
parent_spec: ../README.md
satisfies:
    - FR-001
    - FR-002
    - FR-003
    - FR-004
    - FR-005
    - FR-006
    - FR-007
    - FR-008
    - FR-009
    - FR-010
    - FR-011
    - FR-012
    - NF-001
    - NF-002
    - NF-003
    - NF-004
    - NF-005
status: done
verified_by:
    - TC-001
    - TC-002
    - TC-003
    - TC-004
    - TC-005
    - TC-006
    - TC-007
---

# T-006: Verification

> **Parent spec**: [CLI glossary](../README.md) · **Status**: todo  
> **Satisfies**: FR-001-FR-012, NF-001-NF-005 · **Depends on**: T-005 · **Verified by**: TC-001-TC-007

## Objective

Run the full verification set and fix any gaps before moving the implementation to review.

## Context

The charter requires `go test -race`, `gofmt`, `go vet`, `golangci-lint`, and `flexspec validate`. Keep fixes scoped to this spec.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/glossary/*` | read/modify | Fix test or implementation gaps |
| `cmd/glossary*` | read/modify | Fix CLI behavior gaps |
| `internal/validate/*` | read/modify | Fix validation gaps |
| `internal/update/**` | read/modify | Fix migration gaps |
| `skills/flexspec/SKILL.md` | read/modify | Fix skill policy gaps |
| `skills/flexspec-glossary-discovery/SKILL.md` | read/modify | Fix discovery skill gaps |
| `.flexspec/charter.md` | read/modify | Ensure delivered behavior is reflected |

## Implementation Steps

1. Run `gofmt -w` on changed Go files.
2. Run focused Go tests for changed packages, then `go test ./...`.
3. Run `go vet ./...` and `golangci-lint run`.
4. Run `flexspec validate`.
5. Review the diff against all requirements and testing criteria; fix any gaps.

## Acceptance Criteria

- [ ] All automated checks pass or any environmental blocker is documented. _(TC-007)_
- [ ] Manual skill review passes TC-005 and TC-006.
- [ ] Charter reflects glossary capability and automatic charter update behavior.
- [ ] No scope drift beyond this spec.

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-007 | verification | Full project checks pass after implementation | repository root |

Run: `gofmt -w <changed go files>` then `go test ./...`, `go vet ./...`, `golangci-lint run`, `flexspec validate`

## Out of Scope

- Product UI.
- Hosted services.
- Unrelated docs or README changes.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Depends on: [`T-005`](T-005-glossary-discovery-skill.md)
