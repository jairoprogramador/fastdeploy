package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"os"
	"sync"
)

type fileRepositoryImpl struct {}

var (
	instanceFileRepository     repository.FileRepository
	instanceOnceFileRepository sync.Once
)

func GetFileRepository() repository.FileRepository {
	instanceOnceFileRepository.Do(func() {
		instanceFileRepository = &fileRepositoryImpl{}
	})
	return instanceFileRepository
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
