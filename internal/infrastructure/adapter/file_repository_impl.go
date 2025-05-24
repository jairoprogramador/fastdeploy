package adapter

import (
	"os"
	"path/filepath"
)

// FileRepository defines the interface for file system operations
type FileRepository interface {
	ExistsFile(pathFile string) (bool, error)
	DeleteFile(pathFile string) error
	WriteFile(pathFile string, content string) error
	OpenFile(pathFile string) (*os.File, error)
	CreateFile(pathFile string) (*os.File, error)

	ExistsDirectory(pathDirectory string) (bool, error)
	GetUserHomeDirectory() (string, error)
	ReadDirectory(pathDirectory string) ([]os.DirEntry, error)

	GetPath(paths ...string) string
	GetPathAbsolute(pathRelative string) (string, error)
}

type fileRepositoryImpl struct{}

func NewFileRepositoryImpl() FileRepository {
	return &fileRepositoryImpl{}
}

func (st *fileRepositoryImpl) WriteFile(pathFile string, content string) error {
	file, err := st.CreateFile(pathFile)
	if err == nil {
		_, err = file.WriteString(content)
	}
	defer file.Close()
	return err
}

func (st *fileRepositoryImpl) CreateFile(pathFile string) (*os.File, error) {
	pathDirectory := filepath.Dir(pathFile)
	existsDirectory, err := st.ExistsDirectory(pathDirectory)
	if err != nil {
		return nil, err
	}
	if !existsDirectory {
		if err := os.MkdirAll(pathDirectory, 0755); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(pathFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

func (st *fileRepositoryImpl) OpenFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

func (st *fileRepositoryImpl) ExistsFile(pathFile string) (bool, error) {
	isPathDirectory, err := st.isPathDirectory(pathFile)
	if err != nil {
		return false, err
	}
	return !isPathDirectory, nil
}

func (st *fileRepositoryImpl) ExistsDirectory(pathDirectory string) (bool, error) {
	return st.isPathDirectory(pathDirectory)
}

func (st *fileRepositoryImpl) DeleteFile(pathFile string) error {
	var err error
	var existsFile bool

	if existsFile, err = st.ExistsFile(pathFile); err == nil && existsFile {
		if err = os.Remove(pathFile); err != nil {
			return err
		}
	}
	return err
}

func (st *fileRepositoryImpl) GetPathAbsolute(relativePath string) (string, error) {
	return filepath.Abs(relativePath)
}

func (st *fileRepositoryImpl) GetUserHomeDirectory() (string, error) {
	return os.UserHomeDir()
}

func (st *fileRepositoryImpl) ReadDirectory(pathDirectory string) ([]os.DirEntry, error) {
	return os.ReadDir(pathDirectory)
}

func (st *fileRepositoryImpl) GetPath(paths ...string) string {
	return filepath.Join(paths...)
}

func (st *fileRepositoryImpl) isPathDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fileInfo.IsDir(), nil
}
