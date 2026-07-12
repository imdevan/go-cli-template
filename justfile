# Justfile setup
# ================================================================================
set shell := ["zsh", "-cu"]

PACKAGE := "go-cli-template"
PACKAGE_BIN := "./bin/" + PACKAGE
PACKAGE_CMD := "./cmd/" + PACKAGE

cli_docs := "bin/go-cli-docs"
cli_pack := "bin/go-cli-package"

_install-docs:
	@if [ ! -f ./bin/go-cli-docs ]; then \
		echo "📥 Installing go-cli-docs..."; \
		GOBIN="{{justfile_directory()}}/bin" "$(go env GOROOT)/bin/go" install github.com/imdevan/go-cli-docs/cmd/go-cli-docs@latest; \
	fi

_install-pack:
	@if [ ! -f ./bin/go-cli-package ]; then \
		echo "📥 Installing go-cli-package..."; \
		GOBIN="{{justfile_directory()}}/bin" "$(go env GOROOT)/bin/go" install github.com/imdevan/go-cli-package/cmd/go-cli-package@latest; \
	fi



# Build
# ================================================================================

build:
	go build -o {{PACKAGE_BIN}} {{PACKAGE_CMD}}
	@size=$(stat -c %s {{PACKAGE_BIN}} 2>/dev/null || stat -f %z {{cli_pack}} 2>/dev/null); \
	echo "Build size: $(awk "BEGIN {printf \"%.2f MB\", $size/1048576}")"

build-run:
	go build -o {{PACKAGE_BIN}} {{PACKAGE_CMD}} && {{PACKAGE_BIN}}

watch:
	@rg --files | entr -r sh -c 'sleep 0.5; go build -o {{PACKAGE_BIN}} {{PACKAGE_CMD}}'

dev-build:
	go build -gcflags "all=-N -l" -o {{PACKAGE_BIN}} {{PACKAGE_CMD}}

# Install local build globally
install:
	install -m 0755 {{PACKAGE_BIN}} /usr/local/bin/{{PACKAGE}}

# Uninstall local build globally
uninstall:
	rm -f /usr/local/bin/{{PACKAGE}}

test:
	go test ./...

test-verbose:
	go test -v ./...

clean:
	rm -rf bin

# Documentation
# ================================================================================

docs-init args="": _install-docs
	{{cli_docs}} init {{args}}

docs-generate args="":  _install-docs
	{{cli_docs}} generate {{args}}

docs-dev args="": _install-docs
	{{cli_docs}} watch {{args}} & cd docs && bun install && bun run dev

docs-build: docs-generate
	@echo "🏗️  Building documentation site..."
	cd docs && NODE_ENV=production bun run build

docs-preview:
	@echo "👀 Previewing built documentation..."
	cd docs && bun run preview

docs-clean:
	@echo "🧹 Cleaning documentation build artifacts..."
	rm -rf docs/dist docs/.astro docs/node_modules docs/src/content/docs/api

# Package management
# ================================================================================

# Rename template based on internal/package/package.toml or pass --package 
sync args="":
  {{cli_pack}} sync {{args}}

# Github tag management

tag-list: build _install-pack
	{{cli_pack}} tag list

tag version="": build _install-pack
	{{cli_pack}} tag create {{version}}

tag-delete version="": build _install-pack
	{{cli_pack}} tag delete {{version}}

# Github release (calls tag)

release version="": build _install-pack
	{{cli_pack}} release {{version}}

# Pipeline init

init-homebrew: build _install-pack
	{{cli_pack}} init homebrew

init-aur: build _install-pack
	{{cli_pack}} init aur

init: build _install-pack
	{{cli_pack}} init all

# Package updates

update-homebrew version="": build _install-pack
	{{cli_pack}} update homebrew {{version}}

update-aur version="": build _install-pack
	{{cli_pack}} update aur {{version}}

update version="": build _install-pack
	{{cli_pack}} update all {{version}}


# Publish

publish-homebrew version="": build
	{{cli_pack}} publish homebrew {{version}}

publish-aur version="": build
	{{cli_pack}} publish aur {{version}}

publish version="": build
	{{cli_pack}} publish all {{version}}


