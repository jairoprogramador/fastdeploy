package file

import (
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"os"
	"path/filepath"
)

const (
	dirPermission  = 0755
	filePermission = 0644
)

type FilePort interface {
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

type fileAdapter struct {
	fileLogger *logger.FileLogger
}

func NewFileAdapter(fileLogger *logger.FileLogger) FilePort {
	return &fileAdapter{
		fileLogger: fileLogger,
	}
}

func (fc *fileAdapter) WriteFile(filePath string, content string) error {
	file, err := fc.CreateFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fc.logError(err)
	}
	return err
}

func (fc *fileAdapter) CreateFile(filePath string) (*os.File, error) {
	dirPath := filepath.Dir(filePath)

	exists, err := fc.ExistsDirectory(dirPath)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := os.MkdirAll(dirPath, dirPermission); err != nil {
			fc.logError(err)
			return nil, err
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filePermission)
	if err != nil {
		fc.logError(err)
	}
	return file, err
}

func (fc *fileAdapter) OpenFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fc.logError(err)
	}
	return file, err
}

func (fc *fileAdapter) ExistsFile(filePath string) (bool, error) {
	isDirectory, err := fc.ExistsDirectory(filePath)
	if err != nil {
		return false, err
	}
	return !isDirectory, nil
}

func (fc *fileAdapter) ExistsDirectory(dirPath string) (bool, error) {
	return fc.isPathDirectory(dirPath)
}

func (fc *fileAdapter) DeleteFile(filePath string) error {
	exists, err := fc.ExistsFile(filePath)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	err = os.Remove(filePath)
	if err != nil {
		fc.logError(err)
	}
	return err
}

func (fc *fileAdapter) GetPathAbsolute(relativePath string) (string, error) {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		fc.logError(err)
	}
	return path, err
}

func (fc *fileAdapter) GetUserHomeDirectory() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		fc.logError(err)
	}
	return dir, err
}

func (fc *fileAdapter) ReadDirectory(dirPath string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fc.logError(err)
	}
	return entries, err
}

func (fc *fileAdapter) GetPath(paths ...string) string {
	return filepath.Join(paths...)
}

func (fc *fileAdapter) isPathDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		fc.logError(err)
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func (fc *fileAdapter) logError(err error) {
	if err != nil {
		fc.fileLogger.Error(err)
	}
}
