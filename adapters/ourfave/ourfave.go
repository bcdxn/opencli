// Package ourfave provides the ability to create OpenCLI documents from existing urfave/cli v3 CLIs.
//
// It attaches a hidden subcommand to an existing urfave/cli root command that, when
// invoked, walks the entire command tree and emits a valid OpenCLI spec document.
//
// Typical usage:
//
//	var rootCmd = &cli.Command{Name: "myapp", Usage: "An awesome CLI"}
//	ourfave.FromCommand(rootCmd)
//	rootCmd.Run(os.Args)
package ourfave

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/spec"
	"github.com/urfave/cli/v3"
)

const defaultOpenCLIVersion = "1.0.0-alpha.13"

// config holds the functional-option state.
type config struct {
	info        *spec.Info
	install     []spec.InstallMethod
	globalFlags []spec.FlagItem
	output      io.Writer
	format      codec.Format
}

func (c *config) getInfo() spec.Info {
	if c.info != nil {
		return *c.info
	}
	return spec.Info{}
}

// Option is a functional option for FromCommand.
type Option func(*config)

// WithInfo sets the top-level spec.Info block.
// If omitted, only binary name (derived from root.Name) will be populated.
func WithInfo(info *spec.Info) Option {
	return func(c *config) { c.info = info }
}

// WithInstallMethods sets the install methods list.
func WithInstallMethods(install []spec.InstallMethod) Option {
	return func(c *config) { c.install = install }
}

// WithGlobalFlags sets flags that apply globally (e.g. --help, --version).
func WithGlobalFlags(flags []spec.FlagItem) Option {
	return func(c *config) { c.globalFlags = flags }
}

// WithFormat sets the output format ("yaml" or "json").
// Defaults to YAML.
func WithFormat(f codec.Format) Option {
	return func(c *config) { c.format = f }
}

// WithOutput sets the io.Writer for the generated spec.
// Defaults to os.Stdout if omitted.
func WithOutput(w io.Writer) Option {
	return func(c *config) { c.output = w }
}

// FromCommand attaches a hidden "__opencli" subcommand to rootCmd.
// Running that subcommand walks the urfave/cli command tree and writes an OpenCLI
// spec document to the configured output (default: stdout).
func FromCommand(rootCmd *cli.Command, opts ...Option) {
	c := defaultConfigWithOptions(opts...)

	rootCmd.Commands = append(rootCmd.Commands, &cli.Command{
		Name:   "__opencli",
		Hidden: true,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			outPath := cmd.String("out")
			if outPath != "" {
				f, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					return fmt.Errorf("unable to write opencli doc to specified file: %w", err)
				}
				defer f.Close()

				c.output = f
			}
			doc := documentFromCommand(rootCmd, c)
			return writeDoc(c.output, doc, c.format)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "The path to the file where the generated spec will be output",
			},
		},
	})
}

// GenerateDocument walks a urfave/cli command tree and returns the corresponding
// OpenCLI spec.Document. This is exported so callers can inspect or modify
// the document before serializing it themselves.
func GenerateDocument(rootCmd *cli.Command, opts ...Option) *spec.Document {
	c := defaultConfigWithOptions(opts...)
	return documentFromCommand(rootCmd, c)
}

func documentFromCommand(rootCmd *cli.Command, c *config) *spec.Document {
	var doc spec.Document
	doc.OpenCLIVersion = defaultOpenCLIVersion
	doc.Info = c.getInfo()
	if doc.Info.Binary == "" {
		doc.Info.Binary = GetBinaryName(rootCmd)
	}
	doc.Install = c.install

	if doc.Info.Description == "" {
		doc.Info.Description = rootCmd.Description
	}
	if doc.Info.Title == "" {
		doc.Info.Title = rootCmd.Usage
	}
	if doc.Info.Version == "" {
		doc.Info.Version = rootCmd.Version
	}

	// Build global section from explicitly provided global flags
	if len(c.globalFlags) > 0 {
		if doc.Global == nil {
			doc.Global = &spec.Global{}
		}
		doc.Global.Flags = c.globalFlags
	}

	// Add global flags derived from non-local (persistent) flags set on the root command
	if pFlags := parseGlobalFlags(rootCmd); len(pFlags) > 0 {
		if doc.Global == nil {
			doc.Global = &spec.Global{}
		}
		doc.Global.Flags = append(doc.Global.Flags, pFlags...)
	}

	// Walk the command tree using urfave/cli's built-in Walk
	doc.Commands = commandFromUrfave(rootCmd, nil, true)

	return &doc
}

