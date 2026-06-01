package codec

import "github.com/bcdxn/opencli/spec"

// rawDocument represents an unmarshalled OpenCLI document without any post-processing/indexing.
// It is used as an intermediate step to building the indexed spec.Document structure.
type rawDocument struct {
	OpenCLIVersion       string                    `json:"opencliVersion" yaml:"opencliVersion"`
	Info                 spec.Info                 `json:"info" yaml:"info"`
	Install              []spec.InstallMethod      `json:"install,omitempty" yaml:"install,omitempty"`
	Global               *spec.Global              `json:"global,omitempty" yaml:"global,omitempty"`
	Commands             map[string]rawCommandItem `json:"commands,omitempty" yaml:"commands,omitempty"`
	MemoizedCommandLines []string
}

type rawCommandItem struct {
	Summary     string              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Aliases     []string            `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Args        []spec.ArgumentItem `json:"args,omitempty" yaml:"args,omitempty"`
	Flags       []spec.FlagItem     `json:"flags,omitempty" yaml:"flags,omitempty"`
	Hidden      bool                `json:"hidden,omitempty" yaml:"hidden,omitempty"`
	Group       bool                `json:"group,omitempty" yaml:"group,omitempty"`
	ExitCodes   []spec.ExitCode     `json:"exitCodes,omitempty" yaml:"exitCodes,omitempty"`
}
