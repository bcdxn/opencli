.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	@echo Running tests...
	@go test ./...


version = $(shell git describe --tags HEAD)


.PHONY: gen-docs
gen-docs: generate
	@go run cmd/ocli/main.go gen docs \
		--format markdown \
		--out ./docs \
		ocli.ocs.yaml
	@go run cmd/ocli/main.go gen docs \
		--format html-embed \
		--out ./web/public \
		ocli.ocs.yaml
	@go run cmd/ocli/main.go gen docs \
		--format man \
		--out ./docs \
		ocli.ocs.yaml
	mkdir -p build && mv docs/ocli.ocs.1 build/ocli.1


.PHONY: gen-examples
gen-examples: generate
	@go run cmd/ocli/main.go gen docs \
		--format markdown \
		--out ./examples/docs \
		./examples/petstore-cli.ocs.yaml
	@go run cmd/ocli/main.go gen docs \
		--format markdown \
		--out ./examples/docs \
		./examples/pleasantries-cli.ocs.yaml
	@go run cmd/ocli/main.go gen cli \
		--framework yargs \
		--out ./examples/code/yargs/pleasantries/src \
		./examples/pleasantries-cli.ocs.yaml
	@go run cmd/ocli/main.go gen cli \
		--framework cobra \
		--out ./examples/code/cobra/pleasantries/internal \
		./examples/pleasantries-cli.ocs.yaml
	@go run cmd/ocli/main.go gen cli \
		--framework urfavecli \
		--out ./examples/code/urfavecli/pleasantries/internal \
		./examples/pleasantries-cli.ocs.yaml

.PHONY: gen-ocli
gen-ocli:
	@go run cmd/ocli/main.go gen cli \
		--framework cobra \
		--out ./internal/cli \
		./ocli.ocs.yaml

.PHONY: release
release: gen-docs gen-examples
	@echo "building OpenCLI version $(version)"
	@goreleaser release --clean --skip=publish

.PHONY: clean
clean:
	rm -rf dist

.PHONY: all
all: test release

.PHONY: copy-wasm-exec
copy-wasm-exec:
	cp -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/public/wasm_exec.js

.PHONY: build-wasm
build-wasm: copy-wasm-exec generate
	GOOS=js GOARCH=wasm go build -o web/public/opencli.wasm ./cmd/wasm/main.go

.PHONY: build-ui
build-ui: build-wasm gen-docs
	cp ./examples/petstore-cli.ocs.yaml ./web/public/petstore-cli.ocs.yaml
	cp ./spec.schema.json ./web/src/spec.schema.json
	cd web && npm ci && npm run build

# Warning: local/manual regression testing ahead
.PHONY: gen-scratch-clis
gen-scratch-clis:
	@go run cmd/ocli/main.go gen cli \
		--framework yargs \
		--out ./scratch/yargs/src \
		./examples/petstore-cli.ocs.yaml
	@go run cmd/ocli/main.go gen cli \
		--framework cobra \
		--out ./scratch/cobra/internal \
		./examples/petstore-cli.ocs.yaml
	@go run cmd/ocli/main.go gen cli \
		--framework urfavecli \
		--out ./scratch/urfavecli/internal \
		./examples/petstore-cli.ocs.yaml
