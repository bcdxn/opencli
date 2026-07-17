import GenHtmlDocs from "../../../views/GenHtmlDocs";

export const metadata = {
  title: "OpenCLI Specification | Docs",
  description:
    "Guides for getting started with the OpenCLI Specification, including markdown and HTML documentation generation.",
  alternates: {
    canonical: "/docs/html-docs",
  },
};

export default function Page() {
  return <GenHtmlDocs />;
}
