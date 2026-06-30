import { useState, useEffect, useRef } from "react";
import ReactMarkdown from "react-markdown";
import rehypeRaw from "rehype-raw";
import remarkGfm from "remark-gfm";
import { generateOCSDocs } from "../wasm/client";
import { useI18n } from "../i18n";
import "./Preview.css";

interface PreviewProps {
  content: string;
  inputFormat: "yaml" | "json";
  wasmReady: boolean;
}

type OutputFormat = "markdown" | "html-page" | "html-embed";
type ViewMode = "rendered" | "raw";

function toInlineScriptText(script: string): string {
  // Prevent accidental early script termination when embedding generated JS in srcDoc.
  return script.replace(/<\/(script)/gi, "<\\/$1");
}

function buildComponentPreviewDoc(componentScript: string): string {
  const safeScript = toInlineScriptText(componentScript);
  return `<!doctype html>
<html lang="en" style="height:100%; padding: 0; margin: 0;">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>OpenCLI Docs Preview</title>
  </head>
  <body style="height:100%; padding: 0; margin: 0;">
    <div id="docs" style="height:100%;"></div>
    <script>${safeScript}</script>
    <script>
      window.OcliDocs({ containerId: "docs" });
    </script>
  </body>
</html>`;
}

export default function Preview({
  content,
  inputFormat,
  wasmReady,
}: PreviewProps) {
  const { t } = useI18n();
  const [outputFormat, setOutputFormat] = useState<OutputFormat>("markdown");
  const [viewMode, setViewMode] = useState<ViewMode>("rendered");
  const [generatedOutput, setGeneratedOutput] = useState("");
  const [genError, setGenError] = useState("");
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (!wasmReady || !content) return;

    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => {
      const result = generateOCSDocs(content, inputFormat, outputFormat);
      if (result.error) {
        setGenError(result.error);
        setGeneratedOutput("");
      } else {
        setGenError("");
        setGeneratedOutput(result.output);
      }
    }, 300);

    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, [content, inputFormat, outputFormat, viewMode, wasmReady]);

  const renderBody = () => {
    if (genError) {
      return <div className="preview-error">{genError}</div>;
    }

    if (!generatedOutput) {
      return null;
    }

    if (outputFormat === "html-page" || outputFormat === "html-embed") {
      if (viewMode === "rendered") {
        const srcDoc =
          outputFormat === "html-embed"
            ? buildComponentPreviewDoc(generatedOutput)
            : generatedOutput;
        return (
          <iframe
            className="preview-iframe"
            srcDoc={srcDoc}
            sandbox="allow-scripts"
            title={t("preview.iframeTitle")}
          />
        );
      }
      // Raw component script source
      return (
        <div className="preview-raw">
          <pre>
            <code>{generatedOutput}</code>
          </pre>
        </div>
      );
    }

    // Markdown
    if (viewMode === "rendered") {
      return (
        <div className="preview-markdown">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            rehypePlugins={[rehypeRaw]}
          >
            {generatedOutput}
          </ReactMarkdown>
        </div>
      );
    }

    // Raw markdown source
    return (
      <div className="preview-raw">
        <pre>
          <code>{generatedOutput}</code>
        </pre>
      </div>
    );
  };

  return (
    <div className="preview-panel">
      <div className="preview-header">
        <div className="preview-header-group">
          <span>{t("preview.output")}</span>
          <button
            className={`toggle-btn${outputFormat === "markdown" ? " active" : ""}`}
            onClick={() => setOutputFormat("markdown")}
          >
            {t("preview.markdown")}
          </button>
          <button
            className={`toggle-btn${outputFormat === "html-page" ? " active" : ""}`}
            onClick={() => setOutputFormat("html-page")}
          >
            {t("preview.htmlPage")}
          </button>
          <button
            className={`toggle-btn${outputFormat === "html-embed" ? " active" : ""}`}
            onClick={() => setOutputFormat("html-embed")}
          >
            {t("preview.htmlEmbed")}
          </button>
        </div>
        <div className="preview-header-group">
          <span>{t("preview.view")}</span>
          <button
            className={`toggle-btn${viewMode === "rendered" ? " active" : ""}`}
            onClick={() => setViewMode("rendered")}
          >
            {t("preview.rendered")}
          </button>
          <button
            className={`toggle-btn${viewMode === "raw" ? " active" : ""}`}
            onClick={() => setViewMode("raw")}
          >
            {t("preview.raw")}
          </button>
        </div>
      </div>

      <div className="preview-body">{renderBody()}</div>
    </div>
  );
}
