---
name: flexspec
description: >
  Run full FlexSpec lifecycle via /flexspec. Use for create/refine/implement/review
  of FlexSpec specs. Workflow is status-driven (author -> implement -> review),
  one phase per prompt unless one-shot. Always scaffold with `flexspec` CLI (`init`,
  `new --template`, `list`, `validate`) and never hand-create spec directories/files.
---

# FlexSpec Lifecycle (`/flexspec`)

A FlexSpec spec is a feature contract before code. State lives in spec `status`.
The `type` field (`feature | bug | chore | refactor | docs | infra | spike |
research`) drives which sections are required and how deep discovery runs.

Templates:
- `templates/flexspec-simple.md`
- `templates/expanded/flexspec-expanded.md`
- `templates/expanded/flexspec-expanded-task.md`

## Core Rules (non-negotiable)

1. Ask, do not assume. Unknowns that can change behavior, scope, tests, security, data, rollout, or UX block planning/implementation.
2. Be exhaustively curious before `planned` for feature/bug; scale discovery to the type (see Discovery Gate Scaling). Think through realistic branches, edge cases, abuse cases, and failure modes, then ask grouped questions for unresolved details.
3. Use CLI for scaffolding only. Never hand-create spec dirs/template copies.
4. Keep scope tied to charter (`.flexspec/charter.md`), especially §7 standards and §8 boundaries.
5. Keep specs as tight as the type allows, but never omit details needed by a future LLM implementer. Token budgets are advisory, not readiness gates.

## Phase Routing by `status`

| Phase | Run when status is | End status |
| --- | --- | --- |
| 1 Author | none / `draft` | `planned` (or `proposed` for spike/research) |
| 2 Implement | `planned` / `in_progress` | `in_review` |
| 3 Review | `in_review` | `complete` |

Spike and research types end at `proposed` after Phase 1 — the plan is the
deliverable, so Phase 2/3 do not run unless the user explicitly converts the
spec to an implementable type.

Default: one phase per `/flexspec` invocation, then stop and ask to continue.

## Run Mode Resolution

Resolve mode before execution:
1. If user passed `--one-shot`: run Author -> Implement -> Review continuously.
2. Else run `flexspec config --json` (or `flexspec config`) and check `always_one_shot`; if `true`, same one-shot behavior.
3. Else: one phase per prompt.

Do not open `.flexspec/config.yaml` manually for `always_one_shot`, `spec_template`, or `specs_dir` - use `flexspec config` / `flexspec config --json`.

One-shot semantics: phases run back-to-back without stopping for handoff, but
**blocking unknowns still pause execution** — ask the user, resolve, then resume
the continuous run. One-shot never skips the "ask, do not assume" rule.

## Template Resolution

Resolve template in this order:
1. `/flexspec --template <simple|expanded>`
2. `flexspec config --json` -> `spec_template` (`simple` or `expanded`)
3. Infer from scope (simple vs expanded)

Only `simple` and `expanded` are valid. Anything else = unset; infer or ask user if borderline.

When scaffolding, always pass explicit template:

```bash
flexspec new <spec-name> --template <simple|expanded>
```

## Type Resolution

Resolve the spec `type` in this order:
1. `/flexspec --type <value>` (user-supplied)
2. Explicit user statement in the request ("this is a chore", "spike to investigate X")
3. Infer from the request shape (see Inference Heuristic below)
4. If borderline, ask the user

Valid `type` values: `feature | bug | chore | refactor | docs | infra | spike | research`.

When scaffolding, pass the type explicitly so frontmatter is correct from the start:

```bash
flexspec new <spec-name> --template <simple|expanded> --type <value>
```

### Inference Heuristic

| Signal | Type |
| --- | --- |
| "fix", "broken", "wrong result", repro steps, regression | `bug` |
| "spike", "investigate", "explore", "evaluate options", "research", "feasibility" | `spike` or `research` |
| "rename", "bump version", "update dependency", "config change", "rotate key" | `chore` |
| "refactor", "restructure", "extract", "simplify" (no behavior change) | `refactor` |
| "document", "README", "guide", "docs" | `docs` |
| "deploy", "migration", "infrastructure", "pipeline", "provision" | `infra` |
| Otherwise | `feature` |

