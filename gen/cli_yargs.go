package gen

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/bcdxn/opencli/spec"
)

// yargsAllCmdsTmplData is the top-level data passed to support-file templates (actions, run, etc.).
type yargsAllCmdsTmplData struct {
	ModuleVersion string
	Binary        string
	BinaryPascal  string
	LeafCommands  []yargsCmdEntry
	RootImport    yargsChildImport   // import stub for the root command (used by run.ts)
	ChildImports  []yargsChildImport // direct children of root (used by run.ts)
	ExitCodes     []spec.ExitCode
}

// yargsCmdEntry holds pre-computed data for a single leaf command (used by actions/params).
type yargsCmdEntry struct {
	MethodName    string
	ArgsTypeName  string
	FlagsTypeName string
	Args          []yargsFieldEntry
	Flags         []yargsFieldEntry
}

// yargsFieldEntry describes one arg or flag field in a TypeScript interface.
type yargsFieldEntry struct {
	FieldName  string
	TSType     string
	TypeName   string // non-empty when a dedicated enum type is emitted
	Choices    []yargsChoiceEntry
	IsVariadic bool
	IsRequired bool
}

// yargsChoiceEntry holds one allowed value for an enum field.
type yargsChoiceEntry struct {
	EnumKey string // e.g. AVAILABLE
	Value   string // e.g. "available"
}

// yargsCommandFileTmplData is the template data for a single generated command module file.
type yargsCommandFileTmplData struct {
	ModuleVersion    string
	Binary           string
	BinaryPascal     string
	FuncName         string // e.g. newPetstorePetAddCmd
	SpecFuncName     string // e.g. getPetstorePetAddCmdHelpData
	OutPath          string // relative output path e.g. gencli/cmd-petstore-pet-add.ts
	Segment          string
	SegmentDSL       string
	IsRoot           bool
	IsGroup          bool
	IsHidden         bool
	MethodName       string
	ArgsTypeName     string
	FlagsTypeName    string
	CommandArgName   string // camelCase local variable for yargs argv type annotation
	Summary          string
	Description      string
	Aliases          []string
	CommandLine      string
	VisibleChildren  bool
	VisibleArgs      bool
	VisibleFlags     bool
	ChildImports     []yargsChildImport
	YargsArgs        []yargsArgEntry
	YargsFlags       []yargsFlagEntry
	SpecArgs         []specArgEntry
	SpecFlags        []specFlagEntry
	CommandModifiers []string
	ArgsModifiers    []string
	FlagsModifiers   []string
}

// yargsChildImport holds data for importing and registering a child command module.
type yargsChildImport struct {
	FuncName string // e.g. newPetstorePetCmd
	FileName string // e.g. cmd-petstore-pet  (no extension, for import path)
	Segment  string
	Summary  string
}

// yargsArgEntry describes how to bind a positional argument in a yargs command.
type yargsArgEntry struct {
	FieldName  string // camelCase field name used in argv (yargs camelises kebab names)
	RawName    string // unmodified name from spec, used for .positional() and the local argv interface
	TSType     string
	IsRequired bool
	TypeName   string // non-empty when the field uses a generated enum type
	Choices    []yargsChoiceEntry
}

// yargsFlagEntry describes how to bind an option/flag in a yargs command.
type yargsFlagEntry struct {
	FieldName    string // camelCase field on argv
	RawName      string // unmodified name from spec, used for .option() and the local argv interface
	TSType       string
	IsRequired   bool
	IsVariadic   bool
	TypeName     string // non-empty when field uses generated enum type
	Choices      []yargsChoiceEntry
	Shorthand    string
	ExtraAliases []string
	Default      string // TypeScript literal or empty
}

//go:embed templates/code/yargs
var yargsTemplateFiles embed.FS

