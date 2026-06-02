package validate

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"sync"

	"github.com/goccy/go-yaml"
	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

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
	var doc interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	// validate using schema
	if err := schema.Validate(doc); err != nil {
		return err
	}
	return nil
}

func ValidateYAML(data []byte) error {
	if err := ensureCompiler(); err != nil {
		return err
	}
	var doc interface{}
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return err
	}
	// validate using schema
	if err := schema.Validate(doc); err != nil {
		return err
	}
	return nil
}