### Per-Type Section Matrix

The single template stays; the agent **omits** inapplicable sections rather than
leaving placeholder text. `NA` = write `Not applicable - this is not a bug fix.` (4/5) or omit.

| Type | §1 | §2 | §3 | §4/5 | §6 | §7 | §8 | §9 | §10 | Discovery | Terminal |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| feature | req | req | req | NA | req | req | req | req | req | Exhaustive | complete |
| bug | req | req | req | req | req | req | req | req | req | Exhaustive (bug-focused) | complete |
| chore | req | req | req | NA | opt | req | light | light | req | Lite (2-3 Q) | complete |
| refactor | req | req | req | NA | req | req | req | req | req | Lite-Med (scope/safety) | complete |
| docs | req | req | req | NA | if flow | req | manual | if appl | req | Lite | complete |
| infra | req | req | req | NA | req | req | req | req | req | Med (ops/rollback) | complete |
| spike | req | req | req | NA | opt | investigative | NA | open Qs | NA | Lite | proposed |
| research | req | req | req | NA | NA | findings | NA | open Qs | NA | Lite | proposed |

`req` = required; `opt` = include if the work has flow; `light` = a few
observable checks, not full TC matrix; `manual` = manual verification steps;
`investigative` = steps describe what to probe, not what to build; `open Qs` =
Section 9 lists open questions the spike/research should answer; `NA` = omit.

## CLI Contract

Run from project root.

| Command | Use |
| --- | --- |
| `flexspec init` | `.flexspec/` or config missing |
| `flexspec config` | read project config (table); prefer `--json` for agents |
| `flexspec new <name> --template <simple\|expanded> [--type <value>]` | create new spec |
| `flexspec list` | discover specs/status/type |
| `flexspec status set <spec> --status <status>` | update spec frontmatter status |
| `flexspec status set <spec> --task <task-file> --status <status>` | update expanded task frontmatter status |
| `flexspec validate` | structural checks after edits / before handoff |
| `flexspec glossary list` | list all glossary terms from `.flexspec/glossary.yaml` |
| `flexspec glossary query <text>` | search glossary terms (exact/alias/substring) |
| `flexspec glossary add <term>` | add or update a glossary term definition |
| `flexspec glossary scan` | scan specs/charter/code for candidate terms (cross-platform, `--json`) |
| `flexspec update` | upgrade CLI, reinstall skills, run migrations (`--dry-run`, `--check`, step flags) |

Forbidden scaffolding actions:
- no manual `specs/NNN-slug` directory creation
- no manual seed `README.md`
- no template copy-paste from `.flexspec/templates`
- no manual sequence numbering
- no manual frontmatter status edits; use `flexspec status set` for spec/task `status`
- no manual frontmatter `type` edits on existing specs; `flexspec update --migrate` backfills missing `type` values

Allowed after `flexspec new`:
- edit CLI-created `README.md` (set `type`, fill sections per matrix, omit NA sections)
- add expanded task files under CLI-created `tasks/`
- keep `task_count` in spec YAML frontmatter and `· **Tasks**: N` in the README metadata line in sync when adding/removing Section 10 tasks or task files (`flexspec validate` warns on drift; `flexspec update --migrate` backfills)

## Charter Gate (Phase 1)

Always read `.flexspec/charter.md` first.

Classify:
- Empty: blank/whitespace only
- Template-only: has `{...}` placeholders or `<!-- ... -->` hints or `status: draft`
- Active: `status: active` and no placeholders/comments

Gate strictness by type:
- **feature / bug**: if missing/empty/template-only, stop and recommend `/flexspec-charter` first. Only continue without full charter if user explicitly insists.
- **chore / refactor / docs / infra / spike / research**: a stub/template-only charter is acceptable. Record the assumption in Section 2 ("Charter not yet active; proceeding with type `<type>`") and continue. Recommend `/flexspec-charter` when the user has time, but do not block.

