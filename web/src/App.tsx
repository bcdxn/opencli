import { Navigate, Route, Routes } from "react-router-dom";
import LandingPage from "./pages/LandingPage";
import EditorPage from "./pages/EditorPage";
import SpecPage from "./pages/SpecPage";

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<LandingPage />} />
      <Route path="/editor" element={<EditorPage />} />
      <Route path="/spec" element={<SpecPage />} />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
