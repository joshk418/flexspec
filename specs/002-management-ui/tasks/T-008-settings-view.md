---
id: T-008
name: Settings view
parent_spec: '../README.md'
status: done
satisfies: [FR-008, FR-009]
depends_on: [T-005, T-004]
verified_by: [TC-005]
---

# T-008: Settings view

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-008, FR-009 · **Depends on**: T-005, T-004

## Objective

Settings page: UI preferences in `localStorage`; FlexSpec config editor with load/save via API.

## Context

UI prefs keys: `boardView` default (`kanban`|`table`), `theme` (`light`|`dark`|`system`). Config section: textarea with current YAML from `GET /api/config` serialized as YAML (client-side `js-yaml` or display JSON and let server accept JSON only — spec says YAML textarea: fetch raw file via new `GET /api/config/raw` **or** reconstruct YAML from JSON; prefer **`GET /api/config/raw`** returning file contents for fidelity).

Add to spec API if needed: `GET /api/config/raw` returns plain text YAML (small addition — note in implementation).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `web/src/pages/SettingsPage.tsx` | modify | full UI |
| `internal/ui/handlers.go` | modify | optional `GET /api/config/raw` |
| `web/src/hooks/usePreferences.ts` | create | localStorage |

## Implementation Steps

1. Section **Appearance**: theme select; default board view select → `localStorage`.
2. Section **FlexSpec config**: load raw YAML (`GET /api/config/raw` or read-only file endpoint); textarea edit.
3. Save button → `PUT /api/config` with parsed YAML (client parses to JSON for PUT body matching Go struct, or send YAML string with `Content-Type: text/yaml` — pick one and document).
4. Show server validation errors inline.
5. Success toast / message on save.

## Acceptance Criteria

- [ ] Theme and board default persist locally _(FR-008)_
- [ ] Invalid config shows error without writing file _(FR-009)_
- [ ] Valid save updates `.flexspec/config.yaml` _(FR-009)_

## Testing

Manual + TC-005.

## Out of Scope

- Charter editor.

## Open Questions

- None.

## References

- `.flexspec/config.yaml` schema in `internal/config/config.go`
