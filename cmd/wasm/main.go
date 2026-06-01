//go:build !test && js && wasm

package main

import (
	"encoding/json"
	"errors"
	"syscall/js"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/gen"
	"github.com/bcdxn/opencli/validate"
)

type jsValidationError struct {
	Message string `json:"message"`
	Path    string `json:"path"`
}

// jsValidateOCS is exposed as window.validateOCS(input, format).
// Returns { valid: bool, errors: [{message, path}] }. Never panics.
func jsValidateOCS(this js.Value, args []js.Value) any {
	failWith := func(msg, path string) any {
		return jsonToJS(map[string]any{
			"valid":  false,
			"errors": []jsValidationError{{Message: msg, Path: path}},
		})
	}

	if len(args) < 2 {
		return failWith("missing arguments: input and format required", "")
	}

	input := []byte(args[0].String())
	format := args[1].String()

	var err error
	switch format {
	case "yaml":
		err = validate.ValidateYAML(input)
	case "json":
		err = validate.ValidateJSON(input)
	default:
		return failWith("unsupported format: must be 'yaml' or 'json'", "")
	}

	if err == nil {
		return jsonToJS(map[string]any{
			"valid":  true,
			"errors": []jsValidationError{},
		})
	}

	var ve *validate.ValidationError
	if errors.As(err, &ve) {
		return jsonToJS(map[string]any{
			"valid":  false,
			"errors": []jsValidationError{{Message: ve.Message, Path: ve.Path}},
		})
	}

	return jsonToJS(map[string]any{
		"valid":  false,
		"errors": []jsValidationError{{Message: err.Error(), Path: ""}},
	})
}

// jsGenerateOCSDocs is exposed as window.generateOCSDocs(input, inputFormat, outputFormat, htmlFlavor).
// Returns { output: string, error: string }. Never panics.
func jsGenerateOCSDocs(this js.Value, args []js.Value) any {
	failWith := func(msg string) any {
		return jsonToJS(map[string]any{"output": "", "error": msg})
	}

	if len(args) < 4 {
		return failWith("missing arguments: input, inputFormat, outputFormat, htmlFlavor required")
	}

	input := []byte(args[0].String())
	inputFormat := args[1].String()
	outputFormat := args[2].String()
	htmlFlavor := args[3].String()

	// Parse input document
	unmarshal := codec.UnmarshalYAML
	switch inputFormat {
	case "yaml":
		unmarshal = codec.UnmarshalYAML
	case "json":
		unmarshal = codec.UnmarshalJSON
	default:
		return failWith("unsupported inputFormat: must be 'yaml' or 'json'")
	}

	doc, err := unmarshal(input)
	if err != nil {
		return failWith(err.Error())
	}

	// Build generation options
	var opts []gen.GenDocsOption
	switch outputFormat {
	case "markdown":
		opts = append(opts, gen.DocsWithFormat(gen.Markdown))
	case "html":
		opts = append(opts, gen.DocsWithFormat(gen.HTML))
		switch htmlFlavor {
		case "page":
			opts = append(opts, gen.DocsWithHTMLFlavor(gen.StandalonePage))
		case "component":
			opts = append(opts, gen.DocsWithHTMLFlavor(gen.EmbeddableComponent))
		default:
			return failWith("unsupported htmlFlavor: must be 'page' or 'component'")
		}
	default:
		return failWith("unsupported outputFormat: must be 'markdown' or 'html'")
	}

	result, err := gen.Docs(doc, opts...)
	if err != nil {
		return failWith(err.Error())
	}

	return jsonToJS(map[string]any{"output": string(result), "error": ""})
}

// jsonToJS marshals v to JSON then parses it into a JS object via JSON.parse.
func jsonToJS(v any) js.Value {
	b, err := json.Marshal(v)
	if err != nil {
		b, _ = json.Marshal(map[string]any{"error": err.Error()})
	}
	return js.Global().Get("JSON").Call("parse", string(b))
}

func main() {
	js.Global().Set("validateOCS", js.FuncOf(jsValidateOCS))
	js.Global().Set("generateOCSDocs", js.FuncOf(jsGenerateOCSDocs))
	select {} // keep WASM alive
}