Charter freshness check before `planned`/`proposed`:
- detect deltas implied by the spec, especially product capabilities, standards, boundaries, glossary terms, interfaces, security, rollout, and operational expectations
- if no deltas: continue silently
- if deltas: update `.flexspec/charter.md` directly and record the change in spec Section 2 under `Charter updates applied automatically`
- only charter conflicts with §7 or §8 are blocking (must ask user before proceeding)
- spike/research types: skip charter freshness (the deliverable is findings, not a product change)

Automatic charter update rules:
- `/flexspec` updates `.flexspec/charter.md` automatically when a spec changes product capabilities, standards, boundaries, or glossary terms.
- Do not ask the user whether to update the charter for in-scope deltas.
- Record the delta in spec Section 2 under `Charter updates applied automatically`.

## Glossary Gate (Phase 1 only)

Always read `.flexspec/glossary.yaml` after the charter.

- List known terms with `flexspec glossary list --json`.
- During authoring, watch for project-specific terms not in the glossary.
- If a term's meaning is clear from context or standard usage, record it silently with `flexspec glossary add <term> --definition <text> --source <source>`.
- If a term is project-specific but unclear, ask the user for the exact meaning before persisting.
- Never fabricate definitions for unclear terms.
- Do not manually edit `.flexspec/glossary.yaml`; always use `flexspec glossary add`.

Phase 2 (implementation) and Phase 3 (review): only **record** clear terms
silently via `flexspec glossary add`. Never interview during implementation or
review — discovery belongs to Phase 1.

## Template Choice Heuristic

Use `simple` for localized/small work (few files, low architectural impact).
Use `expanded` for cross-cutting/large work (multiple subsystems, many tasks, new architecture/data model/interfaces).
If borderline, ask user.

---

# Phase 1: Author

Goal: complete, unambiguous, testable spec on disk (per the type's section
matrix); move status to `planned`, or `proposed` for spike/research types.

## Phase 1 Workflow

1. Read charter, user request, and relevant repo context.
2. Read `.flexspec/glossary.yaml` and note known terms.
3. Resolve `type` (see Type Resolution) and template (see Template Resolution).
4. Initialize if needed (`flexspec init`).
5. Scaffold with CLI (`flexspec new <name> --template <simple|expanded> --type <value>`).
6. Run the Discovery Gate scaled to type before finalizing design details.
7. If the request includes UI work, run the UI Interview Gate too.
8. Fill CLI-created spec files per the per-type section matrix (omit NA sections; do not re-scaffold).
9. Surface unknowns; ask user in grouped questions; resolve all blocking items.
10. Map answers into the sections the type requires.
11. Run readiness checks (sections per matrix, IDs, tests, mappings).
12. Run charter freshness check: update charter automatically for in-scope deltas; only §7/§8 conflicts are blocking. Spike/research skip this.
13. Run glossary gate: record clear terms, ask for unclear ones.
14. Set `status` with `flexspec status set <spec> --status <planned|proposed>` (specs are authored in `draft`).
15. End phase; summarize and ask user to run `/flexspec` again (unless one-shot). For spike/research, the spec is complete at `proposed` — summarize findings and open questions.

## Discovery Gate Scaling (Phase 1)

Run the Discovery Gate for every new or refined spec before `planned`/`proposed`,
scaled to the type.

### Exhaustive (feature, bug)

Process:
- First, mine the charter, glossary, existing code, tests, docs, and user request for answers.
- Then think creatively and adversarially: enumerate happy paths, alternate paths, edge cases, error states, abuse cases, security boundaries, data anomalies, concurrency races, operational failures, migration/rollback concerns, accessibility/performance needs, and test or rollout gaps.
- Ask only questions whose answers can change scope, behavior, implementation order, requirements, tests, security, data handling, UX, or rollout.
- Prefer grouped multiple-choice or short-answer questions. Use multi-select when several options may apply.
- Do not proceed to `planned` with unresolved design decisions. If a point is clear from repo evidence or a documented standard, record it as a confirmed assumption instead of asking.

