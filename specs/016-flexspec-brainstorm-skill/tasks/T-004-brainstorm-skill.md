---
blocks:
    - T-005
depends_on:
    - T-003
id: T-004
name: /flexspec-brainstorm skill
parent_spec: ../README.md
satisfies:
    - FR-003
    - FR-004
    - FR-007
    - FR-008
status: done
verified_by: []
---

# T-004: /flexspec-brainstorm skill

> **Parent spec**: [Brainstorm skill and CLI scaffolding](../README.md) · **Status**: todo
> **Satisfies**: FR-003, FR-004, FR-007, FR-008 · **Depends on**: T-003 · **Verified by**: manual review against §6.1 trace table
> **Blocks**: T-005

## Objective

Author `skills/flexspec-brainstorm/SKILL.md`: an agent skill that runs an in-depth, non-blocking pre-spec interview and persists findings to a brainstorm doc scaffolded by `flexspec brainstorm new`.

## Context

Model this on `skills/flexspec-charter/SKILL.md` (interview workflow, structured-question-tool preference, section-by-section question bank) and `skills/flexspec-migrate/SKILL.md` (Core Rules block, CLI-only scaffolding rule, `Non-destructive default` rule). Both auto-trigger from ambient phrases with no `disable-model-invocation` frontmatter field — this skill follows the same convention.

Key differences from `/flexspec`'s Discovery Gate: this skill is **exploratory, not gating**. Unlike `/flexspec`'s Definition of Ready (which blocks `planned` on open questions), a brainstorm doc is allowed to end with unresolved items in its "Open Questions & Risks" section — the point is to think things through, not to reach a final decision. Do not invent a "brainstorm status" or block anything.

The skill must respect these boundaries (parent spec §2, §3, FR-004):
- Reads `.flexspec/charter.md` and `.flexspec/glossary.yaml` for context only — never writes to either.
- Writes only inside `.flexspec/brainstorms/` (via the CLI command from T-003) — never touches `<specs_dir>`, `README.md`, or `AGENTS.md`.
- Never auto-deletes or silently overwrites an existing brainstorm doc — if `.flexspec/brainstorms/<slug>.md` already exists for the resolved slug, ask the user whether to continue that doc (read + re-interview / append) or pick a different name before calling `flexspec brainstorm new --force`.
- No frontmatter `status` field is ever added or managed (FR-007).

Interview topics mirror `templates/brainstorm.md`'s 11 sections (T-001): problem & goal, users & context, workflow & edge cases, data & interfaces, security & abuse cases, failure handling & concurrency, performance & scale, operational considerations, alternatives & tradeoffs considered, open questions & risks, decisions & direction. Ask grouped multiple-choice/short-answer questions using the runtime's structured question tool when available (same guidance as `flexspec-charter`); batch related questions; write synthesized findings (not a raw transcript) into the corresponding template section as the interview progresses.

**Research behavior (FR-008, parent spec §6.1 steps 11–16):** when the user's own answers raise a specific candidate library, package, or service (the user's words are the trigger — no separate permission prompt needed, unlike `flexspec-migrate`'s web-lookup-permission rule), use whatever web research tool the runtime already exposes (e.g. `WebSearch`/`WebFetch`) to look it up and compare it against alternatives the user mentioned or the skill is aware of from repo context. Write findings plus sources into the "Alternatives & Tradeoffs Considered" section. If no such tool is available in the current runtime, say so and ask the user to paste in whatever comparison info they have — never invent library capabilities, pricing, or comparison claims. Do not research every generic term; only concrete candidates the user is actually considering.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `skills/flexspec-brainstorm/SKILL.md` | create | New skill definition |
| `skills/flexspec-charter/SKILL.md` | read | Interview-workflow and structured-question-tool pattern reference |
| `skills/flexspec-migrate/SKILL.md` | read | Core Rules block and CLI-only-scaffolding pattern reference |
| `templates/brainstorm.md` | read | Section list this skill's interview must cover (from T-001) |

### Workflow / Requirement Mapping

| Parent Section | Mapping |
| --- | --- |
| Workflow graph steps | §6.1 (full flow, steps 1–18, including research branch steps 11–16) |
| Implementation plan steps | §7.4 step 4 |
| Requirements | FR-003, FR-004, FR-007, FR-008 |
| Tests | none automated; verified by manual walkthrough against §6.1 |

## Implementation Steps

