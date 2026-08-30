# Implementing an `ActionsInterface` command function

When you generate a CLI from an OpenCLI spec, codegen emits an `ActionsInterface` whose methods are the real work of your program. This document shows how to read the three kinds of data available inside one such implementation:

- **positional args**: passed as a generated `<Method>Args` struct parameter (method name + `Args`)
- **command flags**: passed as a generated `<Method>Flags` struct parameter (method name + `Flags`)
- **global (root-level) flags**: _not_ a method parameter; retrieved from the
  `context.Context` in Go frameworks, or via module accessors in Yargs

The examples below are based on this spec fragment:

```yaml
info:
  binary: example
global:
  flags:
    - name: debug
      aliases: [v]
      type: boolean
    - name: timeout
      type: integer
      default: 30
commands:
  # Root group; the binary segment is part of every method name.
  example {command} [flags]:
    kind: group

  example cmd <recipient> --count <integer> [flags]:
    args:
      - name: recipient
        type: string
    flags:
      - name: count
        type: integer
```

Each leaf gets one method named by PascalCase-joining every path segment _including the binary_ — so `example cmd` becomes `ExampleCmd`, with parameter types `ExampleCmdArgs` / `ExampleCmdFlags`. In Go (identical for Cobra and urfave/cli v3):

```go
type ActionsInterface interface {
	ExampleCmd(ctx context.Context, args ExampleCmdArgs, flags ExampleCmdFlags) error
	// other actions...
}
```

and its parameter types:

```go
type ExampleCmdArgs struct {
	Recipient string
}

type ExampleCmdFlags struct {
	Count int64
}

type GlobalFlags struct {
	Debug        bool
	Timeout      int64
}
```

## Go (Cobra & urfave/cli v3) — global flags on the context

The generated command handler builds a `GlobalFlags` value and injects it into the context before calling your action, so you never receive globals as an argument. Instead, retrieve them with the exported helper from the same package:

```go
// ExampleCmd implements gencli.ActionsInterface for the "example cmd" leaf.
func (a *app) ExampleCmd(
  ctx context.Context,
  args gencli.ExampleCmdArgs,
  flags gencli.ExampleCmdFlags,
) error {
	// Positional arguments arrive as a struct parameter.
	recipient := args.Recipient
	// Command-level flags also arrive as a struct parameter.
	count := flags.Count
	// Global flags are NOT a method parameter — read them from ctx.
	global := gencli.GlobalFlagsFromContext(ctx)
	if global.Debug {
		fmt.Fprintf(a.IOStreams().ErrOut(), "sending %d message(s) to %s\n", count, recipient)
	}
	timeout := time.Duration(global.Timeout) * time.Second

	// do stuff...
	return nil
}
```

Notes:

- `GlobalFlagsFromContext` returns zero values if no globals were set (e.g. when calling an action directly from a test without wrapping the context). Use `gencli.WithGlobalFlags(ctx, g)` in tests when you want to reproduce the context setup performed by the generated runtime handler.
- The parameter types and context helpers all live in the _generated_ package (in its params file). Your implementation lives elsewhere, so everything is referenced through that package's import name — `gencli.ExampleCmdArgs`, `gencli.GlobalFlagsFromContext(ctx)`, etc.

## TypeScript (Yargs) — module-level accessors

Yargs has no context object, so codegen uses an exported pair of accessor functions instead. The generated handler calls `setGlobalFlags(...)` immediately before invoking your action; inside the action you read them with `getGlobalFlags()`:

```ts
// actions.ts and params.ts are generated side by side; import both.
import type { ActionsInterface } from "./actions";
import {
  getGlobalFlags,
  type ExampleCmdArgs,
  type ExampleCmdFlags,
} from "./params";

export class App implements ActionsInterface {
  async ExampleCmd(
    args: ExampleCmdArgs,
    flags: ExampleCmdFlags,
  ): Promise<void> {
    // Positional arguments arrive as a struct parameter (fields may be undefined).
    const recipient = args.recipient;
    // Command-level flags also arrive as a struct parameter.
    const count = flags.count ?? 1; // number | undefined — apply your own fallback
    // Global flags come from the module accessor, not the method signature.
    const global = getGlobalFlags();
    if (global.debug) {
      console.error(`sending ${count} message(s) to ${recipient}`);
    }
    const timeoutMs = (global.timeout ?? 30) * 1000;

    // do stuff...
  }
}
```

Notes:

- Optional Yargs parameter fields are typed `T | undefined`, so guard optional
  values with `??` or an explicit check.
- The accessor is a module-level singleton set per invocation by the generated
  handler — it's safe under normal sequential CLI use (the same pattern codegen
  already uses for config loading).

## Quick reference

| Data            | Go                                     | JS/TS                              |
| --------------- | -------------------------------------- | ---------------------------------- |
| Positional args | `args <Method>Args` method parameter   | `args: <Method>Args` parameter     |
| Command flags   | `flags <Method>Flags` method parameter | `flags: <Method>Flags` parameter   |
| Global flags    | `gencli.GlobalFlagsFromContext(ctx)`   | `getGlobalFlags()` from `./params` |
