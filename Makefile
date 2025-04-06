.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	go test ./...

version = $(shell git describe --tags HEAD)

.PHONY: build
build: generate
	go run cmd/ocli/main.go gen docs \
		--spec-file ./internal/cli/cli.ocs.yaml \
		--output-dir ./docs \
		--format markdown \
		--dryrun=false
	
	go run cmd/ocli/main.go gen cli \
		--spec-file ./internal/cli/cli.ocs.yaml \
		--output-dir ./internal/cli \
		--framework urfavecli \
		--go-package cli \
		--dryrun=false

.PHONY: examples
examples: build
	go run cmd/ocli/main.go gen docs \
		--spec-file ./examples/cli.ocs.yaml \
		--output-dir ./examples/markdown-docs \
		--format markdown \
		--dryrun=false
	
	go run cmd/ocli/main.go gen cli \
		--spec-file ./examples/cli.ocs.yaml \
		--output-dir ./examples/yargs \
		--framework yargs \
		--module-type cjs \
		--dryrun=false
	
	go run cmd/ocli/main.go gen cli \
		--spec-file ./examples/cli.ocs.yaml \
		--output-dir ./examples/urfavecli/cli \
		--framework urfavecli \
		--go-package cli \
		--dryrun=false

.PHONY: release
release: examples
	
	
	echo "building OpenCLI version $(version)::::"
	goreleaser release --clean --skip=publish

.PHONY: clean
clean:
	rm -rf dist

.PHONY: all
all: test release