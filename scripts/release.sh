#!/usr/bin/env bash
# Update all package files for a release
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [ -z "${VERSION}" ]; then
	VERSION=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "version")
	echo "📦 Using version from package.toml: v${VERSION}"
else
	echo "📦 Using provided version: ${VERSION}"
fi

PACKAGE_NAME=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "package_name")
NAME=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "name")
PACKAGE_NAME="${PACKAGE_NAME:-$NAME}"

echo ""
echo "🚀 Starting release process for version ${VERSION}..."
echo ""
echo "Step 1: Updating AUR PKGBUILD..."
"${ROOT_DIR}/scripts/update_aur_pkgbuild.sh" "${VERSION}"
echo ""
echo "Step 2: Updating Homebrew formula..."
"${ROOT_DIR}/scripts/update_homebrew_formula.sh" "${VERSION}"
echo ""
echo "✅ Release packages updated!"
echo ""
echo "Next steps:"
echo "1. Test AUR package:"
echo "   cd aur-${PACKAGE_NAME} && makepkg -si"
echo ""
echo "2. Test Homebrew formula:"
echo "   brew install --build-from-source homebrew-${PACKAGE_NAME}/Formula/${PACKAGE_NAME}.rb"
echo ""
echo "3. Deploy to AUR:"
echo "   just deploy-aur ${VERSION}"
echo ""
echo "4. Deploy to Homebrew:"
echo "   just deploy-homebrew ${VERSION}"
