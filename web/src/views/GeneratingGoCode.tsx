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

// ── Code Generation (Go) Page ──────────────────────────────────────────────────

function GeneratingGoCodePage() {
  return (
    <>
      <h2 className="guide-section__title">
        Generating A Go CLI From OpenCLI Specs
      </h2>
      <p className="guide-section__subtitle">
        Turn a declarative OpenCLI Specification into framework-specific,
        production-ready CLI code — then implement only the business logic.
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
            JSON) file. For this walkthrough we'll use the
          </p>
          <a
            href="https://github.com/bcdxn/opencli/blob/main/examples/petstore-cli.ocs.yaml"
            target="_blank"
            rel="noreferrer"
          >
            petstore-cli.ocs.yaml
          </a>{" "}
          example from the{" "}
          <a
            href="https://github.com/bcdxn/opencli"
            target="_blank"
            rel="noreferrer"
          >
            OpenCLI GitHub repository
          </a>
          , modeled after the classic Swagger petstore API so the concepts feel
          familiar.
          <p>
            The document describes a CLI for managing pets, orders, and users.
            Here's a portion of what it looks like:
          </p>
          <HighlightedCodeBlock
            language="yaml"
            lines={[
              `opencliVersion: 1.0.0-alpha.12`,
              ``,
              `info:`,
              `  title: PetStore CLI`,
              `  summary: An example CLI Document describing operations a petstore CLI may provide.`,
              `  # ...`,
              ``,
              `commands:`,
              `  petstore pet add [flags]:`,
              `    summary: Add a new pet to the store`,
              `    args:`,
              `      - name: path-to-req-body`,
              `        type: string`,
              `        summary: The path to a JSON file containing the new pet payload`,
              `        required: false`,
              `    flags:`,
              `      - name: name`,
              `        aliases:`,
              `          - n`,
              `        type: string`,
              `        summary: The name of the pet`,
              `      - name: photo-urls`,
              `        aliases:`,
              `          - p`,
              `        type: string`,
              `        summary: A list of photo URLs to display for the pet`,
              `        description: |`,
              `          Provide this flag multiple times to set multiple photo URLs.`,
              `        variadic: true`,
              `      - name: status`,
              `        type: string`,
              `        summary: The pet status in the store`,
              `        choices:`,
              `          - value: available`,
              `          - value: pending`,
              `          - value: sold`,
              `      - name: tag`,
              `        type: string`,
              `        summary: Tag to assign to the pet for grouping/sorting`,
              `        description: |`,
              `          Provide this flag multiple times to add multiple tags.`,
              `        variadic: true`,
              ``,
              `  petstore pet find-by-status [flags]:`,
              `    summary: Find pets by status`,
              `    flags:`,
              `      - name: status`,
              `        type: string`,
              `        summary: The status to filter pets by`,
              `        choices:`,
              `          - value: available`,
              `          - value: pending`,
              `          - value: sold`,
              `        required: true`,
              ``,
              `  # ...`,
            ]}
          />
          <p>
            You can find the full example document{" "}
            <a
              href="https://github.com/bcdxn/opencli/blob/main/examples/petstore-cli.ocs.yaml"
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
          <p>Set up a fresh Go module and pull in the petstore spec:</p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ mkdir petstore && cd petstore`,
              `$ go mod init petstore`,
              `# pull in the full example ocs file (or use your own)`,
              `$ curl -O https://raw.githubusercontent.com/bcdxn/opencli/refs/heads/main/examples/petstore-cli.ocs.yaml`,
            ]}
          />

          <p>
            That's it for setup — one spec file, one module. Now we're ready to
            generate code.
          </p>
        </div>
      </div>

      {/* Step 4: Generate Boilerplate Code */}
      <div className="guide-step">
        <div className="guide-step__number">4</div>
        <div className="guide-step__content">
          <h4>Generate Boilerplate Code</h4>
          <p>
            A single <span className="guide-inline-code">ocli gen cli</span>{" "}
            command produces all the scaffolding. We'll generate a
            urfave/cli-based CLI here, but the same process works for Cobra. If
            you want to see a JS/TS example checkout the{" "}
            <a href="/docs/generating-ts-code">Generating TS Code</a> docs.
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ ocli gen cli \\`,
              `  --framework urfavecli \\`,
              `  --out ./internal \\`,
              `  ./petstore-cli.ocs.yaml`,
              `# → Reading spec:       ./petstore-cli.ocs.yaml`,
              `# → Generating CLI code:    framework=cobra, output=./internal`,
              `# ✓ CLI Code written to: ./internal`,
            ]}
          />

          <p>Then resolve dependencies:</p>

          <HighlightedCodeBlock language="sh" lines={["$ go mod tidy"]} />

          <p>
            All generated code is encapsulated in the{" "}
            <span className="guide-inline-code">gencli</span> package. Each
            command gets its own file, plus supporting files for bootstrapping,
            error handling, and I/O management:
          </p>

          <HighlightedCodeBlock
            language="plain"
            lines={[
              `go.mod`,
              `petstore-cli.ocs.yaml`,
              `internal/`,
              `\u2514\u2500\u2500 gencli/`,
              `    \u251c\u2500\u2500 actions.gen.go    Actions interface & command signatures`,
              `    \u251c\u2500\u2500 errors.gen.go     CLI error types & exit codes`,
              `    \u251c\u2500\u2500 help.gen.go       Default help/usage messaging`,
              `    \u251c\u2500\u2500 iostreams.gen.go  Standard I/O streams abstraction`,
              `    \u251c\u2500\u2500 params.gen.go     Command flags & parameter types`,
              `    \u251c\u2500\u2500 run.go            CLI entry point (Run function)`,
              `    \u2514\u2500\u2500 cmd_...           Generated Cobra command definitions`,
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
            {" "}
            Let's take a look at the all-important{" "}
            <span className="guide-inline-code">
              internal/gencli/actions.gen.go
            </span>{" "}
            Below shows an example of a method from that interface.
          </p>

          <HighlightedCodeBlock
            language="go"
            lines={[
              `// internal/gencli/actions.gen.go`,
              `type ActionsInterface interface {`,
              `  // ...`,
              `  func NewCmdPetstorePetAdd(`,
              `    ctx context.Context,`,
              `    args PetstorePetAddArgs,`,
              `    flags PetstorePetAddFlags,`,
              `  ) error {`,
              `    // ...`,
              `  }`,
              `}`,
            ]}
          />

          <div className="guide-callout">
            <p>
              Look at your generated{" "}
              <span className="guide-inline-code">
                internal/gencli/actions.gen.go
              </span>{" "}
              to see the full interface we'll need to implement.
            </p>
          </div>

          <p>
            Notice that the methods we need to implement have no
            framework-dependencies injected. We could reuse our same
            <span className="guide-inline-code">ActionsInterface</span>{" "}
            implementation for multiple frameworks within the same language
            (e.g. cobra and urfave/cli within Go).
          </p>
          <p>
            The generated types for{" "}
            <span className="guide-inline-code">args</span> and{" "}
            <span className="guide-inline-code">flags</span> are strongly typed,
            so you get compile-time safety — no more typos in flag names or
            mismatched types.
          </p>

          <p>
            Next we can take a look at the generated command files{" "}
            <span className="guide-inline-code">
              internal/gencli/cmd_*.gen.go
            </span>
            . Each generated command file adapts our ActionsInterface methods,
            handling the framework specifics of parsing args and flags and
            passing them to our framework-<i>agnostic</i> implementations.
          </p>

          <p>
            If you're interested, you can look at a generated file to see how
            the
            <span className="guide-inline-code">Action</span> handler delegates
            to the corresponding function on our struct implementing the{" "}
            <span className="guide-inline-code">ActionsInterface</span> shown
            below. But in general you can treat these generated command file as
            black boxes.
          </p>

          <HighlightedCodeBlock
            language="go"
            lines={[
              `cmd := &cli.Command{`,
              `  Name:        "add",`,
              `  Usage:       "Add a new pet to the store",`,
              `  Action: func(ctx context.Context, c *cli.Command) error {`,
              `    // ... parse and validate args/flags`,
              `    return a.PetstorePetAdd(ctx, cmdArgs, cmdFlags)`,
              `  },`,
              `}`,
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
            This is where you write your actual business logic. Create a type
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
            Start by creating a new package for your implementation to keep it
            separate from the generated code in the{" "}
            <span className="guide-inline-code">gencli</span> package:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ mkdir -p ./internal/cliapp`,
              `$ touch ./internal/cliapp/actions.go`,
            ]}
          />

          <p>
            Define your <span className="guide-inline-code">Actions</span> type:
          </p>

          <HighlightedCodeBlock
            language="go"
            lines={[
              `package cliapp`,
              ``,
              `import (`,
              `  "context"`,
              `  "fmt"`,
              ``,
              `  "petstore/internal/gencli"`,
              `  "github.com/bcdxn/opencli/spec"`,
              `)`,
              ``,
              `func NewActions(version string) Actions {`,
              `  return Actions{version: version}`,
              `}`,
              ``,
              `type Actions struct {`,
              `  version string`,
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
            language="go"
            lines={[
              `func (a Actions) PetstoreList(ctx context.Context) error {`,
              `  fmt.Println("listing all resources...")`,
              `  return nil`,
              `}`,
              ``,
              `func (a Actions) PetstorePetAdd(`,
              `  ctx context.Context,`,
              `  args gencli.PetstorePetAddArgs,`,
              `  flags gencli.PetstorePetAddFlags,`,
              `) error {`,
              `  fmt.Printf("adding pet: name=%s, status=%s, tags=%v\\n",`,
              `    flags.Name, flags.Status, flags.Tag)`,
              `  return nil`,
              `}`,
              ``,
              `func (a Actions) PetstorePetUpdate(`,
              `  ctx context.Context,`,
              `  args gencli.PetstorePetUpdateArgs,`,
              `  flags gencli.PetstorePetUpdateFlags,`,
              `) error {`,
              `  fmt.Printf("updating pet with data from: %s\\n", args.PathToReqBody)`,
              `  return nil`,
              `}`,
              ``,
              `// ... implement remaining methods to satisfy ActionsInterface ...`,
            ]}
          />

          <div className="guide-callout">
            <p>
              You can download a full example implementation{" "}
              <a href="/assets/code/actions.go">here</a>.
            </p>
          </div>

          <p>
            Finally, wire up the helper methods using sensible defaults provided
            by the generated code (or replace them with custom implementations
            if you need tailored behavior):
          </p>

          <HighlightedCodeBlock
            language="go"
            lines={[
              `func (a Actions) HelpFunc(cmd *spec.CommandItem) {`,
              `  gencli.DefaultHelpFunc(a, cmd)`,
              `}`,
              ``,
              `func (a Actions) UsageFunc(cmd *spec.CommandItem) error {`,
              `  gencli.DefaultUsageFunc(a, cmd)`,
              `  return nil`,
              `}`,
              ``,
              `func (a Actions) IOStreams() gencli.IOStreams {`,
              `  return gencli.DefaultIOS()`,
              `}`,
              ``,
              `func (a Actions) Version() string {`,
              `  return a.version`,
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
            <span className="guide-inline-code">main.go</span>:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[`$ mkdir -p cmd/petstore`, `$ touch cmd/petstore/main.go`]}
          />

          <HighlightedCodeBlock
            language="go"
            lines={[
              `package main`,
              ``,
              `import (`,
              `  "context"`,
              `  "os"`,
              ``,
              `  "petstore/internal/cliapp"`,
              `  "petstore/internal/gencli"`,
              `)`,
              ``,
              `var version = "DEV"`,
              ``,
              `func main() {`,
              `  actions := cliapp.NewActions(version)`,
              `  code := gencli.Run(context.Background(), actions)`,
              `  os.Exit(code)`,
              `}`,
            ]}
          />

          <p>
            Just three lines of substance, and critically — no framework
            dependencies in your user-land code.
          </p>
        </div>
      </div>

      {/* Step 7: Try it Out */}
      <div className="guide-step">
        <div className="guide-step__number">7</div>
        <div className="guide-step__content">
          <h4>Try It Out</h4>
          <p>That's the entire application. Let's run it:</p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ go run cmd/petstore/main.go --help`,
              `# An example CLI Document describing operations a petstore CLI may provide.`,
              `#`,
              `# USAGE:`,
              `#   petstore {command} <arguments> [flags]`,
              `#`,
              `# AVAILABLE COMMANDS`,
              `#   list  List all endpoints available`,
              `#   pet   A collection of commands for managing pets`,
              `#   store A collection of commands for store operations`,
              `#   user  A collection of commands for user management`,
            ]}
          />

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ go run cmd/petstore/main.go pet add --name fluffy --status available --tag dog`,
              `# adding pet: name=fluffy, status=available, tags=[dog]`,
            ]}
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
            <li key="code-generation-go">
              <a href="/docs/getting-started" className="guide-nav__link">
                Getting Started
              </a>
              <a href="/docs/markdown-docs" className="guide-nav__link">
                Markdown Docs
              </a>
              <a href="/docs/html-docs" className="guide-nav__link">
                HTML Docs
              </a>
              <a
                href="/docs/code-generation-go"
                className="guide-nav__link is-active"
              >
                Code Generation (Go)
              </a>
            </li>
          </ul>
        </nav>

        {/* Main content */}
        <main className="guide-main">
          <GeneratingGoCodePage key="code-generation-go" />
        </main>
      </div>
    </div>
  );
}
