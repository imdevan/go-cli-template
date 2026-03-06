#!/usr/bin/env bash
# Initialize AUR repository
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PACKAGE_TOML="${ROOT_DIR}/internal/package/package.toml"

# Source shared utilities
. "${ROOT_DIR}/scripts/lib.sh"

# Read package metadata
NAME="$(parse_toml_key "${PACKAGE_TOML}" "name")"
PACKAGE_NAME="$(parse_toml_key "${PACKAGE_TOML}" "package_name")"
# Fall back to name if package_name is not set
PACKAGE_NAME="${PACKAGE_NAME:-$NAME}"
DESCRIPTION="$(parse_toml_key "${PACKAGE_TOML}" "description")"
HOMEPAGE="$(parse_toml_key "${PACKAGE_TOML}" "homepage")"
REPO_URL="$(parse_toml_key "${PACKAGE_TOML}" "repository")"
AUTHOR="$(parse_toml_key "${PACKAGE_TOML}" "author")"

AUR_DIR="${ROOT_DIR}/aur-${PACKAGE_NAME}"

echo "📦 Initializing AUR repository..."
echo "   Package name: ${PACKAGE_NAME}"
echo "   Binary name: ${NAME}"
echo "   Location: ${AUR_DIR}"

# Create AUR directory if it doesn't exist
if [[ -d "${AUR_DIR}" ]]; then
  echo "⚠️  AUR directory already exists: ${AUR_DIR}"
  read -p "Do you want to reinitialize it? (y/N) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
  fi
  rm -rf "${AUR_DIR}"
fi

mkdir -p "${AUR_DIR}"

# Initialize git repository
cd "${AUR_DIR}"
git init
git branch -M master # AUR uses master branch

# Create .gitignore
cat >"${AUR_DIR}/.gitignore" <<EOF
*.tar.gz
*.tar.xz
*.zip
pkg/
src/
*.pkg.tar.*
EOF

# Create initial PKGBUILD
cat >"${AUR_DIR}/PKGBUILD" <<EOF
# Maintainer: ${AUTHOR}
pkgname=${PACKAGE_NAME}
_binname=${NAME}
pkgver=0.1.0
pkgrel=1
pkgdesc="${DESCRIPTION}"
arch=('x86_64' 'aarch64')
url="${HOMEPAGE}"
license=('MIT')
depends=()
makedepends=('go')
source=("\${_binname}-\${pkgver}.tar.gz::${REPO_URL}archive/refs/tags/v\${pkgver}.tar.gz")
sha256sums=('REPLACE_WITH_ACTUAL_SHA256')

build() {
  cd "\${_binname}-\${pkgver}"
  export CGO_ENABLED=0
  export GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw"
  go build -ldflags="-s -w" -o \${_binname} ./cmd/\${_binname}
}

package() {
  cd "\${_binname}-\${pkgver}"
  install -Dm755 \${_binname} "\${pkgdir}/usr/bin/\${_binname}"
  if [ -f LICENSE ]; then
    install -Dm644 LICENSE "\${pkgdir}/usr/share/licenses/\${pkgname}/LICENSE"
  fi
}
EOF

# Create README
cat >"${AUR_DIR}/README.md" <<EOF
# AUR Package for ${PACKAGE_NAME}

This is the AUR (Arch User Repository) package for [${NAME}](${HOMEPAGE}).

The package is named \`${PACKAGE_NAME}\` but installs the binary as \`${NAME}\`.

## Installation

### Using an AUR helper (recommended)

\`\`\`bash
yay -S ${PACKAGE_NAME}
# or
paru -S ${PACKAGE_NAME}
\`\`\`

### Manual installation

\`\`\`bash
git clone https://aur.archlinux.org/${PACKAGE_NAME}.git
cd ${PACKAGE_NAME}
makepkg -si
\`\`\`

## Usage

After installation, use the \`${NAME}\` command:

\`\`\`bash
${NAME} --help
\`\`\`

## Updating

\`\`\`bash
yay -Syu ${PACKAGE_NAME}
# or
paru -Syu ${PACKAGE_NAME}
\`\`\`

## Uninstall

\`\`\`bash
sudo pacman -R ${PACKAGE_NAME}
\`\`\`

## Maintainer

${AUTHOR}
EOF

# Initial commit
git add PKGBUILD README.md .gitignore
git commit -m "Initial commit: AUR package for ${PACKAGE_NAME}"

echo ""
echo "✅ AUR repository initialized at: ${AUR_DIR}"
echo ""
echo "Next steps:"
echo "1. Register an AUR account: https://aur.archlinux.org/register"
echo "2. Add your SSH key to AUR: https://aur.archlinux.org/account"
echo "3. Push the package:"
echo "   cd ${AUR_DIR}"
echo "   git remote add aur ssh://aur@aur.archlinux.org/${PACKAGE_NAME}.git"
echo "   git push -u aur master"
echo "4. Update PKGBUILD with actual release SHA256 using:"
echo "   just update-aur-pkgbuild VERSION"
echo ""
echo "📚 AUR submission guidelines: https://wiki.archlinux.org/title/AUR_submission_guidelines"
