#!/usr/bin/env bash
# Create and push git tag
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [ -z "${VERSION}" ]; then
	VERSION=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "version")
	echo "📦 Using version from package.toml: ${VERSION}"
fi

# Remove 'v' prefix if present
VERSION="${VERSION#v}"
TAG="v${VERSION}"

echo "🏷️  Creating and pushing tag ${TAG}..."

if git rev-parse "${TAG}" >/dev/null 2>&1; then
	echo "⚠️  Tag ${TAG} already exists"
	read -p "Do you want to delete and recreate it? (y/N) " -n 1 -r
	echo
	if [[ $REPLY =~ ^[Yy]$ ]]; then
		git tag -d "${TAG}"
		git push origin ":refs/tags/${TAG}" 2>/dev/null || true
		echo "✅ Deleted existing tag"
	else
		echo "Aborted."
		exit 1
	fi
fi

git tag -a "${TAG}" -m "Release ${TAG}"
git push origin "${TAG}"
echo "✅ Tag ${TAG} created and pushed"
