---
name: flexspec
description: >
  Drive the full FlexSpec spec-driven development lifecycle, invoked with /flexspec.
  Use when the user runs /flexspec or asks to create, draft, refine, implement, or
  review a FlexSpec spec. The skill is a status-driven, three-phase state machine —
  author the spec, implement it, then review the diffs for spec coverage and AI
  "slop" — advancing one phase per prompt (unless --one-shot). Covers choosing
  between the simple and expanded templates, where specs are written on disk, each
  section's meaning and required format (mermaid, FR/NF/T/TC IDs), the per-task file
  format for expanded specs, the mandatory rule of resolving every unknown with the
  user before implementation, and a slop-detection review checklist.
---

# FlexSpec Lifecycle (`/flexspec`)

A spec is the agreed definition of a feature, written **before** code. This skill
drives the whole lifecycle — authoring the spec, implementing it, and reviewing the
result — using the spec's `status` field as the state.

Templates: `templates/flexspec-simple.md`, `templates/expanded/flexspec-expanded.md`,
`templates/expanded/flexspec-expanded-task.md`. Structure/metadata reference:
`templates/README.md`.

## How `/flexspec` Runs: One Phase Per Prompt

The lifecycle has **three phases**. By default, each `/flexspec` invocation does
**exactly one phase, then stops** and hands back to the user. This is deliberate:
splitting the work across separate prompt sessions keeps context fresh, prevents the
agent from racing ahead, and gives the user a checkpoint between defining, building,
and reviewing.

**Determine the current phase from the active spec's `status`:**

| Phase | Invoke when status is | Does | Ends at status |
| --- | --- | --- | --- |
| **1 · Author** | none yet / `initial` / `refined` | Scaffold the spec dir, author + refine the spec, resolve all unknowns. | `planned` |
| **2 · Implement** | `planned` / `in_progress` | Write the code/tasks defined in the spec. | `in_review` |
| **3 · Review** | `in_review` | Diff-review against the spec + scan for AI slop. | `complete` |

