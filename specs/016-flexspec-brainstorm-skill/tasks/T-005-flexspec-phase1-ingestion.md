---
blocks:
    - T-006
depends_on:
    - T-004
id: T-005
name: /flexspec Phase 1 brainstorm ingestion step
parent_spec: ../README.md
satisfies:
    - FR-005
    - FR-006
status: done
verified_by: []
---

# T-005: /flexspec Phase 1 brainstorm ingestion step

> **Parent spec**: [Brainstorm skill and CLI scaffolding](../README.md) · **Status**: todo
> **Satisfies**: FR-005, FR-006 · **Depends on**: T-004 · **Verified by**: manual review against §6.2 trace table
> **Blocks**: T-006

## Objective

Add a brainstorm auto-detect/ingest step to `skills/flexspec/SKILL.md`'s Phase 1 workflow, inserted before the Discovery Gate runs, so an existing brainstorm doc is treated as pre-answered discovery context instead of being re-interviewed from scratch.

## Context

This is a **self-referential edit**: `skills/flexspec/SKILL.md` is the file that defines the `/flexspec` skill currently being followed to author this very spec. The change must fit into the existing "Phase 1 Workflow" numbered list (currently steps 1–15, see the live skill file) without disturbing unrelated steps, and into the "Charter Gate" / "Glossary Gate" / "Discovery Gate Scaling" sections only if this new step interacts with them (it does not modify those gates — it runs before/alongside the Discovery Gate and only changes which questions get asked).

Insert the new behavior as an explicit sub-section, e.g. "## Brainstorm Ingestion Gate (Phase 1)", placed after the Glossary Gate section and before the Discovery Gate Scaling section (both gates already read supporting files before design work begins, so this fits the same position). Reference it from "Phase 1 Workflow" step list (add a step between "Read glossary" and "Run the Discovery Gate").

Matching algorithm (parent spec §6.2, §3 alternate flows) — write this precisely so a future agent applies it consistently:
- Scan `.flexspec/brainstorms/*.md` (skip entirely, no error, if the directory doesn't exist).
- Compute the candidate slug from the user's request the same way `flexspec new` would (same slugification concept), and treat a brainstorm filename as a candidate if its slug exactly matches, is a prefix/suffix of, or shares significant keyword overlap with the request's slug.
- Also treat any brainstorm doc the user explicitly names or describes in their request (e.g. "using the export-feature brainstorm") as a candidate, taking priority over the heuristic match.
- Zero candidates: proceed with the standard Discovery Gate unchanged; do not mention brainstorms in the spec.
- Exactly one candidate: read it (read-only — never edit or delete it); treat each topic it already answers as resolved discovery input; only ask Discovery Gate questions for topics the doc leaves unanswered or unclear; cite it in the new spec's Section 2 under a line such as "Brainstorm reference: `.flexspec/brainstorms/<slug>.md`".
- Multiple candidates: ask the user (one grouped question) which doc to use, or none.
- This step runs for every type (feature/bug/chore/etc.), not just feature/bug — a brainstorm doc can precede any kind of work, but its main value is for the Exhaustive discovery types.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `skills/flexspec/SKILL.md` | modify | Add "Brainstorm Ingestion Gate" section + Phase 1 Workflow step + DoR note |
| `skills/flexspec-brainstorm/SKILL.md` | read | Confirms the doc shape/sections being ingested (from T-004) |

### Workflow / Requirement Mapping

| Parent Section | Mapping |
| --- | --- |
| Workflow graph steps | §6.2 (full flow, steps 1–10) |
| Implementation plan steps | §7.4 step 5 |
| Requirements | FR-005, FR-006 |
| Tests | none automated; verified by manual walkthrough against §6.2 |

## Implementation Steps

1. In `skills/flexspec/SKILL.md`, locate the "Phase 1 Workflow" numbered list; insert a new step after "Read `.flexspec/glossary.yaml` and note known terms" and before "Resolve `type`": "Run the Brainstorm Ingestion Gate: scan `.flexspec/brainstorms/` for a matching doc before design work begins."
2. Add a new top-level section "## Brainstorm Ingestion Gate (Phase 1)" (positioned after "Glossary Gate (Phase 1 only)" and before "Template Choice Heuristic", matching the existing gate-section style) documenting the matching algorithm from Context above verbatim (candidates, zero/one/multiple branches, read-only guarantee, citation format).
3. Update "Section 2 Reasons For Change" authoring requirements (in the "Authoring Requirements" section) to mention that when a brainstorm doc was ingested, its path is cited under a "Brainstorm reference" line, consistent with the parent spec's own Section 2.
4. Add a bullet to the relevant Definition of Ready checklist(s) (feature/bug at minimum) noting that if a brainstorm doc was ingested, its citation appears in Section 2 — this is descriptive, not a new blocking gate item (ingestion is optional; absence of a brainstorm doc is never a DoR failure).
5. Do not modify the Discovery Gate Scaling tables' question counts or areas — the gate still runs; this step only changes which of its questions have already-known answers.

## Acceptance Criteria

- [ ] `skills/flexspec/SKILL.md` contains a "Brainstorm Ingestion Gate (Phase 1)" section documenting the zero/one/multiple-match algorithm and the read-only guarantee. _(FR-005, FR-006)_
- [ ] The gate is positioned before the Discovery Gate runs in the Phase 1 Workflow list.
- [ ] Section 2 authoring requirements mention the "Brainstorm reference" citation line.
- [ ] No existing Phase 1 step, gate, or DoR checklist item was removed or renumbered incorrectly.
- [ ] All tests in "Testing" below pass.

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| n/a | manual | Live-session walkthrough: run `/flexspec-brainstorm` on a topic (T-004), then `/flexspec` on a matching request in the same or a new session; confirm the brainstorm doc is detected, cited, left unmodified, and Discovery Gate questions already answered in it are not re-asked | interactive agent session |

Run: manual verification only (skill-instruction file; no automated harness in this repo).

## Out of Scope

- Do not build any Go/CLI matching logic — the matching algorithm is agent-instruction text executed by the LLM at authoring time, not code.
- Do not change how `/flexspec` handles specs that have no brainstorm doc (zero-candidate path is a no-op, must not alter existing behavior).
- Do not touch Phase 2 (Implement) or Phase 3 (Review) sections — ingestion is Phase 1 only, matching the existing rule that glossary interviews are also Phase-1-only.

## Open Questions

None.

## References

- Parent spec: [`../README.md`](../README.md) §3, §6.2, §7.4 step 5
- Related tasks: T-004 (produces the docs this step ingests), T-006 (verifies end-to-end)
- File being edited: `skills/flexspec/SKILL.md` (this very skill definition)
