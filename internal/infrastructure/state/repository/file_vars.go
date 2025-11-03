package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
)

type FileVarsRepository struct {
	pathStateRoot string
}

func NewFileVarsRepository(pathStateRoot string) ports.VariablesRepository {
	return &FileVarsRepository{
		pathStateRoot: pathStateRoot,
	}
}

func (r *FileVarsRepository) FindByStep(namesRequest dto.NamesParams, runParams dto.RunParams) (map[string]string, error) {
	varsFilePath := r.getPathFileVars(namesRequest, runParams)

	data, err := os.ReadFile(varsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, fmt.Errorf("error al leer el archivo de variables: %w", err)
	}

	var varsMap map[string]string
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&varsMap); err != nil {
		return nil, fmt.Errorf("error al deserializar las variables: %w", err)
	}

	return varsMap, nil
}

func (r *FileVarsRepository) SaveByStep(namesRequest dto.NamesParams, runParams dto.RunParams, vars map[string]string) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(vars); err != nil {
		return fmt.Errorf("error al serializar las variables a formato gob: %w", err)
	}

	varsFilePath := r.getPathFileVars(namesRequest, runParams)

	if err := os.MkdirAll(filepath.Dir(varsFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para las variables: %w", err)
	}

	return os.WriteFile(varsFilePath, buffer.Bytes(), 0644)
}

func (r *FileVarsRepository) getPathFileVars(namesRequest dto.NamesParams, runParams dto.RunParams) string {
	pathStateEnvironment := filepath.Join(r.pathStateRoot, namesRequest.ProjectName(), namesRequest.RepositoryName(), runParams.Environment())
	return filepath.Join(pathStateEnvironment, fmt.Sprintf("vars%s.gob", runParams.StepName()))
}
