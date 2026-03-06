set shell := ["zsh", "-cu"]

build:
	go build -o bin/bookmark ./cmd/bookmark
	@size=$(stat -c %s bin/bookmark 2>/dev/null || stat -f %z bin/bookmark 2>/dev/null); \
	echo "Build size: $(awk "BEGIN {printf \"%.2f MB\", $size/1048576}")"

build-run:
	go build -o bin/bookmark ./cmd/bookmark && ./bin/bookmark

watch:
	@rg --files | entr -r sh -c 'sleep 0.5; go build -o bin/bookmark ./cmd/bookmark'

dev-build:
	go build -gcflags "all=-N -l" -o bin/bookmark ./cmd/bookmark

cross-platform:
	./scripts/build.sh

build-aur:
	./scripts/build_aur.sh

install:
	install -m 0755 bin/bookmark /usr/local/bin/bookmark

uninstall:
	rm -f /usr/local/bin/bookmark

test:
	go test ./...

test-verbose:
	go test -v ./...

sync:
	./scripts/sync.sh

clean:
	rm -rf bin

# Documentation tasks
docs-init:
	@echo "📦 Installing documentation dependencies..."
	cd docs && bun install

docs-generate:
	@echo "📝 Generating API documentation from Go packages..."
	./scripts/docs_generate.sh

docs-dev:
	@echo "🚀 Starting documentation development server..."
	@just docs-generate
	cd docs && bun run dev

docs-build:
	@echo "🏗️  Building documentation site..."
	@just docs-generate
	cd docs && NODE_ENV=production bun run build

docs-preview:
	@echo "👀 Previewing built documentation..."
	cd docs && bun run preview

docs-clean:
	@echo "🧹 Cleaning documentation build artifacts..."
	rm -rf docs/dist docs/.astro docs/node_modules docs/src/content/docs/api

# Package distribution tasks
init-homebrew-tap:
	@echo "🍺 Initializing Homebrew tap repository..."
	./scripts/init_homebrew_tap.sh

init-aur-repo:
	@echo "📦 Initializing AUR repository..."
	./scripts/init_aur_repo.sh

update-homebrew-formula VERSION:
	@echo "🍺 Updating Homebrew formula to version {{VERSION}}..."
	./scripts/update_homebrew_formula.sh {{VERSION}}

update-aur-pkgbuild VERSION:
	@echo "📦 Updating AUR PKGBUILD to version {{VERSION}}..."
	./scripts/update_aur_pkgbuild.sh {{VERSION}}

# Git tag management
tag VERSION="":
	./scripts/tag.sh {{VERSION}}

tag-delete VERSION="":
	./scripts/tag_delete.sh {{VERSION}}

tag-list:
	@echo "📋 Available tags:"
	@git tag -l --sort=-v:refname | head -20

# Release management
release VERSION="":
	./scripts/release.sh {{VERSION}}

deploy-aur VERSION="":
	./scripts/deploy_aur.sh {{VERSION}}

deploy-homebrew VERSION="":
	./scripts/deploy_homebrew.sh {{VERSION}}

deploy-all VERSION="":
	./scripts/deploy_all.sh {{VERSION}}

