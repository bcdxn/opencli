package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcdxn/opencli/internal/cli/gencli"
)

// MockFileWriter implements gencli.FileWriter for testing
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
	ios, _, output, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

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
	ios, _, output, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

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
	ios, _, output, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

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
	ios, _, output, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

	// Verify - should return error
	if err == nil {
		t.Errorf("expected error for invalid spec with FailOnErr=true, got nil")
	}

	_, ok := err.(*gencli.ValidationError)
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
	ios, _, output, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := gencli.OcliCheckFlags{FailOnErr: false}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

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
	ios, _, _, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: "/nonexistent/path/to/file.json"}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err := actions.OcliCheck(t.Context(), args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for nonexistent file, got nil")
	}

	cliErr, ok := err.(*gencli.ValidationError)
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
	ios, _, _, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpdir}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for directory path, got nil")
	}

	cliErr, ok := err.(*gencli.ValidationError)
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
	ios, _, _, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for unsupported format, got nil")
	}

	cliErr, ok := err.(*gencli.ValidationError)
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
	ios, _, _, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: filePath}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

	// Verify
	if err == nil {
		t.Errorf("expected error for permission denied, got nil")
	}

	_, ok := err.(*gencli.ValidationError)
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
	ios, _, output, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliCheckArgs{PathToSpec: tmpfile.Name()}
	flags := gencli.OcliCheckFlags{FailOnErr: true}

	// Execute
	err = actions.OcliCheck(t.Context(), args, flags)

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

func TestOcliGenDocsHTMLComponentWritesJSAsset(t *testing.T) {
	specDir := t.TempDir()
	specPath := filepath.Join(specDir, "test-cli.ocs.yaml")
	if err := os.WriteFile(specPath, []byte(validSpecYAML), 0644); err != nil {
		t.Fatalf("failed to write temp spec file: %v", err)
	}

	outDir := filepath.Join(specDir, "docs")
	ios, _, _, _ := gencli.TestIOS()
	actions := Actions{IOS: ios}

	args := gencli.OcliGenDocsArgs{PathToSpec: specPath}
	flags := gencli.OcliGenDocsFlags{
		Format: "html-embed",
		Out:    outDir,
	}

	err := actions.OcliGenDocs(t.Context(), args, flags)
	if err != nil {
		t.Fatalf("unexpected error generating html embed docs: %v", err)
	}

	assetPath := filepath.Join(outDir, "ocli-docs.js")
	asset, err := os.ReadFile(assetPath)
	if err != nil {
		t.Fatalf("expected embed asset %q to be created: %v", assetPath, err)
	}

	if !contains(string(asset), "global.OcliDocs = OcliDocs") {
		t.Fatalf("expected embed asset to expose OcliDocs initializer")
	}

	legacyHTMLPath := filepath.Join(outDir, "test-cli.ocs.html")
	if _, err := os.Stat(legacyHTMLPath); err == nil {
		t.Fatalf("did not expect legacy embed html output at %q", legacyHTMLPath)
	}
}

// Helper function to check if a string contains a substring
func contains(str, substr string) bool {
	return bytes.Contains([]byte(str), []byte(substr))
}
