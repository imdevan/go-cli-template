---
title: Configuration
description: Configure go-cli-template for your needs
---

## Configuration File

go-cli-template uses TOML for configuration. The config file is located at:

```
$XDG_CONFIG_HOME/go-cli-template/config.toml
```

On most systems, this resolves to:
- Linux: `~/.config/go-cli-template/config.toml`
- macOS: `~/Library/Application Support/go-cli-template/config.toml`

## Initialize Configuration

Generate a default configuration file:

```bash
go-cli-template config init
```

## View Configuration

View your current configuration:

```bash
go-cli-template config
```

## Edit Configuration

Open the configuration file in your default editor:

```bash
go-cli-template config --edit
```

## Configuration Options

See the [example-config.toml](https://github.com/yourusername/go-cli-template/blob/main/example-config.toml) file in the repository for all available options and their descriptions.

## Environment Variables

You can override configuration with environment variables:

```bash
export GO_CLI_TEMPLATE_CONFIG=/path/to/custom/config.toml
```
