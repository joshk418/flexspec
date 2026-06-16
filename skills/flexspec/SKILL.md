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

Templates:
- `templates/flexspec-simple.md`
- `templates/expanded/flexspec-expanded.md`
- `templates/expanded/flexspec-expanded-task.md`

## Core Rules (non-negotiable)

1. Ask, do not assume. Unknowns that can change behavior, scope, tests, security, data, rollout, or UX block planning/implementation.
2. Be exhaustively curious before `planned`: think through every realistic branch, result, edge case, abuse case, and failure mode, then ask grouped questions for unresolved details.
3. Use CLI for scaffolding only. Never hand-create spec dirs/template copies.
4. Keep scope tied to charter (`.flexspec/charter.md`), especially §7 standards and §8 boundaries.
5. Keep spec/token budgets tight, but never omit details needed by a future LLM implementer.

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

Do not open `.flexspec/config.yaml` manually for `always_one_shot`, `spec_template`, or `specs_dir` - use `flexspec config` / `flexspec config --json`.

One-shot still must ask for blocking unknowns.

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
- keep `task_count` in spec YAML frontmatter and `· **Tasks**: N` in the README metadata line in sync when adding/removing Section 10 tasks or task files (`flexspec validate` warns on drift; `flexspec update --migrate` backfills)

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
- detect deltas implied by the spec, especially product capabilities, standards, boundaries, glossary terms, interfaces, security, rollout, and operational expectations
- if no deltas: continue silently
- if deltas: update `.flexspec/charter.md` directly and record the change in spec Section 2 under `Charter updates applied automatically`
- only charter conflicts with §7 or §8 are blocking (must ask user before proceeding)

Automatic charter update rules:
- `/flexspec` updates `.flexspec/charter.md` automatically when a spec changes product capabilities, standards, boundaries, or glossary terms.
- Do not ask the user whether to update the charter for in-scope deltas.
- Record the delta in spec Section 2 under `Charter updates applied automatically`.

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
6. Run the Exhaustive Discovery Gate before finalizing design details.
7. If the request includes UI work, run the UI Interview Gate too.
8. Fill CLI-created spec files (do not re-scaffold).
9. Surface unknowns; ask user in grouped questions; resolve all blocking items.
10. Map answers into summary, reasons, use cases, workflow graph, implementation plan, tests, requirements, and tasks.
11. Run readiness checks (sections, IDs, tests, mappings, token budgets).
12. Run charter freshness check: update charter automatically for in-scope deltas; only §7/§8 conflicts are blocking.
13. Run glossary gate: record clear terms, ask for unclear ones.
14. Set `status` to `planned` with `flexspec status set <spec> --status planned` (specs are authored in `draft`).
15. End phase; summarize and ask user to run `/flexspec` again (unless one-shot).

## Exhaustive Discovery Gate (Phase 1)

Run this gate for every new or refined spec before `planned`.

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

## Authoring Requirements (both templates)

Section 1 Summary:
- problem, target outcome, who/what affected
- explicit in-scope and out-of-scope

Section 2 Reasons For Change:
- driver, value, consequences if unchanged
- assumptions, risks, charter updates, glossary updates
- open questions must be `None` before `planned`

Section 3 Intended Use Case:
- actors, entry points, preconditions, primary flow
- alternate/error/security/data/UX/operational cases relevant to the change

Sections 4 and 5 Bug Results:
- for bugs: expected result and actual result, precise enough to reproduce/verify
- for non-bugs: explicitly state `Not applicable - this is not a bug fix.`

Section 6 Workflow Graph:
- high-level concept flow from entry point through services/libraries/data/external systems to outcomes
- valid `mermaid` graph/sequence plus matching workflow trace table - see Workflow Graph Quality Bar
- include material success, failure, validation, permission, retry, and fallback branches

