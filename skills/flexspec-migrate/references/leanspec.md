# LeanSpec

Official site: [lean-spec.dev](https://lean-spec.dev). CLI: `leanspec` (also `lean-spec`). Config: `.lean-spec/config.json`.

## Detection signature

**Strong match:**

- File `.lean-spec/config.json` at repo root

**Supporting signals:**

- `specsDir` in config (default `specs/`)
- Spec folders matching `NNN-feature-name/` with `README.md` frontmatter
- Optional `.lean-spec/templates/`

**Note:** Config directory is `.lean-spec/` (hyphen), not `.leanspec/`.

## Layout

```text
.
├── .lean-spec/
│   ├── config.json          # specsDir, structure, templates, frontmatter schema
│   └── templates/           # optional custom templates
└── specs/                   # default; override via config specsDir
    └── 042-my-feature/
        ├── README.md        # required main spec (YAML frontmatter)
        ├── DESIGN.md        # optional sub-spec
        ├── IMPLEMENTATION.md  # or PLAN.md
        ├── TESTING.md       # or TEST.md
        └── API.md           # optional
```

Single-file specs may also exist as flat markdown depending on template; check `config.json` → `structure.defaultFile` (usually `README.md`).

## Migratable unit

Each subdirectory under `<specsDir>/` that contains the main spec file (`README.md` by default).

## Template inference

| Condition | FlexSpec |
| --- | --- |
| Only `README.md` (or single main file) | `simple` |
| Sub-spec files present (`DESIGN.md`, `IMPLEMENTATION.md`, `TESTING.md`, etc.) | `expanded` |

Run `leanspec capabilities -o json` if CLI available to confirm status vocabulary before mapping.

## Field mapping → FlexSpec

| LeanSpec (README.md) | FlexSpec target |
| --- | --- |
| Frontmatter: `name`, `description`, `tags`, `priority`, `created` | FlexSpec frontmatter via edit after `flexspec new` (name/description/priority/tags/created) |
| Problem / Goal / Overview | §1 Summary |
| Solution / Approach | §1 + §2.1 |
| Success Criteria / Acceptance Criteria checklists | §2.3 FR + §4 TC placeholders |
| Non-Goals / Out of Scope | §1 out-of-scope |
| Links to sub-specs | §2.1 file table |
| `DESIGN.md` | §2.1 Architecture; expanded root + optional task |
| `IMPLEMENTATION.md` / `PLAN.md` | §3 tasks source |
| `TESTING.md` / `TEST.md` | §4 Testing Criteria |
| `API.md` | §2.1 + §2.3 FR for API behavior |

Use `leanspec view <id>` output if CLI available to enrich mapping.

## Status map

LeanSpec frontmatter `status` (common markdown adapter values):

| LeanSpec status | FlexSpec |
| --- | --- |
| `draft`, `backlog`, `todo` | `draft` |
| `ready`, `refined`, `approved` | `planned` |
| `in-progress`, `in_progress`, `active` | `in_progress` |
| `review`, `in-review` | `in_review` |
| `done`, `complete`, `completed`, `archived` | `complete` |

If unknown: run `leanspec capabilities -o json` and map semantic status field enums. Default **`draft`**.

**Migration default:** set **`draft`** unless user asks to preserve mapped status (see main SKILL.md).

## Slug naming

From folder: `042-my-feature` → `my-feature` (drop numeric prefix).

## Unmapped / notes

- `AGENTS.md` at repo root (LeanSpec template) → not a spec; optional pointer in charter.
- Relationships (`leanspec link` parent/depends_on) → §5 Other or §3 task depends_on in expanded specs (best-effort).
- Custom frontmatter fields from `config.json` → §5 Other.
