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

// cliAllCommandsTmplData is the template data passed to gencli/* templates.
type cliAllCommandsTmplData struct {
	ModuleVersion string
	Binary        string
	BinaryPascal  string
	LeafCommands  []cliCmdEntry
	ExitCodes     []spec.ExitCode
	GlobalFlags   []cobraFlagEntry
}

// cliCmdEntry holds the pre-computed template data for a single leaf command.
type cliCmdEntry struct {
	MethodName    string
	ArgsTypeName  string
	FlagsTypeName string
	Args          []cliFieldEntry
	Flags         []cliFieldEntry
}

// cliFieldEntry holds the pre-computed template data for a single arg/flag struct field.
type cliFieldEntry struct {
	FieldName string
	GoType    string
	TypeName  string // generated type name when Choices is non-empty
	Choices   []cliChoiceEntry
}

// cliChoiceEntry holds a single enumerated choice for a generated string type.
type cliChoiceEntry struct {
	ConstName string // e.g. "PetstorePetAddStatusAvailable"
	Value     string // e.g. "available"
}

// cobraCommandFileTmplData is the template data passed to command.tmpl.
// All generated cobra files belong to package gencli, so no cross-package imports are needed.
type cobraCommandFileTmplData struct {
	commandFileCoreTmplData
	PackageName  string // always "gencli"
	ChildImports []subCmdImport
	CobraArgs    []cobraArgEntry
	CobraFlags   []cobraFlagEntry
}

// subCmdImport holds data needed to call a child command constructor (same package, no import).
type subCmdImport struct {
	FuncName string // NewCmdPet, NewCmdPetAdd, ...
	Segment  string // original segment name, for getSpec*Cmd
	Summary  string // for getSpec*Cmd
}

// cobraArgEntry describes how to bind a positional argument in a cobra command.
type cobraArgEntry struct {
	FieldName  string
	Position   int
	IsRequired bool
	TypeName   string // non-empty when the struct field uses a generated type (needs cast)
}

// cobraFlagEntry describes how to bind a flag in a cobra command.
type cobraFlagEntry struct {
	FieldName    string
	VarName      string
	FlagName     string
	GoType       string
	CobraBindFn  string
	Default      string
	Summary      string
	TypeName     string   // non-empty when the struct field uses a generated type (needs cast)
	Shorthand    string   // first single-char alias, or empty string
	ExtraAliases []string // aliases not used as shorthand; mapped via SetNormalizeFunc
}

//go:embed templates/code/cobra
var cobraTemplateFiles embed.FS

