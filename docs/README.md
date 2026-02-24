# Documentation

This directory contains the Astro Starlight documentation site for go-cli-template.

## Setup

```bash
# Install dependencies
just docs-init

# Generate API docs and start dev server
just docs-dev
```

## Commands

- `just docs-init` - Install npm dependencies
- `just docs-generate` - Generate API docs from Go packages
- `just docs-dev` - Start development server with hot reload
- `just docs-build` - Build production site
- `just docs-preview` - Preview production build
- `just docs-clean` - Clean build artifacts

## Structure

```
docs/
├── src/
│   ├── content/
│   │   └── docs/
│   │       ├── index.mdx          # Homepage
│   │       ├── guides/            # Manual guides
│   │       └── api/               # Auto-generated API docs
│   ├── assets/                    # Images and static assets
│   └── styles/                    # Custom CSS
├── astro.config.mjs               # Astro configuration
└── package.json                   # Dependencies
```

## Writing Documentation

### Manual Pages

Create `.md` or `.mdx` files in `src/content/docs/guides/`:

```markdown
---
title: Your Guide Title
description: A brief description
---

Your content here...
```

### API Documentation

API docs are auto-generated from Go package comments using gomarkdoc. To update:

1. Add/update Go doc comments in your code
2. Run `just docs-generate`

## Deployment

Documentation is automatically deployed to GitHub Pages on push to main via `.github/workflows/docs.yml`.

To enable:
1. Go to repository Settings → Pages
2. Set Source to "GitHub Actions"
3. Push to main branch
