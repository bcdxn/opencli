.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test: generate
	go test ./...

version = $(shell git describe --tags HEAD)

.PHONY: build
build: generate
	go run cmd/cobra/main.go gen docs \
		--out ./docs \
		--format markdown \
		opencli.ocs.yaml

.PHONY: examples
examples: build
	go run cmd/cobra/main.go gen docs \
		--out ./examples/docs \
		--format markdown \
		./examples/petstore-cli.ocs.yaml
	
	go run cmd/cobra/main.go gen docs \
		--out ./examples/docs \
		--format markdown \
		./examples/pleasantries-cli.ocs.yaml

.PHONY: release
release: examples
	echo "building OpenCLI version $(version)::::"
	goreleaser release --clean --skip=publish

.PHONY: clean
clean:
	rm -rf dist

.PHONY: all
all: test release

.PHONY: html-docs
markdown-docs:
	go run cmd/cobra/main.go gen docs \
		--out ./docs \
		--format markdown \
		opencli.ocs.yaml


.PHONY: copy-wasm-exec
copy-wasm-exec:
	cp -f "$$(go env GOROOT)/lib/wasm/wasm_exec.js" web/public/wasm_exec.js

.PHONY: build-wasm
build-wasm: copy-wasm-exec generate
	GOOS=js GOARCH=wasm go build -o web/public/opencli.wasm ./cmd/wasm/main.go

.PHONY: build-ui
build-ui: build-wasm
	cd web && npm ci && npm run build

.PHONY: dev-wasm
dev-wasm: copy-wasm-exec
	GOOS=js GOARCH=wasm go build -o web/public/opencli.wasm ./cmd/wasm/main.go

.PHONY: dev
dev: dev-wasm
	cd web && npm run dev