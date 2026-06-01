package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/gen"
	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
	"github.com/bcdxn/opencli/validate"
)

// Actions implements the cliutils Actions interface and can be passed via the cliutils.Factory
type Actions struct {
	IOS *cliutils.IOStreams
}

func (a Actions) OcliGenDocs(args cliutils.OcliGenDocsArgs, flags cliutils.OcliGenDocsFlags) error {
	// Validate file exists and is not a directory
	info, err := os.Stat(args.PathToSpec)
	if err != nil {
		if os.IsNotExist(err) {
			return cliutils.NewValidationError(fmt.Sprintf("file not found: %s", args.PathToSpec))
		}
		return cliutils.NewValidationError(fmt.Sprintf("cannot access file: %s (%v)", args.PathToSpec, err))
	}
	if info.IsDir() {
		return cliutils.NewValidationError(fmt.Sprintf("path is a directory, not a file: %s", args.PathToSpec))
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
		return cliutils.NewValidationError(fmt.Sprintf("unsupported spec format: %s (only .json, .yaml, .yml are supported)", ext))
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
		return cliutils.NewValidationError(fmt.Sprintf("unsupported docs format: %q (supported: %s)", flags.Format, strings.Join(keys, ", ")))
	}

	fmt.Fprintf(a.IOS.Out, "\n→ Reading spec:       %s\n", args.PathToSpec)

	// Read and parse the spec document
	data, err := os.ReadFile(args.PathToSpec)
	if err != nil {
		return cliutils.NewValidationError(fmt.Sprintf("cannot read file: %s (%v)", args.PathToSpec, err))
	}
	doc, err := decode(data)
	if err != nil {
		return cliutils.NewValidationError(fmt.Sprintf("failed to parse spec: %v", err))
	}

	fmt.Fprintf(a.IOS.Out, "→ Generating docs:    format=%s, output=%s\n", flags.Format, flags.OutputDir)

	// Ensure output directory exists
	if err := os.MkdirAll(flags.OutputDir, 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", flags.OutputDir, err)
	}

	opts := make([]gen.GenDocsOption, 0)

	opts = append(opts, gen.DocsWithFormat(meta.format))

	if meta.format == gen.HTML {
		switch strings.ToLower(flags.HTMLFlavor) {
		case "component", "embedded", "embeddable":
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

	// Write output file (strip spec extension, append doc format extension)
	baseName := strings.TrimSuffix(filepath.Base(args.PathToSpec), filepath.Ext(args.PathToSpec))
	outFile := filepath.Join(flags.OutputDir, baseName+meta.fileExt)
	if err := os.WriteFile(outFile, output, 0644); err != nil {
		return fmt.Errorf("failed to write output file %s: %w", outFile, err)
	}

	fmt.Fprintf(a.IOS.Out, "✓ Documentation written to: %s\n\n", outFile)

	return nil
}

// OcliCheck implements the `ocli check` command and uses the `validate` package to validate/check the specified document.
func (a Actions) OcliCheck(args cliutils.OcliCheckArgs, flags cliutils.OcliCheckFlags) error {
	// Verify file exists and is readable
	info, err := os.Stat(args.PathToSpec)
	if err != nil {
		if os.IsNotExist(err) {
			return cliutils.NewValidationError(fmt.Sprintf("file not found: %s", args.PathToSpec))
		}
		return cliutils.NewValidationError(fmt.Sprintf("cannot access file: %s (%v)", args.PathToSpec, err))
	}

	if info.IsDir() {
		return cliutils.NewValidationError(fmt.Sprintf("path is a directory, not a file: %s", args.PathToSpec))
	}

	// Read file content
	data, err := os.ReadFile(args.PathToSpec)
	if err != nil {
		return cliutils.NewValidationError(fmt.Sprintf("cannot read file: %s (%v)", args.PathToSpec, err))
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
		return cliutils.NewValidationError(fmt.Sprintf("unsupported file format: %s (only .json, .yaml, .yml are supported)", ext))
	}

	// Output results
	fmt.Fprintf(a.IOS.Out, "\n✓ Checking %s\n", args.PathToSpec)
	fmt.Fprintf(a.IOS.Out, "  Format: %s\n\n", strings.TrimPrefix(ext, "."))

	if validationErr != nil {
		fmt.Fprintf(a.IOS.Out, "✗ Validation failed:\n")
		for _, line := range strings.Split(validationErr.Error(), "\n") {
			fmt.Fprintf(a.IOS.Out, "  %s\n", line)
		}
		fmt.Fprintln(a.IOS.Out)

		if flags.FailOnErr {
			return cliutils.NewValidationError("document validation failed")
		}
	} else {
		fmt.Fprintf(a.IOS.Out, "✓ Document is valid\n\n")
	}

	return nil
}
