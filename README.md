# OpenCLI Specification

[![Go Reference](https://pkg.go.dev/badge/github.com/bcdxn/opencli.svg)](https://pkg.go.dev/github.com/bcdxn/opencli)
![ocli-badge](https://img.shields.io/badge/OpenCLI_Spec-Compliant-brightgreen?link=https%3A%2F%2Fgithub.com%2Fbcdxn%2Fopencli)

Define your CLI in a declarative, language-agnostic document that can be used to generate documentation and boilerplate code.

_Like OpenAPI Spec, but for your CLIs_

## OpenCLI Specs Benefits

- Promote contract first development
- Decouple implementation of commands from the CLI Framework
- Automatically generate documentation your CLI
- Automatically generate CLI framework-specific code

## Example

Let's describe the following CLI

```sh
$ pleasantries greet John --language=english
# hello John
$ pleasantries farewell Jane --language=spanish
# adios Jane
```

The CLI above can be described using an OpenCLI Specification Document like:

```yaml
opencliVersion: 1.0.0-alpha.3

info:
  title: Pleasantries
  summary: A fun CLI to greet or bid farewell
  version: 1.0.0
  binary: pleasantries
      
commands:
  pleasantries {command} <name> [flags]:
    group: true

  pleasantries greet <name> [flags]:
    summary: "Say hello"
    arguments:
      - name: "name"
        summary: "A name to include the greeting"
        required: true
        type: "string"
    flags:
      - name: "language"
        summary: "The language of the greeting"
        type: "string"
        choices:
          - value: "english"
          - value: "spanish"
        default: "english"

  pleasantries farewell <name> [flags]:
    summary: "Say goodbye"
    arguments:
      - name: "name"
        summary: "A name to include in the farewell"
        required: true
        type: "string"
    flags:
      - name: "language"
        summary: "The language of the greeting"
        type: "string"
        choices:
          - value: "english"
          - value: "spanish"
        default: "english"
```

## OpenCLI CLI 🤯

See a full example of an OpenCLI Document in action.

- The spec that defines the OpenCLI CLI [here](https://github.com/bcdxn/opencli/blob/main/internal/cli/cli.ocs.yaml)
- The markdown documentation automatically generated from the spec [here](https://github.com/bcdxn/opencli/blob/main/docs/docs.gen.md)
- The CLI interface and boilerplate code generated from the spec [here](https://github.com/bcdxn/opencli/blob/main/internal/cli/cli_interface.gen.go) and [here](https://github.com/bcdxn/opencli/blob/main/internal/cli/cli.gen.go)

## Releases

Start using OpenCLI Specification Documents to describe your CLIs. Head over to the [releases page](https://github.com/bcdxn/opencli/releases) to download the CLI for your system.

## Inspiration

* [OpenAPI Specification](https://swagger.io/specification/)
* Code generation tools like:
  - [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
  - [ogen](https://ogen.dev)
* Stripe's [amazing looking CLI documentation](https://docs.stripe.com/cli)
