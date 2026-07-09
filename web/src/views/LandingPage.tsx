"use client";

import Link from "next/link";
import SiteHeader from "../components/SiteHeader";
import "./LandingPage.css";
import SiteFooter from "../components/SiteFooter";
import React from "react";
import { useI18n } from "../i18n";

const yamlSampleEn = `opencliVersion: 1.0.0-alpha.12

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

const yamlSampleCh = `opencliVersion: 1.0.0-alpha.12

info:
  title: 寒暄
  summary: 亲切的问候
  version: 1.0.0
  binary: pleasantries

commands:
  pleasantries {command} <args> [flags]:
    group: true

  pleasantries greet <args> [flags]:
    summary: "问好"
    args:
      - name: "name"
        required: true
        type: "string"
    flags:
      - name: "language"
        type: "string"
        default: "chinese"`;

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
  const { t, locale } = useI18n();
  const activeYamlSample = locale === "zh-CN" ? yamlSampleCh : yamlSampleEn;

  return (
    <div className="landing-route">
      <SiteHeader />

      <main className="landing-main">
        <section className="hero-shell">
          <div className="hero-copy">
            <div className="hero-kicker">{t("landing.kicker")}</div>
            <h1>{t("landing.title")}</h1>
            <h2>{t("landing.subtitle")}</h2>
            <p>{t("landing.description")}</p>

            <div className="hero-actions">
              <Link className="button primary" href="/editor">
                {t("landing.cta.editor")}
              </Link>
              <Link className="button secondary" href="/docs">
                {t("landing.cta.docs")}
              </Link>
            </div>

            <div className="hero-tags">
              <span>{t("landing.tag.contract")}</span>
              <span>{t("landing.tag.framework")}</span>
              <span>{t("landing.tag.agent")}</span>
            </div>
          </div>

          <section
            className="terminal-shot"
            aria-label={t("landing.terminalAria")}
          >
            <div className="terminal-titlebar">
              <span />
              <span />
              <span />
            </div>
            <div className="terminal-body">
              <pre>
                <code>
                  {activeYamlSample
                    .split("\n")
                    .map((line, i) => tokenizeYamlLine(line, i))}
                </code>
              </pre>
            </div>
          </section>
        </section>

        <section className="benefits">
          <article>
            <h2>{t("landing.benefit.contract.title")}</h2>
            <p>{t("landing.benefit.contract.copy")}</p>
          </article>
          <article>
            <h2>{t("landing.benefit.scaling.title")}</h2>
            <p>{t("landing.benefit.scaling.copy")}</p>
          </article>
          <article>
            <h2>{t("landing.benefit.agent.title")}</h2>
            <p>{t("landing.benefit.agent.copy")}</p>
          </article>
        </section>
      </main>
      <SiteFooter />
    </div>
  );
}
