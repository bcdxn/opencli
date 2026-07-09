package validate

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/spec"
	"github.com/goccy/go-yaml"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// ValidationError represents a logical validation error in the spec
type ValidationError struct {
	Message string
	Path    string
}

func (e *ValidationError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s (at %s)", e.Message, e.Path)
	}
	return e.Message
}

// schemaValidationError wraps a *jsonschema.ValidationError to produce clean,
// human-readable output without exposing internal schema file paths.
type schemaValidationError struct {
	causes []*jsonschema.ValidationError
}

var schemaPrinter = message.NewPrinter(language.English)

func (e *schemaValidationError) Error() string {
	var sb strings.Builder
	for i, cause := range e.causes {
		if i > 0 {
			sb.WriteByte('\n')
		}
		writeSchemaError(&sb, cause, 0)
	}
	return sb.String()
}

// writeSchemaError recursively formats a ValidationError, skipping Group
// wrappers that only add noise.
func writeSchemaError(sb *strings.Builder, e *jsonschema.ValidationError, depth int) {
	// Skip group-only wrappers and descend directly into their causes.
	if _, ok := e.ErrorKind.(*kind.Group); ok {
		for i, cause := range e.Causes {
			if depth > 0 || i > 0 {
				sb.WriteByte('\n')
			}
			writeSchemaError(sb, cause, depth)
		}
		return
	}

	if depth > 0 {
		sb.WriteByte('\n')
		for i := 0; i < depth; i++ {
			sb.WriteString("  ")
		}
	}
	sb.WriteString("- ")
	if len(e.InstanceLocation) > 0 {
		fmt.Fprintf(sb, "at '/%s': ", strings.Join(e.InstanceLocation, "/"))
	}
	sb.WriteString(e.ErrorKind.LocalizedString(schemaPrinter))

	for _, cause := range e.Causes {
		writeSchemaError(sb, cause, depth+1)
	}
}

// wrapSchemaError converts a *jsonschema.ValidationError into the cleaner
// schemaValidationError type. Other error types are returned unchanged.
func wrapSchemaError(err error) error {
	verr, ok := err.(*jsonschema.ValidationError)
	if !ok {
		return err
	}
	return &schemaValidationError{causes: verr.Causes}
}

//go:generate mkdir -p out
//go:generate cp ../spec.schema.json ./out/spec.schema.json
//go:embed out/spec.schema.json
var schemaBytes []byte

var (
	compilerOnce sync.Once
	compilerErr  error
	schema       *jsonschema.Schema
)

func ensureCompiler() error {
	compilerOnce.Do(func() {
		schemaContents, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaBytes))
		if err != nil {
			compilerErr = err
			return
		}

		c := jsonschema.NewCompiler()
		if err := c.AddResource("spec.schema.json", schemaContents); err != nil {
			compilerErr = err
			return
		}
		s, err := c.Compile("spec.schema.json")
		if err != nil {
			compilerErr = err
			return
		}
		schema = s
	})
	return compilerErr
}

func ValidateJSON(data []byte) error {
	if err := ensureCompiler(); err != nil {
		return err
	}
	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("error unmarshalling document: %w", err)
	}
	// validate using schema
	if err := schema.Validate(doc); err != nil {
		return wrapSchemaError(err)
	}

	// Unmarshal into typed spec.Document and run logical validations
	specDoc, err := codec.UnmarshalJSON(data)
	if err != nil {
		return fmt.Errorf("error parsing document: %w", err)
	}

	return runLogicalValidations(specDoc)
}

func ValidateYAML(data []byte) error {
	if err := ensureCompiler(); err != nil {
		return err
	}
	var doc any
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}
	// validate using schema
	if err := schema.Validate(doc); err != nil {
		return wrapSchemaError(err)
	}

	// Unmarshal into typed spec.Document and run logical validations
	specDoc, err := codec.UnmarshalYAML(data)
	if err != nil {
		return fmt.Errorf("error parsing document: %w", err)
	}

	return runLogicalValidations(specDoc)
}

