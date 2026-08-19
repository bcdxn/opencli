package gen

import (
	"bytes"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
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
	GlobalFlags   []yargsFlagEntry
	// Config file paths from global.config (only formats that are declared)
	ConfigJSON string
	ConfigTOML string
	ConfigYAML string
	// HasAltSources is true if any flag declares an alternative source, which
	// requires emitting gencli/config.ts and calling loadConfig at startup.
	HasAltSources bool
	// HasFileAltSource is true if any flag declares a $FILE alternative source,
	// which requires the generated config code to import a JSONPath library.
	HasFileAltSource bool
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
	Default    any
}

// yargsChoiceEntry holds one allowed value for an enum field.
type yargsChoiceEntry struct {
	EnumKey string // e.g. AVAILABLE
	Value   string // e.g. "available"
}

// yargsCommandFileTmplData is the template data for a single generated command module file.
type yargsCommandFileTmplData struct {
	commandFileCoreTmplData
	SegmentDSL     string
	CommandArgName string // camelCase local variable for yargs argv type annotation
	ChildImports   []yargsChildImport
	YargsArgs      []yargsArgEntry
	YargsFlags     []yargsFlagEntry
	GlobalFlags    []yargsFlagEntry // non-help/version global flags, shared by all leaf commands
	ConfigImports  []string         // unique resolver names from gencli/config.ts used by this command's alt-source flags
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
	Default    string //Typescript literal or empty
}

