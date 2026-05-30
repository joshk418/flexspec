---
id: T-001
name: UI server and read API
parent_spec: '../README.md'
status: done
satisfies: [FR-001, FR-003, FR-013, FR-014, NF-002]
depends_on: []
verified_by: [TC-001, TC-002, TC-003]
---

# T-001: UI server and read API

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-001, FR-003, FR-013, FR-014, NF-002 · **Verified by**: TC-001, TC-002, TC-003

## Objective

Add `internal/ui` with an HTTP server exposing read-only JSON for specs and config, plus a placeholder static handler (wired fully in T-010).

## Context

Reuse `config.Load` and `spec.List` from project root (cwd). Register routes on Go 1.22+ `http.ServeMux` with method patterns. Default listen `127.0.0.1:3000`.

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `internal/ui/server.go` | create | `Server` struct, `ListenAndServe`, shutdown |
| `internal/ui/handlers.go` | create | `handleHealth`, `handleSpecs`, `handleSpecDetail`, `handleConfigGet` |
| `internal/ui/server_test.go` | create | httptest for handlers |
| `internal/config/config.go` | modify | optional: export fields only, no Save yet |

## Implementation Steps

1. Define `Server{ Root string, Addr string, SpecFS http.FileSystem }` and `New(root, addr)`.
2. `GET /api/health` → `{"ok":true}`.
3. `GET /api/specs` → JSON encode `spec.List(root, cfg)` with stable struct tags (`id`, `dir`, `name`, `description`, `status`, `spec_type`, `tasks`).
4. `GET /api/specs/{dir}` → read `README.md` body (after frontmatter) + task files with bodies for expanded specs; 404 if missing.
5. `GET /api/config` → JSON from `config.Load`.
6. If config missing on server start, return error from `New` / `Run` with message mentioning `flexspec init`.
7. Stub `handleStatic` — serve empty or minimal placeholder until T-010 embed exists.

## Acceptance Criteria

- [ ] Handlers return correct JSON for fixture project with `001-*` spec _(FR-003)_
- [ ] Server binds to configured host (default loopback) _(FR-014)_
- [ ] Missing `.flexspec/config.yaml` prevents server start with clear error _(FR-013)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-002 | integration | `/api/health` 200 | `internal/ui/server_test.go` |
| TC-003 | integration | `/api/specs` lists fixture spec | `internal/ui/server_test.go` |

Run: `go test -race ./internal/ui/...`

## Out of Scope

- fsnotify, SSE, write APIs, `cmd/ui.go`, React.

## Open Questions

- None.

## References

- Parent spec: [§2.4](../README.md#24-external-interfaces)
- `internal/spec/spec.go` — `List`, `ParseSpecMeta`
