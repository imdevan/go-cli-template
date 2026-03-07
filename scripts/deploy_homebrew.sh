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
REPO_URL=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "repository")
GITHUB_USER="$(echo "${REPO_URL}" | sed -E 's|https://github.com/([^/]+)/.*|\1|')"

echo "🍺 Deploying Homebrew formula for version ${VERSION}..."

cd "${ROOT_DIR}/homebrew-${PACKAGE_NAME}"

# Add GitHub remote if not configured
if ! git remote get-url origin &>/dev/null; then
	echo "Adding GitHub remote..."
	git remote add origin "git@github.com:${GITHUB_USER}/homebrew-${PACKAGE_NAME}.git"
fi

# Create GitHub repository if it doesn't exist
if ! gh repo view "$(git remote get-url origin | sed -E 's|git@github.com:||; s|\.git$||')" &>/dev/null 2>&1; then
	echo "Creating GitHub repository..."
	gh repo create "homebrew-${PACKAGE_NAME}" --public
fi

git add "Formula/${PACKAGE_NAME}.rb"

# Only commit if there are staged changes
if ! git diff --cached --quiet; then
	git commit -m "Update ${PACKAGE_NAME} to v${VERSION}"
fi

# Use -u on first push to set upstream tracking
if git rev-parse --abbrev-ref --symbolic-full-name @{u} &>/dev/null 2>&1; then
	git push
else
	echo "First push — setting upstream..."
	git push -u origin main
fi

echo "✅ Homebrew formula deployed!"
echo "   Install with:"
echo "   brew tap ${GITHUB_USER}/${PACKAGE_NAME} && brew install ${PACKAGE_NAME}"
