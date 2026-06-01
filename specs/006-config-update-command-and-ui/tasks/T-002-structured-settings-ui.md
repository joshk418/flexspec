---
id: "T-002"
name: "Structured settings UI"
parent_spec: "../README.md"
status: done
satisfies: [FR-007, FR-008, NF-003, NF-004]
depends_on: [T-001]
verified_by: [TC-004, TC-005]
---

# T-002: Structured settings UI

> **Parent spec**: [Config update command and UI](../README.md) · **Status**: todo  
> **Satisfies**: FR-007, FR-008, NF-003, NF-004 · **Depends on**: T-001 · **Verified by**: TC-004, TC-005

## Objective

Replace Settings page YAML editing with structured config controls in a two-column table: names first, values second.

## Context

`SettingsPage.tsx` currently fetches raw YAML through `fetchConfigRaw`, parses/saves YAML through `saveConfigYAML`, and renders a textarea. Backend already exposes `GET /api/config` and `PUT /api/config` using `config.Config`.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `ui/src/api/client.ts` | modify | Add structured fetch/save helpers; remove normal YAML edit path |
| `ui/src/pages/SettingsPage.tsx` | modify | Render config table with typed controls |
| `ui/src/components/Select.tsx` | read | Reuse for enum/boolean controls |
| `ui/src/index.css` | modify | Add shared form/table styles if useful |
| `internal/ui/handlers.go` | read/modify | Keep JSON GET/PUT stable; raw endpoint may remain |
| `internal/ui/server_test.go` | modify | Add config API round-trip tests |

## Implementation Steps

1. In `client.ts`, expose `fetchConfig(): Promise<ProjectConfig>` using `GET /api/config`.
2. Replace `saveConfigYAML` with `saveConfig(config: ProjectConfig): Promise<ProjectConfig>` using `PUT /api/config`.
3. In `SettingsPage.tsx`, replace `yaml` state with `config` draft state.
4. Render a table with headers `Name` and `Value`; rows: `specs_dir`, `always_one_shot`, `spec_template`.
5. Use text input for `specs_dir`, `Select` for boolean values, and `Select` for template values: `Infer` (`""`), `Simple`, `Expanded`.
6. Save full structured config; show existing success/error messages. On error, retain draft values.
7. Add API tests for `GET /api/config`, valid `PUT /api/config`, and invalid `PUT` returning `400`.

## Acceptance Criteria

- [ ] Settings page no longer requires editing YAML text for config. _(FR-008)_
- [ ] Config section is a table with `Name` then `Value`. _(FR-007)_
- [ ] `specs_dir`, `always_one_shot`, and `spec_template` are editable through typed controls. _(FR-007)_
- [ ] Validation errors display and preserve the user's draft values. _(FR-008)_
- [ ] Existing API JSON shape remains compatible. _(NF-004)_
- [ ] Controls are keyboard-accessible native inputs/buttons/select-backed controls. _(NF-003)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-004 | build/manual | UI renders Name/Value table and no YAML textarea | `ui/src/pages/SettingsPage.tsx` |
| TC-005 | unit | config API JSON round-trip and invalid config `400` | `internal/ui/server_test.go` |

Run: `go test -race ./...` and `npm run build` from `ui/`.

## Out of Scope

- Adding a React test framework.
- Preserving raw YAML comments in UI writes.
- Creating config keys dynamically from unknown YAML.

## Open Questions

- None.

## References

- Parent spec: [`../README.md`](../README.md)
- Related tasks: T-001, T-003
