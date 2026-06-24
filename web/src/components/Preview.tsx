import { useState, useEffect, useRef } from "react";
import ReactMarkdown from "react-markdown";
import rehypeRaw from "rehype-raw";
import remarkGfm from "remark-gfm";
import { generateOCSDocs } from "../wasm/client";
import "./Preview.css";

interface PreviewProps {
  content: string;
  inputFormat: "yaml" | "json";
  wasmReady: boolean;
}

type OutputFormat = "markdown" | "html";
type ViewMode = "rendered" | "raw";

function toInlineScriptText(script: string): string {
  // Prevent accidental early script termination when embedding generated JS in srcDoc.
  return script.replace(/<\/(script)/gi, "<\\/$1");
}

function buildComponentPreviewDoc(componentScript: string): string {
  const safeScript = toInlineScriptText(componentScript);
  return `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>OpenCLI Docs Preview</title>
  </head>
  <body>
    <div id="docs"></div>
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
  const [outputFormat, setOutputFormat] = useState<OutputFormat>("markdown");
  const [viewMode, setViewMode] = useState<ViewMode>("rendered");
  const [generatedOutput, setGeneratedOutput] = useState("");
  const [genError, setGenError] = useState("");
  const debounceRef = useRef<ReturnType<typeof setTimeout>>();

  // HTML output now uses embed flavor which emits an embeddable JS initializer.
  const htmlFlavor = outputFormat === "html" ? "embed" : "page";

  useEffect(() => {
    if (!wasmReady || !content) return;

    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => {
      const result = generateOCSDocs(
        content,
        inputFormat,
        outputFormat,
        htmlFlavor,
      );
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
  }, [content, inputFormat, outputFormat, viewMode, wasmReady, htmlFlavor]);

  const renderBody = () => {
    if (genError) {
      return <div className="preview-error">{genError}</div>;
    }

    if (!generatedOutput) {
      return null;
    }

    if (outputFormat === "html") {
      if (viewMode === "rendered") {
        const srcDoc = buildComponentPreviewDoc(generatedOutput);
        return (
          <iframe
            className="preview-iframe"
            srcDoc={srcDoc}
            sandbox="allow-scripts"
            title="HTML Preview"
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
          <span>Output:</span>
          <button
            className={`toggle-btn${outputFormat === "markdown" ? " active" : ""}`}
            onClick={() => setOutputFormat("markdown")}
          >
            Markdown
          </button>
          <button
            className={`toggle-btn${outputFormat === "html" ? " active" : ""}`}
            onClick={() => setOutputFormat("html")}
          >
            HTML
          </button>
        </div>
        <div className="preview-header-group">
          <span>View:</span>
          <button
            className={`toggle-btn${viewMode === "rendered" ? " active" : ""}`}
            onClick={() => setViewMode("rendered")}
          >
            Rendered
          </button>
          <button
            className={`toggle-btn${viewMode === "raw" ? " active" : ""}`}
            onClick={() => setViewMode("raw")}
          >
            Raw
          </button>
        </div>
      </div>

      <div className="preview-body">{renderBody()}</div>
    </div>
  );
}
