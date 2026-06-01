package cliutils

import "github.com/spf13/cobra"

const (
	ExitCodeOK              int = 0
	ExitCodeInternalErr     int = 1
	ExitCodeBadUserInputErr int = 2
)

type CLIError struct {
	Code    int
	Message string
	Command *cobra.Command
	ArgsDef map[string]string // <arg-name>: arg description key/value pairs
	UseLine string
}

func (err CLIError) Error() string {
	return err.Message
}

// ValidationError represents an error that occurred during action validation
type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func InternalError(cmd *cobra.Command, argsDef map[string]string, useLine, message string) *CLIError {
	return &CLIError{
		Code:    ExitCodeInternalErr,
		Message: message,
		Command: cmd,
		ArgsDef: argsDef,
		UseLine: useLine,
	}
}

func BadUserInput(cmd *cobra.Command, argsDef map[string]string, useLine, message string) *CLIError {
	return &CLIError{
		Code:    ExitCodeBadUserInputErr,
		Message: message,
		Command: cmd,
		ArgsDef: argsDef,
		UseLine: useLine,
	}
}
