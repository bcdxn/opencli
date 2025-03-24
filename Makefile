.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	go test ./...

version = $(shell git describe --tags HEAD)

build: generate
	go build -o dist/ocli -ldflags "-X github.com/bcdxn/opencli/internal/cli.Version=$(version)"

.PHONY: clean
clean:
	rm -rf dist

.PHONY: all
all: test build