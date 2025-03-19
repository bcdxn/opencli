package ocli

import (
	"encoding/json"
	"os"
	"regexp"
	"sort"

	"github.com/bcdxn/opencli/internal/oclidoc"
	"gopkg.in/yaml.v3"
)

// ocli.go contains the OpenCLI domain types.

// OpenCliDocument represents the OpenCLI document.
type OpenCliDocument struct {
	OpenCLIVersion string
	Info           Info
	Install        []Install
	Commands       []Command
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

// Command represents n OpenCLI command.
type Command struct {
	Line        string // The full command line as defined in the OpenCLI Spec Document
	Name        string // The command part of the command line
	Params      string // The args/flags part of the command line
	Summary     string
	Description string
	Arguments   []Argument
	Flags       []Flag
	Hidden      bool
	Group       bool
}

// Argument represents an OpenCLI command argument.
type Argument struct {
	Name        string
	Summary     string
	Description string
	Type        string
	Kind        string
	Choices     []Choice
	Required    bool
}

// Flag represents an OpenCLI command flag.
type Flag struct {
	Name        string
	Alias       string
	Summary     string
	Description string
	Type        string
	Kind        string
	Choices     []Choice
	Hidden      bool
	Required    bool
}

type Choice struct {
	Value       string
	Description string
}

// UnmarshalYAML ummarshalls the given YAML file into an OpenCliDocument domain object.
func UnmarshalYAML(path string) (OpenCliDocument, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return OpenCliDocument{}, err
	}
	// validate the document
	err = ValidateDocumentYAML(contents)
	if err != nil {
		return OpenCliDocument{}, err
	}
	var doc oclidoc.OpenCliDocument
	err = yaml.Unmarshal(contents, &doc)
	if err != nil {
		return OpenCliDocument{}, err
	}
	// return the domain-oriented struct
	return docFromUnmarshalled(doc)
}

// UnmarshalJSON ummarshalls the given JSON file into an OpenCliDocument domain object.
func UnmarshalJSON(path string) (OpenCliDocument, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return OpenCliDocument{}, err
	}
	// validate the document
	err = ValidateDocumentJSON(contents)
	if err != nil {
		return OpenCliDocument{}, err
	}
	var doc oclidoc.OpenCliDocument
	err = json.Unmarshal(contents, &doc)
	if err != nil {
		return OpenCliDocument{}, err
	}
	// return the domain-oriented struct
	return docFromUnmarshalled(doc)
}

/* Private helper fucntions
------------------------------------------------------------------------------------------------- */

// fromUnmarshalled translates the raw unmarshalled struct to the domain struct
func docFromUnmarshalled(doc oclidoc.OpenCliDocument) (OpenCliDocument, error) {
	domainDoc := OpenCliDocument{
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
	// Add commands
	for cmd, cmdObj := range doc.Commands {
		domainDoc.Commands = append(domainDoc.Commands, translateCommand(doc, cmd, cmdObj))
	}

	// Sort command by `Line` property to ensure a stable order
	sort.Slice(domainDoc.Commands, func(i, j int) bool {
		return domainDoc.Commands[i].Name < domainDoc.Commands[j].Name
	})

	return domainDoc, nil
}

func translateCommand(doc oclidoc.OpenCliDocument, cmd string, cmdObj oclidoc.Command) Command {
	binRE := regexp.MustCompile(`^([^\S\r\n]*` + doc.Info.Binary + `[^\S\r\n]+)?`)
	lineParamsRE := regexp.MustCompile(`[^\S\r\n][^A-Za-z]`)
	line := binRE.ReplaceAllString(cmd, doc.Info.Binary+" ")
	name := binRE.ReplaceAllString(cmd, "")
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
	// return the translated command
	return Command{
		Description: cmdObj.Description,
		Group:       cmdObj.Group,
		Hidden:      cmdObj.Hidden,
		Line:        line,
		Name:        name,
		Summary:     cmdObj.Summary,
		Arguments:   args,
		Flags:       flags,
	}
}

// NonHiddenCommands returns true if there are any commands where Hidden is false.
func (d OpenCliDocument) VisibleCommands() bool {
	visible := false

	for _, cmd := range d.Commands {
		if !cmd.Hidden {
			visible = true
			break
		}
	}

	return visible
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

func translateArgument(arg oclidoc.Argument) Argument {
	domainArg := Argument{
		Name:        arg.Name,
		Summary:     arg.Summary,
		Description: arg.Description,
		Type:        arg.Type,
		Kind:        arg.Kind,
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
		Summary:     flag.Summary,
		Description: flag.Description,
		Type:        flag.Type,
		Kind:        flag.Kind,
		Required:    flag.Required,
		Hidden:      flag.Hidden,
	}

	for _, choice := range flag.Choices {
		domainFlag.Choices = append(domainFlag.Choices, Choice{
			Value:       choice.Value,
			Description: choice.Description,
		})
	}

	return domainFlag
}
