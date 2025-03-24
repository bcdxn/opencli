package oclispec

import (
	"testing"
)

// TestValidateDocumentJSON tests the ValidateDocument function by passing different documents in
// and validating those documents against all supported versions of the specification.
func TestValidateDocumentJSON(t *testing.T) {
	tests := []struct {
		name           string
		opencliVersion string
		document       []byte
		wantErr        bool
	}{
		{
			name: "valid document",
			document: []byte(`{
				"opencliVersion": "1.0.0-alpha.0",
				"info": {
					"binary": "test",
					"title": "Test OpenCLI Specification",
					"version": "1.0.0"
				},
				"commands": {}
			}`),
			wantErr: false,
		},
		{
			name: "unsupported version",
			document: []byte(`{
				"opencliVersion": "0.0.5",
				"info": {
					"binary": "test",
					"title": "Test OpenCLI Specification",
					"version": "1.0.0"
				},
				"commands": {}
			}`),
			wantErr: true,
		},
		{
			name:     "malformed document",
			document: []byte(`{"invalid": "d`),
			wantErr:  true,
		},
		{
			name:     "invalid document",
			document: []byte(`{"opencliVersion":"1.0.0-alpha.0","key": "value"}`),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDocumentJSON(tt.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDocumentJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateDocumentYAML tests the ValidateDocument function by passing different documents in
// and validating those documents against all supported versions of the specification.
func TestValidateDocumentYAML(t *testing.T) {
	tests := []struct {
		name           string
		opencliVersion string
		document       []byte
		wantErr        bool
	}{
		{
			name: "valid document",
			document: []byte(`
opencliVersion: "1.0.0-alpha.0"
info:
  binary: "test"	
  title: "Test OpenCLI Specification"
  version: "1.0.0"
commands: {}
      `),
			wantErr: false,
		},
		{
			name: "unsupported version",
			document: []byte(`
opencliVersion: "0.0.5"
info:
  binary: "test"
  title: "Test OpenCLI Specification"
  version: "1.0.0"
commands: {}
      `),
			wantErr: true,
		},
		{
			name:     "malformed document",
			document: []byte(`{"invalid": "d`),
			wantErr:  true,
		},
		{
			name: "invalid document",
			document: []byte(`opencliVersion: "1.0.0-alpha.0"
key: "value"
      `),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDocumentYAML(tt.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDocumentYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestVersions tests the Versions function by comparing the expected versions with the actual versions.
func TestVersions(t *testing.T) {
	expectedVersions := []string{"1.0.0-alpha.0"}
	actualVersions := Versions()

	if len(actualVersions) != len(expectedVersions) {
		t.Errorf("Expected %d versions, got %d", len(expectedVersions), len(actualVersions))
		return
	}

	for i, version := range actualVersions {
		if version != expectedVersions[i] {
			t.Errorf("Expected version %s, got %s", expectedVersions[i], version)
		}
	}
}
