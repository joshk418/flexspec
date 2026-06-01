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

Contributors need Node.js 22+ only to build the embedded UI under `ui/`. End users of the released binary do not need Node.

## Installing FlexSpec Skills

Agent skills live under [`skills/`](skills/):

| Skill               | Path                                                           | Command             |
| ------------------- | -------------------------------------------------------------- | ------------------- |
| Spec lifecycle      | [`skills/flexspec/`](skills/flexspec/SKILL.md)                 | `/flexspec`         |
| Application charter | [`skills/flexspec-charter/`](skills/flexspec-charter/SKILL.md) | `/flexspec-charter` |

Install both into your coding agent with [`npx skills`](https://github.com/vercel-labs/skills):

```bash
npx skills add joshk418/flexspec
```

This auto-detects your installed agents (Cursor, Claude Code, Codex, and 50+ others) and installs all skills in the repo. Useful variants:

```bash
# Install for all projects instead of just the current one
npx skills add joshk418/flexspec --global

# Target a specific agent
npx skills add joshk418/flexspec --agent cursor

# Preview the skills in the repo without installing
npx skills add joshk418/flexspec --list
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
| `flexspec list`                                     | Compact table of spec directory identifiers, statuses, and task counts |
| `flexspec list --json`                              | Same data as JSON (scripts, CI)                                     |
| `flexspec validate`                                 | Check config, charter, templates, and specs for structural problems |
| `flexspec ui`                                       | Start local management UI (default http://127.0.0.1:3000)           |
| `flexspec status set <spec> --status <s>`           | Update spec or task frontmatter status (`--task` for task files)      |

```bash
flexspec --help
flexspec init
flexspec config
flexspec config set spec_template expanded
flexspec new my-feature --template simple
flexspec list
flexspec validate
flexspec ui --no-open
flexspec status set 001-my-feature --status in_progress
```

`flexspec ui` flags: `--port`, `--host` (default `127.0.0.1`), `--open` / `--no-open`.

`flexspec validate` prints findings as `severity`, `path`, `rule`, `message` (tab-separated), then a summary. It exits **0** when there are no errors and **1** when any error-severity finding exists (warnings alone do not fail). If config is missing, it reports that and skips deeper checks.

## License

MIT. See [LICENSE](LICENSE).
