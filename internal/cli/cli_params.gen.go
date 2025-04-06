// Code generated by ocli-codegen DO NOT EDIT.
// Package cli provides primitives to interact with the Open CLI Spec-Compliant CLI.

package cli

/* CLI argument types
------------------------------------------------------------------------------------------------- */

// OcliSpecificationCheckArgs holds the parsed arguments that will be injected into the command handler implementation.
type OcliSpecificationCheckArgs struct {
  PathToSpec string
}

/* CLI flag types
------------------------------------------------------------------------------------------------- */

// OcliGenerateCliFlags holds the parsed flags that will be injected into the command handler implementation.
type OcliGenerateCliFlags struct {
  SpecFile string
  OutputDir string
  Framework string
  GoPackage string
  ModuleType string
  Dryrun bool
}

// OcliGenerateDocsFlags holds the parsed flags that will be injected into the command handler implementation.
type OcliGenerateDocsFlags struct {
  SpecFile string
  OutputDir string
  Format string
  Footer bool
  Dryrun bool
}

func validateChoices(choices []string, val string) bool {
	for _, choice := range choices {
		if choice == val {
			return true
		}
	}

	return false
}
