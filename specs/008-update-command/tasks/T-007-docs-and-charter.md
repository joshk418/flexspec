---
id: T-007
name: Docs, charter, verify
parent_spec: ../README.md
status: done
satisfies: [NF-003]
depends_on: [T-006]
verified_by: [TC-011]
---

# T-007: Docs, charter, verify

> **Parent**: [Update command](../README.md) · **Status**: todo
> **Satisfies**: NF-003 · **Depends on**: T-006 · **Verified by**: TC-011

## Objective

Document `flexspec update`, update the charter (§4/§5/§8/§9/§11), and run the verification pipeline.

## Context

Root `README.md` has a command table. Charter at `.flexspec/charter.md`. The §8 carve-out (update may modify its own binary + installed skills) was approved.

## Files In Scope

| File | Action |
| --- | --- |
| `README.md` | modify |
| `.flexspec/charter.md` | modify |

## Implementation Steps

1. README: add `flexspec update` to the command table; document default-all behavior, `--cli/--skills/--migrate`, `--dry-run`, `--check`, `--only`, `--force`; note `--skills` needs Node/`npx`; note breaking status migration.
2. Charter §4: capability bullet for `flexspec update` (self-update + migrations).
3. Charter §5: note `--skills` uses `npx` (Node) for that step only; otherwise runtime Node-free.
4. Charter §8: carve-out sentence (update may modify own binary + installed skills).
5. Charter §9: glossary "Migration", "Update", "Self-update".
6. Charter §11: revision row citing `008-update-command`.
7. Run `go test -race ./...`, `go vet ./...`, `gofmt -l .` (TC-011); `go run . update --dry-run` to sanity-check the plan on this repo.

## Acceptance Criteria

- [ ] README documents `flexspec update` _(NF-003)_
- [ ] Charter §4/§5/§8/§9/§11 updated _(charter process)_
- [ ] CI-equivalent checks pass _(TC-011)_

## Testing

`go test -race ./...`; manual `go run . update --dry-run`.

## Out of Scope

New steps/migrations beyond parent spec.

## Open Questions

None.

## References

- Parent §5 charter freshness
