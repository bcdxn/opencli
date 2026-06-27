package gencli

import (
	"github.com/bcdxn/opencli/spec"
)

type ActionsInterface interface {
	OcliGenDocs(args OcliGenDocsArgs, flags OcliGenDocsFlags) error
	OcliCheck(args OcliCheckArgs, flags OcliCheckFlags) error
	HelpFunc(cmd *spec.CommandItem)
	UsageFunc(cmd *spec.CommandItem) error
	IOStreams() IOStreams
}
