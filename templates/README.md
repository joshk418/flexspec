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
| `flexspec config` | Show project config (KEY / VALUE table) |
| `flexspec config --json` | Machine-readable project config |
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
| `status` | `draft` · `planned` · `in_progress` · `in_review` · `complete` | Current lifecycle stage. |
| `created` | datetime | When the spec was created. |
| `implementation_start` | datetime | When implementation began. |
| `implementation_finished` | datetime | When implementation completed. |
| `priority` | `low` · `medium` · `high` · `critical` | Relative importance. |
| `tags` | list | Free-form labels for grouping/search. |

### Status Lifecycle

| Status | Meaning |
| --- | --- |
| `draft` | Authoring: summary/design forming; open questions may remain (interview happens here). |
| `planned` | Open questions resolved; implementation plan + task list finalized. |
| `in_progress` | Implementation underway. |
| `in_review` | Implementation complete, under review. |
| `complete` | Merged and verified against testing criteria. |

A spec must not advance to `planned` while open questions remain in Section 2.

## Spec Sections

Both templates share the same top-level sections; the expanded template adds
data/interface depth and moves tasks into separate files.

| # | Section | Simple | Expanded (adds) |
| --- | --- | --- | --- |
| 1 | Summary | Problem, outcome, affected users/systems, scope boundaries. | Same, for a larger feature. |
| 2 | Reasons For Change | Driver, value, consequences, assumptions, risks, charter/glossary updates, open questions. | Same, with rollout/migration risks as needed. |
| 3 | Intended Use Case | Actors, entry points, preconditions, primary/alternate/security/data edge cases. | Same, with operational cases. |
| 4 | Expected Result (bugs only) | Expected behavior for bug specs; otherwise explicitly not applicable. | Same. |
| 5 | Actual Result (bugs only) | Observed broken behavior for bug specs; otherwise explicitly not applicable. | Same. |
| 6 | Workflow Graph | High-level concept flow from entry point through services/libraries/data/external systems to outcomes. | Multiple graph/table pairs when needed. |
| 7 | Implementation Plan | Files/interfaces plus ordered implementation steps. | Adds data model, persistence, and external interface detail. |
| 8 | Test Plan | Tests proving each requirement; everything must be testable. | Also maps each test to the implementing task. |
| 9 | Functional and Non-Functional Requirements | Stable FR/NF IDs, specific and testable. | Same. |
| 10 | Tasks | Task number, name, description, blocks, blocked by, requirement mapping. | Task index + per-task files in `tasks/`. |

### Expanded Task Files

Each `tasks/T-XXX-<slug>.md` is self-contained so an agent can execute it without
drifting: frontmatter (`id`, `parent_spec`, `status`, `satisfies`, `depends_on`,
`verified_by`, `blocks`), Objective, Context, Files In Scope, Workflow /
Requirement Mapping, Implementation Steps, Acceptance Criteria, Testing, Out of
Scope, Open Questions, References.

## Workflow Graph Conventions

FlexSpec workflow graphs document the **conceptual flow** for humans and LLM
reviewers. **Section 6** always requires **mermaid + a markdown table** with
matching step numbers. It should show the entry point, decisions, services,
libraries, data stores, external systems, success outcomes, and material failure
outcomes without turning into a whole-file architecture diagram.

| Section | Required? | Purpose | Diagram | Table |
| --- | --- | --- | --- | --- |
| **Section 6 Workflow Graph** | Yes | Conceptual end-to-end flow | `flowchart` preferred; `sequenceDiagram` OK for request/response flows | Step, Boundary, What Happens, Input / Condition, Outcome, FR/NF |

**Boundary format:** route, command, UI action, event, job, service/library,
database, queue, third-party integration, or visible outcome.

**Interaction labels:** verb + payload/result, for example `validates request`,
`calls billing provider`, `writes audit row`, `returns actionable error`.

**Linkage:** every Section 6 step should be covered by Section 7 implementation
steps and Section 10 tasks; FR/NF IDs belong on graph rows where behavior is
satisfied or constrained.

Authoring rules: `skills/flexspec/SKILL.md` -> **Workflow Graph Quality Bar**.

## ID Conventions

Stable IDs let sections, tasks, and tests cross-reference each other:

| Prefix | Applies to | Example |
| --- | --- | --- |
| `FR-` | Functional requirement | `FR-001` |
| `NF-` | Non-functional requirement | `NF-001` |
| `T-` | Implementation task | `T-001` |
| `TC-` | Test case | `TC-001` |
