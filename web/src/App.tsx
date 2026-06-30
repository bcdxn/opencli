import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import LandingPage from "./pages/LandingPage";
import EditorPage from "./pages/EditorPage";
import DocsPage from "./pages/DocsPage";
import SchemaPage from "./pages/SchemaPage";
import { useEffect } from "react";
import { I18nProvider } from "./i18n";
import "./App.css";

export default function App() {
  const baseUrl = import.meta.env.BASE_URL;
  const basename = baseUrl === "./" ? "/" : baseUrl;

  useEffect(() => {
    let seo = document.getElementById("seo-content");
    if (seo) {
      seo.style.display = "none";
    }
  });

  return (
    <I18nProvider>
      <BrowserRouter basename={basename}>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/editor" element={<EditorPage />} />
          <Route path="/docs" element={<DocsPage />} />
          <Route path="/schema" element={<SchemaPage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </I18nProvider>
  );
}
