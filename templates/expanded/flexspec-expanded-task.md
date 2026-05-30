---
id: 'T-000'
name: ''
parent_spec: '../README.md'
status: [todo,in_progress,in_review,blocked,done]
satisfies: []        # requirement IDs, e.g. [FR-001, NF-001]
depends_on: []       # task IDs that must complete first, e.g. [T-001]
verified_by: []      # test IDs from the spec, e.g. [TC-001]
---

<!--
Keep this file under ~1000 tokens (~130 lines). Be terse — link to the parent
spec instead of restating it. If the task can't fit, split it into smaller tasks.
-->

# {id}: {task title}

> **Parent spec**: [{spec name}](../README.md) · **Status**: {status}
> **Satisfies**: {FR/NF ids} · **Depends on**: {task ids} · **Verified by**: {TC ids}

## Objective

<!--
One or two sentences: exactly what this task accomplishes and the requirement(s)
it satisfies. The reader should know "done" means after this line.
-->

{objective}

## Context

<!--
Everything the agent needs to start WITHOUT reading the rest of the codebase.
Summarize the relevant existing behavior, patterns to follow, and constraints.
Link to the parent spec sections for the bigger picture, but make this task
self-contained enough to avoid context rot.
-->

{context}

### Files In Scope

<!-- Only the files this task reads or changes. Keep it tight. -->

| File | Action | Notes |
| --- | --- | --- |
| `path/to/file` | create / modify / read | What changes / why referenced |

## Implementation Steps

<!--
Ordered, concrete steps. Each step should be small and verifiable. Reference the
exact functions, types, and files. Avoid vague instructions — an LLM should be
able to follow these literally without inventing design decisions. If a decision
is unresolved, it belongs in Open Questions (and the task is not ready to start).
-->

1. {step}
2. {step}
3. {step}

## Acceptance Criteria

<!--
Checklist defining "done" for this task. Each item should be objectively
verifiable and tie back to a requirement or test where possible.
-->

- [ ] {criterion} _(FR-001)_
- [ ] {criterion}
- [ ] All tests in "Testing" below pass.

## Testing

<!--
The tests that prove this task is complete, mapped to the spec's Testing
Criteria (TC ids). State what to write, where, and how to run them. If the work
cannot be tested as planned, stop and rework the parent spec's implementation
plan.
-->

| Test ID | Type | What it asserts | Location |
| --- | --- | --- | --- |
| TC-001 | unit / integration / e2e | {assertion} | `path/to/test` |

Run: `{command to run these tests}`

## Out of Scope

<!--
Explicitly list nearby work this task must NOT do, to prevent scope creep and
drift into other tasks.
-->

- {thing not to do here}

## Open Questions

<!--
Anything unresolved blocking this task. MUST be empty before the task moves to
`in_progress`. Ask the user to resolve these before starting.
-->

- {blocking question, if any}

## References

<!-- Links to parent spec sections, related tasks, docs, or prior art. -->

- Parent spec: [`../README.md`](../README.md)
- Related tasks: {T-00X}
