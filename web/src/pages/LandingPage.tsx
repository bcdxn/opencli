import { Link } from "react-router-dom";
import SiteHeader from "../components/SiteHeader";
import "./LandingPage.css";
import React from "react";

const yamlSample = `opencliVersion: 1.0.0-alpha.8

info:
  title: Pleasantries CLI
  summary: A fun CLI to greet the caller
  version: 1.0.0
  binary: pleasantries

commands:
  pleasantries {command} <args> [flags]:
    group: true

  pleasantries greet <args> [flags]:
    summary: "Say hello"
    args:
      - name: "name"
        required: true
        type: "string"
    flags:
      - name: "language"
        type: "string"
        default: "english"`;

function renderYamlValue(value: string): React.JSX.Element {
  if (value.startsWith('"') || value.startsWith("'")) {
    return <span className="yaml-string">{value}</span>;
  }
  if (value === "true" || value === "false") {
    return <span className="yaml-bool">{value}</span>;
  }
  return <span className="yaml-value">{value}</span>;
}

function tokenizeYamlLine(line: string, i: number): React.JSX.Element {
  if (line.trim() === "") {
    return <span key={i}>{"\n"}</span>;
  }

  const indentMatch = line.match(/^(\s*)/);
  const indent = indentMatch ? indentMatch[1] : "";
  let rest = line.slice(indent.length);

  if (rest.startsWith("#")) {
    return (
      <span key={i}>
        {indent}
        <span className="yaml-comment">{rest}</span>
        {"\n"}
      </span>
    );
  }

  let bullet: React.JSX.Element | null = null;
  if (rest.startsWith("- ")) {
    bullet = <span className="yaml-bullet">{"- "}</span>;
    rest = rest.slice(2);
  }

  const kvMatch = rest.match(/^([^:]+)(:)(\s*)(.*)$/);
  if (kvMatch) {
    const [, key, colon, space, value] = kvMatch;
    return (
      <span key={i}>
        {indent}
        {bullet}
        <span className="yaml-key">{key}</span>
        <span className="yaml-punct">{colon}</span>
        {space}
        {value ? renderYamlValue(value) : null}
        {"\n"}
      </span>
    );
  }

  return (
    <span key={i}>
      {line}
      {"\n"}
    </span>
  );
}

export default function LandingPage() {
  return (
    <div className="landing-route">
      <SiteHeader />

      <main className="landing-main">
        <section className="hero-shell">
          <div className="hero-copy">
            <div className="hero-kicker">OPENCLI</div>
            <h1>Contract-First CLI Design</h1>
            <h2>Define your interface. Automate the rest.</h2>
            <p>
              OpenCLI is an open document specification for command-line tools.
              Author verifiable, human-readable CLI definitions. Generate
              documentation and framework-specific code from a single spec.
            </p>

            <div className="hero-actions">
              <Link className="button primary" to="/editor">
                Open Live Editor
              </Link>
              <Link className="button secondary" to="/docs">
                View CLI Docs
              </Link>
            </div>

            <div className="hero-tags">
              <span>contract-first</span>
              <span>framework-agnostic</span>
              <span>agent-ready</span>
            </div>
          </div>

          <section
            className="terminal-shot"
            aria-label="CLI screenshot placeholder"
          >
            <div className="terminal-titlebar">
              <span />
              <span />
              <span />
            </div>
            <div className="terminal-body">
              <pre>
                <code>
                  {yamlSample
                    .split("\n")
                    .map((line, i) => tokenizeYamlLine(line, i))}
                </code>
              </pre>
            </div>
          </section>
        </section>

        <section className="benefits">
          <article>
            <h2>Contract-First Development</h2>
            <p>
              Promote stable command contracts and keep implementation details
              decoupled from CLI frameworks.
            </p>
          </article>
          <article>
            <h2>Tooling That Scales</h2>
            <p>
              Validate specs, generate docs, and standardize outputs across
              teams and languages.
            </p>
          </article>
          <article>
            <h2>Better Agent Understanding</h2>
            <p>
              Give LLMs and automation tools a structured, explicit model of
              your CLI surface area.
            </p>
          </article>
        </section>
      </main>
    </div>
  );
}
