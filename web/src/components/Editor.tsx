"use client";

import { useEffect, useRef, useState } from "react";
import { EditorView, basicSetup } from "codemirror";
import { Compartment } from "@codemirror/state";
import { yaml } from "@codemirror/lang-yaml";
import { json } from "@codemirror/lang-json";
import { validateOCS } from "../wasm/client";
import type { ValidationError } from "../wasm/types";
import { useI18n } from "../i18n";
import "./Editor.css";

interface EditorProps {
  content: string;
  format: "yaml" | "json";
  wasmReady: boolean;
  onContentChange: (content: string) => void;
  onFormatChange: (format: "yaml" | "json") => void;
}

export default function Editor({
  content,
  format,
  wasmReady,
  onContentChange,
  onFormatChange,
}: EditorProps) {
  const { t } = useI18n();
  const editorRef = useRef<HTMLDivElement>(null);
  const viewRef = useRef<EditorView | null>(null);
  const langCompartmentRef = useRef<Compartment | null>(null);
  // Keep a stable ref to onContentChange to avoid stale closures in the editor listener
  const onContentChangeRef = useRef(onContentChange);
  useEffect(() => {
    onContentChangeRef.current = onContentChange;
  }, [onContentChange]);

  const [errors, setErrors] = useState<ValidationError[]>([]);
  const [valid, setValid] = useState<boolean | null>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // Create the CodeMirror editor once on mount
  useEffect(() => {
    if (!editorRef.current) return;

    const langCompartment = new Compartment();
    langCompartmentRef.current = langCompartment;

    const view = new EditorView({
      doc: "",
      extensions: [
        basicSetup,
        langCompartment.of(yaml()),
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            onContentChangeRef.current(update.state.doc.toString());
          }
        }),
      ],
      parent: editorRef.current,
    });

    viewRef.current = view;

    return () => {
      view.destroy();
      viewRef.current = null;
      langCompartmentRef.current = null;
    };
  }, []);

  // Update language extension when format changes
  useEffect(() => {
    if (!viewRef.current || !langCompartmentRef.current) return;
    viewRef.current.dispatch({
      effects: langCompartmentRef.current.reconfigure(
        format === "yaml" ? yaml() : json(),
      ),
    });
  }, [format]);

  // Sync external content changes into the editor (e.g., initial petstore YAML load)
  useEffect(() => {
    if (!viewRef.current) return;
    const currentDoc = viewRef.current.state.doc.toString();
    if (currentDoc !== content) {
      viewRef.current.dispatch({
        changes: { from: 0, to: currentDoc.length, insert: content },
      });
    }
  }, [content]);

  // Debounced validation whenever content, format, or wasmReady changes
  useEffect(() => {
    if (!wasmReady) return;
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => {
      const result = validateOCS(content, format);
      setValid(result.valid);
      setErrors(result.errors);
    }, 300);
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [content, format, wasmReady]);

  return (
    <div className="editor-panel">
      <div className="editor-header">
        <span>{t("editor.format")}</span>
        <button
          className={`toggle-btn${format === "yaml" ? " active" : ""}`}
          onClick={() => onFormatChange("yaml")}
        >
          YAML
        </button>
        <button
          className={`toggle-btn${format === "json" ? " active" : ""}`}
          onClick={() => onFormatChange("json")}
        >
          JSON
        </button>
      </div>

      <div className="editor-body" ref={editorRef} />

      <div className="editor-footer">
        {valid === null && (
          <div className="validation-valid">{t("editor.waiting")}</div>
        )}
        {valid === true && (
          <div className="validation-valid">✓ {t("editor.valid")}</div>
        )}
        {valid === false && errors.length > 0 && (
          <div className="validation-errors">
            {errors.map((e, i) => (
              <div key={i} className="validation-error-item">
                {e.path ? (
                  <>
                    <span className="error-path">{e.path}</span>: {e.message}
                  </>
                ) : (
                  e.message
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
