---
name: flexspec-migrate
description: >
  Convert specs from other SDD tools (Spec Kit, OpenSpec, LeanSpec, Spec Kitty,
  Kiro, Tessl) into FlexSpec format via /flexspec-migrate. Detects source dirs,
  maps content, scaffolds with flexspec CLI. Use when migrating to flexspec.
---

# FlexSpec Migration (`/flexspec-migrate`)

Convert specs from another spec-driven-development tool into FlexSpec. Embedded
tool docs live in `references/` — read those before mapping; do not web-search
for supported tools.

## Triggers

- `/flexspec-migrate` (optional: `[tool]` and `[path]`)
- "migrate my specs to flexspec", "convert from openspec/speckit/leanspec", etc.

## Core Rules

1. **Ask, do not assume.** Confirm scope and destructive actions before writes.
2. **CLI only for scaffolding.** `flexspec new`, `flexspec status set` — never hand-create spec dirs, sequence numbers, or frontmatter status.
3. **Non-destructive default.** Keep source files unless user explicitly confirms deletion for this run.
4. **No web without permission.** If detection fails, ask user to allow web lookup; if denied or still unknown, ask for the source directory path.
5. **Draft end-state.** Port content into `draft` specs; do not fabricate §2.2/§3.1 code maps or TC rows. User runs `/flexspec` next to finish design.
6. **Never invent content.** Unmapped fields → report in migration summary + note in spec §5 Other.

## Prerequisites

Run from project root. If FlexSpec not initialized:

```bash
flexspec init
flexspec config --json   # specs_dir, spec_template
```

## Workflow

### 1. Detect source tool(s)

Scan repo root for signatures (see **Detection summary**). If user passed `[tool]`
or `[path]`, validate against the matching reference doc.

Multiple tools may coexist — inventory each separately.

If no match: ask permission to web-search docs → if denied/unknown, ask user for
spec directory path → treat as **generic** (`references/generic.md`).

If signature matches but layout differs (common for Spec Kitty / Tessl variants):
fall back to generic interview for that path.

### 2. Inventory specs

List migratable units per tool (feature folders, change folders, `.spec.md` files,
LeanSpec spec dirs, etc.). Read `references/<tool>.md` **Layout** section.

Present inventory to user; require multi-select confirmation of which specs to migrate.

### 3. Ask per-run policies

Before any writes, confirm:

1. **Which specs** (from step 2).
2. **Originals:** keep (default) or delete after successful migration? Deletion only after explicit yes.

### 4. Migrate each selected spec

For each source spec:

**a. Infer template**

| Source shape | FlexSpec template |
| --- | --- |
| Single main doc, no task files / task list | `simple` |
| Multiple docs (plan/design/tasks) or `tasks/` dir or task checklist file | `expanded` |

Borderline → ask user.

**b. Scaffold**

```bash
flexspec new <slug> --template <simple|expanded>
```

Slug: kebab-case from source folder/title; avoid collision with `flexspec list`.

**c. Map content** (per `references/<tool>.md` field table)

Fill CLI-created files only:

| FlexSpec section | Source (typical) |
| --- | --- |
| §1 Summary | overview, problem, goals, in/out scope |
| §2.1 Architecture | plan/design file prose + file lists (best-effort) |
| §2.2 Code Map | **Leave placeholder** — "Complete via `/flexspec` Phase 1" |
| §2.3 FR/NF | requirements, user stories, acceptance criteria, NFRs |
| §3 Tasks | task items → T-XXX list or expanded `tasks/` files |
| §4 Testing | port acceptance/checklist items only if explicit; else placeholder |
| §5 Other | unmapped content, migration notes, "run `/flexspec` to reach planned" |

Renumber IDs as new FR-/NF-/T-/TC- sequences. Preserve source titles/IDs in prose where helpful.

**d. Set status**

Map source status → FlexSpec (see **Status map**). Default `draft` when unknown.

```bash
flexspec status set <spec-id> --status <status>
```

Per spec decision: migrated specs normally end **`draft`** even if source was further along — user completes design via `/flexspec`. If source was `complete`, still use `draft` unless user asks to preserve mapped status.

**e. Expanded only:** create task files under `tasks/` from source tasks; use `flexspec status set` for task status when source task state exists.

### 5. Optional cleanup

If user chose delete: remove only confirmed source paths after all migrations succeed. Never delete `.flexspec/`, FlexSpec `specs_dir`, or unrelated project files.

### 6. Report

```markdown
## Migration report

| Source | New FlexSpec ID | Template | Status | Notes |
| --- | --- | --- | --- | --- |
| `<tool>:<path>` | `NNN-slug` | simple/expanded | draft | unmapped: … |

**Next steps**
1. Review migrated specs in `<specs_dir>/`.
2. Run `/flexspec` on each spec to complete code maps, testing criteria, and reach `planned`.
3. Optionally run `/flexspec-charter` if migration revealed product deltas.
4. When satisfied, archive or delete original tool dirs (if kept).
```

Run `flexspec validate` after migration.

## Detection summary

| Tool | Primary signature | Spec location | Reference |
| --- | --- | --- | --- |
| GitHub Spec Kit | `.specify/` + root `specs/` | `specs/NNN-feature/` | `references/speckit.md` |
| OpenSpec | `openspec/` | `openspec/changes/<name>/` (active work) and/or `openspec/specs/<domain>/` | `references/openspec.md` |
| LeanSpec | `.lean-spec/config.json` | `<specsDir>/NNN-name/` (default `specs/`) | `references/leanspec.md` |
| Spec Kitty | `.kittify/` + `kitty-specs/` | `kitty-specs/NNN-feature/` | `references/speckitty.md` |
| Kiro | `.kiro/specs/` | `.kiro/specs/<feature>/` | `references/kiro.md` |
| Tessl | `.tessl/` and/or `*.spec.md` | `specs/*.spec.md` (flat) | `references/tessl.md` |
| Unknown | — | user-provided | `references/generic.md` |

## Status map (common values → FlexSpec)

| Source (examples) | FlexSpec |
| --- | --- |
| draft, initial, specify, proposed, planned (authoring) | `draft` |
| refined, approved, ready, spec-complete | `planned` |
| in-progress, in_progress, implementing, active, underway | `in_progress` |
| review, in-review, validating, analyzing | `in_review` |
| complete, done, archived, shipped, accepted | `complete` |

Tool-specific values: see each reference doc. When ambiguous → `draft` + note in §5.

## Forbidden

- Hand-create `specs/NNN-*` directories or copy FlexSpec templates manually
- Edit spec `status` in frontmatter by hand (use `flexspec status set`)
- Fabricate code maps, mermaid traces, or test criteria during migration
- Web-search without user permission
- Delete source files without explicit per-run confirmation

## References index

Read the matching file before mapping:

- `references/speckit.md` — GitHub Spec Kit
- `references/openspec.md` — OpenSpec
- `references/leanspec.md` — LeanSpec
- `references/speckitty.md` — Spec Kitty
- `references/kiro.md` — Kiro IDE
- `references/tessl.md` — Tessl SDD tile
- `references/generic.md` — unknown tools / custom layouts
