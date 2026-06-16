# GitHub Spec Kit

Official repo: [github/spec-kit](https://github.com/github/spec-kit). Slash commands: `/speckit.specify`, `/speckit.plan`, `/speckit.tasks`.

## Detection signature

**Strong match** (both):

- Directory `.specify/` at repo root (contains `memory/`, `scripts/`, `templates/`)
- Directory `specs/` at repo root with subfolders matching `NNN-*`

**Weak match:** only `specs/NNN-*/spec.md` without `.specify/` — confirm with user (may be manual Spec Kit layout or unrelated).

**Not Spec Kit:** specs only under `.specify/` (that's config, not feature specs).

## Layout

```text
.
├── .specify/
│   ├── memory/constitution.md
│   ├── scripts/          # bash or powershell helpers
│   └── templates/        # spec-template.md, plan-template.md, tasks-template.md
└── specs/
    └── 001-feature-name/
        ├── spec.md           # requirements, user stories, success criteria
        ├── plan.md           # architecture, stack, phases (after /speckit.plan)
        ├── tasks.md          # phased task checklist (after /speckit.tasks)
        ├── research.md       # optional
        ├── data-model.md     # optional
        ├── quickstart.md     # optional
        └── contracts/        # optional API/interface defs
```

**Migratable unit:** each `specs/NNN-feature-name/` folder.

**Do not migrate:** `.specify/` itself (tooling/templates). Optionally note `constitution.md` content for charter follow-up, not as a FlexSpec spec.

## Template inference

| Condition | FlexSpec |
| --- | --- |
| Only `spec.md` exists | `simple` |
| `plan.md` and/or `tasks.md` exist | `expanded` |

## Field mapping → FlexSpec

| Spec Kit file / section | FlexSpec target |
| --- | --- |
| `spec.md` — Overview, User Stories, Requirements, Success Criteria, Edge Cases | Section 1 + Section 3 + Section 9 FR/NF |
| `spec.md` — Out of scope / Non-goals | Section 1 out-of-scope |
| `plan.md` — Technical Context, Architecture, Project Structure | Section 7 implementation plan (prose + file table best-effort) |
| `plan.md` — Constitution Check, Constraints | Section 9 NF-* |
| `data-model.md` | Section 7.2 data model placeholder |
| `contracts/*` | Section 7 interfaces; Section 9 FR for API contracts |
| `tasks.md` — phases, task lines, `[P]` markers | Section 10 task table or expanded `tasks/T-XXX-*.md` |
| `research.md`, `quickstart.md` | Section 2 notes (links/summary); do not invent design from empty sections |
| User story priorities P1/P2/P3 | Section 10 task ordering; priority in frontmatter if clear |

## Status map

Spec Kit feature folders typically have **no per-feature status file**. Infer from git branch state or file completeness:

| Signal | FlexSpec status |
| --- | --- |
| Only `spec.md` | `draft` |
| `plan.md` present, no `tasks.md` | `draft` |
| `tasks.md` present, unchecked tasks | `planned` or `draft` (default `draft`) |
| Implementation implied complete (user says so) | `draft` at migration (user runs `/flexspec` to re-validate) |

Default: **`draft`**.

## Slug naming

From folder name: `001-create-taskify` → `create-taskify` (drop numeric prefix for `flexspec new` slug; CLI assigns next sequence).

## Unmapped / notes

- Git feature branches created by Spec Kit are not migrated.
- `constitution.md` → suggest `/flexspec-charter`, not a feature spec.
- Constitution violations / gate sections -> Section 2 notes only.
