package validator

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

//go:generate cp -r ../../schemas/ ./distschemas
//go:embed distschemas
var schemas embed.FS

// ValidateDocument validates the given document against the OpenCLI Specification.
// Any validation errors are returned. If no error is returned, the document is valid.
func ValidateDocument(document []byte) error {
	// First attemp to unmarshal the document to ensure we're dealing with a valid JSON document
	var docObj documentVersion
	if err := json.Unmarshal(document, &docObj); err != nil {
		return fmt.Errorf("error unmarshalling JSON document: %w", err)
	}
	// Select the version of the OpenCLI Specification to validate the given document against
	if docObj.Version == "" {
		return fmt.Errorf("missing 'openCliVersion' field in document")
	}
	// Read and unmarshal the JSON schema file for the selected OpenCLI Spec version
	contents, err := schemas.ReadFile(fmt.Sprintf("distschemas/%s_specification.schema.json", docObj.Version))
	if err != nil {
		return fmt.Errorf("unsupported OpenCLI version %s; use one of %v", docObj.Version, Versions())
	}
	schema, err := jsonschema.UnmarshalJSON(bytes.NewReader(contents))
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON schema: %w", err)
	}
	// Create a new JSON Schema compiler and add the OpenCLI Spec JSON Schema to it
	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", schema); err != nil {
		return fmt.Errorf("error compiling JSON schema: %w", err)
	}
	// Compile the JSON Schema
	v, err := c.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("error compiling schema: %w", err)
	}
	// Unmarshal the given document and validate it against the JSON Schema
	inst, err := jsonschema.UnmarshalJSON(strings.NewReader(string(document)))
	if err != nil {
		return fmt.Errorf("error unmarshalling document: %w", err)
	}
	// Return the result of the JSON Schema validation against the Open CLI Spec
	return v.Validate(inst)
}

// Versions returns a list of supported OpenCLI Specification versions.
func Versions() []string {
	specSchemaFiles, _ := schemas.ReadDir("distschemas")

	var versionStrings []string
	for _, entry := range specSchemaFiles {
		versionStrings = append(versionStrings, strings.Split(entry.Name(), "_")[0])
	}

	return versionStrings
}

/* Private types and functions
------------------------------------------------------------------------------------------------- */

type documentVersion struct {
	Version string `json:"openCliVersion"`
}
