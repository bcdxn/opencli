package validate_test

import (
	_ "embed"
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
