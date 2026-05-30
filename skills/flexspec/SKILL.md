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
| 1 Author | none / `initial` / `refined` | `planned` |
| 2 Implement | `planned` / `in_progress` | `in_review` |
| 3 Review | `in_review` | `complete` |

Default: one phase per `/flexspec` invocation, then stop and ask to continue.

## Run Mode Resolution

Resolve mode before execution:
1. If user passed `--one-shot`: run Author -> Implement -> Review continuously.
2. Else if `.flexspec/config.yaml` has `always_one_shot: true`: same one-shot behavior.
3. Else: one phase per prompt.

One-shot still must ask for blocking unknowns.

## Template Resolution

Resolve template in this order:
1. `/flexspec --template <simple|expanded>`
2. `.flexspec/config.yaml` -> `spec_template: simple|expanded`
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
| `flexspec new <name> --template <simple\|expanded>` | create new spec |
| `flexspec list` | discover specs/status |
| `flexspec validate` | structural checks after edits / before handoff |

Forbidden scaffolding actions:
- no manual `specs/NNN-slug` directory creation
- no manual seed `README.md`
- no template copy-paste from `.flexspec/templates`
- no manual sequence numbering

Allowed after `flexspec new`:
- edit CLI-created `README.md`
- add expanded task files under CLI-created `tasks/`

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
- if deltas: ask whether charter needs update (`yes/no/partial`)

Delta gating:
- non-one-shot: do not set `planned` until delta question answered; deferred answer must be recorded in spec §5 Other
- one-shot: only charter conflicts with §7 or §8 are blocking (must ask). Other deltas become "charter follow-up" note in spec §5 Other.

Do not auto-edit charter unattended in one-shot.

## Template Choice Heuristic

Use `simple` for localized/small work (few files, low architectural impact).
Use `expanded` for cross-cutting/large work (multiple subsystems, many tasks, new architecture/data model/interfaces).
If borderline, ask user.

---

# Phase 1: Author

Goal: complete, unambiguous, testable spec on disk; move status to `planned`.

## Phase 1 Workflow

1. Read charter, user request, and relevant repo context.
2. Choose template using resolution/heuristic rules.
3. Initialize if needed (`flexspec init`).
4. Scaffold with CLI (`flexspec new <name> --template <simple|expanded>`).
5. Fill CLI-created spec files (do not re-scaffold).
6. Surface unknowns; ask user in grouped questions; resolve all blocking items.
7. Run readiness checks (sections, IDs, tests, mappings, token budgets).
8. Run charter freshness check and resolve/defer per mode rules.
9. Set `status` through `refined` to `planned`.
10. End phase; summarize and ask user to run `/flexspec` again (unless one-shot).

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
- implementation map (§3.1): build order + which §2.2 steps each task enables, with task execution table
- tasks:
  - simple: in-file `T-XXX` list with satisfies mapping
  - expanded: index table + separate task files in `tasks/`

## Code Map Quality Bar (Phase 1)

Code maps document **code execution** for human and LLM reviewers — step through the plan like a debugger, without opening the repo. Before `planned`, §2.2 and §3.1 must each include a **mermaid diagram + markdown table** that match.

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

### §3.1 Implementation Code Map (build order + execution enablement)

**Diagram (required)**
- Task nodes: `T-XXX :: file :: symbol` (primary symbols changed).
- Solid edges: build / `depends_on` order from §3.2.
- Dotted edges: `enables §2.2 step N` (or range) — what becomes runnable when the task lands.
- Parallel branches OK; merge before integration tasks.

**Task execution table (required)**:

| Task | Build after | Implements §2.2 steps | Symbols added/changed | Execution unlocked |

- Every §2.2 step owned by ≥1 task; every §2.1 file on a task row.
- Symbols in tasks must match §2.2 trace `Location` column.

**Reviewer test**: After T-00X, which execution steps from §2.2 can run in a dev environment?

### Anti-patterns (reject and rewrite)

- Architecture-only boxes; no numbered/call-ordered execution.
- Mermaid without a matching trace table (or mismatched step numbers).
- Generic nodes; edges with no verb/payload; task graph with only `T-001 → T-002`.
- §2.2 steps with no owning task in §3.1; §2.1 files absent from §3.1 table.

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
- [ ] Required sections complete, with valid mermaid blocks passing Code Map Quality Bar.
- [ ] FR/NF specific and testable.
- [ ] Tasks mapped to requirements.
- [ ] Every FR mapped to >=1 TC.
- [ ] No blocking open questions remain.
- [ ] Charter read; no unresolved conflict with §7/§8.
- [ ] Charter delta question resolved or deferred per mode rules.
- [ ] Optional/project habit: `flexspec validate` has no errors.
- [ ] Spec `status: planned`.

---

# Phase 2: Implement

Run when status is `planned` or `in_progress`.

1. Read spec `README.md` and expanded `tasks/` files if present.
2. Set spec status to `in_progress`.
3. Implement in dependency order (`depends_on` + plan map).
4. For expanded specs, update task status `todo -> in_progress -> done`.
5. Stay within spec scope/files and each task's "Out of Scope".
6. Satisfy all `FR`/`NF`; implement tests required by `TC` mappings.
7. Run project verification (tests/build) and `flexspec validate` when project uses it.
8. Set spec status `in_review`; summarize completed requirements/tasks; stop and ask to continue (unless one-shot).

If unresolved spec gap appears: ask user; do not guess.

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

Reference: [Slop Code Taxonomy](https://zbmowrey.com/blog/slop-code-taxonomy/)

## Phase 3 Exit

- Pass: set `status: complete`, record `implementation_finished`, summarize.
- Fail: keep `status: in_review`, list required fixes, ask how to proceed (or auto-fix in one-shot).
