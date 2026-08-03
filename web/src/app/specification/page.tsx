import SpecPage from "../../views/SpecPage";

export const metadata = {
  title: "OpenCLI Specification | Specification Schema",
  description:
    "Details the OpenCLI specification schema, and supported fields.",
  alternates: {
    canonical: "/specification",
  },
  keywords: ["opencli", "open cli", "opencli specification"],
};

export default function Page() {
  return <SpecPage />;
}
