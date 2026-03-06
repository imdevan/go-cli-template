#!/usr/bin/env bash
# Deploy Homebrew formula
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [ -z "${VERSION}" ]; then
	VERSION=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "version")
fi

PACKAGE_NAME=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "package_name")
NAME=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "name")
PACKAGE_NAME="${PACKAGE_NAME:-$NAME}"

echo "🍺 Deploying Homebrew formula for version ${VERSION}..."

cd "${ROOT_DIR}/homebrew-${PACKAGE_NAME}"
git add "Formula/${PACKAGE_NAME}.rb"
git commit -m "Update ${PACKAGE_NAME} to v${VERSION}"
git push

echo "✅ Homebrew formula deployed!"
