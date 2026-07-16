// Package ocobra provides the ability to create OpenCLI documents from existing Cobra CLIs.
//
// It attaches a hidden subcommand to an existing Cobra root command that, when
// invoked, walks the entire command tree and emits a valid OpenCLI spec document.
//
// Typical usage:
//
//	var rootCmd = &cobra.Command{Use: "myapp", Short: "An awesome CLI"}
//	opencobra.FromCommand(rootCmd)
//	rootCmd.Execute()
package ocobra

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
// If omitted, only binary name (derived from root.Use) will be populated.
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
// Running that subcommand walks the Cobra command tree and writes an OpenCLI
// spec document to the configured output (default: stdout).
func FromCommand(rootCmd *cobra.Command, opts ...Option) {
	c := defaultConfigWithOptions(opts...)
	var flagOut string

	genCmd := &cobra.Command{
		Use:    "__opencli",
		Short:  "Generates OpenCLI specification",
		Hidden: true,
		Args:   cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flagOut != "" {
				// The user has passed a file to output to
				f, err := os.OpenFile(flagOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					return fmt.Errorf("unable to write opencli doc to specified file: %w", err)
				}
				defer f.Close()

				c.output = f
			}
			doc := documentFromCommand(rootCmd, c)
			return writeDoc(c.output, doc, c.format)
		},
	}

	genCmd.Flags().StringVarP(&flagOut, "out", "o", "", "The path to the directory where the generated code will be output")

	rootCmd.AddCommand(genCmd)
}

// GenerateDocument walks a Cobra command tree and returns the corresponding
// OpenCLI spec.Document. This is exported so callers can inspect or modify
// the document before serializing it themselves.
func GenerateDocument(rootCmd *cobra.Command, opts ...Option) *spec.Document {
	c := defaultConfigWithOptions(opts...)
	return documentFromCommand(rootCmd, c)
}

func documentFromCommand(rootCmd *cobra.Command, c *config) *spec.Document {
	var doc spec.Document
	doc.OpenCLIVersion = defaultOpenCLIVersion
	doc.Info = c.getInfo()
	if doc.Info.Binary == "" {
		doc.Info.Binary = GetBinaryName(rootCmd)
	}
	doc.Install = c.install

	if doc.Info.Description == "" {
		doc.Info.Description = rootCmd.Long
	}
	if doc.Info.Title == "" {
		doc.Info.Title = rootCmd.Short
	}
	if doc.Info.Version == "" {
		doc.Info.Version = rootCmd.Version
	}

	// Build global section
	if len(c.globalFlags) > 0 {
		if doc.Global == nil {
			doc.Global = &spec.Global{}
		}
		doc.Global.Flags = c.globalFlags
	}

	// add global flags derived from persistent flags set on the root
	if pFlags := parseGlobalFlags(rootCmd); len(pFlags) > 0 {
		if doc.Global == nil {
			doc.Global = &spec.Global{}
		}
		doc.Global.Flags = append(doc.Global.Flags, pFlags...)
	}

	// Walk the command tree
	doc.Commands = commandFromCobra(rootCmd, nil)

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

// commandFromCobra converts a cobra.Command (and its children) into a
// spec.CommandItem tree.
func commandFromCobra(cmd *cobra.Command, parent *spec.CommandItem) *spec.CommandItem {
	segment := cmd.Name()
	fullLine := segment
	if parent != nil {
		// we're not at the root and must join previous commands to get the full command line
		fullLine = strings.Join([]string{parent.CommandLine, segment}, " ")
	}

	// Parse arguments once and share with getModifiers to avoid double work.
	args := ParseUse(cmd.Use)
	cmdMods, argMods, flagMods, passthroughMods := getModifiers(cmd, args)

	item := &spec.CommandItem{
		Segment:                  segment,
		CommandLine:              fullLine,
		Summary:                  cmd.Short,
		Description:              cmd.Long,
		Hidden:                   cmd.Hidden,
		CommandModifiers:         cmdMods,
		ArgsModifiers:            argMods,
		FlagsModifiers:           flagMods,
		PassthroughArgsModifiers: passthroughMods,
		VisibleArgs:              len(argMods) > 0,
		VisibleChildren:          len(cmdMods) > 0,
		VisibleFlags:             len(flagMods) > 0,
		VisiblePassthroughArgs:   len(passthroughMods) > 0,
		Examples: []spec.Example{
			{
				Title:   "Example",
				Content: cmd.Example,
			},
		},
	}

	// Aliases
	item.Aliases = cmd.Aliases

	// Determine kind
	hasRun := cmd.Run != nil || cmd.RunE != nil
	item.Kind = spec.CommandKindGroup
	if hasRun {
		// Both run and subcommands → action (per user preference)
		item.Kind = spec.CommandKindAction
	}

	item.Args = args

	// Convert persistent + local flags
	item.Flags = parseFlags(cmd)

	// Recurse into children
	for _, child := range cmd.Commands() {
		if child.Use == "__opencli" {
			continue // skip our own genator command
		}
		childItem := commandFromCobra(child, item)
		item.Commands = append(item.Commands, childItem)
	}

	return item
}

func getModifiers(cmd *cobra.Command, args []spec.ArgumentItem) ([]string, []string, []string, []string) {
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

	if !cmd.HasAvailableFlags() {
		flagsModifier = []string{}
	}

	cmdsModifier := []string{"{commands}"}
	if !cmd.HasAvailableSubCommands() {
		cmdsModifier = []string{}
	}

	return cmdsModifier, argsModifier, flagsModifier, passthroughModifier
}

// parseFlags merges persistent flags and local flags into a single slice since OpenCLI spec does
// not have a concept of persisted flags other than top-level 'global' flags.
func parseFlags(cmd *cobra.Command) []spec.FlagItem {
	var items []spec.FlagItem

	// Persistent flags are inherited; collect them from the root down to avoid duplicates.
	var seen map[string]bool
	root := cmd.Root()
	// We don't want to add persistent flags from the root command because those will be
	// added via the 'globalFlags' section
	if root.Parent() != nil {
		seen = make(map[string]bool)
		// We only want persistent flags from ancestors, not the command itself
		// (local flags are handled separately).
		var walk func(c *cobra.Command)
		walk = func(c *cobra.Command) {
			if c == nil {
				// base case, we were at the root and there is no parent
				return
			}
			// 'visit' the persisted flags set on the current command
			c.PersistentFlags().VisitAll(func(f *pflag.Flag) {
				if !f.Hidden {
					seen[f.Name] = true
				}
				if f.Hidden || f.Name == "help" || f.Name == "version" {
					return
				}
				items = append(items, flagToItem(f))
			})

			// check persisted flags on parent recursively
			walk(c.Parent())
		}
		walk(cmd.Parent())
	}

	// Add persistent flags defined on this command (that we've not already seen) unless
	// this command is the root command, in which its persistent flags will be added in
	// the `globalFlags` sections.
	if cmd.Parent() != nil {
		cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
			if seen != nil && seen[f.Name] {
				return // already added by ancestor
			}
			if f.Name == "help" || f.Name == "version" {
				return
			}
			items = append(items, flagToItem(f))
		})
	}

	// Local flags
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// if seen != nil && seen[f.Name] {
		// 	return // persistent flag already added
		// }
		if f.Name == "help" || f.Name == "version" {
			return
		}
		items = append(items, flagToItem(f))
	})

	return items
}

