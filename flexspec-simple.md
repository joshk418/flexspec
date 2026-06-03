---
name: ""
description: ""
status: [draft, planned, in_progress, in_review, complete]
created: "{datetime}"
implementation_start: "datetime"
implementation_finished: "datetime"
priority: [low, medium, high, critical]
tags: []
spec_type: simple
---

# {name}

> **Status**: {status} · **Priority**: {priority} · **Created**: {date}

<!--
SIMPLE SPEC — keep this file under ~2000 tokens. If it won't fit, the work is
likely too large for a simple spec; use the expanded template instead.
-->

## 1. Summary

<!--
Detailed, high-level overview of what this spec delivers and why.
Cover: the problem, the intended outcome, who/what it affects, and the scope
boundaries (what is explicitly out of scope). Write so a reader unfamiliar with
the feature understands the goal before reading the design.
-->

{summary}

## 2. Design

<!-- All moving parts and flows of the feature, detailed end to end. -->

### 2.1 Architecture / Technical Plan

<!--
Detailed description of how the spec will be implemented. Reference concrete
files, packages, and components an implementer (human or LLM) must touch or read.
List every relevant file in the table below.
-->

{architecture_overview}

| File / Component | Type                       | Role in this spec             |
| ---------------- | -------------------------- | ----------------------------- |
| `path/to/file`   | new / modified / reference | What it does / why it matters |

### 2.2 Code Map

<!--
CODE EXECUTION MAP — how the running program moves through this feature.

A human or LLM reviewer should be able to step through execution like a debugger:
who runs first, what each symbol does, what data crosses each boundary, and where
control branches (success vs error). Do not restate §2.1 as a box diagram.

Required (diagram + trace table):
1. Mermaid showing ordered execution (prefer `sequenceDiagram` with `autonumber`
   for call order; use `flowchart` when loops/branches are clearer). Include
   `alt`/`opt`/`else` (or labeled branch edges) for non-happy paths when they
   exist.
2. Execution trace table below the diagram — one row per step, same numbering.

Diagram rules:
- Every step: `path/to/file :: symbol` (or route/CLI/event if symbol TBD).
- Edge/participant labels: verb + payload (`calls create(dto)`, `returns 201`,
  `throws ValidationError`, `reads row`).
- `subgraph` or sequence participants for boundaries (client, app, data, external).
- Map FR-XXX / NF-XXX on steps where behavior is satisfied or constrained.

Trace table columns (required): Step | Location | Executes | Input / condition |
Output / side effect | FR/NF

Avoid: generic nodes; steps with no executable symbol; diagram without matching
trace rows. Note uncertain symbols in §5 Other.

Replace examples below with repo-accurate execution.
-->

```mermaid
sequenceDiagram
    autonumber
    participant Client
    participant Route as routes/order.ts
    participant Handler as handlers/order.ts::orderHandler
    participant Service as services/order.ts::OrderService.create
    participant Repo as repos/order.ts::OrderRepository.insert
    participant DB as orders

    Client->>Route: POST /api/orders (body)
    Route->>Handler: dispatch request
    Handler->>Service: create(dto) FR-001
    alt invalid dto
        Service-->>Handler: ValidationError NF-001
        Handler-->>Client: 400
    else valid
        Service->>Service: validate(dto)
        Service->>Repo: insert(order)
        Repo->>DB: INSERT
        DB-->>Repo: row
        Repo-->>Service: Order
        Service-->>Handler: 201 + id FR-001
        Handler-->>Client: JSON body
    end
```

| Step | Location | Executes | Input / condition | Output / side effect | FR/NF |
| --- | --- | --- | --- | --- | --- |
| 1 | `routes/order.ts` | route match | `POST` + JSON body | dispatches to handler | — |
| 2 | `handlers/order.ts :: orderHandler` | handler entry | request DTO | calls `create` | — |
| 3 | `services/order.ts :: OrderService.create` | business logic | dto | validates or errors | FR-001, NF-001 |
| 4 | `repos/order.ts :: OrderRepository.insert` | persistence | `Order` entity | SQL INSERT | FR-001 |
| 5 | `orders` (table) | store row | INSERT | row returned | — |
| 6 | `handlers/order.ts :: orderHandler` | response | `Order` | `201` + JSON to client | FR-001 |

