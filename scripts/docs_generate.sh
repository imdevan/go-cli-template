#!/usr/bin/env bash
set -euo pipefail

# Generate API documentation from Go packages using gomarkdoc
# Usage: ./docs_generate.sh [--dev]
#   --dev: Use '/' as base for local development

PACKAGE_FILE="internal/package/package.toml"
DOCS_API_DIR="docs/src/content/docs/api"
DOCS_CONFIG="docs/config.mjs"

# Parse package.toml
parse_toml() {
  local key=$1
  grep "^$key = " "$PACKAGE_FILE" | sed 's/^[^=]*= *"\(.*\)"$/\1/'
}

echo "📦 Reading package metadata..."
PROJECT_NAME=$(parse_toml "name")
DESCRIPTION=$(parse_toml "description")
DOCS_SITE=$(parse_toml "docs_site")
DOCS_BASE=$(parse_toml "docs_base")
REPOSITORY=$(parse_toml "repository")

# Use defaults if repository is empty
if [ -z "$REPOSITORY" ]; then
  REPOSITORY="https://github.com/yourusername/${PROJECT_NAME}"
fi

echo "🔧 Updating docs config..."

# Update docs/config.mjs with values from package.toml
if [ -f "$DOCS_CONFIG" ]; then
  cat >"$DOCS_CONFIG" <<EOF
const stage = process.env.NODE_ENV || "dev"
const isProduction = stage === "production"
const basePath = isProduction ? "$DOCS_BASE" : ""

export default {
  url: isProduction ? "$DOCS_SITE$DOCS_BASE" : "http://localhost:4321",
  basePath,
  github: "$REPOSITORY",
  githubDocs: "$REPOSITORY",
  title: "$PROJECT_NAME",
  description: "$DESCRIPTION",
}
EOF
  echo "  ✓ Updated config.mjs with package metadata"
fi

echo "🔧 Checking for gomarkdoc..."
if ! command -v gomarkdoc &>/dev/null; then
  echo "📦 Installing gomarkdoc..."
  go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
fi

echo "🧹 Cleaning old API docs..."
rm -rf "$DOCS_API_DIR"
mkdir -p "$DOCS_API_DIR"

echo "📝 Generating API documentation..."

# Generate docs for each internal package
for pkg in internal/*/; do
  pkg_name=$(basename "$pkg")

  # Skip test utilities and adapters subdirectories
  if [[ "$pkg_name" == "testutil" ]]; then
    continue
  fi

  echo "  - Processing $pkg_name..."

  # Generate to temp file first
  gomarkdoc --output "/tmp/${pkg_name}.md" "./$pkg" 2>/dev/null || {
    echo "    ⚠️  No exported symbols in $pkg_name"
    continue
  }

  # Add frontmatter and content (skip HTML comment and blank lines)
  {
    echo "---"
    echo "title: ${pkg_name}"
    echo "description: API documentation for the ${pkg_name} package"
    echo "---"
    echo ""
    # Skip HTML comment and any frontmatter that gomarkdoc added
    sed -n '/^# /,$p' "/tmp/${pkg_name}.md"
  } >"$DOCS_API_DIR/${pkg_name}.md"
done

# Generate docs for adapters
echo "  - Processing adapters..."
mkdir -p "$DOCS_API_DIR/adapters"

for adapter in internal/adapters/*/; do
  adapter_name=$(basename "$adapter")
  echo "    - Processing adapters/$adapter_name..."

  # Generate to temp file first
  gomarkdoc --output "/tmp/adapter_${adapter_name}.md" "./$adapter" 2>/dev/null || {
    echo "      ⚠️  No exported symbols in $adapter_name"
    continue
  }

  # Add frontmatter and content (skip HTML comment and blank lines)
  {
    echo "---"
    echo "title: adapters/${adapter_name}"
    echo "description: API documentation for the ${adapter_name} adapter"
    echo "---"
    echo ""
    # Skip HTML comment and any frontmatter that gomarkdoc added
    sed -n '/^# /,$p' "/tmp/adapter_${adapter_name}.md"
  } >"$DOCS_API_DIR/adapters/${adapter_name}.md"
done

echo "✅ API documentation generated successfully!"
echo "📁 Output: $DOCS_API_DIR"
