# FlexSpec Templates

This folder holds the markdown templates FlexSpec uses to scaffold spec files. A
spec is a structured document agreed on **before** code is written, so humans and
AI coding agents share the same definition of done.

## Templates

| File | Location after `flexspec init` | Use when |
| --- | --- | --- |
| `charter.md` | `.flexspec/charter.md` | Application-wide product context (vision, capabilities, boundaries). Filled via `/flexspec-charter`; not a feature spec. |
| `flexspec-simple.md` | `.flexspec/templates/` | Small, focused changes (bug fixes, copy/styling tweaks, adding one form). |
| `expanded/flexspec-expanded.md` | `.flexspec/templates/expanded/` | Large features spanning multiple layers/subsystems (new endpoint sets + DB + UI, auth systems, adding a test suite). |
| `expanded/flexspec-expanded-task.md` | `.flexspec/templates/expanded/` | One self-contained task file under an expanded spec's `tasks/` directory. |

> `charter.md` is embedded from this folder but written to `.flexspec/charter.md` only — it is **not** copied into `.flexspec/templates/`.

> Authoring guidance lives in agent skills, not here: **charter** →
> `skills/flexspec-charter/SKILL.md`; **specs** → `skills/flexspec/SKILL.md`.
> This README is the static reference for template structure and metadata.

## CLI commands

| Command | Purpose |
| --- | --- |
| `flexspec init` | Scaffold `.flexspec/`, config, charter, and templates |
| `flexspec new <name> --template <simple\|expanded>` | Create `NNN-slug/README.md` (and `tasks/` for expanded) |
| `flexspec list` | List specs, statuses, and task counts |
| `flexspec list --json` | Machine-readable spec list |
| `flexspec validate` | Structural validation of config, charter, templates, and specs |
| `flexspec ui` | Local management UI (board, spec browser, settings) |
| `flexspec status set <spec> --status <s>` | Update spec or task status in frontmatter |

## Where Specs Live

Specs are written to a user-configured specs directory. Each spec gets its own
folder named `NNN-feature-slug` (zero-padded sequence + short slug), and the spec
itself is always the `README.md` inside that folder.

```
<specs_dir>/001-user-auth/          # simple spec
  README.md

<specs_dir>/002-billing-export/     # expanded spec
  README.md
  tasks/
    T-001-create-schema.md
    T-002-add-endpoints.md
    T-003-build-ui.md
```

For expanded specs, every implementation task is a separate file under `tasks/`
so each unit of work stays focused enough for an agent to complete without
context rot.

## Frontmatter Metadata

Every spec starts with YAML frontmatter:

| Field | Values | Meaning |
| --- | --- | --- |
| `name` | string | Human-readable spec title. |
| `description` | string | Short summary shown in listings (e.g. management UI, `flexspec list --json`). |
| `status` | `initial` · `refined` · `planned` · `in_progress` · `in_review` · `complete` | Current lifecycle stage. |
| `created` | datetime | When the spec was created. |
| `implementation_start` | datetime | When implementation began. |
| `implementation_finished` | datetime | When implementation completed. |
| `priority` | `low` · `medium` · `high` · `critical` | Relative importance. |
| `tags` | list | Free-form labels for grouping/search. |

### Status Lifecycle

| Status | Meaning |
| --- | --- |
| `initial` | Draft created; summary/design still forming, open questions remain. |
| `refined` | All open questions resolved; design agreed. |
| `planned` | Implementation plan + task list finalized. |
| `in_progress` | Implementation underway. |
| `in_review` | Implementation complete, under review. |
| `complete` | Merged and verified against testing criteria. |

A spec must not advance past `refined` while open questions remain in Section 5.

## Spec Sections

Both templates share the same top-level sections; the expanded template adds
design depth and moves tasks into separate files.

| # | Section | Simple | Expanded (adds) |
| --- | --- | --- | --- |
| 1 | Summary | Overview, scope, outcome. | Same, for a larger feature. |
| 2 | Design | Architecture + file map, **code execution** map (§2.2: diagram + trace table), FR/NF. | Adds Data Model (`erDiagram`) and External Interfaces. |
| 3 | Implementation Plan | **Build + execution enablement** map (§3.1: diagram + task table) + tasks. | Task list is an index table; each task is its own file in `tasks/`. |
| 4 | Testing Criteria | Tests proving each requirement; everything must be testable. | Also maps each test to the implementing task. |
| 5 | Other | Open questions, assumptions, risks, observations. | Same, plus rollout/migration notes. |

### Expanded Task Files

Each `tasks/T-XXX-<slug>.md` is self-contained so an agent can execute it without
drifting: frontmatter (`id`, `parent_spec`, `status`, `satisfies`, `depends_on`,
`verified_by`), Objective, Context, Files In Scope, Implementation Steps,
Acceptance Criteria, Testing, Out of Scope, Open Questions, References.

## Code Map Conventions

FlexSpec code maps document **code execution** for humans and LLM reviewers. Each map section requires **mermaid + a markdown table** with matching step/task numbers.

| Section | Purpose | Diagram | Table |
| --- | --- | --- | --- |
| **§2.2 Code Map** | Runtime execution (debugger-style) | `sequenceDiagram` + `autonumber` preferred; `alt`/`opt` for branches | Execution trace: Step, Location, Executes, Input, Output, FR/NF |
| **§3.1 Implementation Code Map** | Build order + what runtime steps unlock | Tasks + symbols; solid = build order; dotted = enables §2.2 step(s) | Task execution: Task, Build after, Implements §2.2 steps, Symbols, Execution unlocked |

**Location format:** `` `path/to/file :: symbol` `` (handler, `Class.method`, route, or CLI command).

**Interaction labels:** verb + payload — `calls create(dto)`, `returns 201`, `throws ValidationError`, `reads row`.

**Linkage:** every §2.2 step owned by ≥1 task; every §2.1 file on §3.1 table; FR/NF on trace rows where applicable.

Authoring rules: `skills/flexspec/SKILL.md` → **Code Map Quality Bar**.

## ID Conventions

Stable IDs let sections, tasks, and tests cross-reference each other:

| Prefix | Applies to | Example |
| --- | --- | --- |
| `FR-` | Functional requirement | `FR-001` |
| `NF-` | Non-functional requirement | `NF-001` |
| `T-` | Implementation task | `T-001` |
| `TC-` | Test case | `TC-001` |