func writeDoc(w io.Writer, doc *spec.Document, format codec.Format) error {
	var data []byte
	var err error
	switch format {
	case codec.FormatJSON:
		data, err = codec.MarshalJSON(doc)
	default:
		data, err = codec.MarshalYAML(doc)
	}
	if err != nil {
		return fmt.Errorf("marshal spec: %w", err)
	}
	_, err = w.Write(data)
	return err
}

// commandFromUrfave converts a urfave/cli Command into a spec.CommandItem.
// It processes the current command and recursively handles children.
func commandFromUrfave(cmd *cli.Command, parent *spec.CommandItem, isRoot bool) *spec.CommandItem {
	segment := cmd.Name
	fullLine := segment
	if parent != nil {
		fullLine = strings.Join([]string{parent.CommandLine, segment}, " ")
	}

	// Parse arguments from cmd.Arguments (declared arguments only)
	args := parseArguments(cmd.Arguments)

	cmdMods, argMods, flagMods, passthroughMods := getModifiers(cmd, args)

	item := &spec.CommandItem{
		Segment:                  segment,
		CommandLine:              fullLine,
		Summary:                  cmd.Usage,
		Description:              cmd.Description,
		Hidden:                   cmd.Hidden,
		CommandModifiers:         cmdMods,
		ArgsModifiers:            argMods,
		FlagsModifiers:           flagMods,
		PassthroughArgsModifiers: passthroughMods,
		VisibleArgs:              len(argMods) > 0,
		VisibleChildren:          len(cmdMods) > 0,
		VisibleFlags:             len(flagMods) > 0,
		VisiblePassthroughArgs:   len(passthroughMods) > 0,
	}

	// Aliases
	item.Aliases = cmd.Aliases

	// Determine kind: if the command has an Action, it's an action; otherwise a group
	if cmd.Action != nil {
		item.Kind = spec.CommandKindAction
	} else {
		item.Kind = spec.CommandKindGroup
	}

	item.Args = args

	// Convert flags
	item.Flags = parseFlags(cmd, isRoot)

	// Recurse into children (skip our own generator command)
	for _, child := range cmd.Commands {
		if child.Name == "__opencli" {
			continue
		}
		childItem := commandFromUrfave(child, item, false)
		item.Commands = append(item.Commands, childItem)
	}

	return item
}

func getModifiers(cmd *cli.Command, args []spec.ArgumentItem) ([]string, []string, []string, []string) {
	argsModifier := []string{}
	passthroughModifier := []string{}

	for _, arg := range args {
		if arg.Passthrough && !arg.Hidden {
			passthroughModifier = []string{"--", "<arguments>"}
		} else if !arg.Hidden {
			argsModifier = []string{"<arguments>"}
		}
	}

	flagsModifier := []string{"[flags]"}
	if len(cmd.Flags) == 0 {
		flagsModifier = []string{}
	}

	cmdsModifier := []string{"{commands}"}
	if len(cmd.Commands) == 0 {
		cmdsModifier = []string{}
	}

	return cmdsModifier, argsModifier, flagsModifier, passthroughModifier
}

// parseFlags converts urfave/cli flags to spec.FlagItem slice.
// When isRoot is true, only local flags are returned; non-local (persistent)
// flags are handled by parseGlobalFlags and go into the global section instead.
func parseFlags(cmd *cli.Command, isRoot bool) []spec.FlagItem {
	var items []spec.FlagItem

	for _, f := range cmd.Flags {
		// Skip help and version flags (they're handled globally)
		if flagIsHelp(f) || flagIsVersion(f) {
			continue
		}

		// At root level, skip non-local flags since they belong in the global section
		if isRoot && !isLocalFlag(f) {
			continue
		}

		items = append(items, flagToItem(f))
	}

	return items
}

// parseGlobalFlags extracts non-local (persistent/inherited) flags from the root command.
// In urfave/cli v3, a flag with Local=false is inherited by subcommands.
func parseGlobalFlags(rootCmd *cli.Command) []spec.FlagItem {
	var items []spec.FlagItem

	for _, f := range rootCmd.Flags {
		if flagIsHelp(f) {
			continue
		}

		// Check if this flag is non-local (i.e., persistent/inherited by subcommands)
		if !isLocalFlag(f) {
			items = append(items, flagToItem(f))
		}
	}

	return items
}

