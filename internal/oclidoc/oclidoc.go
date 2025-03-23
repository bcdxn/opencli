package oclidoc

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
	Variadic    bool     `json:"variadic" yaml:"variadic"`
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
	Variadic    bool     `json:"variadic" yaml:"variadic"`
	Choices     []Choice `json:"choices" yaml:"choices"`
	Hidden      bool     `json:"hidden" yaml:"hidden"`
	Required    bool     `json:"required" yaml:"required"`
	Default     any      `json:"default" yaml:"default"`
}

// Choice represents an enumeration of an argument/flag
type Choice struct {
	Value       string `json:"value" yaml:"value"`
	Description string `json:"description" yaml:"description"`
}
