#!/usr/bin/env bash
set -euo pipefail

# Generate API documentation from Go packages using gomarkdoc
# Usage: ./docs_generate.sh [--dev]
#   --dev: Use '/' as base for local development

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_FILE="${ROOT_DIR}/internal/package/package.toml"
DOCS_API_DIR="${ROOT_DIR}/docs/src/content/docs/api"
DOCS_CONFIG="${ROOT_DIR}/docs/config.mjs"
DOCS_SIDEBAR="${ROOT_DIR}/docs/sidebar.mjs"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

echo "📦 Reading package metadata..."
PROJECT_NAME=$(parse_toml_key "$PACKAGE_FILE" "name")
CMD_DIR="${ROOT_DIR}/cmd/${PROJECT_NAME}"
DESCRIPTION=$(parse_toml_key "$PACKAGE_FILE" "description")
DOCS_SITE=$(parse_toml_key "$PACKAGE_FILE" "docs_site")
DOCS_BASE=$(parse_toml_key "$PACKAGE_FILE" "docs_base")
REPOSITORY=$(parse_toml_key "$PACKAGE_FILE" "repository")

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

export default {
  url: isProduction ? "$DOCS_SITE" : "http://localhost:4321",
  basePath:  isProduction ? "$DOCS_BASE" : "/",
  github: "$REPOSITORY",
  githubDocs: "$REPOSITORY",
  title: "$PROJECT_NAME",
  description: "$DESCRIPTION",
}
EOF
  echo "  ✓ Updated config.mjs with package metadata"
fi

echo "🔧 Generating sidebar configuration..."

# Detect commands from cmd directory
COMMANDS=""
if [ -d "$CMD_DIR" ]; then
  for cmd_file in "$CMD_DIR"/*.go; do
    # Skip test files, main.go, and root.go
    if [[ "$cmd_file" == *"_test.go" ]] || [[ "$cmd_file" == *"/main.go" ]] || [[ "$cmd_file" == *"/root.go" ]]; then
      continue
    fi

    # Extract command name from filename (e.g., config.go -> config, delete_cmd.go -> delete)
    cmd_name=$(basename "$cmd_file" .go | sed 's/_cmd$//')

    # Convert underscores to spaces for display (e.g., config_init -> config init)
    cmd_display=$(echo "$cmd_name" | sed 's/_/ /g')

    # Convert underscores to hyphens for URL (e.g., config_init -> config-init)
    cmd_url=$(echo "$cmd_name" | sed 's/_/-/g')

    COMMANDS="${COMMANDS}            { label: '${cmd_display}', link: '/commands/${cmd_url}' },
"
  done
fi

# Detect API packages
API_PACKAGES=""
if [ -d "${ROOT_DIR}/internal" ]; then
  for pkg in "${ROOT_DIR}/internal"/*/; do
    pkg_name=$(basename "$pkg")
    # Skip testutil
    if [[ "$pkg_name" == "testutil" ]]; then
      continue
    fi
    # Check if directory contains Go files (is a package, not just a folder)
    if ls "$pkg"*.go >/dev/null 2>&1; then
      API_PACKAGES="${API_PACKAGES}            { label: '${pkg_name}', link: '/api/${pkg_name}' },
"
    fi
  done
fi

