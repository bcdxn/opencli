"use client";

import { useEffect } from "react";
import SiteHeader from "../components/SiteHeader";
import "./DocsPage.css";

export default function DocsPage() {
  useEffect(() => {
    if (window.OcliDocs) {
      window.OcliDocs({ containerId: "docs" });
    }
  });
  return (
    <div className="docs-page">
      <SiteHeader />
      <main className="docs-main">
        <div id="docs"></div>
      </main>
    </div>
  );
}
