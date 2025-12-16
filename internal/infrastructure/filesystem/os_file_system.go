package filesystem

import "os"

// OSFileSystem is a concrete implementation of the FileSystem port that uses
// the standard 'os' package for file operations.
type OSFileSystem struct{}

// NewOSFileSystem creates a new OSFileSystem.
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// ReadFile reads the content of a file at the given path.
func (fs *OSFileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile writes data to a file at the given path with default permissions (0644).
func (fs *OSFileSystem) WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
