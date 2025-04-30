package oclicode

import (
	"bytes"
	"os"
	"testing"
	"text/template"

	"github.com/bcdxn/opencli/oclispec"
)

func TestTemplatesUrfave(t *testing.T) {
	tests := []struct {
		name     string // name of the test
		template string // path of template file
		expected string // path of expected rendered content
		data     any
	}{
		{
			name:     "Command Name root",
			template: "./templates/urfavecli/cmd_props_name.tmpl",
			expected: "./testdata/urfavecli/cmd_props_name_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "my",
				Command: oclispec.Command{
					Name: "my",
				},
			},
		},
		{
			name:     "Command Name Nested",
			template: "./templates/urfavecli/cmd_props_name.tmpl",
			expected: "./testdata/urfavecli/cmd_props_name_1.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
				},
			},
		},
		{
			name:     "Command Usage Text",
			template: "./templates/urfavecli/cmd_props_usage_text.tmpl",
			expected: "./testdata/urfavecli/cmd_props_usage_text_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
					Line: "my neat nested <arguments> [flags]",
				},
			},
		},
		{
			name:     "Command Description with Summary",
			template: "./templates/urfavecli/cmd_props_description.tmpl",
			expected: "./testdata/urfavecli/cmd_props_description_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name:        "my neat nested",
					Summary:     "A test summary",
					Description: "A test description",
				},
			},
		},
		{
			name:     "Command Description with Description",
			template: "./templates/urfavecli/cmd_props_description.tmpl",
			expected: "./testdata/urfavecli/cmd_props_description_1.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name:        "my neat nested",
					Description: "A test description",
				},
			},
		},
		{
			name:     "Command Description without summary/description",
			template: "./templates/urfavecli/cmd_props_description.tmpl",
			expected: "./testdata/urfavecli/cmd_props_description_2.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
				},
			},
		},
		{
			name:     "Command Aliases",
			template: "./templates/urfavecli/cmd_props_aliases.tmpl",
			expected: "./testdata/urfavecli/cmd_props_aliases_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name:    "my neat nested",
					Aliases: []string{"a", "b", "c"},
				},
			},
		},
		{
			name:     "Command Aliases None",
			template: "./templates/urfavecli/cmd_props_aliases.tmpl",
			expected: "./testdata/urfavecli/cmd_props_aliases_1.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
				},
			},
		},
		// issues with urfave cli argument configs; so we've removed the config generation and corresponding test.
		// {
		// 	name:     "Command Arguments",
		// 	template: "./templates/urfavecli/cmd_props_args.tmpl",
		// 	expected: "./testdata/urfavecli/cmd_props_args_0.txt",
		// 	data: oclispec.CommandTrieNode{
		// 		Name: "nested",
		// 		Command: oclispec.Command{
		// 			Name: "my neat nested",
		// 			Arguments: []oclispec.Argument{
		// 				{
		// 					Name:        "one",
		// 					Type:        "string",
		// 					Summary:     "a test summary",
		// 					Description: "a test description",
		// 				},
		// 				{
		// 					Name:        "two",
		// 					Type:        "string",
		// 					Summary:     "a test summary",
		// 					Description: "a test description",
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		{
			name:     "Command Flags",
			template: "./templates/urfavecli/cmd_props_flags.tmpl",
			expected: "./testdata/urfavecli/cmd_props_flags_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
					Flags: []oclispec.Flag{
						{
							Name:        "one-flag",
							Type:        "string",
							Summary:     "A test summary",
							Description: "A test description",
							Aliases:     []string{"o", "f"},
							Default: oclispec.DefaultValue{
								IsSet:  true,
								String: "default_const",
							},
							AltSources: []oclispec.AlternativeSource{
								{
									Type:                "env",
									EnvironmentVariable: "SOME_VAR",
								},
								{
									Type: "file",
									File: oclispec.FileSource{
										Name:     "yaml",
										Format:   "yaml",
										Path:     "cfg.yaml",
										Property: "some.var",
									},
								},
								{
									Type: "file",
									File: oclispec.FileSource{
										Name:     "json",
										Format:   "json",
										Path:     "cfg.json",
										Property: "some.var",
									},
								},
								{
									Type: "file",
									File: oclispec.FileSource{
										Name:     "toml",
										Format:   "toml",
										Path:     "cfg.toml",
										Property: "some.var",
									},
								},
							},
						},
						{
							Name:        "two-flag",
							Type:        "boolean",
							Description: "A test description",
						},
						{
							Name:     "three-flag",
							Type:     "string",
							Variadic: true,
							Hidden:   true,
						},
					},
				},
			},
		},
		{
			name:     "Command Subcommands",
			template: "./templates/urfavecli/cmd_props_commands.tmpl",
			expected: "./testdata/urfavecli/cmd_props_commands_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
				},
				Commands: []*oclispec.CommandTrieNode{
					{
						Name: "one",
						Command: oclispec.Command{
							Name: "my neat nested one",
						},
					},
					{
						Name: "two",
						Command: oclispec.Command{
							Name: "my neat nested two",
						},
					},
					{
						Name: "three",
						Command: oclispec.Command{
							Name: "my neat nested three",
						},
					},
				},
			},
		},
		{
			name:     "Command Subcommands None",
			template: "./templates/urfavecli/cmd_props_commands.tmpl",
			expected: "./testdata/urfavecli/cmd_props_commands_1.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
				},
			},
		},
		{
			name:     "Command Action Args",
			template: "./templates/urfavecli/cmd_action_args.tmpl",
			expected: "./testdata/urfavecli/cmd_action_args_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
					Arguments: []oclispec.Argument{
						{
							Name:     "first-arg",
							Type:     "string",
							Required: true,
						},
						{
							Name:     "second-arg",
							Type:     "string",
							Required: true,
							Choices: []oclispec.Choice{
								{
									Value: "c1",
								},
								{
									Value: "c2",
								},
								{
									Value: "c3",
								},
							},
						},
						{
							Name:     "third-arg",
							Type:     "string",
							Variadic: true,
							Choices: []oclispec.Choice{
								{
									Value: "c1",
								},
								{
									Value: "c2",
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "Command Name Nested",
			template: "./templates/urfavecli/cmd_action_flags.tmpl",
			expected: "./testdata/urfavecli/cmd_action_flags_0.txt",
			data: oclispec.CommandTrieNode{
				Name: "nested",
				Command: oclispec.Command{
					Name: "my neat nested",
					Flags: []oclispec.Flag{
						{
							Name:     "flag-one",
							Type:     "boolean",
							Required: true,
						},
						{
							Name: "flag-two",
							Type: "string",
							Choices: []oclispec.Choice{
								{
									Value: "c1",
								},
								{
									Value: "c2",
								},
							},
						},
						{
							Name:     "flag-three",
							Type:     "string",
							Variadic: true,
							Required: true,
							Choices: []oclispec.Choice{
								{
									Value: "c1",
								},
								{
									Value: "c2",
								},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpl := MustLoadTemplate(t, test.template)
			expected := MustLoadExpected(t, test.expected)

			rendered := bytes.NewBuffer([]byte{})
			tmpl.Execute(rendered, test.data)

			if rendered.String() != expected {
				t.Errorf("rendered template did not match expected output:\nexpected:\n%s\nactual:\n%s", expected, rendered.String())
			}
		})
	}
}

func MustLoadTemplate(t *testing.T, path string) *template.Template {
	tmpl, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading template file %s - %v", path, err)
	}

	parseTmpl, err := template.New("test_tmpl").Funcs(funcmap()).Parse(string(tmpl))
	if err != nil {
		t.Fatalf("Error parsing template file %s - %v", path, err)
	}
	return parseTmpl
}

func MustLoadExpected(t *testing.T, path string) string {
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("error loading expected content %s - %v", path, err)
	}

	return string(contents)
}