# Detect API adapters
API_ADAPTERS=""
if [ -d "${ROOT_DIR}/internal/adapters" ]; then
  for adapter in "${ROOT_DIR}/internal/adapters"/*/; do
    adapter_name=$(basename "$adapter")
    API_ADAPTERS="${API_ADAPTERS}              { label: '${adapter_name}', link: '/api/adapters/${adapter_name}' },
"
  done
fi

# Check if API Reference should be included
# Include if: NODE_ENV is not "production" OR package_name is "go-cli-template"
INCLUDE_API_REFERENCE=false
if [ "${NODE_ENV:-}" != "production" ] || [ "${PROJECT_NAME}" = "go-cli-template" ]; then
  INCLUDE_API_REFERENCE=true
fi

# Build API Reference section if needed
API_REFERENCE_SECTION=""
if [ "$INCLUDE_API_REFERENCE" = true ]; then
  API_REFERENCE_SECTION="  {
    label: 'API Reference',
    items: [
${API_PACKAGES}      {
        label: 'Adapters',
        items: [
${API_ADAPTERS}        ],
      },
    ],
  },"
fi

# Generate sidebar.mjs with dynamic environment check
cat >"$DOCS_SIDEBAR" <<EOF
import config from './config.mjs';

const apiReference = {
  label: 'API Reference',
  items: [
${API_PACKAGES}    {
      label: 'Adapters',
      items: [
${API_ADAPTERS}      ],
    },
  ],
};

const sidebar = [
  {
    label: '${PROJECT_NAME}',
    link: '/',
  },
  {
    label: 'Install',
    link: '/install',
  },
  {
    label: 'Commands',
    items: [
      { label: '${PROJECT_NAME}', link: '/commands/${PROJECT_NAME}' },
${COMMANDS}    ],
  },
  {
    label: 'Configuration',
    link: '/configuration',
  },
];

// Add API Reference if not in production or if this is go-cli-template
const isProduction = process.env.NODE_ENV === 'production';
const projectName = '${PROJECT_NAME}';
if (!isProduction || projectName === 'go-cli-template') {
  sidebar.push(apiReference);
}

export default sidebar;
EOF

echo "  ✓ Generated sidebar.mjs with detected commands and API packages"

echo "📝 Generating content pages..."

DOCS_CONTENT_DIR="docs/src/content/docs"

# Generate index page from README.md
if [ -f "README.md" ]; then
  convert_with_frontmatter "README.md" "${DOCS_CONTENT_DIR}/index.md" \
    "${PROJECT_NAME}" "${DESCRIPTION}"
  echo "  ✓ Generated index.md from README.md"
fi

# Generate install page from INSTALL.md
if [ -f "INSTALL.md" ]; then
  convert_with_frontmatter "INSTALL.md" "${DOCS_CONTENT_DIR}/install.md" \
    "Install" "Installation instructions for ${PROJECT_NAME}"
  echo "  ✓ Generated install.md from INSTALL.md"
fi

# Generate configuration page from CONFIG.md
if [ -f "CONFIG.md" ]; then
  convert_with_frontmatter "CONFIG.md" "${DOCS_CONTENT_DIR}/configuration.md" \
    "Configuration" "Configuration options for ${PROJECT_NAME}"
  echo "  ✓ Generated configuration.md from CONFIG.md"
fi

# Create commands directory
mkdir -p "${DOCS_CONTENT_DIR}/commands"

# Generate root command page from root.go
if [ -f "${CMD_DIR}/root.go" ]; then
  # For root command, use the description from package.toml
  ROOT_SHORT="${DESCRIPTION}"

  # Extract godoc comment for root command (supports both // and /* */ style)
  ROOT_GODOC=$(awk '
    /^\/\*$/ {
      in_block = 1
      comment = ""
      next
    }
    in_block && /\*\// {
      in_block = 0
      print comment
      exit
    }
    in_block {
      if (comment == "") {
        comment = $0
      } else {
        comment = comment "\n" $0
      }
    }
  ' "${CMD_DIR}/root.go")

  cat >"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md" <<EOF
---
title: ${PROJECT_NAME}
description: ${ROOT_SHORT}
---

${ROOT_SHORT}

## Usage

