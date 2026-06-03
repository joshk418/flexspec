---
created: "2026-06-03"
description: Make §3.1 implementation code maps optional in FlexSpec skill and templates; require only for high-complexity work.
implementation_finished: "2026-06-03"
implementation_start: "2026-06-03"
name: optional-implementation-code-map
priority: medium
spec_type: simple
status: complete
tags:
    - templates
    - skills
    - token-efficiency
task_count: 5
---

# Optional implementation code map

> **Status**: complete · **Priority**: medium · **Created**: 2026-06-03 · **Tasks**: 5

## 1. Summary

FlexSpec currently scaffolds every new spec with a full **§3.1 Implementation Code Map**
(mermaid diagram + task execution table). For many simple features, the **§3.2 Task
List** already conveys build order, files, and requirement mapping — the §3.1 visual
duplicates that and wastes tokens in the spec file.

**Problem:** Spec token budgets and author fatigue increase without adding review
value when implementation order is linear and obvious from the task list.

**Outcome:** Newly scaffolded specs omit §3.1 by default. The `/flexspec` skill
requires §3.1 only when a documented **complexity heuristic** applies (auth systems,
large refactors, cross-cutting work, parallel task tracks, etc.). §2.2 design code
maps stay required — only the **implementation** map becomes conditional.

**In scope:** `skills/flexspec/SKILL.md`; repo-root `templates/flexspec-simple.md`,
`templates/expanded/flexspec-expanded.md`, `templates/README.md`; mirrored
`.flexspec/templates/` copies; optional charter §4 wording for the behavior.

**Out of scope:** Retroactive edits to existing specs (`001`–`009`); changes to
`/flexspec-migrate` or other skills; `flexspec validate` structural rules (no §3.1
checks today); migrations that rewrite user specs; bundled skill reinstall beyond
normal `flexspec update --skills`.

## 2. Design

### 2.1 Architecture / Technical Plan

Documentation-only change: templates are embedded from `templates/` via `main.go`
`embed.FS` and copied to `.flexspec/templates/` on `init`. The flexspec skill is
shipped under `skills/flexspec/` and installed globally by `flexspec update --skills`.
No Go code changes unless a future validate rule is added (explicitly out of scope).

| File / Component | Type | Role in this spec |
| --- | --- | --- |
| `skills/flexspec/SKILL.md` | modified | Code Map Quality Bar, Phase 1 DoR, §3 authoring rules, complexity heuristic |
| `templates/flexspec-simple.md` | modified | Remove default §3.1 scaffold; optional guidance comment |
| `templates/expanded/flexspec-expanded.md` | modified | Same as simple |
| `templates/README.md` | modified | Code map conventions: §3.1 optional + when to include |
| `.flexspec/templates/flexspec-simple.md` | modified | Keep in sync with `templates/` |
| `.flexspec/templates/expanded/flexspec-expanded.md` | modified | Keep in sync |
| `.flexspec/templates/README.md` | modified | Keep in sync with `templates/README.md` |
| `.flexspec/charter.md` | modified (optional) | §4 note that §3.1 is optional for low-complexity specs |

### 2.2 Code Map

```mermaid
sequenceDiagram
    autonumber
    participant Author as Agent /flexspec Phase 1
    participant Skill as skills/flexspec/SKILL.md
    participant CLI as flexspec new
    participant Embed as templates/*.md embed.FS
    participant Proj as .flexspec/templates/*.md
    participant Spec as specs/NNN-*/README.md

    Author->>Skill: read Code Map + heuristic FR-002
    Author->>CLI: flexspec new name --template simple|expanded
    CLI->>Embed: read template bytes
    Embed-->>CLI: template without §3.1 body
    CLI->>Spec: write README.md scaffold
    Author->>Author: fill §2.2 + §3.2; add §3.1 only if heuristic
    Author->>Proj: sync .flexspec/templates on init/update (existing behavior)
```

