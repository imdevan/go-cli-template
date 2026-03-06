---
title: go-cli-template
description: A generic CLI tool template built with Go, Cobra, and Bubble Tea. This template provides a foundation for building interactive command-line applications with a clean architecture and modern UI components.
---

A generic CLI tool template built with Go, Cobra, and Bubble Tea. This template provides a foundation for building interactive command-line applications with a clean architecture and modern UI components.

## Usage

```bash
go-cli-template [alias]
go-cli-template [command]
```

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `-c, --config` | string | config file path |
| `-v, --version` | bool | print version information |

## Available Commands

- [`completion`](/commands/completion) - Generate shell completion scripts
- [`config`](/commands/config) - View or edit configuration
- [`config init`](/commands/config-init) - Generate a default config file

## Source

See [root.go](https://github.com/imdevan/go-cli-template//blob/main/cmd/go-cli-template/root.go) for implementation details.
