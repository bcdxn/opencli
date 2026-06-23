# Contributing

Once a release candidate for the spec and initial tooling has been released, contributions may be accepted.

## Repo Layout

- `spec/`, `opencli.ocs.yaml`, `spec.schema.json`: Core OpenCLI spec types, canonical example spec, and JSON Schema.
- `cmd/`: Entry points for executables (Cobra CLI and WASM target).
- `internal/`: Internal CLI implementation details and supporting utilities.
- `codec/`: Spec encode/decode logic and fixtures.
- `validate/`: Spec validation package and tests.
- `gen/`: Documentation generators and templates (Markdown/HTML).
- `docs/`: Generated or curated documentation artifacts.
- `examples/`: Example spec files for local testing and demos.
- `web/`: Vite/React web app for editing/previewing specs.

## Motivation

1. Contract first development - focus on the ergonomics of your CLI before you write any code
2. Remove the redundant work of writing docs pages for your CLI - Generate docs for your CLI automatically and keep them from going stale.
3. Separate your services from the CLI framework - generate the framework-specific boilerplate that invokes your implementation, keeping your code framework-agnostic.
4. Help LLMs to quickly understand CLI capabilities with less context usage than simply parsing the code

## Goals

- [x] Create a spec
- [x] Create a JSON-Schema to validate the spec
- [x] Generate a documentation from a spec-compliant file
  - [x] Markdown documentation
  - [x] HTML documentation
- [ ] Generate CLI boilerplate from a spec-compliant file
  - [ ] [urfave/cli](https://cli.urfave.org)
  - [ ] [cobra](https://cobra.dev)
  - [ ] [yargs](https://yargs.js.org)
- [x] Generate a static docs site

## Testing and Building

### Makefile

Some files rely on copied/generated files that must be in place before tests can run.
Ensure those prerequesites are taken care of by using the Makefile targets

```sh
go generate ./...
```

#### Run Tests

```sh
make test
```

#### Build

```sh
make build
```

```sh
make build-ui
```

#### Release

```sh
make release
```
