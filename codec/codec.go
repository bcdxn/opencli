// Package codec provides capabilities for marshalling and unmarshalling
// OpenCLI documents.
//
// It supports JSON and YAML formats and automatically validates specs against
// the official OpenCLI schema definitions.
package codec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/bcdxn/opencli/internal/ds"
	"github.com/bcdxn/opencli/spec"
	yaml "github.com/goccy/go-yaml"
)

// Format represents a supported serialization format for OpenCLI documents.
type Format string

const (
	// FormatYAML is the YAML serialization format.
	FormatYAML Format = "yaml"
	// FormatJSON is the JSON serialization format.
	FormatJSON Format = "json"
)

// UnmarshalJSON decodes a JSON-encoded OpenCLI document into a spec.Document.
//
// Example:
//
//	doc, err := codec.UnmarshalJSON([]byte(`{"opencliVersion": "0.1.0", ...}`))
//	if err != nil {
//		log.Fatal(err)
//	}
func UnmarshalJSON(data []byte) (*spec.Document, error) {
	if data == nil {
		return nil, errors.New("spec document is nil")
	}

	var rawDoc rawDocument
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		return nil, err
	}

	doc, err := buildSpecDoc(&rawDoc)
	if err != nil {
		return nil, fmt.Errorf("error building OpenCLI spec document: %w", err)
	}

	return doc, nil
}

// UnmarshalYAML decodes a YAML-encoded OpenCLI document into a spec.Document.
//
// Example:
//
//	doc, err := codec.UnmarshalYAML([]byte("opencliVersion: 0.1.0\n..."))
//	if err != nil {
//	    log.Fatal(err)
//	}
func UnmarshalYAML(data []byte) (*spec.Document, error) {
	if data == nil {
		return nil, errors.New("spec document is nil")
	}

	var rawDoc rawDocument
	if err := yaml.Unmarshal(data, &rawDoc); err != nil {
		return nil, err
	}

	doc, err := buildSpecDoc(&rawDoc)
	if err != nil {
		return nil, fmt.Errorf("error building OpenCLI spec document: %w", err)
	}

	return doc, nil
}

