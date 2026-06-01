---
id: "T-003"
name: "Docs and charter updates"
parent_spec: "../README.md"
status: done
satisfies: [FR-009]
depends_on: [T-001, T-002]
verified_by: [TC-006]
---

# T-003: Docs and charter updates

> **Parent spec**: [Config update command and UI](../README.md) · **Status**: todo  
> **Satisfies**: FR-009 · **Depends on**: T-001, T-002 · **Verified by**: TC-006

## Objective

Update product docs so CLI/UI config write behavior is discoverable and charter stays current.

## Context

Charter is active and currently describes `flexspec config` as read-only. `README.md` has command table rows for `flexspec config` and `flexspec config --json`. User confirmed README may be modified as needed; do not remove it.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `README.md` | modify | Add `flexspec config set <key> <value>` and settings UI wording |
| `.flexspec/charter.md` | modify | Update §4, §9, and §11 revision history |

## Implementation Steps

1. In `README.md`, add command-table row for `flexspec config set <key> <value>`.
2. Update the `flexspec ui` feature/usage wording to mention structured settings editing if needed.
3. In `.flexspec/charter.md` §4, describe config as readable and writable through CLI set and structured UI settings.
4. In §9 glossary, update `Config` definition from read-only to read/update wording.
5. Add §11 revision row dated `2026-06-01`, source `006-config-update-command-and-ui`.

## Acceptance Criteria

- [ ] README documents config set syntax. _(FR-009)_
- [ ] README still exists and keeps unrelated content intact. _(FR-009)_
- [ ] Charter §4 and §9 reflect new config write capability. _(FR-009)_
- [ ] Charter §11 records this spec. _(FR-009)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-006 | manual | README command table and charter capability/glossary mention config updates | `README.md`, `.flexspec/charter.md` |

Run: `go run . validate`

## Out of Scope

- Rewriting README structure.
- Editing docs unrelated to config write behavior.
- Removing README or replacing it wholesale.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Related tasks: T-001, T-002
