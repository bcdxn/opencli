// This page is a hold over from original site structure. CLI reference docs have officially moved
// to /reference, but we'll keep this here so we don't tank SEO by removing the indexed page.
import ReferencePage from "../../views/ReferencePage";

export const metadata = {
  title: "OpenCLI Specification | CLI Docs",
  description: "Learn how to use the OpenCLI Specification tooling.",
  alternates: {
    canonical: "/docs",
  },
};

export default function Page() {
  return <ReferencePage />;
}
