import { useEffect, useState } from "react";
import SiteHeader from "../components/SiteHeader";
import { useI18n } from "../i18n";
import "./SchemaPage.css";

export default function SpecPage() {
  const { t } = useI18n();
  const [schema, setSchema] = useState(t("schema.loading"));

  useEffect(() => {
    setSchema(t("schema.loading"));
  }, [t]);

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
        setSchema(t("schema.failed"));
      });
  }, [t]);

  return (
    <div className="spec-page">
      <SiteHeader />
      <main className="spec-main">
        <section
          className="spec-disclaimer"
          aria-label={t("schema.disclaimerAria")}
        >
          <p className="spec-disclaimer-title">{t("schema.note.title")}</p>
          <p className="spec-disclaimer-copy">{t("schema.note.copy1")}</p>
          <p className="spec-disclaimer-copy">
            {t("schema.note.copy2.prefix")}{" "}
            <a href={`${import.meta.env.BASE_URL}docs#ocli-check`}>
              {t("schema.note.copy2.link")}
            </a>{" "}
            {t("schema.note.copy2.suffix")}
          </p>
        </section>
        <pre>{schema}</pre>
      </main>
    </div>
  );
}
