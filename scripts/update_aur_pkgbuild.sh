#!/usr/bin/env bash
# Update AUR PKGBUILD with new version
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_TOML="${ROOT_DIR}/internal/package/package.toml"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [[ -z "${VERSION}" ]]; then
  VERSION="$(parse_toml_key "${PACKAGE_TOML}" "version")"
fi

# Remove 'v' prefix if present
VERSION="${VERSION#v}"

# Read package metadata
NAME="$(parse_toml_key "${PACKAGE_TOML}" "name")"
PACKAGE_NAME="$(parse_toml_key "${PACKAGE_TOML}" "package_name")"
# Fall back to name if package_name is not set
PACKAGE_NAME="${PACKAGE_NAME:-$NAME}"
REPO_URL="$(parse_toml_key "${PACKAGE_TOML}" "repository")"
DESCRIPTION="$(parse_toml_key "${PACKAGE_TOML}" "description")"
HOMEPAGE="$(parse_toml_key "${PACKAGE_TOML}" "homepage")"
AUTHOR="$(parse_toml_key "${PACKAGE_TOML}" "author")"

AUR_DIR="${ROOT_DIR}/aur-${PACKAGE_NAME}"
PKGBUILD_PATH="${AUR_DIR}/PKGBUILD"

if [[ ! -d "${AUR_DIR}" ]]; then
  echo "❌ AUR repository not found at: ${AUR_DIR}"
  echo "Run 'just init-aur-repo' first"
  exit 1
fi

# Helper to get SHA256 of binary assets
get_platform_sha() {
  local platform="$1"
  local ext="${2:-tar.gz}"
  local local_file="${ROOT_DIR}/dist/v${VERSION}/${NAME}-${platform}.${ext}"
  local clean_repo="${REPO_URL%/}"
  local repo_path="${clean_repo#https://github.com/}"
  local url="https://github.com/${repo_path}/releases/download/v${VERSION}/${NAME}-${platform}.${ext}"
  
  if [[ -f "${local_file}" ]]; then
    sha256sum "${local_file}" | awk '{print $1}'
  else
    echo "📥 Downloading to calculate hash for ${platform}..." >&2
    if ! sha=$(download_and_hash "${url}"); then
      echo "❌ Failed to download: ${url}" >&2
      exit 1
    fi
    echo "${sha}"
  fi
}

echo "🔍 Calculating SHA256 hashes for binary assets..."
LINUX_AMD64_SHA=$(get_platform_sha "linux-amd64")
LINUX_ARM64_SHA=$(get_platform_sha "linux-arm64")

CLEAN_REPO="${REPO_URL%/}"
REPO_PATH="${CLEAN_REPO#https://github.com/}"
LINUX_AMD64_URL="https://github.com/${REPO_PATH}/releases/download/v${VERSION}/${NAME}-linux-amd64.tar.gz"
LINUX_ARM64_URL="https://github.com/${REPO_PATH}/releases/download/v${VERSION}/${NAME}-linux-arm64.tar.gz"

# Update PKGBUILD
cat >"${PKGBUILD_PATH}" <<EOF
# Maintainer: ${AUTHOR}
pkgname=${PACKAGE_NAME}
_binname=${NAME}
pkgver=${VERSION}
pkgrel=1
pkgdesc="${DESCRIPTION}"
arch=('x86_64' 'aarch64')
url="${HOMEPAGE}"
license=('MIT')
depends=()

source_x86_64=("\${_binname}-linux-amd64-\${pkgver}.tar.gz::${LINUX_AMD64_URL}")
source_aarch64=("\${_binname}-linux-arm64-\${pkgver}.tar.gz::${LINUX_ARM64_URL}")
sha256sums_x86_64=('${LINUX_AMD64_SHA}')
sha256sums_aarch64=('${LINUX_ARM64_SHA}')

package() {
  if [ "\${CARCH}" = "x86_64" ]; then
    install -Dm755 "\${srcdir}/\${_binname}-linux-amd64" "\${pkgdir}/usr/bin/\${_binname}"
  elif [ "\${CARCH}" = "aarch64" ]; then
    install -Dm755 "\${srcdir}/\${_binname}-linux-arm64" "\${pkgdir}/usr/bin/\${_binname}"
  fi
}
EOF

# Generate .SRCINFO
cd "${AUR_DIR}"
if command -v makepkg &>/dev/null; then
  makepkg --printsrcinfo >.SRCINFO
  echo "✅ Generated .SRCINFO"
else
  echo "⚠️  makepkg not found, skipping .SRCINFO generation"
  echo "   You'll need to run 'makepkg --printsrcinfo > .SRCINFO' manually"
fi

echo "✅ Updated PKGBUILD: ${PKGBUILD_PATH}"
echo ""
echo "Next steps:"
echo "1. Test the package locally:"
echo "   cd ${AUR_DIR} && makepkg -si"
echo "2. Deploy to AUR:"
echo "   just deploy-aur ${VERSION}"
