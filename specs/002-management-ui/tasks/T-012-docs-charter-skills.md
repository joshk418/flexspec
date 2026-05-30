---
id: T-012
name: Charter, skills, and repo documentation
parent_spec: '../README.md'
status: done
satisfies: [FR-015, FR-016]
depends_on: [T-003, T-009, T-010]
verified_by: [TC-010]
---

# T-012: Charter, skills, and repo documentation

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-015, FR-016 · **Depends on**: T-003, T-009, T-010 · **Verified by**: TC-010

## Objective

Update in-repo charter, agent skills, README, and embedded templates so new commands (`flexspec ui`, `list --json`, `status set`) and UI workflows are documented for humans and agents.

## Context

`flexspec init` copies `templates/charter.md` → `.flexspec/charter.md`; both must stay aligned. Skills under `skills/` ship via `npx skills add` — keep in sync with `.agents/skills` conventions. Do **not** edit user-installed copies outside this repo.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `templates/charter.md` | modify | §4 capabilities, §6 diagram, §9 glossary, §10 close UI question, §11 revision row |
| `.flexspec/charter.md` | modify | Same deltas as template (this repo's live charter) |
| `skills/flexspec/SKILL.md` | modify | CLI table, optional UI workflow in Phase 2 |
| `skills/flexspec-charter/SKILL.md` | modify | Mention `flexspec ui` for viewing charter (read-only) |
| `README.md` | modify | Features, usage table, dev build note (`make build-ui`) |
| `templates/README.md` | modify | CLI command table |
| `.flexspec/templates/README.md` | modify | Match `templates/README.md` if drifted |
| `cmd/root.go` | modify | Short help lists `ui` (optional) |

## Implementation Steps

1. **Charter §4** — Move "management UI" from Planned → Available; list `flexspec ui`, `flexspec list --json`, `flexspec status set`.
2. **Charter §6** — Extend mermaid: `cli -->|ui| UIServer[local HTTP + embedded UI]`.
3. **Charter §9** — Add glossary: `Management UI`, `flexspec ui`, SSE/live board (brief).
4. **Charter §10** — Remove or resolve non-blocking UI scope question; reference spec `002-management-ui`.
5. **Charter §11** — Row: `2026-05-30` · management UI + CLI · `002-management-ui`.
6. **`skills/flexspec/SKILL.md`** — Extend CLI table; add subsection **Optional: local dashboard** under Phase 2: suggest `flexspec ui --no-open` while implementing; agents may use `flexspec status set` or edit frontmatter; UI does not replace `/flexspec` lifecycle.
7. **`skills/flexspec-charter/SKILL.md`** — Note charter can be viewed in browser when UI is running (no in-UI charter editor v1).
8. **`README.md`** — Usage rows for `ui`, `list --json`, `status set`; contributor note: Node required only to build `web/`.
9. **`templates/README.md`** — Same CLI rows as README (init scaffolds this file).
10. Sync `.flexspec/charter.md` and `.flexspec/templates/README.md` with template sources.
11. Grep repo for stale "planned management UI" / CLI-only-four-commands prose; fix stragglers.

## Acceptance Criteria

- [ ] Charter lists UI as available, not planned _(FR-016)_
- [ ] `skills/flexspec/SKILL.md` documents all new commands with accurate flags _(FR-015)_
- [ ] README usage table matches implemented CLI _(FR-015)_
- [ ] `templates/charter.md` and `.flexspec/charter.md` agree on §4–§6 _(FR-016)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-010 | manual | Grep + read-through checklist | See parent spec |

Checklist:

```bash
rg 'flexspec ui' README.md skills/ templates/ .flexspec/
rg 'list --json|status set' README.md skills/
rg 'planned.*management UI|management UI.*planned' -i .
```

## Out of Scope

- Publishing skills to npm registry (separate release process).
- Updating external docs sites.

## Open Questions

- None.

## References

- Parent spec: [§5 Other](../README.md#5-other) (former charter follow-up)
- `specs/001-cli-validate` — prior charter sync pattern
