package adapter

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model/logger"
	"os"
	"path/filepath"
)

const (
	dirPermission  = 0755
	filePermission = 0644
)

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

type osFileController struct {
	fileLogger *logger.FileLogger
}

func NewOsFileController(fileLogger *logger.FileLogger) FileController {
	return &osFileController{
		fileLogger: fileLogger,
	}
}

func (fc *osFileController) WriteFile(filePath string, content string) error {
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

func (fc *osFileController) CreateFile(filePath string) (*os.File, error) {
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

func (fc *osFileController) OpenFile(filePath string) (*os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fc.logError(err)
	}
	return file, err
}

func (fc *osFileController) ExistsFile(filePath string) (bool, error) {
	isDirectory, err := fc.ExistsDirectory(filePath)
	if err != nil {
		return false, err
	}
	return !isDirectory, nil
}

func (fc *osFileController) ExistsDirectory(dirPath string) (bool, error) {
	return fc.isPathDirectory(dirPath)
}

func (fc *osFileController) DeleteFile(filePath string) error {
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

func (fc *osFileController) GetPathAbsolute(relativePath string) (string, error) {
	path, err := filepath.Abs(relativePath)
	if err != nil {
		fc.logError(err)
	}
	return path, err
}

func (fc *osFileController) GetUserHomeDirectory() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		fc.logError(err)
	}
	return dir, err
}

func (fc *osFileController) ReadDirectory(dirPath string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fc.logError(err)
	}
	return entries, err
}

func (fc *osFileController) GetPath(paths ...string) string {
	return filepath.Join(paths...)
}

func (fc *osFileController) isPathDirectory(path string) (bool, error) {
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

func (fc *osFileController) logError(err error) {
	if err != nil {
		fc.fileLogger.Error(err)
	}
}
