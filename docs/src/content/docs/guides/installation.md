---
title: Installation
description: How to install go-cli-template
---

## From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/go-cli-template.git
cd go-cli-template

# Build and install
just build
just install
```

## Using Go Install

```bash
go install github.com/yourusername/go-cli-template/cmd/go-cli-template@latest
```

## Homebrew (macOS/Linux)

```bash
brew tap yourusername/tap
brew install go-cli-template
```

## AUR (Arch Linux)

```bash
yay -S go-cli-template
```

## Verify Installation

```bash
go-cli-template --version
```

## Shell Completion

### Bash

```bash
go-cli-template completion bash > /etc/bash_completion.d/go-cli-template
```

### Zsh

```bash
go-cli-template completion zsh > "${fpath[1]}/_go-cli-template"
```

### Fish

```bash
go-cli-template completion fish > ~/.config/fish/completions/go-cli-template.fish
```

### PowerShell

```powershell
go-cli-template completion powershell | Out-String | Invoke-Expression
```
