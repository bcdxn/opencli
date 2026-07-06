"use client";

import Link from "next/link";
import { useI18n } from "../i18n";
import "./SiteFooter.css";

export default function SiteFooter() {
  const { locale, setLocale } = useI18n();

  return (
    <footer className="site-footer" aria-label="Language selector">
      <div className="site-footer-inner">
        <Link className="footer-brand" href="/">
          OpenCLI
        </Link>
        <span className="footer-sep" aria-hidden="true">
          |
        </span>
        <button
          type="button"
          className={`footer-locale${locale === "en" ? " active" : ""}`}
          onClick={() => setLocale("en")}
        >
          en
        </button>
        <span className="footer-sep" aria-hidden="true">
          |
        </span>
        <button
          type="button"
          className={`footer-locale${locale === "zh-CN" ? " active" : ""}`}
          onClick={() => setLocale("zh-CN")}
        >
          简体中文
        </button>
      </div>
    </footer>
  );
}
