set shell := ["zsh", "-cu"]

build:
	go build -o bin/go-cli-template ./cmd/go-cli-template

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