\`\`\`bash
${PROJECT_NAME} [alias]
${PROJECT_NAME} [command]
\`\`\`
EOF

  # Add godoc description if available
  if [ -n "$ROOT_GODOC" ]; then
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "## Description" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "$ROOT_GODOC" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
  fi

  # Extract flags from root.go
  root_flags=$(awk '
    /Flags\(\)\..*VarP?\(/ {
      line = $0
      if (match(line, /Flags\(\)\.(Bool|String|Int)VarP?\([^,]+, *"([^"]+)", *"([^"]*)", *[^,]+, *"([^"]+)"\)/, arr)) {
        flag_type = arr[1]
        flag_long = arr[2]
        flag_short = arr[3]
        flag_desc = arr[4]

        if (flag_short != "") {
          flag_col = "-" flag_short ", --" flag_long
        } else {
          flag_col = "--" flag_long
        }

        type_col = tolower(flag_type)

        print "| `" flag_col "` | " type_col " | " flag_desc " |"
      }
    }
  ' "${CMD_DIR}/root.go")

  # Add flags table if flags were found
  if [ -n "$root_flags" ]; then
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "## Flags" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "| Flag | Type | Description |" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "|------|------|-------------|" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
    echo "$root_flags" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
  fi

  echo "" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
  echo "## Available Commands" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
  echo "" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"

  # List all subcommands
  for cmd_file in "$CMD_DIR"/*.go; do
    if [[ "$cmd_file" == *"_test.go" ]] || [[ "$cmd_file" == *"/main.go" ]] || [[ "$cmd_file" == *"/root.go" ]]; then
      continue
    fi

    cmd_name=$(basename "$cmd_file" .go | sed 's/_cmd$//')
    cmd_display=$(echo "$cmd_name" | sed 's/_/ /g')
    cmd_url=$(echo "$cmd_name" | sed 's/_/-/g')

    # Extract Short description - handle both quoted strings and variables
    cmd_short=$(awk '/Short:/ {
      if (match($0, /Short: *"([^"]*)"/, arr)) {
        print arr[1]
      }
    }' "$cmd_file" | head -1)

    echo "- [\`${cmd_display}\`](/commands/${cmd_url}) - ${cmd_short}" >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md"
  done

  cat >>"${DOCS_CONTENT_DIR}/commands/${PROJECT_NAME}.md" <<EOF

## Source

See [root.go](${REPOSITORY}/blob/main/cmd/${PROJECT_NAME}/root.go) for implementation details.
EOF

  echo "  ✓ Generated commands/${PROJECT_NAME}.md"
fi

# Generate documentation for each command
for cmd_file in "$CMD_DIR"/*.go; do
  # Skip test files, main.go, and root.go
  if [[ "$cmd_file" == *"_test.go" ]] || [[ "$cmd_file" == *"/main.go" ]] || [[ "$cmd_file" == *"/root.go" ]]; then
    continue
  fi

  # Extract command name from filename
  cmd_name=$(basename "$cmd_file" .go | sed 's/_cmd$//')
  cmd_display=$(echo "$cmd_name" | sed 's/_/ /g')
  cmd_url=$(echo "$cmd_name" | sed 's/_/-/g')

  # Extract command information from the Go file using awk for better parsing
  cmd_use=$(awk '/Use:/ {
    if (match($0, /Use: *"([^"]*)"/, arr)) {
      print arr[1]
    }
  }' "$cmd_file" | head -1)

  cmd_short=$(awk '/Short:/ {
    if (match($0, /Short: *"([^"]*)"/, arr)) {
      print arr[1]
    }
  }' "$cmd_file" | head -1)

  # Extract godoc comment (supports both // and /* */ style)
  cmd_godoc=$(awk '
    /^\/\*$/ {
      in_block = 1
      comment = ""
      next
    }
    in_block && /\*\// {
      in_block = 0
      print comment
      exit
    }
    in_block {
      if (comment == "") {
        comment = $0
      } else {
        comment = comment "\n" $0
      }
    }
  ' "$cmd_file")

  # Use display name if Use is empty
  if [ -z "$cmd_use" ]; then
    cmd_use="$cmd_display"
  fi

  # Generate command documentation
  cat >"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md" <<EOF
---
title: ${cmd_display}
description: ${cmd_short}
---

${cmd_short}

## Usage

\`\`\`bash
${PROJECT_NAME} ${cmd_use}
\`\`\`
EOF

  # Add godoc description if available
  if [ -n "$cmd_godoc" ]; then
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "## Description" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "$cmd_godoc" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
  fi

  # Extract flags from the command file
  flags=$(awk '
    /Flags\(\)\..*VarP?\(/ {
      # Extract flag information
      line = $0
      # Try to match the pattern for flag definitions
      if (match(line, /Flags\(\)\.(Bool|String|Int)VarP?\([^,]+, *"([^"]+)", *"([^"]*)", *[^,]+, *"([^"]+)"\)/, arr)) {
        flag_type = arr[1]
        flag_long = arr[2]
        flag_short = arr[3]
        flag_desc = arr[4]
        
        # Build flag column
        if (flag_short != "") {
          flag_col = "-" flag_short ", --" flag_long
        } else {
          flag_col = "--" flag_long
        }
        
        # Determine type
        type_col = tolower(flag_type)
        
        print "| `" flag_col "` | " type_col " | " flag_desc " |"
      }
    }
  ' "$cmd_file")

  # Add flags table if flags were found
  if [ -n "$flags" ]; then
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "## Flags" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "| Flag | Type | Description |" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "|------|------|-------------|" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
    echo "$flags" >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md"
  fi

  # Add source link
  cat >>"${DOCS_CONTENT_DIR}/commands/${cmd_url}.md" <<EOF

## Source

See [$(basename "$cmd_file")](${REPOSITORY}/blob/main/cmd/${PROJECT_NAME}/$(basename "$cmd_file")) for implementation details.
EOF

  echo "  ✓ Generated commands/${cmd_url}.md"
done

echo "🔧 Checking for gomarkdoc..."
if ! command -v gomarkdoc &>/dev/null; then
  echo "📦 Installing gomarkdoc..."
  go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
fi

echo "📝 Generating API documentation..."

# Generate into a temp directory to avoid Astro seeing a partially-written api/ dir
DOCS_API_TEMP="$(mktemp -d)"
trap 'rm -rf "$DOCS_API_TEMP"' EXIT

# Generate docs for each internal package
for pkg in internal/*/; do
  pkg_name=$(basename "$pkg")

  # Skip test utilities and adapters subdirectories
  if [[ "$pkg_name" == "testutil" ]]; then
    continue
  fi

  echo "  - Processing $pkg_name..."

  # Generate to temp file first
  gomarkdoc \
    --output "${DOCS_API_TEMP}/${pkg_name}.raw.md" \
    --template-file "file=${ROOT_DIR}/docs/templates/file.gotxt" \
    --footer $'## Source\n\nSee [internal/'"${pkg_name}"$'/]('"${REPOSITORY}"$'/blob/main/internal/'"${pkg_name}"$'/) for implementation details.' \
    "./$pkg" 2>/dev/null || {
    echo "    ⚠️  No exported symbols in $pkg_name"
    continue
  }

  # Add frontmatter and content (skip HTML comment and any frontmatter that gomarkdoc added)
  {
    echo "---"
    echo "title: ${pkg_name}"
    echo "description: API documentation for the ${pkg_name} package"
    echo "---"
    echo ""
    sed '1,/^# /d' "${DOCS_API_TEMP}/${pkg_name}.raw.md"
  } >"${DOCS_API_TEMP}/${pkg_name}.md"
  rm -f "${DOCS_API_TEMP}/${pkg_name}.raw.md"
done

# Generate docs for adapters
echo "  - Processing adapters..."
mkdir -p "${DOCS_API_TEMP}/adapters"

for adapter in internal/adapters/*/; do
  adapter_name=$(basename "$adapter")
  echo "    - Processing adapters/$adapter_name..."

  # Generate to temp file first
  gomarkdoc \
    --output "${DOCS_API_TEMP}/adapters/${adapter_name}.raw.md" \
    --template-file "file=${ROOT_DIR}/docs/templates/file.gotxt" \
    --footer $'## Source\n\nSee [internal/adapters/'"${adapter_name}"$'/]('"${REPOSITORY}"$'/blob/main/internal/adapters/'"${adapter_name}"$'/) for implementation details.' \
    "./$adapter" 2>/dev/null || {
    echo "      ⚠️  No exported symbols in $adapter_name"
    continue
  }

  # Add frontmatter and content (skip HTML comment and any frontmatter that gomarkdoc added)
  {
    echo "---"
    echo "title: adapters/${adapter_name}"
    echo "description: API documentation for the ${adapter_name} adapter"
    echo "---"
    echo ""
    sed '1,/^# /d' "${DOCS_API_TEMP}/adapters/${adapter_name}.raw.md"
  } >"${DOCS_API_TEMP}/adapters/${adapter_name}.md"
  rm -f "${DOCS_API_TEMP}/adapters/${adapter_name}.raw.md"
done

# Atomically swap in the newly generated api/ directory
rm -rf "$DOCS_API_DIR"
mv "$DOCS_API_TEMP" "$DOCS_API_DIR"
trap - EXIT

echo "✅ API documentation generated successfully!"
echo "📁 Output: $DOCS_API_DIR"
