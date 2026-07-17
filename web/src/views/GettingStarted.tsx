"use client";

import React, { useState, useCallback } from "react";
import SyntaxHighlighter from "react-syntax-highlighter";
import { nord } from "react-syntax-highlighter/dist/esm/styles/hljs";
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
          {copied ? "✓ Copied" : "Copy"}
        </button>
      </div>
      <div className="body">
        <SyntaxHighlighter
          language="yaml"
          style={nord}
          customStyle={{ background: "#1A1D24" }}
        >
          {lines.join("\n")}
        </SyntaxHighlighter>
      </div>
    </div>
  );
}

// ── Getting Started Page ──────────────────────────────────────────────────────

function GettingStartedPage() {
  return (
    <>
      <h2 className="guide-section__title">Getting Started</h2>
      <p className="guide-section__subtitle">
        Learn how to create your first OpenCLI Specification document, validate
        it, and generate documentation — all from scratch.
      </p>

      {/* Step 1 */}
      <div className="guide-step">
        <div className="guide-step__number">1</div>
        <div className="guide-step__content">
          <h4>Create your spec file</h4>
          <p>
            Start by creating a YAML file that describes your CLI. We'll build a
            simple example called{" "}
            <span className="guide-inline-code">pleasantries</span> — a CLI that
            can greet or bid farewell to someone.
          </p>
          <p>
            Create a file called{" "}
            <span className="guide-inline-code">pleasantries-cli.ocs.yaml</span>{" "}
            and begin with the version and metadata:
          </p>

          <HighlightedCodeBlock
            language="yaml"
            lines={[
              `opencliVersion: "1.0.0-alpha.13"`,
              ``,
              `# Metadata about your CLI"`,
              `info:`,
              `  title: "Pleasantries"`,
              `  summary: "A fun CLI to greet or bid farewell"`,
              `  version: "1.0.0"`,
              `  binary: "pleasantries"`,
            ]}
          />

          <p>
            The <span className="guide-inline-code">opencliVersion</span> field
            specifies which version of the spec you're using. The{" "}
            <span className="guide-inline-code">info</span> block contains
            human-readable metadata about your CLI — its name, description,
            version, and the binary name.
          </p>
        </div>
      </div>

      {/* Step 2 */}
      <div className="guide-step">
        <div className="guide-step__number">2</div>
        <div className="guide-step__content">
          <h4>Define your commands</h4>
          <p>
            Next, add the <span className="guide-inline-code">commands</span>{" "}
            section. Each command is a key that describes the invocation
            pattern, followed by its properties:
          </p>

          <HighlightedCodeBlock
            language="yaml"
            lines={[
              `commands:`,
              `  # Group command - acts a container of subcommands`,
              `  pleasantries {command} <arguments> [flags]:`,
              `    kind: "group"`,
              "",
              `  # The 'greet' command`,
              `  pleasantries greet <name> [flags]:`,
              `    summary: "Say hello"`,
              `    args:`,
              `    - name: "name"`,
              `      summary: "A name to include in the greeting"`,
              `      required: true`,
              `      type: "string"`,
              `    flags:`,
              `    - name: "language"`,
              `      summary: "The language of the greeting"`,
              `      type: "string"`,
              `      choices:`,
              `      - value: "english"`,
              `      - value: "spanish"`,
              `      default: "english"`,
            ]}
          />

          <p>
            Here we define a <span className="guide-inline-code">greet</span>{" "}
            command that takes a required{" "}
            <span className="guide-inline-code">&lt;name&gt;</span> argument and
            an optional <span className="guide-inline-code">--language</span>{" "}
            flag with choices.
          </p>
          <p>
            Notice that we also defined the{" "}
            <span className="guide-inline-code">pleasantries</span> which acts
            as a 'grouping' command. Defining grouping commands is optional -
            the code generation and docs generation will walk the command tree
            either way - but they can be a good way to help readability of the
            OpenCLI document and add additional comments.
          </p>
        </div>
      </div>

      {/* Step 3 */}
      <div className="guide-step">
        <div className="guide-step__number">3</div>
        <div className="guide-step__content">
          <h4>Add examples</h4>
          <p>Make your spec more useful by adding example invocations:</p>

          <HighlightedCodeBlock
            language="yaml"
            lines={[
              `    examples:`,
              `    - title: "greet the user"`,
              `      content:`,
              `        $ pleasantries greet --language english John`,
              `        # Hello, John`,
            ]}
          />

          <p>
            The complete spec with both commands (
            <span className="guide-inline-code">greet</span> and{" "}
            <span className="guide-inline-code">farewell</span>) is available in
            the{" "}
            <a
              href="https://github.com/bcdxn/opencli/blob/main/examples/pleasantries-cli.ocs.yaml"
              target="_blank"
              rel="noreferrer"
            >
              examples directory
            </a>
            .
          </p>
        </div>
      </div>

      {/* Step 4 */}
      <div className="guide-step">
        <div className="guide-step__number">4</div>
        <div className="guide-step__content">
          <h4>Install the OpenCLI CLI</h4>
          <p>
            To validate your spec, install the{" "}
            <span className="guide-inline-code">ocli</span> tool:
          </p>

          <HighlightedCodeBlock
            language="shell"
            lines={[`$ go install github.com/bcdxn/opencli/cmd/ocli@latest`]}
          />

          <p>
            This requires Go 1.21+ installed on your system. Once installed,
            you'll have the <span className="guide-inline-code">ocli</span>{" "}
            command available.
          </p>
        </div>
      </div>

      {/* Step 5 */}
      <div className="guide-step">
        <div className="guide-step__number">5</div>
        <div className="guide-step__content">
          <h4>Validate your spec</h4>
          <p>
            Run the <span className="guide-inline-code">check</span> command to
            validate your document against the OpenCLI Specification:
          </p>

          <HighlightedCodeBlock
            language="shell"
            lines={[
              `$ ocli check ./pleasantries-cli.ocs.yaml`,
              `# ✓ pleasantries-cli.ocs.yaml is valid`,
            ]}
          />

          <p>
            If there are any issues, the CLI will report them with line numbers
            and descriptions. A valid spec means you're ready to generate
            documentation or code.
          </p>

          <div className="guide-callout">
            <p>
              <strong>Tip:</strong> The{" "}
              <span className="guide-inline-code">check</span> command validates
              both the JSON Schema structure and additional rules that can't be
              expressed in schema alone. Always run it before generating docs or
              code.
            </p>
          </div>
        </div>
      </div>

      {/* Next steps */}
      <h3>What's next?</h3>
      <ul className="guide-section__subtitle">
        <li>
          Generate <a href="/docs/markdown-docs">Markdown documentation</a> for
          your CLI
        </li>
        <li>
          Generate <a href="/docs/html-docs">HTML documentation</a> that looks
          great in a browser
        </li>
      </ul>
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
            <li key="getting-started">
              <a
                href="/docs/getting-started"
                className="guide-nav__link is-active"
              >
                Getting Started
              </a>
              <a href="/docs/markdown-docs" className="guide-nav__link">
                Markdown Docs
              </a>
              <a href="/docs/html-docs" className="guide-nav__link">
                HTML Docs
              </a>
            </li>
          </ul>
        </nav>

        {/* Main content */}
        <main className="guide-main">
          <GettingStartedPage key="getting-started" />
        </main>
      </div>
    </div>
  );
}
