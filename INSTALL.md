# Installation

This template is currently setup to build and deploy to homebrew and AUR. 

Because that is what I use so that that is what I have capacity to test at the moment. 

## Homebrew
```bash
brew install imdevan/go-cli-template/go-cli-template
```

## Arch (AUR)
```bash
yay -S go-cli-template
```

## GitHub Release

Download the latest binary for your platform from the [releases page](https://github.com/imdevan/go-cli-template/releases).

```bash
# Linux (amd64)
curl -L https://github.com/imdevan/go-cli-template/releases/latest/download/go-cli-template-linux-amd64.tar.gz | tar -xz
sudo mv go-cli-template-linux-amd64 /usr/local/bin/go-cli-template
```

```bash
# macOS (Apple Silicon)
curl -L https://github.com/imdevan/go-cli-template/releases/latest/download/go-cli-template-darwin-arm64.tar.gz | tar -xz
sudo mv go-cli-template-darwin-arm64 /usr/local/bin/go-cli-template
```

```bash
# macOS (Intel)
curl -L https://github.com/imdevan/go-cli-template/releases/latest/download/go-cli-template-darwin-amd64.tar.gz | tar -xz
sudo mv go-cli-template-darwin-amd64 /usr/local/bin/go-cli-template
```

## Manual
```bash
gh repo clone imdevan/go-cli-template
cd go-cli-template
just build
sudo just install
```
