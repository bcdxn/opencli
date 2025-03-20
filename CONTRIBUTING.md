# Contributing

Once a release candidate for the spec and initial tooling has been released, contributions may be accepted.

## Repo Layout

- `/spec` - The JSON Schema specification files for the supported versions of OpenCLI.
- `/pkg` - public packages meant to be use used by other projects; API stability follows semantic versioning and semantic import versioning.
- `/cmd` - entrypoints to runnable programs/apps; these programs are typically built and distributed as binaries and should not be imported into other codebases.
- `/internal` - internal packages not meant to be distributed or used imported into other codebases; no API stability is guaranteed.

## Motivation

1. Contract first development - focus on the ergonomics of your CLI before you write any code
2. Remove the redundant work of writing docs pages for your CLI - Generate docs for your CLI automatically and keep them from going stale.
3. Separate your services from the CLI framework - generate the framework-specific boilerplate that invokes your implementation, keeping your code framework-agnostic.

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
