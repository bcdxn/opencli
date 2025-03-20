# OpenCLI Specification

[![Go Reference](https://pkg.go.dev/badge/github.com/bcdxn/opencli.svg)](https://pkg.go.dev/github.com/bcdxn/opencli)
![ocli-badge](https://img.shields.io/badge/OpenCLI_Spec-Compliant-brightgreen?link=https%3A%2F%2Fgithub.com%2Fbcdxn%2Fopencli)

Define your CLI in a declarative, language-agnostic document that can be used to generate documentation and boilerplate code.

_Like OpenAPI Spec, but for your CLIs_

## Capabilities

- Unmarshal and validate validate OpenCLI Spec files
- Generate CLI boilerplate code for common CLI frameworks
- Generate CLI documentation in various formats

## Inspiration

* [OpenAPI Specification](https://swagger.io/specification/)
* Code generation tools like:
  - [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen)
  - [ogen](https://ogen.dev)
* Stripe's [amazing looking CLI documentation](https://docs.stripe.com/cli)
