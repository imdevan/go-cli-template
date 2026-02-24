set shell := ["zsh", "-cu"]

build:
	go build -o bin/go-cli-template ./cmd/go-cli-template
	@ls -lh bin/go-cli-template | awk '{print "Build size: " $$5}'

build-run:
	go build -o bin/go-cli-template ./cmd/go-cli-template && ./bin/go-cli-template

watch:
	@rg --files | entr -r sh -c 'sleep 0.5; go build -o bin/go-cli-template ./cmd/go-cli-template'

dev-build:
	go build -gcflags "all=-N -l" -o bin/go-cli-template ./cmd/go-cli-template

cross-platform:
	./scripts/build.sh

build-aur:
	./scripts/build_aur.sh

install:
	install -m 0755 bin/go-cli-template /usr/local/bin/go-cli-template

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
	cd docs && npm install

docs-generate:
	@echo "📝 Generating API documentation from Go packages..."
	./scripts/docs_generate.sh

docs-dev:
	@echo "🚀 Starting documentation development server..."
	@just docs-generate
	cd docs && npm run dev

docs-build:
	@echo "🏗️  Building documentation site..."
	@just docs-generate
	cd docs && npm run build

docs-preview:
	@echo "👀 Previewing built documentation..."
	cd docs && npm run preview

docs-clean:
	@echo "🧹 Cleaning documentation build artifacts..."
	rm -rf docs/dist docs/.astro docs/node_modules docs/src/content/docs/api
