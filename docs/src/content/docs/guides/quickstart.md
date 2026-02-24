---
title: Quick Start
description: Get started with go-cli-template in minutes
---

## Basic Usage

```bash
# Run the main command
go-cli-template

# View configuration
go-cli-template config

# Initialize configuration
go-cli-template config init
```

## Configuration

Configuration file location: `$XDG_CONFIG_HOME/go-cli-template/config.toml`

Example configuration:

```toml
# See example-config.toml for all available options
[general]
theme = "default"
```

## Development

```bash
# Build the project
just build

# Run tests
just test

# Watch for changes
just watch

# Build for multiple platforms
just cross-platform
```

## Customization

This template is designed to be customized:

1. Edit `package.toml` with your project details
2. Run `just sync` to sync changes across all files
3. Review changes with `git diff`
4. Build and test: `just build && just test`

The `package.toml` file is the single source of truth for project metadata.
