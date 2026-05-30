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

## Usage

```bash
flexspec --help
```

More commands and workflows will be documented here as they are added.

## License

Copyright © 2026 Josh Kyte
