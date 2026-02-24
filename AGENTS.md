# Repository Guidelines


## Project Structure & Module Organization
- `cmd/` holds the Cobra CLI entrypoint and subcommands (e.g., `main.go`, `root.go`).
- `internal/` contains app wiring, workflows, domain models, adapters, and utilities.
- `scripts/` provides build/release helpers; `.github/workflows/` hosts release automation.
- `example-config.toml`, `README.md`, and `INSTALL.md` document usage and configuration.

## Build, Test
- `just build` builds the current platform binary.
- `just test` runs all Go tests.

## Coding Style & Naming Conventions
- Use Go standard formatting (`gofmt`) and idiomatic Go style.
- Indentation: tabs in Go source, 2 spaces in Markdown/TOML.
- File naming: `*_test.go` for tests, `snake_case` for scripts.
- Package naming: lowercase, no underscores; exported identifiers in PascalCase.
- Use internal/ui/ packages for ui elements
- Use https://github.com/charmbracelet/ packages as much as possible for input and output
- Use bubble tea tui inline options. Do not create a full screen tui

## Testing Guidelines
- Unit tests live alongside packages under `internal/` (e.g., `internal/template/processor_test.go`).
- Property-based tests use `gopter` where appropriate.
- Run tests with `just test

## Agent-Specific Instructions
- Use just over just see justfile
- When adding new CLI flags or templates, update `README.md` and `example-config.toml`.
- Prefer Bubble Tea/Bubbles for interactive UI elements.
