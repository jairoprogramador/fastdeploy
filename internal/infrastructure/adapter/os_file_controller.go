package adapter

import (
	"os"
	"path/filepath"
)

// File permission constants
const (
	dirPermission  = 0755
	filePermission = 0644
)

// FileController defines the interface for file system operations
type FileController interface {
	ExistsFile(filePath string) (bool, error)
	DeleteFile(filePath string) error
	WriteFile(filePath string, content string) error
	OpenFile(filePath string) (*os.File, error)
	CreateFile(filePath string) (*os.File, error)

	ExistsDirectory(dirPath string) (bool, error)
	GetUserHomeDirectory() (string, error)
	ReadDirectory(dirPath string) ([]os.DirEntry, error)

	GetPath(paths ...string) string
	GetPathAbsolute(relativePath string) (string, error)
}

type osFileController struct{}

// NewOsFileController creates a new file controller instance
func NewOsFileController() FileController {
	return &osFileController{}
}

// WriteFile writes content to a file, creating it if it doesn't exist
func (fc *osFileController) WriteFile(filePath string, content string) error {
	file, err := fc.CreateFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

// CreateFile creates a new file, ensuring its directory exists
func (fc *osFileController) CreateFile(filePath string) (*os.File, error) {
	dirPath := filepath.Dir(filePath)

	// Ensure directory exists
	exists, err := fc.ExistsDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := os.MkdirAll(dirPath, dirPermission); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePermission)
}

// OpenFile opens an existing file for reading
func (fc *osFileController) OpenFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

// ExistsFile checks if a file exists at the given path
func (fc *osFileController) ExistsFile(filePath string) (bool, error) {
	isDirectory, err := fc.isPathDirectory(filePath)
	if err != nil {
		return false, err
	}
	return !isDirectory, nil
}

// ExistsDirectory checks if a directory exists at the given path
func (fc *osFileController) ExistsDirectory(dirPath string) (bool, error) {
	return fc.isPathDirectory(dirPath)
}

// DeleteFile removes a file if it exists
func (fc *osFileController) DeleteFile(filePath string) error {
	exists, err := fc.ExistsFile(filePath)
	if err != nil {
		return err
	}

	if !exists {
		return nil // File doesn't exist, nothing to delete
	}

	return os.Remove(filePath)
}

// GetPathAbsolute converts a relative path to an absolute path
func (fc *osFileController) GetPathAbsolute(relativePath string) (string, error) {
	return filepath.Abs(relativePath)
}

// GetUserHomeDirectory returns the current user's home directory
func (fc *osFileController) GetUserHomeDirectory() (string, error) {
	return os.UserHomeDir()
}

// ReadDirectory returns the contents of a directory
func (fc *osFileController) ReadDirectory(dirPath string) ([]os.DirEntry, error) {
	return os.ReadDir(dirPath)
}

// GetPath joins path elements into a single path
func (fc *osFileController) GetPath(paths ...string) string {
	return filepath.Join(paths...)
}

// isPathDirectory checks if a path points to a directory
func (fc *osFileController) isPathDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fileInfo.IsDir(), nil
}
