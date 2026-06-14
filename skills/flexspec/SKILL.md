---
name: flexspec
description: >
  Run full FlexSpec lifecycle via /flexspec. Use for create/refine/implement/review
  of FlexSpec specs. Workflow is status-driven (author -> implement -> review),
  one phase per prompt unless one-shot. Always scaffold with `flexspec` CLI (`init`,
  `new --template`, `list`, `validate`) and never hand-create spec directories/files.
---

# FlexSpec Lifecycle (`/flexspec`)

A FlexSpec spec is feature contract before code. State lives in spec `status`.

Templates:
- `templates/flexspec-simple.md`
- `templates/expanded/flexspec-expanded.md`
- `templates/expanded/flexspec-expanded-task.md`

## Core Rules (non-negotiable)

1. Ask, do not assume. Unknowns block planning/implementation.
2. Use CLI for scaffolding only. Never hand-create spec dirs/template copies.
3. Keep scope tied to charter (`.flexspec/charter.md`), especially §7 standards and §8 boundaries.
4. Keep spec/token budgets tight.

## Phase Routing by `status`

| Phase | Run when status is | End status |
| --- | --- | --- |
| 1 Author | none / `draft` | `planned` |
| 2 Implement | `planned` / `in_progress` | `in_review` |
| 3 Review | `in_review` | `complete` |

Default: one phase per `/flexspec` invocation, then stop and ask to continue.

## Run Mode Resolution

Resolve mode before execution:
1. If user passed `--one-shot`: run Author -> Implement -> Review continuously.
2. Else run `flexspec config --json` (or `flexspec config`) and check `always_one_shot`; if `true`, same one-shot behavior.
3. Else: one phase per prompt.

Do not open `.flexspec/config.yaml` manually for `always_one_shot`, `spec_template`, or `specs_dir` — use `flexspec config` / `flexspec config --json`.

One-shot still must ask for blocking unknowns.

## Template Resolution

Resolve template in this order:
1. `/flexspec --template <simple|expanded>`
2. `flexspec config --json` -> `spec_template` (`simple` or `expanded`)
3. Infer from scope (simple vs expanded)

Only `simple` and `expanded` valid. Anything else = unset; infer or ask user if borderline.

When scaffolding, always pass explicit template:

```bash
flexspec new <spec-name> --template <simple|expanded>
```

## CLI Contract

Run from project root.

| Command | Use |
| --- | --- |
| `flexspec init` | `.flexspec/` or config missing |
| `flexspec config` | read project config (table); prefer `--json` for agents |
| `flexspec new <name> --template <simple\|expanded>` | create new spec |
| `flexspec list` | discover specs/status |
| `flexspec status set <spec> --status <status>` | update spec frontmatter status |
| `flexspec status set <spec> --task <task-file> --status <status>` | update expanded task frontmatter status |
| `flexspec validate` | structural checks after edits / before handoff |
| `flexspec glossary list` | list all glossary terms from `.flexspec/glossary.yaml` |
| `flexspec glossary query <text>` | search glossary terms (exact/alias/substring) |
| `flexspec glossary add <term>` | add or update a glossary term definition |
| `flexspec update` | upgrade CLI, reinstall skills, run migrations (`--dry-run`, `--check`, step flags) |

Forbidden scaffolding actions:
- no manual `specs/NNN-slug` directory creation
- no manual seed `README.md`
- no template copy-paste from `.flexspec/templates`
- no manual sequence numbering
- no manual frontmatter status edits; use `flexspec status set` for spec/task `status`

Allowed after `flexspec new`:
- edit CLI-created `README.md`
- add expanded task files under CLI-created `tasks/`
- keep `task_count` in spec YAML frontmatter and `· **Tasks**: N` in the README metadata line in sync when adding/removing §3.2 bullets or task files (`flexspec validate` warns on drift; `flexspec update --migrate` backfills)

## Charter Gate (Phase 1)

Always read `.flexspec/charter.md` first.

Classify:
- Empty: blank/whitespace only
- Template-only: has `{...}` placeholders or `<!-- ... -->` hints or `status: draft`
- Active: `status: active` and no placeholders/comments

