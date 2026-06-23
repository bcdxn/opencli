import { useEffect, useRef, useState } from "react";
import SiteHeader from "../components/SiteHeader";
import { loadWasm } from "../wasm/client";
import Editor from "../components/Editor";
import Preview from "../components/Preview";
import "./EditorPage.css";

export default function EditorPage() {
  const [wasmReady, setWasmReady] = useState(false);
  const [content, setContent] = useState("");
  const [format, setFormat] = useState<"yaml" | "json">("yaml");
  const [editorWidth, setEditorWidth] = useState(50);
  const panelsRef = useRef<HTMLDivElement>(null);
  const isDraggingRef = useRef(false);

  useEffect(() => {
    loadWasm()
      .then(() => setWasmReady(true))
      .catch(console.error);
  }, []);

  useEffect(() => {
    fetch(`${import.meta.env.BASE_URL}petstore-cli.ocs.yaml`)
      .then((r) => r.text())
      .then((text) => setContent(text))
      .catch(console.error);
  }, []);

  const handleDividerMouseDown = () => {
    isDraggingRef.current = true;
  };

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isDraggingRef.current || !panelsRef.current) return;

      const panelsRect = panelsRef.current.getBoundingClientRect();
      const newWidth = ((e.clientX - panelsRect.left) / panelsRect.width) * 100;
      setEditorWidth(Math.max(20, Math.min(80, newWidth)));
    };

    const handleMouseUp = () => {
      isDraggingRef.current = false;
    };

    document.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("mouseup", handleMouseUp);

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  const previewWidth = 100 - editorWidth;

  return (
    <div className="editor-route">
      <SiteHeader />

      <div className="app">
        {!wasmReady && (
          <div className="loading-overlay">
            <span>Loading WASM engine...</span>
          </div>
        )}
        <div className="panels" ref={panelsRef}>
          <div className="panel" style={{ flex: `0 1 ${editorWidth}%` }}>
            <Editor
              content={content}
              format={format}
              wasmReady={wasmReady}
              onContentChange={setContent}
              onFormatChange={setFormat}
            />
          </div>
          <div
            className="panels-divider"
            onMouseDown={handleDividerMouseDown}
          />
          <div className="panel" style={{ flex: `0 1 ${previewWidth}%` }}>
            <Preview
              content={content}
              inputFormat={format}
              wasmReady={wasmReady}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
