#!/usr/bin/env bash
# Initialize Homebrew tap repository
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

# Extract GitHub username from repository URL
GITHUB_USER="$(echo "${REPO_URL}" | sed -E 's|https://github.com/([^/]+)/.*|\1|')"

TAP_NAME="homebrew-${PACKAGE_NAME}"
TAP_DIR="${ROOT_DIR}/${TAP_NAME}"

echo "🍺 Initializing Homebrew tap repository..."
echo "   Tap name: ${TAP_NAME}"
echo "   Package name: ${PACKAGE_NAME}"
echo "   Binary name: ${NAME}"
echo "   Location: ${TAP_DIR}"

# Create tap directory if it doesn't exist
if [[ -d "${TAP_DIR}" ]]; then
  echo "⚠️  Tap directory already exists: ${TAP_DIR}"
  read -p "Do you want to reinitialize it? (y/N) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
  fi
  rm -rf "${TAP_DIR}"
fi

mkdir -p "${TAP_DIR}/Formula"

# Initialize git repository
cd "${TAP_DIR}"
git init
git branch -M main

# Create README
cat >"${TAP_DIR}/README.md" <<EOF
# Homebrew Tap for ${PACKAGE_NAME}

This is the official Homebrew tap for [${NAME}](${HOMEPAGE}).

The formula is named \`${PACKAGE_NAME}\` but installs the binary as \`${NAME}\`.

## Installation

\`\`\`bash
brew tap ${GITHUB_USER}/${PACKAGE_NAME}
brew install ${PACKAGE_NAME}
\`\`\`

## Usage

After installation, use the \`${NAME}\` command:

\`\`\`bash
${NAME} --help
\`\`\`

## Updating

\`\`\`bash
brew update
brew upgrade ${PACKAGE_NAME}
\`\`\`

## Uninstall

\`\`\`bash
brew uninstall ${PACKAGE_NAME}
brew untap ${GITHUB_USER}/${PACKAGE_NAME}
\`\`\`
EOF

# Create initial formula template
CLASS_NAME="$(echo "${PACKAGE_NAME}" | sed 's/-/ /g; s/\b\(.\)/\u\1/g; s/ //g')"
cat >"${TAP_DIR}/Formula/${PACKAGE_NAME}.rb" <<EOF
class ${CLASS_NAME} < Formula
  desc "${DESCRIPTION}"
  homepage "${HOMEPAGE}"
  url "${REPO_URL}archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_ACTUAL_SHA256"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(ldflags: "-s -w", output: bin/"${NAME}"), "./cmd/${NAME}"
  end

  test do
    assert_match "v0.1.0", shell_output("#{bin}/${NAME} --version")
  end
end
EOF

# Create .gitignore
cat >"${TAP_DIR}/.gitignore" <<EOF
.DS_Store
*.swp
*.swo
*~
EOF

# Initial commit
git add .
git commit -m "Initial commit: Homebrew tap for ${PACKAGE_NAME}"

echo ""
echo "✅ Homebrew tap initialized at: ${TAP_DIR}"
echo ""
echo "Next steps:"
echo "1. Create a GitHub repository: https://github.com/new"
echo "   Repository name: ${TAP_NAME}"
echo "2. Push the tap:"
echo "   cd ${TAP_DIR}"
echo "   git remote add origin git@github.com:${GITHUB_USER}/${TAP_NAME}.git"
echo "   git push -u origin main"
echo "3. Update the formula with actual release SHA256 using:"
echo "   just update-homebrew-formula VERSION"
