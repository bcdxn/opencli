import { Link } from "react-router-dom";
import "../App.css";
import "./LandingPage.css";

export default function LandingPage() {
  return (
    <div className="landing-route">
      <header>
        <div className="header-logo">
          <Link to="/">OPENCLI</Link>
        </div>
        <nav className="header-nav">
          <Link to="/editor">Editor</Link>
          <Link to="/spec">Schema</Link>
          <a
            className="github-link"
            href="https://github.com/bcdxn/opencli"
            target="_blank"
            rel="noreferrer"
            aria-label="Open the OpenCLI repository on GitHub"
          >
            <svg viewBox="0 0 24 24" aria-hidden="true" focusable="false">
              <path d="M12 2C6.47 2 2 6.48 2 12.07c0 4.46 2.87 8.24 6.84 9.58.5.09.68-.22.68-.48 0-.24-.01-.86-.01-1.69-2.78.62-3.37-1.37-3.37-1.37-.46-1.2-1.12-1.52-1.12-1.52-.92-.64.07-.63.07-.63 1.02.07 1.56 1.06 1.56 1.06.91 1.58 2.39 1.12 2.97.86.09-.67.35-1.12.63-1.38-2.22-.26-4.56-1.13-4.56-5.03 0-1.11.39-2.02 1.03-2.73-.1-.26-.45-1.32.1-2.75 0 0 .84-.27 2.75 1.04.8-.23 1.65-.35 2.5-.35s1.7.12 2.5.35c1.91-1.31 2.75-1.04 2.75-1.04.55 1.43.2 2.49.1 2.75.64.71 1.03 1.62 1.03 2.73 0 3.91-2.35 4.77-4.58 5.02.36.32.68.95.68 1.92 0 1.39-.01 2.51-.01 2.85 0 .26.18.58.69.48A10.1 10.1 0 0 0 22 12.07C22 6.48 17.53 2 12 2Z" />
            </svg>
            <span>GitHub</span>
          </a>
        </nav>
      </header>

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
