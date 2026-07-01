---
blocks: []
depends_on:
    - T-005
id: T-006
name: Verification pass
parent_spec: ../README.md
satisfies:
    - NF-003
    - NF-004
status: done
verified_by:
    - TC-005
---

# T-006: Verification pass

> **Parent spec**: [Brainstorm skill and CLI scaffolding](../README.md) · **Status**: todo
> **Satisfies**: NF-003, NF-004 · **Depends on**: T-005 · **Verified by**: TC-005
> **Blocks**: none

## Objective

Run the full CI check suite plus manual CLI and migration smoke tests to confirm the feature works end-to-end and does not regress existing behavior (`flexspec validate`, `flexspec list`, existing template migrations).

## Context

Charter §7 CI gate: `go test -race`, `gofmt` clean, `go vet`, `golangci-lint`. This task also exercises the operational claim in parent spec §3/FR-002: that `flexspec update --migrate` backfills `.flexspec/templates/brainstorm.md` into a project that was `init`'d before this feature shipped, using the pre-existing `templates-resync` migration with no new migration code. This is the one behavior in this spec that can only be proven by simulating an "old" project state, so it gets an explicit manual test case (TC-005) here rather than a unit test.

Also confirm the two declared non-goals hold: `flexspec validate`'s `requiredTemplates` list (`internal/validate/flexspec.go`) was NOT extended to require `brainstorm.md` (NF-004), and `flexspec list` output is unchanged by the presence of `.flexspec/brainstorms/*.md` files (FR-007).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| (repo-wide) | read | Run test/lint/vet across the whole module |
| `internal/brainstorm/*` | read | Confirm T-002 tests pass |
| `cmd/brainstorm*.go` | read | Confirm T-003 tests pass |
| `internal/validate/flexspec.go` | read | Confirm `requiredTemplates` unchanged |

### Workflow / Requirement Mapping

| Parent Section | Mapping |
| --- | --- |
| Workflow graph steps | §6.1 step 7 (template-missing branch), §6.2 (full ingestion flow, verified manually) |
| Implementation plan steps | §7.4 step 6 |
| Requirements | NF-003, NF-004 |
| Tests | TC-005 |

## Implementation Steps

1. Run `go test ./... -race` from the repo root; all packages, including new `internal/brainstorm` and `cmd` tests, must pass.
2. Run `gofmt -l .`; output must be empty.
3. Run `go vet ./...`; must exit clean.
4. Run `golangci-lint run`; must exit clean.
5. Run `flexspec validate` against this repo's own `.flexspec/`; confirm no new errors or warnings related to brainstorm files, and confirm `brainstorm.md` is not flagged as a required-but-missing template.
6. Manual smoke test (new project): in a scratch temp directory, run `flexspec init`, then `flexspec brainstorm new demo-feature`; confirm `.flexspec/brainstorms/demo-feature.md` exists with the template content; run it again without `--force` (expect an error) and with `--force` (expect success).
7. Manual smoke test (TC-005, pre-existing project): in a second scratch temp directory, simulate an "old" `.flexspec/` by running `flexspec init` and then deleting `.flexspec/templates/brainstorm.md`; run `flexspec update --check`; confirm the `templates-resync` migration reports `brainstorm.md` as pending; run `flexspec update --migrate`; confirm `.flexspec/templates/brainstorm.md` is created and `flexspec brainstorm new` now works.
8. Confirm `flexspec list` output (with or without `--json`) is unaffected by the presence of files under `.flexspec/brainstorms/` (create one, re-run `flexspec list`, diff output before/after).
9. Update this spec's status to `in_review` via `flexspec status set 016-flexspec-brainstorm-skill --status in_review` once all checks pass (per `/flexspec` Phase 2 exit).

## Acceptance Criteria

- [ ] `go test ./... -race`, `gofmt -l .` (empty), `go vet ./...`, `golangci-lint run` all pass. _(NF-003)_
- [ ] `flexspec validate` passes with no new findings; `brainstorm.md` is confirmed absent from `requiredTemplates`. _(NF-004)_
- [ ] New-project smoke test (init -> brainstorm new -> overwrite guard -> `--force`) passes.
- [ ] Pre-existing-project smoke test (TC-005: missing template -> `update --check` flags it -> `update --migrate` fixes it) passes.
- [ ] `flexspec list` output is unchanged by the presence of brainstorm docs.
- [ ] All tests in "Testing" below pass.

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-005 | manual | `flexspec update --check`/`--migrate` backfills `.flexspec/templates/brainstorm.md` into a pre-existing project via the existing `templates-resync` migration, with no new migration code | scratch temp directory (manual) |

Run: `go test ./... -race && gofmt -l . && go vet ./... && golangci-lint run && flexspec validate`

## Out of Scope

- Do not write a new automated migration test for `brainstorm.md` specifically — `internal/migrate/templates_resync_test.go` already covers the generic resync behavior; this task only needs to confirm it applies here (manual TC-005), not duplicate that package's test suite.
- Do not add brainstorm-aware checks to `internal/validate` (that would contradict NF-004).

## Open Questions

None.

## References

- Parent spec: [`../README.md`](../README.md) §7.4 step 6, §8 TC-005
- Related tasks: all prior tasks (T-001–T-005)
