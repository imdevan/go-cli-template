# go-cli-template

<img width="480" height="270" alt="screenshot-2026-02-23_16-30-13" src="https://github.com/user-attachments/assets/65386b56-f06f-47be-9063-5c947b30dc51" />

A generic CLI tool template built with Go, Cobra, and Bubble Tea. This template provides a foundation for building interactive command-line applications with a clean architecture and modern UI components.

## Features

- Interactive list with filtering
- [Built on go Cobra](https://github.com/spf13/cobra)
- Configuration management with TOML
- Styles, build scripts, and tests to get you started.
- [Inline Bubble Tea TUI components](https://github.com/charmbracelet/bubbletea)
- Homebrew and aur package management with TOML too!
- Automatic documentation with [gomarkdoc](https://github.com/princjef/gomarkdoc) and [astro starlight](https://starlight.astro.build/)
- XDG Base Directory support
- Shell completion (bash, zsh, fish, powershell)

## Quick start

```bash
just build-run
```

## Commands

```bash
go-cli-template                 # Root command (placerholder shows folder content)
go-cli-template config          # View or edit configuration
go-cli-template config init     # Generate default config file
go-cli-template completion      # Generate shell completion scripts
```

## Development

```bash
just sync            # Sync project from package.toml
just build           # Build the binary
just build-run       # Build and run the binary
just dev-build       # Build with debug symbols
just test            # Run tests
just test-verbose    # Run tests with verbose output
just watch           # Watch for changes and rebuild
just cross-platform  # Build for multiple platforms
just install         # Install to /usr/local/bin
just clean           # Remove build artifacts
```

## Configuration

Configuration file location: `$XDG_CONFIG_HOME/go-cli-template/config.toml`

See `example-config.toml` for available configuration options.

## Installation

See `INSTALL.md` for installation options.

## Customization

This template is designed to be customized for your specific CLI tool needs:

1. Edit `package.toml` with your project details (name, module, description, etc.)
2. Run `just sync` to sync changes across all files
3. Review changes with `git diff`
4. Build and test: `just build && just test`

The `package.toml` file is the single source of truth for project metadata. The sync script will update:
- Go module name in `go.mod` and all import paths
- Binary name in justfile and build scripts
- Config paths in `internal/utils/paths.go`
- Completion examples
- README description
- Version in root.go

## Architecture

- `cmd/`              - CLI entrypoint and commands
- `internal/config`   - Configuration management
- `internal/domain`   - Domain models
- `internal/ui`       - Bubble Tea UI components
- `internal/utils`    - Utility functions
- `internal/adapters` - External service adapters (editor, clipboard)

# Thank you!

This project was made by deconstructing a another cli project of mine [Prompter](http://devan.gg/prompter-cli/). Check it out if you like fiddling with coding agents and want a more vim centric way of managing your prompting!


