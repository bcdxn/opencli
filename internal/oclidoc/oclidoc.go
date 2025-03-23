package oclidoc

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// oclidoc.go contains the OpenCLI domain types used when unmarshalling an OpenCLI Document.

// OpenCliDocument represents the OpenCLI document.
type OpenCliDocument struct {
	OpenCLIVersion string             `json:"opencliVersion" yaml:"opencliVersion"`
	Info           Info               `json:"info" yaml:"info"`
	Install        []Install          `json:"install" yaml:"install"`
	Commands       map[string]Command `json:"commands" yaml:"commands"`
}

// Info represents the metadata about the CLI described by the OpenCLI document.
type Info struct {
	Title       string  `json:"title" yaml:"title"`
	Summary     string  `json:"summary" yaml:"summary"`
	Description string  `json:"description" yaml:"description"`
	License     License `json:"license" yaml:"license"`
	Contact     Contact `json:"contact" yaml:"contact"`
	Binary      string  `json:"binary" yaml:"binary"`
	Version     string  `json:"version" yaml:"version"`
}

// License represents the license information for the CLI described by the OpenCLI document.
type License struct {
	Name   string `json:"name" yaml:"name"`
	SpdxID string `json:"spdxId" yaml:"spdxId"`
	URL    string `json:"url" yaml:"url"`
}

// Contact represents contact information for maintainers of the CLI described by the OpenCLI document.
type Contact struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"email" yaml:"email"`
	URL   string `json:"url" yaml:"url"`
}

// Install represents information about ways a user can install the CLI.
type Install struct {
	Name        string `json:"name" yaml:"name"`
	Command     string `json:"command" yaml:"command"`
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description" yaml:"description"`
}

// Command represents an OpenCLI command.
type Command struct {
	Summary     string     `json:"summary" yaml:"summary"`
	Description string     `json:"description" yaml:"description"`
	Aliases     []string   `json:"aliases" yaml:"aliases"`
	Arguments   []Argument `json:"arguments" yaml:"arguments"`
	Flags       []Flag     `json:"flags" yaml:"flags"`
	Hidden      bool       `json:"hidden" yaml:"hidden"`
	Group       bool       `json:"group" yaml:"group"`
}

// Argument represents an OpenCLI command argument.
type Argument struct {
	Name        string   `json:"name" yaml:"name"`
	Summary     string   `json:"summary" yaml:"summary"`
	Description string   `json:"description" yaml:"description"`
	Type        string   `json:"type" yaml:"type"`
	Variadic    Variadic `json:"variadic" yaml:"variadic"`
	Choices     []Choice `json:"choices" yaml:"choices"`
	Required    bool     `json:"required" yaml:"required"`
	Default     any      `json:"default" yaml:"default"`
}

// Flag represents an OpenCLI command flag.
type Flag struct {
	Name        string   `json:"name" yaml:"name"`
	Aliases     []string `json:"aliases" yaml:"aliases"`
	Hint        string   `json:"hint" yaml:"hint"`
	Summary     string   `json:"summary" yaml:"summary"`
	Description string   `json:"description" yaml:"description"`
	Type        string   `json:"type" yaml:"type"`
	Variadic    Variadic `json:"variadic" yaml:"variadic"`
	Choices     []Choice `json:"choices" yaml:"choices"`
	Hidden      bool     `json:"hidden" yaml:"hidden"`
	Required    bool     `json:"required" yaml:"required"`
	Default     any      `json:"default" yaml:"default"`
}

type Variadic struct {
	Enabled bool
	Sep     string
}

// UnmarshalJSON handles custom unmarshalling logic
// The 'variadic' property in the JSON may hold either a bool literal or a string literal, and it must be marshalled into a struct.
func (v *Variadic) UnmarshalJSON(bs []byte) error {
	// first try unmarshalling a plain bool representation
	var boolDeclaration bool
	err := json.Unmarshal(bs, &boolDeclaration)
	if err == nil {
		if boolDeclaration {
			v.Enabled = true
			v.Sep = "," // default separator is a comma
		}
		return nil
	}
	// bool unmarshal failed, try string representation
	var strDeclaration string
	err = json.Unmarshal(bs, &strDeclaration)
	if err == nil {
		if strDeclaration != "" {
			v.Enabled = true
			v.Sep = strDeclaration
		}
		return nil
	}

	return err
}

// UnmarshalYAML handles custom unmarshalling logic.
// The 'variadic' property in the YAML may hold either a bool literal or a string literal and it must be marshalled into a struct.
func (v *Variadic) UnmarshalYAML(node *yaml.Node) error {
	// first try unmarshalling a plain bool representation
	var boolDeclaration bool
	err := node.Decode(&boolDeclaration)
	// err := yaml.Unmarshal([]byte(node.Value), &boolDeclaration)
	if err == nil {
		if boolDeclaration {
			v.Enabled = true
			v.Sep = "," // default separator is a comma
		}
		return nil
	}
	// bool unmarshal failed, try string representation
	var strDeclaration string
	err = node.Decode(&strDeclaration)
	// err = yaml.Unmarshal([]byte(node.Value), &strDeclaration)
	if err == nil {
		if strDeclaration != "" {
			v.Enabled = true
			v.Sep = strDeclaration
		}
		return nil
	}
	return err
}

// Choice represents an enumeration of an argument/flag
type Choice struct {
	Value       string `json:"value" yaml:"value"`
	Description string `json:"description" yaml:"description"`
}