| Step | Location | Executes | Input / condition | Output / side effect | FR/NF |
| --- | --- | --- | --- | --- | --- |
| 1 | `skills/flexspec/SKILL.md` | author reads rules | feature scope | knows §3.1 optional | FR-001, FR-002 |
| 2 | `internal/spec/create.go` :: create | `flexspec new` | template name | spec dir + README from embed | FR-003 |
| 3 | `templates/*.md` | embed read | init/new | bytes served to CLI | FR-003 |
| 4 | `specs/.../README.md` | author edits | planned spec | §3.2 always; §3.1 if complex | FR-001, FR-002, FR-004 |
| 5 | `.flexspec/templates/` | file copy on init | user run init | local template mirror | FR-003 |

### 2.3 Requirements

**Functional**

- **FR-001** — `/flexspec` Phase 1 must treat §3.1 as **optional** by default; §3.2 Task List remains required for all specs.
- **FR-002** — The skill must define a **complexity heuristic** listing when §3.1 is required (e.g. auth/security-critical flows, large refactors, cross-cutting multi-subsystem work, parallel `depends_on` branches, multiple §2.2 execution paths, expanded specs with many interdependent tasks, or user-requested visual map).
- **FR-003** — `flexspec new` scaffolds must **not** include a filled §3.1 example (diagram + table); optional HTML guidance may explain when to add §3.1.
- **FR-004** — When §3.1 is omitted, §3.2 tasks must still cover every §2.1 file and every §2.2 step (via task text: files, `depends_on`, and §2.2 step references); record omission reason in §5 Other (e.g. `§3.1 omitted: linear 3-task CLI change`).
- **FR-005** — When §3.1 is included, existing Code Map Quality Bar for implementation maps (mermaid + matching task execution table, linkage to §2.2) remains unchanged.
- **FR-006** — `templates/README.md` and `.flexspec/templates/README.md` document optional §3.1 and point to the skill heuristic.

**Non-Functional**

- **NF-001** — Simple template scaffold token footprint must drop materially vs current (no default §3.1 mermaid/table block).
- **NF-002** — No retroactive migration of existing specs or installed migrate skill behavior.

## 3. Implementation Plan

### 3.2 Task List

- **T-001** — Update `skills/flexspec/SKILL.md`: split §2.2 (required) vs §3.1 (conditional); add complexity heuristic; update DoR, anti-patterns, and §3 Implementation Plan bullets. _(satisfies: FR-001, FR-002, FR-004, FR-005)_
- **T-002** — Update `templates/flexspec-simple.md` and `templates/expanded/flexspec-expanded.md`: remove default §3.1 diagram/table; add short optional §3.1 comment block; keep §3.2 as primary implementation section. _(satisfies: FR-003, NF-001)_
- **T-003** — Mirror template + README changes under `.flexspec/templates/`. _(satisfies: FR-003, FR-006)_
- **T-004** — Update `templates/README.md` code-map table and linkage notes for optional §3.1. _(satisfies: FR-006)_
- **T-005** — Run `flexspec new` smoke spec + `flexspec validate`; confirm scaffold has no §3.1 body and skill/template docs agree. _(satisfies: FR-003, TC-001, TC-002)_

## 4. Testing Criteria

| Test ID | Verifies | Description | Type |
| --- | --- | --- | --- |
| TC-001 | FR-003, NF-001 | `flexspec new tmp-smoke --template simple` → README has §3.2, no §3.1 mermaid/table example block | manual |
| TC-002 | FR-001, FR-002 | Skill text states §3.1 optional + lists ≥4 complexity triggers | manual |
| TC-003 | FR-004 | Sample omitted-§3.1 spec: §3.2 tasks reference §2.2 steps/files; §5 notes omission | manual |
| TC-004 | FR-005 | Skill still requires full §3.1 diagram+table when heuristic triggers | manual |
| TC-005 | FR-006 | `templates/README.md` marks §3.1 optional with pointer to skill | manual |
| TC-006 | NF-002 | No edits under `specs/001`–`009` or `skills/flexspec-migrate` | manual |

## 5. Other

- **Assumption:** §2.2 design code maps remain mandatory; user request targets implementation maps only.
- **Charter follow-up:** Consider §4 bullet that §3.1 is optional for low-complexity specs (aligns with charter token-efficiency goal). Deferred unless user requests charter edit during implementation.
- **Risk:** Authors skip §3.1 on genuinely complex work — mitigated by explicit heuristic + Phase 1 reviewer test in skill.
- **Rollout:** Only affects specs created after templates/skills ship; users run `flexspec update --skills` to refresh global skill.
