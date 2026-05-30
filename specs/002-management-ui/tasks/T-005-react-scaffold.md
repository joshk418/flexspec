---
id: T-005
name: React app scaffold
parent_spec: '../README.md'
status: done
satisfies: [FR-003, NF-001]
depends_on: [T-001]
verified_by: [TC-002]
---

# T-005: React app scaffold

> **Parent spec**: [Management UI](../README.md) · **Status**: todo
> **Satisfies**: FR-003, NF-001 · **Depends on**: T-001 · **Verified by**: TC-002

## Objective

Create `web/` Vite + React + TypeScript app with routing, nav shell, and typed API client.

## Context

Dev: Vite on `:5173` proxies `/api` → `http://127.0.0.1:3000`. Prod: built to `web/dist` for embed. Keep dependencies lean (react-router, no heavy UI kit required — CSS modules or minimal Tailwind optional).

### Files In Scope

| File | Action | Notes |
| --- | --- | --- |
| `web/package.json` | create | scripts: `dev`, `build` |
| `web/vite.config.ts` | create | proxy `/api` |
| `web/src/main.tsx` | create | entry |
| `web/src/App.tsx` | create | layout + `<Routes>` |
| `web/src/api/client.ts` | create | `fetchSpecs`, `fetchSpec`, `subscribeEvents` |
| `web/src/pages/BoardPage.tsx` | create | stub |
| `web/src/pages/SpecsPage.tsx` | create | stub |
| `web/src/pages/SettingsPage.tsx` | create | stub |

## Implementation Steps

1. `npm create vite@latest web -- --template react-ts` (or equivalent manual scaffold).
2. Routes: `/board`, `/specs`, `/specs/:dir`, `/settings`; redirect `/` → `/board`.
3. Top nav links between three views.
4. API client with base URL `''` (same origin) in prod.
5. `EventSource` helper for `/api/events` calling optional callback `onSpecsChanged`.
6. Placeholder pages render "TODO" until T-006–T-008.

## Acceptance Criteria

- [ ] `npm run build` produces `web/dist/index.html` _(NF-001 build-time only)_
- [ ] Dev proxy reaches Go API _(FR-003 dev path)_

## Testing

Manual: `go run . ui` + `cd web && npm run dev` with proxy.

## Out of Scope

- Full board/settings implementation.

## Open Questions

- None.

## References

- Parent spec: [UI routes](../README.md#24-external-interfaces)
