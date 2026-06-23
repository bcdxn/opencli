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
            <h1>
              Design your CLI once, then generate docs and tooling from it.
            </h1>
            <p>
              OpenCLI is a contract-first specification for command line tools.
              Describe your CLI in a document format that is both human readable
              and machine readable. Validate it, and generate consistent
              documentation.
            </p>

            <div className="hero-actions">
              <Link className="button primary" to="/editor">
                Open Live Editor
              </Link>
              <Link className="button secondary" to="/spec">
                View Specification
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
