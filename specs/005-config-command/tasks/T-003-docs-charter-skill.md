---
id: T-003
name: Docs charter and skill
parent_spec: '../README.md'
status: done
satisfies: [FR-008, NF-002]
depends_on: [T-002]
verified_by: [TC-006]
---

# T-003: Docs, charter, and skill

> **Parent spec**: [Config command](../README.md) · **Status**: todo
> **Satisfies**: FR-008, NF-002 · **Depends on**: T-002 · **Verified by**: TC-006

## Objective

Record `flexspec config` everywhere humans and agents discover CLI commands; steer agents away from manual `config.yaml` reads.

## Context

Mirror T-012 from spec 002: update repo `README.md`, `.flexspec/charter.md` §4, `templates/README.md`, `.flexspec/templates/README.md`, and `skills/flexspec/SKILL.md` CLI table + Phase 1/2 guidance ("run `flexspec config` or `flexspec config --json`" instead of opening `.flexspec/config.yaml`).

Do **not** add example config keys or `comment_aggressiveness`.

## Files In Scope

- `README.md`
- `.flexspec/charter.md`
- `templates/README.md`
- `.flexspec/templates/README.md`
- `skills/flexspec/SKILL.md`

## Implementation Steps

1. README usage table: row for `flexspec config` / `--json`; one-line example in bash block.
2. Charter §4 Available today: add `flexspec config` (`--json`).
3. Both template README CLI tables: same row.
4. Skill: add to mandatory CLI table; short note under config resolution (replace "read config.yaml" with `flexspec config --json` for `always_one_shot` / `spec_template` checks).
5. Charter §11 revision row if project convention requires it for §4 edits.

## Acceptance Criteria

- [ ] All five files mention the command consistently
- [ ] Skill explicitly prefers CLI over manual YAML for agents
- [ ] No new config keys documented beyond existing three

## Testing

- Manual: grep `flexspec config` in listed files
- `go test` still green

## Out of Scope

- `flexspec-charter` skill (no config command needed there unless charter interview references config — skip unless already listing all CLI commands)

## Open Questions

(none)

## References

- `specs/002-management-ui/tasks/T-012-docs-charter-skills.md`
- Parent FR-008