func genCLICobra(doc *spec.Document, opts *genCLIOptions) (map[string][]byte, error) {
	out := make(map[string][]byte)

	binary := doc.Info.Binary
	binaryPascal := toPascalCase(binary)

	var leafCommands []cliCmdEntry
	var cmdFiles []cobraCommandFileTmplData

	rootCmd := doc.Commands
	walkCmdTree(doc, rootCmd, binary, binaryPascal, opts.ModuleVersion, []string{}, &leafCommands, &cmdFiles)

	if rootCmd.Summary == "" {
		cmdFiles[0].Summary = doc.Info.Summary
	}
	if rootCmd.Description == "" {
		cmdFiles[0].Description = doc.Info.Description
	}

	var exitCodes []spec.ExitCode
	var globalFlags []cobraFlagEntry
	if doc.Global != nil {
		exitCodes = doc.Global.ExitCodes
		for _, flag := range doc.Global.Flags {
			if flag.Name == "help" || flag.Name == "version" {
				continue
			}
			shorthand, extraAliases := splitAliases(flag.Aliases)
			globalFlags = append(globalFlags, cobraFlagEntry{
				FieldName:    toPascalCase(flag.Name),
				VarName:      "flag" + toPascalCase(flag.Name),
				FlagName:     flag.Name,
				GoType:       toGoType(flag.Type, flag.Variadic),
				CobraBindFn:  cobraBindFn(flag.Type, flag.Variadic),
				Default:      cobraDefaultVal(flag.Type, flag.Variadic),
				Summary:      flag.Summary,
				Shorthand:    shorthand,
				ExtraAliases: extraAliases,
			})
		}
	}

	allCmdsData := cliAllCommandsTmplData{
		ModuleVersion: opts.ModuleVersion,
		Binary:        binary,
		BinaryPascal:  binaryPascal,
		LeafCommands:  leafCommands,
		ExitCodes:     exitCodes,
		GlobalFlags:   globalFlags,
	}

	funcMap := cobraTemplateFuncMap()

	type gencliFile struct {
		outPath  string
		tmplPath string
	}
	gencliFiles := []gencliFile{
		{"gencli/actions.gen.go", "templates/code/cobra/gencli/actions.tmpl"},
		{"gencli/errors.gen.go", "templates/code/cobra/gencli/errors.tmpl"},
		{"gencli/help.gen.go", "templates/code/cobra/gencli/help.tmpl"},
		{"gencli/iostreams.gen.go", "templates/code/cobra/gencli/iostreams.tmpl"},
		{"gencli/params.gen.go", "templates/code/cobra/gencli/params.tmpl"},
	}

	for _, f := range gencliFiles {
		content, err := renderCobraTemplate(f.tmplPath, funcMap, allCmdsData)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", f.outPath, err)
		}
		formatted, err := format.Source(content)
		if err != nil {
			return nil, fmt.Errorf("formatting %s: %w\nsource:\n%s", f.outPath, err, content)
		}
		out[f.outPath] = formatted
	}

	runContent, err := renderCobraTemplate("templates/code/cobra/gencli/run.tmpl", funcMap, allCmdsData)
	if err != nil {
		return nil, fmt.Errorf("rendering gencli/run.go: %w", err)
	}
	formattedRun, err := format.Source(runContent)
	if err != nil {
		return nil, fmt.Errorf("formatting gencli/run.go: %w\nsource:\n%s", err, runContent)
	}
	out["gencli/run.go"] = formattedRun

	for _, cmdFile := range cmdFiles {
		content, err := renderCobraTemplate("templates/code/cobra/gencli/command.tmpl", funcMap, cmdFile)
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

// walkCmdTree recursively collects template data for all commands in the tree.
func walkCmdTree(
	doc *spec.Document,
	cmd *spec.CommandItem,
	binary, binaryPascal string,
	moduleVersion string,
	parentSegments []string,
	leafCommands *[]cliCmdEntry,
	cmdFiles *[]cobraCommandFileTmplData,
) {
	segments := appendSegment(parentSegments, cmd.Segment)

	isGroup := cmd.Group || len(cmd.Commands) > 0
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
				GoType:    "string", // positional args are always strings from cobra
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
	var cobraArgs []cobraArgEntry
	var cobraFlags []cobraFlagEntry

	for i, arg := range cmd.Args {
		argTypeName := ""
		if len(arg.Choices) > 0 {
			argTypeName = methodName + toPascalCase(arg.Name)
		}
		specArgs = append(specArgs, specArgEntry{Name: arg.Name, Summary: arg.Summary})
		cobraArgs = append(cobraArgs, cobraArgEntry{
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
		shorthand, extraAliases := splitAliases(flag.Aliases)
		specFlags = append(specFlags, specFlagEntry{Name: flag.Name, Summary: flag.Summary})
		cobraFlags = append(cobraFlags, cobraFlagEntry{
			FieldName:    toPascalCase(flag.Name),
			VarName:      "flag" + toPascalCase(flag.Name),
			FlagName:     flag.Name,
			GoType:       toGoType(flag.Type, flag.Variadic),
			CobraBindFn:  cobraBindFn(flag.Type, flag.Variadic),
			Default:      cobraDefaultVal(flag.Type, flag.Variadic),
			Summary:      flag.Summary,
			TypeName:     flagTypeName,
			Shorthand:    shorthand,
			ExtraAliases: extraAliases,
		})
	}

	childImports := buildChildImports(cmd.Commands, segments)
	cmdCore.FuncName = commandFuncName(segments)
	cmdCore.SpecFuncName = getSpecFuncName(segments)
	cmdCore.OutPath = commandOutPath(segments)
	cmdCore.CommandLine = strings.Join(append([]string{binary}, segments...), " ")
	cmdCore.SpecArgs = specArgs
	cmdCore.SpecFlags = specFlags

	cmdFile := cobraCommandFileTmplData{
		commandFileCoreTmplData: cmdCore,
		PackageName:             "gencli",
		ChildImports:            childImports,
		CobraArgs:               cobraArgs,
		CobraFlags:              cobraFlags,
	}
	*cmdFiles = append(*cmdFiles, cmdFile)

	for _, subcmd := range cmd.Commands {
		walkCmdTree(doc, subcmd, binary, binaryPascal, moduleVersion, segments, leafCommands, cmdFiles)
	}
}

// buildChildImports returns same-package call stubs for each direct child command.
func buildChildImports(cmds []*spec.CommandItem, parentSegments []string) []subCmdImport {
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

func renderCobraTemplate(tmplPath string, funcMap template.FuncMap, data any) ([]byte, error) {
	content, err := cobraTemplateFiles.ReadFile(tmplPath)
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

func cobraTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"goString": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
		"hasExtraAliases": func(flags []cobraFlagEntry) bool {
			for _, f := range flags {
				if len(f.ExtraAliases) > 0 {
					return true
				}
			}
			return false
		},
	}
}

// cobraBindFn returns the cobra Flags() method name for the given spec type and variadic flag.
// Always returns the 'P' variant so that a shorthand alias can be supplied as the third argument.
func cobraBindFn(t string, variadic bool) string {
	if variadic {
		switch t {
		case "integer":
			return "Int64SliceVarP"
		case "boolean":
			return "BoolSliceVarP"
		case "number":
			return "Float64SliceVarP"
		default:
			return "StringArrayVarP"
		}
	}
	switch t {
	case "integer":
		return "Int64VarP"
	case "boolean":
		return "BoolVarP"
	case "number":
		return "Float64VarP"
	default:
		return "StringVarP"
	}
}

// cobraDefaultVal returns the Go literal for the zero/empty default value of a cobra flag.
func cobraDefaultVal(t string, variadic bool) string {
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
		return "0"
	default:
		return `""`
	}
}

// commandFuncName returns the NewCmd* function name for the given segment path.
// [] -> "NewCmdRoot", ["pet"] -> "NewCmdPet", ["pet","add"] -> "NewCmdPetAdd"
func commandFuncName(segments []string) string {
	if len(segments) == 0 {
		return "NewCmdRoot"
	}
	parts := make([]string, len(segments))
	for i, s := range segments {
		parts[i] = toPascalCase(s)
	}
	return "NewCmd" + strings.Join(parts, "")
}

// getSpecFuncName returns the getSpec*Cmd function name for the given segment path.
// [] -> "getSpecRootCmd", ["pet","add"] -> "getSpecPetAddCmd"
func getSpecFuncName(segments []string) string {
	if len(segments) == 0 {
		return "getSpecRootCmd"
	}
	parts := make([]string, len(segments))
	for i, s := range segments {
		parts[i] = toPascalCase(s)
	}
	return "getSpec" + strings.Join(parts, "") + "Cmd"
}

// commandOutPath returns the output file path (under gencli/) for the given segment path.
// [] -> "gencli/root.go", ["pet"] -> "gencli/pet.go", ["pet","add"] -> "gencli/pet_add.go"
func commandOutPath(segments []string) string {
	if len(segments) == 0 {
		return "gencli/cmd_root.gen.go"
	}
	parts := make([]string, len(segments))
	for i, s := range segments {
		parts[i] = toGoPackageName(s)
	}
	return "gencli/cmd_" + strings.Join(parts, "_") + ".gen.go"
}