func genCLIYargs(doc *spec.Document, opts *genCLIOptions) (map[string][]byte, error) {
	out := make(map[string][]byte)

	binary := doc.Info.Binary
	binaryPascal := toPascalCase(binary)

	var leafCommands []yargsCmdEntry
	var cmdFiles []yargsCommandFileTmplData

	rootCmd := doc.Commands
	walkYargsCmdTree(doc, rootCmd, binary, binaryPascal, opts.ModuleVersion, []string{}, &leafCommands, &cmdFiles)

	if rootCmd.Summary == "" && len(cmdFiles) > 0 {
		cmdFiles[0].Summary = doc.Info.Summary
	}
	if rootCmd.Description == "" && len(cmdFiles) > 0 {
		cmdFiles[0].Description = doc.Info.Description
	}

	rootChildImports := yargsBuildChildImports(rootCmd.Commands, []string{binary})

	allCmdsData := yargsAllCmdsTmplData{
		ModuleVersion: opts.ModuleVersion,
		Binary:        binary,
		BinaryPascal:  binaryPascal,
		LeafCommands:  leafCommands,
		ChildImports:  rootChildImports,
		RootImport: yargsChildImport{
			FuncName: yargsCommandFuncName([]string{binary}),
			FileName: yargsCommandFileName([]string{binary}),
			Segment:  binary,
			Summary:  doc.Info.Summary,
		},
	}
	if doc.Global != nil {
		allCmdsData.ExitCodes = doc.Global.ExitCodes
	}

	funcMap := yargsTemplateFuncMap()

	type gencliFile struct {
		outPath  string
		tmplPath string
	}
	supportFiles := []gencliFile{
		{"gencli/actions.ts", "templates/code/yargs/gencli/actions.tmpl"},
		{"gencli/params.ts", "templates/code/yargs/gencli/params.tmpl"},
		{"gencli/errors.ts", "templates/code/yargs/gencli/errors.tmpl"},
		{"gencli/help.ts", "templates/code/yargs/gencli/help.tmpl"},
		{"gencli/types.ts", "templates/code/yargs/gencli/types.tmpl"},
		{"gencli/run.ts", "templates/code/yargs/gencli/run.tmpl"},
	}
	for _, f := range supportFiles {
		content, err := renderYargsTemplate(f.tmplPath, funcMap, allCmdsData)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", f.outPath, err)
		}
		out[f.outPath] = content
	}

	sort.Slice(cmdFiles, func(i, j int) bool { return cmdFiles[i].OutPath < cmdFiles[j].OutPath })
	for _, cmdFile := range cmdFiles {
		content, err := renderYargsTemplate("templates/code/yargs/gencli/command.tmpl", funcMap, cmdFile)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", cmdFile.OutPath, err)
		}
		out[cmdFile.OutPath] = content
	}

	return out, nil
}

