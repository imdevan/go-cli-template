#!/usr/bin/env bash
# Shared utility functions for build scripts

# Parse a key from a TOML file
# Usage: parse_toml_key <file> <key>
parse_toml_key() {
  local file=$1
  local key=$2
  grep "^${key} = " "$file" | sed 's/^[^=]*= *"\(.*\)"$/\1/'
}

# Download a file and return its SHA256 hash
# Usage: download_and_hash <url>
download_and_hash() {
  local url=$1
  local temp_file=$(mktemp)

  if ! curl -sL "$url" -o "$temp_file"; then
    rm -f "$temp_file"
    return 1
  fi

  sha256sum "$temp_file" | awk '{print $1}'
  rm -f "$temp_file"
}
# Add YAML frontmatter to a markdown file
# Usage: add_frontmatter <output_file> <title> <description>
add_frontmatter() {
  local output_file=$1
  local title=$2
  local description=$3
  
  {
    echo "---"
    echo "title: ${title}"
    echo "description: ${description}"
    echo "---"
    echo ""
  } > "$output_file"
}

# Add frontmatter and skip first heading from source file
# Usage: convert_with_frontmatter <source_file> <output_file> <title> <description>
convert_with_frontmatter() {
  local source_file=$1
  local output_file=$2
  local title=$3
  local description=$4
  
  {
    echo "---"
    echo "title: ${title}"
    echo "description: ${description}"
    echo "---"
    echo ""
    # Skip the first heading from source and output the rest
    sed '1{/^# /d;}' "$source_file"
  } > "$output_file"
}
