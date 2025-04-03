package oclifile

import (
	"encoding/json"
	"errors"

	"gopkg.in/yaml.v3"
)

// oclidoc.go contains the OpenCLI domain types used when unmarshalling an OpenCLI Document file.

// OpenCliDocument represents the OpenCLI document.
type OpenCliDocument struct {
	OpenCLIVersion string             `json:"opencliVersion" yaml:"opencliVersion"`
	Info           Info               `json:"info" yaml:"info"`
	Install        []Install          `json:"install" yaml:"install"`
	Global         Global             `json:"global" yaml:"global"`
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

type Global struct {
	ExitCodes []ExitCode `json:"exitCodes" yaml:"exitCodes"`
	Flags     []Flag     `json:"flags" yaml:"flags"`
}

// ExitCode represents a possible exit code from the CLI.
type ExitCode struct {
	Code        int    `json:"code" yaml:"code"`
	Status      string `json:"status" yaml:"status"`
	Summary     string `json:"summary" yaml:"summary"`
	Description string `json:"description" yaml:"description"`
}

// Command represents an OpenCLI command.
type Command struct {
	Line        string     // Not unmarshalled, but will be added during translation process
	Summary     string     `json:"summary" yaml:"summary"`
	Description string     `json:"description" yaml:"description"`
	Aliases     []string   `json:"aliases" yaml:"aliases"`
	Arguments   []Argument `json:"arguments" yaml:"arguments"`
	Flags       []Flag     `json:"flags" yaml:"flags"`
	Hidden      bool       `json:"hidden" yaml:"hidden"`
	Group       bool       `json:"group" yaml:"group"`
	ExitCodes   []ExitCode `json:"exitCodes" yaml:"exitCodes"`
}

// Argument represents an OpenCLI command argument.
type Argument struct {
	Name        string       `json:"name" yaml:"name"`
	Summary     string       `json:"summary" yaml:"summary"`
	Description string       `json:"description" yaml:"description"`
	Type        string       `json:"type" yaml:"type"`
	Variadic    bool         `json:"variadic" yaml:"variadic"`
	Choices     []Choice     `json:"choices" yaml:"choices"`
	Required    bool         `json:"required" yaml:"required"`
	Default     DefaultValue `json:"default" yaml:"default"`
}

// Flag represents an OpenCLI command flag.
type Flag struct {
	Name        string             `json:"name" yaml:"name"`
	Aliases     []string           `json:"aliases" yaml:"aliases"`
	Hint        string             `json:"hint" yaml:"hint"`
	Summary     string             `json:"summary" yaml:"summary"`
	Description string             `json:"description" yaml:"description"`
	Type        string             `json:"type" yaml:"type"`
	Variadic    bool               `json:"variadic" yaml:"variadic"`
	Choices     []Choice           `json:"choices" yaml:"choices"`
	Hidden      bool               `json:"hidden" yaml:"hidden"`
	Required    bool               `json:"required" yaml:"required"`
	Default     DefaultValue       `json:"default" yaml:"default"`
	AltSources  AlternativeSources `json:"alternativeSources" yaml:"alternativeSources"`
}

// Choice represents an enumeration of an argument/flag
type Choice struct {
	Value       string `json:"value" yaml:"value"`
	Description string `json:"description" yaml:"description"`
}

type DefaultValue struct {
	IsSet  bool
	Bool   bool
	String string
}

func (v *DefaultValue) UnmarshalJSON(bs []byte) error {
	// first try unmarshalling a bool value
	var b bool
	err := json.Unmarshal(bs, &b)
	if err == nil {
		v.Bool = b
		v.IsSet = true
		return nil
	}
	// next try unmarshalling a string value
	var s string
	err = json.Unmarshal(bs, &s)
	if err == nil {
		v.String = s
		v.IsSet = true
		return nil
	}

	// The value was neither a string nor a bool and is therefore not allowed
	return errors.New("expected bool or string but found neither")
}

func (v *DefaultValue) UnmarshalYAML(node *yaml.Node) error {
	// first try unmarshalling a bool value
	var b bool
	err := yaml.Unmarshal([]byte(node.Value), &b)
	if err == nil {
		v.Bool = b
		v.IsSet = true
		return nil
	}
	// next try unmarshalling a string value
	var s string
	err = yaml.Unmarshal([]byte(node.Value), &s)
	if err == nil {
		v.String = s
		v.IsSet = true
		return nil
	}
	// The value was neither a string nor a bool and is therefore not allowed
	return errors.New("expected bool or string but found neither")
}

type AlternativeSource struct {
	Type                string
	EnvironmentVariable string
	File                FileSource
}

type AlternativeSources []AlternativeSource

type rawMessageYAML struct {
	unmarshal func(interface{}) error
}

func (msg *rawMessageYAML) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.unmarshal = unmarshal
	return nil
}

func (msg *rawMessageYAML) Unmarshal(v interface{}) error {
	res := msg.unmarshal(v)
	return res
}

func (s *AlternativeSources) UnmarshalJSON(bs []byte) error {
	var sources []json.RawMessage
	err := json.Unmarshal(bs, &sources)

	if err != nil {
		return err
	}

	for _, source := range sources {
		var altSrc AlternativeSource
		err := json.Unmarshal(source, &altSrc)
		if err != nil {
			return err
		} else {
			*s = append(*s, altSrc)
		}
	}

	return nil
}

func (s *AlternativeSources) UnmarshalYAML(node *yaml.Node) error {
	var sources []rawMessageYAML
	err := node.Decode(&sources)
	if err != nil {
		return err
	}

	for _, source := range sources {
		var altSrc AlternativeSource
		err := source.Unmarshal(&altSrc)
		// // err := yaml.Unmarshal(source, &altSrc)
		if err != nil {
			return err
		} else {
			*s = append(*s, altSrc)
		}
	}

	return nil
}

func (s *AlternativeSource) UnmarshalJSON(bs []byte) error {
	// first try unmarshalling an environment variable source
	var env EnvironmentVariableSource
	err := json.Unmarshal(bs, &env)
	if err != nil {
		s.EnvironmentVariable = env.Name
		s.Type = env.Type
		return nil
	}

	// next try unmarshalling a file source
	var file FileSource
	err = json.Unmarshal(bs, &file)
	if err != nil {
		s.File = file
		s.Type = file.Type
		return nil
	}

	return errors.New("expected EnvironmentVariable or File alternative source but found neither")
}

func (s *AlternativeSource) UnmarshalYAML(node *yaml.Node) error {
	// first try unmarshalling an environment variable source
	var env EnvironmentVariableSource
	err := node.Decode(&env)
	if err == nil && env.Type == "env" {
		s.EnvironmentVariable = env.Name
		s.Type = env.Type
		return nil
	}

	// next try unmarshalling a file source
	var file FileSource
	err = node.Decode(&file)
	if err == nil && file.Type == "file" {
		s.File = file
		s.Type = file.Type
		return nil
	}

	return errors.New("expected EnvironmentVariable or File alternative source but found neither")
}

type EnvironmentVariableSource struct {
	Type string `json:"type" yaml:"type"`
	Name string `json:"name" yaml:"name"`
}

type FileSource struct {
	Type     string `json:"type" yaml:"type"`
	Format   string `json:"format" yaml:"format"`
	Path     string `json:"path" yaml:"path"`
	Property string `json:"property" yaml:"property"`
}
