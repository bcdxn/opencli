package docs

import (
	"github.com/bcdxn/opencli/internal/cli/gencli"
	"github.com/bcdxn/opencli/spec"
	"github.com/spf13/cobra"
)

func NewCmdDocs(a gencli.ActionsInterface) *cobra.Command {
	// Flag vars
	var formatChoice string
	var htmlFlavor string
	var noBadge bool
	var noFooter bool
	var outputDir string
	// Command def
	command := &cobra.Command{
		Use:   "docs",
		Short: "Generate documentation from an OpenCLI Spec document",
		Long:  "The generate docs command will generate documentation from an OpenCLI Spec document. You can specify the format of the documentation to be generated using the `--format` flag.",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return gencli.BadUserInput(c, "missing required argument <spec-file>")
			}

			docsArgs := gencli.OcliGenDocsArgs{
				PathToSpec: args[0],
			}

			docsFlags := gencli.OcliGenDocsFlags{
				Format:     formatChoice,
				HTMLFlavor: htmlFlavor,
				OutputDir:  outputDir,
				NoBadge:    noBadge,
				NoFooter:   noFooter,
			}

			return a.OcliGenDocs(docsArgs, docsFlags)
		},
	}
	// Command config
	command.SilenceErrors = true
	command.SilenceUsage = true
	// Flag bindings
	command.Flags().StringVarP(&formatChoice, "format", "f", "markdown", "docs format (markdown, html)")
	command.Flags().StringVar(&htmlFlavor, "html-flavor", "page", "html flavor (page, embed)")
	command.Flags().StringVarP(&outputDir, "out", "o", "./docs", "output directory path")
	command.Flags().BoolVar(&noBadge, "no-badge", false, "do not include the OpenCLI badge")
	command.Flags().BoolVar(&noFooter, "no-footer", false, "do not include the 'created by OpenCLI' footer")
	// Help
	command.SetHelpFunc(func(c *cobra.Command, args []string) {
		a.HelpFunc(&spec.CommandItem{})
	})

	return command
}
