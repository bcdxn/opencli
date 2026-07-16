package ocobra

import (
	"bytes"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/bcdxn/opencli/spec"
	"github.com/spf13/cobra"
)

func TestGenerateDocument_Basic(t *testing.T) {
	root := &cobra.Command{
		Use:   "myapp",
		Short: "An awesome CLI",
		Long:  "A longer description of myapp",
	}

	child := &cobra.Command{
		Use:   "greet <name>",
		Short: "Say hello",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	root.AddCommand(child)

	doc := GenerateDocument(root)

	if doc.OpenCLIVersion != defaultOpenCLIVersion {
		t.Fatalf("unexpected version: %s", doc.OpenCLIVersion)
	}
	if doc.Info.Binary != "myapp" {
		t.Fatalf("expected binary 'myapp', got %q", doc.Info.Binary)
	}
	if doc.Commands == nil {
		t.Fatal("commands should not be nil")
	}
	if doc.Commands.Segment != "myapp" {
		t.Fatalf("expected root segment 'myapp', got %q", doc.Commands.Segment)
	}
	if len(doc.Commands.Commands) != 1 {
		t.Fatalf("expected 1 child, got %d", len(doc.Commands.Commands))
	}
	if doc.Commands.Commands[0].Segment != "greet" {
		t.Fatalf("expected child 'greet', got %q", doc.Commands.Commands[0].Segment)
	}
}

func TestGenerateDocument_GroupVsAction(t *testing.T) {
	root := &cobra.Command{Use: "app"}

	groupCmd := &cobra.Command{
		Use: "group",
	}
	groupCmd.AddCommand(&cobra.Command{Use: "sub"})
	actionCmd := &cobra.Command{
		Use: "action",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	bothCmd := &cobra.Command{
		Use: "both",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	bothCmd.AddCommand(&cobra.Command{Use: "nested"})

	root.AddCommand(groupCmd, actionCmd, bothCmd)

	doc := GenerateDocument(root)

	// Cobra sorts subcommands alphabetically, so find by name
	cmdMap := make(map[string]*spec.CommandItem)
	for _, c := range doc.Commands.Commands {
		cmdMap[c.Segment] = c
	}

	if cmdMap["group"].Kind != spec.CommandKindGroup {
		t.Errorf("group should be group kind, got %s", cmdMap["group"].Kind)
	}
	if cmdMap["action"].Kind != spec.CommandKindAction {
		t.Errorf("action should be action kind, got %s", cmdMap["action"].Kind)
	}
	if cmdMap["both"].Kind != spec.CommandKindAction {
		t.Errorf("both (run+subs) should be action kind, got %s", cmdMap["both"].Kind)
	}
}

func TestGenerateDocument_WithInfo(t *testing.T) {
	info := &spec.Info{
		Title:   "My App",
		Version: "2.0.0",
		Binary:  "myapp",
	}

	root := &cobra.Command{Use: "myapp"}
	doc := GenerateDocument(root, WithInfo(info))

	if doc.Info.Title != "My App" {
		t.Errorf("expected title 'My App', got %q", doc.Info.Title)
	}
	if doc.Info.Version != "2.0.0" {
		t.Errorf("expected version '2.0.0', got %q", doc.Info.Version)
	}
}

func TestGenerateDocument_WithInstallMethods(t *testing.T) {
	install := []spec.InstallMethod{
		{Name: "brew", Command: "brew install myapp"},
	}

	root := &cobra.Command{Use: "myapp"}
	doc := GenerateDocument(root, WithInstallMethods(install))

	if len(doc.Install) != 1 {
		t.Fatalf("expected 1 install method, got %d", len(doc.Install))
	}
	if doc.Install[0].Name != "brew" {
		t.Errorf("expected install name 'brew', got %q", doc.Install[0].Name)
	}
}

func TestGenerateDocument_WithGlobalFlags(t *testing.T) {
	flags := []spec.FlagItem{
		{Name: "verbose", Type: "boolean"},
	}

	root := &cobra.Command{Use: "myapp"}
	doc := GenerateDocument(root, WithGlobalFlags(flags))

	if doc.Global == nil {
		t.Fatal("global should not be nil")
	}
	if len(doc.Global.Flags) != 1 {
		t.Fatalf("expected 1 global flag, got %d", len(doc.Global.Flags))
	}
}

func TestGenerateDocument_Flags(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmd.Flags().String("output", "", "output file")
	cmd.Flags().BoolP("verbose", "v", false, "enable verbose mode")

	doc := GenerateDocument(cmd)
	if len(doc.Commands.Flags) != 2 {
		t.Fatalf("expected 2 flags, got %d", len(doc.Commands.Flags))
	}

	flagMap := make(map[string]spec.FlagItem)
	for _, f := range doc.Commands.Flags {
		flagMap[f.Name] = f
	}

	outFlag, ok := flagMap["output"]
	if !ok {
		t.Fatal("missing 'output' flag")
	}
	if outFlag.Type != "string" {
		t.Errorf("expected type 'string', got %q", outFlag.Type)
	}

	verbFlag, ok := flagMap["verbose"]
	if !ok {
		t.Fatal("missing 'verbose' flag")
	}
	if verbFlag.Type != "boolean" {
		t.Errorf("expected type 'boolean', got %q", verbFlag.Type)
	}
	if len(verbFlag.Aliases) == 0 || verbFlag.Aliases[0] != "v" {
		t.Errorf("expected alias 'v', got %v", verbFlag.Aliases)
	}
}

func TestGenerateDocument_Arguments(t *testing.T) {
	cmd := &cobra.Command{
		Use: "copy <src> <dst>",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	doc := GenerateDocument(cmd)
	if len(doc.Commands.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(doc.Commands.Args))
	}
	if doc.Commands.Args[0].Name != "src" {
		t.Errorf("expected arg 'src', got %q", doc.Commands.Args[0].Name)
	}
	if !doc.Commands.Args[0].Required {
		t.Error("src should be required")
	}
}

func TestGenerateDocument_OptionalArguments(t *testing.T) {
	cmd := &cobra.Command{
		Use: "show [name]",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	doc := GenerateDocument(cmd)
	if len(doc.Commands.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(doc.Commands.Args))
	}
	if doc.Commands.Args[0].Required {
		t.Error("[name] should NOT be required")
	}
}

func TestGenerateDocument_VariadicArguments(t *testing.T) {
	cmd := &cobra.Command{
		Use: "rm <files>...",
		Run: func(cmd *cobra.Command, args []string) {},
	}

	doc := GenerateDocument(cmd)
	if len(doc.Commands.Args) != 1 {
		t.Fatalf("expected 1 arg, got %d", len(doc.Commands.Args))
	}
	if !doc.Commands.Args[0].Variadic {
		t.Error("files should be variadic")
	}
}

func TestGenerateDocument_YAMLOutput(t *testing.T) {
	root := &cobra.Command{
		Use:   "myapp",
		Short: "A test CLI",
	}
	buf := new(bytes.Buffer)
	doc := GenerateDocument(root, WithOutput(buf))

	if doc == nil {
		t.Fatal("document should not be nil")
	}
}

func TestFromCommand_AttachesHiddenSubcommand(t *testing.T) {
	root := &cobra.Command{Use: "myapp"}
	FromCommand(root)

	found := false
	for _, c := range root.Commands() {
		if c.Use == "__opencli" {
			found = true
			if !c.Hidden {
				t.Error("__opencli should be hidden")
			}
			break
		}
	}
	if !found {
		t.Error("__opencli subcommand not found")
	}
}

func TestParseDefault(t *testing.T) {
	tests := []struct {
		typ                string
		variadic           bool
		input              string
		expected           any
		expectedMarshalled string
	}{
		// {"string", false, "", "", "{}"},
		{"boolean", false, "true", true, `{"default":true}`},
		{"boolean", false, "false", nil, `{}`},
		{"integer", false, "42", int64(42), `{"default":42}`},
		{"number", false, "3.14", float64(3.14), `{"default":3.14}`},
		{"string", false, "hello", "hello", `{"default":"hello"}`},
		{"string", true, "[hello,there]", []any{"hello", "there"}, `{"default":["hello","there"]}`},
		{"boolean", true, "[true,false]", []any{true, false}, `{"default":[true,false]}`},
		{"integer", true, "[0,1,2]", []any{int64(0), int64(1), int64(2)}, `{"default":[0,1,2]}`},
		{"number", true, "[1.2,2.2]", []any{float64(1.2), float64(2.2)}, `{"default":[1.2,2.2]}`},
	}

	for _, tc := range tests {
		got := parseDefault(tc.typ, tc.variadic, tc.input)
		if strSlice, ok := got.([]any); ok {
			if !slices.Equal(strSlice, tc.expected.([]any)) {
				t.Errorf("parseDefault(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		} else {
			if got != tc.expected {
				t.Errorf("parseDefault(%q) = %v, want %v", tc.input, got, tc.expected)
			}
		}

		unmarshalled := Unmarshalled{Default: tc.expected}
		gotMarshalled, err := json.Marshal(unmarshalled)
		if err != nil {
			t.Fatal("error marshalling test struct", err)
		}
		if string(gotMarshalled) != tc.expectedMarshalled {
			t.Errorf("marshalled(%q) = %v, want %v", tc.expected, string(gotMarshalled), tc.expectedMarshalled)
		}
	}
}

type Unmarshalled struct {
	Default any `json:"default,omitempty"`
}

func TestGetBinaryName(t *testing.T) {
	cmd := &cobra.Command{Use: "myapp [flags]"}
	if GetBinaryName(cmd) != "myapp" {
		t.Errorf("expected 'myapp', got %q", GetBinaryName(cmd))
	}
}

func TestBuildCommandLine_Nested(t *testing.T) {
	root := &cobra.Command{Use: "ocli"}
	gen := &cobra.Command{Use: "gen [flags]"}
	docs := &cobra.Command{Use: "docs <format>"}
	root.AddCommand(gen)
	gen.AddCommand(docs)

	doc := GenerateDocument(root)

	// Find the docs command
	for _, c := range doc.Commands.Commands {
		if c.Segment == "gen" {
			for _, cc := range c.Commands {
				if cc.Segment != "docs" {
					t.Errorf("expected 'docs', got %q", cc.Segment)
				}
				if !strings.Contains(cc.CommandLine, "ocli") {
					t.Error("command line should contain root")
				}
				if !strings.Contains(cc.CommandLine, "gen") {
					t.Error("command line should contain parent")
				}
			}
		}
	}
}
