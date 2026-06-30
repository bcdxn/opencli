import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from "react";

export type Locale = "en" | "zh-CN";

const STORAGE_KEY = "opencli.locale";

const messages: Record<Locale, Record<string, string>> = {
  en: {
    "language.label": "Language",
    "language.english": "English",
    "language.mandarin": "Mandarin (Simplified)",
    "nav.editor": "Editor",
    "nav.docs": "CLI Docs",
    "nav.schema": "Schema",
    "nav.githubAria": "Open the OpenCLI repository on GitHub",

    "landing.kicker": "OPENCLI",
    "landing.title": "Contract-First CLI Design",
    "landing.subtitle": "Define your interface. Automate the rest.",
    "landing.description":
      "OpenCLI is an open document specification for command-line tools. Author verifiable, human-readable CLI definitions. Generate documentation and framework-specific code from a single spec.",
    "landing.cta.editor": "Open Live Editor",
    "landing.cta.docs": "View CLI Docs",
    "landing.tag.contract": "contract-first",
    "landing.tag.framework": "framework-agnostic",
    "landing.tag.agent": "agent-ready",
    "landing.terminalAria": "CLI screenshot placeholder",
    "landing.benefit.contract.title": "Contract-First Development",
    "landing.benefit.contract.copy":
      "Promote stable command contracts and keep implementation details decoupled from CLI frameworks.",
    "landing.benefit.scaling.title": "Tooling That Scales",
    "landing.benefit.scaling.copy":
      "Validate specs, generate docs, and standardize outputs across teams and languages.",
    "landing.benefit.agent.title": "Better Agent Understanding",
    "landing.benefit.agent.copy":
      "Give LLMs and automation tools a structured, explicit model of your CLI surface area.",

    "editor.loadingWasm": "Loading WASM engine...",
    "editor.sample.label": "Sample",
    "editor.sample.en": "Petstore (English)",
    "editor.sample.zh": "Petstore (Mandarin)",
    "editor.format": "Format:",
    "editor.waiting": "Waiting for WASM...",
    "editor.valid": "Valid",

    "preview.output": "Output:",
    "preview.view": "View:",
    "preview.markdown": "Markdown",
    "preview.htmlPage": "HTML Page",
    "preview.htmlEmbed": "HTML Embed",
    "preview.rendered": "Rendered",
    "preview.raw": "Raw",
    "preview.iframeTitle": "HTML Preview",

    "schema.loading": "Loading schema...",
    "schema.failed": "Failed to load spec.schema.json",
    "schema.disclaimerAria": "Validation disclaimer",
    "schema.note.title": "Validation note",
    "schema.note.copy1":
      "This JSON Schema helps enforce structural validity, but not every validation rule can be represented in JSON Schema alone.",
    "schema.note.copy2.prefix": "For full spec validation, run the",
    "schema.note.copy2.link": "OpenCLI CLI",
    "schema.note.copy2.suffix": "check command.",
  },
  "zh-CN": {
    "language.label": "语言",
    "language.english": "English",
    "language.mandarin": "简体中文",
    "nav.editor": "编辑器",
    "nav.docs": "CLI 文档",
    "nav.schema": "模式",
    "nav.githubAria": "在 GitHub 上打开 OpenCLI 仓库",

    "landing.kicker": "OPENCLI",
    "landing.title": "契约优先的 CLI 设计",
    "landing.subtitle": "先定义接口，其余自动化完成。",
    "landing.description":
      "OpenCLI 是一个面向命令行工具的开放文档规范。你可以编写可验证、可读性高的 CLI 定义，并从同一份规范生成文档和框架代码。",
    "landing.cta.editor": "打开在线编辑器",
    "landing.cta.docs": "查看 CLI 文档",
    "landing.tag.contract": "契约优先",
    "landing.tag.framework": "框架无关",
    "landing.tag.agent": "智能体开发",
    "landing.terminalAria": "CLI 截图占位",
    "landing.benefit.contract.title": "契约优先开发",
    "landing.benefit.contract.copy":
      "让命令契约稳定，并将实现细节与 CLI 框架解耦。",
    "landing.benefit.scaling.title": "可扩展工具链",
    "landing.benefit.scaling.copy":
      "校验规范、生成文档，并在跨团队与跨语言场景中统一输出。",
    "landing.benefit.agent.title": "更好的智能体理解",
    "landing.benefit.agent.copy":
      "为 LLM 与自动化工具提供结构化且明确的 CLI 表面模型。",

    "editor.loadingWasm": "正在加载 WASM 引擎...",
    "editor.sample.label": "示例",
    "editor.sample.en": "Petstore（英文）",
    "editor.sample.zh": "Petstore（中文）",
    "editor.format": "格式：",
    "editor.waiting": "等待 WASM 就绪...",
    "editor.valid": "有效",

    "preview.output": "输出：",
    "preview.view": "视图：",
    "preview.markdown": "Markdown",
    "preview.htmlPage": "HTML 页面",
    "preview.htmlEmbed": "HTML 组件",
    "preview.rendered": "渲染",
    "preview.raw": "原始",
    "preview.iframeTitle": "HTML 预览",

    "schema.loading": "正在加载 schema...",
    "schema.failed": "加载 spec.schema.json 失败",
    "schema.disclaimerAria": "校验说明",
    "schema.note.title": "校验说明",
    "schema.note.copy1":
      "这个 JSON Schema 可以帮助校验结构有效性，但并非所有校验规则都能仅用 JSON Schema 表达。",
    "schema.note.copy2.prefix": "如需完整规范校验，请运行",
    "schema.note.copy2.link": "OpenCLI CLI",
    "schema.note.copy2.suffix": "的 check 命令。",
  },
};

interface I18nContextValue {
  locale: Locale;
  setLocale: (locale: Locale) => void;
  t: (key: string) => string;
}

const I18nContext = createContext<I18nContextValue | null>(null);

function detectInitialLocale(): Locale {
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === "en" || stored === "zh-CN") {
    return stored;
  }

  const preferred =
    navigator.languages && navigator.languages.length > 0
      ? navigator.languages[0]
      : navigator.language;
  if (preferred.toLowerCase().startsWith("zh")) {
    return "zh-CN";
  }

  return "en";
}

export function I18nProvider({ children }: PropsWithChildren) {
  const [locale, setLocale] = useState<Locale>(detectInitialLocale);

  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, locale);
    document.documentElement.lang = locale;
  }, [locale]);

  const value = useMemo<I18nContextValue>(() => {
    return {
      locale,
      setLocale,
      t: (key: string) => messages[locale][key] ?? messages.en[key] ?? key,
    };
  }, [locale]);

  return <I18nContext.Provider value={value}>{children}</I18nContext.Provider>;
}

export function useI18n(): I18nContextValue {
  const value = useContext(I18nContext);
  if (!value) {
    throw new Error("useI18n must be used within I18nProvider");
  }
  return value;
}
