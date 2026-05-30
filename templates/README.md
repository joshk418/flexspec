# Flexspec Templates

This folder holds the markdown templates Flexspec uses to scaffold spec files. A
spec is a structured document agreed on **before** code is written, so humans and
AI coding agents share the same definition of done.

## Templates

| File | Mode | Use when |
| --- | --- | --- |
| `flexspec-simple.md` | Simple (single file) | Small, focused changes (bug fixes, copy/styling tweaks, adding one form). |
| `expanded/flexspec-expanded.md` | Expanded (root spec) | Large features spanning multiple layers/subsystems (new endpoint sets + DB + UI, auth systems, adding a test suite). |
| `expanded/flexspec-expanded-task.md` | Expanded (task) | One self-contained task file under an expanded spec's `tasks/` directory. |

> Authoring guidance (how to fill each section, how to choose between simple and
> expanded, and the required question-asking workflow) lives in the **flexspec
> skill** at `skills/flexspec/SKILL.md`, not here. This README is the static
> reference for the template structure and metadata.

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
| 2 | Design | Architecture + file map, mermaid code map, FR/NF requirements. | Adds Data Model (`erDiagram`) and External Interfaces. |
| 3 | Implementation Plan | Code map + ID'd task list in-file. | Task list is an index table; each task is its own file in `tasks/`. |
| 4 | Testing Criteria | Tests proving each requirement; everything must be testable. | Also maps each test to the implementing task. |
| 5 | Other | Open questions, assumptions, risks, observations. | Same, plus rollout/migration notes. |

### Expanded Task Files

Each `tasks/T-XXX-<slug>.md` is self-contained so an agent can execute it without
drifting: frontmatter (`id`, `parent_spec`, `status`, `satisfies`, `depends_on`,
`verified_by`), Objective, Context, Files In Scope, Implementation Steps,
Acceptance Criteria, Testing, Out of Scope, Open Questions, References.

## ID Conventions

Stable IDs let sections, tasks, and tests cross-reference each other:

| Prefix | Applies to | Example |
| --- | --- | --- |
| `FR-` | Functional requirement | `FR-001` |
| `NF-` | Non-functional requirement | `NF-001` |
| `T-` | Implementation task | `T-001` |
| `TC-` | Test case | `TC-001` |
