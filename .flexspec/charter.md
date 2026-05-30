---
product_name: "FlexSpec"
version: "0.1"
last_updated: "2026-05-30"
status: active
---

# FlexSpec

> **Charter status**: active · **Version**: 0.1 · **Last updated**: 2026-05-30

## 1. Product overview

**One-liner:** A spec-driven development CLI (Go) for generating and tracking feature specifications via markdown templates, with optional adapters for external issue trackers.

**Problem:** Most spec-driven development tools cause prompt and context fatigue by introducing too many files, which erodes token efficiency over time. FlexSpec keeps workflows simple — a single file per spec for simple tasks, and a more expanded multi-file structure for complex tasks, but only expanded enough to add the context that complexity demands. Users keep the "flexibility": they can override which spec structure is created and freely modify templates and configuration so the system fits how they work.

**Intended outcome:** Teams adopt a full spec-driven workflow that documents changes across an ever-evolving project. As features are added, the spec corpus and charter stay current, keeping humans and AI agents informed and aligned. FlexSpec skills ensure agents update only the charter, never external files outside the system.

## 2. Vision and goals

**North star:** Keep humans and AI agents aligned on intent without context fatigue.

**Success criteria:**

- **Adoption** — installs, `flexspec init` runs, and GitHub stars trend up.
- **Reduced agent drift** — fewer off-spec edits per implementation.
- **Retention** — projects keep adding specs after 30 days.

## 3. Users and stakeholders

FlexSpec serves both solo developers and teams, but the desired outcome is **adoption within teams**.

| Persona | Role | Primary needs |
| --- | --- | --- |
| Solo developer | Builds with AI coding agents | Quick spec creation, focused implementation |
| Development team | Shared spec discipline | Consistent specs, living docs, agent alignment |
| AI coding agent | Consumer of specs (Cursor, Codex, Claude, Pi, Zed, etc.) | Clear, structured context to stay on-track during implementation |
| Engineering leads / Product managers | Oversight & process | Visibility into what is being built and why |

**Jobs to be done:**

- Quick spec creation for a new feature.
- Focused implementation that stays within the spec's intent.
- Keep documentation and charter current as the project evolves.

## 4. Capabilities

**Available today:**

- Spec scaffolding — simple (single-file) and expanded (multi-file) templates.
- Charter management — product-wide context authored via `/flexspec-charter`.
- CLI — `flexspec init`, `flexspec new`, `flexspec list` (`--json` for machine output), `flexspec validate`, `flexspec ui` (local management dashboard), `flexspec status set` (update spec/task status in frontmatter).
- Management UI — `flexspec ui` serves an embedded React app: kanban/table board by spec status, spec browser with markdown rendering, settings for UI prefs and `.flexspec/config.yaml`; live refresh via filesystem watch (SSE).
- Agent skills — `/flexspec` (spec lifecycle) and `/flexspec-charter` (application charter).
- Configuration and template overrides — users control spec structure via config (`spec_template`) and a per-spec skill flag (`--template`); templates are freely editable.

**Planned:**

- Adapters for external systems (Jira, Shortcut, GitHub Issues, and more).

## 5. Technical context

- **Language/runtime:** Go 1.26.2.
- **CLI framework:** `spf13/cobra`.
- **Config/data:** YAML (`gopkg.in/yaml.v3`); markdown-first spec and charter files.
- **Templates:** bundled via `embed.FS` and scaffolded on `init`.
- **Distribution:** `go install github.com/joshk418/flexspec@latest`; skills installed via `npx skills`.

**Constraints agents must respect:**

- Go ≥ 1.26 floor (CI uses the `go.mod` version).
- Minimal dependencies — Cobra, `yaml.v3`, and `fsnotify` (UI file watch); avoid heavy new deps. React UI is embedded at build time; end users do not need Node at runtime.
- Skills write only inside `.flexspec/` and the configured spec directory; agents may modify code files during implementation but must not touch `README`, `AGENTS.md`, or related docs unless explicitly instructed.
- `init` never clobbers user edits unless `--force` is passed.
- Cross-platform — build paths with `filepath`.
- CI gate: `go test -race`, `gofmt`, `go vet`, `golangci-lint`.

