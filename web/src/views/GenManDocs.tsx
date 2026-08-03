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
          {copied ? "✓ Copied" : "Copy"}
        </button>
      </div>
      <div className="body">
        <SyntaxHighlighter
          language="bash"
          style={nord}
          customStyle={{ background: "#1A1D24" }}
        >
          {lines.join("\n")}
        </SyntaxHighlighter>
      </div>
    </div>
  );
}

// ── Man Pages Docs Page ───────────────────────────────────────────────────────

function ManDocsPage() {
  return (
    <>
      <h2 className="guide-section__title">Generate Man Pages</h2>
      <p className="guide-section__subtitle">
        Generate proper roff/troff man pages from your OpenCLI Specification —
        ready to ship with your CLI so users get native <code>man</code> support
        out of the box.
      </p>

      <div className="guide-image">
        <img src="/img/demo.gif" alt="OpenCLI man page generation demo" />
      </div>

      {/* Step 1 */}
      <div className="guide-step__vertical-space"></div>
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

      {/* Step 2 */}
      <div className="guide-step">
        <div className="guide-step__number">2</div>
        <div className="guide-step__content">
          <h4>Generate a man page</h4>
          <p>
            Run the <span className="guide-inline-code">gen docs</span> command
            with the <span className="guide-inline-code">--format man</span>{" "}
            flag:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ ocli gen docs \\`,
              `  --format man \\`,
              `  --out ./man \\`,
              `  ./pleasantries-cli.ocs.yaml`,
            ]}
          />

          <p>
            This produces a <span className="guide-inline-code">.1</span> file
            (e.g.,{" "}
            <span className="guide-inline-code">pleasantries-cli.ocs.1</span>)
            in the output directory. The man page follows standard roff/troff
            conventions with properly formatted sections for name, synopsis,
            description, commands, arguments, flags, and examples.
          </p>

          <h4>What's in the generated man page</h4>
          <p>A generated man page includes:</p>
          <ul>
            <li>
              <span className="guide-inline-code">NAME</span> — CLI name and
              summary
            </li>
            <li>
              <span className="guide-inline-code">SYNOPSIS</span> — Usage
              pattern
            </li>
            <li>
              <span className="guide-inline-code">DESCRIPTION</span> — Full CLI
              description with capabilities
            </li>
            <li>
              <span className="guide-inline-code">GLOBAL OPTIONS</span> — Shared
              flags across all commands
            </li>
            <li>
              <span className="guide-inline-code">COMMANDS</span> — Each command
              with its arguments, flags, and examples
            </li>
          </ul>

          <div className="guide-callout">
            <p>
              <strong>Why man pages?</strong> Man pages are the Unix standard
              for CLI documentation. Generating them from your OpenCLI spec
              means your docs stay in sync with your spec — no more manually
              maintained roff files that drift from reality.
            </p>
          </div>
        </div>
      </div>

      {/* Step 3 */}
      <div className="guide-step">
        <div className="guide-step__number">3</div>
        <div className="guide-step__content">
          <h4>Use the man page</h4>
          <p>
            You have a few options for working with your generated man page:
          </p>

          <h4>Option 1 — Preview without installing</h4>
          <p>
            Quick preview straight from the output directory using the{" "}
            <span className="guide-inline-code">man</span> command with a
            relative path:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[`$ man ./man/pleasantries-cli.ocs.1`]}
          />

          <p>
            No installation needed — great for validating output before
            shipping.
          </p>

          <h4>Option 2 — Install system-wide</h4>
          <p>
            Copy the man page to your system's man directory so it's available
            globally:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ sudo cp ./man/pleasantries-cli.ocs.1 /usr/local/share/man/man1/pleasantries.1`,
              `$ man pleasantries`,
            ]}
          />

          <p>
            On macOS and Linux, the{" "}
            <span className="guide-inline-code">man</span> command will
            automatically pick it up after installation.
          </p>

          <h4>Option 3 — Ship with your GoReleaser Homebrew tap</h4>
          <p>
            If you distribute your CLI via{" "}
            <a href="https://goreleaser.com" target="_blank" rel="noreferrer">
              GoReleaser
            </a>
            , include the man page in your Homebrew formula so it's installed
            alongside the binary:
          </p>

          <HighlightedCodeBlock
            language="yaml"
            lines={[
              "# .goreleaser.yaml",
              "archives:",
              "  - files:",
              "    - README.md",
              "    - LICENSE",
              "    - pleasantries-cli.ocs.1",
            ]}
          />

          <p>
            With this config, users who install via{" "}
            <span className="guide-inline-code">brew install pleasantries</span>{" "}
            will get the man page installed automatically — no extra steps
            required.
          </p>

          <div className="guide-callout">
            <p>
              <strong>Tip:</strong> If your CLI has multiple commands, each gets
              its own section in the generated man page — no need to split them
              into separate files. The whole spec becomes one cohesive man page.
            </p>
          </div>
        </div>
      </div>

      {/* Next steps */}
      <div className="whats-next">
        <h3>What's next?</h3>
        <ul>
          <li>
            Generate <a href="/docs/markdown-docs">Markdown documentation</a>{" "}
            for GitHub or static sites
          </li>
          <li>
            Generate <a href="/docs/html-docs">HTML documentation</a> that looks
            great in a browser
          </li>
          <li>
            Learn how to{" "}
            <a href="/docs/code-generation-go">generate CLI code</a> from your
            spec
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
            <li key="getting-started">
              <a href="/docs/getting-started" className="guide-nav__link">
                Getting Started
              </a>
              <a href="/docs/markdown-docs" className="guide-nav__link">
                Markdown Docs
              </a>
              <a href="/docs/html-docs" className="guide-nav__link">
                HTML Docs
              </a>
              <a href="/docs/man-pages" className="guide-nav__link is-active">
                Man Pages
              </a>
              <a href="/docs/code-generation-go" className="guide-nav__link">
                Code Generation (Go)
              </a>
            </li>
          </ul>
        </nav>

        {/* Main content */}
        <main className="guide-main">
          <ManDocsPage key="man-pages" />
        </main>
      </div>
    </div>
  );
}
