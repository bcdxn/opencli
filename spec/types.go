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
	Title       string   `json:"title" yaml:"title"`
	Summary     string   `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	License     *License `json:"license,omitempty" yaml:"license,omitempty"`
	Contact     *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	Binary      string   `json:"binary" yaml:"binary"`
	Version     string   `json:"version" yaml:"version"`
}

type License struct {
	Name   string `json:"name" yaml:"name"`
	SpdxID string `json:"spdxId,omitempty" yaml:"spdxId,omitempty"`
	URL    string `json:"url,omitempty" yaml:"url,omitempty"`
}

type Contact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
}

type InstallMethod struct {
	Name        string `json:"name" yaml:"name"`
	Command     string `json:"command,omitempty" yaml:"command,omitempty"`
	URL         string `json:"url,omitempty" yaml:"url,omitempty"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type Global struct {
	ExitCodes []ExitCode    `json:"exitCodes,omitempty" yaml:"exitCodes,omitempty"`
	Config    Configuration `json:"config" yaml:"config"`
	Flags     []FlagItem    `json:"flags,omitempty" yaml:"flags,omitempty"`
}

type ExitCode struct {
	Code        int    `json:"code" yaml:"code"`
	Status      string `json:"status" yaml:"status"`
	Summary     string `json:"summary" yaml:"summary"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

type Configuration struct {
	JSON string `json:"json,omitempty" yaml:"json,omitempty"`
	TOML string `json:"toml,omitempty" yaml:"toml,omitempty"`
	YAML string `json:"yaml,omitempty" yaml:"yaml,omitempty"`
}

type CommandItem struct {
	Segment        string // the single command segment, e.g.: petstore pets add --> `add`
	Derived        bool   // Command was derived from internal segments when unmarshalling a spec-compliant OpenCLI doc and not explicitly declared
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
	Examples       []Example
	Commands       []*CommandItem
	// Properties set during post processing of unmarshalling
	VisibleArgs          bool
	VisibleFlags         bool
	VisibleChildren      bool
	VisibleChildrenArgs  bool
	VisibleChildrenFlags bool
	CommandModifiers     []string
	ArgsModifiers        []string
	FlagsModifiers       []string
}

type ArgumentItem struct {
	Name        string   `json:"name" yaml:"name"`
	Type        string   `json:"type,omitempty" yaml:"type,omitempty"`
	Variadic    bool     `json:"variadic" yaml:"variadic"`
	MinItems    int      `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems    int      `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	Summary     string   `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool     `json:"required" yaml:"required"`
	Default     any      `json:"default" yaml:"default"`
	Hidden      bool     `json:"hidden" yaml:"hidden"`
	Choices     []Choice `json:"choices,omitempty" yaml:"choices,omitempty"`
}

type FlagItem struct {
	Name        string              `json:"name" yaml:"name"`
	Aliases     []string            `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Type        string              `json:"type" yaml:"type"`
	Variadic    bool                `json:"variadic" yaml:"variadic"`
	MinItems    int                 `json:"minItems,omitempty" yaml:"minItems,omitempty"`
	MaxItems    int                 `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	Summary     string              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool                `json:"required" yaml:"required"`
	Default     any                 `json:"default,omitempty" yaml:"default,omitempty"`
	Hidden      bool                `json:"hidden" yaml:"hidden"`
	Choices     []Choice            `json:"choices,omitempty" yaml:"choices,omitempty"`
	AltSources  []AlternativeSource `json:"alternativeSources,omitempty" yaml:"alternativeSources,omitempty"`
}

type Choice struct {
	// Value can be a string, float64, or bool due to the schema's mixed types.
	Value       any    `json:"value"`
	Description string `json:"description,omitempty"`
}

type AlternativeSource struct {
	Type     string `json:"type" yaml:"type"`
	Property string `json:"property" yaml:"property"`
}

type Example struct {
	Title   string `json:"title" yaml:"title"`
	Content string `json:"content" yaml:"content"`
}
