# Contributing

Once a 1.0.0 release candidate for the spec and initial tooling has been released, contributions may be accepted.

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

Before raising pull requests, you must run the commands below to ensure proper regression. The automated workflows will also run the same regression suite.

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

## Style Guidelines

This repository tries to adhere, generally to the [Google Go Style Guide](https://google.github.io/styleguide/go/) and [idiomatic Go](https://go.dev/doc/effective_go) where appropropriate.

## Code Contributions

To obtain OpenCLI Spec locally, fork this repository and work from branches using the naming scheme:

- `feature/...` - branches that provide new feature implementation or enhancements to existing features
- `bugfix/...` - branches that provide fixes to bugs

New code will only be accepted with accompanying unit tests, which will be assessed for quality during the pull request.

## Repo Layout

The code, examples, schema, and web editor are all implemented within this same repository so you can evaluate the ecosystem end to end without context switching. If you are only exploring, start with `README.md` and `examples/`; if you are validating behavior, use the build and test targets in this document.

- `spec/`, `opencli.ocs.yaml`, `spec.schema.json`: Core OpenCLI spec types, canonical example spec, and JSON Schema.
- `cmd/`: Entry points for executables (Cobra CLI and WASM target).
- `internal/`: Internal CLI implementation details and supporting utilities.
- `codec/`: Spec encode/decode logic and fixtures.
- `validate/`: Spec validation package and tests.
- `gen/`: Documentation generators and templates (Markdown/HTML).
- `docs/`: Generated or curated documentation artifacts.
- `examples/`: Example spec files for local testing and demos.
- `web/`: Vite/React web app for editing/previewing specs.

## Feedback

Because this project is open source, feedback is encouraged through GitHub Issues. Please open a bug report for incorrect behavior, broken docs, or unexpected output, and open an enhancement request for proposed features or UX improvements. Clear reproduction steps, expected vs. actual results, environment details, and sample spec files are the most helpful details you can include.
