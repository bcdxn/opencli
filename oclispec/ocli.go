package oclispec

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"slices"

	"github.com/bcdxn/opencli/internal/oclifile"
	"gopkg.in/yaml.v3"
)

// ocli.go contains the OpenCLI domain types.

// Document represents the OpenCLI document.
type Document struct {
	OpenCLIVersion string
	Info           Info
	Install        []Install
	Global         Global
	CommandTrie    *CommandTrie
}

// Info represents the metadata about the CLI described by the OpenCLI document.
type Info struct {
	Title       string
	Summary     string
	Description string
	License     License
	Contact     Contact
	Binary      string
	Version     string
}

// License represents the license information for the CLI described by the OpenCLI document.
type License struct {
	Name   string
	SpdxID string
	URL    string
}

// Contact represents contact information for maintainers of the CLI described by the OpenCLI document.
type Contact struct {
	Name  string
	Email string
	URL   string
}

// Install represents information about ways a user can install the CLI.
type Install struct {
	Name        string
	Command     string
	URL         string
	Description string
}

// Global contains information that applies to the CLI regardless of command context.
type Global struct {
	ExitCodes []ExitCode
	Flags     []Flag
}

// ExitCode represents a possible exit code of a CLI command
type ExitCode struct {
	Code        int
	Status      string
	Summary     string
	Description string
}

// Command represents n OpenCLI command.
type Command struct {
	Line                 string // The full command line as defined in the OpenCLI Spec Document
	Name                 string // The command part of the command line
	LeafName             string // The final command in the command line
	Params               string // The args/flags part of the command line
	Aliases              []string
	Summary              string
	Description          string
	Arguments            []Argument
	Flags                []Flag
	Hidden               bool
	Group                bool
	CmdSpecificExitCodes []ExitCode
	ExitCodes            []ExitCode
}

// Argument represents an OpenCLI command argument.
type Argument struct {
	Name        string
	Summary     string
	Description string
	Type        string
	Variadic    bool
	Choices     []Choice
	Required    bool
	Default     DefaultValue
}

// Flag represents an OpenCLI command flag.
type Flag struct {
	Name        string
	Aliases     []string
	Hint        string
	Summary     string
	Description string
	Type        string
	Variadic    bool
	Choices     []Choice
	Hidden      bool
	Required    bool
	Default     DefaultValue
	AltSources  []AlternativeSource
}

type Choice struct {
	Value       string
	Description string
}

type DefaultValue struct {
	IsSet  bool
	Bool   bool
	String string
}

type AlternativeSource struct {
	Type                string
	EnvironmentVariable string
	File                FileSource
}

type FileSource struct {
	Format   string
	Path     string
	Property string
}

// Arguments returns true if any of the commands have arguments.
func (d Document) Arguments() bool {
	var argsHelper func(node *CommandTrieNode) bool
	argsHelper = func(node *CommandTrieNode) bool {
		if len(node.Command.Arguments) > 0 {
			return true
		}

		return slices.ContainsFunc(node.Commands, argsHelper)
	}

	return argsHelper(d.CommandTrie.Root)
}

// Flags returns true if any of the commands have flags.
func (d Document) Flags() bool {
	var flagsHelper func(node *CommandTrieNode) bool
	flagsHelper = func(node *CommandTrieNode) bool {
		if len(node.Command.Flags) > 0 {
			return true
		}

		return slices.ContainsFunc(node.Commands, flagsHelper)
	}

	return flagsHelper(d.CommandTrie.Root)
}

// VisibleFlags returns true if any of the commands have visible flags.
func (d Document) VisibleFlags() bool {
	var flagsHelper func(node *CommandTrieNode) bool
	flagsHelper = func(node *CommandTrieNode) bool {
		for _, flag := range node.Command.Flags {
			if !flag.Hidden {
				return true
			}
		}

		return slices.ContainsFunc(node.Commands, flagsHelper)
	}

	return flagsHelper(d.CommandTrie.Root)
}

// EnumeratedArgs returns true if any fixed arguments on any commands contain enumerated values.
func (d Document) FixedEnumeratedArgs() bool {
	var helper func(node *CommandTrieNode) bool
	helper = func(node *CommandTrieNode) bool {
		for _, arg := range node.Command.Arguments {
			if len(arg.Choices) > 0 && !arg.Variadic {
				return true
			}
		}

		return slices.ContainsFunc(node.Commands, helper)
	}

	return helper(d.CommandTrie.Root)
}

