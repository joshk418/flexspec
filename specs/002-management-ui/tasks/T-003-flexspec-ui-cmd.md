---
id: T-003
name: flexspec ui command
parent_spec: '../README.md'
status: done
satisfies: [FR-001, FR-002, FR-013, NF-006]
depends_on: [T-001, T-002]
verified_by: [TC-001]
---

# T-003: flexspec ui command

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-001, FR-002, FR-013, NF-006 · **Depends on**: T-001, T-002 · **Verified by**: TC-001

## Objective

Wire `flexspec ui` Cobra command: start server, port fallback, optional browser open, block until interrupt.

## Context

Follow `cmd/list.go` pattern for cwd root. Use `github.com/pkg/browser` **only if** charter allows new dep — prefer stdlib: exec `open` / `xdg-open` / `rundll32` via `runtime.GOOS` to avoid extra dep (document in task if using stdlib only).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `cmd/ui.go` | create | Cobra command + flags |
| `cmd/root.go` | modify | `AddCommand(uiCmd)` |

## Implementation Steps

1. Flags: `--port` (default 3000), `--host` (default `127.0.0.1`), `--open` (default true), `--no-open`.
2. Try ports `port`..`port+20`; print `FlexSpec UI at http://host:port`.
3. Construct `ui.Server` from T-001/T-002 and call blocking `Run()` until SIGINT/SIGTERM.
4. If `--open` and not `--no-open`, open URL in default browser (stdlib OS exec).
5. Surface config-missing error from server init.

## Acceptance Criteria

- [ ] `flexspec ui --no-open` starts server without browser _(FR-002, NF-006)_
- [ ] Port conflict tries next port _(FR-001)_
- [ ] Missing init → non-zero exit + message _(FR-013)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-001 | integration | ui fails without `.flexspec/` | `cmd/ui_test.go` |

Run: `go test -race ./cmd/...`

## Out of Scope

- Background daemon / `--kill` (future).

## Open Questions

- None.

## References

- LeanSpec / Spec Kitty: `ui --port`, `--no-open` conventions