// runLogicalValidations performs logical validations beyond schema compliance
func runLogicalValidations(doc *spec.Document) error {
	// Extract defined config files from global section
	definedConfigFiles := make(map[string]bool)
	if doc.Global != nil {
		if doc.Global.Config.JSON != "" {
			definedConfigFiles["json"] = true
		}
		if doc.Global.Config.TOML != "" {
			definedConfigFiles["toml"] = true
		}
		if doc.Global.Config.YAML != "" {
			definedConfigFiles["yaml"] = true
		}
	}

	// Validate all commands recursively
	if doc.Commands != nil {
		if err := validateCommandTree(doc.Commands, definedConfigFiles); err != nil {
			return err
		}
	}

	return nil
}

// validateCommandTree recursively validates a command and its subcommands
func validateCommandTree(cmd *spec.CommandItem, definedConfigFiles map[string]bool) error {
	// Validate this command
	if err := validateCommand(cmd, definedConfigFiles); err != nil {
		return err
	}

	// Validate all subcommands
	if cmd.Commands != nil {
		for _, subcmd := range cmd.Commands {
			if err := validateCommandTree(subcmd, definedConfigFiles); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateCommand performs validations on a single command
func validateCommand(cmd *spec.CommandItem, definedConfigFiles map[string]bool) error {
	// Validate positional arguments ordering
	if len(cmd.Args) > 0 {
		if err := validateArgumentOrdering(cmd); err != nil {
			return err
		}
		if err := validateArgumentConstraints(cmd); err != nil {
			return err
		}
	}

	if len(cmd.Flags) > 0 {
		if err := validateFlagFileReferences(cmd, definedConfigFiles); err != nil {
			return err
		}
		if err := validateFlagConstraints(cmd); err != nil {
			return err
		}
	}

	// Validate group commands don't have args or flags
	if cmd.Kind == spec.CommandKindGroup {
		if len(cmd.Args) > 0 {
			return &ValidationError{
				Message: "group command cannot have arguments",
				Path:    fmt.Sprintf("command '%s'", cmd.CommandLine),
			}
		}
		if len(cmd.Flags) > 0 {
			return &ValidationError{
				Message: "group command cannot have flags",
				Path:    fmt.Sprintf("command '%s'", cmd.CommandLine),
			}
		}
	}

	return nil
}

// validateArgumentOrdering ensures required positional args don't come after optional ones
func validateArgumentOrdering(cmd *spec.CommandItem) error {
	seenOptional := false

	for i, arg := range cmd.Args {
		isRequired := arg.Required

		// If we've seen an optional argument, required args after it are invalid
		if seenOptional && isRequired {
			return &ValidationError{
				Message: fmt.Sprintf("required positional argument '%s' cannot come after optional arguments", arg.Name),
				Path:    fmt.Sprintf("command '%s', arg index %d", cmd.CommandLine, i),
			}
		}

		if !isRequired {
			seenOptional = true
		}
	}

	return nil
}

// validateArgumentConstraints checks that minItems/maxItems are only used with variadic args
func validateArgumentConstraints(cmd *spec.CommandItem) error {
	for i, arg := range cmd.Args {
		// Check minItems/maxItems are only used with variadic
		if (arg.MinItems > 0 || arg.MaxItems > 0) && !arg.Variadic {
			var field string
			if arg.MinItems > 0 {
				field = "minItems"
			} else {
				field = "maxItems"
			}
			return &ValidationError{
				Message: fmt.Sprintf("argument '%s' has %s but is not variadic", arg.Name, field),
				Path:    fmt.Sprintf("command '%s', args[%d]", cmd.CommandLine, i),
			}
		}

		// Check minItems <= maxItems if both are specified
		if arg.Variadic && arg.MinItems > 0 && arg.MaxItems > 0 && arg.MinItems > arg.MaxItems {
			return &ValidationError{
				Message: fmt.Sprintf("argument '%s' has minItems (%d) greater than maxItems (%d)", arg.Name, arg.MinItems, arg.MaxItems),
				Path:    fmt.Sprintf("command '%s', args[%d]", cmd.CommandLine, i),
			}
		}
	}
	return nil
}

// validateFlagFileReferences checks that $FILE references exist in global config
func validateFlagFileReferences(cmd *spec.CommandItem, definedConfigFiles map[string]bool) error {
	for i, flag := range cmd.Flags {
		if err := validateFileReferences(cmd.CommandLine, flag.Name, "flag", i, flag.AltSources, definedConfigFiles); err != nil {
			return err
		}
	}
	return nil
}

// validateFileReferences checks that $FILE references exist in global config
func validateFileReferences(cmdLine, itemName, itemType string, index int, altSources []spec.AlternativeSource, definedConfigFiles map[string]bool) error {
	for j, source := range altSources {
		if source.Type == "$FILE" {
			// Verify that at least one config file is defined globally
			if len(definedConfigFiles) == 0 {
				return &ValidationError{
					Message: fmt.Sprintf("%s '%s' references $FILE but no config files are defined in global.config", itemType, itemName),
					Path:    fmt.Sprintf("command '%s', %s[%d].alternativeSources[%d]", cmdLine, itemType, index, j),
				}
			}
		}
	}
	return nil
}

// validateFlagConstraints checks for duplicate flag names/aliases and variadic+required constraints
func validateFlagConstraints(cmd *spec.CommandItem) error {
	// Check for duplicate flag names and aliases
	seen := make(map[string]int)

	for i, flag := range cmd.Flags {
		flagName := flag.Name
		if flagName == "" {
			continue
		}

		if prevIdx, exists := seen[flagName]; exists {
			return &ValidationError{
				Message: fmt.Sprintf("duplicate flag name '%s' (also defined at index %d)", flagName, prevIdx),
				Path:    fmt.Sprintf("command '%s', flags[%d]", cmd.CommandLine, i),
			}
		}
		seen[flagName] = i

		// Check aliases
		for _, alias := range flag.Aliases {
			if alias == "" {
				continue
			}

			if prevIdx, exists := seen[alias]; exists {
				return &ValidationError{
					Message: fmt.Sprintf("duplicate flag alias '%s' (already defined at index %d)", alias, prevIdx),
					Path:    fmt.Sprintf("command '%s', flags[%d].aliases", cmd.CommandLine, i),
				}
			}
			seen[alias] = i
		}

		// Validate variadic flags aren't required
		if flag.Variadic && flag.Required {
			return &ValidationError{
				Message: fmt.Sprintf("variadic flag '%s' cannot be marked as required (variadic flags can be provided 0 or more times)", flagName),
				Path:    fmt.Sprintf("command '%s', flags[%d]", cmd.CommandLine, i),
			}
		}

		// Check minItems/maxItems are only used with variadic
		if (flag.MinItems > 0 || flag.MaxItems > 0) && !flag.Variadic {
			var field string
			if flag.MinItems > 0 {
				field = "minItems"
			} else {
				field = "maxItems"
			}
			return &ValidationError{
				Message: fmt.Sprintf("flag '%s' has %s but is not variadic", flagName, field),
				Path:    fmt.Sprintf("command '%s', flags[%d]", cmd.CommandLine, i),
			}
		}

		// Check minItems <= maxItems if both are specified
		if flag.Variadic && flag.MinItems > 0 && flag.MaxItems > 0 && flag.MinItems > flag.MaxItems {
			return &ValidationError{
				Message: fmt.Sprintf("flag '%s' has minItems (%d) greater than maxItems (%d)", flagName, flag.MinItems, flag.MaxItems),
				Path:    fmt.Sprintf("command '%s', flags[%d]", cmd.CommandLine, i),
			}
		}
	}

	return nil
}
