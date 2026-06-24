import { useEffect, useState } from "react";
import SiteHeader from "../components/SiteHeader";
import "./SpecPage.css";

export default function SpecPage() {
  const [schema, setSchema] = useState("Loading schema...");

  useEffect(() => {
    fetch(`${import.meta.env.BASE_URL}spec.schema.json`)
      .then((res) => res.text())
      .then((text) => {
        try {
          const parsed = JSON.parse(text);
          setSchema(JSON.stringify(parsed, null, 2));
        } catch {
          setSchema(text);
        }
      })
      .catch(() => {
        setSchema("Failed to load spec.schema.json");
      });
  }, []);

  return (
    <div className="spec-page">
      <SiteHeader />
      <main className="spec-main">
        <section className="spec-disclaimer" aria-label="Validation disclaimer">
          <p className="spec-disclaimer-title">Validation note</p>
          <p className="spec-disclaimer-copy">
            This JSON Schema helps enforce structural validity, but not every
            validation rule can be represented in JSON Schema alone.
          </p>
          <p className="spec-disclaimer-copy">
            For full spec validation, run the <a href={`${import.meta.env.BASE_URL}docs#ocli-check`}>OpenCLI CLI</a> check command.
          </p>
        </section>
        <pre>{schema}</pre>
      </main>
    </div>
  );
}
