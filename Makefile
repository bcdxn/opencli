.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	go test ./...

version = $(shell git describe --tags HEAD)

.PHONY: examples
examples: generate
	echo "updating examples::::"
	echo "generate example docs"
	go run cmd/ocli/main.go gen docs \
		--spec-file ./examples/cli.ocs.yaml \
		--output-dir ./examples/markdown-docs \
		--format markdown \
		--dryrun=false
	
	echo "generate example yargs CLI"
	go run cmd/ocli/main.go gen cli \
		--spec-file ./examples/cli.ocs.yaml \
		--output-dir ./examples/yargs \
		--framework yargs \
		--module-type cjs \
		--dryrun=false
	
	echo "generate example Urfave/CLI"
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