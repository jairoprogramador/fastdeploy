package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"os"
	"path/filepath"
)

type fileRepositoryImpl struct{}

func NewFileRepositoryImpl() repository.FileRepository {
	return &fileRepositoryImpl{}
}

func (st *fileRepositoryImpl) WriteFile(pathFile string, content string) error {
	return filesystem.WriteFile(pathFile, content)
}

func (st *fileRepositoryImpl) CreateFile(filePath string) (*os.File, error) {
	return filesystem.CreateFile(filePath)
}

func (st *fileRepositoryImpl) OpenFile(filePath string) (*os.File, error) {
	return filesystem.OpenFile(filePath)
}

func (st *fileRepositoryImpl) ExistsFile(pathFile string) bool {
	exists, _ := filesystem.ExistsFile(pathFile)
	return exists
}

func (st *fileRepositoryImpl) ExistsDirectory(pathDirectory string) bool {
	exists, _ := filesystem.ExistsDirectory(pathDirectory)
	return exists
}

func (st *fileRepositoryImpl) DeleteFile(path string) error {
	return filesystem.RemoveFile(path)
}

func (st *fileRepositoryImpl) GetAbsolutePath(relativePath string) (string, error) {
	absPath, err := filepath.Abs(relativePath)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func (st *fileRepositoryImpl) GetUserHomeDirectory() (string, error) {
	return os.UserHomeDir()
}

func (st *fileRepositoryImpl) CreateDirectoryAll(path string, perm os.FileMode) error {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			return nil
		}
		return &os.PathError{Op: "mkdirall", Path: path, Err: os.ErrExist}
	}
	if !os.IsNotExist(err) {
		return err
	}
	return os.MkdirAll(path, perm)
}
