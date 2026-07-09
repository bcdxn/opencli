package validate_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/bcdxn/opencli/validate"
)

//go:generate mkdir -p out
//go:generate cp -r ../examples/petstore-cli.ocs.json ./out/petstore-cli.ocs.json
//go:generate cp -r ../examples/petstore-cli.ocs.yaml ./out/petstore-cli.ocs.yaml
//go:generate cp -r ../examples/pleasantries-cli.ocs.yaml ./out/pleasantries-cli.ocs.yaml

//go:embed out/petstore-cli.ocs.yaml
var petstoreYAML []byte

//go:embed out/petstore-cli.ocs.json
var petstoreJSON []byte

//go:embed out/pleasantries-cli.ocs.yaml
var pleasantriesYAML []byte

func TestValidateJSON(t *testing.T) {
	err := validate.ValidateJSON(petstoreJSON)
	if err != nil {
		t.Fatalf("should have validated successfully but found err: %v", err)
	}
}

func TestValidateYAML(t *testing.T) {
	err := validate.ValidateYAML(petstoreYAML)
	if err != nil {
		t.Fatalf("should have validated successfully but found err: %v", err)
	}

	err = validate.ValidateYAML(pleasantriesYAML)
	if err != nil {
		t.Fatalf("should have validated successfully but found err: %v", err)
	}
}

func TestValidateJSON_InvalidJSON(t *testing.T) {
	err := validate.ValidateJSON([]byte("{"))
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
	if !strings.Contains(err.Error(), "error unmarshalling document") {
		t.Fatalf("expected unmarshal error, got: %v", err)
	}
}

func TestValidateYAML_LogicalValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr string
	}{
		{
			name: "group command with flags",
			input: replaceOnce(t, petstoreYAML,
				"  petstore {command} <arguments> [flags]:\n    kind: group",
				"  petstore {command} <arguments> [flags]:\n    kind: group\n    flags:\n    - name: verbose\n      type: boolean"),
			wantErr: "group command cannot have flags",
		},
		{
			name: "required positional arg after optional",
			input: replaceOnce(t, petstoreYAML,
				"    - name: path-to-user-body\n      type: string\n      summary: The path to a JSON file containing the user payload\n      required: false\n    flags:",
				"    - name: path-to-user-body\n      type: string\n      summary: The path to a JSON file containing the user payload\n      required: false\n    - name: user-id\n      type: string\n      summary: A required user id\n      required: true\n    flags:"),
			wantErr: "required positional argument 'user-id' cannot come after optional arguments",
		},
		{
			name: "arg minItems without variadic",
			input: replaceOnce(t, pleasantriesYAML,
				"      - name: \"name\"\n        summary: \"A name to include in the greeting\"\n        required: true\n        type: \"string\"",
				"      - name: \"name\"\n        summary: \"A name to include in the greeting\"\n        required: true\n        type: \"string\"\n        minItems: 1"),
			wantErr: "argument 'name' has minItems but is not variadic",
		},
		{
			name: "variadic flag cannot be required",
			input: replaceOnce(t, petstoreYAML,
				"      variadic: true",
				"      variadic: true\n      required: true"),
			wantErr: "variadic flag 'photo-urls' cannot be marked as required",
		},
		{
			name: "duplicate flag alias",
			input: replaceOnce(t, pleasantriesYAML,
				"    examples:",
				"      - name: \"lang\"\n        aliases:\n          - \"language\"\n        summary: \"Duplicate alias\"\n        type: \"string\"\n    examples:"),
			wantErr: "duplicate flag alias 'language'",
		},
		{
			name: "file source without global config",
			input: replaceOnce(t, pleasantriesYAML,
				"        default: \"english\"",
				"        default: \"english\"\n        alternativeSources:\n          - type: \"$FILE\"\n            property: \"$.greet.language\""),
			wantErr: "references $FILE but no config files are defined in global.config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.ValidateYAML(tt.input)
			if err == nil {
				t.Fatal("expected validation error but got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

func replaceOnce(t *testing.T, input []byte, old, new string) []byte {
	t.Helper()
	original := string(input)
	updated := strings.Replace(original, old, new, 1)
	if updated == original {
		t.Fatalf("failed to apply test mutation; snippet not found: %q", old)
	}
	return []byte(updated)
}
