#!/usr/bin/env bash
# Deploy all packages
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [ -z "${VERSION}" ]; then
	VERSION=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "version")
	echo "📦 Using version from package.toml: v${VERSION}"
fi

PACKAGE_NAME=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "package_name")
NAME=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "name")
PACKAGE_NAME="${PACKAGE_NAME:-$NAME}"
REPOSITORY=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "repository")
GITHUB_USER="$(echo "${REPOSITORY}" | sed -E 's|https://github.com/([^/]+)/.*|\1|')"

echo "🚀 Deploying all packages for version ${VERSION}..."
echo ""

"${ROOT_DIR}/scripts/deploy_aur.sh" "${VERSION}"
echo ""
"${ROOT_DIR}/scripts/deploy_homebrew.sh" "${VERSION}"
echo ""

echo "✅ All packages deployed!"
echo ""
echo "AUR: https://aur.archlinux.org/packages/${PACKAGE_NAME}"
echo "Homebrew: https://github.com/${GITHUB_USER}/homebrew-${PACKAGE_NAME}"