// EnumeratedArgs returns true if any variadic arguments on any commands contain enumerated values.
func (d Document) VariadicEnumeratedArgs() bool {
	var helper func(node *CommandTrieNode) bool
	helper = func(node *CommandTrieNode) bool {
		for _, arg := range node.Command.Arguments {
			if len(arg.Choices) > 0 && arg.Variadic {
				return true
			}
		}

		return slices.ContainsFunc(node.Commands, helper)
	}

	return helper(d.CommandTrie.Root)
}

// EnumeratedFlags returns true if any fixed type flags on any commands contain enumerated values.
func (d Document) FixedEnumeratedFlags() bool {
	var helper func(node *CommandTrieNode) bool
	helper = func(node *CommandTrieNode) bool {
		for _, flag := range node.Command.Flags {
			if len(flag.Choices) > 0 && !flag.Variadic {
				return true
			}
		}

		return slices.ContainsFunc(node.Commands, helper)
	}

	return helper(d.CommandTrie.Root)
}

// EnumeratedFlags returns true if any variadic type flags on any commands contain enumerated values.
func (d Document) VariadicEnumeratedFlags() bool {
	var helper func(node *CommandTrieNode) bool
	helper = func(node *CommandTrieNode) bool {
		for _, flag := range node.Command.Flags {
			if len(flag.Choices) > 0 && flag.Variadic {
				return true
			}
		}

		return slices.ContainsFunc(node.Commands, helper)
	}

	return helper(d.CommandTrie.Root)
}

func (d Document) InternalCliErrorCode() int {
	for _, ec := range d.Global.ExitCodes {
		if ec.Status == "INTERNAL_CLI_ERROR" {
			return ec.Code
		}
	}
	return 1
}

func (d Document) BadUserInputErrorCode() int {
	for _, ec := range d.Global.ExitCodes {
		if ec.Status == "BAD_USER_INPUT_ERROR" {
			return ec.Code
		}
	}
	return 2
}

func (d Document) AltSources() bool {
	var helper func(node *CommandTrieNode) bool
	helper = func(node *CommandTrieNode) bool {
		for _, flag := range node.Command.Flags {
			if len(flag.AltSources) > 0 {
				return true
			}
		}

		return slices.ContainsFunc(node.Commands, helper)
	}

	return helper(d.CommandTrie.Root)
}

// NonHiddenFlags returns true if there are any flags for the given command where Hidden isfalse.
func (cmd Command) VisibleFlags() bool {
	visible := false

	for _, flag := range cmd.Flags {
		if !flag.Hidden {
			visible = true
			break
		}
	}

	return visible
}

// EnumeratedArgs returns true if any fixed arguments on the command contain enumerated values.
func (cmd Command) FixedEnumeratedArgs() bool {
	for _, arg := range cmd.Arguments {
		if len(arg.Choices) > 0 && !arg.Variadic {
			return true
		}
	}

	return false
}

// EnumeratedArgs returns true if any variadic arguments on the command contain enumerated values.
func (cmd Command) VariadicEnumeratedArgs() bool {
	for _, arg := range cmd.Arguments {
		if len(arg.Choices) > 0 && arg.Variadic {
			return true
		}
	}

	return false
}

// EnumeratedFlags returns true if any fixed type flags on the command contain enumerated values.
func (cmd Command) FixedEnumeratedFlags() bool {
	for _, flag := range cmd.Flags {
		if len(flag.Choices) > 0 && !flag.Variadic {
			return true
		}
	}

	return false
}

// EnumeratedFlags returns true if any variadic type flags on the command contain enumerated values.
func (cmd Command) VariadicEnumeratedFlags() bool {
	for _, flag := range cmd.Flags {
		if len(flag.Choices) > 0 && flag.Variadic {
			return true
		}
	}

	return false
}

func (cmd Command) InternalCliErrorCode() int {
	for _, ec := range cmd.ExitCodes {
		if ec.Status == "INTERNAL_CLI_ERROR" {
			return ec.Code
		}
	}
	return 1
}

func (cmd Command) BadUserInputErrorCode() int {
	for _, ec := range cmd.ExitCodes {
		if ec.Status == "BAD_USER_INPUT_ERROR" {
			return ec.Code
		}
	}
	return 2
}

func (cmd Command) UnauthenticatedErrorCode() int {
	for _, ec := range cmd.ExitCodes {
		if ec.Status == "UNAUTHENTICATED_ERROR" {
			return ec.Code
		}
	}
	return 3
}

func (cmd Command) UnauthorizedErrorCode() int {
	for _, ec := range cmd.ExitCodes {
		if ec.Status == "UNAUTHORIZED_ERROR" {
			return ec.Code
		}
	}
	return 4
}

