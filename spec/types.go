package spec

// Document represents the top-level OpenCLI document.
type Document struct {
	OpenCLIVersion string
	Info           Info
	Install        []InstallMethod
	Global         *Global
	Commands       *CommandItem
}

type Info struct {
	Title       string  `json:"title" yaml:"title"`
	Summary     string  `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	License     License `json:"license,omitempty" yaml:"license,omitempty"`
	Contact     Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	Binary      string  `json:"binary" yaml:"binary"`
	Version     string  `json:"version" yaml:"version"`
}

type License struct {
	Name   string `json:"name" yaml:"name"`
	SpdxID string `json:"spdxId,omitempty" yaml:"spdxId,omitempty"`
	URL    string `json:"url,omitempty" yaml:"url,omitempty"`
}

type Contact struct {
	Name  string `json:"name" yaml:"name"`
	Email string `json:"email" yaml:"email"`
	URL   string `json:"url" yaml:"url"`
}

type InstallMethod struct {
	Name        string `json:"name" yaml:"name"`
	Command     string `json:"command,omitempty" yaml:"command,omitempty"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type Global struct {
	ExitCodes   []ExitCode  `json:"exitCodes,omitempty" yaml:"exitCodes,omitempty"`
	ConfigFiles ConfigFiles `json:"configFiles" yaml:"configFiles"`
	Flags       []FlagItem  `json:"flags,omitempty" yaml:"flags,omitempty"`
}

type ExitCode struct {
	Code    int    `json:"code" yaml:"code"`
	Status  string `json:"status" yaml:"status"`
	Summary string `json:"summary" yaml:"summary"`
}

type ConfigFiles struct {
	Json string `json:"json,omitempty" yaml:"json,omitempty"`
	Toml string `json:"toml,omitempty" yaml:"toml,omitempty"`
	Yaml string `json:"yaml,omitempty" yaml:"yaml,omitempty"`
}

type CommandItem struct {
	Segment        string // the single command segment, e.g.: petstore pets add --> `add`
	CommandLineRaw string // the literal command line string defined in the spec-compliant document
	CommandLine    string // the parsed command line without {commands}, <arguments>, [flags], etc.
	Summary        string
	Description    string
	Aliases        []string
	Args           []ArgumentItem
	Flags          []FlagItem
	Hidden         bool
	Group          bool
	ExitCodes      []ExitCode
	Commands       []*CommandItem
	// Properties set during post processing of unmarshalling
	Children         bool
	ChildrenArgs     bool
	ChildrenFlags    bool
	CommandModifiers []string
	ArgsModifiers    []string
	FlagModifiers    []string
}

type ArgumentItem struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type,omitempty" yaml:"type,omitempty"`
	Variadic bool   `json:"variadic,omitempty" yaml:"variadic,omitempty"`
	Summary  string `json:"summary,omitempty" yaml:"summary,omitempty"`
	Required bool   `json:"required,omitempty" yaml:"required,omitempty"`
}

type FlagItem struct {
	Name    string      `json:"name" yaml:"name"`
	Aliases []string    `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Type    string      `json:"type" yaml:"type"`
	Summary string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Default interface{} `json:"default,omitempty" yaml:"default,omitempty"`
	Hidden  bool        `json:"hidden,omitempty" yaml:"hidden,omitempty"`
}
