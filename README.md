# OpenCLI Specification

_Like OpenAPI Spec, but for your CLIs_

## Motivation

1. Contract first development - focus on the ergonomics of your CLI before you write any code
2. Remove the redundant work of writing docs pages for your CLI - Generate docs for your CLI automatically and keep them from going stale.
3. Separate your services from the CLI framework - generate the framework-specific boilerplate that invokes your implementation.

## Goals

- [x] Create a spec
- [x] Create a JSON-Schema to validate the spec
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

## Repo Layout

- `/spec` - The JSON Schema specification files for the supported version of OpenCLI.
- `/pkg` - public packages meant to be use used by other projects; API stability follows semantic versioning and semantic import versioning.
- `/cmd` - entrypoints to runnable programs/apps; these programs are typically built and distributed as binaries and should not be imported into other codebases.
- `/internal` - internal packages not meant to be distributed or used imported into other codebases; no API stability is guaranteed.

### `/pkg` Public Packages

#### `opencli/validator`

Offers the ability to validate OpenCLI Spec files.

### Usage

```go
package main

import (
  "fmt"

  "github.com/bcdxn/opencli/pkg/validator"
)
  

func main() {
  // check what versions are available
  versions := validator.Versions()
  fmt.Println(versions) // [1.0.0-alpha.0]

  // validate an document to see if it is compliant with the OpenCLI Specification
  err := validator.ValidateDocument([]byte(`
    {
      "opencliVersion": "1.0.0-alpha.0",
      "info": {
        "binary": "test",
        "title": "Test OpenCLI Specification",
        "version": "1.0.0"
      },
      "commands": {}
    }
  `))
  if err != nil {
    panic(err)
  }
}

```

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