func (cmd Command) CanceledErrorCode() int {
	for _, ec := range cmd.ExitCodes {
		if ec.Status == "CANCELED_ERROR" {
			return ec.Code
		}
	}
	return 5
}

func (cmd Command) NotImplementedCode() int {
	for _, ec := range cmd.ExitCodes {
		if ec.Status == "NOT_IMPLEMENTED_ERROR" {
			return ec.Code
		}
	}
	return 6
}

// UnmarshalYAML ummarshalls the given YAML file into an Document domain object.
func UnmarshalYAML(path string) (Document, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	// validate the document
	err = ValidateDocumentYAML(contents)
	if err != nil {
		return Document{}, err
	}
	var doc oclifile.OpenCliDocument
	err = yaml.Unmarshal(contents, &doc)
	if err != nil {
		return Document{}, err
	}
	// return the domain-oriented struct
	return docFromUnmarshalled(doc)
}

// UnmarshalJSON ummarshalls the given JSON file into an Document domain object.
func UnmarshalJSON(path string) (Document, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return Document{}, err
	}
	// validate the document
	err = ValidateDocumentJSON(contents)
	if err != nil {
		return Document{}, err
	}
	var doc oclifile.OpenCliDocument
	err = json.Unmarshal(contents, &doc)
	if err != nil {
		return Document{}, err
	}
	// return the domain-oriented struct
	return docFromUnmarshalled(doc)
}

/* Command Trie Data Structure
------------------------------------------------------------------------------------------------- */

var wsRE = regexp.MustCompile(`[^\S\r\n]+`)

// CommandTrie is a data structure similar to a [Trie](https://en.wikipedia.org/wiki/Trie).
// Instead of each node containing letters in a string, it contains a command in the command line.
// This allows us to ensure that grouping commands that may not otherwise be defined in the OpenCLI Doc are accounted for.
// Grouping commands are commands that are not executable and exist only as an internal node in the trie (not a leaf).
type CommandTrie struct {
	doc  oclifile.OpenCliDocument
	Root *CommandTrieNode
}

// CommandTrieNode represents a hierarchical view of the CLI command structure.
type CommandTrieNode struct {
	Name     string
	Command  Command
	Commands []*CommandTrieNode
}

func (t *CommandTrie) Insert(cmdLine string, cmd oclifile.Command) {
	cmdLineNoParams := regexp.MustCompile(`[^\S\r\n][^A-Za-z]`).Split(cmdLine, -1)[0]
	argsRE := regexp.MustCompile(`<.*>`)
	flagsRE := regexp.MustCompile(`\[.*\]`)
	cmds := wsRE.Split(cmdLineNoParams, -1)

	if t.Root == nil {
		rootParams := []string{}
		if len(cmds) > 0 {
			rootParams = append(rootParams, "{command}")
		}
		if argsRE.MatchString(cmdLine) {
			rootParams = append(rootParams, "<arguments>")
		}
		if flagsRE.MatchString(cmdLine) {
			rootParams = append(rootParams, "[flags]")
		}

		t.Root = &CommandTrieNode{
			Name: cmds[0],
			Command: Command{
				Line:        strings.Join(append([]string{cmds[0]}, rootParams...), " "),
				Name:        cmds[0],
				LeafName:    cmds[0],
				Params:      strings.Join(rootParams, " "),
				Summary:     t.doc.Info.Summary,
				Description: t.doc.Info.Description,
				Group:       true,
			},
		}
	}

	nameSegments := []string{cmds[0]}

	node := t.Root
	for i := 1; i < len(cmds); i++ {
		cmdIndex := node.indexOfSubcommand(cmds[i])

		if cmdIndex < 0 {
			newNode := &CommandTrieNode{
				Name: cmds[i],
			}
			// add derived data; this may be overwritten by command definitions in the file
			cmdParams := []string{}
			isGroup := false
			if i < len(cmds)-1 { // if there are other commands left in the line, we're at an internal node in the trie
				cmdParams = append(cmdParams, "{command}")
				isGroup = true
			}
			if argsRE.MatchString(cmdLine) {
				cmdParams = append(cmdParams, "<arguments>")
			}
			if flagsRE.MatchString(cmdLine) {
				cmdParams = append(cmdParams, "[flags]")
			}
			nameSegments = append(nameSegments, cmds[i])
			newNode.Command = Command{
				Line:     strings.Join(append(nameSegments, cmdParams...), " "),
				Name:     strings.Join(nameSegments, " "),
				LeafName: cmds[i],
				Params:   strings.Join(cmdParams, " "),
				Group:    isGroup,
			}
			node.Commands = append(node.Commands, newNode)
			cmdIndex = len(node.Commands) - 1
		}

		node = node.Commands[cmdIndex]
	}

	node.Command = translateCommand(cmdLine, cmd)
}

