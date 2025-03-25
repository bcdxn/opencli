package oclispec

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
	"sort"
	"text/template"

	"github.com/bcdxn/opencli/internal/oclidoc"
	"gopkg.in/yaml.v3"
)

// ocli.go contains the OpenCLI domain types.

// Document represents the OpenCLI document.
type Document struct {
	OpenCLIVersion string
	Info           Info
	Install        []Install
	Global         Global
	Commands       []Command
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
	Default     any
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
}

type Choice struct {
	Value       string
	Description string
}

type DefaultValue struct {
	Bool   bool
	String string
}

// NonHiddenCommands returns true if there are any commands where Hidden is false.
func (d Document) VisibleCommands() bool {
	for _, cmd := range d.Commands {
		if !cmd.Hidden {
			return true
		}
	}

	return false
}

// Arguments returns true if any of the commands have arguments.
func (d Document) Arguments() bool {
	for _, cmd := range d.Commands {
		if len(cmd.Arguments) > 0 {
			return true
		}
	}

	return false
}

// Flags returns true if any of the commands have flags.
func (d Document) Flags() bool {
	for _, cmd := range d.Commands {
		if len(cmd.Flags) > 0 {
			return true
		}
	}

	return false
}

// VisibleFlags returns true if any of the commands have visible flags.
func (d Document) VisibleFlags() bool {
	for _, cmd := range d.Commands {
		for _, flag := range cmd.Flags {
			if !flag.Hidden {
				return true
			}
		}
	}

	return false
}

// EnumeratedArgs returns true if any fixed arguments on any commands contain enumerated values.
func (d Document) FixedEnumeratedArgs() bool {
	for _, cmd := range d.Commands {
		for _, arg := range cmd.Arguments {
			if len(arg.Choices) > 0 && !arg.Variadic {
				return true
			}
		}
	}

	return false
}

// EnumeratedArgs returns true if any variadic arguments on any commands contain enumerated values.
func (d Document) VariadicEnumeratedArgs() bool {
	for _, cmd := range d.Commands {
		for _, arg := range cmd.Arguments {
			if len(arg.Choices) > 0 && arg.Variadic {
				return true
			}
		}
	}

	return false
}

// EnumeratedFlags returns true if any fixed type flags on any commands contain enumerated values.
func (d Document) FixedEnumeratedFlags() bool {
	for _, cmd := range d.Commands {
		for _, flag := range cmd.Flags {
			if len(flag.Choices) > 0 && !flag.Variadic {
				return true
			}
		}
	}

	return false
}

