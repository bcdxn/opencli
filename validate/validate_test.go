package validate_test

import (
	_ "embed"
	"testing"

	"github.com/bcdxn/opencli/validate"
)

//go:generate mkdir -p gen
//go:generate cp -r ../examples/petstore-cli.ocs.json ./gen/petstore-cli.ocs.json
//go:generate cp -r ../examples/petstore-cli.ocs.yaml ./gen/petstore-cli.ocs.yaml
//go:generate cp -r ../examples/pleasantries-cli.ocs.yaml ./gen/pleasantries-cli.ocs.yaml

//go:embed gen/petstore-cli.ocs.yaml
var petstoreYAML []byte

//go:embed gen/petstore-cli.ocs.json
var petstoreJSON []byte

//go:embed gen/pleasantries-cli.ocs.yaml
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
