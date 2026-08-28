"use client";

import React, { useState, useCallback } from "react";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { nord } from "react-syntax-highlighter/dist/esm/styles/prism";
import SiteHeader from "../components/SiteHeader";
import "./Docs.css";

// ── Highlighted Code Block (with syntax highlighting) ─────────────────────────

function HighlightedCodeBlock({
  lines,
  language,
}: {
  lines: React.ReactNode[];
  language: string;
}) {
  const [copied, setCopied] = useState(false);

  // Extract plain text for copy
  const getPlainText = useCallback(() => {
    const tempDiv = document.createElement("div");
    lines.forEach((line) => {
      if (typeof line === "string") tempDiv.textContent += line;
      else if (React.isValidElement(line)) {
        const children = line.props; //?.children;
        if (Array.isArray(children)) {
          children.forEach((c) => {
            if (typeof c === "string") tempDiv.textContent += c;
          });
        } else if (typeof children === "string") {
          tempDiv.textContent += children;
        }
      }
    });
    return tempDiv.textContent || "";
  }, [lines]);

  const handleCopy = useCallback(() => {
    navigator.clipboard.writeText(getPlainText()).then(() => {
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    });
  }, [getPlainText]);

  return (
    <div className="guide-code-block">
      <div className="guide-code-block__header">
        <span className="guide-code-block__lang">{language}</span>
        <button
          className={`guide-code-block__copy${copied ? " copied" : ""}`}
          onClick={handleCopy}
        >
          {copied ? "\u2713 Copied" : "Copy"}
        </button>
      </div>
      <div className="body">
        <SyntaxHighlighter
          language={language}
          style={nord}
          customStyle={{ background: "#1A1D24" }}
        >
          {lines.join("\n")}
        </SyntaxHighlighter>
      </div>
    </div>
  );
}

// ── Code Generation (TS) Page ───────────────────────────────────────────────