If missing/empty/template-only:
- stop and recommend `/flexspec-charter` first
- only continue without full charter if user explicitly insists

Charter freshness check before `planned`:
- detect deltas implied by spec (§2, §3, §4, §5-§6, §7, §8, §9)
- if no deltas: continue silently
- if deltas: update `.flexspec/charter.md` directly and record the change in spec §5 Other
- only charter conflicts with §7 or §8 are blocking (must ask user before proceeding)

Automatic charter update rules:
- `/flexspec` updates `.flexspec/charter.md` automatically when a spec changes product capabilities, standards, boundaries, or glossary terms.
- Do not ask the user whether to update the charter for in-scope deltas.
- Record the delta in spec §5 Other under "Charter updates applied automatically".

## Glossary Gate (Phase 1 and Phase 2)

Always read `.flexspec/glossary.yaml` after the charter.

- List known terms with `flexspec glossary list --json`.
- During authoring and implementation, watch for project-specific terms that are not in the glossary.
- If a term's meaning is clear from context or standard usage, record it silently with `flexspec glossary add <term> --definition <text> --source <source>`.
- If a term is project-specific but unclear, ask the user for the exact meaning before persisting.
- Never fabricate definitions for unclear terms.
- Do not manually edit `.flexspec/glossary.yaml`; always use `flexspec glossary add`.

## Template Choice Heuristic

Use `simple` for localized/small work (few files, low architectural impact).
Use `expanded` for cross-cutting/large work (multiple subsystems, many tasks, new architecture/data model/interfaces).
If borderline, ask user.

---

# Phase 1: Author

Goal: complete, unambiguous, testable spec on disk; move status to `planned`.

## Phase 1 Workflow

1. Read charter, user request, and relevant repo context.
2. Read `.flexspec/glossary.yaml` and note known terms.
3. Choose template using resolution/heuristic rules.
4. Initialize if needed (`flexspec init`).
5. Scaffold with CLI (`flexspec new <name> --template <simple|expanded>`).
6. Classify the request for ambiguity. If goals, scope, success criteria, constraints, or key unknowns are missing or vague, run the Ambiguity Interview Gate before filling design details.
7. If the request includes UI work, run the UI Interview Gate before filling design details.
8. Fill CLI-created spec files (do not re-scaffold).
9. Surface unknowns; ask user in grouped questions; resolve all blocking items.
10. Map ambiguity interview answers into clarified goals, scope, success criteria, constraints, requirements, tasks, tests, and §5 assumptions/risks.
11. For UI specs, map UI interview answers into requirements, tasks, testing criteria, and §5 assumptions/risks.
12. Run readiness checks (sections, IDs, tests, mappings, token budgets).
13. Run charter freshness check: update charter automatically for in-scope deltas; only §7/§8 conflicts are blocking.
14. Run glossary gate: record clear terms, ask for unclear ones.
15. Set `status` to `planned` with `flexspec status set <spec> --status planned` (specs are authored in `draft`).
16. End phase; summarize and ask user to run `/flexspec` again (unless one-shot).

## Ambiguity Interview Gate (Phase 1)

Run this gate when the request is vague or underspecified, regardless of whether
it includes UI. Ambiguous signals include missing goals, missing scope
boundaries, missing success criteria, missing constraints, unknown user
workflows, undefined data models, or unspecified integration points. For tiny
copy/style fixes or other obviously trivial edits, you may skip both gates and
record the skip rationale in §5 Other.

Use the integrated structured question system available in the current agent
runtime when it exists (for example Cursor `AskQuestion`, or equivalent Claude,
Codex, or other agent multiple-choice tools). Prefer grouped multiple-choice
questions, with multi-select when several choices can apply. If no structured
question tool exists, ask the same options in concise text and record the fallback
in §5 Other.

Ask only the groups needed for the feature, but cover all high-risk unknowns
before `planned`:

