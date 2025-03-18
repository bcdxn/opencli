# OpenCLI Specification

[![Go Reference](https://pkg.go.dev/badge/github.com/bcdxn/opencli.svg)](https://pkg.go.dev/github.com/bcdxn/opencli)
![ocli-badge](https://img.shields.io/badge/OpenCLI_Spec-Compliant-brightgreen?link=https%3A%2F%2Fgithub.com%2Fbcdxn%2Fopencli)

_Like OpenAPI Spec, but for your CLIs_

A declarative specification for your CLI that can be used to generate documentation and boilerplate code

## Capabilities

- Unmarshal and validate validate OpenCLI Spec files
- Generate CLI boilerplate code for common CLI frameworks
- Generate CLI documentation in various formats

## Motivation

1. Contract first development - focus on the ergonomics of your CLI before you write any code
2. Remove the redundant work of writing docs pages for your CLI - Generate docs for your CLI automatically and keep them from going stale.
3. Separate your services from the CLI framework - generate the framework-specific boilerplate that invokes your implementation.

## Goals

- [x] Create a spec
- [x] Create a JSON-Schema to validate the spec
- [x] Generate a Markdown documentation file from a spec-compliant file
- [ ] Generate CLI boilerplate from a spec-compliant file
  - [ ] [urfave/cli](https://github.com/urfave/cli)
- [ ] Generate a static docs site
- [ ] Add support for additional CLI frameworks
  - [ ] [spf13/cobra](https://github.com/spf13/cobra)
  - [ ] [yargs](https://www.npmjs.com/package/yargs)
  - [ ] [oclif](https://www.npmjs.com/package/yargs)
  - [ ] ...

## Inspiration

* [OpenAPI Specification](https://swagger.io/specification/)
* Code generation tools like:
  - [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
  - [ogen](https://ogen.dev)
* Stripe's [amazing looking CLI documentation](https://docs.stripe.com/cli)

## Repo Layout

- `/spec` - The JSON Schema specification files for the supported versions of OpenCLI.
- `/pkg` - public packages meant to be use used by other projects; API stability follows semantic versioning and semantic import versioning.
- `/cmd` - entrypoints to runnable programs/apps; these programs are typically built and distributed as binaries and should not be imported into other codebases.
- `/internal` - internal packages not meant to be distributed or used imported into other codebases; no API stability is guaranteed.

## Testing

#### Run all go:generate directives

Some files rely on copied/generated files that must be in place before tests can run.
Ensure those files are in the right place by executing all `go:generate` directives.

```sh
go generate ./...
```

#### Run all Tests

```sh
go test ./...
```
