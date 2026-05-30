# Flexspec

Flexspec is a spec-driven development CLI for generating and tracking feature specifications. Use markdown files and templates, or connect adapters for issue trackers like Jira, Shortcut, and GitHub Issues.

## Why Flexspec

Spec-driven development helps teams and AI coding agents agree on what to build before writing code. Flexspec focuses on keeping specs clear, structured, and easy to maintain.

Compared with tools such as [Spec Kit](https://github.com/github/spec-kit), Spec Kitty, and [OpenSpec](https://github.com/Fission-AI/OpenSpec), Flexspec gives you two modes:

- **Simple specs** — Create easy-to-read markdown files from provided templates for small, focused features.
- **Linked specs** — Build larger features as multi-file specifications with explicit links between related docs.

That range lets you start lightweight and scale up when a feature needs more detail. Well-defined specs help agents follow requirements closely and reduce LLM drift away from the intended outcome.

## Features

- Generate and manage spec files from the command line
- Markdown-first workflow with reusable templates
- Adapter support for external systems (Jira, Shortcut, GitHub Issues, and more)
- Single-file specs for quick features
- Multi-file, linked specifications for complex features

## Installation

```bash
go install github.com/joshk418/flexspec@latest
```

Or clone and build locally:

```bash
git clone https://github.com/joshk418/flexspec.git
cd flexspec
go build -o flexspec .
```

## Installing Flexspec Skills

Agent skills live under [`skills/`](skills/):

| Skill | Path | Command |
| --- | --- | --- |
| Spec lifecycle | [`skills/flexspec/`](skills/flexspec/SKILL.md) | `/flexspec` |
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

Once installed, reload your agent before invoking `/flexspec` or `/flexspec-charter`.

## Usage

```bash
flexspec --help
```

More commands and workflows will be documented here as they are added.

## License

Copyright © 2026 Josh Kyte