At the **end of each phase**, update the spec `status`, summarize what was done, and
**ask the user whether to continue to the next phase** (e.g. "Spec is ready — run
`/flexspec` again to start implementation?"). Do not roll straight into the next
phase.

### `--one-shot` Mode

If the user invokes `/flexspec --one-shot` (handy for Claude / unattended runs), do
**not** stop between phases — run Author → Implement → Review back-to-back to
completion in a single session.

- One-shot still **must resolve all unknowns first** (Phase 1). If blocking questions
  remain that you cannot answer from the repo/request, stop and ask anyway —
  correctness beats autonomy. Once answered, continue automatically.
- Update `status` as each phase completes and report progress between phases, but keep
  going without waiting for a prompt.

### `always_one_shot` Config

Before deciding the run mode, read `always_one_shot` from `.flexspec/config.yaml`.
If it is `true`, behave as if `--one-shot` were passed even when the flag is
absent — run all phases back-to-back. An explicit `--one-shot` flag also forces
one-shot mode regardless of config.

If no flag is given **and** `always_one_shot` is `false` (or config is missing),
default to the one-phase-per-prompt behavior above.

### `spec_template` Config and `--template` Flag

By default `/flexspec` **infers** which template to use (simple vs expanded) from
the work's size — see "Choosing a Template". Two overrides exist, highest precedence
first:

1. **`--template <simple|expanded>` flag** — an explicit `/flexspec --template expanded`
   forces that template for this run, overriding the config and inference.
2. **`spec_template` in `.flexspec/config.yaml`** — if set to `simple` or `expanded`,
   use it. It has **no default**: when the key is blank or absent, fall back to
   inference so the normal workflow still chooses simple or expanded on its own.

Only `simple` and `expanded` are valid; treat any other value as unset and infer (or
ask the user if borderline).

---

# Phase 1 · Author the Spec

Goal: produce a complete, unambiguous, testable spec on disk and advance status to
`planned`. This is the authoring half of the skill; everything from here through the
"Definition of Ready" checklist is Phase 1 detail.

## The Most Important Rule: Resolve All Unknowns First

Writing the spec is secondary. **Eliminating ambiguity is the primary job.** A spec
built on guesses produces wrong implementations.

- As you draft, surface every unknown, assumption, ambiguity, and decision point.
- **Ask the user direct questions** for anything you cannot determine with certainty
  from the repo or their request. Batch related questions; do not drip one at a time.
- Do **not** invent requirements, file paths, behaviors, or constraints to fill gaps.
  If you are unsure, ask.
- A spec **cannot** be marked `refined` (or advance further) and **implementation
  must not start** while any open question remains (Section 5 of the spec, or a
  task's Open Questions).
- Only once the user has answered everything and no blocking questions remain is the
  spec ready.

When in doubt: ask, don't assume.

## Application charter

Every feature spec must align with the **application charter** at `.flexspec/charter.md`.
The charter holds product-wide vision, capabilities, constraints, and boundaries — not
individual feature details.

- **Always read** `.flexspec/charter.md` at the start of Phase 1 (before repo exploration).
- **Classify charter state** using the same sentinels as `/flexspec-charter`:
  - **Empty** — zero-length or whitespace only.
  - **Template-only** — `{` placeholders, `<!--` guidance comments, or frontmatter `status: draft`.
  - **Active** — `status: active`, no placeholders or guidance comments.
- **If missing, empty, or template-only**: stop and recommend `/flexspec-charter` first.
  Offer a minimal inline charter pass only if the user **insists** on proceeding without it.
- **When authoring specs**: align Summary, NFRs, architecture, and conventions with the
  charter; cite charter sections where helpful. **Flag conflicts** between the feature
  request and charter §8 (boundaries) or §7 (standards) — ask which wins.
- **Charter freshness (required)** — During and after spec authoring, detect material
  deltas vs. the charter: §2 vision/goals, §3 users, §4 capabilities, §5–§6 stack/architecture,
  §7 conventions, §8 boundaries, §9 glossary.

## Where Specs Are Written

Specs are created on disk, not just in chat. The `flexspec` CLI owns scaffolding.

- **Project setup** — `flexspec init` bootstraps a project: it creates `.flexspec/`
  (with `config.yaml` and `charter.md`) and the `templates/` directory. If a project
  isn't initialized yet, run `flexspec init` first.
- **Location** — the specs directory is defined by `specs_dir` in
  `.flexspec/config.yaml`. Resolve it from there; never hard-code a path. If config
  is missing/unreadable, run `flexspec init` or **ask the user**.
- **One directory per spec** — `flexspec new` creates the spec folder named with a
  zero-padded auto-incrementing sequence number and a short slug: `NNN-<spec_name>`
  (e.g. `001-user-auth`). The number auto-increments from existing folders.
- **The spec is `README.md`** inside that folder. For expanded specs the CLI also
  creates the `tasks/` directory.

Simple spec layout:

```
<specs_dir>/001-feature-slug/
  README.md        <- filled from flexspec-simple.md
```

Expanded spec layout:

```
<specs_dir>/002-feature-slug/
  README.md        <- filled from flexspec-expanded.md
  tasks/
    T-001-<slug>.md  <- filled from flexspec-expanded-task.md
    T-002-<slug>.md
    ...
```

## Choosing a Template

Pick the template that matches the size and surface area of the work. **An explicit
`--template` flag or `spec_template` config wins** (see the override rules above); if
the user forces simple or expanded by either means, honor it. Otherwise infer:

**Use the simple template** for small, focused changes — typically one file/area, no
new architecture:
- Bug fixes, copy/text changes, styling tweaks (e.g. button size).
- Adding a single form to an existing page.
- Small, localized updates with few requirements.

**Use the expanded template** for large features that introduce significant new code
or span multiple layers/subsystems:
- A new set of endpoints with corresponding DB tables and frontend components.
- Adding an auth system.
- Introducing a test suite to an existing application.
- Anything with many tasks, multiple components, or cross-cutting changes.

If it is genuinely borderline, state your recommendation and **ask the user** which
they prefer before proceeding.

## Workflow

1. **Gather context.** Read `.flexspec/charter.md` first (see Application charter).
   Then read the user's request and explore the repo (relevant files, existing patterns,
   constraints). Note what you know vs. what you don't.
2. **Choose the template** (simple vs expanded) per the rules above; confirm with the
   user if borderline or if they have a preference.
3. **Scaffold the spec directory with the CLI.** Run `flexspec new` to create the
   spec folder under the specs directory from `.flexspec/config.yaml`. The CLI
   auto-increments the sequence number from existing folders, creates the
   `NNN-<spec_name>` directory with a `README.md` from the matching template, and (for
   expanded specs) a `tasks/` directory. Always use the CLI so numbering stays
   consistent — don't hand-create spec folders. If the project isn't initialized, run
   `flexspec init` first.
4. **Fill the spec files.** Populate `<folder>/README.md` from the chosen template and
   fill the frontmatter (`name`, `priority`, `tags`, `status: initial`, `created`).
   For expanded specs, author one task file per task in `tasks/`.
5. **Identify unknowns and ask.** List every open question, then ask the user. Wait
   for answers before finalizing dependent sections.
6. **Fill every section** per the guidance below, replacing all `{placeholders}` and
   removing the `<!-- -->` guidance comments.
7. **Self-check** against the Definition of Ready. If anything is only testable with
   rework, rework the implementation plan.
8. **Charter freshness check** (before `planned`):
   1. Detect concrete charter deltas implied by this spec (by section: §2, §3, §4, §5–§6, §7, §8, §9).
   2. **Deltas-only prompt:**
      - **No deltas** → state "No charter changes detected" and continue. **Do not ask.**
      - **Deltas detected** → list bullets by section, then ask: "Based on this spec, does
        `.flexspec/charter.md` need to be updated?" (yes / no / partial).
   3. **If yes or partial** (non-one-shot): recommend `/flexspec-charter` with the delta list,
      or apply targeted charter edits **only after user confirms** each change; append a §11
      revision row referencing the spec slug.
   4. **Gating (non-one-shot)** — When deltas exist, do not set `planned` until the charter
      question is answered. If the user defers, record a note in the spec's §5 Other, then proceed.
   5. **One-shot (`--one-shot` / `always_one_shot`)** — Do not block on the charter prompt:
      - Deltas that **conflict** with charter §8 or §7 are **blocking** → stop and ask even in one-shot.
      - Other deltas → record a "charter follow-up" note in spec §5 Other and continue; do not edit
        the charter unattended.
9. **Finalize the spec.** When no open questions remain, the design is agreed, the task list
   is complete, and the charter freshness check is resolved, move `status` through `refined` to `planned`.
10. **End the phase.** Summarize the spec and **ask the user to continue**: "Spec
    `NNN-slug` is planned and ready. Run `/flexspec` again to begin implementation."
    Stop here unless running `--one-shot`.

## Section-by-Section Guide (both templates)

### 1. Summary
Detailed high-level overview. State the problem, the intended outcome, who/what it
affects, and explicit scope boundaries (including what is out of scope). A reader
should understand the goal before reading the design.

### 2. Design

**Architecture / Technical Plan** — Detailed account of how the feature will be
built. Reference concrete files, packages, and components. Every relevant file goes
in the markdown table (`File / Component`, `Type` = new/modified/reference, `Role`).
Include anything an LLM implementer would need to reference.

**Code Map** — One or more `mermaid` diagrams showing components and how they relate
(data flow, call relationships, boundaries). Always real mermaid in a
```` ```mermaid ```` block, never prose substitutes.

**Expanded only — Data Model** — Schemas/tables/entities created or changed, with an
`erDiagram` and a table of fields/keys/migrations. State explicitly if no persistent
data is touched.

**Expanded only — External Interfaces** — API endpoints, CLI commands, events, UI
routes/components, and third-party integrations the feature exposes or consumes.

**Requirements** — Functional (`FR-001`…, what the system must do) and Non-Functional
(`NF-001`…, performance/security/reliability/UX). Specific and verifiable — each
should map to a test later.

### 3. Implementation Plan

**Implementation Code Map** — A `mermaid` diagram showing how the pieces/tasks build
on each other (dependencies / order).

**Task List**
- *Simple:* a list where each task has a stable ID (`T-001`…) and cites the
  requirement(s) it satisfies, e.g. `_(satisfies: FR-001, NF-001)_`. Tasks stay in
  the single file.
- *Expanded:* an index table (`Task`, `File`, `Satisfies`, `Depends on`, `Summary`).
  Each task is authored as its own file under `tasks/` (see task guide below). Keep
  tasks small enough that an LLM can complete one without losing context.

### 4. Testing Criteria
Every piece of functionality must be testable. Define tests proving each requirement
(`TC-001`…), mapping each to the requirement it `Verifies` (and, for expanded, the
task that implements it). **If something cannot be tested, the implementation plan is
wrong — rework it until all functionality is testable.** Every FR should be covered by
at least one TC.

### 5. Other
Open questions, assumptions, risks, rollout/migration notes, observations. This is
where unresolved items live *during* drafting. Before finalizing, every **open
question must be answered** (move the decision into the relevant section). Remaining
items should be non-blocking notes only.

## Authoring Task Files (expanded specs)

Each task file (`tasks/T-XXX-<slug>.md`) must let an LLM execute that task standalone,
without re-reading the whole codebase and without drifting into other tasks.

**Keep each task file under ~1000 tokens (roughly 130 lines / 750 words).** This
keeps large task lists token-efficient and each unit of work focused. If a task
cannot be described that tightly, it is too big — split it into multiple tasks
(`T-00X`, `T-00Y`) rather than bloating one file. Be terse: link to parent-spec
sections instead of restating them, keep tables minimal, and cut prose.

Fill:

- **Frontmatter** — `id`, `name`, `parent_spec`, `status`, `satisfies` (FR/NF IDs),
  `depends_on` (task IDs), `verified_by` (TC IDs).
- **Objective** — one or two sentences; what "done" means.
- **Context** — enough background, patterns, and constraints to start cold. Link to
  parent spec sections for the big picture, but keep the task self-contained.
- **Files In Scope** — only the files this task reads/changes.
- **Implementation Steps** — ordered, concrete, literal steps referencing exact
  functions/types/files. No unresolved design decisions (those go to Open Questions).
- **Acceptance Criteria** — objectively verifiable checklist tied to FR/TC IDs.
- **Testing** — the TC tests this task satisfies, where they live, and how to run.
- **Out of Scope** — nearby work this task must NOT do.
- **Open Questions** — must be empty before the task starts; ask the user to resolve.
- **References** — parent spec and related tasks.

## Format Rules

- Replace all `{placeholders}`; remove `<!-- -->` guidance comments from final files.
- IDs are stable and never reused: `FR-`, `NF-`, `T-`, `TC-`.
- Mermaid sections must contain valid `mermaid` fenced blocks.
- Cross-reference IDs across sections (tasks cite requirements, tests cite
  requirements and tasks).
- Keep frontmatter valid YAML; set `status` to reflect reality.
- The spec file is always `README.md` inside the `NNN-feature-slug` folder.
- **Token budgets** — keep the spec `README.md` token-efficient: a **simple spec
  under ~2000 tokens**, an **expanded spec under ~3500 tokens**, and each expanded
  **task file under ~1000 tokens**. If a spec exceeds its budget, it is likely
  scoped too large (split it) or too verbose (tighten prose, lean on tables and
  references rather than restating context).

## Definition of Ready (pre-implementation checklist)

- [ ] Correct template chosen (simple vs expanded) for the work's size.
- [ ] Spec written to `<specs_dir>/NNN-feature-slug/README.md`.
- [ ] All `{placeholders}` and guidance comments removed.
- [ ] Spec `README.md` within budget (simple <~2000 tokens, expanded <~3500).
- [ ] Summary states scope and out-of-scope.
- [ ] Architecture file table lists every touched/referenced file.
- [ ] All required mermaid maps present and valid (Design + Implementation; data
      model for expanded if data is touched).
- [ ] FR/NF requirements are specific and each is testable.
- [ ] Every task has a `T-` ID and cites the requirement(s) it satisfies.
- [ ] Expanded: each task has its own file in `tasks/`, under ~1000 tokens, with no
      open questions.
- [ ] Every functional requirement is covered by at least one `TC-` test.
- [ ] No open/blocking questions remain anywhere.
- [ ] Charter read; spec does not contradict charter §7 (standards) or §8 (boundaries).
- [ ] Charter freshness checked — if deltas exist, update question asked and resolved (or deferred in §5 Other).
- [ ] `status` set to `planned` (Phase 1 complete).

---

# Phase 2 · Implement the Spec

Invoked with `/flexspec` when status is `planned` (or `in_progress` to resume). Goal:
write the code that fulfills the spec, and advance status to `in_review`.

1. **Load the spec.** Read the spec `README.md` and, for expanded specs, the `tasks/`
   files. Treat the spec as the source of truth — do not re-litigate decisions or
   invent new scope. If you hit an unresolved question or a gap the spec does not
   cover, **stop and ask the user** (or, in `--one-shot`, ask only if genuinely
   blocking); do not guess.
2. **Set status** to `in_progress`.
3. **Work tasks in dependency order.** For expanded specs, follow `depends_on` and the
   implementation code map. Implement one task at a time; keep each task's changes
   scoped to its "Files In Scope" and respect its "Out of Scope". Update each task
   file's `status` (`todo` → `in_progress` → `done`) as you go.
4. **Implement to the requirements, not just "make it run".** Satisfy every `FR-`/`NF-`
   and write the tests named in the Testing Criteria (`TC-` ids). Match existing
   codebase patterns and conventions.
5. **Verify locally.** Run the tests/build. Everything in scope must pass before the
   phase ends.
6. **End the phase.** Set spec `status` to `in_review`, summarize what was built and
   which tasks/requirements are done, and **ask the user to continue**: "Implementation
   complete. Run `/flexspec` again to review the diffs." Stop unless `--one-shot`.

---

# Phase 3 · Review Implementation (Spec Coverage + Slop Scan)

Invoked with `/flexspec` when status is `in_review`. Goal: prove the implementation
fully satisfies the spec **and** that the code is not low-quality "AI slop". On pass,
advance status to `complete`; otherwise, report gaps and loop back.

Review the full diff (e.g. `git diff` against the base branch / spec start) against
the spec.

## 3a. Spec Coverage

- **Every requirement met.** Walk each `FR-`/`NF-` and confirm the diff implements it.
- **Every task done.** For expanded specs, confirm each `T-` task's Acceptance Criteria
  are satisfied and its status is `done`.
- **Every test present and passing.** Each `TC-` exists, runs, and actually asserts the
  behavior it claims (not a stub). Run the suite.
- **No scope drift.** Nothing was built outside the spec; nothing in scope was skipped.
- **No gaps.** List anything missing or partial. If gaps exist, do not mark complete —
  report them and (in `--one-shot`) fix them, then re-review.

## 3b. AI Slop Scan

AI-generated code fails in consistent, taxonomizable ways. Scan the diff against the
five root causes and their signals from the
[Slop Code Taxonomy](https://zbmowrey.com/blog/slop-code-taxonomy/). Flag every hit
with file/line and a fix.

| Root cause | What it produces | Signals to hunt in the diff |
| --- | --- | --- |
| **Context blindness** (no memory of the whole system) | Structural incoherence | Duplicated logic that should reuse existing helpers; inconsistent patterns vs. the rest of the codebase; config/env drift; reinventing utilities that already exist. |
| **Pattern mimicry without judgment** (copies training patterns blindly) | Security + supply-chain + concurrency bugs | SQL injection via string concatenation; missing authz/input validation on trust boundaries; **hallucinated/nonexistent packages** (~20% of AI package suggestions are fake — verify every new dependency resolves); insecure defaults; race conditions / unsynchronized shared state. |
| **Volume outpaces capacity** (more code than can be reviewed) | Comprehension debt | Large unexplained changes; code the author can't justify; churn-prone "written to be rewritten" code; over-broad diffs. |
| **Optimized for "does it run"** (targets happy path) | Bloat, weak runtime quality, fake tests | Dead/unreachable code and needless abstraction; missing error handling and edge cases; **tests that pass but assert nothing** ("green but nothing works"); ignored performance/accessibility. |
| **Institutional exposure** (org-level risk) | License/compliance gaps | **Leaked secrets/keys** in code or config; copied code with incompatible licenses; missing audit/logging where required. |

Quick red-flag checklist (any hit = not slop-free):

- [ ] Every new dependency is real and actually used (no hallucinated packages).
- [ ] No secrets, tokens, or keys committed.
- [ ] No SQL/command/template injection; trust boundaries validate + authorize input.
- [ ] No duplicated logic that should reuse existing code; patterns match the codebase and charter §7 where applicable.
- [ ] Errors and edge cases handled; no happy-path-only code.
- [ ] Tests assert real behavior (not stubs/tautologies) and would fail if the code broke.
- [ ] No dead code, needless abstraction, or unexplained bloat.
- [ ] Concurrency/shared state is safe.
- [ ] Diff is comprehensible — a reviewer could explain every change.

## End of Phase 3

- If coverage is complete and the slop scan is clean: set `status` to `complete`,
  record `implementation_finished`, and summarize.
- If issues remain: keep `status` at `in_review`, present the gaps/slop findings
  with fixes. In `--one-shot`, fix and re-review automatically; otherwise ask the
  user how to proceed.
