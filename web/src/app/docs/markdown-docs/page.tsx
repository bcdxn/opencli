import GenMarkdownDocs from "../../../views/GenMarkdownDocs";

export const metadata = {
  title: "OpenCLI Specification | Docs",
  description:
    "Guides for getting started with the OpenCLI Specification, including markdown and HTML documentation generation.",
  alternates: {
    canonical: "/docs/markdown-docs",
  },
};

export default function Page() {
  return <GenMarkdownDocs />;
}
