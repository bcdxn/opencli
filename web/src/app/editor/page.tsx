import EditorPage from "../../views/EditorPage";

export const metadata = {
  title: "OpenCLI Specification | Document Editor and Preview",
  description:
    "Create, edit, and validate OpenCLI specification documents in your browser. Preview generated documentation live.",
  alternates: {
    canonical: "https://opencli.dev/editor",
  },
};

export default function Page() {
  return <EditorPage />;
}