// yargsFlagEntry describes how to bind an option/flag in a yargs command.
type yargsFlagEntry struct {
	FieldName      string // camelCase field on argv
	RawName        string // unmodified name from spec, used for .option() and the local argv interface
	TSType         string
	IsRequired     bool
	IsVariadic     bool
	TypeName       string // non-empty when field uses generated enum type
	Choices        []yargsChoiceEntry
	Shorthand      string
	ExtraAliases   []string
	Default        string // TypeScript literal or empty
	Summary        string // flag summary (used for global option registration)
	VariadicCoerce string // pre-rendered JS arrow coercing array elements for non-string variadics ("" otherwise)
	AltSources     []spec.AlternativeSource
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

	var globalFlags []yargsFlagEntry
	configJSON, configTOML, configYAML := "", "", ""
	if doc.Global != nil {
		configJSON = doc.Global.Config.JSON
		configTOML = doc.Global.Config.TOML
		configYAML = doc.Global.Config.YAML
		for _, flag := range doc.Global.Flags {
			if flag.Name == "help" || flag.Name == "version" {
				continue
			}
			shorthand, extraAliases := splitAliases(flag.Aliases)
			globalFlags = append(globalFlags, yargsFlagEntry{
				FieldName:    toCamelCase(flag.Name),
				RawName:      flag.Name,
				TSType:       toTSType(flag.Type, flag.Variadic),
				IsRequired:   flag.Required,
				IsVariadic:   flag.Variadic,
				Shorthand:    shorthand,
				ExtraAliases: extraAliases,
				Default:      yargsDefaultVal(flag.Default),
				Summary:      flag.Summary,
				AltSources:   flag.AltSources,
			})
		}
	}

	// Track alternative-source usage across all flags so we know whether to emit
	// gencli/config.ts (and call loadConfig from run.tmpl). A $FILE source
	// additionally requires the generated config code to import a JSONPath library.
	hasAltSources, hasFileAltSource := scanYargsAltSources(globalFlags)
	for i := range cmdFiles {
		cmdHasAlt, cmdHasFile := scanYargsAltSources(cmdFiles[i].YargsFlags)
		hasAltSources = hasAltSources || cmdHasAlt
		hasFileAltSource = hasFileAltSource || cmdHasFile
	}

	allCmdsData := yargsAllCmdsTmplData{
		ModuleVersion:    opts.ModuleVersion,
		Binary:           binary,
		BinaryPascal:     binaryPascal,
		LeafCommands:     leafCommands,
		ChildImports:     rootChildImports,
		GlobalFlags:      globalFlags,
		ConfigJSON:       configJSON,
		ConfigTOML:       configTOML,
		ConfigYAML:       configYAML,
		HasAltSources:    hasAltSources,
		HasFileAltSource: hasFileAltSource,
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
		// Only emitted when at least one flag declares alternative sources, so that
		// specs without this feature produce no extra files or imports.
	}
	if hasAltSources {
		supportFiles = append(supportFiles, gencliFile{"gencli/config.ts", "templates/code/yargs/gencli/config.tmpl"})
	}
	supportFiles = append(supportFiles,
		gencliFile{"gencli/params.ts", "templates/code/yargs/gencli/params.tmpl"},
		gencliFile{"gencli/errors.ts", "templates/code/yargs/gencli/errors.tmpl"},
		gencliFile{"gencli/help.ts", "templates/code/yargs/gencli/help.tmpl"},
		gencliFile{"gencli/types.ts", "templates/code/yargs/gencli/types.tmpl"},
		gencliFile{"gencli/run.ts", "templates/code/yargs/gencli/run.tmpl"},
	)
	for _, f := range supportFiles {
		content, err := renderYargsTemplate(f.tmplPath, funcMap, allCmdsData)
		if err != nil {
			return nil, fmt.Errorf("rendering %s: %w", f.outPath, err)
		}
		out[f.outPath] = content
	}

	sort.Slice(cmdFiles, func(i, j int) bool { return cmdFiles[i].OutPath < cmdFiles[j].OutPath })
	for i := range cmdFiles {
		cmdFiles[i].GlobalFlags = globalFlags
	}
	// Collect the resolver functions each command file needs from gencli/config.ts,
	// based on its own flags plus any alt-source global flags (which every leaf
	// command handler also resolves).
	hasGlobalAlt, _ := scanYargsAltSources(globalFlags)
	for i := range cmdFiles {
		seen := make(map[string]bool)
		var imports []string
		addResolver := func(f yargsFlagEntry) {
			if len(f.AltSources) == 0 {
				return
			}
			if r, ok := yargsAltSourceResolvers[f.TSType]; ok && !seen[r] {
				seen[r] = true
				imports = append(imports, r)
			}
		}
		for _, f := range cmdFiles[i].YargsFlags {
			addResolver(f)
		}
		if hasGlobalAlt {
			for _, f := range globalFlags {
				addResolver(f)
			}
		}
		cmdFiles[i].ConfigImports = imports
	}
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
				Default:    arg.Default,
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
				Default:    flag.Default,
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
			FieldName:      toCamelCase(flag.Name),
			RawName:        flag.Name,
			TSType:         toTSType(flag.Type, flag.Variadic),
			IsRequired:     flag.Required,
			IsVariadic:     flag.Variadic,
			VariadicCoerce: yargsVariadicCoerce(flag.Type),
			TypeName:       flagTypeName,
			Choices:        choices,
			Shorthand:      shorthand,
			ExtraAliases:   extraAliases,
			Default:        yargsDefaultVal(flag.Default),
			AltSources:     flag.AltSources,
		})
	}

	childImports := yargsBuildChildImports(cmd.Commands, segments)
	cmdCore.FuncName = yargsCommandFuncName(segments)
	cmdCore.SpecFuncName = yargsSpecFuncName(segments)
	cmdCore.OutPath = yargsCommandOutPath(segments)
	cmdCore.CommandLine = strings.Join(segments, " ")
	cmdCore.SpecArgs = specArgs
	cmdCore.SpecFlags = specFlags

	cmdFile := yargsCommandFileTmplData{
		commandFileCoreTmplData: cmdCore,
		SegmentDSL:              yargsCommandDSL(cmd),
		CommandArgName:          toCamelCase(strings.Join(segments, "-")) + "Args",
		ChildImports:            childImports,
		YargsArgs:               yargsArgs,
		YargsFlags:              yargsFlags,
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
		// resolveFlagValue returns the expression that yields a flag's value in
		// generated command code. Without alternative sources it reads the parsed
		// argv field directly; with alt-sources it calls the type-specific resolver
		// from gencli/config.ts, which prefers an explicit CLI value and otherwise
		// falls back to env/config in declared order.
		"resolveFlagValue": func(f yargsFlagEntry) string {
			if len(f.AltSources) == 0 {
				return plainYargsFlagExpr(f)
			}
			resolver, ok := yargsAltSourceResolvers[f.TSType]
			if !ok {
				return plainYargsFlagExpr(f)
			}
			names := formatTSAltNames(yargsAltSourceNames(f))
			expr := fmt.Sprintf("%s(argv, %s, %s)", resolver, names, formatTSAltSources(f.AltSources))
			if f.TypeName != "" {
				expr = fmt.Sprintf("%s(%s)", f.TypeName, expr)
			}
			return expr
		},
	}
}

// plainYargsFlagExpr returns the expression for a flag with no alternative sources:
// the parsed argv field cast to its generated type (choices enum when one exists, else
// the base TS type). This mirrors the pre-alt-source template exactly so that specs
// without alternative sources produce byte-identical output.
func plainYargsFlagExpr(f yargsFlagEntry) string {
	if f.TypeName != "" {
		return fmt.Sprintf("argv.%s as %s", f.FieldName, f.TypeName)
	}
	return fmt.Sprintf("argv.%s as %s", f.FieldName, f.TSType)
}