Question areas to consider:

| Area | Examples to resolve |
| --- | --- |
| Goal and scope | exact outcome, must-have vs excluded work, compatibility promises, success metric |
| User/workflow | actors, entry points, triggers, preconditions, primary path, alternate paths |
| Bug details | expected result, actual result, reproduction, affected versions/data, regression window |
| Data/state | entities, validation, missing/duplicate/stale data, migrations, backfill, retention |
| Interfaces | routes, CLI flags, events, APIs, payloads, errors, versioning, third-party contracts |
| Security/privacy | authn/authz, trust boundaries, injection, secrets, PII, audit, rate limiting |
| Failure handling | dependency outage, timeout, partial success, retry, idempotency, rollback |
| Concurrency | races, shared state, transactions, locking, eventual consistency |
| UX/accessibility | loading/empty/error/success states, keyboard/focus, contrast, reduced motion |
| Performance/scale | input size, latency, throughput, memory, caching, pagination |
| Operations | logging, metrics, feature flags, deployment, migration, rollback, supportability |
| Testing | unit/integration/e2e/manual checks, fixtures, negative tests, acceptance criteria |

High bar: after the gate, a future LLM should not need to ask product/design/security questions to implement the spec correctly.

### Medium (infra, refactor)

Focus on: scope boundaries, rollback/safety, what must not change, migration
order, verification that behavior is preserved. Skip UX/accessibility unless
the infra change surfaces to users. Ask 4-8 grouped questions.

### Lite (chore, docs, spike, research)

Focus on: exact target (what file/version/name), success criterion, any
side effects or compatibility concerns. Ask 2-3 grouped questions. For
spike/research: what question are we answering, what would count as a
satisfying answer, time/depth budget, and what artifacts (notes, prototype,
report) the user expects. Do not block on non-essential unknowns — record
them as open questions in Section 9 for spike/research.

## UI Interview Gate (Phase 1)

Run this gate when the request creates or changes user-facing UI: pages, screens,
visual components, forms, auth, onboarding, settings, dashboards, navigation,
marketing surfaces, or complete UI builds. For tiny copy/style fixes, ask only if
style intent is unclear.

Use the integrated structured question system available in the current agent
runtime when it exists. Prefer grouped multiple-choice questions, with
multi-select when several choices can apply. If no structured question tool
exists, ask the same options in concise text and record the fallback in Section 2.

Ask only the groups needed for the feature, but cover all high-risk unknowns
before `planned`:

| Area | Example options to offer |
| --- | --- |
| Visual identity | reuse existing app style, polished SaaS, playful/illustrated, minimal utility, dense admin |
| Layout system | spacious marketing, balanced app, compact dashboard, mobile-first; grid, spacing scale, breakpoints |
| Component library | existing components to reuse, design-system constraints, third-party library, custom components |
| Component details | icons on primary/secondary buttons, avatars/illustrations, cards vs flat sections, dividers/borders |
| Motion & feedback | transitions, hover/focus/active states, loading states, skeletons, toasts/snackbars, inline validation |
| UX flows & content | empty/zero/error/success states, user flow steps, confirmation patterns, copy tone, error messaging |
| Accessibility & input | labels, focus order, contrast, keyboard shortcuts, reduced motion, screen-reader text, touch targets |
| App fit | routes/screens to match, navigation patterns, existing conventions, dark/light/system mode |

For auth and form-heavy UI, explicitly decide password visibility toggles, icon
usage, field validation timing, submit/loading behavior, forgot/reset links, and
social/provider button treatment when relevant.

Before setting `status: planned`, translate answers into:
- Section 3 use cases and UI states.
- Section 6 workflow branches for user-visible state transitions.
- Section 7 file/component plan naming existing UI patterns to reuse.
- Section 8 tests or manual checks for states, accessibility, and responsive behavior.
- Section 9 FR/NF requirements for visible behavior and accessibility constraints.
- Section 10 tasks that implement chosen states and interactions.
- Section 2 assumptions/risks for deferred style or product choices.

