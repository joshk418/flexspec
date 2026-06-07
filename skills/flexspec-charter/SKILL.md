---
name: flexspec-charter
description: >
  Interview the user and create or update .flexspec/charter.md for application-wide
  context. Use when the user runs /flexspec-charter or asks to define product vision,
  scope, stack, or conventions before writing FlexSpec specs.
---

# FlexSpec Application Charter (`/flexspec-charter`)

The charter is the **durable product context** for every FlexSpec spec. It captures
vision, capabilities, technical constraints, and boundaries at the application level —
not individual features.

Template: embedded at `templates/charter.md`, scaffolded to `.flexspec/charter.md` by
`flexspec init`. Structure reference: same file after init.

## Invocation

This skill **auto-triggers** from ambient context (e.g. "define our product vision",
"update the application charter") — same as `/flexspec`. No `disable-model-invocation`.

## The Most Important Rule: Ask, Don't Assume

Do not invent product facts, stack choices, or boundaries. Batch 2–4 related questions
per round. If the user defers an answer, record it in §10 as blocking or non-blocking
and do not mark the charter `active` while blocking items remain.

Use the integrated structured question system available in the current agent runtime
when it exists (for example Cursor `AskQuestion`, or equivalent Claude, Codex, or
other agent multiple-choice tools). Prefer grouped multiple-choice prompts; use
multi-select when several answers can apply. If no structured question tool exists,
ask the same choices in concise text and record any deferred or ambiguous answers
in §10.

---

## Workflow

### 1. Prerequisites

- If `.flexspec/` is missing, run `flexspec init` (or ask the user to).
- Charter path: **`.flexspec/charter.md`** (never under `.flexspec/templates/`).

### 2. Load and classify state

Read the existing charter and classify:

| State | Sentinels |
| --- | --- |
| **Empty** | Zero-length or whitespace only. |
| **Template-only** | Contains `{` placeholders or `<!--` guidance comments, or frontmatter `status: draft`. |
| **Active** | Frontmatter `status: active`, no `{` placeholders, no `<!--` hints. |

### 3. Choose interview mode

- **Create** — empty or template-only: walk §1–§10 in order.
- **Update** — active charter: re-interview only empty, stale, or user-requested sections.
- **Triggered update** — `/flexspec` handed off a delta list (sections + bullets + spec slug):
  interview only those sections, merge, append §11 row citing the spec (e.g. `001-user-auth`).
- **Full refresh** — only when the user explicitly asks to redo the whole charter.

### 4. Interview (question bank by section)

Map answers into the template sections. Sample questions:

| Section | Questions |
| --- | --- |
| §1 Product overview | Product name? One-liner? What problem does it solve? What outcome defines success? |
| §2 Vision and goals | North star in one sentence? 2–3 measurable success criteria? |
| §3 Users | Who are the primary personas? What jobs must the product enable? |
| §4 Capabilities | What major capability domains exist today? What domains are planned? For FlexSpec itself, include CLI commands (`init`, `new`, `list`, `validate`) and skills. |
| §5 Technical context | Stack, hosting, key integrations, hard constraints (languages, compliance)? |
| §6 Architecture | Main components and boundaries? Optional: sketch data/control flow. |
| §7 Standards | Testing expectations? Security defaults? Naming/repo conventions agents must follow? |
| §8 Product boundaries | What will this product **not** do (global non-goals)? |
| §9 Glossary | Domain terms agents must use consistently? |
| §10 Assumptions / questions | What are we assuming? What is still unknown (blocking vs not)? |

### 4.1 Glossary discovery handoff

During charter creation, full refresh, or updates that touch §1, §3, §4, §5, §6,
§7, §8, or §9, run the glossary discovery workflow before final handoff.

