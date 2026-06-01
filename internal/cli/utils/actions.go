package cliutils

type CLIActions interface {
	OcliGenDocs(args OcliGenDocsArgs, flags OcliGenDocsFlags) error
	OcliCheck(args OcliCheckArgs, flags OcliCheckFlags) error
}

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

type OcliCheckArgs struct {
	PathToSpec string
}

type OcliCheckFlags struct {
	FailOnErr bool
}
