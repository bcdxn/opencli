import { HashRouter, Navigate, Route, Routes } from "react-router-dom";
import LandingPage from "./pages/LandingPage";
import EditorPage from "./pages/EditorPage";
import SpecPage from "./pages/SpecPage";

export default function App() {
  const baseUrl = import.meta.env.BASE_URL;
  const basename = baseUrl === "./" ? "/" : baseUrl;

  return (
    <HashRouter basename={basename}>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/editor" element={<EditorPage />} />
        <Route path="/spec" element={<SpecPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </HashRouter>
  );
}