- Invoke or follow `flexspec-glossary-discovery` after the main charter interview so the glossary is built from the same product language.
- Use `flexspec glossary list --json` to load known terms before discovery.
- Persist clear or user-confirmed terms only through `flexspec glossary add`.
- If a term is project-specific but unclear, ask the user for the exact meaning before persisting.
- Do not block charter activation on optional glossary additions unless an unclear term is needed to understand §1-§8.
- Keep `flexspec-glossary-discovery` available as a standalone skill for later glossary refreshes or manual updates.

For charter creation, charter refresh, or `/flexspec`-triggered deltas involving
UI, ask a short structured UI standards round before writing §4, §7, §8, and §10.
Keep it product-wide; do not embed individual feature specs.

| UI area | Multiple-choice prompt examples |
| --- | --- |
| Product style | existing app style, polished SaaS, playful/illustrated, minimal utility, dense admin, custom |
| Design-system source | established component library, local components, bespoke per feature, undecided |
| Interaction defaults | icons on buttons, password show/hide toggles, inline validation, skeleton loading, toast/banner feedback |
| State requirements | empty, loading, error, success, permission/unauthenticated, responsive/mobile |
| Accessibility baseline | WCAG AA contrast, visible focus, labels/help text, reduced motion, screen-reader text |
| Boundaries | avoid new visual language, avoid animation, avoid new dependencies, require user approval for major redesigns |

When `/flexspec` proposes a UI delta, treat the delta list as the agenda and ask
only the missing product-level UI choices. Preserve existing charter sections that
were not part of the requested update.

### 5. Write the charter

- Replace all `{placeholders}`; remove all `<!-- -->` guidance comments.
- Set frontmatter: `product_name`, `version`, `last_updated` (ISO date), `status`.
- Set `status: active` when §1–§8 are sufficiently filled and §10 has **no blocking** open questions.
- Keep total size ~1500–2500 tokens; link to repo paths instead of pasting large trees.
- **Do not** embed individual feature spec content — only product-level facts.
- On update runs, **preserve** user edits in sections you did not re-interview.
- Append a §11 revision row: date, short summary, source (`/flexspec-charter` or spec slug).

### 6. Completion checklist

Before handoff, confirm:

- [ ] Product name and one-liner populated (§1).
- [ ] Vision / north star and success criteria (§2).
- [ ] At least one persona (§3).
- [ ] Capability map reflects current product scope (§4).
- [ ] Technical context and constraints (§5).
- [ ] Product boundaries / non-goals (§8).
- [ ] Glossary discovery run or intentionally skipped with reason.
- [ ] No blocking open questions in §10.
- [ ] Frontmatter `status: active`, no `{` or `<!--` left in body.

### 7. Handoff

Tell the user to run `/flexspec` for feature specs and `flexspec validate` to confirm
the project structure is sound. Optionally mention `flexspec ui` to browse the charter
and specs in a local browser (read-only for charter in v1). The charter is now the
application context source for Phase 1 authoring.

---

## Rules

- Charter is **product-wide**; feature specs live under `<specs_dir>/NNN-slug/`.
- Prefer tables and bullets over long prose.
- If `.flexspec/charter.md` is missing but templates exist, run `flexspec init` before interviewing.
- When `/flexspec` proposes charter deltas, treat the delta list as the interview agenda — do not re-ask about unchanged sections unless the user wants a full refresh.
- When charter work introduces domain terminology, build or refresh `.flexspec/glossary.yaml` through `flexspec-glossary-discovery`; keep the discovery skill runnable on its own for manual glossary maintenance.

---

## Section reference (template §1–§11)

| § | Purpose |
| --- | --- |
| 1 | Product overview — problem, outcome, one-liner |
| 2 | Vision and goals — north star, success criteria |
| 3 | Users and stakeholders — personas, jobs-to-be-done |
| 4 | Capabilities — product-level feature domains |
| 5 | Technical context — stack, deployment, integrations, constraints |
| 6 | Architecture — optional mermaid, boundaries |
| 7 | Standards and conventions — testing, security, patterns |
| 8 | Product boundaries — global non-goals |
| 9 | Domain glossary |
| 10 | Assumptions and open questions |
| 11 | Revision history — dated changelog |
