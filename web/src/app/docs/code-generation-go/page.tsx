import GeneratingGoCode from "../../../views/GeneratingGoCode";

export const metadata = {
  title: "OpenCLI Specification | Docs",
  description:
    "Learn how to generate framework-specific CLI code from an OpenCLI Specification document.",
  alternates: {
    canonical: "/docs/code-generation-go",
  },
  keywords: [
    "opencli",
    "open cli",
    "opencli specification",
    "generate cli go code",
    "generate cobra cli",
    "generate urfave/cli cli",
  ],
};

export default function Page() {
  return <GeneratingGoCode />;
}
