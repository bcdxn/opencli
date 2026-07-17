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

// ── HTML Docs Page ────────────────────────────────────────────────────────────

function HtmlDocsPage() {
  return (
    <>
      <h2 className="guide-section__title">Generate HTML Docs</h2>
      <p className="guide-section__subtitle">
        Generate beautiful, interactive HTML documentation from your OpenCLI
        Specification that can stand alone or be embedded in any web page.
      </p>

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

      {/* Step 2 */}
      <div className="guide-step">
        <div className="guide-step__number">2</div>
        <div className="guide-step__content">
          <h4>Standalone HTML Page</h4>
          <p>
            Generate elegant documentation in a self-contained{" "}
            <span className="guide-inline-code">index.html</span> file that can
            be deployed to any static site, e.g. GitHub pages.
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ ocli gen docs`,
              `  --format html-page`,
              `  --out ./web`,
              `  ./pleasantries-cli.ocs.yaml`,
            ]}
          />

          <p>
            This produces a fully styled HTML page in
            <span className="guide-inline-code">
              ./web/pleasantries-cli.ocs.html
            </span>{" "}
            ready to open in any browser:
          </p>

          <div className="guide-image">
            <img
              src="/img/gen-html-docs-screenshot.png"
              alt="OpenCLI HTML documentation preview"
            />
          </div>
        </div>
      </div>

      {/* Step 3 */}
      <div className="guide-step">
        <div className="guide-step__number">3</div>
        <div className="guide-step__content">
          <h4>Embeddable HTML</h4>
          <p>
            You can also generate a lightweight JavaScript bundle that renders
            docs inside an existing web page. This is useful when you want to
            embed CLI documentation alongside other content:
          </p>

          <HighlightedCodeBlock
            language="sh"
            lines={[
              `$ ocli gen docs`,
              `  --format html-embed`,
              `  --out ./assets`,
              `  ./pleasantries-cli.ocs.yaml`,
            ]}
          />

          <p>
            This writes a single{" "}
            <span className="guide-inline-code">ocli-docs.js</span> file.
            Include it in any HTML page to embed your docs:
          </p>

          <HighlightedCodeBlock
            language="html"
            lines={[
              `<html>`,
              `  <head>`,
              `    <script src="./assets/ocli-docs.js"></script>`,
              `  <head>`,
              `  <body>`,
              `    <div id="docs"></div>`,
              `    <script>`,
              `      window.OcliDocs({ containerId: "docs" });`,
              `    </script>`,
              `  </body>`,
              `</html>`,
            ]}
          />

          <div className="guide-callout">
            <p>
              <strong>How it works:</strong> The{" "}
              <span className="guide-inline-code">html-embed</span> format
              bundles everything into a single JavaScript file — no build step,
              no dependencies. Just drop the script tag and call{" "}
              <span className="guide-inline-code">window.OcliDocs()</span> with
              the container ID where you want the docs to render.
            </p>
          </div>

          <h3>Comparing formats</h3>
          <ul>
            <li>
              <span className="guide-inline-code">html-page</span> — A complete,
              styled HTML file. Best for standalone documentation sites.
            </li>
            <li>
              <span className="guide-inline-code">html-embed</span> — A
              JavaScript bundle you can embed in any existing page. Best for
              integrating CLI docs into a larger website.
            </li>
          </ul>
        </div>
      </div>

      {/* Next steps */}
      <div className="whats-next">
        <h3>What's next?</h3>
        <ul>
          <li>
            Try out a <a href="/editor">live preview</a> of the HTML
            documentation generation
          </li>
          <li>
            See a real world example of generated HTML docs on the{" "}
            <a href="/reference">OCLI Reference Page</a>
          </li>
          <li>
            Learn how to <a href="/docs/code-generation-go">generate code</a>{" "}
            from an Open CLI document
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
              <a href="/docs/html-docs" className="guide-nav__link is-active">
                HTML Docs
              </a>
              <a href="/docs/code-generation-go" className="guide-nav__link">
                Code Generation (Go)
              </a>
            </li>
          </ul>
        </nav>

        {/* Main content */}
        <main className="guide-main">
          <HtmlDocsPage key="html-docs" />
        </main>
      </div>
    </div>
  );
}
