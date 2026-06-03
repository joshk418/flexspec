# Generic / unknown SDD tool

Use when:

- No embedded signature matches (see main SKILL.md detection table)
- User names a tool not in `references/`
- Signature matched but on-disk layout differs from the reference doc
- User denied web search but supplied a custom directory path

## Interview (before migrating)

Ask user (batch 2–4 questions):

1. **Source directory path** — absolute or repo-relative path containing specs.
2. **Migratable unit** — what counts as one spec? (single file, folder, issue id, etc.)
3. **Main content files** — which files hold requirements, design, tasks, tests?
4. **Status field** — where is status stored? (frontmatter key, filename, external tracker)
5. **Template guess** — single doc vs multi-file/tasks? (confirm simple vs expanded)

Optional: if user **grants web permission**, search official docs once for directory layout, then record findings in migration report (do not add new reference file unless user asks).

## Detection (custom path)

Given user path `<PATH>`:

1. List top-level entries (files + dirs).
2. Identify repeating pattern (numbered folders, markdown with frontmatter, etc.).
3. Build inventory table for user confirmation.

## Template inference

| Pattern | FlexSpec |
| --- | --- |
| One markdown file per feature, no task breakdown files | `simple` |
| Folder with 2+ docs OR `tasks/` OR checklist task file | `expanded` |

## Generic field mapping → FlexSpec

Map by **semantic role**, not filename (filenames vary by tool):

| Semantic role | Common filenames | FlexSpec target |
| --- | --- | --- |
| Intent / problem / scope | `spec.md`, `README.md`, `requirements.md`, `proposal.md` | §1 Summary |
| Out of scope | non-goals section, `OUT_OF_SCOPE.md` | §1 out-of-scope |
| Requirements | FR lists, user stories, EARS, scenarios | §2.3 FR-* |
| Constraints / NFR | performance, security sections | §2.3 NF-* |
| Design / architecture | `plan.md`, `design.md`, `ARCHITECTURE.md` | §2.1 |
| Tasks / work items | `tasks.md`, `TASKS.md`, `tasks/*`, checklists | §3 |
| Tests / acceptance | `TESTING.md`, acceptance criteria, test links | §4 (explicit only) |
| Metadata | YAML frontmatter, `meta.json` | FlexSpec frontmatter + §5 |
| Unknown sections | — | §5 Other verbatim summary |

**Never fabricate** §2.2 code maps or TC rows without explicit source content.

## Status map (fallback)

Ask user for their tool's status values, then map:

| Typical source meaning | FlexSpec |
| --- | --- |
| not started / idea / draft | `draft` |
| approved / ready to build | `planned` |
| building / active | `in_progress` |
| review / QA | `in_review` |
| done / shipped | `complete` |

If user cannot provide mapping → **`draft`** for all migrated specs.

## Workflow reminder

Still use FlexSpec CLI only:

```bash
flexspec new <slug> --template <simple|expanded>
# edit CLI-created README.md (+ tasks/ if expanded)
flexspec status set <spec-id> --status draft
flexspec validate
```

## Slug naming

From folder or file basename; kebab-case; strip numeric prefixes (`001-`, `042-`); avoid collisions with `flexspec list`.

## Report extras

For generic migrations, report must include:

- Detected layout description (for future reference)
- Any user-provided status mapping table
- Recommendation to contribute a new `references/<tool>.md` if user will migrate again
