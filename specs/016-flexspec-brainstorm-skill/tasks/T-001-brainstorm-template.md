---
blocks:
    - T-002
depends_on: []
id: T-001
name: Embedded brainstorm template + docs
parent_spec: ../README.md
satisfies:
    - FR-002
status: done
verified_by: []
---

# T-001: Embedded brainstorm template + docs

> **Parent spec**: [Brainstorm skill and CLI scaffolding](../README.md) · **Status**: todo
> **Satisfies**: FR-002 · **Depends on**: none · **Verified by**: (docs task, no dedicated test)
> **Blocks**: T-002

## Objective

Author `templates/brainstorm.md`, the embedded template that later tasks scaffold into `.flexspec/brainstorms/<slug>.md`, and document it in `templates/README.md`.

## Context

`templates/` holds every embedded template (`charter.md`, `flexspec-simple.md`, `expanded/*`, `glossary.yaml`). `main.go` uses `//go:embed all:templates`, so any new file placed here is automatically embedded — no Go code change needed to register it.

Two existing mechanisms will pick this file up automatically once it exists:
- `cmd/init.go`'s `copyTemplates` walks the embedded tree and copies every file except `charter.md`/`glossary.yaml` into `.flexspec/templates/`, preserving structure. `brainstorm.md` will be copied there on `flexspec init`.
- `internal/migrate/templates_resync.go`'s `templatesResyncMigration` walks the same embedded tree (excluding only `charter.md`) and backfills missing files into `.flexspec/templates/` for existing projects on `flexspec update --migrate`. No new migration ID is needed.

Unlike `charter.md` (scaffolded once to a fixed path, filled by a skill interview) or `flexspec-simple.md`/`flexspec-expanded.md` (read fresh by `flexspec new` on every invocation), `brainstorm.md` follows the **spec-template pattern**: it lives at `.flexspec/templates/brainstorm.md` after init, is user-editable, and is read fresh by `internal/brainstorm.Create` (T-002) on every `flexspec brainstorm new` call.

Template content should mirror the Discovery Gate's topic areas from `skills/flexspec/SKILL.md` (goal/scope, user/workflow, data/state, interfaces, security/privacy, failure handling, concurrency, UX/accessibility, performance/scale, operations, testing) plus an alternatives/tradeoffs section — this shape is what lets `/flexspec` Phase 1 (T-005) map ingested brainstorm content directly onto spec sections. Use `{placeholder}` / `<!-- -->` guidance-comment conventions like `templates/charter.md`, since the `/flexspec-brainstorm` skill (T-004), not the CLI, fills the content.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `templates/brainstorm.md` | create | New embedded template, placeholder-driven like `templates/charter.md` |
| `templates/README.md` | modify | Add a template row, CLI command row, and a short "Brainstorm Docs" note |

### Workflow / Requirement Mapping

| Parent Section | Mapping |
| --- | --- |
| Workflow graph steps | §6.1 step 9 (scaffold from template) |
| Implementation plan steps | §7.4 step 1 |
| Requirements | FR-002 |
| Tests | none (docs/template authoring; covered indirectly by T-002/T-006) |

## Implementation Steps

1. Create `templates/brainstorm.md` with YAML frontmatter (`name: '{name}'`, `created: '{date}'` — no `status` field, per FR-007/NF-004 non-goal of lifecycle tracking).
2. Add a title line `# Brainstorm: {name}` and a metadata line `> **Created**: {date}`.
3. Add a short `<!-- -->` guidance comment explaining this is pre-spec exploration, filled by `/flexspec-brainstorm`, and that open questions are allowed to remain (not a Definition-of-Ready gate).
4. Add numbered sections with `{placeholder}` bodies: `1. Problem & Goal`, `2. Users & Context`, `3. Workflow & Edge Cases`, `4. Data & Interfaces`, `5. Security & Abuse Cases`, `6. Failure Handling & Concurrency`, `7. Performance & Scale`, `8. Operational Considerations`, `9. Alternatives & Tradeoffs Considered`, `10. Open Questions & Risks`, `11. Decisions & Direction`.
5. In `templates/README.md`: add a row to the templates table — `brainstorm.md` | `.flexspec/templates/` | Pre-spec exploration; scaffolded per-session by `flexspec brainstorm new <name>`, filled via `/flexspec-brainstorm`; not a feature spec.
6. In `templates/README.md`: add `flexspec brainstorm new <name> [--force]` to the CLI commands table with purpose "Create `.flexspec/brainstorms/<slug>.md` for pre-spec exploration".
7. In `templates/README.md`: add a short "Brainstorm Docs" subsection (mirroring the "Where Specs Live" subsection) stating brainstorm docs live flat under `.flexspec/brainstorms/<slug>.md`, have no status lifecycle, and are not listed by `flexspec list`/`flexspec validate`/the UI board.

## Acceptance Criteria

- [ ] `templates/brainstorm.md` exists with frontmatter + 11 numbered sections listed above. _(FR-002)_
- [ ] `templates/brainstorm.md` has no `status` frontmatter field.
- [ ] `templates/README.md` documents the new template, the CLI command, and the `.flexspec/brainstorms/` location.
- [ ] All tests in "Testing" below pass.

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| n/a | manual | `gofmt`/build still pass after adding a non-Go embedded file (sanity check only) | repo root |

Run: `go build ./...`

## Out of Scope

- Do not implement `internal/brainstorm` or the CLI command here (T-002/T-003).
- Do not modify `cmd/init.go` or `internal/migrate/templates_resync.go` — both already generically walk the embedded `templates/` tree and need no code change for a new file.
- Do not touch the root `README.md` or `AGENTS.md` (charter §8 boundary).

## Open Questions

None.

## References

- Parent spec: [`../README.md`](../README.md) §6.1, §7.1, §7.4 step 1
- Related tasks: T-002 (reads this template), T-005 (maps ingested content back onto spec sections)
