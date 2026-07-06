"use client";

import { useEffect, useRef, useState } from "react";
import SiteHeader from "../components/SiteHeader";
import { loadWasm } from "../wasm/client";
import Editor from "../components/Editor";
import Preview from "../components/Preview";
import { useI18n, type Locale } from "../i18n";
import "./EditorPage.css";

const sampleFiles: Record<Locale, string> = {
  en: "petstore-cli.ocs.yaml",
  "zh-CN": "petstore-cli.zh-cn.ocs.yaml",
};

export default function EditorPage() {
  const { locale, t } = useI18n();
  const [wasmReady, setWasmReady] = useState(false);
  const [content, setContent] = useState("");
  const [format, setFormat] = useState<"yaml" | "json">("yaml");
  const [sampleLocale, setSampleLocale] = useState<Locale>(locale);
  const [editorWidth, setEditorWidth] = useState(50);
  const panelsRef = useRef<HTMLDivElement>(null);
  const isDraggingRef = useRef(false);

  useEffect(() => {
    loadWasm()
      .then(() => setWasmReady(true))
      .catch(console.error);
  }, []);

  useEffect(() => {
    setSampleLocale(locale);
  }, [locale]);

  useEffect(() => {
    fetch(`/${sampleFiles[sampleLocale]}`)
      .then((r) => r.text())
      .then((text) => {
        setContent(text);
        setFormat("yaml");
      })
      .catch(console.error);
  }, [sampleLocale]);

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
        <div className="editor-sample-bar">
          <label htmlFor="sample-language">{t("editor.sample.label")}</label>
          <select
            id="sample-language"
            value={sampleLocale}
            onChange={(event) => setSampleLocale(event.target.value as Locale)}
          >
            <option value="en">{t("editor.sample.en")}</option>
            <option value="zh-CN">{t("editor.sample.zh")}</option>
          </select>
        </div>
        {!wasmReady && (
          <div className="loading-overlay">
            <span>{t("editor.loadingWasm")}</span>
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