| Area | Example options to offer |
| --- | --- |
| Goal clarity | specific user outcome, business metric, engineering outcome, exploration/spike |
| Scope boundaries | in-scope features, out-of-scope features, MVP vs future, files not to touch |
| Success criteria | observable behavior, acceptance threshold, performance target, test expectation |
| Constraints | existing tech, deadlines, dependencies, backwards compatibility, offline/air-gapped |
| Unknowns | data model, external API, auth/permissions, deployment, analytics, migration |
| Risk/rollback | blast radius, feature flag, rollback plan, user communication |

Before setting `status: planned`, translate answers into:

- Clarified in-scope and out-of-scope statements in §1.
- Concrete FR/NF requirements.
- §2.1 file/component plan.
- §3 tasks with dependency order.
- §4 tests mapped to requirements.
- §5 assumptions/risks for any unresolved or deferred items.

## UI Interview Gate (Phase 1)

Run this gate when the request creates or changes user-facing UI: pages, screens,
visual components, forms, auth, onboarding, settings, dashboards, navigation,
marketing surfaces, or complete UI builds. For tiny copy/style fixes, ask only if
style intent is unclear.

Use the integrated structured question system available in the current agent
runtime when it exists (for example Cursor `AskQuestion`, or equivalent Claude,
Codex, or other agent multiple-choice tools). Prefer grouped multiple-choice
questions, with multi-select when several choices can apply. If no structured
question tool exists, ask the same options in concise text and record the fallback
in §5 Other.

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

- FR/NF requirements for visible behavior and accessibility constraints.
- §2.1 file/component plan naming existing UI patterns to reuse.
- §3 tasks that implement the chosen states and interactions.
- §4 tests or manual checks for states, accessibility, and responsive behavior.
- §5 assumptions/risks for any deferred style or product choices.

## Authoring Requirements (both templates)

Section 1 Summary:
- problem, target outcome, who/what affected
- explicit in-scope and out-of-scope

Section 2 Design:
- architecture plan with concrete files/components
- markdown file table (`File / Type / Role`) covering all touched/referenced files
- valid `mermaid` **code execution** map (§2.2) + execution trace table — see Code Map Quality Bar
- requirements:
  - functional: `FR-XXX`
  - non-functional: `NF-XXX`
  - each specific and testable

Expanded-only Section 2:
- Data Model (with `erDiagram` if persistent data touched; explicit "none" otherwise)
- External Interfaces (APIs/routes/events/CLI/integrations)

Section 3 Implementation Plan:
- **§3.2 Task List (required):** build order, files touched, requirement mapping; when §3.1 is omitted, each task cites `depends_on`, §2.1 files, and §2.2 step range(s)
- **§3.1 Implementation Code Map (optional):** visual build-order map — include only when the **§3.1 complexity heuristic** applies; when included, mermaid diagram + task execution table per quality bar below
- tasks:
  - simple: in-file `T-XXX` list with satisfies mapping
  - expanded: index table + separate task files in `tasks/`

## Code Map Quality Bar (Phase 1)

Code maps document **code execution** for human and LLM reviewers — step through the plan like a debugger, without opening the repo. Before `planned`, **§2.2** must include a **mermaid diagram + execution trace table** that match. **§3.1** is optional by default; when the complexity heuristic requires it (or the author includes it), §3.1 must include a **mermaid diagram + task execution table** that match.

### §2.2 Design Code Map (runtime execution)

**Diagram (required)**
- Prefer `sequenceDiagram` with `autonumber` for call order; use `flowchart` when loops/async are clearer.
- Every step: `path/to/file :: symbol` (or route/CLI/event if symbol TBD).
- Label interactions with verb + payload (`calls create(dto)`, `returns 201`, `throws ErrX`, `reads []byte`).
- Model branches with `alt`/`opt`/`else` (sequence) or labeled branch edges (flowchart) for errors and key conditionals.
- Tie `FR-XXX` / `NF-XXX` to steps where behavior is satisfied or constrained.
- **Expanded**: multiple diagram+table pairs when multiple execution paths exist (CLI vs worker, etc.).

**Execution trace table (required)** — same step numbers as diagram:

| Step | Location | Executes | Input / condition | Output / side effect | FR/NF |

Keep ≤12 rows per path in simple specs; split or add a second diagram/table if longer.

**Reviewer test**: Can you answer "what runs at step N, with what input, producing what output?" from the table alone?