// scanYargsAltSources reports whether any flag in the slice declares an alternative
// source (hasAlt), and specifically whether any of them is a $FILE source
// (hasFile). A $FILE source requires the generated config code to import a JSONPath library.
func scanYargsAltSources(flags []yargsFlagEntry) (hasAlt, hasFile bool) {
	for _, f := range flags {
		if len(f.AltSources) > 0 {
			hasAlt = true
		}
		for _, src := range f.AltSources {
			if src.Type == "$FILE" {
				hasFile = true
			}
		}
	}
	return hasAlt, hasFile
}

// yargsAltSourceResolvers maps a TypeScript flag type to the generated resolver
// function in gencli/config.ts that falls back to alternative sources when the
// flag was not provided on the command line.
var yargsAltSourceResolvers = map[string]string{
	"string":    "resolveStringFlag",
	"number":    "resolveNumberFlag",
	"boolean":   "resolveBoolFlag",
	"string[]":  "resolveStringSliceFlag",
	"number[]":  "resolveNumberSliceFlag",
	"boolean[]": "resolveBoolSliceFlag",
}

// formatTSAltSources formats a slice of spec.AlternativeSource as a TypeScript array
// literal of the generated AltSource type, for use in generated template code.
func formatTSAltSources(sources []spec.AlternativeSource) string {
	if len(sources) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(sources))
	for _, s := range sources {
		parts = append(parts, fmt.Sprintf("{ type: %q, property: %q }", s.Type, s.Property))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// yargsAltSourceNames returns every option spelling that can reference the flag on the
// command line: its camelCase field name (always first — used to read argv), its raw
// spec name when different (yargs accepts kebab-case too), and all aliases. The generated
// wasSetOnCli scanner matches process.argv tokens against this list, mirroring how pflag's
// Changed and urfave/cli's IsSet account for shorthands and aliases.
func yargsAltSourceNames(f yargsFlagEntry) []string {
	seen := make(map[string]bool)
	var names []string
	add := func(s string) {
		if s != "" && !seen[s] {
			seen[s] = true
			names = append(names, s)
		}
	}
	add(f.FieldName)
	add(f.RawName)
	add(f.Shorthand)
	for _, a := range f.ExtraAliases {
		add(a)
	}
	return names
}

// formatTSAltNames formats the accepted option-name list as a TypeScript string array literal.
func formatTSAltNames(names []string) string {
	parts := make([]string, 0, len(names))
	for _, n := range names {
		parts = append(parts, fmt.Sprintf("%q", n))
	}
	return "[" + strings.Join(parts, ", ") + "]"
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

// yargsVariadicCoerce returns a pre-rendered JavaScript arrow function that maps
// the array elements of a non-string variadic flag to their proper JS type, or ""
// when no coercion is needed (non-variadics and string variadics). Yargs has no
// native typed-array support: with only `type: "array"`, elements arrive as
// strings (and the parser's numeric heuristic misses non-integers), so we coerce
// explicitly. The Array.isArray guard keeps absent flags undefined rather than
// turning them into empty arrays.
func yargsVariadicCoerce(t string) string {
	switch t {
	case "integer", "number":
		return "(v) => Array.isArray(v) ? v.map(Number) : v"
	case "boolean":
		return `(v) => Array.isArray(v) ? v.map((x) => x === true || x === "true") : v`
	default:
		return ""
	}
}

// yargsDefaultVal returns a TypeScript literal default value for a flag. The
// codec normalizes defaults to string/int64/float64/bool (and typed slices for
// variadic flags); the other scalar cases are kept as defensive fallbacks.
func yargsDefaultVal(val any) string {
	switch v := val.(type) {
	case []string:
		parts := make([]string, len(v))
		for i, s := range v {
			parts[i] = fmt.Sprintf("\"%s\"", strings.ReplaceAll(s, "\"", "\\\""))
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case []int64:
		parts := make([]string, len(v))
		for i, n := range v {
			parts[i] = strconv.FormatInt(n, 10)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case []float64:
		parts := make([]string, len(v))
		for i, f := range v {
			parts[i] = strconv.FormatFloat(f, 'f', -1, 64)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case []bool:
		parts := make([]string, len(v))
		for i, b := range v {
			parts[i] = strconv.FormatBool(b)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case string:
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(fmt.Sprintf("%s", val), "\"", "\\\""))
	case int, int32, int64, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case bool:
		return fmt.Sprintf("%t", val)
	default:
		return ""
	}
}
