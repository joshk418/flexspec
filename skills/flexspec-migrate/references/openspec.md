# OpenSpec

Official repo: [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec). Workflow: `/opsx:propose`, `/opsx:apply`, `/opsx:archive`. CLI: `openspec init`.

## Detection signature

**Strong match:**

- Directory `openspec/` at repo root

Common children:

- `openspec/specs/` — canonical capability specs
- `openspec/changes/` — active change workspaces
- `openspec/config.yaml` — optional project config
- `openspec/AGENTS.md` or `openspec/project.md` — optional agent context

## Layout

```text
openspec/
├── config.yaml              # optional: schema, context, rules
├── AGENTS.md                # optional
├── specs/                   # source of truth (current system behavior)
│   └── <domain>/
│       └── spec.md
└── changes/                 # proposed work (one folder per change)
    ├── <change-name>/
    │   ├── proposal.md      # why / what
    │   ├── design.md        # how / architecture
    │   ├── tasks.md         # implementation checklist
    │   └── specs/           # delta specs (ADDED/MODIFIED/REMOVED)
    │       └── <domain>/
    │           └── spec.md
    └── archive/             # completed changes (historical)
        └── YYYY-MM-DD-<name>/
```

## What to migrate

Ask user which mode:

| Mode | Migratable units | Typical use |
| --- | --- | --- |
| **Active changes** (default for in-flight work) | Each `openspec/changes/<change-name>/` (not `archive/`) | Features being built |
| **Canonical specs** | Each `openspec/specs/<domain>/` | Documenting current system behavior |
| **Both** | All of the above | Full conversion |

**Skip by default:** `openspec/changes/archive/` unless user wants historical import.

## Template inference

| Condition | FlexSpec |
| --- | --- |
| Single `spec.md` or delta only, no `tasks.md` | `simple` |
| `proposal.md` + `design.md` + `tasks.md` or delta `specs/` | `expanded` |

## Field mapping → FlexSpec

| OpenSpec artifact | FlexSpec target |
| --- | --- |
| `proposal.md` — intent, scope, approach | Section 1 + Section 2 |
| `proposal.md` — out of scope | Section 1 out-of-scope |
| `design.md` | Section 7 implementation plan |
| `specs/<domain>/spec.md` or change delta `specs/**/spec.md` — Requirements, Scenarios (Given/When/Then) | Section 9 FR-*; scenarios -> TC placeholders in Section 8 or Section 2 |
| Delta sections `## ADDED Requirements`, `MODIFIED`, `REMOVED` | Section 1 note "delta from canonical"; FR in Section 9; REMOVED -> Section 2 |
| `tasks.md` — `- [ ]` checklist items | Section 10 task table or expanded task files |
| `config.yaml` context/rules | Section 2 (charter follow-up); do not copy wholesale into every spec |

## Status map

OpenSpec changes do not use a single standard status file. Infer:

| Signal | FlexSpec status |
| --- | --- |
| Change folder in `changes/` (not archive), tasks mostly unchecked | `draft` or `planned` (default **`draft`**) |
| `tasks.md` partially checked | `in_progress` if user confirms active work |
| Folder in `changes/archive/` | `complete` content but migrate as **`draft`** for FlexSpec lifecycle |
| Canonical `openspec/specs/` only | `draft` (behavior docs, not lifecycle) |

Default: **`draft`**.

## Slug naming

- Change: `add-dark-mode` → `add-dark-mode`
- Domain: `auth-login` → `auth-login`

Avoid collision with existing FlexSpec ids.

## Unmapped / notes

- Delta merge semantics (ADDED/MODIFIED/REMOVED) — preserve verbatim in Section 2 or Section 9.
- `openspec/AGENTS.md` → not a feature spec; optional charter/process note.
- After migration, user may retire `openspec/` directory manually if confirmed.
