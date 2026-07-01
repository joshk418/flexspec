---
blocks:
    - T-003
depends_on:
    - T-001
id: T-002
name: internal/brainstorm package
parent_spec: ../README.md
satisfies:
    - FR-001
    - NF-001
    - NF-002
status: done
verified_by:
    - TC-001
    - TC-002
    - TC-003
---

# T-002: internal/brainstorm package

> **Parent spec**: [Brainstorm skill and CLI scaffolding](../README.md) · **Status**: todo
> **Satisfies**: FR-001, NF-001, NF-002 · **Depends on**: T-001 · **Verified by**: TC-001, TC-002, TC-003
> **Blocks**: T-003

## Objective

Implement `internal/brainstorm.Create`, which scaffolds `.flexspec/brainstorms/<slug>.md` from the project-local `.flexspec/templates/brainstorm.md`, creating the `brainstorms/` directory on demand and refusing to overwrite an existing file unless forced.

## Context

Follow the plain-`os`-based style used by `internal/glossary` (no `fileSystem` interface abstraction — that pattern is specific to `internal/spec` for its sequence-number scanning needs, not required here). Tests use `t.TempDir()`.

Mirror `internal/spec/create.go`'s `templatePathFor`/read/write shape, but simpler: no sequence numbering, no `tasks/` subdirectory, no `type` frontmatter injection. Reuse `spec.Slugify` (already exported from `internal/spec`) instead of duplicating slug sanitization — do not copy the regex logic.

Reference behavior from `internal/spec/create.go`:
- Template path: `.flexspec/templates/brainstorm.md` (relative to project root).
- Missing-template error should mirror `spec.Create`'s "template not found; run `flexspec init` first" pattern, but must additionally mention `flexspec update --migrate` (existing projects that already ran `init` before this feature shipped need `update --migrate`, not `init`, per spec §3 Operational cases).
- Target path: `.flexspec/brainstorms/<slug>.md` — create `.flexspec/brainstorms/` with `MkdirAll` if absent (same as `specs_dir` being "created on demand").
- Overwrite guard: `Stat` the target first; if it exists and `force` is false, return an error without touching the file; if `force` is true, overwrite.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/brainstorm/brainstorm.go` | create | `Create(root, slug string, force bool) (Result, error)` |
| `internal/brainstorm/brainstorm_test.go` | create | Table-driven tests |
| `internal/spec/create.go` | read | Reference implementation pattern; reuse `spec.Slugify` from this package |

### Workflow / Requirement Mapping

| Parent Section | Mapping |
| --- | --- |
| Workflow graph steps | §6.1 steps 6–9 |
| Implementation plan steps | §7.4 step 2 |
| Requirements | FR-001, NF-001, NF-002 |
| Tests | TC-001, TC-002, TC-003 |

## Implementation Steps

1. Create `internal/brainstorm/brainstorm.go` with package doc comment (one line, no narration) and imports (`fmt`, `os`, `path/filepath`, `github.com/joshk418/flexspec/internal/spec` for `spec.Slugify`).
2. Define `type Result struct { Slug string; Path string }`.
3. Define constants: `flexspecDir = ".flexspec"`, `templatesDir = "templates"`, `brainstormTemplate = "brainstorm.md"`, `brainstormsDir = "brainstorms"`, `dirPerm = 0o755`, `filePerm = 0o644`.
4. Implement `func Create(root, name string, force bool) (Result, error)`:
   - `slug, err := spec.Slugify(name)` — propagate error unchanged.
   - Build `templatePath := filepath.Join(root, flexspecDir, templatesDir, brainstormTemplate)`; `Stat` it; if missing, return an error: `fmt.Errorf("brainstorm template not found at %s; run \`flexspec init\` (new projects) or \`flexspec update --migrate\` (existing projects) first", templatePath)`.
   - Read template bytes via `os.ReadFile`.
   - Build `brainstormsPath := filepath.Join(root, flexspecDir, brainstormsDir)`; `os.MkdirAll(brainstormsPath, dirPerm)`.
   - Build `targetPath := filepath.Join(brainstormsPath, slug+".md")`.
   - `Stat(targetPath)`: if it exists and `!force`, return `fmt.Errorf("brainstorm doc %s already exists; use --force to overwrite", targetPath)` without writing.
   - `os.WriteFile(targetPath, data, filePerm)`.
   - Return `Result{Slug: slug, Path: targetPath}, nil`.
5. Wrap every returned error with `%w` where an underlying error is wrapped (per charter §7 code conventions); document the exported `Create` function and `Result` type.

## Acceptance Criteria

- [ ] `Create` writes `.flexspec/brainstorms/<slug>.md` with the template's exact bytes when the target doesn't exist. _(FR-001)_
- [ ] `Create` creates `.flexspec/brainstorms/` on demand when it doesn't exist. _(FR-001)_
- [ ] `Create` returns an error and does not modify the existing file when the target exists and `force` is false. _(NF-001)_
- [ ] `Create` overwrites the existing file when `force` is true.
- [ ] `Create` returns an actionable error mentioning both `flexspec init` and `flexspec update --migrate` when the template file is missing.
- [ ] `Create` reuses `spec.Slugify`; no duplicated slug-sanitization logic exists in this package.
- [ ] All tests in "Testing" below pass.

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-001 | unit | `Create` scaffolds the file + directory from the template | `internal/brainstorm/brainstorm_test.go` |
| TC-002 | unit | Overwrite guard: errors without `force`, succeeds with `force` | `internal/brainstorm/brainstorm_test.go` |
| TC-003 | unit | Actionable error when `.flexspec/templates/brainstorm.md` is missing | `internal/brainstorm/brainstorm_test.go` |

Run: `go test ./internal/brainstorm/... -race`

## Out of Scope

- No Cobra/CLI wiring here (T-003).
- No slug-matching or ingestion logic for `/flexspec` (that lives in the skill file, T-005, not Go code).
- Do not add a `fileSystem` interface abstraction — plain `os` calls matching `internal/glossary`'s style are sufficient here.

## Open Questions

None.

## References

- Parent spec: [`../README.md`](../README.md) §7.1, §7.4 step 2, §8 TC-001–TC-003
- Related tasks: T-001 (template this reads), T-003 (CLI command wrapping this)
- Pattern reference: `internal/spec/create.go`, `internal/glossary/glossary.go`
