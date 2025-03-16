package validator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

//go:generate cp -r ../../spec/ ./distschemas
//go:embed distschemas
var schemas embed.FS

// ValidateDocumentJSON validates the given JSON document against the OpenCLI Specification.
// Any validation errors are returned.
// If no error is returned, the document is valid.
func ValidateDocumentJSON(document []byte) error {
	schema, err := schemaValidator(document)
	if err != nil {
		return err
	}
	// Unmarshal the given document and validate it against the JSON Schema
	var specDoc any
	err = json.Unmarshal(document, &specDoc)
	if err != nil {
		return fmt.Errorf("error unmarshalling document: %w", err)
	}
	// Return the result of the JSON Schema validation against the Open CLI Spec
	return schema.Validate(specDoc)
}

// MustValidateDocumentJSON is a wrapper of ValidateDocumentJSON that panics if the document is invalid.
func MustValidateDocumentJSON(document []byte) {
	if err := ValidateDocumentJSON(document); err != nil {
		panic(err)
	}
}

// ValidateDocumentYAML validates the given YAML document against the OpenCLI Specification.
// Any validation errors are returned.
// If no error is returned, the document is valid.
func ValidateDocumentYAML(document []byte) error {
	schema, err := schemaValidator(document)
	if err != nil {
		return err
	}
	// Unmarshal the given document and validate it against the JSON Schema
	var specDoc any
	err = yaml.Unmarshal(document, &specDoc)
	if err != nil {
		return fmt.Errorf("error unmarshalling document: %w", err)
	}
	// Return the result of the JSON Schema validation against the Open CLI Spec
	return schema.Validate(specDoc)
}

// MustValidateDocumentYAML is a wrapper of ValidateDocumentYAML that panics if the document is invalid.
func MustValidateDocumentYAML(document []byte) {
	if err := ValidateDocumentYAML(document); err != nil {
		panic(err)
	}
}

// Versions returns a list of supported OpenCLI Specification versions.
func Versions() []string {
	specSchemaFiles, _ := schemas.ReadDir("distschemas")

	var versionStrings []string
	for _, entry := range specSchemaFiles {
		version := strings.Split(entry.Name(), "_")[0]
		versionStrings = append(versionStrings, version)
	}

	return versionStrings
}

/* Private types and functions
------------------------------------------------------------------------------------------------- */

type documentVersion struct {
	Version string `json:"openCliVersion"`
}

// schemaValidator returns a JSON Schema that can be used for validation.
func schemaValidator(document []byte) (*jsonschema.Schema, error) {
	// First attempt to unmarshal the document to ensure we're dealing with a valid JSON document
	var docObj documentVersion
	if err := json.Unmarshal(document, &docObj); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON document: %w", err)
	}
	// Select the version of the OpenCLI Specification to validate the given document against
	if docObj.Version == "" {
		return nil, fmt.Errorf("missing 'openCliVersion' field in document")
	}
	// Read and unmarshal the JSON schema file for the selected OpenCLI Spec version
	contents, err := schemas.ReadFile(fmt.Sprintf("distschemas/%s_specification.schema.json", docObj.Version))
	if err != nil {
		return nil, fmt.Errorf("unsupported OpenCLI version %s; use one of %v", docObj.Version, Versions())
	}
	schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(contents))
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON schema: %w", err)
	}
	// Create a new JSON Schema compiler and add the OpenCLI Spec JSON Schema to it
	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", schema); err != nil {
		return nil, fmt.Errorf("error compiling JSON schema: %w", err)
	}
	// Compile the JSON Schema
	v, err := c.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("error compiling schema: %w", err)
	}
	return v, nil
}
