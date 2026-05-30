---
name: flexspec
description: >
  Drive the full FlexSpec spec-driven development lifecycle, invoked with /flexspec.
  Use when the user runs /flexspec or asks to create, draft, refine, implement, or
  review a FlexSpec spec. The skill is a status-driven, three-phase state machine —
  author the spec, implement it, then review the diffs for spec coverage and AI
  "slop" — advancing one phase per prompt (unless --one-shot). **Always scaffold
  specs with the `flexspec` CLI** (`init`, `new --template`, `list`, `validate`) — never
  hand-create spec directories or copy templates. Covers choosing between the simple
  and expanded templates, where specs are written on disk, each section's meaning
  and required format (mermaid, FR/NF/T/TC IDs), the per-task file format for
  expanded specs, a mandatory discovery interview that resolves every unknown in the
  user's request before drafting design, and a slop-detection review checklist.
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
| **1 · Author** | none yet / `initial` / `refined` | Discovery interview → `flexspec new` → author spec; resolve all unknowns before `planned`. | `planned` |
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

When scaffolding, pass the same template to the CLI:

```bash
flexspec new <spec-name> --template <simple|expanded>
```

The `/flexspec --template` flag and `spec_template` config govern skill behavior;
`flexspec new --template` is what actually writes the scaffold on disk.

---

# Phase 1 · Author the Spec

Goal: produce a complete, unambiguous, testable spec on disk and advance status to
`planned`. This is the authoring half of the skill; everything from here through the
"Definition of Ready" checklist is Phase 1 detail.

## The Most Important Rule: Resolve All Unknowns First

Writing the spec is secondary. **Eliminating ambiguity is the primary job.** A spec
built on guesses produces wrong implementations.

- **Interview before you invent.** Run the Discovery Interview (below) on the user's
  current request *before* filling Design, Implementation Plan, or task files.
- As you draft, surface every unknown, assumption, ambiguity, and decision point.
- **Ask the user direct questions** for anything you cannot determine with certainty
  from the repo or their request. Batch related questions; do not drip one at a time.
- Do **not** invent requirements, file paths, behaviors, or constraints to fill gaps.
  If you are unsure, ask.
