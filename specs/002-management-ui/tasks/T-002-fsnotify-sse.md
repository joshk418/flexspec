---
id: T-002
name: Filesystem watch and SSE
parent_spec: '../README.md'
status: done
satisfies: [FR-010, NF-002]
depends_on: [T-001]
verified_by: [TC-006]
---

# T-002: Filesystem watch and SSE

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-010, NF-002 · **Depends on**: T-001 · **Verified by**: TC-006

## Objective

Watch `specs_dir` and `.flexspec/config.yaml` for changes; broadcast debounced `specs-changed` events on `GET /api/events` (SSE).

## Context

Add `github.com/fsnotify/fsnotify` to `go.mod`. Debounce 500ms–2s to coalesce rapid agent saves. Hub pattern: register clients on SSE connect, broadcast on debounced fire.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/ui/watch.go` | create | Start/stop watcher with server lifecycle |
| `internal/ui/events.go` | create | SSE handler, client registry, broadcast |
| `internal/ui/watch_test.go` | create | debounce + event emission |
| `go.mod` | modify | `fsnotify` dependency |

## Implementation Steps

1. On server start, watch `filepath.Join(root, cfg.SpecsDir)` recursively and `filepath.Join(root, ".flexspec/config.yaml")`.
2. Ignore non-`.md` / irrelevant paths if needed for noise reduction (still catch README + tasks).
3. Implement debounced callback calling `hub.Broadcast("specs-changed")`.
4. `GET /api/events`: `Content-Type: text/event-stream`, flush `event: specs-changed\ndata: {}\n\n` on broadcast; handle client disconnect.
5. Stop watcher on server shutdown.

## Acceptance Criteria

- [ ] Editing a spec README triggers SSE within debounce window _(FR-010)_
- [ ] Only one new Go dependency (`fsnotify`) _(NF-002)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-006 | integration | write file → receive SSE | `internal/ui/watch_test.go` |

Run: `go test -race ./internal/ui/...`

## Out of Scope

- WebSocket (SSE only v1).

## Open Questions

- None.

## References

- Parent spec: [FR-010](../README.md)