// flagToItem converts a pflag.Flag to a spec.FlagItem.
func flagToItem(f *pflag.Flag) spec.FlagItem {
	t, v := resolveFlagType(f)
	item := spec.FlagItem{
		Name:     f.Name,
		Summary:  f.Usage,
		Type:     t,
		Variadic: v,
		Default:  parseDefault(t, v, f.DefValue),
		Hidden:   f.Hidden,
	}
	if f.Shorthand != "" {
		item.Aliases = []string{f.Shorthand}
	}

	return item
}

func parseGlobalFlags(rootCmd *cobra.Command) []spec.FlagItem {
	var items []spec.FlagItem

	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if f.Hidden || f.Name == "help" { // help is already included by default
			return
		}
		items = append(items, flagToItem(f))
	})

	return items
}

// resolveFlagType maps a pflag value type to an OpenCLI type string and if the flag is variadic
func resolveFlagType(f *pflag.Flag) (string, bool) {
	typ := f.Value.Type()
	switch typ {
	case "boolSlice":
		return "bool", true
	case "intSlice":
		return "integer", true
	case "floatSlice":
		return "number", true
	case "stringSlice":
		return "string", true
	case "bool":
		return "boolean", false
	case "string", "count":
		return "string", false
	case "int", "uint", "int8", "int16", "int32", "int64",
		"uint8", "uint16", "uint32", "uint64":
		return "integer", false
	case "float32", "float64":
		return "number", false
	case "duration", "ips", "ipNet", "ipSlice", "cidr", "bytes":
		return "string", false
	default:
		return "string", false
	}
}

// parseDefault attempts to convert the string default value to the
// appropriate Go type so JSON/YAML serialization is correct.
func parseDefault(typ string, variadic bool, val string) any {
	if variadic && (val == "[]" || val == "") {
		return nil
	} else if variadic {
		return strToTypSlice(typ, val)
	}

	return parseDefaultScalar(typ, val, false)
}

func parseDefaultScalar(t string, v string, variadic bool) any {
	// within slices, we don't want to simply return nil because then we end up with
	// `"default": [null]` entries
	if variadic {
		if v == "" {
			return v
		}
		if v == "false" {
			return false
		}
		if v == "0" {
			return int64(0)
		}
		if v == "[]" {
			return []any{}
		}
	}
	// handle zero values, Since this is an interface{} type, we need to return nil instead of
	// standard zero values so that it doesn't get marshalled on omitempty=true properties
	if v == "" {
		return nil
	}
	if v == "false" {
		return nil
	}
	if v == "0" {
		return nil
	}
	if v == "[]" {
		return nil
	}

	switch t {
	case "number":
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	case "integer":
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	case "boolean":
		if v == "true" {
			return true
		} else {
			return nil
		}
	}

	// default string
	return v
}

// strToTypeSlice takes the stringified value and coverts it to a slice of the appropriate value
//
// e.g.:
//
//	"[1,2,3]" --> []int64{1, 2, 3}
//	"[one,two,three]" --> []string{"one", "two", "three"}
func strToTypSlice(t string, s string) []any {
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")

	slc := strings.Split(s, ",")

	d := []any{}
	for _, n := range slc {
		d = append(d, parseDefaultScalar(t, n, true))
	}
	return d
}

// GetBinaryName extracts the binary name from a root command's Use field.
// This is a convenience helper for constructing spec.Info.
func GetBinaryName(rootCmd *cobra.Command) string {
	return strings.Fields(rootCmd.Use)[0]
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
