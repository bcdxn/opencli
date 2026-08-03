import GenManDocs from "../../../views/GenManDocs";

export const metadata = {
  title: "OpenCLI Specification | Man Pages",
  description:
    "Generate roff/troff man pages from your OpenCLI Specification for native Unix man support.",
  alternates: {
    canonical: "/docs/man-pages",
  },
  keywords: [
    "opencli",
    "open cli",
    "opencli specification",
    "generate cli man pages",
  ],
};

export default function Page() {
  return <GenManDocs />;
}