Section 7 Implementation Plan:
- concrete files/components/interfaces to create, modify, or read
- ordered implementation steps with dependencies, files/symbols, and requirement mapping
- expanded only: Data Model / Persistence and External Interfaces subsections must be complete, or explicitly `None`

Section 8 Test Plan:
- `TC-XXX` test criteria map to requirements (and task where relevant)
- every FR covered by at least one TC
- negative/security/error/edge-case tests included when behavior has branches
- if untestable, implementation plan must be reworked

Section 9 Functional and Non-Functional Requirements:
- functional: `FR-XXX`
- non-functional: `NF-XXX`
- each specific, observable, and testable

Section 10 Tasks:
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
- keep each task under ~1000 tokens (split large tasks)

## Format and Budget Rules

- Replace all `{placeholders}` and remove `<!-- -->` guidance comments.
- Keep valid YAML frontmatter.
- IDs are stable and never reused: `FR-`, `NF-`, `T-`, `TC-`.
- Cross-reference IDs consistently across workflow, implementation plan, tests, requirements, and tasks.
- Spec lives at `<specs_dir>/NNN-feature-slug/README.md`.
- Token targets:
  - simple README: <~2500
  - expanded README: <~4000
  - each expanded task: <~1000

## Definition of Ready (Phase 1 exit)

- [ ] Scaffold done via CLI (`flexspec init`/`flexspec new` as needed).
- [ ] Correct template chosen.
- [ ] Exhaustive Discovery Gate completed; all design-changing unknowns resolved.
- [ ] All placeholders/comments removed.
- [ ] Required sections complete; Section 6 passes Workflow Graph Quality Bar.
- [ ] FR/NF specific and testable.
- [ ] Section 7 steps map to files/symbols and requirements.
- [ ] Tasks map to requirements, implementation steps, and workflow steps.
- [ ] Every FR mapped to >=1 TC.
- [ ] No blocking open questions remain in Section 2 or task files.
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
4. Implement in dependency order (`depends_on`, `blocks`, Section 7 ordered steps, and Section 10 task table).
5. For expanded specs, update task status `todo -> in_progress -> done` with `flexspec status set <spec> --task <task-file> --status <status>`.
6. When adding or removing implementation tasks, update spec `task_count` frontmatter and the README metadata `**Tasks**` segment to match.
7. Stay within spec scope/files and each task's "Out of Scope".
8. Satisfy all `FR`/`NF`; implement tests required by `TC` mappings.
9. Run glossary gate during implementation: record clear terms, ask for unclear ones.
10. Follow the **Code Comment Policy** below during all implementation.
11. Run project verification (tests/build) and `flexspec validate` when project uses it.
12. Set spec status `in_review` with `flexspec status set <spec> --status in_review`; summarize completed requirements/tasks; stop and ask to continue (unless one-shot).

If unresolved spec gap appears: ask user; do not guess. If the gap reveals a missing branch/security/data/test case, update the spec before implementing it.

## Code Comment Policy

- **Minimize comments.** Avoid writing code comments unless required by linter, project config, or user convention. Follow existing repo comment style - if sparse, stay sparse; if doc-comments on public APIs, match that.
- **Keep comments short.** One sentence or a few words max. Long explanatory blocks reduce readability. If code needs a paragraph to explain, refactor the code for clarity instead.
- **Never reference FlexSpec artifacts in code.** No spec names, directories, task IDs (`T-001`), requirement IDs (`FR-001`), or any FlexSpec artifact in source code. Specs live outside the codebase and may be archived or absent - embedding references creates confusion and stale pointers.

---

# Phase 3: Review (Coverage + Slop Scan)

Run when status is `in_review`. Review full diff against spec.

## Coverage Checks

- every `FR`/`NF` implemented
- every expanded `T` task done and acceptance criteria met
- every `TC` exists, runs, and asserts real behavior
- every Section 6 workflow branch is implemented or explicitly out of scope
- every Section 7 implementation step is completed
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
