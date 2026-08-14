package gen

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/bcdxn/opencli/spec"
)

// urfaveCliAllCommandsTmplData is the template data passed to urfave gencli/* templates.
type urfaveCliAllCommandsTmplData struct {
	ModuleVersion string
	Binary        string
	BinaryPascal  string
	LeafCommands  []cliCmdEntry
	ExitCodes     []spec.ExitCode
	GlobalFlags   []urfaveCliFlagEntry
	// Config file paths from global.config (only formats that are declared)
	ConfigJSON string
	ConfigTOML string
	ConfigYAML string
	// HasConfig is true if any config file path is declared
	HasConfig bool
}

// urfaveCliCommandFileTmplData is the template data passed to command.tmpl.
// All generated urfave files belong to package gencli, so no cross-package imports are needed.
type urfaveCliCommandFileTmplData struct {
	commandFileCoreTmplData
	PackageName  string // always "gencli"
	ChildImports []subCmdImport
	UrfaveArgs   []urfaveCliArgEntry
	UrfaveFlags  []urfaveCliFlagEntry
}

// urfaveCliArgEntry describes how to bind a positional argument in an urfave command.
type urfaveCliArgEntry struct {
	FieldName  string
	Position   int
	IsRequired bool
	TypeName   string // non-empty when the struct field uses a generated type (needs cast)
}

// urfaveCliAltSource represents an alternative source for a flag value ($ENV or $FILE).
type urfaveCliAltSource struct {
	Type     string // "$ENV" or "$FILE"
	Property string // env var name or JSONPath in config file
}

// urfaveCliFlagEntry describes how to bind a flag in an urfave command.
type urfaveCliFlagEntry struct {
	FieldName    string
	FlagName     string
	GoType       string
	UrfaveFlag   string // e.g. "cli.StringFlag", "cli.Int64Flag"
	Default      string // Go literal for the default value
	Summary      string
	TypeName     string               // non-empty when the struct field uses a generated type (needs cast)
	Aliases      []string             // all aliases (urfave uses Aliases []string, not separate shorthand)
	Accessor     string               // e.g. "String", "Int64", "Bool", "Float64", "StringSlice", etc.
	AltSources   []urfaveCliAltSource // alternative sources from flag's alternativeSources
	HasAltSource bool                 // true if AltSources is non-empty
}

//go:embed templates/code/urfavecli
var urfaveCliTemplateFiles embed.FS