// walkYargsCmdTree recursively collects template data for all commands in the tree.
func walkYargsCmdTree(
	doc *spec.Document,
	cmd *spec.CommandItem,
	binary, binaryPascal string,
	moduleVersion string,
	parentSegments []string,
	leafCommands *[]yargsCmdEntry,
	cmdFiles *[]yargsCommandFileTmplData,
) {
	segments := make([]string, len(parentSegments)+1)
	copy(segments, parentSegments)
	segments[len(parentSegments)] = cmd.Segment

	isGroup := cmd.Group || len(cmd.Commands) > 0
	methodName := buildMethodName(segments)

	if !isGroup {
		entry := yargsCmdEntry{
			MethodName:    methodName,
			ArgsTypeName:  methodName + "Args",
			FlagsTypeName: methodName + "Flags",
		}
		for _, arg := range cmd.Args {
			fe := yargsFieldEntry{
				FieldName:  toCamelCase(arg.Name),
				TSType:     toTSType(arg.Type, false),
				IsRequired: arg.Required,
			}
			if len(arg.Choices) > 0 {
				fe.TypeName = methodName + toPascalCase(arg.Name)
				fe.TSType = fe.TypeName
				for _, c := range arg.Choices {
					valStr := fmt.Sprintf("%v", c.Value)
					fe.Choices = append(fe.Choices, yargsChoiceEntry{
						EnumKey: strings.ToUpper(strings.ReplaceAll(toGoPackageName(valStr), "-", "_")),
						Value:   valStr,
					})
				}
			}
			entry.Args = append(entry.Args, fe)
		}
		for _, flag := range cmd.Flags {
			fe := yargsFieldEntry{
				FieldName:  toCamelCase(flag.Name),
				TSType:     toTSType(flag.Type, flag.Variadic),
				IsRequired: flag.Required,
				IsVariadic: flag.Variadic,
			}
			if len(flag.Choices) > 0 && (flag.Type == "string" || flag.Type == "") && !flag.Variadic {
				fe.TypeName = methodName + toPascalCase(flag.Name)
				fe.TSType = fe.TypeName
				for _, c := range flag.Choices {
					valStr := fmt.Sprintf("%v", c.Value)
					fe.Choices = append(fe.Choices, yargsChoiceEntry{
						EnumKey: strings.ToUpper(strings.ReplaceAll(toGoPackageName(valStr), "-", "_")),
						Value:   valStr,
					})
				}
			}
			entry.Flags = append(entry.Flags, fe)
		}
		*leafCommands = append(*leafCommands, entry)
	}

	var yargsArgs []yargsArgEntry
	var yargsFlags []yargsFlagEntry
	var specArgs []specArgEntry
	var specFlags []specFlagEntry

	for _, arg := range cmd.Args {
		argTypeName := ""
		var choices []yargsChoiceEntry
		if len(arg.Choices) > 0 {
			argTypeName = methodName + toPascalCase(arg.Name)
			for _, c := range arg.Choices {
				valStr := fmt.Sprintf("%v", c.Value)
				choices = append(choices, yargsChoiceEntry{
					EnumKey: strings.ToUpper(strings.ReplaceAll(toGoPackageName(valStr), "-", "_")),
					Value:   valStr,
				})
			}
		}
		specArgs = append(specArgs, specArgEntry{Name: arg.Name, Summary: arg.Summary})
		yargsArgs = append(yargsArgs, yargsArgEntry{
			FieldName:  toCamelCase(arg.Name),
			RawName:    arg.Name,
			TSType:     toTSType(arg.Type, false),
			IsRequired: arg.Required,
			TypeName:   argTypeName,
			Choices:    choices,
		})
	}

	for _, flag := range cmd.Flags {
		flagTypeName := ""
		var choices []yargsChoiceEntry
		if len(flag.Choices) > 0 && (flag.Type == "string" || flag.Type == "") && !flag.Variadic {
			flagTypeName = methodName + toPascalCase(flag.Name)
			for _, c := range flag.Choices {
				valStr := fmt.Sprintf("%v", c.Value)
				choices = append(choices, yargsChoiceEntry{
					EnumKey: strings.ToUpper(strings.ReplaceAll(toGoPackageName(valStr), "-", "_")),
					Value:   valStr,
				})
			}
		}
		shorthand := ""
		extraAliases := []string{}
		for _, a := range flag.Aliases {
			if len(a) == 1 && shorthand == "" {
				shorthand = a
			}
			extraAliases = append(extraAliases, a)
		}
		specFlags = append(specFlags, specFlagEntry{Name: flag.Name, Summary: flag.Summary, Aliases: extraAliases})
		yargsFlags = append(yargsFlags, yargsFlagEntry{
			FieldName:    toCamelCase(flag.Name),
			RawName:      flag.Name,
			TSType:       toTSType(flag.Type, flag.Variadic),
			IsRequired:   flag.Required,
			IsVariadic:   flag.Variadic,
			TypeName:     flagTypeName,
			Choices:      choices,
			Shorthand:    shorthand,
			ExtraAliases: extraAliases,
			Default:      yargsDefaultVal(flag.Type, flag.Variadic),
		})
	}

	childImports := yargsBuildChildImports(cmd.Commands, segments)

	cmdFile := yargsCommandFileTmplData{
		ModuleVersion:    moduleVersion,
		Binary:           binary,
		BinaryPascal:     binaryPascal,
		FuncName:         yargsCommandFuncName(segments),
		SpecFuncName:     yargsSpecFuncName(segments),
		OutPath:          yargsCommandOutPath(segments),
		Segment:          cmd.Segment,
		SegmentDSL:       yargsCommandDSL(cmd),
		IsRoot:           len(parentSegments) == 0,
		IsGroup:          isGroup,
		IsHidden:         cmd.Hidden,
		MethodName:       methodName,
		ArgsTypeName:     methodName + "Args",
		FlagsTypeName:    methodName + "Flags",
		CommandArgName:   toCamelCase(strings.Join(segments, "-")) + "Args",
		Summary:          cmd.Summary,
		Description:      strings.TrimRight(cmd.Description, "\n"),
		Aliases:          cmd.Aliases,
		CommandLine:      strings.Join(segments, " "),
		VisibleChildren:  cmd.VisibleChildren,
		VisibleArgs:      cmd.VisibleArgs,
		VisibleFlags:     cmd.VisibleFlags,
		ChildImports:     childImports,
		YargsArgs:        yargsArgs,
		YargsFlags:       yargsFlags,
		SpecArgs:         specArgs,
		SpecFlags:        specFlags,
		CommandModifiers: cmd.CommandModifiers,
		ArgsModifiers:    cmd.ArgsModifiers,
		FlagsModifiers:   cmd.FlagsModifiers,
	}
	*cmdFiles = append(*cmdFiles, cmdFile)

	for _, subcmd := range cmd.Commands {
		walkYargsCmdTree(doc, subcmd, binary, binaryPascal, moduleVersion, segments, leafCommands, cmdFiles)
	}
}

