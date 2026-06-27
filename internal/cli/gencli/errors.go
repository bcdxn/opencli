package gencli

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

func InternalError(cmd *cobra.Command, message string) *CLIError {
	return &CLIError{
		Code:    ExitCodeInternalErr,
		Message: message,
		Command: cmd,
	}
}

func BadUserInput(cmd *cobra.Command, message string) *CLIError {
	return &CLIError{
		Code:    ExitCodeBadUserInputErr,
		Message: message,
		Command: cmd,
	}
}