func genCLIUrfaveCli(doc *spec.Document, opts *genCLIOptions) (map[string][]byte, error) {
	out := make(map[string][]byte)

	binary := doc.Info.Binary
	binaryPascal := toPascalCase(binary)

	var leafCommands []cliCmdEntry
	var cmdFiles []urfaveCliCommandFileTmplData

	rootCmd := doc.Commands
	walkUrfaveCliCmdTree(doc, rootCmd, binary, binaryPascal, opts.ModuleVersion, []string{}, &leafCommands, &cmdFiles)

	if rootCmd.Summary == "" {
		cmdFiles[0].Summary = doc.Info.Summary
	}
	if rootCmd.Description == "" {
		cmdFiles[0].Description = doc.Info.Description
	}

	var exitCodes []spec.ExitCode
	var globalFlags []urfaveCliFlagEntry
	var configJSON, configTOML, configYAML string
	hasConfig := false
	if doc.Global != nil {
		exitCodes = doc.Global.ExitCodes
		configJSON = doc.Global.Config.JSON
		configTOML = doc.Global.Config.TOML
		configYAML = doc.Global.Config.YAML
		hasConfig = configJSON != "" || configTOML != "" || configYAML != ""
		for _, flag := range doc.Global.Flags {
			if flag.Name == "help" || flag.Name == "version" {
				continue
			}
			altSources := urfaveCliAltSources(flag.AltSources)
			globalFlags = append(globalFlags, urfaveCliFlagEntry{
				FieldName:    toPascalCase(flag.Name),
				FlagName:     flag.Name,
				GoType:       toGoType(flag.Type, flag.Variadic),
				UrfaveFlag:   urfaveCliFlagStruct(flag.Type, flag.Variadic),
				Default:      urfaveCliDefaultVal(flag.Default, flag.Type, flag.Variadic),
				Summary:      flag.Summary,
				Aliases:      flag.Aliases,
				Accessor:     urfaveCliAccessor(flag.Type, flag.Variadic),
				AltSources:   altSources,
				HasAltSource: len(altSources) > 0,
			})
		}
	}

	allCmdsData := urfaveCliAllCommandsTmplData{
		ModuleVersion: opts.ModuleVersion,
		Binary:        binary,
		BinaryPascal:  binaryPascal,
		LeafCommands:  leafCommands,
		ExitCodes:     exitCodes,
		GlobalFlags:   globalFlags,
		ConfigJSON:    configJSON,
		ConfigTOML:    configTOML,
		ConfigYAML:    configYAML,
		HasConfig:     hasConfig,
	}

	funcMap := urfaveCliTemplateFuncMap()

	type gencliFile struct {
		outPath  string
		tmplPath string
	}
	gencliFiles := []gencliFile{
		{"gencli/actions.gen.go", "templates/code/urfavecli/gencli/actions.tmpl"},
		{"gencli/config.gen.go", "templates/code/urfavecli/gencli/config.tmpl"},
		{"gencli/errors.gen.go", "templates/code/urfavecli/gencli/errors.tmpl"},
		{"gencli/help.gen.go", "templates/code/urfavecli/gencli/help.tmpl"},
		{"gencli/iostreams.gen.go", "templates/code/urfavecli/gencli/iostreams.tmpl"},
		{"gencli/params.gen.go", "templates/code/urfavecli/gencli/params.tmpl"},
		{"gencli/run.gen.go", "templates/code/urfavecli/gencli/run.tmpl"},
	}

	for _, f := range gencliFiles {
		content, err := renderUrfaveCliTemplate(f.tmplPath, funcMap, allCmdsData)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", f.outPath, err)
		}
		formatted, err := format.Source(content)
		if err != nil {
			return nil, fmt.Errorf("formatting %s: %w\nsource:\n%s", f.outPath, err, content)
		}
		out[f.outPath] = formatted
	}

	for _, cmdFile := range cmdFiles {
		content, err := renderUrfaveCliTemplate("templates/code/urfavecli/gencli/command.tmpl", funcMap, cmdFile)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", cmdFile.OutPath, err)
		}
		formatted, err := format.Source(content)
		if err != nil {
			return nil, fmt.Errorf("formatting %s: %w\nsource:\n%s", cmdFile.OutPath, err, content)
		}
		out[cmdFile.OutPath] = formatted
	}

	return out, nil
}