### §3.1 complexity heuristic

**Include §3.1** when any of the following apply:

- Auth, security, or permission-critical flows with multi-step enforcement
- Large refactor or migration spanning many files with strict build ordering
- Cross-cutting work across multiple subsystems (e.g. ≥3 top-level packages or bounded contexts)
- Parallel implementation tracks (branching `depends_on`, or tasks that unblock different §2.2 paths)
- Multiple §2.2 execution paths (CLI vs worker, API vs batch, etc.) where task→step mapping is non-obvious
- Expanded spec with many interdependent tasks (typically ≥5 tasks with non-linear dependencies)
- User explicitly requests a visual implementation map

**Omit §3.1 (default)** when the §3.2 task list alone conveys build order — e.g. simple template, ≤4 tasks in a single linear chain, short §2.2 trace (≤6 steps) with obvious layer-by-layer mapping. Record in §5 Other: `§3.1 omitted: <reason>`.

**When §3.1 is omitted**, §3.2 must still satisfy linkage:

- Every §2.2 step owned by ≥1 task (cite step numbers in task text)
- Every §2.1 file listed on a task row or bullet
- Each task states files touched, `depends_on` (if any), and §2.2 steps implemented

### §3.1 Implementation Code Map (when required or included)

**Diagram (required when §3.1 present)**
- Task nodes: `T-XXX :: file :: symbol` (primary symbols changed).
- Solid edges: build / `depends_on` order from §3.2.
- Dotted edges: `enables §2.2 step N` (or range) — what becomes runnable when the task lands.
- Parallel branches OK; merge before integration tasks.

**Task execution table (required when §3.1 present)**:

| Task | Build after | Implements §2.2 steps | Symbols added/changed | Execution unlocked |

- Every §2.2 step owned by ≥1 task (via §3.1 table and/or §3.2); every §2.1 file on a task row.
- Symbols in tasks must match §2.2 trace `Location` column.

**Reviewer test**: After T-00X, which execution steps from §2.2 can run in a dev environment?

### Anti-patterns (reject and rewrite)

- Architecture-only boxes; no numbered/call-ordered execution.
- Mermaid without a matching trace table (or mismatched step numbers).
- Generic nodes; edges with no verb/payload; task graph with only `T-001 → T-002`.
- §2.2 steps with no owning task in §3.2 or §3.1; §2.1 files absent from task entries.
- §3.1 included without meeting the complexity heuristic (token waste); or omitted on complex work when the task list cannot show build order and §2.2 linkage.
- Scaffold-style §3.1 example diagrams left in a finished spec.

### When symbols are unknown

Use nearest concrete anchor and record in §5 Other. Trace table still required with best-known `Location` and `Executes` columns.

Section 4 Testing Criteria:
- `TC-XXX` test criteria map to requirements (and task where relevant)
- every FR covered by at least one TC
- if untestable, implementation plan must be reworked

Section 5 Other:
- open questions/assumptions/risks/rollout notes
- all blocking open questions must be resolved before `planned`

## Expanded Task File Contract

Path pattern: `tasks/T-XXX-<slug>.md`

Each task file must include:
- frontmatter: `id`, `name`, `parent_spec`, `status`, `satisfies`, `depends_on`, `verified_by`
- Objective
- Context
- Files In Scope
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
- keep each task under ~1000 tokens (split large tasks)

## Format and Budget Rules

- Replace all `{placeholders}` and remove `<!-- -->` guidance comments.
- Keep valid YAML frontmatter.
- IDs are stable and never reused: `FR-`, `NF-`, `T-`, `TC-`.
- Cross-reference IDs consistently across sections/tasks/tests.
- Spec lives at `<specs_dir>/NNN-feature-slug/README.md`.
- Token targets:
  - simple README: <~2000
  - expanded README: <~3500
  - each expanded task: <~1000

## Definition of Ready (Phase 1 exit)

