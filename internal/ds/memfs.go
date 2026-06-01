package ds

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// MemFS satisfies the fs.FS interface
type MemFS map[string][]byte

func (m MemFS) Open(name string) (fs.File, error) {
	// Clean the path relative to the root (io/fs paths cannot start with /)
	cleaned := path.Clean(name)
	if cleaned == "." {
		return nil, fs.ErrNotExist
	}

	content, exists := m[cleaned]
	if !exists {
		return nil, fs.ErrNotExist
	}

	return &MemFile{
		name:   cleaned,
		reader: bytes.NewReader(content),
		size:   int64(len(content)),
	}, nil
}

// Persist writes all in-memory files safely to the specified target directory on disk.
func (m MemFS) Persist(targetDir string) error {
	// Resolve absolute path of target destination to defend against traversal attacks
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for target dir: %w", err)
	}

	// Iterate through all the map files
	for relativePath, content := range m {
		// Clean and normalize the path to remove leading slashes or backtracking elements
		cleanedRel := path.Clean(relativePath)
		cleanedRel = strings.TrimLeft(cleanedRel, "/")

		// Derive final host filepath using platform-native file separators
		localPath, err := filepath.Localize(cleanedRel)
		if err != nil {
			return fmt.Errorf("failed to localize the path for the OS: %w", err)
		}
		finalPath := filepath.Join(absTarget, localPath)

		// Directory Traversal Guardrail: Ensure final path stays inside the target boundary
		if !strings.HasPrefix(finalPath, absTarget+string(filepath.Separator)) && finalPath != absTarget {
			return fmt.Errorf("security block: file path %q attempts to escape target directory", relativePath)
		}

		// Build necessary parent folders (e.g., config/ or logs/) if they do not exist
		parentDir := filepath.Dir(finalPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory tree %q: %w", parentDir, err)
		}

		// Write the actual file data to the local disk
		if err := os.WriteFile(finalPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %q: %w", finalPath, err)
		}
	}

	return nil
}

// MemFile satisfies the fs.File interface
type MemFile struct {
	name   string
	reader *bytes.Reader
	size   int64
}

func (f *MemFile) Stat() (fs.FileInfo, error) { return f, nil }
func (f *MemFile) Read(b []byte) (int, error) { return f.reader.Read(b) }
func (f *MemFile) Close() error               { return nil }

// Dummy methods below satisfy the fs.FileInfo metadata interface
func (f *MemFile) Name() string       { return f.name }
func (f *MemFile) Size() int64        { return f.size }
func (f *MemFile) Mode() fs.FileMode  { return 0444 } // Read-only
func (f *MemFile) ModTime() time.Time { return time.Now() }
func (f *MemFile) IsDir() bool        { return false }
func (f *MemFile) Sys() any           { return nil }
