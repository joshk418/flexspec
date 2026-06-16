# Spec Kitty

Official docs: [docs.spec-kitty.ai](https://docs.spec-kitty.ai). Repo: [Priivacy-ai/spec-kitty](https://github.com/Priivacy-ai/spec-kitty). Commands: `/spec-kitty.specify`, `.plan`, `.tasks`, etc.

## Detection signature

**Strong match** (both):

- Directory `.kittify/` at repo root (config + templates)
- Directory `kitty-specs/` with `NNN-feature-name/` subfolders

**Weak match:** only `kitty-specs/` — confirm with user (may be copied specs without full Spec Kitty install).

**Supporting signals:**

- `.worktrees/` (execution workspaces)
- Agent slash-command dirs (`.claude/`, `.cursor/`, etc.) with `spec-kitty` commands

If layout does not match below → use `references/generic.md` + ask user.

## Layout

```text
.
├── .kittify/                    # Spec Kitty configuration and templates
├── kitty-specs/
│   └── 001-user-authentication/
│       ├── meta.json            # feature metadata, mission
│       ├── spec.md              # requirements, user stories, acceptance criteria
│       ├── plan.md              # architecture, design decisions
│       ├── research.md          # optional
│       ├── data-model.md        # optional
│       ├── tasks.md             # work package index + checkboxes
│       ├── checklists/
│       │   └── requirements.md
│       └── tasks/               # flat WP prompt files
│           ├── WP01-setup.md
│           └── WP02-api.md
└── .worktrees/                  # git worktrees (do not migrate)
```

## Migratable unit

Each `kitty-specs/NNN-feature-name/` folder.

**Skip:** `.worktrees/`, `.kittify/` (tooling).

## Template inference

| Condition | FlexSpec |
| --- | --- |
| `spec.md` only (no `plan.md`, no `tasks/`) | `simple` |
| `plan.md` and/or `tasks.md` and/or `tasks/WP*.md` | `expanded` |

## Field mapping → FlexSpec

| Spec Kitty file | FlexSpec target |
| --- | --- |
| `meta.json` — title, mission, metadata | frontmatter name/description; Section 1 Summary intro |
| `spec.md` | Section 1 + Section 9 FR/NF |
| `checklists/requirements.md` | Section 9 FR + Section 8 TC seeds |
| `plan.md` | Section 7 implementation plan |
| `research.md`, `data-model.md` | Section 7 references; Section 2 for research notes |
| `tasks.md` — WP index, checkboxes | Section 10 task table |
| `tasks/WP*.md` — work packages | expanded `tasks/T-XXX-<slug>.md` (map WP01 → T-001) |
| WP frontmatter `lane:` status | task status via `flexspec status set --task` |

## Status map

From `meta.json`, WP frontmatter `lane`, or `tasks.md` checkboxes:

| Spec Kitty signal | FlexSpec |
| --- | --- |
| spec only, no plan | `draft` |
| plan exists, tasks not generated | `draft` |
| `lane: planned` on WPs | `planned` (task-level) |
| `lane: in_progress` / active worktree | `in_progress` |
| `lane: done` / all checkboxes complete | `complete` (migrate spec as **`draft`** by default) |

Default spec status: **`draft`**.

## Slug naming

`001-user-authentication` → `user-authentication`.

## Unmapped / notes

- `.worktrees/<feature>-lane-*` — execution sandboxes; do not migrate.
- Contracts/, quickstart.md (if present) — same as Spec Kit mapping.
- Spec Kitty is Spec Kit–adjacent; if both `.specify/` and `kitty-specs/` exist, treat as **two tools** and let user pick which tree to migrate.
