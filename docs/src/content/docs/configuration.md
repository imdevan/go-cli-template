---
title: Configuration
description: Configuration options for go-cli-template
---


Configuration file location: `$XDG_CONFIG_HOME/go-cli-template/config.toml`

## Configuration Options

The following options can be set in your configuration file:

### General Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `editor` | string | `nvim` | Editor to use for editing bookmarks and config files |
| `interactive_default` | bool | `false` | Start in interactive mode by default when no arguments are provided |

### Display Settings

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `list_spacing` | string | `space` | List item spacing. Options: `compact` (title only), `tight` (title + description, no margin), `space` (default, with spacing) |

### Colors

Colors support named, numeric, or hex values (e.g., `7`, `13`, `"#ff8800"`).

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `headings` | string | `15` | Color for headings |
| `primary` | string | `02` | Primary color |
| `secondary` | string | `06` | Secondary color |
| `text` | string | `07` | Text color |
| `text_highlight` | string | `06` | Highlighted text color |
| `description_highlight` | string | `05` | Highlighted description color |
| `tags` | string | `13` | Tags color |
| `flags` | string | `12` | Flags color |
| `muted` | string | `08` | Muted text color |
| `border` | string | `08` | Border color |

## Example Configuration

```toml
# General
editor = "nvim"

# CLI behavior
interactive_default = true

# UI
# list_spacing options: compact (title only), tight (title + description, no margin), space (default, with spacing)
list_spacing = "space"

# Colors
# Colors support named, numeric, or hex values (ex: 7, 13, "#ff8800").
headings = "15"
primary = "02"
secondary = "06"
text = "07"
text_highlight = "06"
description_highlight = "05"
tags = "13"
flags = "12"
muted = "08"
border = "08"
```

## Initializing Configuration

To create a new configuration file with default values:

```bash
bookmark config init
```

To overwrite an existing configuration:

```bash
bookmark config init --force
```

To create and immediately open in your editor:

```bash
bookmark config init --editor
```

## Editing Configuration

To edit your configuration file:

```bash
bookmark config
```

This will open the config file in your configured editor.
