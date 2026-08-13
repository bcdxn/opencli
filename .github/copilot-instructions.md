# OpenCLI Repository - Copilot Instructions

## Repository Overview

OpenCLI Specification is a declarative, language-agnostic document specification for describing Command Line Interfaces (CLIs), similar to OpenAPI for APIs. It enables contract-first CLI development, CLI documentation generation, and boilerplate CLI code generation for a variety of frameworks similar to oapi-codegen.

## Repository Layout

### Core Packages

**`codec/`** - Marshaling/unmarshaling of OpenCLI spec-compliant documents

- `codec.go` - Core encode/decode logic for OpenCLI documents
- `types.go` - Intermediate raw document types
- Handles YAML/JSON serialization with post-processing

**`spec/`** - OpenCLI Specification types

- `types.go` - Core data structures: Document, Info, CommandItem, ArgumentItem, FlagItem, etc.
- Defines the OpenCLI schema types used throughout the codebase

**`validate/`** - Validation of an OpenCLI spec

- `validate.go` - JSON Schema + logical constraint validation
- `spec.schema.json` - JSON Schema for OpenCLI documents
- Validates both schema compliance and business rules

**`gen/`** - Code & documentation generation from an OpenCLI spec

- `cli.go` - Main CLI generation entry point
- `cli_cobra.go` - Cobra framework code generation
- `cli_urfave_cli.go` - urfave/cli v3 code generation
- `cli_yargs.go` - Yargs (Node.js) code generation
- `docs.go` - Documentation generation
- `docs_md.go` - Markdown docs
- `docs_man.go` - Man page generation
- `docs_html.go` - HTML docs
- `templates/` - Go templates for code generation

### CLI Application

**`cmd/ocli/`** - Main CLI binary

- `main.go` - Entry point, uses internal CLI app

**`cmd/wasm/`** - WebAssembly build for browser editor

- `main.go` - WASM entry point

**`internal/cli/`** - CLI implementation

- `app/` - Cobra CLI command logic implementation
- `gencli/` - generated code (OpenCLI CLI is itself generated from an OpenCLI doc)

### Adapters

**`adapters/ocobra/`** - Cobra adapter

- `ocobra.go` - Generates OpenCLI docs from existing Cobra CLIs
- `parser.go` - Command tree parsing

**`adapters/ourfave/`** - Urfave CLI adapter

- `ourfave.go` - Generates OpenCLI docs from Urfave CLI apps

### Data Structures

**`internal/ds/`** - Internal data structures

- `ordered-map.go` - Ordered map implementation
- `memfs.go` - In-memory filesystem for testing

### Examples & Tests

**`examples/`** - Example spec-compliant documents with their generated docs and code

- `petstore-cli.ocs.yaml` - Petstore API CLI example
- `pleasantries-cli.ocs.yaml` - Simple greetings CLI
- `code/` - Generated code examples for Cobra/Yargs

**`testdata/`** - Test fixtures

- Expected outputs for codec, generation, validation

### Documentation

**`docs/`** - Generated documentation

- `opencli.ocs.md` - OpenCLI CLI `ocli` OpenCLI spec-compliant document, i.e. this document describes the `ocli` CLI

**`web/`** - Next.js documentation site

- `src/` - React/TypeScript frontend
- `public/` - Static assets, WASM runtime
- Live editor at opencli.dev

### Configuration

- `go.mod` - Go module definition
- `go.work` - Workspace configuration
- `Makefile` - Build targets (test, gen-docs, release)
- `opencli.ocs.yaml` - Self-documenting spec for this project
- `spec.schema.json` - JSON Schema that defines the OpenCLI Specification. This is the central pillar of the repository. All of the other packages and documentation serves this specification.

## Key Functionality

1. **Spec Definition**: YAML/JSON documents describing CLI commands, args, flags
2. **Validation**: Schema + logical validation via `validate` package
3. **Code Generation**: Generate boilerplate for Cobra, Urfave CLI, Yargs
4. **Docs Generation**: Markdown, Man pages, HTML documentation
5. **Adapters**: Extract specs from existing CLIs (Cobra, Urfave)
6. **Web Editor**: Browser-based spec editor with live preview

## Development Workflow

- `make test` - Run all tests
- `make gen-docs` - Generate documentation
- `make build-wasm` - Build WASM for web editor
- `go generate ./...` - Generate code

## Important Notes

- The project is self-documenting: `opencli.ocs.yaml` describes the CLI itself
- Generated code uses `gencli` package and directory naming convention to avoid conflicts. We can safely regenerate the code in the repo without fear of clobbering any of the logic
- Web app uses Next.js 16 App Router with static export
