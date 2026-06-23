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
        <pre>{schema}</pre>
      </main>
    </div>
  );
}
