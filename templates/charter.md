---
product_name: "{product_name}"
version: "0.1"
last_updated: "{last_updated}"
status: draft
---

# {product_name}

> **Charter status**: {status} · **Version**: {version} · **Last updated**: {last_updated}

<!--
APPLICATION CHARTER — product-wide context for all FlexSpec specs. Keep under
~1500–2500 tokens. Fill via /flexspec-charter. Remove guidance comments and
{placeholders} when active. Do not embed individual feature specs here.
-->

## 1. Product overview

<!--
Name, one-liner, problem, intended outcome. A reader should grasp what this
application is before reading any feature spec.
-->

**One-liner:** {one_liner}

**Problem:** {problem}

**Intended outcome:** {intended_outcome}

## 2. Vision and goals

<!--
North star and measurable success criteria for the product as a whole.
-->

**North star:** {north_star}

**Success criteria:**

{success_criteria}

## 3. Users and stakeholders

<!--
Primary personas and jobs-to-be-done. Use a table when helpful.
-->

| Persona | Role | Primary needs |
| --- | --- | --- |
| {persona_1} | {role_1} | {needs_1} |

**Jobs to be done:**

{jobs_to_be_done}

## 4. Capabilities

<!--
High-level capability map (domains/features at product level). Not spec IDs.
Update when new product areas emerge from shipped specs.
-->

{capabilities}

## 5. Technical context

<!--
Stack, deployment, integrations, hard constraints agents must respect.
-->

{technical_context}

## 6. Architecture

<!--
Optional system-context view. Use mermaid when it clarifies boundaries.
-->

{architecture_description}

```mermaid
{architecture_diagram}
```

**Boundaries:** {architecture_boundaries}

## 7. Standards and conventions

<!--
Testing, security, naming, patterns — defaults for all specs and implementations.
-->

{standards_and_conventions}

## 8. Product boundaries

<!--
Global non-goals (distinct from per-spec out-of-scope). What this product will
not do or defer intentionally.
-->

{product_boundaries}

## 9. Domain glossary

<!--
Terms and definitions used across the product.
-->

| Term | Definition |
| --- | --- |
| {term_1} | {definition_1} |

## 10. Assumptions and open questions

<!--
Blocking items must be resolved before charter is active. Non-blocking notes OK.
-->

**Assumptions:**

{assumptions}

**Open questions (blocking):**

{open_questions_blocking}

**Open questions (non-blocking):**

{open_questions_nonblocking}

## 11. Revision history

<!--
Append a row after each /flexspec-charter session or confirmed charter update
from a spec. Include originating spec slug when applicable.
-->

| Date | Summary | Source |
| --- | --- | --- |
| {revision_date} | {revision_summary} | {revision_source} |
