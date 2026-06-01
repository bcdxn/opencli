package docs

import (
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/help"
	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
	"github.com/spf13/cobra"
)

func NewCmdDocs(f *cliutils.Factory) *cobra.Command {
	var argsDef map[string]string
	var formatChoice string
	var htmlFlavor string
	var noBadge bool
	var noFooter bool
	var outputDir string
	useLine := "ocli gen docs <arguments> [flags]"

	command := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation from an OpenCLI Spec document",
		Long:  "The generate docs command will generate documentation from an OpenCLI Spec document. You can specify the format of the documentation to be generated using the `--format` flag.",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cliutils.BadUserInput(c, argsDef, useLine, "missing required argument <spec-file>")
			}

			docsArgs := cliutils.OcliGenDocsArgs{
				PathToSpec: args[0],
			}

			docsFlags := cliutils.OcliGenDocsFlags{
				Format:     formatChoice,
				HTMLFlavor: htmlFlavor,
				OutputDir:  outputDir,
				NoBadge:    noBadge,
				NoFooter:   noFooter,
			}

			return f.Actions.OcliGenDocs(docsArgs, docsFlags)
		},
	}

	command.SilenceErrors = true
	command.SilenceUsage = true
	// Arguments
	argsDef = make(map[string]string)
	argsDef["spec-file"] = "path to spec file"
	// Flags
	command.Flags().StringVarP(&formatChoice, "format", "f", "markdown", "docs format (markdown, html)")
	command.Flags().StringVar(&htmlFlavor, "html-flavor", "page", "html flavor (page, component)")
	command.Flags().StringVarP(&outputDir, "out", "o", "./docs", "output directory path")
	command.Flags().BoolVar(&noBadge, "no-badge", false, "do not include the OpenCLI badge")
	command.Flags().BoolVar(&noFooter, "no-footer", false, "do not include the 'created by OpenCLI' footer")
	// Help
	w := f.IOStreams.Out
	command.SetHelpFunc(func(c *cobra.Command, _ []string) {
		help.HelpFunc(w, command, argsDef, useLine)
	})

	return command
}