/* Private receiver methods
------------------------------------------------------------------------------------------------- */

func (n *CommandTrieNode) indexOfSubcommand(name string) int {
	for i, cmd := range n.Commands {
		if cmd.Name == name {
			return i
		}
	}
	return -1
}

/* Private helper functions
------------------------------------------------------------------------------------------------- */

// fromUnmarshalled translates the raw unmarshalled struct to the domain struct
func docFromUnmarshalled(doc oclifile.OpenCliDocument) (Document, error) {
	// First ensure all command lines are 'absolute paths' that start with the binary
	err := validateCommandLines(doc)
	if err != nil {
		return Document{}, err
	}

	domainDoc := Document{
		Info: Info{
			Binary: doc.Info.Binary,
			Contact: Contact{
				Email: doc.Info.Contact.Email,
				Name:  doc.Info.Contact.Name,
				URL:   doc.Info.Contact.URL,
			},
			Description: doc.Info.Description,
			License: License{
				Name:   doc.Info.License.Name,
				SpdxID: doc.Info.License.SpdxID,
				URL:    doc.Info.License.URL,
			},
			Summary: doc.Info.Summary,
			Title:   doc.Info.Title,
			Version: doc.Info.Version,
		},
		OpenCLIVersion: doc.OpenCLIVersion,
	}
	// Add install methods
	for _, install := range doc.Install {
		domainDoc.Install = append(domainDoc.Install, Install{
			Command:     install.Command,
			Description: install.Description,
			Name:        install.Name,
			URL:         install.URL,
		})
	}
	// Add globals, e.g. flags, exit codes, etc.
	domainDoc.Global = translateGlobal(doc)
	// Build hierarchical CommandTrie
	trie, err := buildCommandTrie(doc)
	if err != nil {
		return Document{}, err
	}

	domainDoc.CommandTrie = trie

	return domainDoc, nil
}

func translateGlobal(doc oclifile.OpenCliDocument) Global {
	exitCodes := []ExitCode{}
	// Add global exit codes
	for _, exitCode := range doc.Global.ExitCodes {
		exitCodes = append(exitCodes, ExitCode{
			Code:        exitCode.Code,
			Status:      exitCode.Status,
			Description: exitCode.Description,
			Summary:     exitCode.Summary,
		})
	}

	sort.Slice(exitCodes, func(i, j int) bool {
		return exitCodes[i].Code <= exitCodes[j].Code
	})

	var globalFlags []Flag
	for _, flag := range doc.Global.Flags {
		globalFlags = append(globalFlags, translateFlag(flag))
	}

	return Global{
		ExitCodes: exitCodes,
		Flags:     globalFlags,
	}
}

func validateCommandLines(doc oclifile.OpenCliDocument) error {
	binRE := regexp.MustCompile(`^` + doc.Info.Binary)

	for key := range doc.Commands {
		if !binRE.MatchString(key) {
			return fmt.Errorf("cmd `%s` must be an absolute path starting with the binary", key)
		}
	}

	return nil
}

func translateCommand(cmdLine string, cmd oclifile.Command) Command {
	name, leaf, params := parseCommandLine(cmdLine)
	// add arguments to command
	var args []Argument
	for _, arg := range cmd.Arguments {
		args = append(args, translateArgument(arg))
	}
	// add flags to command
	var flags []Flag
	for _, flag := range cmd.Flags {
		flags = append(flags, translateFlag(flag))
	}
	// add exit codes to command
	cmdSpecificExitCodes := translateCmdExitCodes(cmd)
	command := Command{
		Aliases:              cmd.Aliases,
		Group:                cmd.Group,
		Hidden:               cmd.Hidden,
		Line:                 cmdLine,
		Name:                 name,
		LeafName:             leaf,
		Params:               params,
		Arguments:            args,
		Flags:                flags,
		CmdSpecificExitCodes: cmdSpecificExitCodes,
	}
	// Set attributes IFF they are not zero values overwriting non-zero values
	if cmd.Description != "" {
		command.Description = cmd.Description
	}
	if cmd.Summary != "" {
		command.Summary = cmd.Summary
	}
	// return the translated command
	return command
}

