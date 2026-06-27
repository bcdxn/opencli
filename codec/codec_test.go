package codec_test

import (
	"bytes"
	_ "embed"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/spec"
)

//go:generate mkdir -p out
//go:generate cp -r ../examples/petstore-cli.ocs.json ./out/petstore-cli.ocs.json
//go:generate cp -r ../examples/petstore-cli.ocs.yaml ./out/petstore-cli.ocs.yaml

var (
	goldenJSON = "testdata/expected.json"
	goldenYAML = "testdata/expected.yaml"
)

//go:embed out/petstore-cli.ocs.yaml
var exampleYAML []byte

//go:embed out/petstore-cli.ocs.json
var exampleJSON []byte

var update = flag.Bool("update", false, "update golden files for CLI generation tests")

func TestUnmarshalYAML(t *testing.T) {
	d, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unmarshal yaml failed: %v", err)
	}
	validateSpec(t, d)
}

func TestUnmarshalJSON(t *testing.T) {
	d, err := codec.UnmarshalJSON(exampleJSON)
	if err != nil {
		t.Fatalf("unmarshal json failed: %v", err)
	}
	validateSpec(t, d)
}

func TestMarshalJSON(t *testing.T) {
	doc, err := codec.UnmarshalJSON(exampleJSON)
	if err != nil {
		t.Fatalf("unmarshal json failed: %v", err)
	}

	actual, err := codec.MarshalJSON(doc)
	if err != nil {
		t.Fatalf("marshal json failed: %v", err)
	}

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenJSON), 0755); err != nil {
			t.Fatalf("failed to create golden dir for %s: %v", goldenJSON, err)
		}
		if err := os.WriteFile(goldenJSON, actual, 0644); err != nil {
			t.Fatalf("failed to write golden file %s: %v", goldenJSON, err)
		}
	}

	expected, err := os.ReadFile("testdata/expected.json")
	if err != nil {
		t.Fatalf("error writing output file: %v", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("marshaled JSON content does not match")
	}
}

func TestMarshalYAML(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unmarshal yaml failed: %v", err)
	}

	actual, err := codec.MarshalYAML(doc)
	if err != nil {
		t.Fatalf("marshal yaml failed: %v", err)
	}

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenYAML), 0755); err != nil {
			t.Fatalf("failed to create golden dir for %s: %v", goldenYAML, err)
		}
		if err := os.WriteFile(goldenYAML, actual, 0644); err != nil {
			t.Fatalf("failed to write golden file %s: %v", goldenYAML, err)
		}
	}

	expected, err := os.ReadFile("testdata/expected.yaml")
	if err != nil {
		t.Fatalf("error writing output file: %v", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("marshaled YAML content does not match")
	}
}

func validateSpec(t *testing.T, d *spec.Document) {
	t.Helper()

	if d.Info.Title != "PetStore CLI" {
		t.Fatalf("unexpected title: %s", d.Info.Title)
	}

	if d.Commands == nil {
		t.Fatal("expected root commands to be built, got nil")
	}

	root := d.Commands
	if root.Segment != "petstore" {
		t.Fatalf("expected root segment petstore, got %s", root.Segment)
	}
	if root.CommandLine != "petstore" {
		t.Fatalf("expected root command line petstore, got %s", root.CommandLine)
	}

	if got := len(root.Commands); got != 4 {
		t.Fatalf("expected 4 root subcommands, got %d", got)
	}

	list := findSubcommand(root, "list")
	if list == nil {
		t.Fatal("expected root to contain list subcommand")
	}
	if list.CommandLine != "petstore list" {
		t.Fatalf("expected petstore list command line, got %s", list.CommandLine)
	}

	pet := findSubcommand(root, "pet")
	if pet == nil {
		t.Fatal("expected root to contain pet subcommand")
	}
	if pet.CommandLine != "petstore pet" {
		t.Fatalf("expected pet command line petstore pet, got %s", pet.CommandLine)
	}

	store := findSubcommand(root, "store")
	if store == nil {
		t.Fatal("expected root to contain store subcommand")
	}
	if store.CommandLine != "petstore store" {
		t.Fatalf("expected store command line petstore store, got %s", store.CommandLine)
	}

	// 'user' is a derived 'group' command that is not explicitly defined in the document
	user := findSubcommand(root, "user")
	if user == nil {
		t.Fatal("expected root to contain user subcommand")
	}
	if user.CommandLine != "petstore user" {
		t.Fatalf("expected user command line petstore user, got %s", user.CommandLine)
	}
	if !user.Group {
		t.Fatal("expected petstore user command to be a 'group'")
	}

	expectedStoreChildren := []string{"order", "inventory"}
	if got := len(store.Commands); got != len(expectedStoreChildren) {
		t.Fatalf("expected %d store subcommands, got %d", len(expectedStoreChildren), got)
	}
	for _, segment := range expectedStoreChildren {
		child := findSubcommand(store, segment)
		if child == nil {
			t.Fatalf("expected store to contain %q subcommand", segment)
		}
	}

	expectedOrderChildren := []string{"place", "get", "delete"}
	order := findSubcommand(store, "order")
	if order == nil {
		t.Fatal("expected store to contain order subcommand")
	}
	if got := len(order.Commands); got != len(expectedOrderChildren) {
		t.Fatalf("expected %d store order subcommands, got %d", len(expectedOrderChildren), got)
	}
	for _, segment := range expectedOrderChildren {
		child := findSubcommand(order, segment)
		if child == nil {
			t.Fatalf("expected store order to contain %q subcommand", segment)
		}
	}

	expectedUserChildren := []string{"login", "create-with-list", "get", "delete", "create", "logout", "update"}
	if got := len(user.Commands); got != len(expectedUserChildren) {
		t.Fatalf("expected %d user subcommands, got %d", len(expectedUserChildren), got)
	}
	for _, segment := range expectedUserChildren {
		child := findSubcommand(user, segment)
		if child == nil {
			t.Fatalf("expected user to contain %q subcommand", segment)
		}
	}

	expectedPetChildren := []string{
		"add",
		"update",
		"find-by-status",
		"find-by-tags",
		"get",
		"update-form",
		"delete",
		"upload-image",
	}
	if got := len(pet.Commands); got != len(expectedPetChildren) {
		t.Fatalf("expected %d pet subcommands, got %d", len(expectedPetChildren), got)
	}

	for _, segment := range expectedPetChildren {
		child := findSubcommand(pet, segment)
		if child == nil {
			t.Fatalf("expected pet to contain %q subcommand", segment)
		}
		if child.Segment != segment {
			t.Fatalf("expected pet subcommand segment %q, got %q", segment, child.Segment)
		}
		if child.CommandLine != "petstore pet "+segment {
			t.Fatalf("expected command line petstore pet %s, got %s", segment, child.CommandLine)
		}
	}
}

func findSubcommand(parent *spec.CommandItem, segment string) *spec.CommandItem {
	for _, cmd := range parent.Commands {
		if cmd.Segment == segment {
			return cmd
		}
	}
	return nil
}
