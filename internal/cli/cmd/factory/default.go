package cmd

import (
	"github.com/bcdxn/opencli/internal/cli/app"
	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
)

func NewCmdFactory(buildVersion string, ios *cliutils.IOStreams, executablePath string) *cliutils.Factory {

	f := &cliutils.Factory{
		BuildVersion:   buildVersion,
		ExecutablePath: executablePath,
		IOStreams:      ios,
		Actions: app.Actions{
			IOS: ios,
		},
	}

	return f
}