// isLocalFlag returns true if the flag is marked as local (not inherited by subcommands).
func isLocalFlag(f cli.Flag) bool {
	// FlagBase has IsLocal() method
	if lf, ok := f.(interface{ IsLocal() bool }); ok {
		return lf.IsLocal()
	}
	// Default: flags are not local (they're inherited)
	return false
}

// flagIsHelp returns true if the flag is a help flag.
func flagIsHelp(f cli.Flag) bool {
	names := f.Names()
	for _, n := range names {
		if n == "help" {
			return true
		}
	}
	return false
}

// flagIsVersion returns true if the flag is a version flag.
func flagIsVersion(f cli.Flag) bool {
	names := f.Names()
	for _, n := range names {
		if n == "version" {
			return true
		}
	}
	return false
}

// flagToItem converts a urfave/cli Flag to a spec.FlagItem.
func flagToItem(f cli.Flag) spec.FlagItem {
	names := f.Names()
	name := names[0]
	var aliases []string
	if len(names) > 1 {
		aliases = names[1:]
	}

	item := spec.FlagItem{
		Name:     name,
		Aliases:  aliases,
		Summary:  getFlagUsage(f),
		Type:     resolveFlagType(f),
		Variadic: isVariadicFlag(f),
		Default:  getFlagDefault(f),
		Hidden:   isFlagHidden(f),
		Required: isFlagRequired(f),
	}

	return item
}

// getFlagUsage returns the usage string from a Flag.
func getFlagUsage(f cli.Flag) string {
	if uf, ok := f.(interface{ GetUsage() string }); ok {
		return uf.GetUsage()
	}
	return ""
}

// isFlagHidden returns true if the flag is hidden.
func isFlagHidden(f cli.Flag) bool {
	if hf, ok := f.(interface{ IsVisible() bool }); ok {
		return !hf.IsVisible()
	}
	return false
}

// isFlagRequired returns true if the flag is required.
func isFlagRequired(f cli.Flag) bool {
	if rf, ok := f.(interface{ IsRequired() bool }); ok {
		return rf.IsRequired()
	}
	return false
}

// getFlagDefault returns the default value string from a Flag.
func getFlagDefault(f cli.Flag) any {
	if df, ok := f.(interface{ GetDefaultText() string }); ok {
		return parseDefaultValue(f, df.GetDefaultText())
	}
	return nil
}

// isVariadicFlag returns true if the flag accepts multiple values.
func isVariadicFlag(f cli.Flag) bool {
	if mf, ok := f.(interface{ IsMultiValueFlag() bool }); ok {
		return mf.IsMultiValueFlag()
	}
	return false
}

// resolveFlagType maps a urfave/cli Flag to an OpenCLI type string.
func resolveFlagType(f cli.Flag) string {
	// For slice/array flags, use SchemaItemsType to get the element type
	if sf, ok := f.(interface{ SchemaItemsType() string }); ok && sf.SchemaItemsType() != "" {
		return mapSchemaType(sf.SchemaItemsType())
	}

	// Try to get the schema type from DocGenerationFlag
	if sf, ok := f.(interface{ SchemaType() string }); ok {
		return mapSchemaType(sf.SchemaType())
	}
	return "string"
}

// mapSchemaType converts urfave/cli schema type strings to OpenCLI types.
// urfave/cli already returns OpenCLI-compatible names ("boolean", "integer", "number"),
// so we pass those through directly alongside the Go-type names.
func mapSchemaType(schemaType string) string {
	switch strings.ToLower(schemaType) {
	case "bool", "boolean":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "integer":
		return "integer"
	case "float", "float32", "float64", "number":
		return "number"
	case "string", "time.Duration", "time.Time", "duration":
		return "string"
	default:
		return "string"
	}
}

// parseDefaultValue attempts to convert the string default value to the
// appropriate Go type so JSON/YAML serialization is correct.
func parseDefaultValue(f cli.Flag, val string) any {
	typ := resolveFlagType(f)
	variadic := isVariadicFlag(f)

	if variadic && (val == "[]" || val == "") {
		return nil
	}

	switch typ {
	case "number":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			if f == 0 {
				return nil
			}
			return f
		}
	case "integer":
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			if i == 0 {
				return nil
			}
			return i
		}
	case "boolean":
		if val == "true" {
			return true
		}
		return nil
	}

	// Default: return string value or nil for empty
	if val == "" {
		return nil
	}
	return val
}

