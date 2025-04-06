# OpenCLI Specification

Define your Command Line Interface (CLI) in a declarative, language-agnostic document that can be used to generate documentation and boilerplate code.

_Like OpenAPI Spec, but for your CLIs_

---

[![Go Reference](https://pkg.go.dev/badge/github.com/bcdxn/opencli.svg)](https://pkg.go.dev/github.com/bcdxn/opencli)
[![Go Report Card](https://goreportcard.com/badge/github.com/bcdxn/opencli)](https://goreportcard.com/report/github.com/bcdxn/opencli)
![OpenCLI Compliant](https://img.shields.io/badge/OpenCLI_Spec-Compliant-brightgreen?link=https%3A%2F%2Fgithub.com%2Fbcdxn%2Fopencli)

## Table of Contents

- [Overview](#overview)
- [Benefits](#benefits)
- [OpenCLI CLI](#opencli-cli)
- [Examples](#examples)
- [The Spec](#the-spec)
- [Releases](#releases)
- [Inspiration](#inspiration)

## Overview

OpenCLI specification is a document specification that can be used to describe CLIs. Spec-compliant documents are meant to be human-readable but the tooling supports documentation generation in a variety of formats.

## Benefits

- Promote contract first development
- Decouple implementation of commands from the CLI Framework
- Automatically generate documentation your CLI
- Automatically generate CLI framework-specific code

## OpenCLI CLI

Use the CLI to validate specs, generate docs and generate boilerplate code.

- [Markdown Docs](https://github.com/bcdxn/opencli/blob/main/docs/docs.gen.md)

## Examples

### Pleasantries CLI

Let's describe the following CLI

```sh
$ pleasantries greet John --language=english
# hello John
$ pleasantries farewell Jane --language=spanish
# adios Jane
```

The CLI above can be described using an OpenCLI Specification Document in YAML (or JSON):

```yaml
# cli.yaml

opencliVersion: 1.0.0-alpha.6

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

From this example we can generate documentation using the follow command:

```sh
ocli gen docs \
  --spec-file ./cli.osc.yaml \
  --output-dir ./docs \
  --format markdown \
  --dryrun=false
```

You can see the generated documentation [here](https://github.com/bcdxn/opencli/blob/main/examples/markdown-docs/docs.gen.md).

Next, we can generate CLI Framework boilerplate code using the following command:

```sh
ocli gen cli \
  --spec-file ./cli.osc.yaml \
  --output-dir ./internal/cli \
  --framework urfavecli \
  --go-package cli \
  --dryrun=false
```

You can see the generated code [here (go)](https://github.com/bcdxn/opencli/blob/main/examples/urfavecli) and [here (js)](https://github.com/bcdxn/opencli/blob/main/examples/yargs).

### OpenCLI CLI

The OpenCLI CLI uses an OpenCLI Spec and the OpenCLI CLI to generate itself 🤯

- The spec that defines the OpenCLI CLI - [here](https://github.com/bcdxn/opencli/blob/main/internal/cli/cli.ocs.yaml)
- The markdown documentation automatically generated from the spec - [here](https://github.com/bcdxn/opencli/blob/main/docs/docs.gen.md)
- The boilerplate code generated from the spec
  - [generated interface](https://github.com/bcdxn/opencli/blob/main/internal/cli/cli_interface.gen.go)
  - [generated framework boilerplate](https://github.com/bcdxn/opencli/blob/main/internal/cli/cli.gen.go)
  - [generated parameter types](https://github.com/bcdxn/opencli/blob/main/internal/cli/cli_params.gen.go)
 
## The Spec

The full spec is described by JSON Schema - https://github.com/bcdxn/opencli/tree/main/spec

## Releases

Start using OpenCLI Specification Documents to describe your CLIs. Head over to the [releases page](https://github.com/bcdxn/opencli/releases) to download the CLI for your system.

## Inspiration

* [OpenAPI Specification](https://swagger.io/specification/)
* Code generation tools like:
  - [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
  - [ogen](https://ogen.dev)
* Stripe's [amazing looking CLI documentation](https://docs.stripe.com/cli)
