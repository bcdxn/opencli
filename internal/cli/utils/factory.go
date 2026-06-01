package cliutils

type Factory struct {
	BuildVersion   string
	ExecutablePath string
	IOStreams      *IOStreams
	Actions        CLIActions
}