func buildSpecDoc(rawDoc *rawDocument) (*spec.Document, error) {
	// Build hierarchical data structure
	var doc spec.Document

	doc.Global = rawDoc.Global
	if doc.Global != nil {
		// Global flags are not part of the command Trie, so postProcessing never
		// reaches them. Normalize their defaults here to keep every flag default
		// in canonical Go types regardless of where it is declared.
		for i := range doc.Global.Flags {
			if err := normalizeOneFlagDefault(&doc.Global.Flags[i]); err != nil {
				return nil, fmt.Errorf("global flag: %w", err)
			}
		}
	}
	doc.Info = rawDoc.Info
	doc.Install = rawDoc.Install
	doc.OpenCLIVersion = rawDoc.OpenCLIVersion

	for _, rawCmd := range rawDoc.Commands.Entries() {
		insertCommand(&doc, rawCmd.Key, rawCmd.Value)
	}
	// run post processing to add/update values after building hierarchical command structure
	if err := postProcessing(&doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

// MarshalJSON encodes a spec.Document into JSON bytes.
//
// Example:
//
//	data, err := codec.MarshalJSON(doc)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(data))
func MarshalJSON(doc *spec.Document) ([]byte, error) {
	if doc == nil {
		return nil, errors.New("spec document is nil")
	}

	rawDoc := convertToRawDoc(doc)
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(rawDoc)
	if err != nil {
		return nil, fmt.Errorf("error marshaling OpenCLI spec document: %w", err)
	}

	return buf.Bytes(), nil
}

// MarshalYAML encodes a spec.Document into YAML bytes.
//
// Example:
//
//	data, err := codec.MarshalYAML(doc)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(data))
func MarshalYAML(doc *spec.Document) ([]byte, error) {
	if doc == nil {
		return nil, errors.New("spec document is nil")
	}

	rawDoc := convertToRawDoc(doc)

	data, err := yaml.MarshalWithOptions(rawDoc, yaml.UseLiteralStyleIfMultiline(true))
	if err != nil {
		return data, fmt.Errorf("error marshaling OpenCLI spec document: %w", err)
	}

	return data, nil
}

var (
	wsRE     = regexp.MustCompile(`[^\S\r\n]+`)         // whitespace
	paramsRE = regexp.MustCompile(`[^\S\r\n][^A-Za-z]`) // end of command/beginning of nested commands, args, flags
)

// insertcommand adds a command in a Trie-like structure within the document.
// This ensures that all segments along a path are represented.
func insertCommand(doc *spec.Document, rawCmdLine string, rawCmd rawCommandItem) {
	cmdLine := paramsRE.Split(rawCmdLine, -1)[0] // the literal command as written in the spec
	cmdSegments := wsRE.Split(cmdLine, -1)
	// root command is nil, add it
	if doc.Commands == nil {
		doc.Commands = &spec.CommandItem{
			Segment:     cmdSegments[0],
			CommandLine: cmdSegments[0],
			Derived:     true,
		}
	}
	// we are at the root node
	if len(cmdSegments) == 1 {
		setCommandFields(doc.Commands, rawCmdLine, rawCmd)
	}
	// add each segment of the command hierarchically
	node := doc.Commands
	cmdLineBuilder := []string{cmdSegments[0]}
	for i := 1; i < len(cmdSegments); i++ {
		cmdIndex := indexOfSubcommand(node, cmdSegments[i])
		cmdLineBuilder = append(cmdLineBuilder, cmdSegments[i])

		if cmdIndex < 0 {
			// the command node does not exist in the Trie, add it as a derived command
			newNode := &spec.CommandItem{
				Segment:     cmdSegments[i],
				CommandLine: strings.Join(cmdLineBuilder, " "),
				Derived:     true,
			}
			node.Commands = append(node.Commands, newNode)
			cmdIndex = len(node.Commands) - 1
		}

		if i == len(cmdSegments)-1 { // reached the end of a full 'command line' in the spec document
			// Populate metadata for an existing node once its full command definition is reached.
			setCommandFields(node.Commands[cmdIndex], rawCmdLine, rawCmd)
		}
		// advance our position in the Trie
		node = node.Commands[cmdIndex]
	}
}

// setCommandFields adds the fields from the raw unmarshalled command onto the given spec.CommandItem.
func setCommandFields(cmd *spec.CommandItem, rawCmdLine string, rawCmd rawCommandItem) {
	cmd.Derived = false
	cmd.CommandLineRaw = rawCmdLine
	cmd.Summary = rawCmd.Summary
	cmd.Description = rawCmd.Description
	cmd.Aliases = rawCmd.Aliases
	cmd.Args = rawCmd.Args
	cmd.Flags = rawCmd.Flags
	cmd.Hidden = rawCmd.Hidden
	cmd.Kind = rawCmd.Kind
	cmd.ExitCodes = rawCmd.ExitCodes
	cmd.Examples = rawCmd.Examples
}

// indexOfSubcommand finds the index of the command in the Trie-like structure's node.
// if the segment has not been added to the Trie yet, it will return -1
func indexOfSubcommand(cmd *spec.CommandItem, segment string) int {
	for i, cmd := range cmd.Commands {
		if cmd.Segment == segment {
			return i
		}
	}
	return -1
}

// postProcessing is applied to the Document which traverses the command Items,
// processing each command in the Trie.
func postProcessing(doc *spec.Document) error {
	if doc.Commands == nil {
		return nil
	}

	return postProcessingDFS(doc.Commands)
}

// postProcessingDFS is a recursive function that processes each command item.
func postProcessingDFS(node *spec.CommandItem) error {
	if node == nil {
		return nil
	}
	// process node

	// 1. If the command is not defined explicitly in the doc, then it must be a grouping
	if node.Derived {
		node.Kind = spec.CommandKindGroup
	}
	// 2. If a command has subcommands it should be indicated
	for _, child := range node.Commands {
		if !child.Hidden {
			node.VisibleChildren = true
			break
		}
	}
	// 3. If a command has non-hidden arguments it should be indicated
	for _, arg := range node.Args {
		if !arg.Hidden && !arg.Passthrough {
			node.VisibleArgs = true
			break
		} else if !arg.Hidden {
			node.VisiblePassthroughArgs = true
			break
		}
	}
	// 4. if a command has non-hidden flags it should be indicated
	for _, flag := range node.Flags {
		if !flag.Hidden {
			node.VisibleFlags = true
			break
		}
	}
	// 5. If a command's subcommands have arguments, it should be indicated
	node.VisibleChildrenArgs = node.VisibleChildren && visibleChildrenArgs(node)
	// 6. If a command's subcommands have flags, it should be indicated
	node.VisibleChildrenFlags = node.VisibleChildren && visibleChildrenFlags(node)
	// 7. Add the arguments modifiers
	addModifiers(node)
	// 8. Normalize flag default values to canonical Go types so downstream
	// consumers (code generation, docs) don't depend on decoder-specific types
	if err := normalizeFlagDefaults(node); err != nil {
		return fmt.Errorf("command %q: %w", node.CommandLine, err)
	}
	// iterate through the node's subcommands and process each child recursively
	for _, child := range node.Commands {
		if err := postProcessingDFS(child); err != nil {
			return err
		}
	}

	return nil
}

// normalizeFlagDefaults coerces every flag default value in a command to a
// canonical Go type based on the declared flag type. Decoders produce different
// concrete types for the same literal (e.g. goccy/go-yaml decodes integers as
// uint64 while encoding/json uses float64), so downstream consumers must not
// rely on the raw decoded type. Canonical forms: string -> string, integer ->
// int64, number -> float64, boolean -> bool. Variadic flags may carry list
// defaults (rejected by schema validation but still decodable); their elements
// are coerced to []string, []int64, []float64, or []bool respectively so the
// emitters can render them without panicking on decoder-specific shapes. A
// default whose value cannot be represented as its declared type is a decode
// error rather than something silently dropped downstream.
func normalizeFlagDefaults(node *spec.CommandItem) error {
	if node == nil {
		return nil
	}

	for i := range node.Flags {
		if err := normalizeOneFlagDefault(&node.Flags[i]); err != nil {
			return err
		}
	}

	return nil
}

// normalizeOneFlagDefault coerces a single flag's default value to its canonical
// Go type based on the declared flag type. It is shared by command flags (via
// normalizeFlagDefaults) and global flags, which live outside the command Trie.
func normalizeOneFlagDefault(flag *spec.FlagItem) error {
	var (
		norm any
		err  error
	)
	switch flag.Type {
	case "string":
		norm, err = toStringDefault(flag.Default, flag.Name)
	case "integer":
		norm, err = toInt64Default(flag.Default, flag.Name)
	case "number":
		norm, err = toFloat64Default(flag.Default, flag.Name)
	case "boolean":
		norm, err = toBoolDefault(flag.Default, flag.Name)
	default:
		// Unknown or unset type: leave the value untouched.
		return nil
	}
	if err != nil {
		return err
	}
	flag.Default = norm

	return nil
}

func toStringDefault(val any, name string) (any, error) {
	switch v := val.(type) {
	case nil:
		return nil, nil
	case []any:
		out := make([]string, 0, len(v))
		for _, e := range v {
			s, ok := coerceString(e)
			if !ok {
				return nil, fmt.Errorf("flag %q has a non-string element in its default list", name)
			}
			out = append(out, s)
		}
		return out, nil
	case string:
		return v, nil
	default:
		s, ok := coerceString(val)
		if !ok {
			return nil, fmt.Errorf("flag %q has a non-string default value", name)
		}
		return s, nil
	}
}

func toInt64Default(val any, name string) (any, error) {
	switch v := val.(type) {
	case nil:
		return nil, nil
	case []any:
		out := make([]int64, 0, len(v))
		for _, e := range v {
			n, ok := coerceInt64(e)
			if !ok {
				return nil, fmt.Errorf("flag %q has a non-integer element in its default list", name)
			}
			out = append(out, n)
		}
		return out, nil
	default:
		n, ok := coerceInt64(val)
		if !ok {
			return nil, fmt.Errorf("flag %q has an invalid integer default value", name)
		}
		return n, nil
	}
}

func toFloat64Default(val any, name string) (any, error) {
	switch v := val.(type) {
	case nil:
		return nil, nil
	case []any:
		out := make([]float64, 0, len(v))
		for _, e := range v {
			f, ok := coerceFloat64(e)
			if !ok {
				return nil, fmt.Errorf("flag %q has a non-numeric element in its default list", name)
			}
			out = append(out, f)
		}
		return out, nil
	default:
		f, ok := coerceFloat64(val)
		if !ok {
			return nil, fmt.Errorf("flag %q has an invalid number default value", name)
		}
		return f, nil
	}
}

func toBoolDefault(val any, name string) (any, error) {
	switch v := val.(type) {
	case nil:
		return nil, nil
	case []any:
		out := make([]bool, 0, len(v))
		for _, e := range v {
			b, ok := coerceBool(e)
			if !ok {
				return nil, fmt.Errorf("flag %q has a non-boolean element in its default list", name)
			}
			out = append(out, b)
		}
		return out, nil
	default:
		b, ok := coerceBool(val)
		if !ok {
			return nil, fmt.Errorf("flag %q has an invalid boolean default value", name)
		}
		return b, nil
	}
}

// coerceString converts a decoded scalar to its string form. Decoders may yield
// strings directly or numeric/boolean scalars for loosely-typed documents.
func coerceString(val any) (string, bool) {
	switch v := val.(type) {
	case nil:
		return "", false
	case string:
		return v, true
	case int64:
		return strconv.FormatInt(v, 10), true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case float64:
		if math.Trunc(v) == v && !math.IsInf(v, 0) {
			return strconv.FormatInt(int64(v), 10), true
		}
		return strconv.FormatFloat(v, 'f', -1, 64), true
	case bool:
		return strconv.FormatBool(v), true
	default:
		return "", false
	}
}

// coerceInt64 converts a decoded scalar to int64. Both decoders may represent
// integers as uint64 (goccy/go-yaml) or float64 (encoding/json).
func coerceInt64(val any) (int64, bool) {
	switch v := val.(type) {
	case nil:
		return 0, false
	case int64:
		return v, true
	case uint64:
		if v > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case float64:
		if v != math.Trunc(v) || v < math.MinInt64 || v > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	default:
		return 0, false
	}
}

// coerceFloat64 converts a decoded scalar to float64.
func coerceFloat64(val any) (float64, bool) {
	switch v := val.(type) {
	case nil:
		return 0, false
	case int64:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// coerceBool converts a decoded scalar to bool. Only real booleans are accepted;
// coercing numbers or strings would silently change the flag's semantics.
func coerceBool(val any) (bool, bool) {
	switch v := val.(type) {
	case nil:
		return false, false
	case bool:
		return v, true
	default:
		return false, false
	}
}

// visibleChildrenArgs returns true if at least one subcommand takes arguments
func visibleChildrenArgs(node *spec.CommandItem) bool {
	if node == nil {
		return false
	}

	for _, child := range node.Commands {
		for _, arg := range child.Args {
			if !arg.Hidden {
				return true
			}
		}
		// recursively check for children
		if visibleChildrenArgs(child) {
			return true
		}
	}

	return false
}

// visibleChildrenFlags returns true if at least one subcommand takes arguments
func visibleChildrenFlags(node *spec.CommandItem) bool {
	if node == nil {
		return false
	}

	for _, child := range node.Commands {
		for _, flag := range child.Flags {
			if !flag.Hidden {
				return true
			}
		}
		// recursively check for children
		if visibleChildrenFlags(child) {
			return true
		}
	}

	return false
}

// argsModifiers returns the list of arguments modifiers that will display in documentation or the
// the spec command line
func addModifiers(node *spec.CommandItem) {

	if node.VisibleChildrenArgs || node.VisibleArgs {
		node.ArgsModifiers = []string{"<arguments>"}
	} else {
		node.ArgsModifiers = []string{}
	}

	if node.VisiblePassthroughArgs {
		node.PassthroughArgsModifiers = []string{"--"}
		for _, arg := range node.Args {
			if arg.Passthrough {
				node.PassthroughArgsModifiers = append(node.PassthroughArgsModifiers, fmt.Sprintf("<%s>", arg.Name))
			}
		}
	} else {
		node.PassthroughArgsModifiers = []string{}
	}

	if node.VisibleChildrenFlags || node.VisibleFlags {
		node.FlagsModifiers = []string{"[flags]"}
	} else {
		node.FlagsModifiers = []string{}
	}

	if node.VisibleChildren {
		node.CommandModifiers = []string{"{command}"}
	} else {
		node.CommandModifiers = []string{}
	}

}

func convertToRawDoc(doc *spec.Document) *rawDocument {
	rawDoc := &rawDocument{
		OpenCLIVersion: doc.OpenCLIVersion,
		Info:           doc.Info,
		Install:        doc.Install,
		Global:         doc.Global,
		Commands:       ds.NewMap[string, rawCommandItem](),
	}

	convertToRawCommands(doc.Commands, rawDoc)

	return rawDoc
}

func convertToRawCommands(cmd *spec.CommandItem, rawDoc *rawDocument) {
	rawCmd := rawCommandItem{
		Summary:     cmd.Summary,
		Description: cmd.Description,
		Aliases:     cmd.Aliases,
		Args:        cmd.Args,
		Flags:       cmd.Flags,
		Hidden:      cmd.Hidden,
		Kind:        cmd.Kind,
		ExitCodes:   cmd.ExitCodes,
	}

	if len(rawCmd.Kind) == 0 {
		// default kind to action if not specified. This may be overridden later if the command is a
		// derived command (i.e. it wasn't declared in the OCS document)
		rawCmd.Kind = spec.CommandKindAction
	}

	cmdLine := cmd.CommandLine
	if len(cmd.CommandModifiers) > 0 {
		cmdLine += " " + strings.Join(cmd.CommandModifiers, " ")
	}
	if len(cmd.ArgsModifiers) > 0 {
		cmdLine += " " + strings.Join(cmd.ArgsModifiers, " ")
	}
	if len(cmd.FlagsModifiers) > 0 {
		cmdLine += " " + strings.Join(cmd.FlagsModifiers, " ")
	}

	rawDoc.Commands.Set(cmdLine, rawCmd)

	for _, child := range cmd.Commands {
		convertToRawCommands(child, rawDoc)
	}
}
