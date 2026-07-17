"use client";

import { useEffect } from "react";
import SiteHeader from "../components/SiteHeader";
import "./ReferencePage.css";

export default function ReferencePage() {
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
