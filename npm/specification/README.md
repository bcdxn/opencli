# @openclidev/specification

The canonical JSON Schema for the [OpenCLI](https://opencli.dev) specification.

Use this package to validate OpenCLI documents or build tools and integrations that consume the specification.

## Installation

```bash
npm install @openclidev/specification
```

## Usage

Import the schema directly:

```typescript
import schema from "@openclidev/specification";
```

Or reference the schema file:

```text
@openclidev/specification/schema/opencli.schema.json
```

The schema is also available via the `schema` export:

```typescript
import schema from "@openclidev/specification/schema";
```

## Source of Truth

The OpenCLI specification and its JSON Schema are maintained in the [OpenCLI repository](https://github.com/bcdxn/opencli). This package is a distribution of the canonical schema for use by tools and implementations across the [ecosystem](https://github.com/bcdxn/opencli#ecosystem).

## License

MIT
