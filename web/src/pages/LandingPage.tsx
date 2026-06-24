import { Link } from "react-router-dom";
import SiteHeader from "../components/SiteHeader";
import "./LandingPage.css";

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
              <p>$ ocli check ./my-cli.ocs.yaml</p>
              <p className="ok">→ Checking ./examples/petstore-cli.ocs.yaml</p>
              <p className="ok">✓ Document is valid</p>
              <p>$ ocli gen docs --format markdown ./my-cli.ocs.yaml</p>
              <p className="ok">→ Reading spec: ./my-cli.ocs.yaml</p>
              <p className="ok">
                → Generating docs: format=markdown, output=./docs
              </p>
              <p className="ok">
                ✓ Documentation written to: docs/my-cli.ocs.html
              </p>
              <p>$ ocli gen cli --framework cobra ./my-cli.ocs.yaml</p>
              <p className="ok">→ Reading spec: ./my-cli.ocs.yaml</p>
              <p className="ok">
                → Generating CLI boilerplate: framework=cobra, output=./cli
              </p>
              <p className="ok">✓ Boilerplate written to: cli/...</p>
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
