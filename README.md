# OpenCLI Specification

_Like OpenAPI Spec, but for your CLIs_

## Motivation

1. Contract first development - focus on the ergonomics of your CLI before you write any code
2. Remove the redundant work of writing docs pages for your CLI - Generate docs for your CLI automatically and keep them from going stale.
3. Separate your services from the CLI framework - generate the framework-specific boilerplate that invokes your implementation.

## Goals

- [ ] Create a spec
- [ ] Create a JSON-Schema to validate the spec
- [ ] Generate a Markdown documentation file from a spec-compliant file
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