package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/bcdxn/opencli/spec"
	yaml "github.com/goccy/go-yaml"
)

// Format indicates input/output format
type Format string

const (
	FormatYAML Format = "yaml"
	FormatJSON Format = "json"
)

// UnmarshalJSON decodes serialized JSON bytes into a Spec.
func UnmarshalJSON(data []byte) (*spec.Document, error) {
	if data == nil {
		return nil, errors.New("spec document is nil")
	}

	var rawDoc rawDocument
	if err := json.Unmarshal(data, &rawDoc); err != nil {
		return nil, err
	}

	doc := buildSpecDoc(rawDoc)

	return doc, nil
}

// UnmarshalYAML decodes serialized YAML bytes into a Spec.
func UnmarshalYAML(data []byte) (*spec.Document, error) {
	if data == nil {
		return nil, errors.New("spec document is nil")
	}

	var rawDoc rawDocument
	if err := yaml.Unmarshal(data, &rawDoc); err != nil {
		return nil, err
	}

	doc := buildSpecDoc(rawDoc)

	return doc, nil
}

func buildSpecDoc(rawDoc rawDocument) *spec.Document {
	// Build hierarchical data structure
	var doc spec.Document

	doc.Global = rawDoc.Global
	doc.Info = rawDoc.Info
	doc.Install = rawDoc.Install
	doc.OpenCLIVersion = rawDoc.OpenCLIVersion

	for rawCmdLine, rawCmd := range rawDoc.Commands {
		// memoize the raw command line string value (without the {commands} <arguments> [flags] modifiers)
		rawDoc.MemoizedCommandLines = append(rawDoc.MemoizedCommandLines, paramsRE.Split(rawCmdLine, -1)[0])
		insertCommand(&doc, rawCmdLine, rawCmd)
	}
	// run post processing to add/update values after building hierarchical command structure
	postProcessing(&doc, &rawDoc)

	return &doc
}

// // MarshalJSON encodes Spec into bytes in the requested format.
// func MarshalJSON(s *spec.Spec) ([]byte, error) {
// 	if s == nil {
// 		return nil, errors.New("spec is nil")
// 	}

// 	data, err := json.MarshalIndent(s, "", "  ")
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshaling spec to JSON: %w", err)
// 	}

// 	return data, nil
// }

// // Marshal encodes Spec into bytes in the requested format.
// func MarshalYAML(s *spec.Spec) ([]byte, error) {
// 	if s == nil {
// 		return nil, errors.New("spec is nil")
// 	}

// 	data, err := yaml.Marshal(s)
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshaling spec to YAML: %w", err)
// 	}

// 	return data, nil
// }

