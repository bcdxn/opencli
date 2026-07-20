package ourfave

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/bcdxn/opencli/spec"
	"github.com/urfave/cli/v3"
)

func TestGenerateDocument_Basic(t *testing.T) {
	root := &cli.Command{
		Name:  "myapp",
		Usage: "An awesome CLI",
	}

	child := &cli.Command{
		Name:  "greet",
		Usage: "Say hello",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
	root.Commands = append(root.Commands, child)

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
	root := &cli.Command{Name: "app"}

	groupCmd := &cli.Command{
		Name: "group",
	}
	groupCmd.Commands = append(groupCmd.Commands, &cli.Command{Name: "sub"})
	actionCmd := &cli.Command{
		Name: "action",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
	bothCmd := &cli.Command{
		Name: "both",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
	bothCmd.Commands = append(bothCmd.Commands, &cli.Command{Name: "nested"})

	root.Commands = append(root.Commands, groupCmd, actionCmd, bothCmd)

	doc := GenerateDocument(root)

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

	root := &cli.Command{Name: "myapp"}
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

	root := &cli.Command{Name: "myapp"}
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

	root := &cli.Command{Name: "myapp"}
	doc := GenerateDocument(root, WithGlobalFlags(flags))

	if doc.Global == nil {
		t.Fatal("global should not be nil")
	}
	if len(doc.Global.Flags) != 1 {
		t.Fatalf("expected 1 global flag, got %d", len(doc.Global.Flags))
	}
}

func TestGenerateDocument_Arguments(t *testing.T) {
	cmd := &cli.Command{
		Name: "copy",
		Arguments: []cli.Argument{
			&cli.StringArg{Name: "src"},
			&cli.StringArg{Name: "dst"},
			&cli.IntArg{Name: "count"},
			&cli.FloatArg{Name: "ratio"},
			&cli.Int8Arg{Name: "i8"},
			&cli.Int16Arg{Name: "i16"},
			&cli.Int32Arg{Name: "i32"},
			&cli.Int64Arg{Name: "i64"},
			&cli.UintArg{Name: "u"},
			&cli.Uint8Arg{Name: "u8"},
			&cli.Uint16Arg{Name: "u16"},
			&cli.Uint32Arg{Name: "u32"},
			&cli.Uint64Arg{Name: "u64"},
			&cli.TimestampArg{Name: "ts"},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}

	doc := GenerateDocument(cmd)
	if len(doc.Commands.Args) != 14 {
		t.Fatalf("expected 14 args, got %d", len(doc.Commands.Args))
	}

	tests := []struct {
		name string
		typ  string
	}{
		{"src", "string"},
		{"dst", "string"},
		{"count", "integer"},
		{"ratio", "number"},
		{"i8", "integer"},
		{"i16", "integer"},
		{"i32", "integer"},
		{"i64", "integer"},
		{"u", "integer"},
		{"u8", "integer"},
		{"u16", "integer"},
		{"u32", "integer"},
		{"u64", "integer"},
		{"ts", "string"},
	}

	for i, tc := range tests {
		if doc.Commands.Args[i].Name != tc.name {
			t.Errorf("arg[%d]: expected name %q, got %q", i, tc.name, doc.Commands.Args[i].Name)
		}
		if doc.Commands.Args[i].Type != tc.typ {
			t.Errorf("arg[%d] (%s): expected type %q, got %q", i, tc.name, tc.typ, doc.Commands.Args[i].Type)
		}
	}
}

func TestGenerateDocument_Flags(t *testing.T) {
	cmd := &cli.Command{
		Name: "test",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "str"},
			&cli.BoolFlag{Name: "bool"},
			&cli.IntFlag{Name: "int"},
			&cli.Int8Flag{Name: "i8"},
			&cli.Int16Flag{Name: "i16"},
			&cli.Int32Flag{Name: "i32"},
			&cli.Int64Flag{Name: "i64"},
			&cli.UintFlag{Name: "uint"},
			&cli.Uint8Flag{Name: "u8"},
			&cli.Uint16Flag{Name: "u16"},
			&cli.Uint32Flag{Name: "u32"},
			&cli.Uint64Flag{Name: "u64"},
			&cli.FloatFlag{Name: "float"},
			&cli.Float32Flag{Name: "f32"},
			&cli.Float64Flag{Name: "f64"},
			&cli.DurationFlag{Name: "dur"},
			&cli.StringSliceFlag{Name: "str-slice"},
			&cli.IntSliceFlag{Name: "int-slice"},
			&cli.Int8SliceFlag{Name: "i8-slice"},
			&cli.Int16SliceFlag{Name: "i16-slice"},
			&cli.Int32SliceFlag{Name: "i32-slice"},
			&cli.Int64SliceFlag{Name: "i64-slice"},
			&cli.UintSliceFlag{Name: "uint-slice"},
			&cli.Uint8SliceFlag{Name: "u8-slice"},
			&cli.Uint16SliceFlag{Name: "u16-slice"},
			&cli.Uint32SliceFlag{Name: "u32-slice"},
			&cli.Uint64SliceFlag{Name: "u64-slice"},
			&cli.FloatSliceFlag{Name: "float-slice"},
			&cli.Float32SliceFlag{Name: "f32-slice"},
			&cli.Float64SliceFlag{Name: "f64-slice"},
		},
	}

	doc := GenerateDocument(cmd)
	if len(doc.Commands.Flags) != 30 { // add 2 for version and help defaults
		t.Fatalf("expected 30 flags, got %d", len(doc.Commands.Flags))
	}

	tests := []struct {
		name     string
		typ      string
		variadic bool
	}{
		{"str", "string", false},
		{"bool", "boolean", false},
		{"int", "integer", false},
		{"i8", "integer", false},
		{"i16", "integer", false},
		{"i32", "integer", false},
		{"i64", "integer", false},
		{"uint", "integer", false},
		{"u8", "integer", false},
		{"u16", "integer", false},
		{"u32", "integer", false},
		{"u64", "integer", false},
		{"float", "number", false},
		{"f32", "number", false},
		{"f64", "number", false},
		{"dur", "string", false},
		{"str-slice", "string", true},
		{"int-slice", "integer", true},
		{"i8-slice", "integer", true},
		{"i16-slice", "integer", true},
		{"i32-slice", "integer", true},
		{"i64-slice", "integer", true},
		{"uint-slice", "integer", true},
		{"u8-slice", "integer", true},
		{"u16-slice", "integer", true},
		{"u32-slice", "integer", true},
		{"u64-slice", "integer", true},
		{"float-slice", "number", true},
		{"f32-slice", "number", true},
		{"f64-slice", "number", true},
	}

	flagMap := make(map[string]spec.FlagItem)
	for _, f := range doc.Commands.Flags {
		flagMap[f.Name] = f
	}

	for _, tc := range tests {
		f, ok := flagMap[tc.name]
		if !ok {
			t.Fatalf("missing flag %q", tc.name)
		}
		if f.Type != tc.typ {
			t.Errorf("flag %s: expected type %q, got %q", tc.name, tc.typ, f.Type)
		}
		if f.Variadic != tc.variadic {
			t.Errorf("flag %s: expected to be variadic but was not", tc.name)
		}
	}
}

func TestGenerateDocument_YAMLOutput(t *testing.T) {
	root := &cli.Command{
		Name:  "myapp",
		Usage: "A test CLI",
	}
	buf := new(bytes.Buffer)
	doc := GenerateDocument(root, WithOutput(buf))

	if doc == nil {
		t.Fatal("document should not be nil")
	}
}

func TestFromCommand_AttachesHiddenSubcommand(t *testing.T) {
	root := &cli.Command{Name: "myapp"}
	FromCommand(root)

	found := false
	for _, c := range root.Commands {
		if c.Name == "__opencli" {
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

func TestGetBinaryName(t *testing.T) {
	cmd := &cli.Command{Name: "myapp"}
	if GetBinaryName(cmd) != "myapp" {
		t.Errorf("expected 'myapp', got %q", GetBinaryName(cmd))
	}
}

func TestBuildCommandLine_Nested(t *testing.T) {
	root := &cli.Command{Name: "ocli"}
	gen := &cli.Command{Name: "gen"}
	docs := &cli.Command{
		Name: "docs",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
	root.Commands = append(root.Commands, gen)
	gen.Commands = append(gen.Commands, docs)

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
