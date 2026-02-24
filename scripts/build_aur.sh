#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMPLATE_PATH="${ROOT_DIR}/aur/PKGBUILD.template"
OUTPUT_DIR="${ROOT_DIR}/dist/aur"
OUTPUT_PATH="${OUTPUT_DIR}/PKGBUILD"
PACKAGE_TOML="${ROOT_DIR}/package.toml"

VERSION="${VERSION:-}"
if [[ -z "${VERSION}" ]]; then
  VERSION="$(git -C "${ROOT_DIR}" describe --tags --abbrev=0 2>/dev/null)"
fi

PKGVER="${VERSION#v}"

# Read repository URL from package.toml
REPO_URL=""
if [[ -f "${PACKAGE_TOML}" ]]; then
  REPO_URL="$(grep '^repository = ' "${PACKAGE_TOML}" | sed 's/^repository = "\(.*\)"$/\1/' | tr -d '\n')"
fi

SOURCE_URL="${AUR_SOURCE_URL:-}"
SHA256="${AUR_SOURCE_SHA256:-}"

if [[ -z "${SOURCE_URL}" ]]; then
  if [[ -z "${REPO_URL}" ]]; then
    echo "Error: repository field in package.toml is empty and AUR_SOURCE_URL is not set." >&2
    exit 1
  fi
  SOURCE_URL="${REPO_URL}/archive/refs/tags/v${PKGVER}.tar.gz"
fi

if [[ -z "${SHA256}" ]]; then
  echo "AUR_SOURCE_SHA256 is required to generate ${OUTPUT_PATH}." >&2
  exit 1
fi

mkdir -p "${OUTPUT_DIR}"

sed \
  -e "s|__PKGVER__|${PKGVER}|g" \
  -e "s|__SOURCE_URL__|${SOURCE_URL}|g" \
  -e "s|__SHA256__|${SHA256}|g" \
  "${TEMPLATE_PATH}" >"${OUTPUT_PATH}"

echo "Wrote ${OUTPUT_PATH}"