1. Create `skills/flexspec-brainstorm/SKILL.md` with frontmatter `name: flexspec-brainstorm` and a `description` stating it interviews the user to explore a feature idea before spec authoring (for skill-matching/auto-trigger discovery), matching the style of existing skill frontmatter descriptions.
2. Add an "Invocation" section stating this skill auto-triggers from ambient phrases ("let's brainstorm X", "help me think through Y before I spec it", "brainstorm a feature") in addition to `/flexspec-brainstorm [topic]`, with no `disable-model-invocation`.
3. Add a "Core Rules" block: (a) Ask, don't assume — but unresolved items are fine, this is exploration not a gate; (b) CLI-only scaffolding — always create the file via `flexspec brainstorm new <name>`, never hand-write `.flexspec/brainstorms/*.md`; (c) read-only charter/glossary — never write to either; (d) writes only inside `.flexspec/brainstorms/`; (e) never silently overwrite an existing doc.
4. Add a "Workflow" section with numbered steps matching parent spec §6.1: resolve topic -> read charter/glossary for context -> resolve slug -> check for existing doc (ask if found) -> run `flexspec brainstorm new <slug>` -> interview across the 11 template sections -> write synthesized findings into the file -> handoff summary.
5. Add a "Question bank" table (one row per template section) with example question prompts per topic, mirroring `flexspec-charter`'s §4 interview question-bank table style — explicitly include security/abuse-case questions and performance/scale questions per the original user request.
6. Add a "Research" subsection (under the Alternatives & Tradeoffs row of the question bank, or as its own short section) documenting the FR-008 behavior verbatim from Context above: trigger (user raises a concrete candidate), tool usage (`WebSearch`/`WebFetch` when available, no new tooling), fallback (no tool available -> ask user, never fabricate), and output location (Alternatives & Tradeoffs Considered section, with sources).
7. Add a "Handoff" section: tell the user the brainstorm doc path, note it has no status/lifecycle tracking, and that running `/flexspec` next will auto-detect and ingest it (citing it in the new spec's Section 2) without re-asking already-answered topics.
8. Add a short "Rules" recap section (mirroring `flexspec-charter`'s trailing Rules block) restating the read-only charter/glossary boundary, the `.flexspec/brainstorms/`-only write boundary, and the research trigger/fallback rule.

## Acceptance Criteria

- [ ] Skill frontmatter and description follow existing skill conventions (checked against `flexspec-charter`/`flexspec-migrate`). _(FR-003)_
- [ ] Skill explicitly states it reads charter/glossary for context only and never writes to them. _(FR-004)_
- [ ] Skill explicitly states no `status` frontmatter field is added to brainstorm docs. _(FR-007)_
- [ ] Skill always scaffolds via `flexspec brainstorm new`, never hand-creates the file.
- [ ] Skill covers all 11 topic areas from `templates/brainstorm.md`, including explicit security and performance question prompts.
- [ ] Skill documents the research trigger (user raises a concrete library/package/service), tool usage, no-tool fallback, and output location (Alternatives & Tradeoffs section with sources). _(FR-008)_
- [ ] Skill states it never fabricates research findings when no web tool is available.
- [ ] All tests in "Testing" below pass.

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| n/a | manual | Live-session walkthrough: invoke `/flexspec-brainstorm` on a sample topic, confirm interview covers all sections, confirm the resulting file matches §6.1's trace table and never touches charter/glossary/specs_dir | interactive agent session |

Run: manual verification only (no automated harness for skill instruction files in this repo, consistent with `flexspec-charter`/`flexspec-migrate`).

## Out of Scope

- Do not build new research/scraping tooling or add dependencies — research uses whatever `WebSearch`/`WebFetch`-style tools the runtime already provides.
- Do not research every term mentioned in passing — only concrete libraries/packages/services the user is actively considering.
- Do not add a "resume brainstorm" CLI flag — resuming an existing doc is skill-level behavior (read + re-interview), not a new CLI capability.
- Do not modify `.flexspec/charter.md`, `.flexspec/glossary.yaml`, or any file under `<specs_dir>`.

## Open Questions

None.

## References

- Parent spec: [`../README.md`](../README.md) §3, §6.1, §7.4 step 4
- Related tasks: T-001 (template sections this interviews against), T-003 (CLI command this invokes), T-005 (`/flexspec` side of the handoff)
- Pattern reference: `skills/flexspec-charter/SKILL.md`, `skills/flexspec-migrate/SKILL.md`