function GeneratingYargsCodePage() {
  return (
    <>
      <h2 className="guide-section__title">
        Generating A TypeScript CLI From OpenCLI Specs
      </h2>
      <p className="guide-section__subtitle">
        Turn a declarative OpenCLI Specification into production-ready CLI code
        using the <a href="https://yargs.js.org">Yargs framework</a>.
      </p>
      <div className="guide-callout">
        <strong>Tip:</strong> Support is currently available for{" "}
        <ul>
          <li>
            <a href="https://cobra.dev" target="_blank" rel="noreferrer">
              Cobra
            </a>
          </li>
          <li>
            <a href="https://cli.urfave.org" target="_blank" rel="noreferrer">
              urfave/cli
            </a>
          </li>
          <li>
            <a href="https://yargs.js.org" target="_blank" rel="noreferrer">
              Yargs
            </a>
          </li>
        </ul>
        <p>
          Want to add support for your favorite CLI framework? Open an{" "}
          <a
            href="https://github.com/bcdxn/opencli/issues"
            target="_blank"
            rel="noreferrer"
          >
            issue
          </a>{" "}
          or submit a{" "}
          <a
            href="https://github.com/bcdxn/opencli/blob/main/CONTRIBUTING.md"
            target="_blank"
            rel="noreferrer"
          >
            pull request
          </a>
          .
        </p>
      </div>

      {/* Step 1 */}
      <div className="guide-step">
        <div className="guide-step__number">1</div>
        <div className="guide-step__content">
          <h4>Install the CLI</h4>
          <p>
            If you haven't already, install the{" "}
            <span className="guide-inline-code">ocli</span> tool:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[`$ go install github.com/bcdxn/opencli/cmd/ocli@latest`]}
          />
        </div>
      </div>

      {/* The OpenCLI Document */}
      <div className="guide-step">
        <div className="guide-step__number">2</div>
        <div className="guide-step__content">
          <h4>Define your OpenCLI Document</h4>
          <p>
            Every OpenCLI-powered project starts with a spec-compliant YAML (or
            JSON) file. For this walkthrough we'll use the{" "}
            <a
              href="https://github.com/bcdxn/opencli/blob/main/examples/pleasantries-cli.ocs.yaml"
              target="_blank"
              rel="noreferrer"
            >
              pleasantries-cli.ocs.yaml
            </a>{" "}
            example from the{" "}
            <a href="https://github.com/bcdxn/opencli">
              OpenCLI GitHub repository
            </a>
            , a small CLI for greeting and bidding farewell to people by name.
          </p>

          <HighlightedCodeBlock
            language="yaml"
            lines={[
              `opencliVersion: 1.0.0-alpha.14`,
              ``,
              `info:`,
              `  title: Pleasantries`,
              `  summary: A fun CLI to greet or bid farewell`,
              `  version: 1.0.0`,
              `  binary: pleasantries`,
              ``,
              `commands:`,
              `  pleasantries {command} <arguments> [flags]:`,
              `    kind: group`,
              ``,
              `  pleasantries greet <name> [flags]:`,
              `    summary: "Say hello"`,
              `    args:`,
              `      - name: "name"`,
              `        summary: "A name to include in the greeting"`,
              `        required: true`,
              `        type: "string"`,
              `    flags:`,
              `      - name: "language"`,
              `        summary: "The language of the greeting"`,
              `        type: "string"`,
              `        choices:`,
              `          - value: "english"`,
              `          - value: "spanish"`,
              `        default: "english"`,
              ``,
              `  pleasantries farewell <name> [flags]:`,
              `    summary: "Say goodbye"`,
              `    # ... same shape as greet, but for farewells`,
            ]}
          />

          <p>
            You can find the full example document{" "}
            <a
              href="https://github.com/bcdxn/opencli/blob/main/examples/pleasantries-cli.ocs.yaml"
              target="_blank"
              rel="noreferrer"
            >
              here
            </a>{" "}
            and explore the complete specification schema at{" "}
            <a href="/specification">opencli.dev/specification</a>.
          </p>
        </div>
      </div>

      {/* Step 3: Initialize the Project */}
      <div className="guide-step">
        <div className="guide-step__number">3</div>
        <div className="guide-step__content">
          <h4>Initialize the Project</h4>
          <p>
            Set up a fresh Node.js project and pull in the pleasantries spec:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ mkdir pleasantries && cd pleasantries`,
              `$ npm init -y`,
              `# pull in the full example ocs file (or use your own)`,
              `$ curl -O https://raw.githubusercontent.com/bcdxn/opencli/refs/heads/main/examples/pleasantries-cli.ocs.yaml`,
            ]}
          />

          <p>Then install the runtime and dev dependencies:</p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ npm i yargs command-line-usage`,
              `$ npm i -D typescript @types/yargs @types/node @types/command-line-usage`,
            ]}
          />

          <div className="guide-callout">
            <p>
              The generated code uses{" "}
              <span className="guide-inline-code">command-line-usage</span> to
              render help and usage output, so it's a required dependency.
            </p>
          </div>

          <p>That's it for setup — one spec file, one package.</p>
        </div>
      </div>

      {/* Step 4: Generate Boilerplate Code */}
      <div className="guide-step">
        <div className="guide-step__number">4</div>
        <div className="guide-step__content">
          <h4>Generate Boilerplate Code</h4>
          <p>
            A single <span className="guide-inline-code">ocli gen cli</span>{" "}
            command produces all the scaffolding:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ ocli gen cli \\`,
              `  --framework yargs \\`,
              `  --out ./src \\`,
              `  ./pleasantries-cli.ocs.yaml`,
              `# → Reading spec:        ./pleasantries-cli.ocs.yaml`,
              `# → Generating CLI code: framework=yargs, output=./src`,
              `# ✓ CLI Code written to: ./src`,
            ]}
          />

          <p>
            All generated code is encapsulated in the{" "}
            <span className="guide-inline-code">gencli</span> directory. Each
            command gets its own file, plus supporting files for bootstrapping,
            error handling, and help rendering:
          </p>

          <HighlightedCodeBlock
            language="plain"
            lines={[
              `package.json`,
              `pleasantries-cli.ocs.yaml`,
              `src/`,
              `\u2514\u2500\u2500 gencli/`,
              `    \u251c\u2500\u2500 actions.ts            ActionsInterface & command signatures`,
              `    \u251c\u2500\u2500 cmd-pleasantries-greet.ts     Generated yargs command definitions`,
              `    \u251c\u2500\u2500 cmd-pleasantries-farewell.ts  Generated yargs command definitions`,
              `    \u251c\u2500\u2500 errors.ts             CLI error types & exit codes`,
              `    \u251c\u2500\u2500 help.ts               Default help/usage rendering`,
              `    \u251c\u2500\u2500 params.ts             Command args, flags & choice enums`,
              `    \u251c\u2500\u2500 types.ts              Shared command metadata types`,
              `    \u2514\u2500\u2500 run.ts                CLI entry point (run function)`,
            ]}
          />

          <div className="guide-callout">
            <p>
              <strong>Key insight:</strong> the generated code defines an{" "}
              <span className="guide-inline-code">ActionsInterface</span>. The
              interface creates a contract that maps methods one-to-one with
              every command in your spec along with some convenience methods.
              Your job is simply to implement that contract and those methods.
            </p>
          </div>

          <p>
            Let's take a look at the all-important{" "}
            <span className="guide-inline-code">src/gencli/actions.ts</span>. It
            defines one method per command, plus helpers for help and usage:
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `// src/gencli/actions.ts`,
              `export interface ActionsInterface {`,
              `  PleasantriesGreet(args: PleasantriesGreetArgs, flags: PleasantriesGreetFlags): Promise<void>;`,
              `  PleasantriesFarewell(args: PleasantriesFarewellArgs, flags: PleasantriesFarewellFlags): Promise<void>;`,
              `  help(cmd: CommandPrintData): void;`,
              `  usage(cmd: CommandPrintData): void;`,
              `}`,
            ]}
          />

          <div className="guide-callout">
            <p>
              Look at your generated{" "}
              <span className="guide-inline-code">src/gencli/actions.ts</span>{" "}
              to see the full interface we'll need to implement.
            </p>
          </div>

          <p>
            Notice that the methods we need to implement have no
            framework-dependencies injected. We could reuse our same{" "}
            <span className="guide-inline-code">ActionsInterface</span>{" "}
            implementation for multiple frameworks within the same language (or
            port it across languages entirely).
          </p>

          <p>
            The generated types for{" "}
            <span className="guide-inline-code">args</span> and{" "}
            <span className="guide-inline-code">flags</span> are strongly typed,
            so you get compile-time safety — no more typos in flag names or
            mismatched types. Flags with choices even become enums:
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `// src/gencli/params.ts`,
              `export enum PleasantriesGreetLanguage {`,
              `  ENGLISH = "english",`,
              `  SPANISH = "spanish",`,
              `}`,
              ``,
              `export interface PleasantriesGreetArgs {`,
              `  name: string;`,
              `}`,
              ``,
              `export interface PleasantriesGreetFlags {`,
              `  language?: PleasantriesGreetLanguage | undefined;`,
              `}`,
            ]}
          />

          <p>
            Next we can take a look at the generated command files, like{" "}
            <span className="guide-inline-code">
              src/gencli/cmd-pleasantries-greet.ts
            </span>
            . Each generated command file adapts our ActionsInterface methods,
            handling the framework specifics of parsing args and flags and
            passing them to our framework-<i>agnostic</i> implementations. If
            you're interested, you can look at a generated file to see how the{" "}
            <span className="guide-inline-code">handler</span> delegates to the
            corresponding method on your class implementing the{" "}
            <span className="guide-inline-code">ActionsInterface</span>. But in
            general you can treat these generated command files as black boxes.
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `handler: async (argv) => {`,
              `  const cmdArgs: PleasantriesGreetArgs = { name: argv.name };`,
              `  const cmdFlags: PleasantriesGreetFlags = { language: argv.language as PleasantriesGreetLanguage };`,
              `  return actions.PleasantriesGreet(cmdArgs, cmdFlags);`,
              `},`,
            ]}
          />
        </div>
      </div>

      {/* Step 5: Implement the Actions Interface */}
      <div className="guide-step">
        <div className="guide-step__number">5</div>
        <div className="guide-step__content">
          <h4>Implement the Actions Interface</h4>
          <p>
            This is where you write your actual business logic. Create a class
            that satisfies{" "}
            <span className="guide-inline-code">ActionsInterface</span>. The
            pattern feels familiar if you've used{" "}
            <a
              href="https://github.com/oapi-codegen/oapi-codegen"
              target="_blank"
              rel="noreferrer"
            >
              oapi-codegen
            </a>{" "}
            with OpenAPI specs.
          </p>

          <p>
            Start by creating a new file for your implementation to keep it
            separate from the generated code in the{" "}
            <span className="guide-inline-code">gencli</span> package:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[`$ touch ./src/actions.ts`]}
          />

          <p>
            Define your <span className="guide-inline-code">Actions</span>{" "}
            class:
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `// src/actions.ts`,
              `import { ActionsInterface } from "./gencli/actions";`,
              `import { CommandPrintData } from "./gencli/types";`,
              `import { PleasantriesGreetArgs, PleasantriesGreetFlags, ... } from "./gencli/params";`,
              `import { defaultHelpFn, defaultUsageFn } from "./gencli/help";`,
              ``,
              `export class Actions implements ActionsInterface {`,
              `  // ... implement each command method below`,
              `}`,
            ]}
          />

          <p>
            Now implement each method to fulfill the interface. For
            demonstration we'll keep the bodies simple — in a real project this
            is where you'd call your API, hit a database, or orchestrate
            whatever your CLI is designed to do:
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `async PleasantriesGreet(args: PleasantriesGreetArgs, flags: PleasantriesGreetFlags): Promise<void> {`,
              `  if (flags.language == "english") {`,
              `    console.log("hello", args.name);`,
              `  } else {`,
              `    console.log("hola", args.name);`,
              `  }`,
              `}`,
              ``,
              `async PleasantriesFarewell(args: PleasantriesFarewellArgs, flags: PleasantriesFarewellFlags): Promise<void> {`,
              `  if (flags.language == "english") {`,
              `    console.log("good bye", args.name);`,
              `  } else {`,
              `    console.log("adios", args.name);`,
              `  }`,
              `}`,
            ]}
          />

          <div className="guide-callout">
            <p>
              You can download a full example implementation{" "}
              <a href="/assets/code/actions.ts">here</a>.
            </p>
          </div>

          <p>
            If your OpenCLI document declares root-level{" "}
            <span className="guide-inline-code">global</span> flags, they're not
            passed to action methods either. Yargs has no context object, so
            codegen instead exports a pair of accessors from{" "}
            <span className="guide-inline-code">params.ts</span>: the generated
            handler calls{" "}
            <span className="guide-inline-code">setGlobalFlags(...)</span>{" "}
            immediately before invoking your action, and you read them with{" "}
            <span className="guide-inline-code">getGlobalFlags()</span> inside:
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `import { getGlobalFlags } from "./gencli/params";`,
              ``,
              `// inside your Actions class...`,
              `async PleasantriesGreet(args: PleasantriesGreetArgs, flags: PleasantriesGreetFlags): Promise<void> {`,
              `  const name = args.name; // positional arg — arrives as a parameter`,
              ``,
              `  // Global (root-level) flags come from the module accessor, not the method signature.`,
              `  const global = getGlobalFlags();`,
              `  if (global.debug) { // e.g., for a root-level --debug flag declared in your spec`,
              `    console.error(\`greeting \${name} in debug mode\`);`,
              `  }`,
              `}`,
            ]}
          />

          <div className="guide-callout">
            <p>
              These types and accessors are only emitted when your document
              declares root-level flags. Yargs parameter fields are typed{" "}
              <span className="guide-inline-code">T | undefined</span> even for
              required args, so guard values with{" "}
              <span className="guide-inline-code">??</span> or an explicit check
              rather than assuming presence (e.g.,{" "}
              <span className="guide-inline-code">global.timeout ?? 30</span>
              ). The accessor is a module-level singleton set per invocation by
              the generated handler — safe under normal sequential CLI use.
            </p>
          </div>

          <p>
            Finally, wire up the helper methods using sensible defaults provided
            by the generated code (or replace them with custom implementations
            if you need tailored behavior):
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `help(cmd: CommandPrintData): void {`,
              `  defaultHelpFn(cmd);`,
              `}`,
              ``,
              `usage(cmd: CommandPrintData): void {`,
              `  defaultUsageFn(cmd);`,
              `}`,
            ]}
          />

          <div className="guide-callout">
            <p>
              <strong>Benefits of this approach:</strong> your spec is the
              contract, your business logic has zero dependencies on any CLI
              framework, and documentation stays in sync with the OpenCLI Spec
              document as the source of truth.
            </p>
          </div>
        </div>
      </div>

      {/* Step 6: Wire Up the Entry Point */}
      <div className="guide-step">
        <div className="guide-step__number">6</div>
        <div className="guide-step__content">
          <h4>Wire Up the Entry Point</h4>
          <p>
            The final piece is a minimal{" "}
            <span className="guide-inline-code">src/index.ts</span>:
          </p>

          <HighlightedCodeBlock
            language="ts"
            lines={[
              `#!/usr/bin/env node`,
              `import yargs from "yargs";`,
              `import { hideBin } from "yargs/helpers";`,
              `import { run } from "./gencli/run";`,
              `import { Actions } from "./actions";`,
              ``,
              `async function main() {`,
              `  const actions = new Actions();`,
              `  await run(yargs(hideBin(process.argv)), actions);`,
              `}`,
              ``,
              `main().catch((err) => {`,
              `  console.error(err);`,
              `  process.exit(1);`,
              `});`,
            ]}
          />

          <p>
            Just a handful of lines of substance, and critically — no framework
            dependencies in your user-land code.
          </p>
        </div>
      </div>

      {/* Step 7: Try it Out */}
      <div className="guide-step">
        <div className="guide-step__number">7</div>
        <div className="guide-step__content">
          <h4>Try It Out</h4>
          <p>That's the entire application. Let's build and run it:</p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ npx tsc   # compiles src/ → dist/`,
              ``,
              `$ node dist/index.js greet John --language spanish`,
              `# hola John`,
            ]}
          />

          <HighlightedCodeBlock
            language="sh"
            lines={[`$ node dist/index.js farewell Alice`, `# good bye Alice`]}
          />

          <p>
            A fully functional CLI with zero framework coupling in your business
            logic. The spec defined the interface,{" "}
            <span className="guide-inline-code">ocli</span> generated the
            scaffolding, and you implemented the business logic.
          </p>
        </div>
      </div>

      {/* Next steps */}
      <div className="whats-next">
        <h3>What's next?</h3>
        <ul>
          <li>
            Explore the full{" "}
            <a href="/specification">OpenCLI Specification reference</a> for a
            deeper understanding of the spec.
          </li>
          <li>
            Browse{" "}
            <a
              href="https://github.com/bcdxn/opencli/tree/main/examples"
              target="_blank"
              rel="noreferrer"
            >
              example specs
            </a>{" "}
            in the repository
          </li>
        </ul>
      </div>
    </>
  );
}

// ── Main Component ────────────────────────────────────────────────────────────

export default function GuidePage() {
  return (
    <div className="guide-page">
      <SiteHeader />
      <div className="guide-layout">
        {/* Left nav */}
        <nav className="guide-nav" aria-label="Guide navigation">
          <p className="guide-nav__heading">Guide</p>
          <ul className="guide-nav__list">
            <li key="code-generation-yargs">
              <a href="/docs/getting-started" className="guide-nav__link">
                Getting Started
              </a>
              <a href="/docs/markdown-docs" className="guide-nav__link">
                Markdown Docs
              </a>
              <a href="/docs/html-docs" className="guide-nav__link">
                HTML Docs
              </a>
              <a href="/docs/man-pages" className="guide-nav__link">
                Man Pages
              </a>
              <a href="/docs/code-generation-go" className="guide-nav__link">
                Code Generation (Go)
              </a>
              <a
                href="/docs/code-generation-yargs"
                className="guide-nav__link is-active"
              >
                Code Generation (TS)
              </a>
            </li>
          </ul>
        </nav>

        {/* Main content */}
        <main className="guide-main">
          <GeneratingYargsCodePage key="code-generation-yargs" />
        </main>
      </div>
    </div>
  );
}
