# Kiro IDE

Official docs: [kiro.dev/docs/specs](https://kiro.dev/docs/specs). Specs live under `.kiro/specs/`.

## Detection signature

**Strong match:**

- Directory `.kiro/specs/` containing one or more feature subfolders

**Supporting signals:**

- `.kiro/steering/` — agent context (not a migratable spec)
- Files named `requirements.md`, `design.md`, `tasks.md` inside feature folders

**Weak match:** flat `.kiro/specs/requirements.md` at specs root (older/alternate layout) — treat as single spec unit.

## Layout

```text
.kiro/
├── specs/
│   └── <feature-name>/           # one folder per feature
│       ├── requirements.md       # user stories, EARS requirements, AC
│       ├── design.md             # architecture, diagrams, data models
│       └── tasks.md              # executable checklist
│   └── (alternate: bugfix.md instead of requirements.md)
└── steering/
    └── context.md                # agent steering — not a feature spec
```

Some projects use `design/` subfolder instead of single `design.md` — glob for `design/**`.

## Migratable unit

Each `.kiro/specs/<feature-name>/` directory with at least one of `requirements.md`, `bugfix.md`, `design.md`, `tasks.md`.

## Template inference

| Condition | FlexSpec |
| --- | --- |
| Only `requirements.md` (or `bugfix.md`) | `simple` |
| `design.md` and/or `tasks.md` present | `expanded` |

## Field mapping → FlexSpec

| Kiro file | FlexSpec target |
| --- | --- |
| `requirements.md` — user stories, EARS (WHEN…SHALL…), functional reqs | Section 1 + Section 3 + Section 9 FR-* |
| `requirements.md` — acceptance criteria | Section 9 FR + Section 8 TC placeholders |
| `bugfix.md` — current/expected behavior | Section 1 plus Sections 4-5 |
| `design.md` — architecture, sequence diagrams, interfaces | Section 6/7 (preserve mermaid if present) |
| `design.md` — testing/error handling | Section 9 NF + Section 8 seeds |
| `tasks.md` — checkboxes, dependencies | Section 10 task table or expanded task files |

Preserve EARS wording in FR bullets where possible.

## Status map

Kiro tracks task status inside `tasks.md` (in-progress/done markers). Feature-level status is often implicit:

| Signal | FlexSpec |
| --- | --- |
| requirements only | `draft` |
| design approved (user confirms) | `planned` |
| tasks partially complete | `in_progress` |
| all tasks done | `complete` (migrate as **`draft`** by default) |

Default: **`draft`**.

## Slug naming

From folder: `.kiro/specs/oauth-login/` → `oauth-login`.

## Unmapped / notes

- `.kiro/steering/` → suggest charter/agent context update, not a FlexSpec spec.
- Kiro "living spec" sync with codebase — migrated spec is a point-in-time snapshot; note date in Section 2.
- Refine/regenerate workflows — not migrated; user uses `/flexspec` after import.