- [ ] Scaffold done via CLI (`flexspec init`/`flexspec new` as needed).
- [ ] Correct template chosen.
- [ ] All placeholders/comments removed.
- [ ] Required sections complete; §2.2 passes Code Map Quality Bar; §3.1 only when heuristic requires or author includes it (full bar when present); §3.1 omission noted in §5 when skipped.
- [ ] FR/NF specific and testable.
- [ ] Tasks mapped to requirements.
- [ ] Every FR mapped to >=1 TC.
- [ ] No blocking open questions remain.
- [ ] Ambiguity Interview Gate run for vague requests; answers mapped into goals, scope, success criteria, constraints, requirements, tasks, tests, and assumptions.
- [ ] UI Interview Gate run for UI specs; answers mapped into requirements, tasks, tests, and assumptions.
- [ ] Charter read; no unresolved conflict with §7/§8.
- [ ] Charter updated automatically for in-scope deltas; any §7/§8 conflicts resolved with user.
- [ ] Optional/project habit: `flexspec validate` has no errors.
- [ ] Spec `status: planned`.

---

# Phase 2: Implement

Run when status is `planned` or `in_progress`.

1. Read spec `README.md` and expanded `tasks/` files if present.
2. Read `.flexspec/glossary.yaml` and note known terms.
3. Set spec status to `in_progress` with `flexspec status set <spec> --status in_progress`.
4. Implement in dependency order (`depends_on` + plan map).
5. For expanded specs, update task status `todo -> in_progress -> done` with `flexspec status set <spec> --task <task-file> --status <status>`.
6. When adding or removing implementation tasks, update spec `task_count` frontmatter and the README metadata `**Tasks**` segment to match.
7. Stay within spec scope/files and each task's "Out of Scope".
8. Satisfy all `FR`/`NF`; implement tests required by `TC` mappings.
9. Run glossary gate during implementation: record clear terms, ask for unclear ones.
10. Follow the **Code Comment Policy** below during all implementation.
11. Run project verification (tests/build) and `flexspec validate` when project uses it.
12. Set spec status `in_review` with `flexspec status set <spec> --status in_review`; summarize completed requirements/tasks; stop and ask to continue (unless one-shot).

If unresolved spec gap appears: ask user; do not guess.

## Code Comment Policy

- **Minimize comments.** Avoid writing code comments unless required by linter, project config, or user convention. Follow existing repo comment style — if sparse, stay sparse; if doc-comments on public APIs, match that.
- **Keep comments short.** One sentence or a few words max. Long explanatory blocks reduce readability. If code needs a paragraph to explain, refactor the code for clarity instead.
- **Never reference FlexSpec artifacts in code.** No spec names, directories, task IDs (`T-001`), requirement IDs (`FR-001`), or any FlexSpec artifact in source code. Specs live outside the codebase and may be archived or absent — embedding references creates confusion and stale pointers.

---

# Phase 3: Review (Coverage + Slop Scan)

Run when status is `in_review`. Review full diff against spec.

## Coverage Checks

- every `FR`/`NF` implemented
- every expanded `T` task done and acceptance criteria met
- every `TC` exists, runs, and asserts real behavior
- no scope drift (nothing required missing, nothing out-of-scope added)

If gaps exist:
- keep `status: in_review`
- report gaps and fixes
- in one-shot, fix then re-review automatically

## Slop Checks (any hit fails review)

- [ ] New dependencies are real, resolvable, and actually used.
- [ ] No secrets/tokens/keys committed.
- [ ] No injection vectors (SQL/command/template); trust boundaries validate + authorize.
- [ ] No unnecessary duplicate logic; codebase conventions respected.
- [ ] Error paths + edge cases handled (not happy-path only).
- [ ] Tests are meaningful (not stubs/tautologies).
- [ ] No dead code, needless abstraction, or unexplained bloat.
- [ ] Concurrency/shared state is safe.
- [ ] Diff is explainable by human reviewer end-to-end.
- [ ] No excessive/narrating code comments; no FlexSpec spec/task references in source code.

Reference: [Slop Code Taxonomy](https://zbmowrey.com/blog/slop-code-taxonomy/)

## Phase 3 Exit

- Pass: set `status: complete` with `flexspec status set <spec> --status complete`, record `implementation_finished`, summarize.
- Fail: keep `status: in_review`, list required fixes, ask how to proceed (or auto-fix in one-shot).
