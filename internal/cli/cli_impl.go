package cli

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/bcdxn/opencli/oclicode"
	"github.com/bcdxn/opencli/oclidocs"
	"github.com/bcdxn/opencli/oclispec"
	urfavecli "github.com/urfave/cli/v3"
)

type Impl struct{}

func (Impl) OcliGenerateCli(ctx context.Context, c *urfavecli.Command, flags OcliGenerateCliFlags) error {
	// unmarshal the document
	doc, err := oclispec.UnmarshalYAML(flags.SpecFile)
	if err != nil {
		return urfavecli.Exit(err.Error(), ExitCodeBadUserInputError)
	}

	files, err := oclicode.Generate(
		doc, oclicode.GoPackage(flags.GoPackage),
		oclicode.Framework(flags.Framework),
		oclicode.ModuleType(flags.ModuleType),
	)
	if err != nil {
		return urfavecli.Exit(err.Error(), ExitCodeInternalCliError)
	}

	if flags.Dryrun {
		fmt.Println("--dryrun enabled; skipping write to file")
		for _, f := range files {
			fmt.Println("\n\n========================================================================")
			fmt.Println("FILE: ", f.Name)
			fmt.Println("------------------------------------------------------------------------")
			fmt.Println(string(f.Contents))
		}
		return nil
	}

	for _, file := range files {
		err = os.WriteFile(path.Join(flags.OutputDir, file.Name), file.Contents, 0644)
		if err != nil {
			return urfavecli.Exit(err.Error(), ExitCodeInternalCliError)
		}
	}

	return nil
}

func (Impl) OcliGenerateDocs(ctx context.Context, c *urfavecli.Command, flags OcliGenerateDocsFlags) error {
	jsonRE := regexp.MustCompile(`(?i)\.json$`)
	yamlRE := regexp.MustCompile(`(?i)\.yaml$`)

	var doc oclispec.Document
	var err error

	if jsonRE.MatchString(flags.SpecFile) {
		doc, err = oclispec.UnmarshalJSON(flags.SpecFile)
	} else if yamlRE.MatchString(flags.SpecFile) {
		doc, err = oclispec.UnmarshalYAML(flags.SpecFile)
	} else {
		return urfavecli.Exit("unsupported OpenCLI Document format - must be one of [JSON, YAML]", ExitCodeBadUserInputError)
	}
	// unmarshal the document
	if err != nil {
		return urfavecli.Exit(err.Error(), ExitCodeBadUserInputError)
	}

	files, err := oclidocs.Generate(doc)
	if err != nil {
		return urfavecli.Exit(err.Error(), ExitCodeInternalCliError)
	}

	if flags.Dryrun {
		fmt.Println("--dryrun enabled; skipping write to file")
		for _, f := range files {
			fmt.Println("\n\n========================================================================")
			fmt.Println("FILE: ", f.Name)
			fmt.Println("------------------------------------------------------------------------")
			fmt.Println(string(f.Contents))
		}
		return nil
	}

	for _, file := range files {
		err = os.WriteFile(path.Join(flags.OutputDir, file.Name), file.Contents, 0644)
		if err != nil {
			return urfavecli.Exit(err.Error(), ExitCodeInternalCliError)
		}
	}

	return nil
}

func (Impl) OcliSpecificationCheck(ctx context.Context, c *urfavecli.Command, args OcliSpecificationCheckArgs) error {
	jsonRE := regexp.MustCompile(`(?i)\.json$`)
	yamlRE := regexp.MustCompile(`(?i)\.yaml$`)

	doc, err := os.ReadFile(args.PathToSpec)
	if err != nil {
		return urfavecli.Exit(err.Error(), ExitCodeInternalCliError)
	}

	if jsonRE.MatchString(args.PathToSpec) {
		err = oclispec.ValidateDocumentJSON(doc)
	} else if yamlRE.MatchString(args.PathToSpec) {
		err = oclispec.ValidateDocumentYAML(doc)
	} else {
		return urfavecli.Exit("unsupported OpenCLI Document format - must be one of [JSON, YAML]", ExitCodeBadUserInputError)
	}

	if err != nil {
		fmt.Println("OpenCLI Document is invalid ❌")
		return urfavecli.Exit(err.Error(), ExitCodeBadUserInputError)
	}

	fmt.Println("OpenCLI Document is valid ✅")
	return nil
}

func (Impl) OcliSpecificationVersions(ctx context.Context, c *urfavecli.Command) error {
	versions := oclispec.Versions()

	fmt.Print("Supported Versions:\n\n")
	for _, v := range versions {
		fmt.Printf("- %s\n", v)
	}

	return nil
}

func formatExtension(format string) string {
	switch format {
	case "markdown":
		return "md"
	default:
		return "md"
	}
}
