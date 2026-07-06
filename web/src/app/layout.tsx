import type { Metadata } from "next";
import Script from "next/script";
import Providers from "./providers";
import "../index.css";

export const metadata: Metadata = {
  title: "OpenCLI Specification | OpenAPI for Command Line Interfaces",
  description:
    "OpenCLI is an open specification for describing command line interfaces. Define CLI commands, arguments, options, validation rules, generate consistent documentation and deterministic CLI framework boilerplate code from machine-readable specification files.",
  robots: "index,follow,max-image-preview:large",
  openGraph: {
    type: "website",
    title: "OpenCLI Specification",
    description:
      "OpenCLI is an open specification for command line interfaces - OpenAPI for CLIs. Define commands, arguments, and options in a machine-readable format, validate specifications, and generate consistent documentation and boilerplate code.",
    url: "https://opencli.dev/",
    images: [
      {
        url: "https://opencli.dev/img/opengraph-image.png",
        width: 1200,
        height: 630,
        type: "image/png",
      },
    ],
  },
  twitter: {
    title: "OpenCLI Specification",
    description: "OpenAPI for command line interfaces.",
    card: "summary_large_image",
    images: ["https://opencli.dev/img/twitter-image.png"],
  },
};

const structuredData = {
  "@context": "https://schema.org",
  "@graph": [
    {
      "@type": "TechArticle",
      "@id": "https://opencli.dev",
      headline: "The OpenCLI Specification (OCS) Standard",
      description:
        "A standard, language-agnostic interface for describing for command line interfaces.",
      inLanguage: "en",
      alternateLocales: ["zh-CN"],
      metadataBase: new URL("https://opencli.dev"),
      proficiencyLevel: "Expert",
      about: [
        {
          "@type": "SoftwareApplication",
          "@id": "https://opencli.dev",
          name: "OpenCLI CLI",
          applicationCategory: "DeveloperApplication",
          operatingSystem: "Linux, macOS, Windows",
          description:
            "Terminal-based console applications and command-line tools.",
        },
      ],
      mainEntity: {
        "@type": "SoftwareSourceCode",
        "@id": "https://opencli.dev",
        name: "CLI Specification Schema & Tools",
        codeRepository: "https://github.com/bcdxn/opencli",
        programmingLanguage: ["JSON", "YAML", "GO"],
        license: "https://spdx.org/licenses/MIT",
        runtimePlatform: "Cross-platform shell environments",
      },
    },
  ],
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
        <link rel="sitemap" href="/sitemap.xml" />
        {/* wasm_exec.js must load before the app bundle to provide the Go class */}
        <Script src="/wasm_exec.js" strategy="beforeInteractive" />
        <Script src="/ocli-docs.js" strategy="beforeInteractive" />
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData) }}
        />
      </head>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
