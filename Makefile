.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	go test ./...

version = $(shell git describe --tags HEAD)

build: generate
	goreleaser release --clean --skip=publish

.PHONY: clean
clean:
	rm -rf dist

.PHONY: all
all: test build