// parseArguments converts urfave/cli Arguments to spec.ArgumentItem slice.
// Only declared arguments (via cmd.Arguments) are considered.
func parseArguments(args []cli.Argument) []spec.ArgumentItem {
	if args == nil {
		return nil
	}

	var items []spec.ArgumentItem
	for _, a := range args {
		item := argumentToItem(a)
		items = append(items, item)
	}
	return items
}

// argumentToItem converts a urfave/cli Argument to a spec.ArgumentItem.
func argumentToItem(a cli.Argument) spec.ArgumentItem {
	item := spec.ArgumentItem{
		Name:    getArgumentName(a),
		Type:    resolveArgumentType(a),
		Summary: a.Usage(),
		Default: getArgumentDefault(a),
	}

	return item
}

// getArgumentName extracts the argument name.
func getArgumentName(a cli.Argument) string {
	// Try common argument types
	switch v := a.(type) {
	case *cli.StringArg:
		return v.Name
	case *cli.IntArg:
		return v.Name
	case *cli.UintArg:
		return v.Name
	case *cli.FloatArg:
		return v.Name
	case *cli.Float32Arg:
		return v.Name
	case *cli.Int8Arg:
		return v.Name
	case *cli.Int16Arg:
		return v.Name
	case *cli.Int32Arg:
		return v.Name
	case *cli.Int64Arg:
		return v.Name
	case *cli.Uint8Arg:
		return v.Name
	case *cli.Uint16Arg:
		return v.Name
	case *cli.Uint32Arg:
		return v.Name
	case *cli.Uint64Arg:
		return v.Name
	case *cli.TimestampArg:
		return v.Name
	}
	// Fallback to a static name
	return "arg"
}

// resolveArgumentType maps a urfave/cli Argument to an OpenCLI type string.
func resolveArgumentType(a cli.Argument) string {
	// Check for schema type support
	if sa, ok := a.(interface{ SchemaType() string }); ok {
		return mapSchemaType(sa.SchemaType())
	}

	// Fall back to type-based detection
	switch a.(type) {
	case *cli.StringArg:
		return "string"
	case *cli.IntArg, *cli.Int8Arg, *cli.Int16Arg, *cli.Int32Arg, *cli.Int64Arg,
		*cli.UintArg, *cli.Uint8Arg, *cli.Uint16Arg, *cli.Uint32Arg, *cli.Uint64Arg:
		return "integer"
	case *cli.FloatArg:
		return "number"
	case *cli.TimestampArg:
		return "string"
	}
	return "string"
}

// getArgumentDefault returns the default value from an Argument.
func getArgumentDefault(a cli.Argument) any {
	val := a.Get()
	if val == nil {
		return nil
	}

	// Convert based on type
	switch v := val.(type) {
	case string:
		if v == "" {
			return nil
		}
		return v
	case bool:
		if !v {
			return nil
		}
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		// Return nil for zero values
		fmtStr := fmt.Sprintf("%v", v)
		if fmtStr == "0" {
			return nil
		}
		return v
	case float32, float64:
		fmtStr := fmt.Sprintf("%v", v)
		if fmtStr == "0" {
			return nil
		}
		return v
	}
	return nil
}

// GetBinaryName extracts the binary name from a root command's Name field.
func GetBinaryName(rootCmd *cli.Command) string {
	if rootCmd == nil || len(rootCmd.Name) < 1 {
		// we were given an empty CLI; return a sane default I guess?
		return "root"
	}
	// In urfave/cli, the Name field may contain spaces for the full command line.
	// We take the first word as the binary name.
	return strings.Fields(rootCmd.Name)[0]
}

func defaultConfigWithOptions(opts ...Option) *config {
	c := &config{
		format: codec.FormatYAML,
		output: os.Stdout,
		globalFlags: []spec.FlagItem{
			{
				Name:    "help",
				Summary: "Show help for command",
				Aliases: []string{"h"},
				Type:    "boolean",
			},
			{
				Name:    "version",
				Summary: "Show CLI version",
				Aliases: []string{"v"},
				Type:    "boolean",
			},
		},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}