### 2.3 Requirements

<!--
Functional requirements (what the system must do) and non-functional
requirements (performance, security, reliability, UX constraints).
Use stable IDs so other sections and tasks can reference them.
-->

**Functional**

- **FR-001** — {requirement}
- **FR-002** — {requirement}

**Non-Functional**

- **NF-001** — {requirement}
- **NF-002** — {requirement}

## 3. Implementation Plan

<!--
Built off the technical plan. Enumerate everything that will be created or
modified to complete the spec. Each task gets a stable internal ID (T-001...)
so the FlexSpec system can reference, track, and order it.
-->

### 3.1 Implementation Code Map

<!--
IMPLEMENTATION + EXECUTION ENABLEMENT MAP.

Show (1) task build order, (2) which §2.2 execution steps each task implements or
unblocks, and (3) symbols/files touched. Reviewers should see how the runtime
path from §2.2 becomes executable as tasks land.

Required:
- Mermaid: task nodes with `T-XXX :: file :: symbol`; solid edges = build order;
  dotted edges = "enables §2.2 step N" (or step range).
- Task execution table below diagram (required).

Task execution table columns: Task | Build after | Implements §2.2 steps |
Symbols added/changed | Execution unlocked |

Every §2.2 step must be owned by ≥1 task. Every §2.1 file on a task row.

Replace examples below.
-->

```mermaid
flowchart TB
    subgraph exec [§2.2 runtime steps enabled]
        e4["steps 4-5: Repo.insert → DB"]
        e6["step 6: handler response"]
    end
    subgraph build [Task build order]
        T001["T-001 :: migrations/001_orders.sql"]
        T002["T-002 :: repos/order.ts :: OrderRepository.insert"]
        T003["T-003 :: services/order.ts :: OrderService.create"]
        T004["T-004 :: handlers/order.ts :: orderHandler"]
        T005["T-005 :: tests/order.test.ts"]
    end

    T001 --> T002 --> T003 --> T004 --> T005
    T002 -.->|enables| e4
    T003 -.->|enables| e4
    T004 -.->|enables| e6
```

| Task | Build after | Implements §2.2 steps | Symbols added/changed | Execution unlocked |
| --- | --- | --- | --- | --- |
| T-001 | — | — | schema `orders` | persistence possible |
| T-002 | T-001 | 4–5 | `OrderRepository.insert` | DB write path runs |
| T-003 | T-002 | 3, 4–5 | `OrderService.create`, `validate` | FR-001 logic + repo calls |
| T-004 | T-003 | 1–2, 6 | `orderHandler` | full HTTP path end-to-end |
| T-005 | T-004 | 1–6 (assert) | `tests/order.test.ts` | TC coverage of trace |

### 3.2 Task List

- **T-001** — {task description} _(satisfies: FR-001)_
- **T-002** — {task description} _(satisfies: FR-002, NF-001)_
- **T-003** — {task description}

## 4. Testing Criteria

<!--
Every piece of functionality must be testable. Define the tests that prove each
requirement is met. If something cannot be tested, rework the implementation
plan (Section 3) until it can. Map tests back to requirement/task IDs.
-->

| Test ID | Verifies | Description        | Type                     |
| ------- | -------- | ------------------ | ------------------------ |
| TC-001  | FR-001   | {what is asserted} | unit / integration / e2e |
| TC-002  | NF-001   | {what is asserted} | unit / integration / e2e |

## 5. Other

<!--
Open questions, assumptions, risks, thoughts, and observations. Open questions
MUST be resolved before status moves to `planned` and implementation begins.
-->

- {open question / note / assumption}