// walkUrfaveCliCmdTree recursively collects template data for all commands in the tree.
func walkUrfaveCliCmdTree(
	doc *spec.Document,
	cmd *spec.CommandItem,
	binary, binaryPascal string,
	moduleVersion string,
	parentSegments []string,
	leafCommands *[]cliCmdEntry,
	cmdFiles *[]urfaveCliCommandFileTmplData,
) {
	segments := appendSegment(parentSegments, cmd.Segment)

	isGroup := cmd.Kind == spec.CommandKindGroup || len(cmd.Commands) > 0
	cmdCore := buildCommandFileCore(
		cmd.Segment,
		cmd.Summary,
		cmd.Description,
		cmd.Aliases,
		cmd.Hidden,
		cmd.VisibleChildren,
		cmd.VisibleArgs,
		cmd.VisibleFlags,
		cmd.CommandModifiers,
		cmd.ArgsModifiers,
		cmd.FlagsModifiers,
		cmd.PassthroughArgsModifiers,
		segments,
		binary,
		binaryPascal,
		moduleVersion,
		len(parentSegments) == 0,
		isGroup,
	)
	methodName := cmdCore.MethodName

	if !isGroup {
		entry := cliCmdEntry{
			MethodName:    methodName,
			ArgsTypeName:  methodName + "Args",
			FlagsTypeName: methodName + "Flags",
		}
		for _, arg := range cmd.Args {
			fe := cliFieldEntry{
				FieldName: toPascalCase(arg.Name),
				GoType:    "string", // positional args are always strings
			}
			if len(arg.Choices) > 0 {
				fe.TypeName = methodName + toPascalCase(arg.Name)
				fe.GoType = fe.TypeName
				for _, c := range arg.Choices {
					valStr := fmt.Sprintf("%v", c.Value)
					fe.Choices = append(fe.Choices, cliChoiceEntry{
						ConstName: fe.TypeName + toPascalCase(valStr),
						Value:     valStr,
					})
				}
			}
			entry.Args = append(entry.Args, fe)
		}

		for _, flag := range cmd.Flags {
			fe := cliFieldEntry{
				FieldName: toPascalCase(flag.Name),
				GoType:    toGoType(flag.Type, flag.Variadic),
			}
			if len(flag.Choices) > 0 && (flag.Type == "string" || flag.Type == "") && !flag.Variadic {
				fe.TypeName = methodName + toPascalCase(flag.Name)
				fe.GoType = fe.TypeName
				for _, c := range flag.Choices {
					valStr := fmt.Sprintf("%v", c.Value)
					fe.Choices = append(fe.Choices, cliChoiceEntry{
						ConstName: fe.TypeName + toPascalCase(valStr),
						Value:     valStr,
					})
				}
			}
			entry.Flags = append(entry.Flags, fe)
		}

		*leafCommands = append(*leafCommands, entry)
	}

	var specArgs []specArgEntry
	var specFlags []specFlagEntry
	var urfaveArgs []urfaveCliArgEntry
	var urfaveFlags []urfaveCliFlagEntry

	for i, arg := range cmd.Args {
		argTypeName := ""
		if len(arg.Choices) > 0 {
			argTypeName = methodName + toPascalCase(arg.Name)
		}
		specArgs = append(specArgs, specArgEntry{Name: arg.Name, Summary: arg.Summary})
		urfaveArgs = append(urfaveArgs, urfaveCliArgEntry{
			FieldName:  toPascalCase(arg.Name),
			Position:   i,
			IsRequired: arg.Required,
			TypeName:   argTypeName,
		})
	}

	for _, flag := range cmd.Flags {
		flagTypeName := ""
		if len(flag.Choices) > 0 && (flag.Type == "string" || flag.Type == "") && !flag.Variadic {
			flagTypeName = methodName + toPascalCase(flag.Name)
		}
		altSources := urfaveCliAltSources(flag.AltSources)
		specFlags = append(specFlags, specFlagEntry{
			Name:    flag.Name,
			Summary: flag.Summary,
			Aliases: flag.Aliases,
		})
		urfaveFlags = append(urfaveFlags, urfaveCliFlagEntry{
			FieldName:    toPascalCase(flag.Name),
			FlagName:     flag.Name,
			GoType:       toGoType(flag.Type, flag.Variadic),
			UrfaveFlag:   urfaveCliFlagStruct(flag.Type, flag.Variadic),
			Default:      urfaveCliDefaultVal(flag.Default, flag.Type, flag.Variadic),
			Summary:      flag.Summary,
			TypeName:     flagTypeName,
			Aliases:      flag.Aliases,
			Accessor:     urfaveCliAccessor(flag.Type, flag.Variadic),
			AltSources:   altSources,
			HasAltSource: len(altSources) > 0,
		})
	}

	childImports := urfaveBuildChildImports(cmd.Commands, segments)
	cmdCore.FuncName = commandFuncName(segments)
	cmdCore.SpecFuncName = getSpecFuncName(segments)
	cmdCore.OutPath = commandOutPath(segments)
	cmdCore.CommandLine = strings.Join(segments, " ")
	cmdCore.SpecArgs = specArgs
	cmdCore.SpecFlags = specFlags

	cmdFile := urfaveCliCommandFileTmplData{
		commandFileCoreTmplData: cmdCore,
		PackageName:             "gencli",
		ChildImports:            childImports,
		UrfaveArgs:              urfaveArgs,
		UrfaveFlags:             urfaveFlags,
	}
	*cmdFiles = append(*cmdFiles, cmdFile)

	for _, subcmd := range cmd.Commands {
		walkUrfaveCliCmdTree(doc, subcmd, binary, binaryPascal, moduleVersion, segments, leafCommands, cmdFiles)
	}
}

