#!/usr/bin/env bash
# Create a GitHub release with cross-platform binaries
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

VERSION="${1:-}"
if [ -z "${VERSION}" ]; then
	VERSION=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "version")
	echo "📦 Using version from package.toml: v${VERSION}"
fi

# Normalize version
VERSION="${VERSION#v}"
TAG="v${VERSION}"

NAME=$(parse_toml_key "${ROOT_DIR}/internal/package/package.toml" "name")
CMD_PKG="./cmd/${NAME}"
DIST_DIR="${ROOT_DIR}/dist/${TAG}"

TARGETS=(
	"linux/amd64"
	"linux/arm64"
	"darwin/amd64"
	"darwin/arm64"
	"windows/amd64"
)

echo "🔨 Building binaries for ${TAG}..."
mkdir -p "${DIST_DIR}"

ASSETS=()
for target in "${TARGETS[@]}"; do
	GOOS="${target%/*}"
	GOARCH="${target#*/}"
	output="${NAME}-${GOOS}-${GOARCH}"
	[ "${GOOS}" = "windows" ] && output="${output}.exe"
	out_path="${DIST_DIR}/${output}"

	echo "  - ${GOOS}/${GOARCH}..."
	GOOS="${GOOS}" GOARCH="${GOARCH}" go build -ldflags="-s -w" -o "${out_path}" "${CMD_PKG}"

	# Create tar.gz (zip for windows)
	if [ "${GOOS}" = "windows" ]; then
		archive="${DIST_DIR}/${NAME}-${GOOS}-${GOARCH}.zip"
		zip -j "${archive}" "${out_path}"
	else
		archive="${DIST_DIR}/${NAME}-${GOOS}-${GOARCH}.tar.gz"
		tar -czf "${archive}" -C "${DIST_DIR}" "${output}"
	fi
	ASSETS+=("${archive}")
done

echo "🚀 Creating GitHub release ${TAG}..."
gh release create "${TAG}" \
	--title "${TAG}" \
	--generate-notes \
	--verify-tag \
	"${ASSETS[@]}"

echo "✅ GitHub release ${TAG} created"
echo "🧹 Cleaning up dist..."
rm -rf "${DIST_DIR}"
