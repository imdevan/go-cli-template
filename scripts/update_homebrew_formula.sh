#!/usr/bin/env bash
# Update Homebrew formula with new version
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_TOML="${ROOT_DIR}/internal/package/package.toml"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [[ -z "${VERSION}" ]]; then
  VERSION=$(parse_toml_key "${PACKAGE_TOML}" "version")
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
GITHUB_USER="$(echo "${REPO_URL}" | sed -E 's|https://github.com/([^/]+)/.*|\1|')"

TAP_DIR="${ROOT_DIR}/homebrew-${PACKAGE_NAME}"
FORMULA_PATH="${TAP_DIR}/Formula/${PACKAGE_NAME}.rb"

if [[ ! -d "${TAP_DIR}" ]]; then
  echo "❌ Homebrew tap not found at: ${TAP_DIR}"
  echo "Run 'just init-homebrew-tap' first"
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
DARWIN_AMD64_SHA=$(get_platform_sha "darwin-amd64")
DARWIN_ARM64_SHA=$(get_platform_sha "darwin-arm64")
LINUX_AMD64_SHA=$(get_platform_sha "linux-amd64")
LINUX_ARM64_SHA=$(get_platform_sha "linux-arm64")

CLEAN_REPO="${REPO_URL%/}"
REPO_PATH="${CLEAN_REPO#https://github.com/}"
DARWIN_AMD64_URL="https://github.com/${REPO_PATH}/releases/download/v${VERSION}/${NAME}-darwin-amd64.tar.gz"
DARWIN_ARM64_URL="https://github.com/${REPO_PATH}/releases/download/v${VERSION}/${NAME}-darwin-arm64.tar.gz"
LINUX_AMD64_URL="https://github.com/${REPO_PATH}/releases/download/v${VERSION}/${NAME}-linux-amd64.tar.gz"
LINUX_ARM64_URL="https://github.com/${REPO_PATH}/releases/download/v${VERSION}/${NAME}-linux-arm64.tar.gz"

# Update formula
CLASS_NAME="$(echo "${PACKAGE_NAME}" | sed 's/-/ /g; s/\b\(.\)/\u\1/g; s/ //g')"

cat >"${FORMULA_PATH}" <<EOF
class ${CLASS_NAME} < Formula
  desc "${DESCRIPTION}"
  homepage "${HOMEPAGE}"
  version "${VERSION}"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "${DARWIN_AMD64_URL}"
      sha256 "${DARWIN_AMD64_SHA}"
    elsif Hardware::CPU.arm?
      url "${DARWIN_ARM64_URL}"
      sha256 "${DARWIN_ARM64_SHA}"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "${LINUX_AMD64_URL}"
      sha256 "${LINUX_AMD64_SHA}"
    elsif Hardware::CPU.arm?
      url "${LINUX_ARM64_URL}"
      sha256 "${LINUX_ARM64_SHA}"
    end
  end

  def install
    binary = OS.mac? ? "${NAME}-darwin-" : "${NAME}-linux-"
    binary += Hardware::CPU.intel? ? "amd64" : "arm64"
    bin.install binary => "${NAME}"
  end

  test do
    assert_match "v${VERSION}", shell_output("#{bin}/${NAME} --version")
  end
end
EOF

echo "✅ Updated formula: ${FORMULA_PATH}"
echo ""
echo "Next steps:"
echo "1. Test the formula locally:"
echo "   brew tap ${GITHUB_USER}/homebrew-${PACKAGE_NAME} ${TAP_DIR}"
echo "   brew install --build-from-source ${PACKAGE_NAME}"
echo "   brew untap ${GITHUB_USER}/homebrew-${PACKAGE_NAME}"
echo "2. Deploy:"
echo "   just deploy-homebrew"
