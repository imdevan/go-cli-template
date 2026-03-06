#!/usr/bin/env bash
# Deploy AUR package
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

echo "📦 Deploying AUR package for version ${VERSION}..."

# Set up SSH command with specific key if AUR_SSH_KEY is set
SSH_CMD="ssh"
if [ -n "${AUR_SSH_KEY:-}" ]; then
	SSH_CMD="ssh -i ${AUR_SSH_KEY}"
	echo "Using SSH key: ${AUR_SSH_KEY}"
fi

cd "${ROOT_DIR}/aur-${PACKAGE_NAME}"
GIT_SSH_COMMAND="${SSH_CMD}" git add PKGBUILD .SRCINFO
git commit -m "Update ${PACKAGE_NAME} to v${VERSION}"
GIT_SSH_COMMAND="${SSH_CMD}" git push

echo "✅ AUR package deployed!"
echo "   https://aur.archlinux.org/packages/${PACKAGE_NAME}"