## Authoring Requirements (apply per the type's section matrix)

Section 1 Summary:
- problem, target outcome, who/what affected
- explicit in-scope and out-of-scope

Section 2 Reasons For Change:
- driver, value, consequences if unchanged
- assumptions, risks, charter updates, glossary updates
- open questions must be `None` before `planned` (feature/bug/refactor/infra/docs/chore); for spike/research, open questions live in Section 9 and may remain open at `proposed`

Section 3 Intended Use Case:
- actors, entry points, preconditions, primary flow
- alternate/error/security/data/UX/operational cases relevant to the change

Sections 4 and 5 Bug Results (bug type only):
- for bugs: expected result and actual result, precise enough to reproduce/verify
- for non-bugs: omit these sections entirely (do not write placeholder "Not applicable")

Section 6 Workflow Graph (required for feature/bug/refactor/infra; optional for chore/docs/spike; omit for research):
- high-level concept flow from entry point through services/libraries/data/external systems to outcomes
- valid `mermaid` graph/sequence plus matching workflow trace table - see Workflow Graph Quality Bar
- include material success, failure, validation, permission, retry, and fallback branches

Section 7 Implementation Plan:
- concrete files/components/interfaces to create, modify, or read
- ordered implementation steps with dependencies, files/symbols, and requirement mapping
- expanded only: Data Model / Persistence and External Interfaces subsections must be complete, or explicitly `None`
- for spike: investigative steps (what to probe, what to build to test)
- for research: findings structure (what to investigate, where to look, how to record)

Section 8 Test Plan (required for feature/bug/refactor/infra; light/manual for chore/docs; omit for spike/research):
- `TC-XXX` test criteria map to requirements (and task where relevant)
- every FR covered by at least one TC (feature/bug only)
- negative/security/error/edge-case tests included when behavior has branches
- if untestable, implementation plan must be reworked

Section 9 Functional and Non-Functional Requirements:
- functional: `FR-XXX`, non-functional: `NF-XXX`
- each specific, observable, and testable (feature/bug/refactor/infra)
- for chore/docs: light observable checks or manual verification steps
- for spike/research: open questions the work should answer (may remain open at `proposed`)

Section 10 Tasks (required for feature/bug/refactor/infra/docs/chore; omit for spike/research):
- simple: in-file task table with task number, name, description, blocks, blocked by, requirements mapping
- expanded: index table plus separate task files in `tasks/`
- every Section 6 step and Section 7 step owned by at least one task
- keep `task_count` frontmatter and `· **Tasks**: N` metadata in sync

## Workflow Graph Quality Bar (Phase 1)

Workflow graphs document **conceptual execution** for human and LLM reviewers. They should show how the change works from trigger to outcome without becoming a file-by-file architecture diagram.

Diagram requirements:
- Prefer `flowchart` for conceptual branches; use `sequenceDiagram` when request/response ordering is clearer.
- Every major step uses a numbered label that matches the trace table.
- Include entry point, validation/permission checks, core service/library calls, data stores, external systems, success outcome, and material failure outcomes.
- Label edges with conditions or verbs (`valid`, `invalid`, `calls provider`, `writes audit`, `returns error`).
- Tie `FR-XXX` / `NF-XXX` to table rows where behavior is satisfied or constrained.
- Expanded specs may use multiple diagram/table pairs when multiple workflows exist (CLI vs UI, API vs worker, migration vs runtime).

Trace table requirements:

| Step | Boundary | What Happens | Input / Condition | Outcome | FR/NF |

Keep simple specs to roughly 10 rows per path; split expanded specs into multiple graph/table pairs when longer.

Reviewer test: Can you answer "what happens at step N, with what input or condition, producing what outcome?" from the table alone?

