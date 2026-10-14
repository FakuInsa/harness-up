# HarnessUp (`harness-up`)

The deterministic engineering harness generator for AI coding agents.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey)](https://github.com/FakuInsa/harness-up)

---

## Description

HarnessUp is a command-line tool (CLI) written in Go that scaffolds and configures engineering harnesses for AI coding agents (such as Google Antigravity, Cursor, Claude Code, Windsurf, OpenCode, and others).

An AI coding agent without a harness tends to hallucinate dependencies, bypass architecture patterns, or lose context across sessions. HarnessUp spins up a deterministic, project-tailored environment in seconds with architecture rules, technical conventions, and local persistent memory to keep the agent focused on any repository.

---

## Features

* **Automatic Environment Detection:** Scans repository files to detect language (C#, TypeScript, Python, Go) and existing test runners.
* **Interactive Terminal Wizard (TUI):** A keyboard-driven terminal interface (built with Bubbletea) that calibrates the harness to your project needs.
* **Context Rot Prevention:** Mitigates token degradation through lazy-loaded skills and local SQLite memory with full-text search (FTS5).
* **Multi-Agent Freedom (Zero Lock-in):** Generates configuration in native formats (.agent/rules, .cursor/rules, CLAUDE.md, .windsurfrules).
* **Idempotent Synchronization:** The sync command updates guidelines without overwriting your manual custom rules.
* **Single Standalone Binary:** Compiles into a lightweight native executable with zero external runtime dependencies and instant startup.

---

## Wizard Flow and Extensible Modules

The main command starts an interactive questionnaire:

```bash
harnessup init
```

### Extensible Decision Tree (Built to Expand)

The wizard engine is decoupled and modular, designed to continuously incorporate new languages, frameworks, and architecture patterns over time:

1. **Target Agents:** Select one or multiple AI clients (Antigravity, Cursor, Claude Code, Windsurf, OpenCode).
2. **Architecture Patterns:** Initial presets including Vertical Slice, Clean Architecture, or Hexagonal (extensible to modular monoliths and custom patterns).
3. **Workflows and Verification:** Configure engineering practices (test-driven development, atomic work-unit commits, pre-commit checks).
4. **Style and Conventions:** Naming rules (PascalCase, camelCase, snake_case), immutability, and boundary constraints.
5. **Persistent Memory Engine:** Choose between local SQLite with FTS5 or structured Markdown files.

---

## Usage

### Installation

```bash
# Using Go
go install github.com/FakuInsa/harness-up/cmd/harnessup@latest

# Or download the precompiled binary for your operating system from the Releases page
```

### Main Commands

```bash
# Interactively initialize the harness in the current project
harnessup init

# Synchronize rules and guidelines based on harnessup.yaml
harnessup sync

# Verify harness health and consistency
harnessup doctor

# Non-interactive execution for CI/CD or automation scripts
harnessup init --agent=antigravity --arch=vertical-slice --non-interactive
```

---

## Contributing

Contributions are welcome to expand the harness catalog, AI client adapters, and language templates.

1. Open an issue explaining the proposed feature, bug report, or new adapter.
2. Fork the repository and create a branch for your change.
3. Submit a Pull Request clearly describing the behavior introduced.

---

## License

Distributed under the [MIT](LICENSE) License. See the file for details.
