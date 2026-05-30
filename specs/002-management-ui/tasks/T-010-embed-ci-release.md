---
id: T-010
name: Embed UI build and release pipeline
parent_spec: '../README.md'
status: done
satisfies: [NF-001, NF-004]
depends_on: [T-005, T-006, T-007, T-008, T-003]
verified_by: [TC-009]
---

# T-010: Embed UI build and release pipeline

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: NF-001, NF-004 · **Depends on**: T-005–T-008, T-003 · **Verified by**: TC-009

## Objective

Embed `web/dist` in Go binary; ensure CI and GoReleaser build UI before compile.

## Context

Pattern mirrors `templates` embed in `main.go`. Use `//go:embed web/dist/*` — build tag or stub `dist` for `go test` without npm (commit minimal placeholder `web/dist/.gitkeep` + CI always builds real dist).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `main.go` | modify | embed `web/dist`, pass to cmd |
| `internal/ui/embed.go` | modify | `http.FileServer` + SPA fallback |
| `.github/workflows/ci.yml` | modify | Node setup + `npm ci && npm run build` in `web/` |
| `.goreleaser.yaml` | modify | before hook builds UI |
| `Makefile` | create | `build-ui`, `build` targets |
| `.gitignore` | modify | ignore `web/node_modules`, optionally keep `web/dist` in CI artifacts only |

## Implementation Steps

1. Wire `cmd.UIFS embed.FS` from main.
2. `handleStatic`: serve files; unknown paths → `index.html`.
3. CI: install Node LTS, cache npm, build `web/` before `go test`.
4. GoReleaser `before: hooks` run UI build script.
5. Document in README: contributors run `make build` or `make build-ui` before local `go build`.
6. For tests without embed: use `httptest` with manual FS or build tag `noui`.

## Acceptance Criteria

- [ ] Released binary serves board UI without Node installed _(NF-001)_
- [ ] CI fails if `web/dist` missing when required _(NF-004)_

## Testing

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-009 | ci | workflow builds web | `.github/workflows/ci.yml` |
| TC-002 | integration | GET / returns index | `internal/ui/server_test.go` |

## Out of Scope

- Publishing separate `@flexspec/ui` npm package.

## Open Questions

- None.

## References

- `main.go` templates embed pattern