Anti-patterns (reject and rewrite):
- whole-file dependency diagrams with no outcome flow
- generic boxes like `backend` -> `database` without decisions or side effects
- mermaid without matching trace rows or mismatched step numbers
- only the happy path when failure/security/data branches can change implementation
- workflow steps with no owning Section 7 implementation step, Section 8 test, or Section 10 task

When details are unknown: ask the user if the answer can change behavior, security, data handling, tests, or scope. Use nearest concrete anchor only for non-blocking implementation details and record the assumption in Section 2.

## Expanded Task File Contract

Path pattern: `tasks/T-XXX-<slug>.md`

Each task file must include:
- frontmatter: `id`, `name`, `parent_spec`, `status`, `satisfies`, `depends_on`, `verified_by`, `blocks`
- Objective
- Context
- Files In Scope
- Workflow / Requirement Mapping
- Implementation Steps
- Acceptance Criteria
- Testing
- Out of Scope
- Open Questions
- References

Task constraints:
- self-contained execution context
- no unresolved design decisions in steps
- open questions must be empty before starting task
- keep each task under ~1000 tokens (advisory; split large tasks when detail demands it)

## Format and Budget Guidance

- Replace all `{placeholders}` and remove `<!-- -->` guidance comments.
- Keep valid YAML frontmatter.
- IDs are stable and never reused: `FR-`, `NF-`, `T-`, `TC-`.
- Cross-reference IDs consistently across workflow, implementation plan, tests, requirements, and tasks.
- Spec lives at `<specs_dir>/NNN-feature-slug/README.md`.
- Token budgets are **advisory guidance**, not readiness gates. Use them to avoid bloat, but never split a task or omit needed detail just to fit a number:
  - simple README: ~2500
  - expanded README: ~4000
  - each expanded task: ~1000
  - chore/docs: often well under the simple target — that is fine
  - complex features: may exceed the target — split into task files rather than omit detail

## Definition of Ready (Phase 1 exit)

The exit checklist varies by type. Apply the matching variant.

### DoR: feature / bug
- [ ] Scaffold done via CLI (`flexspec init`/`flexspec new` as needed).
- [ ] Correct template and `type` chosen.
- [ ] Discovery Gate (Exhaustive) completed; all design-changing unknowns resolved.
- [ ] All placeholders/comments removed.
- [ ] Required sections complete (per matrix); Section 6 passes Workflow Graph Quality Bar.
- [ ] FR/NF specific and testable.
- [ ] Section 7 steps map to files/symbols and requirements.
- [ ] Tasks map to requirements, implementation steps, and workflow steps.
- [ ] Every FR mapped to >=1 TC.
- [ ] No blocking open questions remain in Section 2 or task files.
- [ ] Charter read; no unresolved conflict with §7/§8.
- [ ] Charter updated automatically for in-scope deltas; any §7/§8 conflicts resolved with user.
- [ ] Optional/project habit: `flexspec validate` has no errors.
- [ ] Spec `status: planned`.

### DoR: refactor / infra
- [ ] Scaffold + template + `type` correct.
- [ ] Discovery Gate (Medium) completed; scope/safety/rollback unknowns resolved.
- [ ] Required sections complete (per matrix); Section 6 passes Quality Bar.
- [ ] Section 7 steps map to files/symbols.
- [ ] Tasks map to implementation steps.
- [ ] No blocking open questions in Section 2.
- [ ] Charter: no §7/§8 conflict.
- [ ] Spec `status: planned`.

### DoR: chore / docs
- [ ] Scaffold + template + `type` correct.
- [ ] Discovery Gate (Lite) completed; 2-3 unknowns resolved.
- [ ] Required sections complete (per matrix); NA sections omitted.
- [ ] No blocking open questions in Section 2.
- [ ] Spec `status: planned`.

### DoR: spike / research (terminal at `proposed`)
- [ ] Scaffold + template + `type` correct.
- [ ] Discovery Gate (Lite) completed; the question-to-answer and success criterion are clear.
- [ ] Section 7 describes what to probe / where to look / how to record findings.
- [ ] Section 9 lists open questions the work should answer (may remain open).
- [ ] No Section 8/10 required.
- [ ] Spec `status: proposed`.

