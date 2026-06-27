package gencli

// Params
type OcliGenDocsArgs struct {
	PathToSpec string
}

type OcliGenDocsFlags struct {
	Format     string
	HTMLFlavor string
	OutputDir  string
	NoBadge    bool
	NoFooter   bool
}
