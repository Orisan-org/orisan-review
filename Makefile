.PHONY: all fmt vet test staticcheck build smoke clean

BINARY=orisan-review

all: fmt vet test build smoke

fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./...

staticcheck:
	@if command -v staticcheck >/dev/null 2>&1; then staticcheck ./...; else echo "staticcheck not installed; skipping"; fi

build:
	go build -o bin/$(BINARY) ./cmd/orisan-review

smoke: build
	./bin/$(BINARY) --help
	./bin/$(BINARY) version
	./bin/$(BINARY) list-rules
	./bin/$(BINARY) list-categories

clean:
	rm -rf bin dist coverage.out
