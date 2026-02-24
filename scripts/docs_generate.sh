#!/usr/bin/env bash
set -euo pipefail

# Generate API documentation from Go packages using gomarkdoc

PACKAGE_FILE="internal/package/package.toml"
DOCS_API_DIR="docs/src/content/docs/api"
ASTRO_CONFIG="docs/astro.config.mjs"

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

echo "🔧 Updating Astro config..."
# Update astro.config.mjs with values from package.toml
if [ -f "$ASTRO_CONFIG" ]; then
  # Create a temporary config with updated values
  awk -v site="$DOCS_SITE" -v base="$DOCS_BASE" -v title="$PROJECT_NAME" -v desc="$DESCRIPTION" -v repo="$REPOSITORY" '
    /^export default defineConfig\({/ {
      print $0
      print "  site: \"" site "\","
      print "  base: \"" base "\","
      next
    }
    /^  site:/ || /^  base:/ { next }
    /title:/ { sub(/title: .*,/, "title: \"" title "\","); print; next }
    /description:/ { sub(/description: .*,/, "description: \"" desc "\","); print; next }
    /github:/ && repo != "" { sub(/github: .*,/, "github: \"" repo "\","); print; next }
    { print }
  ' "$ASTRO_CONFIG" > "${ASTRO_CONFIG}.tmp" && mv "${ASTRO_CONFIG}.tmp" "$ASTRO_CONFIG"
  echo "  ✓ Updated site, base, title, and description"
fi

echo "🔧 Checking for gomarkdoc..."
if ! command -v gomarkdoc &> /dev/null; then
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
    } > "$DOCS_API_DIR/${pkg_name}.md"
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
    } > "$DOCS_API_DIR/adapters/${adapter_name}.md"
done

echo "✅ API documentation generated successfully!"
echo "📁 Output: $DOCS_API_DIR"
