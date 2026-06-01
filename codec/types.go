package codec

import (
	"github.com/bcdxn/opencli/internal/ds"
	"github.com/bcdxn/opencli/spec"
)

// rawDocument represents an unmarshalled OpenCLI document without any post-processing/indexing.
// It is used as an intermediate step to building the indexed spec.Document structure.
type rawDocument struct {
	OpenCLIVersion string                          `json:"opencliVersion" yaml:"opencliVersion"`
	Info           spec.Info                       `json:"info" yaml:"info"`
	Install        []spec.InstallMethod            `json:"install,omitempty" yaml:"install,omitempty"`
	Global         *spec.Global                    `json:"global,omitempty" yaml:"global,omitempty"`
	Commands       *ds.Map[string, rawCommandItem] `json:"commands" yaml:"commands"`
}

type rawCommandItem struct {
	Summary     string              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Aliases     []string            `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Args        []spec.ArgumentItem `json:"args,omitempty" yaml:"args,omitempty"`
	Flags       []spec.FlagItem     `json:"flags,omitempty" yaml:"flags,omitempty"`
	Hidden      bool                `json:"hidden" yaml:"hidden"`
	Group       bool                `json:"group" yaml:"group"`
	ExitCodes   []spec.ExitCode     `json:"exitCodes,omitempty" yaml:"exitCodes,omitempty"`
}
