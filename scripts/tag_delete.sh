#!/usr/bin/env bash
# Delete git tag locally and remotely
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [ -z "${VERSION}" ]; then
	VERSION=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "version")
fi

VERSION="${VERSION#v}"
TAG="v${VERSION}"

echo "🗑️  Deleting tag ${TAG}..."
git tag -d "${TAG}" 2>/dev/null || echo "Local tag not found"
git push origin ":refs/tags/${TAG}" 2>/dev/null || echo "Remote tag not found"
echo "✅ Tag ${TAG} deleted"
