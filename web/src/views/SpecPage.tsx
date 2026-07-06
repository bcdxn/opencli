"use client";

import { useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import SiteHeader from "../components/SiteHeader";
import { useI18n } from "../i18n";
import rawSchema from "../spec.schema.json";
import "./SpecPage.css";

// ── Types ─────────────────────────────────────────────────────────────────────

interface SchemaProp {
  type?: string | string[];
  $ref?: string;
  description?: string;
  format?: string;
  enum?: (string | number | boolean)[];
  items?: SchemaProp;
  properties?: Record<string, SchemaProp>;
  patternProperties?: Record<string, SchemaProp>;
  required?: string[];
  oneOf?: SchemaProp[];
  anyOf?: { required?: string[] }[];
  default?: string | number | boolean;
}

interface SchemaDoc extends SchemaProp {
  title?: string;
  $defs?: Record<string, SchemaProp>;
}

const schema = rawSchema as unknown as SchemaDoc;

// ── Helpers ───────────────────────────────────────────────────────────────────

function toId(name: string): string {
  return name.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
}

interface TypeInfo {
  label: string;
  linkedId: string | null;
  isArray: boolean;
}

function resolveTypeDisplay(prop: SchemaProp): TypeInfo {
  if (prop.$ref) {
    const name = prop.$ref.replace("#/$defs/", "");
    return { label: name, linkedId: toId(name), isArray: false };
  }
  if (prop.type === "array") {
    const items = prop.items;
    if (items?.$ref) {
      const name = items.$ref.replace("#/$defs/", "");
      return { label: name, linkedId: toId(name), isArray: true };
    }
    const itemType = items?.type;
    return {
      label: itemType
        ? Array.isArray(itemType)
          ? itemType.join(" | ")
          : itemType
        : "object",
      linkedId: null,
      isArray: true,
    };
  }
  if (prop.type === "object" && prop.patternProperties) {
    const vals = Object.values(prop.patternProperties);
    if (vals[0]?.$ref) {
      const name = vals[0].$ref.replace("#/$defs/", "");
      return {
        label: `Map<string, ${name}>`,
        linkedId: toId(name),
        isArray: false,
      };
    }
  }
  if (prop.oneOf) {
    return {
      label: prop.oneOf
        .map((s) => {
          if (Array.isArray(s.type)) return s.type[0] ?? "any";
          return s.type ?? "any";
        })
        .join(" | "),
      linkedId: null,
      isArray: false,
    };
  }
  if (prop.type) {
    return {
      label: Array.isArray(prop.type) ? prop.type.join(" | ") : prop.type,
      linkedId: null,
      isArray: false,
    };
  }
  return { label: "object", linkedId: null, isArray: false };
}

function resolveDescription(prop: SchemaProp): string {
  if (prop.description) return prop.description;
  const defs = schema.$defs ?? {};
  if (prop.$ref) {
    const name = prop.$ref.replace("#/$defs/", "");
    return defs[name]?.description ?? "";
  }
  if (prop.type === "array" && prop.items?.$ref) {
    const name = prop.items.$ref.replace("#/$defs/", "");
    return defs[name]?.description ?? "";
  }
  return "";
}

// ── FieldRow ──────────────────────────────────────────────────────────────────

interface FieldRow {
  name: string;
  typeInfo: TypeInfo;
  description: string;
  required: "required" | "conditional" | "optional";
  enumValues?: (string | number | boolean)[];
}

function buildFields(schemaProp: SchemaProp): FieldRow[] {
  if (!schemaProp.properties) return [];

  const topRequired = new Set<string>(schemaProp.required ?? []);

  const anyOfGroups = schemaProp.anyOf?.map(
    (s) => new Set<string>(s.required ?? []),
  );

  const anyOfAlwaysRequired =
    anyOfGroups && anyOfGroups.length > 0
      ? anyOfGroups.reduce(
          (acc, g) => new Set([...acc].filter((x) => g.has(x))),
        )
      : new Set<string>();

  const anyOfEverRequired = new Set<string>();
  if (anyOfGroups) {
    for (const g of anyOfGroups) g.forEach((r) => anyOfEverRequired.add(r));
  }

  const effectiveRequired = new Set([...topRequired, ...anyOfAlwaysRequired]);
  const conditionalFields = new Set(
    [...anyOfEverRequired].filter((r) => !effectiveRequired.has(r)),
  );

  return Object.entries(schemaProp.properties)
    .filter(([k]) => !k.startsWith("^"))
    .map(([name, prop]) => ({
      name,
      typeInfo: resolveTypeDisplay(prop),
      description: resolveDescription(prop),
      required: effectiveRequired.has(name)
        ? "required"
        : conditionalFields.has(name)
          ? "conditional"
          : "optional",
      enumValues: prop.enum,
    }));
}

// ── Markdown renderer ──────────────────────────────────────────────────────────

function Md({ children }: { children: string }) {
  if (!children) return null;
  return (
    <div className="spec-md">
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{children}</ReactMarkdown>
    </div>
  );
}

// ── Sub-components ────────────────────────────────────────────────────────────

function TypeBadge({ info }: { info: TypeInfo }) {
  const inner = info.linkedId ? (
    <a href={`#${info.linkedId}`} className="spec-type-link">
      {info.label}
    </a>
  ) : (
    <code className="spec-type-code">{info.label}</code>
  );

  if (info.isArray) {
    return (
      <span className="spec-type">
        <span className="spec-type-bracket">[</span>
        {inner}
        <span className="spec-type-bracket">]</span>
      </span>
    );
  }
  return <span className="spec-type">{inner}</span>;
}

function FieldTable({ fields }: { fields: FieldRow[] }) {
  return (
    <div className="spec-table-wrap">
      <table className="spec-table">
        <thead>
          <tr>
            <th>Field Name</th>
            <th>Type</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          {fields.map((f) => (
            <tr key={f.name}>
              <td className="spec-table-name-cell">
                <code className="spec-field-name">{f.name}</code>
              </td>
              <td className="spec-table-type-cell">
                <TypeBadge info={f.typeInfo} />
                {f.required === "required" && (
                  <span className="spec-badge spec-badge--required">
                    required
                  </span>
                )}
                {f.required === "conditional" && (
                  <span className="spec-badge spec-badge--conditional">
                    conditional
                  </span>
                )}
              </td>
              <td className="spec-table-desc-cell">
                <Md>{f.description}</Md>
                {f.enumValues && (
                  <ul className="spec-enum-inline">
                    {f.enumValues.map((v) => (
                      <li key={String(v)}>
                        <code>{String(v)}</code>
                      </li>
                    ))}
                  </ul>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

interface ObjectSectionProps {
  id: string;
  title: string;
  description?: string;
  schemaProp: SchemaProp;
}

function ObjectSection({
  id,
  title,
  description,
  schemaProp,
}: ObjectSectionProps) {
  const fields = buildFields(schemaProp);
  const enumValues = schemaProp.enum;

  return (
    <section id={id} className="spec-ref-section">
      <h2 className="spec-ref-section__title">{title}</h2>
      {description && (
        <div className="spec-ref-section__desc">
          <Md>{description}</Md>
        </div>
      )}

      {enumValues && (
        <>
          <h3 className="spec-ref-subsection">Enum Values</h3>
          <ul className="spec-enum-list">
            {enumValues.map((v) => (
              <li key={String(v)}>
                <code>{String(v)}</code>
              </li>
            ))}
          </ul>
        </>
      )}

      {fields.length > 0 && (
        <>
          <h3 className="spec-ref-subsection">Fields</h3>
          <FieldTable fields={fields} />
        </>
      )}
    </section>
  );
}

// ── Main Page ─────────────────────────────────────────────────────────────────

export default function SpecPage() {
  const { t } = useI18n();
  const [activeId, setActiveId] = useState("version");

  useEffect(() => {
    const sections =
      document.querySelectorAll<HTMLElement>(".spec-ref-section");
    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries
          .filter((e) => e.isIntersecting)
          .sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);
        if (visible.length > 0) {
          setActiveId(visible[0].target.id);
        }
      },
      { rootMargin: "-10% 0px -70% 0px", threshold: 0 },
    );
    sections.forEach((s) => observer.observe(s));
    return () => observer.disconnect();
  }, []);

  const defs = schema.$defs ?? {};
  const versionEnum = schema.properties?.opencliVersion?.enum;
  const version =
    versionEnum && versionEnum.length > 0
      ? String(versionEnum[versionEnum.length - 1])
      : "unknown";

  return (
    <div className="spec-page">
      <SiteHeader />
      <div className="spec-ref-layout">
        {/* Left nav */}
        <nav className="spec-ref-nav" aria-label="Schema navigation">
          <p className="spec-ref-nav__heading">On this page</p>
          <ul className="spec-ref-nav__list">
            <li>
              <a
                href="#version"
                className={`spec-ref-nav__link${
                  activeId === "version" ? " is-active" : ""
                }`}
              >
                Version
              </a>
            </li>
            <li>
              <a
                href="#opencli-specification"
                className={`spec-ref-nav__link${
                  activeId === "opencli-specification" ? " is-active" : ""
                }`}
              >
                OpenCLI Specification
              </a>
            </li>
          </ul>
          <p className="spec-ref-nav__subheading">Schema Objects</p>
          <ul className="spec-ref-nav__list">
            {Object.keys(defs).map((name) => {
              const id = toId(name);
              return (
                <li key={id}>
                  <a
                    href={`#${id}`}
                    className={`spec-ref-nav__link${
                      activeId === id ? " is-active" : ""
                    }`}
                  >
                    {name}
                  </a>
                </li>
              );
            })}
          </ul>
        </nav>

        {/* Main content */}
        <main className="spec-ref-main">
          {/* Disclaimer */}
          <section
            className="spec-disclaimer"
            aria-label={t("schema.disclaimerAria")}
          >
            <p className="spec-disclaimer-title">{t("schema.note.title")}</p>
            <p className="spec-disclaimer-copy">{t("schema.note.copy1")}</p>
            <p className="spec-disclaimer-copy">
              {t("schema.note.copy2.prefix")}{" "}
              <a href={`/docs#ocli-check`}>{t("schema.note.copy2.link")}</a>{" "}
              {t("schema.note.copy2.suffix")}
            </p>
          </section>

          {/* Version section */}
          <section id="version" className="spec-ref-section">
            <h2 className="spec-ref-section__title">Version</h2>
            <p className="spec-ref-section__desc">
              The current version of the OpenCLI Specification is{" "}
              <code className="spec-version-badge">{version}</code>.
            </p>
          </section>

          {/* Root object */}
          <ObjectSection
            id="opencli-specification"
            title="OpenCLI Specification"
            description={schema.description}
            schemaProp={schema as SchemaProp}
          />

          {/* $defs */}
          {Object.entries(defs).map(([name, defSchema]) => (
            <ObjectSection
              key={name}
              id={toId(name)}
              title={name}
              description={defSchema.description}
              schemaProp={defSchema}
            />
          ))}
        </main>
      </div>
    </div>
  );
}
