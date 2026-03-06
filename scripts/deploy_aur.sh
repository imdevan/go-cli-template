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

# Add AUR remote if not configured
if ! git remote get-url origin &>/dev/null; then
	echo "Adding AUR remote..."
	git remote add origin "ssh://aur@aur.archlinux.org/${PACKAGE_NAME}.git"
fi

GIT_SSH_COMMAND="${SSH_CMD}" git add PKGBUILD .SRCINFO

# Only commit if there are staged changes
if ! git diff --cached --quiet; then
	git commit -m "Update ${PACKAGE_NAME} to v${VERSION}"
fi

# Use -u on first push to set upstream tracking
if git rev-parse --abbrev-ref --symbolic-full-name @{u} &>/dev/null 2>&1; then
	GIT_SSH_COMMAND="${SSH_CMD}" git push
else
	echo "First push — setting upstream..."
	GIT_SSH_COMMAND="${SSH_CMD}" git push -u origin master
fi

echo "✅ AUR package deployed!"
echo "   https://aur.archlinux.org/packages/${PACKAGE_NAME}"
