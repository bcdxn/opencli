package app

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
)

// MockFileWriter implements cliutils.FileWriter for testing
type MockFileWriter struct {
	*bytes.Buffer
}

func (m *MockFileWriter) Fd() uintptr {
	return 1 // stdout file descriptor
}

// validSpecJSON is a minimal valid OpenCLI specification in JSON format
const validSpecJSON = `{
  "opencliVersion": "1.0.0-alpha.8",
  "info": {
    "title": "Test CLI",
    "binary": "test",
    "version": "1.0.0"
  },
  "global": {
    "config": {
      "json": "~/.test/config.json"
    }
  },
  "commands": {
    "test": {
      "summary": "Test command"
    }
  }
}`

// validSpecYAML is a minimal valid OpenCLI specification in YAML format
const validSpecYAML = `opencliVersion: 1.0.0-alpha.8
info:
  title: Test CLI
  binary: test
  version: 1.0.0
global:
  config:
    json: ~/.test/config.json
commands:
  test:
    summary: Test command
`

// invalidSpec is an invalid OpenCLI specification (missing required fields)
const invalidSpec = `{
  "opencliVersion": "1.0.0-alpha.8"
}`

func TestOcliCheckValidJSON(t *testing.T) {
	// Create temporary file with valid JSON spec
	tmpfile, err := os.CreateTemp("", "valid-spec-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(validSpecJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify
	if err != nil {
		t.Errorf("expected no error for valid JSON spec, got: %v", err)
	}

	outputStr := output.String()
	if !contains(outputStr, "✓ Checking") {
		t.Errorf("expected success indicator in output, got: %s", outputStr)
	}
	if !contains(outputStr, "✓ Document is valid") {
		t.Errorf("expected valid document message, got: %s", outputStr)
	}
	if !contains(outputStr, "Format: json") {
		t.Errorf("expected JSON format in output, got: %s", outputStr)
	}
}

func TestOcliCheckValidYAML(t *testing.T) {
	// Create temporary file with valid YAML spec
	tmpfile, err := os.CreateTemp("", "valid-spec-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(validSpecYAML); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify
	if err != nil {
		t.Errorf("expected no error for valid YAML spec, got: %v", err)
	}

	outputStr := output.String()
	if !contains(outputStr, "✓ Checking") {
		t.Errorf("expected success indicator in output, got: %s", outputStr)
	}
	if !contains(outputStr, "✓ Document is valid") {
		t.Errorf("expected valid document message, got: %s", outputStr)
	}
	if !contains(outputStr, "Format: yaml") {
		t.Errorf("expected YAML format in output, got: %s", outputStr)
	}
}

func TestOcliCheckValidYML(t *testing.T) {
	// Create temporary file with valid YML spec (alternate extension)
	tmpfile, err := os.CreateTemp("", "valid-spec-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(validSpecYAML); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify
	if err != nil {
		t.Errorf("expected no error for valid YML spec, got: %v", err)
	}

	outputStr := output.String()
	if !contains(outputStr, "Format: yml") {
		t.Errorf("expected YML format in output, got: %s", outputStr)
	}
}

func TestOcliCheckInvalidSpecFailOnErr(t *testing.T) {
	// Create temporary file with invalid spec
	tmpfile, err := os.CreateTemp("", "invalid-spec-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(invalidSpec); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify - should return error
	if err == nil {
		t.Errorf("expected error for invalid spec with FailOnErr=true, got nil")
	}

	_, ok := err.(*cliutils.ValidationError)
	if !ok {
		t.Errorf("expected CLIError, got %T", err)
	}

	outputStr := output.String()
	if !contains(outputStr, "✗ Validation failed") {
		t.Errorf("expected validation failed message, got: %s", outputStr)
	}
}

func TestOcliCheckInvalidSpecNoFailOnErr(t *testing.T) {
	// Create temporary file with invalid spec
	tmpfile, err := os.CreateTemp("", "invalid-spec-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(invalidSpec); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := cliutils.OcliCheckFlags{FailOnErr: false}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify - should NOT return error
	if err != nil {
		t.Errorf("expected no error for invalid spec with FailOnErr=false, got: %v", err)
	}

	outputStr := output.String()
	if !contains(outputStr, "✗ Validation failed") {
		t.Errorf("expected validation failed message, got: %s", outputStr)
	}
}

func TestOcliCheckFileNotFound(t *testing.T) {
	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: "/nonexistent/path/to/file.json"}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err := actions.OcliCheck(args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for nonexistent file, got nil")
	}

	cliErr, ok := err.(*cliutils.ValidationError)
	if !ok {
		t.Errorf("expected CLIError, got %T", err)
	}

	if !contains(cliErr.Message, "file not found") {
		t.Errorf("expected 'file not found' in error message, got: %s", cliErr.Message)
	}
}

func TestOcliCheckDirectoryPath(t *testing.T) {
	// Create temporary directory
	tmpdir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpdir}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for directory path, got nil")
	}

	cliErr, ok := err.(*cliutils.ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	if !contains(cliErr.Message, "is a directory") {
		t.Errorf("expected 'is a directory' in error message, got: %s", cliErr.Message)
	}
}

func TestOcliCheckUnsupportedFormat(t *testing.T) {
	// Create temporary file with unsupported extension
	tmpfile, err := os.CreateTemp("", "spec-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(validSpecJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for unsupported format, got nil")
	}

	cliErr, ok := err.(*cliutils.ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	if !contains(cliErr.Message, "unsupported file format") {
		t.Errorf("expected 'unsupported file format' in error message, got: %s", cliErr.Message)
	}
}

func TestOcliCheckPermissionDenied(t *testing.T) {
	// Create temporary file
	tmpfile, err := os.CreateTemp("", "spec-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	filePath := tmpfile.Name()

	if _, err := tmpfile.WriteString(validSpecJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Remove read permissions
	if err := os.Chmod(filePath, 0000); err != nil {
		t.Fatalf("failed to change file permissions: %v", err)
	}
	defer func() {
		os.Chmod(filePath, 0644) // Restore permissions for cleanup
		os.Remove(filePath)
	}()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: filePath}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for permission denied, got nil")
	}

	_, ok := err.(*cliutils.ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestOcliCheckOutputFormatting(t *testing.T) {
	// Create temporary file with valid spec
	tmpfile, err := os.CreateTemp("", "spec-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(validSpecJSON); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	// Setup
	output := &MockFileWriter{Buffer: &bytes.Buffer{}}
	ios := &cliutils.IOStreams{Out: output}
	actions := Actions{IOS: ios}

	args := cliutils.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := cliutils.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(args, flags)

	// Verify output format
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	outputStr := output.String()

	// Check for proper formatting
	if !contains(outputStr, "✓ Checking") {
		t.Error("missing checking indicator")
	}
	if !contains(outputStr, "Format:") {
		t.Error("missing format line")
	}
	if !contains(outputStr, "✓ Document is valid") {
		t.Error("missing valid document message")
	}

	// Verify file path is included
	filename := filepath.Base(tmpfile.Name())
	if !contains(outputStr, filename) {
		t.Errorf("expected filename '%s' in output, got: %s", filename, outputStr)
	}
}

// Helper function to check if a string contains a substring
func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}