// EnumeratedFlags returns true if any variadic type flags on any commands contain enumerated values.
func (d Document) VariadicEnumeratedFlags() bool {
	for _, cmd := range d.Commands {
		for _, flag := range cmd.Flags {
			if len(flag.Choices) > 0 && flag.Variadic {
				return true
			}
		}
	}

	return false
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
	var doc oclidoc.OpenCliDocument
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
	var doc oclidoc.OpenCliDocument
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

type CommandTrie struct {
	Root *CommandTrieNode
}

// CommandTrieNode represents a hierarchical view of the CLI command structure.
type CommandTrieNode struct {
	Name     string
	Command  Command
	Commands []*CommandTrieNode
}

func (t *CommandTrie) Insert(command Command) {
	cmds := wsRE.Split(command.Name, -1)

	if t.Root == nil {
		t.Root = &CommandTrieNode{
			Name: cmds[0],
			Command: Command{
				Name:  cmds[0],
				Group: true,
			},
		}
	}

	node := t.Root
	for i := 1; i < len(cmds); i++ {
		cmdIndex := node.indexOfSubcommand(cmds[i])

		if cmdIndex < 0 {
			node.Commands = append(node.Commands, &CommandTrieNode{
				Name: cmds[i],
			})
			cmdIndex = len(node.Commands) - 1
		}

		node = node.Commands[cmdIndex]
	}

	node.Command = command
}

/* Private receiver methods
------------------------------------------------------------------------------------------------- */

// rootCommandLine uses a template to render the root-level command usage line.
// e.g.: `ocli {command} <arguments> [flags]`.`
func (d Document) rootCommandLine() (string, error) {
	rootLine := bytes.NewBuffer([]byte{})
	rootTmpl := template.Must(template.New("root_line.tmpl").Parse(`{{.Info.Binary}}{{if .VisibleCommands}} {command}{{end}}{{if .Arguments}} <arguments>{{end}}{{if .VisibleFlags}} [flags]{{end}}`))

	err := rootTmpl.Execute(rootLine, d)
	if err != nil {
		return "", nil
	}

	return rootLine.String(), nil
}

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
func docFromUnmarshalled(doc oclidoc.OpenCliDocument) (Document, error) {
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
	domainDoc.Global = translateGlobal(doc)
	// Add commands
	for cmd, cmdObj := range doc.Commands {
		domainDoc.Commands = append(domainDoc.Commands, translateCommand(doc, domainDoc, cmd, cmdObj))
	}
	// Sort command by `Line` property to ensure a stable order
	sort.Slice(domainDoc.Commands, func(i, j int) bool {
		return domainDoc.Commands[i].Name < domainDoc.Commands[j].Name
	})
	// Build hierarchical CommandTrie
	domainDoc.buildCommandTrie()

	return domainDoc, nil
}

func translateGlobal(doc oclidoc.OpenCliDocument) Global {
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

	return Global{
		ExitCodes: exitCodes,
	}
}

func translateCommand(doc oclidoc.OpenCliDocument, domainDoc Document, cmd string, cmdObj oclidoc.Command) Command {
	binRE := regexp.MustCompile(`^([^\S\r\n]*` + doc.Info.Binary + `[^\S\r\n]+)?`)
	lineParamsRE := regexp.MustCompile(`[^\S\r\n][^A-Za-z]`)
	line := binRE.ReplaceAllString(cmd, doc.Info.Binary+" ")
	name := line
	name = lineParamsRE.Split(name, -1)[0]
	// add arguments to command
	var args []Argument
	for _, arg := range cmdObj.Arguments {
		args = append(args, translateArgument(arg))
	}
	// add flags to command
	var flags []Flag
	for _, flag := range cmdObj.Flags {
		flags = append(flags, translateFlag(flag))
	}
	// add exit codes to command
	cmdSpecificExitCodes := translateCmdExitCodes(cmdObj)
	mergedExitCodes := mergeGlobalCmdExitCodes(domainDoc, cmdSpecificExitCodes)
	// return the translated command
	return Command{
		Aliases:              cmdObj.Aliases,
		Description:          cmdObj.Description,
		Group:                cmdObj.Group,
		Hidden:               cmdObj.Hidden,
		Line:                 line,
		Name:                 name,
		Summary:              cmdObj.Summary,
		Arguments:            args,
		Flags:                flags,
		CmdSpecificExitCodes: cmdSpecificExitCodes,
		ExitCodes:            mergedExitCodes,
	}
}

func translateArgument(arg oclidoc.Argument) Argument {
	domainArg := Argument{
		Name:        arg.Name,
		Summary:     arg.Summary,
		Description: arg.Description,
		Type:        arg.Type,
		Variadic:    arg.Variadic,
		Required:    arg.Required,
	}

	for _, choice := range arg.Choices {
		domainArg.Choices = append(domainArg.Choices, Choice{
			Value:       choice.Value,
			Description: choice.Description,
		})
	}

	return domainArg
}

func translateFlag(flag oclidoc.Flag) Flag {
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

	return domainFlag
}

func translateCmdExitCodes(cmdObj oclidoc.Command) []ExitCode {
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

func (d *Document) buildCommandTrie() error {
	rootCmdLine, err := d.rootCommandLine()
	if err != nil {
		return err
	}

	trie := &CommandTrie{}

	trie.Insert(Command{
		Name:        d.Info.Binary,
		Line:        rootCmdLine,
		Summary:     d.Info.Summary,
		Description: d.Info.Description,
		Group:       true,
	})

	for _, cmd := range d.Commands {
		trie.Insert(cmd)
	}

	d.CommandTrie = trie

	return nil
}
