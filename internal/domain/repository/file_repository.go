package repository

import (
	"os"
)

type FileRepository interface {
	ExistsFile(pathFile string) bool
	ExistsDirectory(pathDirectory string) bool
	DeleteFile(pathFile string) error
	WriteFile(pathFile string, content string) error
	OpenFile(filePath string) (*os.File, error)
	CreateFile(filePath string) (*os.File, error)
	GetAbsolutePath(relativePath string) (string, error)
	GetUserHomeDirectory() (string, error)
	CreateDirectoryAll(path string, perm os.FileMode) error
}
