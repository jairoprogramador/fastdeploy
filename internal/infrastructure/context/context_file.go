package context

import (
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/context/port"
)

var CONTEXT_RETURN_DEFAULT = make(map[string]string)

type ContextFile struct{}

func NewContextFile() port.ContextPort {
	return &ContextFile{}
}

func (fr *ContextFile) Load(pathFileContext string) (map[string]string, error) {
	if _, err := os.Stat(pathFileContext); os.IsNotExist(err) {
		return CONTEXT_RETURN_DEFAULT, nil
	}

	fileContext, err := os.Open(pathFileContext)
	if err != nil {
		return CONTEXT_RETURN_DEFAULT, err
	}
	defer fileContext.Close()

	var contextFound map[string]string
	decoder := gob.NewDecoder(fileContext)

	if err := decoder.Decode(&contextFound); err != nil {
		return CONTEXT_RETURN_DEFAULT, err
	}

	return contextFound, nil
}

func (fr *ContextFile) Save(pathFileContext string, context map[string]string) error {
	pathDirContext := filepath.Dir(pathFileContext)
	if err := os.MkdirAll(pathDirContext, 0755); err != nil {
		return err
	}

	fileContext, err := os.Create(pathFileContext)
	if err != nil {
		return err
	}
	defer fileContext.Close()

	encoderContext := gob.NewEncoder(fileContext)
	if err := encoderContext.Encode(context); err != nil {
		return err
	}

	return nil
}
/*
func (fr *ContextFile) Exists(projectName, environment string) (bool, error) {
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

func (pr *ContextFile) getFilePath(projectName, environment string) (string, error) {
	homeDir, err := pr.getHomeDirPath()
	if err != nil {
		return "", err
	}

	var nameFile = fmt.Sprintf("%s%s", environment, CONTEXT_FILE_NAME)

	return filepath.Join(homeDir, projectName, nameFile), nil
}

func (pr *ContextFile) getHomeDirPath() (string, error) {
	if fastDeployHome := os.Getenv("FASTDEPLOY_HOME"); fastDeployHome != "" {
		return fastDeployHome, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDir), nil
}
 */