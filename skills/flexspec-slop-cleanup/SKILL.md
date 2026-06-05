---
name: flexspec-slop-cleanup
description: Detect and optionally fix AI slop in code — semantic duplication, complexity inflation, test-implementation coupling, architectural drift, and useless comments. Reviews the current git diff by default, or specific files when paths are passed. Use when the user says "check for slop", "find AI slop", "slop review", "is this slop", "clean up the AI code", or invokes /flexspec-slop-cleanup. Pass --fix to apply changes; default is report-only.
---

# AI Slop Review

Catch the failure modes AI-generated code is most prone to: it compiles, passes tests, and looks polished — but duplicates logic, over-engineers simple problems, mirrors implementation in tests, drifts from existing patterns, and narrates itself with comments. Patterns are drawn from Larridin's AI Slop analysis and the repo's `AGENTS.md` rules.

## Inputs

- **No args** → review the current git diff (staged + unstaged vs. `main`).
- **Path args** (`/flexspec-slop-cleanup src/foo.ts src/bar.ts`) → review those files.
- **`--fix`** → apply changes after review. Without it, report only.

## Workflow

1. **Gather scope.** Run `git diff main...HEAD` plus `git diff` (unstaged) for diff mode, or read passed files. Use `git ls-files` + grep to discover surrounding patterns the new code should match.
2. **Scan for each pattern below.** For every finding, record: file:line, which pattern, the specific evidence, and the suggested fix.
3. **Report.** Group findings by pattern. For each, include a one-line "why this is slop" tied to the evidence, not a generic restatement.
4. **If `--fix`** → apply the fixes via Edit. Skip findings where the fix requires judgment the user must make (e.g. "this should arguably not exist at all").

## Patterns to detect

**1. Semantic duplication.** New code reimplements logic that already exists. Check: utility functions, formatters, validation, API client helpers, hooks. Grep the codebase for similar function names, signatures, and string literals before flagging. Fix: import the existing helper; delete the duplicate.

**2. Complexity inflation.** A simple requirement met with elaborate machinery. Red flags: a class wrapping one function, options objects with one consumer, abstract base classes with one subclass, try/catch around code that cannot throw, defensive validation on internal callers, feature flags or compat shims for code that has no prior version. Fix: inline, delete the abstraction, trust the caller.

**3. Test-implementation coupling.** Tests assert _how_ code works, not _what_ it does. Red flags: tests that mock the unit under test, assertions on internal call order/counts when behavior is what matters, snapshot tests of implementation detail, tests that would still pass if the function returned wrong values, 1:1 mirroring of branches in the impl. Fix: rewrite to assert observable behavior (return value, side effect on a real dependency, rendered output).

**4. Pattern mimicry / architectural drift.** New code uses a pattern foreign to the surrounding module. Examples: a 4th data-fetching style when 3 exist, a new state container when Zustand is already chosen, manual fetch when `axiosClient` is the convention, relative imports when `@/` is enforced, a new module shape when `src/modules/*/index.ts` has a contract. Cross-reference `AGENTS.md` and neighboring files. Fix: rewrite to match the existing pattern.

**5. Useless comments.** Comments that restate code, reference the current task ("added for X flow"), narrate what is being done, or wrap obvious behavior. Allowed: comments explaining a non-obvious _why_ (hidden constraint, workaround, surprising invariant), and linter-required doc comments (e.g. Go exported identifiers). Fix: delete the comment. Keep the rare ones that pay rent.

## Reporting format

```
AI slop review — <N> findings across <M> files

[duplication] src/lib/format.ts:42
  `formatCurrency` reimplements src/utils/money.ts:formatMoney (same logic, same locale handling)
  → import formatMoney; delete formatCurrency

[complexity] src/api/users.ts:18
  `UserFetcher` class wraps a single fetch call with no shared state
  → replace with a plain async function

[test-coupling] src/components/__tests__/Button.test.tsx:55
  asserts internal useState call order; would pass even if onClick never fires
  → assert that onClick handler runs on user click

[drift] src/modules/foo/api.ts:1
  uses fetch + manual auth headers; project convention is axiosClient (src/lib/axios-client.ts)
  → switch to axiosClient

[comments] src/hooks/useThing.ts:12
  `// set the loading state to true` restates the next line
  → delete
```

End with: `Run with --fix to apply N of M findings automatically. K findings need human judgment and were not auto-fixable.`

## Calibration

- Be specific about _what_ and _where_. "This file has slop" is useless; "line 42 duplicates X at path Y" is actionable.
- If a finding is genuinely ambiguous (might be intentional), say so — don't pad the report with low-confidence flags.
- Three real findings beats ten speculative ones.
- Don't flag comments the linter or language requires (Go package/export docs, JSDoc that powers tooling, license headers).