var (
	wsRE     = regexp.MustCompile(`[^\S\r\n]+`)         // whitespace
	paramsRE = regexp.MustCompile(`[^\S\r\n][^A-Za-z]`) // end of command/beginning of nested commands, args, flags
	cmdsRE   = regexp.MustCompile(`{[^}]+}`)            // nested commands
	argsRE   = regexp.MustCompile(`<[^>]+>`)            // arguments
	flagsRE  = regexp.MustCompile(`\[[^\]]+]`)          // flags
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
		}
	}
	// add each segment of the command hierarchically
	node := doc.Commands
	cmdLineBuilder := []string{cmdSegments[0]}
	for i := 1; i < len(cmdSegments); i++ {
		cmdIndex := indexOfSubcommand(node, cmdSegments[i])
		cmdLineBuilder = append(cmdLineBuilder, cmdSegments[i])

		if cmdIndex < 0 {
			// the command node does not exist in the Trie, add it
			newNode := &spec.CommandItem{
				Segment:        cmdSegments[i],
				CommandLine:    strings.Join(cmdLineBuilder, " "),
				CommandLineRaw: rawCmdLine,
				Summary:        rawCmd.Summary,
				Description:    rawCmd.Description,
				Aliases:        rawCmd.Aliases,
				Args:           rawCmd.Args,
				Flags:          rawCmd.Flags,
				Hidden:         rawCmd.Hidden,
				Group:          rawCmd.Group,
				ExitCodes:      rawCmd.ExitCodes,
			}
			node.Commands = append(node.Commands, newNode)
			cmdIndex = len(node.Commands) - 1
		} else if i == len(cmdSegments)-1 { // leaf node
			// We may have created the node in the Trie as part of a longer command line string
			// Populate metadata for an existing node once its full command definition is reached.
			node.Commands[cmdIndex].CommandLineRaw = rawCmdLine
			node.Commands[cmdIndex].Summary = rawCmd.Summary
			node.Commands[cmdIndex].Description = rawCmd.Description
			node.Commands[cmdIndex].Aliases = rawCmd.Aliases
			node.Commands[cmdIndex].Args = rawCmd.Args
			node.Commands[cmdIndex].Flags = rawCmd.Flags
			node.Commands[cmdIndex].Hidden = rawCmd.Hidden
			node.Commands[cmdIndex].Group = rawCmd.Group
			node.Commands[cmdIndex].ExitCodes = rawCmd.ExitCodes
		}
		// advance our position in the Trie
		node = node.Commands[cmdIndex]
	}
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
func postProcessing(doc *spec.Document, rawDoc *rawDocument) {
	if doc.Commands == nil {
		return
	}

	postProcessingDFS(doc.Commands, rawDoc)
}

// postProcessingDFS is a recursive function that processes each command item.
func postProcessingDFS(node *spec.CommandItem, rawDoc *rawDocument) {
	if node == nil {
		return
	}

	// process node

	// 1. If the command is not defined explicitly in the doc, then it must be a grouping
	node.Group = !slices.Contains(rawDoc.MemoizedCommandLines, node.CommandLine)
	// 2. If a command has subcommands it should be indicated
	node.Children = len(node.Commands) > 0
	// 3. If a command's subcommands have arguments, it should be indicated
	node.ChildrenArgs = node.Children && childrenArgs(node)
	// 4. If a command's subcommands have flags, it should be indicated
	node.ChildrenFlags = node.Children && childrenFlags(node)
	// 5. Add the arguments modifiers
	addModifiers(node)

	// iterate through the nodes subcommands and process each child recursively
	for _, child := range node.Commands {
		postProcessingDFS(child, rawDoc)
	}
}

// childrenArgs returns true if at least one subcommand takes arguments
func childrenArgs(node *spec.CommandItem) bool {
	if node == nil {
		return false
	}

	for _, child := range node.Commands {
		if len(child.Args) > 0 {
			return true
		}
		// recursively check for children
		if childrenArgs(child) {
			return true
		}
	}

	return false
}

// childrenArgs returns true if at least one subcommand takes arguments
func childrenFlags(node *spec.CommandItem) bool {
	if node == nil {
		return false
	}

	for _, child := range node.Commands {
		if len(child.Flags) > 0 {
			return true
		}
		// recursively check for children
		if childrenFlags(child) {
			return true
		}
	}

	return false
}

// argsModifiers returns the list of arguments modifiers that will display in documentation or the
// the spec command line
func addModifiers(node *spec.CommandItem) {
	argModifiers := []string{}

	if !node.Children {
		for _, arg := range node.Args {
			argModifiers = append(argModifiers, fmt.Sprintf("<%s>", arg.Name))
		}
	} else if node.ChildrenArgs {
		argModifiers = append(argModifiers, "<arguments>")
	}

	if node.ChildrenFlags || len(node.Flags) > 0 {
		node.FlagModifiers = []string{"[flags]"}
	} else {
		node.FlagModifiers = []string{}
	}

	if node.Children {
		node.CommandModifiers = []string{"{command}"}
	} else {
		node.CommandModifiers = []string{}
	}

	node.ArgsModifiers = argModifiers
}
