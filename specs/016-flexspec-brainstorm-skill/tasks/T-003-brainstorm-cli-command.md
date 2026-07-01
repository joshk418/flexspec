---
blocks:
    - T-004
depends_on:
    - T-002
id: T-003
name: flexspec brainstorm new CLI command
parent_spec: ../README.md
satisfies:
    - FR-001
status: done
verified_by:
    - TC-004
---

# T-003: flexspec brainstorm new CLI command

> **Parent spec**: [Brainstorm skill and CLI scaffolding](../README.md) · **Status**: todo
> **Satisfies**: FR-001 · **Depends on**: T-002 · **Verified by**: TC-004
> **Blocks**: T-004

## Objective

Wire `internal/brainstorm.Create` into a new `flexspec brainstorm new <name> [--force]` Cobra command, following the `cmd/glossary.go` parent-command-with-subcommand pattern and `cmd/new.go`'s output style.

## Context

`cmd/glossary.go` shows the pattern for a parent Cobra command (`glossaryCmd`) with child subcommands registered in `init()` via `rootCmd.AddCommand` / `glossaryCmd.AddCommand`. `cmd/new.go` shows the output style for a scaffolding command: print `Created spec %s`, then indented `path:`/`template:` lines to `cmd.OutOrStdout()`.

This command is a **new pattern for this repo**: `cmd/new.go` has no companion `cmd/new_test.go` (its logic is tested via `internal/spec/create_test.go` instead), but charter §7 states "one test file per source file" as the testing standard. This task follows the written standard rather than that one outlier, so `cmd/brainstorm_test.go` is in scope.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `cmd/brainstorm.go` | create | `brainstormCmd` (parent) + `brainstormNewCmd` (subcommand), `--force` flag |
| `cmd/brainstorm_test.go` | create | Cobra command execution test in a temp project dir |
| `cmd/glossary.go` | read | Reference for parent/subcommand registration pattern |
| `cmd/new.go` | read | Reference for output style and `os.Getwd()` + arg-joining pattern |

### Workflow / Requirement Mapping

| Parent Section | Mapping |
| --- | --- |
| Workflow graph steps | §6.1 step 6 |
| Implementation plan steps | §7.3, §7.4 step 3 |
| Requirements | FR-001 |
| Tests | TC-004 |

## Implementation Steps

1. Create `cmd/brainstorm.go` with `brainstormCmd = &cobra.Command{Use: "brainstorm", Short: "Manage pre-spec brainstorm docs"}` (parent, no `RunE` — mirrors `glossaryCmd`).
2. Add `brainstormNewCmd = &cobra.Command{Use: "new [name]", Short: "Create a new brainstorm doc", Args: cobra.MinimumNArgs(1)}` with a `Long` description referencing `.flexspec/brainstorms/<slug>.md` and noting it has no status lifecycle and is not shown by `flexspec list`.
3. In `RunE`: resolve `root, err := os.Getwd()`; join args with spaces (`strings.Join(args, " ")`) as the raw name, same as `cmd/new.go`; call `brainstorm.Create(root, name, brainstormForce)`.
4. On success, print `Created brainstorm doc\n  path: %s\n` (mirroring `cmd/new.go`'s two-line output).
5. Add a package-level `var brainstormForce bool` and register `brainstormNewCmd.Flags().BoolVar(&brainstormForce, "force", false, "Overwrite an existing brainstorm doc")`.
6. In `func init()`: `rootCmd.AddCommand(brainstormCmd)`, `brainstormCmd.AddCommand(brainstormNewCmd)`.
7. Create `cmd/brainstorm_test.go`: build a temp dir, write a minimal `.flexspec/templates/brainstorm.md` fixture (e.g. `--- \nname: '{name}'\n---\n# {name}\n`), `os.Chdir` into it (restore via `t.Cleanup`), execute `brainstormNewCmd` via `cmd.SetArgs` + `cmd.Execute()` (or call `RunE` directly with a constructed command, matching whichever pattern `cmd/status_test.go` or `cmd/init_test.go` uses — check that file first for the repo's preferred Cobra-test invocation style), and assert: (a) the target file exists with expected content, (b) running again without `--force` returns an error, (c) running again with `--force` succeeds and overwrites, (d) running with zero args returns a cobra arg-count error.

## Acceptance Criteria

- [ ] `flexspec brainstorm new <name>` creates `.flexspec/brainstorms/<slug>.md` and prints the created path. _(FR-001)_
- [ ] `flexspec brainstorm new <name>` without `--force` errors on an existing file; with `--force` it overwrites.
- [ ] `flexspec brainstorm new` with no args returns a usage/arg-count error.
- [ ] Command registered under `flexspec brainstorm new`, discoverable via `flexspec brainstorm --help` and `flexspec --help`.
- [ ] All tests in "Testing" below pass.

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-004 | integration | End-to-end CLI behavior: create, overwrite guard, `--force`, missing-arg error | `cmd/brainstorm_test.go` |

Run: `go test ./cmd/... -race -run TestBrainstorm`

## Out of Scope

- No `flexspec brainstorm list` subcommand (explicit non-goal, spec §1 out of scope — agents read `.flexspec/brainstorms/` directly).
- No changes to `flexspec validate` or `flexspec list` (FR-007, NF-004).
- No frontmatter `status` handling — `flexspec status set` must not be extended to accept brainstorm docs.

## Open Questions

None.

## References

- Parent spec: [`../README.md`](../README.md) §7.3, §7.4 step 3, §8 TC-004
- Related tasks: T-002 (wrapped by this command), T-004 (skill that invokes this command)
- Pattern reference: `cmd/glossary.go`, `cmd/new.go`
