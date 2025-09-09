package service

import (
	"encoding/gob"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/context/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

const CONTEXT_FILE_NAME = "context.gob"

type FileRepository struct{}

func NewFileRepository() port.Repository {
	return &FileRepository{}
}

func (fr *FileRepository) Load(projectName string) (service.Context, error) {
	filePath, err := fr.getFilePath(projectName)
	if err != nil {
		return service.NewDataContext(), err
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return service.NewDataContext(), nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return service.NewDataContext(), err
	}
	defer file.Close()

	var data map[string]string
	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return service.NewDataContext(), err
	}

	context := service.NewDataContext()
	if data != nil {
		context.SetAll(data)
	}

	return context, nil
}

func (fr *FileRepository) Save(projectName string, data service.Context) error {
	filePath, err := fr.getFilePath(projectName)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data.GetAll()); err != nil {
		return err
	}

	return nil
}

func (fr *FileRepository) Exists(projectName string) (bool, error) {
	filePath, err := fr.getFilePath(projectName)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (pr *FileRepository) getFilePath(projectName string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	directoryPath := filepath.Join(currentUser.HomeDir, constants.FastDeployDir)

	return filepath.Join(directoryPath, projectName, CONTEXT_FILE_NAME), nil
}
