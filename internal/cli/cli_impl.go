package cli

import (
	"context"
	"errors"
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

func (Impl) OcliGenerateCli(ctx context.Context, c *urfavecli.Command, args OcliGenerateCliArgs, flags OcliGenerateCliFlags) error {
	// unmarshal the document
	doc, err := oclispec.UnmarshalYAML(args.PathToSpec)
	if err != nil {
		return err
	}

	files, err := oclicode.Generate(doc, oclicode.GoPackage(flags.GoPackage), oclicode.Framework(flags.Framework))
	if err != nil {
		return err
	}

	for _, file := range files {
		err = os.WriteFile(path.Join(args.PathToOutputDir, file.Name), file.Contents, 0644)
		if err != nil {
			return err
		}
	}

	return err
}

func (Impl) OcliGenerateDocs(ctx context.Context, c *urfavecli.Command, args OcliGenerateDocsArgs, flags OcliGenerateDocsFlags) error {
	jsonRE := regexp.MustCompile(`(?i)\.json$`)
	yamlRE := regexp.MustCompile(`(?i)\.yaml$`)

	var doc oclispec.Document
	var err error

	if jsonRE.MatchString(args.PathToSpec) {
		doc, err = oclispec.UnmarshalJSON(args.PathToSpec)
	} else if yamlRE.MatchString(args.PathToSpec) {
		doc, err = oclispec.UnmarshalYAML(args.PathToSpec)
	} else {
		return errors.New("unsupported OpenCLI Document format - must be one of [JSON, YAML]")
	}
	// unmarshal the document
	if err != nil {
		return err
	}

	docs := oclidocs.Generate(doc)

	if flags.Dryrun {
		fmt.Println("--dryrun enabled; skipping write to file")
		fmt.Println("---")
		fmt.Println("docs.gen.md")
		fmt.Println(string(docs))
		fmt.Println("---")
		return nil
	}

	err = os.WriteFile(path.Join(args.PathToOutputDir, "docs.gen."+formatExtension(flags.Format)), docs, 0644)
	return err
}

func (Impl) OcliSpecificationCheck(ctx context.Context, c *urfavecli.Command, args OcliSpecificationCheckArgs) error {
	jsonRE := regexp.MustCompile(`(?i)\.json$`)
	yamlRE := regexp.MustCompile(`(?i)\.yaml$`)

	doc, err := os.ReadFile(args.PathToSpec)
	if err != nil {
		return err
	}

	if jsonRE.MatchString(args.PathToSpec) {
		err = oclispec.ValidateDocumentJSON(doc)
	} else if yamlRE.MatchString(args.PathToSpec) {
		err = oclispec.ValidateDocumentYAML(doc)
	} else {
		return errors.New("unsupported OpenCLI Document format - must be one of [JSON, YAML]")
	}

	if err != nil {
		fmt.Println("OpenCLI Document is invalid ❌")
		return err
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
