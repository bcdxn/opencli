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

	for _, rawCmd := range rawDoc.Commands.Entries() {
		// memoize the raw command line string value (without the {commands} <arguments> [flags] modifiers)
		rawDoc.MemoizedCommandLines = append(rawDoc.MemoizedCommandLines, paramsRE.Split(rawCmd.Key, -1)[0])
		insertCommand(&doc, rawCmd.Key, rawCmd.Value)
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
				Segment:     cmdSegments[i],
				CommandLine: strings.Join(cmdLineBuilder, " "),
			}
			node.Commands = append(node.Commands, newNode)
			cmdIndex = len(node.Commands) - 1
		}

		if i == len(cmdSegments)-1 { // leaf node
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
	node.Group = node.Group || !slices.Contains(rawDoc.MemoizedCommandLines, node.CommandLine)
	// 2. If a command has subcommands it should be indicated
	for _, child := range node.Commands {
		if !child.Hidden {
			node.VisibleChildren = true
			break
		}
	}
	// 3. If a command has non-hidden arguments it should be indicated
	for _, arg := range node.Args {
		if !arg.Hidden {
			node.VisibleArgs = true
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
	// 3. If a command's subcommands have arguments, it should be indicated
	node.VisibleChildrenArgs = node.VisibleChildren && visibleChildrenArgs(node)
	// 4. If a command's subcommands have flags, it should be indicated
	node.VisibleChildrenFlags = node.VisibleChildren && visibleChildrenFlags(node)
	// 5. Add the arguments modifiers
	addModifiers(node)

	// iterate through the nodes subcommands and process each child recursively
	for _, child := range node.Commands {
		postProcessingDFS(child, rawDoc)
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

	if !node.VisibleChildren {
		argModifiers := []string{}
		for _, arg := range node.Args {
			if !arg.Hidden {
				argModifiers = append(argModifiers, fmt.Sprintf("<%s>", arg.Name))
			}
		}
		node.ArgsModifiers = argModifiers
	} else if node.VisibleChildrenArgs {
		node.ArgsModifiers = []string{"<arguments>"}
	}

	if node.VisibleChildrenFlags || node.VisibleFlags {
		node.FlagModifiers = []string{"[flags]"}
	} else {
		node.FlagModifiers = []string{}
	}

	if node.VisibleChildren {
		node.CommandModifiers = []string{"{command}"}
	} else {
		node.CommandModifiers = []string{}
	}

}
