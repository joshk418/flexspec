---
name: flexspec-glossary-discovery
description: >
  Scan a FlexSpec project for unknown project-specific terms, ask the user for
  exact meanings when unclear, and record confirmed definitions through the
  CLI. Trigger with "flexspec-glossary-discovery", "flexspec-glossary", or
  "discover glossary terms". Also used by /flexspec-charter during charter work.
triggers:
  - flexspec-glossary-discovery
  - flexspec-glossary
  - discover glossary terms
  - scan project glossary
---

# Glossary Discovery (`flexspec-glossary-discovery`)

A focused repository scan for project-specific vocabulary that should be recorded
in `.flexspec/glossary.yaml`.

This skill can run standalone, and `/flexspec-charter` should invoke or follow
this workflow during charter creation, full refresh, or terminology-heavy charter
updates so the charter and glossary are built together.

## Workflow

1. **Read known terms**
   - Run `flexspec glossary list --json` to load existing glossary entries.
   - Build a case-insensitive set of known terms and aliases.

2. **Scan candidate sources**
   - Run `flexspec glossary scan --json` to get candidate terms from specs, charter, code identifiers, and docs. This is cross-platform (no ripgrep dependency).
   - If the CLI subcommand is unavailable (older flexspec), fall back to `rg -o '[A-Z][a-zA-Z]{3,}'` when ripgrep is present, or `grep -oE '[A-Z][a-zA-Z]{3,}'` / `findstr` on systems without ripgrep.
   - The scan includes config key names, package names, and domain identifiers, and excludes common language keywords and standard library names.

3. **Rank and filter**
   - Count frequencies; rank repeated or domain-like terms higher.
   - Exclude terms already in the glossary (exact or alias match).
   - Exclude vendor/library names unless they appear project-specific.
   - Keep top candidates (cap at ~20 to stay token-efficient).

4. **Interview loop**
   - For each unclear candidate, ask: *"What does `<term>` mean in this project?"*
   - Offer: definition text, category, aliases, source.
   - If the user skips, record as skipped.
   - If the meaning is clear from context, skip the interview and record directly.

5. **Persist confirmed terms**
   - Use `flexspec glossary add <term> --definition <text> [--alias <a>] [--category <c>] [--source discovery]`.
   - Never manually edit `.flexspec/glossary.yaml`.

6. **Report**
   - Summarize: added terms, skipped terms, still-ambiguous terms.
   - Suggest re-running discovery after major feature additions.

## Filtering rules

- The `flexspec glossary scan` subcommand applies these exclusions natively; the fallback `rg`/`grep` path requires manual exclusion.
- Exclude common words: `true`, `false`, `null`, `error`, `string`, `int`, `return`, `func`, `function`, `class`, `interface`, `struct`, `package`, `import`, `var`, `const`, `let`, `this`, `self`, `static`, `public`, `private`.
- Exclude standard library and framework names unless the project extends them.
- Exclude file extensions and generic abbreviations.

## Interview style

- Ask one term at a time when few candidates remain.
- Use grouped multiple-choice when the runtime supports it and many candidates are unclear.
- Keep definitions under 200 characters when possible.
- Always include `--source discovery` for terms added by this skill.

## Out of scope

- Embeddings or network searches.
- Writing project documentation outside FlexSpec metadata.
- Modifying source code.