// parseCommandLine returns three distinct segments of the command line
// name - the absolute command line path without the parameters
// leaf - the final command segment in the absolute command line path
// params - the `{command} <arguments> [flags]` part of the command line if present
func parseCommandLine(cmdLine string) (string, string, string) {
	lineParamsRE := regexp.MustCompile(`[^\S\r\n][^A-Za-z]`)
	name := cmdLine
	params := ""
	split := lineParamsRE.Split(name, -1)
	if len(split) > 0 {
		name = split[0]
	}
	if len(split) > 1 {
		params = split[1]
	}
	commandSegments := regexp.MustCompile(`[^\S\r\n]`).Split(name, -1)
	leaf := commandSegments[len(commandSegments)-1]

	return name, leaf, params
}

func translateArgument(arg oclifile.Argument) Argument {
	domainArg := Argument{
		Name:        arg.Name,
		Summary:     arg.Summary,
		Description: arg.Description,
		Type:        arg.Type,
		Variadic:    arg.Variadic,
		Required:    arg.Required,
		Default: DefaultValue{
			IsSet:  arg.Default.IsSet,
			Bool:   arg.Default.Bool,
			String: arg.Default.String,
		},
	}

	for _, choice := range arg.Choices {
		domainArg.Choices = append(domainArg.Choices, Choice{
			Value:       choice.Value,
			Description: choice.Description,
		})
	}

	return domainArg
}

func translateFlag(flag oclifile.Flag) Flag {
	domainFlag := Flag{
		Name:        flag.Name,
		Aliases:     flag.Aliases,
		Hint:        flag.Hint,
		Summary:     flag.Summary,
		Description: flag.Description,
		Type:        flag.Type,
		Variadic:    flag.Variadic,
		Required:    flag.Required,
		Hidden:      flag.Hidden,
		Default: DefaultValue{
			IsSet:  flag.Default.IsSet,
			Bool:   flag.Default.Bool,
			String: flag.Default.String,
		},
	}

	for _, choice := range flag.Choices {
		domainFlag.Choices = append(domainFlag.Choices, Choice{
			Value:       choice.Value,
			Description: choice.Description,
		})
	}

	for _, src := range flag.AltSources {
		domainFlag.AltSources = append(domainFlag.AltSources, AlternativeSource{
			Type:                src.Type,
			EnvironmentVariable: src.EnvironmentVariable,
			File: FileSource{
				Format:   src.File.Format,
				Path:     src.File.Path,
				Property: src.File.Property,
			},
		})
	}

	return domainFlag
}

func translateCmdExitCodes(cmdObj oclifile.Command) []ExitCode {
	var exitCodes []ExitCode
	statuses := map[string]struct{}{}
	for _, ec := range cmdObj.ExitCodes {
		statuses[ec.Status] = struct{}{}
		exitCodes = append(exitCodes, ExitCode{
			Code:        ec.Code,
			Status:      ec.Status,
			Summary:     ec.Summary,
			Description: ec.Description,
		})
	}
	// return sorted array of exit codes
	sort.Slice(exitCodes, func(i, j int) bool {
		return exitCodes[i].Code <= exitCodes[j].Code
	})
	return exitCodes
}

func mergeGlobalCmdExitCodes(doc Document, cmdExitCodes []ExitCode) []ExitCode {
	var exitCodes []ExitCode
	statuses := map[string]struct{}{}
	// add cmd-specific exit codes
	for _, ec := range cmdExitCodes {
		statuses[ec.Status] = struct{}{}
		exitCodes = append(exitCodes, ec)
	}
	// add global exit codes
	for _, ec := range doc.Global.ExitCodes {
		if _, ok := statuses[ec.Status]; !ok {
			exitCodes = append(exitCodes, ec)
		}
	}
	// return sorted array of exit codes
	sort.Slice(exitCodes, func(i, j int) bool {
		return exitCodes[i].Code <= exitCodes[j].Code
	})
	return exitCodes
}

func buildCommandTrie(doc oclifile.OpenCliDocument) (*CommandTrie, error) {
	trie := &CommandTrie{doc: doc}
	// Fist sort the command lines alphabetically so we have a deterministic order in our Trie
	var commandLines []string
	for cmdLine := range doc.Commands {
		commandLines = append(commandLines, cmdLine)
	}
	sort.Strings(commandLines)
	// Insert each command line into the trie
	for _, cmdLine := range commandLines {
		trie.Insert(cmdLine, doc.Commands[cmdLine])
	}
	// return the constructed trie
	return trie, nil
}