// urfaveBuildChildImports returns same-package call stubs for each direct child command.
func urfaveBuildChildImports(cmds []*spec.CommandItem, parentSegments []string) []subCmdImport {
	var imports []subCmdImport
	for _, cmd := range cmds {
		childSegments := make([]string, len(parentSegments)+1)
		copy(childSegments, parentSegments)
		childSegments[len(parentSegments)] = cmd.Segment
		imports = append(imports, subCmdImport{
			FuncName: commandFuncName(childSegments),
			Segment:  cmd.Segment,
			Summary:  cmd.Summary,
		})
	}
	return imports
}

func renderUrfaveCliTemplate(tmplPath string, funcMap template.FuncMap, data any) ([]byte, error) {
	content, err := urfaveCliTemplateFiles.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("reading template %s: %w", tmplPath, err)
	}

	t, err := template.New(filepath.Base(tmplPath)).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing template %s: %w", tmplPath, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("executing template %s: %w", tmplPath, err)
	}

	return buf.Bytes(), nil
}

func urfaveCliTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"goString": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
		"resolveFlagValue": func(f urfaveCliFlagEntry) string {
			if !f.HasAltSource {
				// No alternative sources — read directly from CLI
				if f.TypeName != "" {
					return fmt.Sprintf("%s(c.%s(%q))", f.TypeName, f.Accessor, f.FlagName)
				}
				return fmt.Sprintf("c.%s(%q)", f.Accessor, f.FlagName)
			}
			// Has alternative sources — read from CLI first, fall back to env/config
			// We use a helper that checks if the CLI value is the zero value
			switch f.GoType {
			case "string":
				srcs := formatAltSources(f.AltSources)
				if f.TypeName != "" {
					return fmt.Sprintf("%s(resolveStringFlag(c.String(%q), %s))", f.TypeName, f.FlagName, srcs)
				}
				return fmt.Sprintf("resolveStringFlag(c.String(%q), %s)", f.FlagName, srcs)
			case "int64":
				srcs := formatAltSources(f.AltSources)
				return fmt.Sprintf("resolveInt64Flag(c.Int64(%q), %s)", f.FlagName, srcs)
			case "bool":
				srcs := formatAltSources(f.AltSources)
				return fmt.Sprintf("resolveBoolFlag(c.Bool(%q), %s)", f.FlagName, srcs)
			case "float64":
				srcs := formatAltSources(f.AltSources)
				return fmt.Sprintf("resolveFloat64Flag(c.Float64(%q), %s)", f.FlagName, srcs)
			case "[]string":
				srcs := formatAltSources(f.AltSources)
				return fmt.Sprintf("resolveStringSliceFlag(c.StringSlice(%q), %s)", f.FlagName, srcs)
			case "[]int64":
				srcs := formatAltSources(f.AltSources)
				return fmt.Sprintf("resolveInt64SliceFlag(c.Int64Slice(%q), %s)", f.FlagName, srcs)
			case "[]bool":
				srcs := formatAltSources(f.AltSources)
				return fmt.Sprintf("resolveBoolSliceFlag(c.BoolSlice(%q), %s)", f.FlagName, srcs)
			case "[]float64":
				srcs := formatAltSources(f.AltSources)
				return fmt.Sprintf("resolveFloat64SliceFlag(c.Float64Slice(%q), %s)", f.FlagName, srcs)
			default:
				if f.TypeName != "" {
					return fmt.Sprintf("%s(c.%s(%q))", f.TypeName, f.Accessor, f.FlagName)
				}
				return fmt.Sprintf("c.%s(%q)", f.Accessor, f.FlagName)
			}
		},
	}
}