- A spec **cannot** be marked `refined` (or advance further) and **implementation
  must not start** while any **blocking** open question remains (Section 5 of the spec,
  or a task's Open Questions).
- Only once the user has answered every blocking question is the spec ready.

When in doubt: ask, don't assume.

## Discovery Interview (required before drafting design)

Before writing Design, Implementation Plan, task files, or concrete FR/NF requirements,
**interview the user** to resolve unknowns in their current request. Same spirit as
`/flexspec-charter`: ask, don't assume.

### When to run

| Timing | Required? |
| --- | --- |
| After reading charter + request + repo context, **before** `flexspec new` | **Yes** — for intent, scope, and template-blocking unknowns |
| After scaffold, **before** filling §2 Design / §3 Implementation / tasks | **Yes** — for technical and behavioral unknowns |
| After first draft, **before** `refined` / `planned` | **Yes** — second pass for draft-induced gaps |

**Do not** run `flexspec new` or write substantive spec sections while **blocking**
questions about what to build remain unanswered.

### Step 1 · Audit the request

Parse the user's message (and any linked issues, screenshots, or prior chat) and
build an explicit **Known vs Unknown** ledger:

- **Known** — stated goals, constraints, examples, explicit out-of-scope, named files/APIs.
- **Unknown** — anything needed to write testable FR/NF or a task list but not stated
  or derivable with certainty from repo + charter.
- **Assumed** — inferences you would make if the user did not answer; **never** promote
  an assumption to a requirement without confirmation.

Flag every gap: vague verbs ("improve", "better UX"), missing boundaries, unstated
users/personas, unspecified error/edge behavior, integration points, migration/rollout,
performance/security expectations, and conflicts with charter §7/§8.

### Step 2 · Classify each unknown

| Class | Meaning | Gate |
| --- | --- | --- |
| **Blocking** | Spec cannot be unambiguous without an answer | Must ask; cannot set `refined` or `planned` |
| **Non-blocking** | Reasonable default exists; user may defer | Record in §5 Other as assumption; proceed only if user confirms or defers explicitly |

When unsure whether something is blocking, treat it as **blocking** and ask.

### Step 3 · Interview the user

- Batch **2–4 related questions** per round (not one-at-a-time drips, not 15-question walls).
- Use **numbered questions** so the user can reply inline (`1. …`, `2. …`).
- For each question, say **why it matters** (which section or requirement it unblocks) when non-obvious.
- **Stop and wait** for answers after each round before drafting sections that depend on them.
- If the user defers a blocking item, record it in §5 Other and **do not** advance past `initial`.
- If the user defers a non-blocking item, record the assumed default in §5 and confirm they accept it.

**Forbidden during interview:** scaffolding a spec whose Summary/Design encodes unconfirmed
guesses; filling `{placeholders}` with invented product behavior; choosing expanded vs simple
based on assumed scope the user never confirmed.

### Step 4 · Question bank (feature spec)

Use the categories below. Skip categories already answered in the request or charter;
**do not skip** a category just because the repo "usually" does X — confirm when the
request is silent.

| Category | Sample questions |
| --- | --- |
| **Problem and outcome** | What problem does this solve? What does done look like for the user or system? |
| **Scope** | What is explicitly **in** scope for this spec? What must **not** be included (this iteration)? |
| **Users and context** | Who triggers this? Any roles, permissions, or environments (dev/staging/prod)? |
| **Behavior and UX** | Happy path step-by-step? Error/empty/loading states? Copy or interaction details that matter? |
| **Data and persistence** | New/changed entities or fields? Migration or backfill? Source of truth? |
| **Interfaces** | API shapes, CLI flags, UI routes, events, webhooks — contract expectations? |
| **Integration** | External services, feature flags, auth, existing modules to reuse vs avoid? |
| **Non-functionals** | Performance, security, accessibility, observability, compatibility targets? |
| **Rollout and ops** | Feature flag, migration path, rollback, config changes, docs? |
| **Template and shape** | Simple vs expanded — confirm if borderline or if task count is unclear. |
| **Acceptance** | How will we verify success — manual checks, automated tests, metrics? |
| **Charter alignment** | Anything that might conflict with charter §7 standards or §8 boundaries? |

Add **request-specific** questions for anything ambiguous in the user's exact wording
(e.g. "support dark mode" → all surfaces or one page? system preference only or toggle?).

### Step 5 · Record and iterate

- During drafting, put unresolved items in §5 Other (spec) or task **Open Questions** — then **stop** and run another interview round.
- Before `refined`: §5 and every task **Open Questions** must have no blocking items;
  decisions belong in Summary, Design, or Acceptance Criteria — not left as questions.
- **Second pass:** after the first full draft, re-read the spec as an implementer cold;
  any hesitation ("maybe we should…") becomes a new unknown → interview again.

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

## CLI Commands (mandatory)

**The `flexspec` CLI owns all spec scaffolding.** Run these commands in the shell from
the project root. Do **not** use file-write tools, `mkdir`, or manual template copies
to create spec directories or initial `README.md` files — that bypasses deterministic
numbering and template selection.

| Command | When | Purpose |
| --- | --- | --- |
| `flexspec init` | `.flexspec/` missing | Bootstrap `.flexspec/`, `config.yaml`, `charter.md`, and embedded templates |
| `flexspec new <name> --template <simple\|expanded>` | Starting a new spec | Create `NNN-<slug>/`, seed `README.md` from the chosen template, and (expanded) `tasks/` |
| `flexspec list` | Discovering existing specs | List specs, statuses, and tasks from frontmatter |
| `flexspec list --json` | Scripts / tooling | Same data as the UI API list endpoint |
| `flexspec validate` | Before `list`/`new`, after editing specs, or in CI | Structural checks on config, charter, templates, and specs; exit 1 on errors |
| `flexspec ui` | Optional visibility during implementation | Local dashboard; `--no-open` to skip browser |
| `flexspec status set <spec> --status <s>` | Terminal status updates | Optional `--task T-001-slug.md` for expanded task files |

### `flexspec init`

Run when the project is not initialized:

```bash
flexspec init
```

Optional flags: `--specs-dir <dir>` (default `specs`), `--always-one-shot`, `--force`.

### `flexspec new` (required for every new spec)

After choosing simple vs expanded, scaffold with **both** a name and template:

```bash
flexspec new <spec-name> --template simple
# or
flexspec new <spec-name> --template expanded
```

- **`<spec-name>`** — short feature name (words become a slug, e.g. `user auth` → `user-auth`).
- **`--template` / `-t`** — **always pass explicitly** when authoring via this skill,
  matching the template you chose (simple or expanded). Do not rely on the CLI default
  (`simple`) or on guessing `spec_template` from config unless the user explicitly
  wants config to drive template choice.
- The CLI prints `Created spec NNN-slug`, `path`, and `template` — use that output as
  the canonical spec location.

**Forbidden during scaffolding:** creating `specs/`, `NNN-slug/`, `README.md`, or
`tasks/` yourself; copying template markdown from `.flexspec/templates/` into a new
path; picking sequence numbers manually.

**Allowed after `flexspec new`:** edit the CLI-created `README.md`; for expanded specs,
add task files under the CLI-created `tasks/` directory (see task guide — the CLI
creates the directory only, not individual task files).

Template resolution order in the CLI (for reference): `--template` flag →
`spec_template` in config → default `simple`.

### `flexspec validate`

Run from the project root to catch broken config, missing templates, or unreadable
spec frontmatter before other commands fail opaquely:

```bash
flexspec validate
```

- Prints one tab-separated line per finding (`severity`, `path`, `rule`, `message`), then a summary count.
- Exit **0** when there are no error-severity findings; exit **1** when one or more errors (warnings alone do not fail).
- If `.flexspec/config.yaml` is missing, reports `config.missing` and skips other checks.
- Optional `--strict` is reserved for future semantic checks (structural-only in v1).

Agents should run `flexspec validate` after scaffolding or editing specs and before
marking implementation complete when validation is part of the project's CI habit.

## Where Specs Live

Specs are created on disk, not just in chat.

- **Location** — `specs_dir` in `.flexspec/config.yaml`. Resolve paths from config;
  never hard-code. If config is missing, run `flexspec init` or **ask the user**.
- **One directory per spec** — `flexspec new` creates `NNN-<spec_name>/` (e.g.
  `001-user-auth`) with auto-incrementing sequence numbers.
- **The spec is `README.md`** inside that folder, pre-seeded from the template by the CLI.

Simple spec layout (after `flexspec new … --template simple`):

```
<specs_dir>/001-feature-slug/
  README.md        <- CLI copies flexspec-simple.md; you edit in place
```

Expanded spec layout (after `flexspec new … --template expanded`):

```
<specs_dir>/002-feature-slug/
  README.md        <- CLI copies flexspec-expanded.md; you edit in place
  tasks/           <- CLI creates empty directory
    T-001-<slug>.md  <- you add, using flexspec-expanded-task.md as reference
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
   constraints). Build the Known vs Unknown ledger (Discovery Interview § Step 1).
2. **Discovery interview — round 1 (intent and scope).** Ask blocking questions about
   problem, outcome, scope, template choice (if borderline), and charter conflicts.
   **Wait for answers.** Do not scaffold or draft Design until blocking intent/scope
   questions are resolved (or explicitly deferred with recorded assumptions in §5).
3. **Choose the template** (simple vs expanded) per the rules above; confirm with the
   user if borderline or if they have a preference.
4. **Initialize if needed, then scaffold with the CLI.**
   - If `.flexspec/config.yaml` is missing → run `flexspec init` in the shell.
   - Run `flexspec new <spec-name> --template <simple|expanded>` with the template
     chosen in step 3. **Do not** create directories or `README.md` yourself.
   - Confirm success from CLI output (`Created spec NNN-slug`, `path`, `template`).
   - Optionally run `flexspec list` to verify the new spec appears.
   - Optionally run `flexspec validate` to confirm the project and new spec parse cleanly.
5. **Discovery interview — round 2 (design and behavior).** Before filling §2 Design,
   §3 Implementation Plan, or task files, ask blocking questions from the question bank
   for anything still unknown (data model, interfaces, edge cases, NFRs, rollout). **Wait
   for answers** before writing sections that depend on them.
6. **Edit the CLI-created spec files (do not re-scaffold).** Open the `README.md` the
   CLI wrote and fill frontmatter (`name`, `priority`, `tags`, `status: initial`,
   `created`) plus every section using **confirmed** decisions only. For expanded specs,
   add one task file per task under the CLI-created `tasks/` directory (read
   `.flexspec/templates/expanded/flexspec-expanded-task.md` for structure — do not recreate
   `README.md` or the spec folder). Unresolved items go in §5 / task Open Questions —
   then stop and interview again; do not guess.
7. **Discovery interview — round 3 (draft review).** Re-read the draft as a cold
   implementer; surface any new unknowns or assumptions. Ask the user; wait for answers.
   Move resolved decisions out of §5 into the proper sections.
8. **Self-check** against the Definition of Ready. If anything is only testable with
   rework, rework the implementation plan (and re-interview if that surfaces new unknowns).
9. **Charter freshness check** (before `planned`):
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
10. **Finalize the spec.** When no **blocking** open questions remain, the design is
    agreed, the task list is complete, and the charter freshness check is resolved, move
    `status` through `refined` to `planned`.
11. **End the phase.** Summarize the spec and **ask the user to continue**: "Spec
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
where unresolved items live *during* drafting and between discovery interview rounds.
Before finalizing, every **blocking** open question must be answered (move the decision
into the relevant section). Remaining items should be non-blocking assumptions the user
explicitly accepted or deferred.

## Authoring Task Files (expanded specs)

The CLI creates an empty `tasks/` directory; **you** add task files there. Do not
recreate the spec directory or main `README.md`. Use
`.flexspec/templates/expanded/flexspec-expanded-task.md` as the structure reference.

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

- [ ] Project initialized via `flexspec init` if `.flexspec/` was missing.
- [ ] Spec scaffolded via `flexspec new <name> --template <simple|expanded>` — not
      hand-created directories or template copies.
- [ ] Correct template chosen (simple vs expanded) for the work's size.
- [ ] Spec written to `<specs_dir>/NNN-feature-slug/README.md` (CLI-created, then edited).
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
- [ ] Discovery interview completed: intent/scope resolved before scaffold; design/behavior
      resolved before §2/§3/tasks; draft review pass done.
- [ ] No blocking open questions remain in §5 or task Open Questions; deferred items
      recorded as explicit assumptions.
- [ ] Charter read; spec does not contradict charter §7 (standards) or §8 (boundaries).
- [ ] Charter freshness checked — if deltas exist, update question asked and resolved (or deferred in §5 Other).
- [ ] Optional: `flexspec validate` reports no errors after spec files are finalized.
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
5. **Verify locally.** Run the tests/build and `flexspec validate` when the project is
   initialized. Everything in scope must pass before the phase ends.
5b. **Optional: local dashboard.** Humans may run `flexspec ui --no-open` in a separate
   terminal for a live board and spec browser; the UI reads the same files on disk and
   refreshes when frontmatter changes. Agents still update status by editing markdown
   or via `flexspec status set`; the UI does not replace `/flexspec` lifecycle rules.
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
