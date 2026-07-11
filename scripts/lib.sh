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