---

# Phase 2: Implement

Run when status is `planned` or `in_progress`.

1. Read spec `README.md` and expanded `tasks/` files if present.
2. Read `.flexspec/glossary.yaml` and note known terms.
3. Set spec status to `in_progress` with `flexspec status set <spec> --status in_progress`.
4. Implement in dependency order (`depends_on`, `blocks`, Section 7 ordered steps, and Section 10 task table).
5. For expanded specs, update task status `todo -> in_progress -> done` with `flexspec status set <spec> --task <task-file> --status <status>`.
6. When adding or removing implementation tasks, update spec `task_count` frontmatter and the README metadata `**Tasks**` segment to match.
7. Stay within spec scope/files and each task's "Out of Scope".
8. Satisfy all `FR`/`NF`; implement tests required by `TC` mappings.
9. Glossary: only **record** clear terms silently via `flexspec glossary add`. Never interview during implementation — discovery belongs to Phase 1.
10. Follow the **Code Comment Policy** below during all implementation.
11. Run project verification (tests/build) and `flexspec validate` when project uses it.
12. Set spec status `in_review` with `flexspec status set <spec> --status in_review`; summarize completed requirements/tasks; stop and ask to continue (unless one-shot).

If unresolved spec gap appears: ask user; do not guess. If the gap reveals a missing branch/security/data/test case, update the spec before implementing it.

## Code Comment Policy

- **Minimize comments.** Avoid writing code comments unless required by linter, project config, or user convention. Follow existing repo comment style - if sparse, stay sparse; if doc-comments on public APIs, match that.
- **Keep comments short.** One sentence or a few words max. Long explanatory blocks reduce readability. If code needs a paragraph to explain, refactor the code for clarity instead.
- **Never reference FlexSpec artifacts in code.** No spec names, directories, task IDs (`T-001`), requirement IDs (`FR-001`), or any FlexSpec artifact in source code. Specs live outside the codebase and may be archived or absent - embedding references creates confusion and stale pointers.

---

# Phase 3: Review (Coverage + Slop Delegation)

Run when status is `in_review`. Review full diff against spec.

## Coverage Checks

- every `FR`/`NF` implemented (feature/bug/refactor/infra)
- every expanded `T` task done and acceptance criteria met
- every `TC` exists, runs, and asserts real behavior (feature/bug)
- every Section 6 workflow branch is implemented or explicitly out of scope (when Section 6 is required for the type)
- every Section 7 implementation step is completed
- no scope drift (nothing required missing, nothing out-of-scope added)

If gaps exist:
- keep `status: in_review`
- report gaps and fixes
- in one-shot, fix then re-review automatically

## Slop Delegation

Phase 3 does **not** run its own slop checklist. It delegates to the
`/flexspec-slop-cleanup` skill, which is the canonical slop reviewer.

Run the slop skill against the spec's files-in-scope (Section 7 file list):

```
/flexspec-slop-cleanup <files-in-scope>
```

- Any finding = review failure. Address each before exiting Phase 3.
- In one-shot mode, `--fix` is allowed: `/flexspec-slop-cleanup --fix <files-in-scope>`.
- The slop skill's "useless comments" pattern enforces the Phase 2 Code Comment Policy at review time.
- If the slop skill is not installed, fall back to the coverage checks above and warn the user to install `/flexspec-slop-cleanup`.

The slop skill supersedes the user-installed `ai-slop` skill. Users who have
`ai-slop` installed should remove it (`npx skills remove ai-slop`) —
`flexspec-slop-cleanup` is the FlexSpec-native equivalent and integrates with
the lifecycle.

## Phase 3 Exit

- Pass: set `status: complete` with `flexspec status set <spec> --status complete`, record `implementation_finished`, summarize.
- Fail: keep `status: in_review`, list required fixes (coverage gaps + slop findings), ask how to proceed (or auto-fix in one-shot).

Spike/research specs never reach Phase 3 — they terminate at `proposed` after Phase 1.
