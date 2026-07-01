---
name: flexspec-brainstorm
description: >
  Interview the user to explore a feature idea before spec authoring — problem
  framing, edge cases, security, performance, alternatives, and open questions.
  Persists a brainstorm doc under .flexspec/brainstorms/ that /flexspec Phase 1
  can auto-detect and ingest. Use when the user runs /flexspec-brainstorm or
  wants to think through an idea before planning.
---

# FlexSpec Brainstorm (`/flexspec-brainstorm`)

Pre-spec exploration. This skill helps a user think through a feature idea —
goals, edge cases, security, performance, alternatives — **before** `/flexspec`
Phase 1 authoring. The output is a durable doc, not a chat transcript, so the
thinking survives `/clear`, compaction, or a new session.

Template: embedded at `templates/brainstorm.md`, scaffolded to
`.flexspec/templates/brainstorm.md` by `flexspec init` (or `flexspec update
--migrate` for existing projects). Scaffold each session's doc with
`flexspec brainstorm new <name> [--force]`.

## Invocation

This skill **auto-triggers** from ambient context ("let's brainstorm X", "help
me think through Y before I spec it", "brainstorm a feature") — same as
`/flexspec-charter` and `/flexspec-migrate`. No `disable-model-invocation`.
Also invoked explicitly via `/flexspec-brainstorm [topic]`.

## Core Rules

1. **Ask, don't assume — but this is exploratory, not gating.** Unlike
   `/flexspec`'s Definition of Ready, a brainstorm doc may end with unresolved
   items in "Open Questions & Risks". The goal is to think things through, not
   reach a final decision. Never invent a brainstorm "status" or block anything
   on incomplete sections.
2. **CLI-only scaffolding.** Always create the file via `flexspec brainstorm
   new <name>`; never hand-write `.flexspec/brainstorms/*.md`.
3. **Read-only charter and glossary.** Read `.flexspec/charter.md` and
   `.flexspec/glossary.yaml` for product context; never write to either file
   from this skill.
4. **Write only inside `.flexspec/brainstorms/`.** Never touch `<specs_dir>`,
   `README.md`, `AGENTS.md`, or any other project file.
5. **Never silently overwrite.** If `.flexspec/brainstorms/<slug>.md` already
   exists for the resolved slug, ask the user whether to continue/refine that
   doc or pick a different name before calling `flexspec brainstorm new
   --force`.
6. **No status/lifecycle tracking.** Brainstorm docs carry no frontmatter
   `status` field and are never registered anywhere `flexspec list` or the
   management UI would surface them.

## Workflow

1. Resolve the topic from the user's request (slash-command argument or
   ambient phrase).
2. Read `.flexspec/charter.md` and `.flexspec/glossary.yaml` for context only
   (skip silently if `.flexspec/` doesn't exist yet — a brainstorm can precede
   `flexspec init`; if so, tell the user to run `flexspec init` before the
   scaffold step below).
3. Resolve a slug from the topic (same slugification concept as `flexspec
   new`).
4. Check whether `.flexspec/brainstorms/<slug>.md` already exists. If it does,
   ask the user: continue/refine that doc, or start a new one under a
   different name (Core Rule 5).
5. Run `flexspec brainstorm new <slug>` (add `--force` only after the user
   confirms continuing/refining the existing doc).
6. Interview the user across the template's topic sections (see Question bank
   below), using the runtime's structured question tool when available;
   prefer grouped multiple-choice questions, multi-select where several
   answers can apply. Batch related questions.
7. When the interview surfaces a candidate library, package, or service (see
   Research below), research it before writing that section.
8. Write synthesized findings — not a raw transcript — into the corresponding
   template section as the interview progresses.
9. Hand off: tell the user the doc path, note it has no status/lifecycle
   tracking, and that running `/flexspec` next will auto-detect and ingest it
   (citing it in the new spec's Section 2) without re-asking already-answered
   topics.

## Question bank (by template section)

| Section | Example questions |
| --- | --- |
| 1. Problem & Goal | What problem are you solving? What does success look like? Who is this for? |
| 2. Users & Context | Who are the actors? What's the entry point? What's the current workaround, if any? |
| 3. Workflow & Edge Cases | What's the happy path? What alternate/error paths matter? What inputs could break it? |
| 4. Data & Interfaces | What data is involved? Any new APIs, routes, CLI flags, or events? Any external contracts? |
| 5. Security & Abuse Cases | Who's authorized? What could a malicious or careless user do? Any secrets, PII, or injection surfaces? |
| 6. Failure Handling & Concurrency | What happens on a dependency outage or timeout? Any races or shared state? |
| 7. Performance & Scale | Expected input size / throughput? Any latency or memory constraints? |
| 8. Operational Considerations | Logging, metrics, rollout, migration, or rollback concerns? |
| 9. Alternatives & Tradeoffs Considered | What other approaches, libraries, or services could work? (see Research below) |
| 10. Open Questions & Risks | What's still unresolved? What are you unsure about? (fine to leave open) |
| 11. Decisions & Direction | What's the current leaning, even if tentative? |

## Research

When the user's own answers raise a specific candidate library, package, or
service they're considering (their words are the trigger — no separate
permission prompt, unlike `/flexspec-migrate`'s tool-detection-failure
trigger), use whatever web research tool the runtime already exposes (e.g.
`WebSearch`/`WebFetch`) to look it up and compare it against alternatives the
user mentioned or that are evident from repo context. Write findings plus
sources into "9. Alternatives & Tradeoffs Considered".

- Build no new research/scraping tooling — use only what the runtime already
  provides.
- If no web research tool is available, say so in that section and ask the
  user to paste in whatever comparison info they have. Never invent library
  capabilities, pricing, or comparison claims.
- Don't research every term mentioned in passing — only concrete
  libraries/packages/services the user is actively considering.

## Handoff

Tell the user:

- The brainstorm doc's path (`.flexspec/brainstorms/<slug>.md`).
- That it has no status/lifecycle tracking and won't appear in `flexspec
  list` or the management UI board.
- That running `/flexspec` next will scan `.flexspec/brainstorms/` for a
  matching doc, ingest it read-only, cite it in the new spec's Section 2, and
  only ask about topics not already answered in it.

---

## Rules recap

- Charter and glossary are **read-only** context for this skill — never
  written to from here.
- This skill writes **only** inside `.flexspec/brainstorms/`.
- Scaffolding is **always** via `flexspec brainstorm new`, never hand-created.
- Research uses only existing runtime tools, triggered by the user's own
  words, and never fabricates findings when no tool is available.
- Unresolved items are acceptable — this is exploration, not a
  Definition-of-Ready gate.
