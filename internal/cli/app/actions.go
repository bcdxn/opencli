package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/gen"
	"github.com/bcdxn/opencli/internal/cli/gencli"
	"github.com/bcdxn/opencli/spec"
	"github.com/bcdxn/opencli/validate"
)

// specCommands are the subcommands tried, in order, when the spec path is a CLI
// binary rather than a document. `__opencli` is the hidden command attached by
// the adapters; `docgen` is the public alternative (see WithPublicCommand).
var specCommands = []string{"__opencli", "docgen"}

// loadSpec reads an OpenCLI document from path and reports its format as a file
// extension (".json", ".yaml", or ".yml"). A .json/.yaml/.yml path is read as a file. Any
// other executable path is treated as a CLI binary: each of specCommands is run in
// turn and the first successful stdout is used.
func loadSpec(path string) ([]byte, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", gencli.NewValidationError(fmt.Sprintf("file not found: %s", path))
		}
		return nil, "", gencli.NewValidationError(fmt.Sprintf("cannot access file: %s (%v)", path, err))
	}
	if info.IsDir() {
		return nil, "", gencli.NewValidationError(fmt.Sprintf("path is a directory, not a file: %s", path))
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".yaml", ".yml":
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, "", gencli.NewValidationError(fmt.Sprintf("cannot read file: %s (%v)", path, err))
		}
		return data, ext, nil
	}

	if info.Mode()&0111 == 0 {
		return nil, "", gencli.NewValidationError(fmt.Sprintf("unsupported file format: %s (only .json, .yaml, .yml, or an executable CLI are supported)", ext))
	}

	var errs []string
	for _, sub := range specCommands {
		out, err := exec.Command(path, sub).Output()
		if err == nil {
			if strings.HasPrefix(strings.TrimSpace(string(out)), "{") {
				return out, ".json", nil
			}
			return out, ".yaml", nil
		}
		errs = append(errs, fmt.Sprintf("%s %s: %v", filepath.Base(path), sub, err))
	}
	return nil, "", gencli.NewValidationError(fmt.Sprintf("cannot get spec from %s (tried %s): %s", path, strings.Join(specCommands, ", "), strings.Join(errs, "; ")))
}

func decoderFor(ext string) func([]byte) (*spec.Document, error) {
	if ext == ".json" {
		return codec.UnmarshalJSON
	}
	return codec.UnmarshalYAML
}

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
	data, ext, err := loadSpec(args.PathToSpec)
	if err != nil {
		return err
	}
	decode := decoderFor(ext)

	// Map the --format flag to a DocFormat; this map is the extension point for new formats
	type docFormatMeta struct {
		format  gen.DocFormat
		fileExt string
	}
	supportedFormats := map[string]docFormatMeta{
		"markdown":   {gen.Markdown, ".md"},
		"html-page":  {gen.HTML_PAGE, ".html"},
		"html-embed": {gen.HTML_EMBED, ".js"},
		"man":        {gen.MAN, ".1"},
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

	// Parse the spec document
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
	data, ext, err := loadSpec(args.PathToSpec)
	if err != nil {
		return err
	}

	var validationErr error
	if ext == ".json" {
		validationErr = validate.ValidateJSON(data)
	} else {
		validationErr = validate.ValidateYAML(data)
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
	data, ext, err := loadSpec(args.PathToSpec)
	if err != nil {
		return err
	}
	decode := decoderFor(ext)

	stdout := a.IOS.Out()
	fmt.Fprintf(stdout, "\n→ Reading spec:       %s\n", args.PathToSpec)

	// Parse the spec document
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
		"cobra":     gen.CobraFramework,
		"yargs":     gen.YargsFramework,
		"urfavecli": gen.UrfaveCliFramework,
	}
	cliFramework, ok := frameworkMap[strings.ToLower(string(flags.Framework))]
	if !ok {
		supported := []string{"cobra", "yargs", "urfavecli"}
		return gencli.NewValidationError(fmt.Sprintf("unsupported CLI framework: %q (supported: %s)", flags.Framework, strings.Join(supported, ", ")))
	}

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
