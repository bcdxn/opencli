package spec_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bcdxn/opencli/spec"
	"gopkg.in/yaml.v3"
)

type schemaVersionDocument struct {
	ID         string `json:"$id"`
	Properties struct {
		OpenCLIVersion struct {
			Enum []string `json:"enum"`
		} `json:"opencliVersion"`
	} `json:"properties"`
}

func TestSchemaVersionArtifacts(t *testing.T) {
	rootSchema, err := os.ReadFile("../spec.schema.json")
	if err != nil {
		t.Fatal(err)
	}

	var document schemaVersionDocument
	if err := json.Unmarshal(rootSchema, &document); err != nil {
		t.Fatal(err)
	}
	if len(document.Properties.OpenCLIVersion.Enum) != 1 {
		t.Fatalf("schema opencliVersion enum has %d values, want exactly one", len(document.Properties.OpenCLIVersion.Enum))
	}
	if got := document.Properties.OpenCLIVersion.Enum[0]; got != spec.SchemaVersion {
		t.Fatalf("generated SchemaVersion is %q, schema declares %q", spec.SchemaVersion, got)
	}

	for _, relativePath := range []string{"../validate/spec.schema.json", "../web/src/spec.schema.json"} {
		copy, err := os.ReadFile(relativePath)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(copy, rootSchema) {
			t.Errorf("%s is not synchronized with ../spec.schema.json", relativePath)
		}
	}
}

func TestOpenCLIVersionDocuments(t *testing.T) {
	var mismatches []string
	root := filepath.Clean("..")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".next", ".scratch", "build", "dist", "node_modules", "out", "scratch":
				return filepath.SkipDir
			}
			return nil
		}

		extension := strings.ToLower(filepath.Ext(path))
		if extension != ".json" && extension != ".yaml" && extension != ".yml" {
			return nil
		}

		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if !bytes.Contains(contents, []byte("opencliVersion")) {
			return nil
		}

		var document map[string]any
		if extension == ".json" {
			err = json.Unmarshal(contents, &document)
		} else {
			err = yaml.Unmarshal(contents, &document)
		}
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}

		version, ok := document["opencliVersion"].(string)
		if ok && version != spec.SchemaVersion {
			relativePath, _ := filepath.Rel(root, path)
			mismatches = append(mismatches, fmt.Sprintf("%s declares %q", relativePath, version))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(mismatches) > 0 {
		t.Fatalf("documents do not use schema version %q:\n%s", spec.SchemaVersion, strings.Join(mismatches, "\n"))
	}
}
