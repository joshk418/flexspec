# Tessl (spec-driven-development tile)

Docs: [docs.tessl.io/use/spec-driven-development-with-tessl](https://docs.tessl.io/use/spec-driven-development-with-tessl). Tile: [tessl-labs/spec-driven-development](https://tessl.io/registry/tessl-labs/spec-driven-development).

Tessl is a package manager for agent skills/tiles. The SDD tile adds methodology rules and expects specs as `.spec.md` files, often under `specs/`. Tessl also uses `.tessl/` for installed plugins (not feature specs).

## Detection signature

**Strong match (any):**

- Directory `.tessl/` at repo root (plugins, generated `RULES.md`)
- One or more `*.spec.md` files, commonly in `specs/`

**Supporting signals:**

- Requirements with `[@test]` link syntax
- YAML frontmatter with `targets:` globs

**Weak match:** only `.tessl/` without `.spec.md` — Tessl installed but no specs written yet; inform user.

If `.spec.md` files live outside `specs/`, use user-provided path + generic mapping.

## Layout

```text
.
├── .tessl/
│   ├── plugins/                 # downloaded tiles (spec-driven-development, etc.)
│   └── RULES.md                 # generated agent rules (gitignored in managed mode)
└── specs/
    ├── auth.spec.md             # flat .spec.md files (common)
    └── payments.spec.md
```

Alternate: nested `specs/<area>/<name>.spec.md` — still one migratable unit per `.spec.md` file.

**Do not migrate:** `.tessl/plugins/**` (tooling/skills content).

## Migratable unit

Each `**/*.spec.md` file (exclude `.tessl/plugins/`).

## Template inference

| Condition | FlexSpec |
| --- | --- |
| Single `.spec.md`, no separate task files | `simple` |
| User indicates multi-file spec layout (rare for Tessl) | `expanded` |

Tessl specs are usually single-file → default **`simple`**.

## Field mapping → FlexSpec

| Tessl `.spec.md` | FlexSpec target |
| --- | --- |
| Frontmatter `name`, `description` | FlexSpec frontmatter |
| Frontmatter `targets` | Section 7 file table (glob targets as reference rows) |
| Body headings — requirements, behavior | Section 1 + Section 9 FR-* |
| `[@test] path/to/test.py` inline links | Section 8 TC-* rows referencing test path (port link text; do not invent tests) |
| Error handling / edge case sections | Section 9 FR/NF |
| Non-goals section | Section 1 out-of-scope |

Example Tessl requirement line:

```markdown
- Invalid credentials return 401
  [@test] ../tests/auth/test_invalid_credentials.py
```

→ **FR-00N** + **TC-00N** with description from bullet; Verification column cites test path.

## Status map

Tessl `.spec.md` files often have **no status frontmatter**. Workflow rule: spec approved before implementation.

| Signal | FlexSpec |
| --- | --- |
| Spec exists, no linked implementation | `draft` |
| User confirms spec was approved pre-code | `planned` |
| Implementation in progress (user confirms) | `in_progress` |
| `spec-verification` passed (user confirms) | `complete` (migrate as **`draft`** by default) |

Default: **`draft`**.

## Slug naming

From filename: `auth.spec.md` → `auth`.

## Unmapped / notes

- `.tessl/RULES.md` and tile skills — not feature specs.
- Tessl also documents OpenSpec-style `openspec/` init via some tiles — if both Tessl `.spec.md` and `openspec/` exist, treat as separate tools; user picks.
- `[@test]` links must remain verbatim in Section 8; flag broken links in migration report.
- Variant layouts less documented — fall back to `generic.md` if frontmatter/schema mismatch.
