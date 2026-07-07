package cli

import (
	"context"
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
		IOS:     gencli.DefaultIOS(),
		version: version,
	}
}

// Actions implements the gencli Actions interface and can be passed via the gencli.Factory
type Actions struct {
	IOS     gencli.IOStreams
	version string
}

func (a Actions) OcliGenDocs(_ context.Context, args gencli.OcliGenDocsArgs, flags gencli.OcliGenDocsFlags) error {
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
		"markdown":   {gen.Markdown, ".md"},
		"html-page":  {gen.HTML_PAGE, ".html"},
		"html-embed": {gen.HTML_EMBED, ".js"},
		"man":        {gen.ManPage, ".1"},
	}
	meta, ok := supportedFormats[strings.ToLower(string(flags.Format))]
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

	fmt.Fprintf(stdout, "→ Generating docs:    format=%s, output=%s\n", flags.Format, flags.Out)

	// Ensure output directory exists
	if err := os.MkdirAll(flags.Out, 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", flags.Out, err)
	}

	opts := make([]gen.GenDocsOption, 0)

	opts = append(opts, gen.DocsWithFormat(meta.format))

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

	// Resolve output file path. HTML embed emits a shared embeddable JS asset.
	baseName := strings.TrimSuffix(filepath.Base(args.PathToSpec), filepath.Ext(args.PathToSpec))
	outFile := filepath.Join(flags.Out, baseName+meta.fileExt)
	if meta.format == gen.HTML_EMBED {
		outFile = filepath.Join(flags.Out, "ocli-docs.js")
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
func (a Actions) OcliCheck(_ context.Context, args gencli.OcliCheckArgs, flags gencli.OcliCheckFlags) error {
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

func (a Actions) OcliGenCli(_ context.Context, args gencli.OcliGenCliArgs, flags gencli.OcliGenCliFlags) error {
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

	fmt.Fprintf(stdout, "→ Generating CLI code:    framework=%s, output=%s\n", flags.Framework, flags.Out)

	// Ensure output directory exists
	if err := os.MkdirAll(flags.Out, 0755); err != nil {
		return fmt.Errorf("cannot create output directory %s: %w", flags.Out, err)
	}

	// Map the --framework flag to a gen.CLIFramework value.
	// This explicit mapping is the extension point for new frameworks.
	frameworkMap := map[string]gen.CLIFramework{
		"cobra": gen.CobraFramework,
		"yargs": gen.YargsFramework,
	}
	cliFramework, ok := frameworkMap[strings.ToLower(string(flags.Framework))]
	if !ok {
		supported := []string{"cobra", "yargs"}
		return gencli.NewValidationError(fmt.Sprintf("unsupported CLI framework: %q (supported: %s)", flags.Framework, strings.Join(supported, ", ")))

	opts := make([]gen.GenCLIOption, 0)
	opts = append(opts, gen.GenCLIWithFramework(cliFramework))

	// Generate CLI code
	output, err := gen.CLI(doc, opts...)
	if err != nil {
		return fmt.Errorf("failed to generate %s CLI code: %w", flags.Framework, err)
	}

	// Resolve output file path
	for fileName, fileContents := range output {
		outFile := filepath.Join(flags.Out, fileName)
		if err := os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
			return fmt.Errorf("cannot create output directory %s: %w", filepath.Dir(outFile), err)
		}
		if err := os.WriteFile(outFile, fileContents, 0644); err != nil {
			return fmt.Errorf("failed to write output file %s: %w", outFile, err)
		}
	}

	fmt.Fprintf(stdout, "✓ CLI Code written to: %s\n\n", flags.Out)

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

func (a Actions) Version() string {
	return a.version
}
