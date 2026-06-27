package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/gen"
	"github.com/bcdxn/opencli/internal/cli/gencli"
	"github.com/bcdxn/opencli/spec"
	"github.com/bcdxn/opencli/validate"
)

// Ensure we conform to the generated ActionsInterface
var _ gencli.ActionsInterface = (*Actions)(nil)

func NewActions(version string) *Actions {
	return &Actions{
		IOS: gencli.DefaultIOS(),
	}
}

// Actions implements the gencli Actions interface and can be passed via the gencli.Factory
type Actions struct {
	IOS gencli.IOStreams
}

func (a Actions) OcliGenDocs(args gencli.OcliGenDocsArgs, flags gencli.OcliGenDocsFlags) error {
	// Validate file exists and is not a directory
	info, err := os.Stat(args.PathToSpec)
	if err != nil {
		if os.IsNotExist(err) {
			return gencli.NewValidationError(fmt.Sprintf("file not found: %s", args.PathToSpec))
		}
		return gencli.NewValidationError(fmt.Sprintf("cannot access file: %s (%v)", args.PathToSpec, err))
	}
	if info.IsDir() {
		return gencli.NewValidationError(fmt.Sprintf("path is a directory, not a file: %s", args.PathToSpec))
	}

	// Select decoder by file extension
	ext := strings.ToLower(filepath.Ext(args.PathToSpec))
	decode := codec.UnmarshalYAML
	switch ext {
	case ".json":
		decode = codec.UnmarshalJSON
	case ".yaml", ".yml":
		// default
	default:
		return gencli.NewValidationError(fmt.Sprintf("unsupported spec format: %s (only .json, .yaml, .yml are supported)", ext))
	}

	// Map the --format flag to a DocFormat; this map is the extension point for new formats
	type docFormatMeta struct {
		format  gen.DocFormat
		fileExt string
	}
	supportedFormats := map[string]docFormatMeta{
		"markdown": {gen.Markdown, ".md"},
		"html":     {gen.HTML, ".html"},
		"man":      {gen.ManPage, ".1"},
	}
	meta, ok := supportedFormats[strings.ToLower(flags.Format)]
	if !ok {
		keys := make([]string, 0, len(supportedFormats))
		for k := range supportedFormats {
			keys = append(keys, k)
		}
		return gencli.NewValidationError(fmt.Sprintf("unsupported docs format: %q (supported: %s)", flags.Format, strings.Join(keys, ", ")))
	}

	stdout := a.IOS.Out()
	fmt.Fprintf(stdout, "\n→ Reading spec:       %s\n", args.PathToSpec)

	// Read and parse the spec document
	data, err := os.ReadFile(args.PathToSpec)
	if err != nil {
		return gencli.NewValidationError(fmt.Sprintf("cannot read file: %s (%v)", args.PathToSpec, err))
	}
	doc, err := decode(data)
	if err != nil {
		return gencli.NewValidationError(fmt.Sprintf("failed to parse spec: %v", err))
	}

	fmt.Fprintf(stdout, "→ Generating docs:    format=%s, output=%s\n", flags.Format, flags.OutputDir)

	// Ensure output directory exists
	if err := os.MkdirAll(flags.OutputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", flags.OutputDir, err)
	}

	opts := make([]gen.GenDocsOption, 0)

	opts = append(opts, gen.DocsWithFormat(meta.format))

	if meta.format == gen.HTML {
		switch strings.ToLower(flags.HTMLFlavor) {
		case "embed":
			opts = append(opts, gen.DocsWithHTMLFlavor(gen.EmbeddableComponent))
		default:
			opts = append(opts, gen.DocsWithHTMLFlavor(gen.StandalonePage))
		}
	}

	if flags.NoBadge {
		opts = append(opts, gen.DocsWithoutBadge())
	}

	if flags.NoFooter {
		opts = append(opts, gen.DocsWithoutFooter())
	}

	// Generate documentation
	output, err := gen.Docs(doc, opts...)
	if err != nil {
		return fmt.Errorf("failed to generate %s docs: %w", flags.Format, err)
	}

	// Resolve output file path. HTML embed flavor is emitted as an embeddable JS asset.
	baseName := strings.TrimSuffix(filepath.Base(args.PathToSpec), filepath.Ext(args.PathToSpec))
	outFile := filepath.Join(flags.OutputDir, baseName+meta.fileExt)
	if meta.format == gen.HTML {
		switch strings.ToLower(flags.HTMLFlavor) {
		case "embed":
			outFile = filepath.Join(flags.OutputDir, "ocli-docs.js")
		}
	}

	if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", filepath.Dir(outFile), err)
	}

	if err := os.WriteFile(outFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outFile, err)
	}

	fmt.Fprintf(stdout, "✓ Documentation written to: %s\n\n", outFile)

	return nil
}

// OcliCheck implements the `ocli check` command and uses the `validate` package to validate/check the specified document.
func (a Actions) OcliCheck(args gencli.OcliCheckArgs, flags gencli.OcliCheckFlags) error {
	// Verify file exists and is readable
	info, err := os.Stat(args.PathToSpec)
	if err != nil {
		if os.IsNotExist(err) {
			return gencli.NewValidationError(fmt.Sprintf("file not found: %s", args.PathToSpec))
		}
		return gencli.NewValidationError(fmt.Sprintf("cannot access file: %s (%v)", args.PathToSpec, err))
	}

	if info.IsDir() {
		return gencli.NewValidationError(fmt.Sprintf("path is a directory, not a file: %s", args.PathToSpec))
	}

	// Read file content
	data, err := os.ReadFile(args.PathToSpec)
	if err != nil {
		return gencli.NewValidationError(fmt.Sprintf("cannot read file: %s (%v)", args.PathToSpec, err))
	}

	// Determine format by file extension
	ext := strings.ToLower(filepath.Ext(args.PathToSpec))
	var validationErr error

	switch ext {
	case ".json":
		validationErr = validate.ValidateJSON(data)
	case ".yaml", ".yml":
		validationErr = validate.ValidateYAML(data)
	default:
		return gencli.NewValidationError(fmt.Sprintf("unsupported file format: %s (only .json, .yaml, .yml are supported)", ext))
	}

	// Output results
	stdout := a.IOS.Out()
	fmt.Fprintf(stdout, "\n✓ Checking %s\n", args.PathToSpec)
	fmt.Fprintf(stdout, "  Format: %s\n\n", strings.TrimPrefix(ext, "."))

	if validationErr != nil {
		fmt.Fprintf(stdout, "✗ Validation failed:\n")
		for _, line := range strings.Split(validationErr.Error(), "\n") {
			fmt.Fprintf(stdout, "  %s\n", line)
		}
		fmt.Fprintln(stdout)

		if flags.FailOnErr {
			return gencli.NewValidationError("document validation failed")
		}
	} else {
		fmt.Fprintf(stdout, "✓ Document is valid\n\n")
	}

	return nil
}

func (a Actions) HelpFunc(cmd *spec.CommandItem) {
	gencli.DefaultHelpFunc(a, cmd)
}

func (a Actions) UsageFunc(cmd *spec.CommandItem) error {
	return gencli.DefaultUsageFunc(a, cmd)
}

func (a Actions) IOStreams() gencli.IOStreams {
	return a.IOS
}