func yargsBuildChildImports(cmds []*spec.CommandItem, parentSegments []string) []yargsChildImport {
	var imports []yargsChildImport
	for _, cmd := range cmds {
		childSegs := make([]string, len(parentSegments)+1)
		copy(childSegs, parentSegments)
		childSegs[len(parentSegments)] = cmd.Segment
		imports = append(imports, yargsChildImport{
			FuncName: yargsCommandFuncName(childSegs),
			FileName: yargsCommandFileName(childSegs),
			Segment:  cmd.Segment,
			Summary:  cmd.Summary,
		})
	}
	return imports
}

func renderYargsTemplate(tmplPath string, funcMap template.FuncMap, data any) ([]byte, error) {
	content, err := yargsTemplateFiles.ReadFile(tmplPath)
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

func yargsTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"tsString": func(s string) string {
			return fmt.Sprintf("%q", s)
		},
		// collectParamImports returns the unique param type names (ArgsTypeName /
		// FlagsTypeName) that are actually needed by at least one leaf command.
		"collectParamImports": func(cmds []yargsCmdEntry) []string {
			seen := make(map[string]bool)
			var result []string
			for _, cmd := range cmds {
				if len(cmd.Args) > 0 && !seen[cmd.ArgsTypeName] {
					seen[cmd.ArgsTypeName] = true
					result = append(result, cmd.ArgsTypeName)
				}
				if len(cmd.Flags) > 0 && !seen[cmd.FlagsTypeName] {
					seen[cmd.FlagsTypeName] = true
					result = append(result, cmd.FlagsTypeName)
				}
			}
			return result
		},
		"joinStrings": strings.Join,
	}
}

// yargsSpecFuncName returns the help-data factory function name for a command.
// ["petstore","pet","add"] -> "getPetstorePetAddCmdHelpData"
func yargsSpecFuncName(segments []string) string {
	parts := make([]string, len(segments))
	for i, s := range segments {
		parts[i] = toPascalCase(s)
	}
	return "get" + strings.Join(parts, "") + "CmdHelpData"
}

// yargsCommandFuncName returns the factory function name for a command (camelCase).
// ["petstore","pet","add"] -> "newPetstorePetAddCmd"
func yargsCommandFuncName(segments []string) string {
	parts := make([]string, len(segments))
	for i, s := range segments {
		parts[i] = toPascalCase(s)
	}
	return "new" + strings.Join(parts, "") + "Cmd"
}

// yargsCommandFileName returns the stem of the output file (no extension).
// ["petstore","pet","add"] -> "cmd-petstore-pet-add"
func yargsCommandFileName(segments []string) string {
	lower := make([]string, len(segments))
	for i, s := range segments {
		lower[i] = strings.ToLower(s)
	}
	return "cmd-" + strings.Join(lower, "-")
}

// yargsCommandOutPath returns the relative output path for a command module.
// ["petstore","pet","add"] -> "gencli/cmd-petstore-pet-add.ts"
func yargsCommandOutPath(segments []string) string {
	return "gencli/" + yargsCommandFileName(segments) + ".ts"
}

// yargsCommandDSL builds the command property for yargs data:
// command aliases + required positional args as DSL tokens.
func yargsCommandDSL(cmd *spec.CommandItem) string {
	aliases := []string{cmd.Segment}
	for _, alias := range cmd.Aliases {
		aliases = append(aliases, alias)
	}

	args := []string{}
	for _, arg := range cmd.Args {
		args = append(args, "<"+arg.Name+">")
	}

	cmdDSL := []string{strings.Join(aliases, "|")}
	cmdDSL = append(cmdDSL, args...)

	return strings.Join(cmdDSL, " ")
}

// yargsDefaultVal returns a TypeScript literal default value for a flag.
func yargsDefaultVal(t string, variadic bool) string {
	if variadic {
		return ""
	}
	switch t {
	case "integer", "number":
		return "0"
	case "boolean":
		return "false"
	default:
		return ""
	}
}
