# FlexSpec

[![CI](https://github.com/joshk418/flexspec/actions/workflows/ci.yml/badge.svg)](https://github.com/joshk418/flexspec/actions/workflows/ci.yml)

FlexSpec is a spec-driven development CLI for generating and tracking feature specifications. Use markdown files and templates, or connect adapters for issue trackers like Jira, Shortcut, and GitHub Issues.

## Why FlexSpec

Spec-driven development helps teams and AI coding agents agree on what to build before writing code. FlexSpec focuses on keeping specs clear, structured, and easy to maintain.

Compared with tools such as [Spec Kit](https://github.com/github/spec-kit), Spec Kitty, and [OpenSpec](https://github.com/Fission-AI/OpenSpec), FlexSpec gives you two modes:

- **Simple specs** — Create easy-to-read markdown files from provided templates for small, focused features.
- **Linked specs** — Build larger features as multi-file specifications with explicit links between related docs.

That range lets you start lightweight and scale up when a feature needs more detail. Well-defined specs help agents follow requirements closely and reduce LLM drift away from the intended outcome.

## Features

- Generate and manage spec files from the command line
- Validate project structure, templates, and spec frontmatter before other commands fail
- Markdown-first workflow with reusable templates
- Adapter support for external systems (Jira, Shortcut, GitHub Issues, and more)
- Single-file specs for quick features
- Multi-file, linked specifications for complex features
- Local management UI (`flexspec ui`) — kanban/table board, spec browser, structured settings for UI prefs and `.flexspec/config.yaml`; live refresh via filesystem watch (SSE).

## Installation

```bash
go install github.com/joshk418/flexspec@latest
```

Or clone and build locally:

```bash
git clone https://github.com/joshk418/flexspec.git
cd flexspec
make build    # builds web UI then compiles flexspec
# or: make build-ui && go build -o flexspec .
```

Contributors need Node.js 22+ only to build the embedded UI under `ui/`. End users of the released binary do not need Node or Go — `flexspec update` downloads prebuilt binaries from GitHub Releases directly.

## Installing FlexSpec Skills

Agent skills live under [`skills/`](skills/):

| Skill               | Path                                                                           | Command                     |
| ------------------- | ------------------------------------------------------------------------------ | --------------------------- |
| Spec lifecycle      | [`skills/flexspec/`](skills/flexspec/SKILL.md)                                 | `/flexspec`                 |
| Application charter | [`skills/flexspec-charter/`](skills/flexspec-charter/SKILL.md)                 | `/flexspec-charter`         |
| Migration           | [`skills/flexspec-migrate/`](skills/flexspec-migrate/SKILL.md)                 | `/flexspec-migrate`         |
| Glossary discovery  | [`skills/flexspec-glossary-discovery/`](skills/flexspec-glossary-discovery/SKILL.md) | `flexspec-glossary-discovery` |
| AI slop review      | [`skills/flexspec-slop-cleanup/`](skills/flexspec-slop-cleanup/SKILL.md)       | `/flexspec-slop-cleanup`    |

Skills are embedded in the flexspec binary and installed by `flexspec update --skills` directly into each detected coding agent's skills directory:

| Agent       | Global path                      | Project path            |
| ----------- | -------------------------------- | ----------------------- |
| Claude Code | `~/.claude/skills/`              | `./.claude/skills/`     |
| Cursor      | `~/.cursor/skills/`              | `./.agents/skills/`     |
| Codex       | `~/.codex/skills/`               | `./.agents/skills/`     |
| OpenCode    | `~/.config/opencode/skills/`     | `./.agents/skills/`     |
| Cline       | `~/.agents/skills/`              | `./.agents/skills/`     |

An agent is "detected" when its config root (e.g. `~/.claude/`) exists; the `skills/` subdir is created on install if missing. Skills are version-pinned to the CLI binary — the exact skills shipped with `flexspec` v0.3.5 are the ones installed by v0.3.5's `update`.

For unsupported agents, the legacy `npx skills` path is still available (requires Node):

```bash
npx skills add joshk418/flexspec --global
# or target a specific agent
npx skills add joshk418/flexspec --agent cursor
```

### Recommended workflow

1. Run `flexspec init` in your project — creates `.flexspec/config.yaml`, `.flexspec/charter.md`, and `.flexspec/templates/`.
2. Run `/flexspec-charter` — interview to fill the application charter (vision, capabilities, constraints).
3. Run `/flexspec` per feature — specs use the charter as product context; when a spec implies charter changes, the agent prompts you to update the charter (deltas only).
4. Run `flexspec validate` after setup or when specs change — catches broken config and unreadable frontmatter (useful in CI; exit code 1 on errors).
5. Optional: run `flexspec ui` — open the local dashboard at http://127.0.0.1:3000 for board and spec visibility while agents work.

Once installed, reload your agent before invoking `/flexspec` or `/flexspec-charter`.

## Usage

From your project root (after `flexspec init`):

| Command                                             | Purpose                                                             |
| --------------------------------------------------- | ------------------------------------------------------------------- |
| `flexspec init`                                     | Create `.flexspec/`, config, charter, and templates                 |
| `flexspec new <name> --template <simple\|expanded>` | Scaffold a new spec under the configured specs directory            |
| `flexspec config`                                   | Show `.flexspec/config.yaml` settings (KEY / VALUE table)           |
| `flexspec config --json`                            | Same config as JSON (scripts, agents)                               |
| `flexspec config set <key> <value>`                 | Update one config key and print the updated table                   |
| `flexspec list`                                     | Compact table of spec directory identifiers, statuses, and task counts (`task_count` frontmatter, or computed) |
| `flexspec list --json`                              | Same data as JSON, including `task_count` (scripts, CI)             |
| `flexspec validate`                                 | Check config, charter, templates, and specs for structural problems |
| `flexspec update`                                   | Check for a newer CLI binary and update; then reinstall skills and run migrations (default: all three) |
| `flexspec update --dry-run`                         | Preview update steps without writing or executing external commands           |
| `flexspec update --check`                           | CI gate: exit 1 when migrations are pending (detect only)                     |
| `flexspec update --cli`                             | Update only the CLI binary (download + swap from GitHub Releases)             |
| `flexspec update --skills`                          | Reinstall embedded skills into each detected agent's skills dir               |
| `flexspec update --migrate`                         | Run only in-project migrations                                                |
| `flexspec ui`                                       | Start local management UI (default http://127.0.0.1:3000)                     |
| `flexspec status set <spec> --status <s>`           | Update spec or task frontmatter status (`--task` for task files)              |

```bash
flexspec --help
flexspec init
flexspec config
flexspec config set spec_template expanded
flexspec new my-feature --template simple
flexspec list
flexspec validate
flexspec update
flexspec update --dry-run
flexspec update --migrate --only status-rename
flexspec update --skills --skills-method npx
flexspec ui --no-open
flexspec status set 001-my-feature --status in_progress
```

`flexspec ui` flags: `--port`, `--host` (default `127.0.0.1`), `--open` / `--no-open`.

`flexspec validate` prints findings in a `SEVERITY / PATH / RULE / MESSAGE` table, then a summary. It exits **0** when there are no errors and **1** when any error-severity finding exists (warnings alone do not fail). If config is missing, it reports that and skips deeper checks.

`flexspec update` runs three steps in order: (1) check the latest GitHub release and, if newer, download the matching prebuilt binary, verify its SHA256 against `checksums.txt`, atomically swap the running executable, and re-exec into the new binary so the remaining steps run under the new code; (2) install the embedded skills into each detected coding agent's skills directory (falls back to `npx skills add --global` if no supported agent is detected); (3) run in-project migrations (spec statuses, `task_count` backfill, template re-sync, config keys, charter checks, glossary, type backfill) — only when inside a `.flexspec/` project. No Go toolchain or Node install is required for the binary or skills steps. Use `--cli`, `--skills`, or `--migrate` to run individual steps. `--skills-method auto|embedded|npx` controls the skills install path (default `auto`: embedded if an agent is detected, else npx). `--no-reexec` downloads and swaps the binary but does not re-exec; re-run `flexspec update --skills --migrate` manually to finish. `--force` re-downloads the binary even when already latest and overwrites differing template files on migrate.

## License

MIT. See [LICENSE](LICENSE).
