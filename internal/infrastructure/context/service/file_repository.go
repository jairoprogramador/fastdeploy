package service

import (
	"encoding/gob"
	"fmt"
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

func (fr *FileRepository) Load(projectName, environment string) (service.Context, error) {
	filePath, err := fr.getFilePath(projectName, environment)
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
	environment, err := data.Get(constants.Environment)
	if err != nil {
		return err
	}
	filePath, err := fr.getFilePath(projectName, environment)
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

func (fr *FileRepository) Exists(projectName, environment string) (bool, error) {
	filePath, err := fr.getFilePath(projectName, environment)
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

func (pr *FileRepository) getFilePath(projectName, environment string) (string, error) {
	homeDir, err := pr.getHomeDirPath()
	if err != nil {
		return "", err
	}

	var nameFile = fmt.Sprintf("%s%s", environment, CONTEXT_FILE_NAME)

	return filepath.Join(homeDir, projectName, nameFile), nil
}

func (pr *FileRepository) getHomeDirPath() (string, error) {
	if fastDeployHome := os.Getenv("FASTDEPLOY_HOME"); fastDeployHome != "" {
		return fastDeployHome, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDir), nil
}