## 6. Architecture

`main` embeds templates and the built web UI (`ui/dist`). The Cobra CLI scaffolds, lists, validates, and optionally serves the management UI (`flexspec ui` → `internal/ui` HTTP server on localhost). Agent skills read the charter and drive the spec lifecycle. Future adapters sit behind a spec-source interface.

```mermaid
flowchart TD
    main[main + embed templates + ui/dist] --> cli[Cobra CLI]
    cli -->|init| fs[.flexspec/: config, charter, templates]
    cli -->|new| specs[specs_dir/NNN-slug/]
    cli -->|list| specs
    cli -->|validate| fs
    cli -->|validate| specs
    cli -->|ui| uiSrv[internal/ui HTTP + SSE]
    uiSrv --> specs
    uiSrv --> fs
    browser[Browser] --> uiSrv
    skills[Agent skills: /flexspec, /flexspec-charter] -->|read| fs
    skills -->|author / implement / review| specs
    adapters[(Adapters — planned)] -.-> skills
```

**Boundaries:** the CLI scaffolds, lists, validates, and serves a **local** UI; skills handle authoring, implementation, and review. The UI is visibility and convenience, not a hosted PM product.

## 7. Standards and conventions

- **Testing:** table-driven tests, one test file per source file, and a single table-driven test per tested function (e.g. `config_test.go`, `metadata_test.go`).
- **CI must pass:** `go test -race`, `gofmt` clean, `go vet`, `golangci-lint`.
- **Code conventions:** one Cobra command per file under `cmd/`; wrap errors with `%w`; document exported functions; no narrating comments.

## 8. Product boundaries

FlexSpec is a tool for managing specifications to keep AI coding agents (Cursor, Codex, Claude, Pi, Zed, etc.) on-track during implementation via the provided skills. It will **not**:

- Be a project-management tool or issue tracker itself.
- Be an AI agent or LLM runtime.
- Modify `README`, `AGENTS.md`, or related documentation files unless explicitly instructed (it does modify code files during implementation).
- Run as a hosted service.

## 9. Domain glossary

| Term | Definition |
| --- | --- |
| Charter | Product-wide context (this file) used by every spec. |
| Spec | A feature specification, simple or expanded, under the configured specs directory. |
| Simple spec | A single-file markdown spec for small, focused features. |
| Expanded spec | A multi-file specification for complex features, with linked task files. |
| Task file | A per-task file within an expanded spec. |
| Adapter | Pluggable connector to an external issue tracker (planned). |
| Phase | A stage in the `/flexspec` lifecycle: author, implement, or review. |
| One-shot | Running all `/flexspec` phases back-to-back without stopping (`always_one_shot` / `--one-shot`). |
| Validate | `flexspec validate` — read-only structural checks on `.flexspec/`, templates, and specs (exit 1 on errors). |
| Management UI | `flexspec ui` — local dashboard (board, spec browser, settings) with live filesystem sync. |
| SSE | Server-sent events from `flexspec ui` when spec files change on disk. |

## 10. Assumptions and open questions

**Assumptions:**

- Users run AI coding agents that support skills.
- Specs and the charter live in git.
- One charter per repository.

**Open questions (blocking):**

- None.

**Open questions (non-blocking):**

- Which issue-tracker adapter ships first.

## 11. Revision history

| Date | Summary | Source |
| --- | --- | --- |
| 2026-05-30 | Initial charter authored — product overview, vision, users, capabilities, technical context, architecture, standards, boundaries, glossary. | /flexspec-charter |
| 2026-05-30 | §4/§6/§9 — document full CLI (`init`, `new`, `list`, `validate`); architecture diagram updated. | 001-cli-validate |
| 2026-05-30 | §4/§5/§6/§9 — management UI (`flexspec ui`), `list --json`, `status set`; architecture + glossary. | 002-management-ui |
| 2026-05-30 | §4 — compact list output by default; `--json` still has full detail. | 003-simplify-list-command |