// urfaveCliFlagStruct returns the urfave/cli v3 flag struct type for the given spec type.
func urfaveCliFlagStruct(t string, variadic bool) string {
	if variadic {
		switch t {
		case "integer":
			return "cli.Int64SliceFlag"
		case "boolean":
			return "cli.BoolSliceFlag"
		case "number":
			return "cli.Float64SliceFlag"
		default:
			return "cli.StringSliceFlag"
		}
	}
	switch t {
	case "integer":
		return "cli.Int64Flag"
	case "boolean":
		return "cli.BoolFlag"
	case "number":
		return "cli.Float64Flag"
	default:
		return "cli.StringFlag"
	}
}

// urfaveCliAccessor returns the method name to read a flag value from *cli.Command.
func urfaveCliAccessor(t string, variadic bool) string {
	if variadic {
		switch t {
		case "integer":
			return "Int64Slice"
		case "boolean":
			return "BoolSlice"
		case "number":
			return "Float64Slice"
		default:
			return "StringSlice"
		}
	}
	switch t {
	case "integer":
		return "Int64"
	case "boolean":
		return "Bool"
	case "number":
		return "Float64"
	default:
		return "String"
	}
}

// urfaveCliZeroValue returns the Go zero-value expression for a flag type.
func urfaveCliZeroValue(t string, variadic bool) string {
	if variadic {
		switch t {
		case "integer":
			return "[]int64{}"
		case "boolean":
			return "[]bool{}"
		case "number":
			return "[]float64{}"
		default:
			return "[]string{}"
		}
	}
	switch t {
	case "integer":
		return "0"
	case "boolean":
		return "false"
	case "number":
		return "0.0"
	default:
		return `""`
	}
}

// urfaveCliDefaultVal returns the Go literal for the default value of an urfave flag.
func urfaveCliDefaultVal(val any, t string, variadic bool) string {
	switch slice := val.(type) {
	// handle slice types first
	case []string:
		var elems []string
		for _, v := range slice {
			elems = append(elems, fmt.Sprintf("%#v", v))
		}
		return fmt.Sprintf("[]string{%s}", strings.Join(elems, ", "))

	case []int:
		var elems []string
		for _, v := range slice {
			elems = append(elems, fmt.Sprintf("%d", v))
		}
		return fmt.Sprintf("[]int64{%s}", strings.Join(elems, ", "))

	case []float64:
		var elems []string
		for _, v := range slice {
			// %g prints the most compact representation of a float
			elems = append(elems, fmt.Sprintf("%g", v))
		}
		return fmt.Sprintf("[]float64{%s}", strings.Join(elems, ", "))

	case []bool:
		var elems []string
		for _, v := range slice {
			elems = append(elems, fmt.Sprintf("%t", v))
		}
		return fmt.Sprintf("[]bool{%s}", strings.Join(elems, ", "))
	// handle non-slice scalars
	case string:
		return fmt.Sprintf("%q", val)
	case int:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case nil:
		return urfaveCliZeroValue(t, variadic)

	default:
		// should never panic because the spec will have been validated before generation is run
		panic(fmt.Sprintf("unsupported type: must be a slice of string, int, float64, or bool - %T", val))
	}
}

// urfaveCliAltSources converts spec.AlternativeSource slices to urfaveCliAltSource slices.
func urfaveCliAltSources(sources []spec.AlternativeSource) []urfaveCliAltSource {
	if len(sources) == 0 {
		return nil
	}
	result := make([]urfaveCliAltSource, 0, len(sources))
	for _, s := range sources {
		result = append(result, urfaveCliAltSource{
			Type:     s.Type,
			Property: s.Property,
		})
	}
	return result
}

// formatAltSources formats a slice of urfaveCliAltSource as a Go composite literal
// for use in generated template code.
func formatAltSources(sources []urfaveCliAltSource) string {
	if len(sources) == 0 {
		return "nil"
	}
	var parts []string
	for _, s := range sources {
		parts = append(parts, fmt.Sprintf("{Type: %q, Property: %q}", s.Type, s.Property))
	}
	return fmt.Sprintf("[]struct{Type,Property string}{%s}", strings.Join(parts, ","))
}
