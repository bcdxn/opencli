# OpenCLI Specification

[![Go Reference](https://pkg.go.dev/badge/github.com/bcdxn/opencli.svg)](https://pkg.go.dev/github.com/bcdxn/opencli)
![ocli-badge](https://img.shields.io/badge/OpenCLI_Spec-Compliant-brightgreen?link=https%3A%2F%2Fgithub.com%2Fbcdxn%2Fopencli)

Define your CLI in a declarative, language-agnostic document that can be used to generate documentation and boilerplate code.

_Like OpenAPI Spec, but for your CLIs_

## Capabilities

- Unmarshal and validate validate OpenCLI Spec files
- Generate CLI boilerplate code for common CLI frameworks
- Generate CLI documentation in various formats

## Example

```yaml
opencliVersion: 1.0.0-alpha.0

info:
  title: Greet
  summary: A fun CLI defined by OpenCLI Spec
  version: 1.0.0
  binary: greet
      
commands:
  greet {command} <arguments> [flags]:
    group: true
  greet me <name> [flags]:
    summary: "Say hello"
    arguments:
      - name: "name"
        summary: "Your name"
        required: false
        type: "string"
    flags:
      - name: "language"
        summary: "The language of the greeting"
        required: false
        type: "string"
        choices:
          - value: "english"
          - value: "french"
          - value: "german"
```

See a full example of an OpenCLI Document [here](https://github.com/bcdxn/opencli/blob/main/internal/cli.ocs.yaml) - the document that defines the OpenCLI CLI 🤯

## Inspiration

* [OpenAPI Specification](https://swagger.io/specification/)
* Code generation tools like:
  - [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
  - [ogen](https://ogen.dev)
* Stripe's [amazing looking CLI documentation](https://docs.stripe.com/cli)
