#!/usr/bin/env bash
set -euo pipefail

# Script to sync project files from package.toml
# package.toml is the source of truth for project metadata

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PACKAGE_FILE="package.toml"

if [ ! -f "$PACKAGE_FILE" ]; then
  echo -e "${RED}Error: $PACKAGE_FILE not found${NC}"
  exit 1
fi

echo -e "${GREEN}Syncing project from package.toml${NC}"
echo "===================================="
echo ""

# Parse package.toml using grep and sed
parse_toml() {
  local key=$1
  grep "^$key = " "$PACKAGE_FILE" | sed 's/^[^=]*= *"\(.*\)"$/\1/'
}

PROJECT_NAME=$(parse_toml "name")
MODULE_NAME=$(parse_toml "module")
DESCRIPTION=$(parse_toml "description")
SHORT_DESC=$(parse_toml "short")
VERSION=$(parse_toml "version")
HOMEPAGE=$(parse_toml "homepage")
AUTHOR=$(parse_toml "author")

if [ -z "$PROJECT_NAME" ]; then
  echo -e "${RED}Error: 'name' is required in $PACKAGE_FILE${NC}"
  exit 1
fi

if [ -z "$MODULE_NAME" ]; then
  echo -e "${RED}Error: 'module' is required in $PACKAGE_FILE${NC}"
  exit 1
fi

echo "Project Name: $PROJECT_NAME"
echo "Module Name:  $MODULE_NAME"
echo "Description:  $DESCRIPTION"
echo "Short:  $SHORT_DESC"
echo "Version:      $VERSION"
echo ""

# Store current values to detect changes
CURRENT_MODULE=$(grep "^module " go.mod | awk '{print $2}')
CURRENT_NAME=$(grep "bin/" justfile | head -1 | sed 's|.*bin/\([^ ]*\).*|\1|')

echo -e "${YELLOW}Syncing files...${NC}"

# Update go.mod
if [ "$CURRENT_MODULE" != "$MODULE_NAME" ]; then
  echo "Updating go.mod module name..."
  sed -i "s|module $CURRENT_MODULE|module $MODULE_NAME|g" go.mod

  # Update all Go import paths
  echo "Updating Go import paths..."
  find . -name "*.go" -type f -exec sed -i "s|$CURRENT_MODULE/|$MODULE_NAME/|g" {} \;
fi

# Update config paths
if [ "$CURRENT_NAME" != "$PROJECT_NAME" ]; then
  echo "Updating config paths..."
  sed -i "s|$CURRENT_NAME|$PROJECT_NAME|g" internal/utils/paths.go

  # Update completion examples
  echo "Updating completion examples..."
  sed -i "s|$CURRENT_NAME|$PROJECT_NAME|g" cmd/$MODULE_NAME/completion.go

  # Update justfile
  echo "Updating justfile..."
  sed -i "s|bin/$CURRENT_NAME|bin/$PROJECT_NAME|g" justfile
  sed -i "s|./cmd/$CURRENT_NAME|./cmd/$PROJECT_NAME|g" justfile

  # Rename cmd directory
  if [ -d "cmd/$CURRENT_NAME" ] && [ "$CURRENT_NAME" != "$PROJECT_NAME" ]; then
    echo "Renaming cmd/$CURRENT_NAME to cmd/$PROJECT_NAME..."
    mv "cmd/$CURRENT_NAME" "cmd/$PROJECT_NAME" 2>/dev/null || true
  fi
fi

# Update README description
if [ -n "$DESCRIPTION" ]; then
  echo "Updating README description..."
  # Update the first description line after the title
  sed -i "2s|.*|$DESCRIPTION|" README.md
fi

# Update version in root.go
if [ -n "$VERSION" ]; then
  echo "Updating version in root.go..."
  find cmd -name "root.go" -type f -exec sed -i "s|version = \".*\"|version = \"$VERSION\"|g" {} \;
fi

# Update name in root.go
if [ -n "$PROJECT_NAME" ]; then
  echo "Updating name in root.go..."
  find cmd -name "root.go" -type f -exec sed -i "s|name    = \".*\"|name    = \"$PROJECT_NAME\"|g" {} \;
fi

# Update short description in root.go
if [ -n "$SHORT_DESC" ]; then
  echo "Updating short description in root.go..."
  find cmd -name "root.go" -type f -exec sed -i "s|short   = \".*\"|short   = \"$SHORT_DESC\"|g" {} \;
fi

echo ""
echo -e "${GREEN}✓ Sync complete!${NC}"
echo ""
echo "Files synced from $PACKAGE_FILE"
echo ""
echo "Next steps:"
echo "1. Review the changes: git diff"
echo "2. Build your project: just build"
echo "3. Run tests: just test"
echo ""
