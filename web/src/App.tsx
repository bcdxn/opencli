import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import LandingPage from "./pages/LandingPage";
import EditorPage from "./pages/EditorPage";
import DocsPage from "./pages/DocsPage";
import SpecPage from "./pages/SpecPage";
import { useEffect } from "react";

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
    <BrowserRouter basename={basename}>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/editor" element={<EditorPage />} />
        <Route path="/docs" element={<DocsPage />} />
        <Route path="/spec" element={<SpecPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
