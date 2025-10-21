package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
)

type VarsRepository struct {
	pathStateEnvironment string
}

func NewVarsRepository(
	pathStateRootFastDeploy string,
	projectName string,
	repositoryName string,
	environment string) (ports.VariablesRepository, error) {

	if pathStateRootFastDeploy == "" {
		return nil, fmt.Errorf("path state root is required")
	}
	if projectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if repositoryName == "" {
		return nil, fmt.Errorf("repository name is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("environment is required")
	}

	pathStateEnvironment := filepath.Join(pathStateRootFastDeploy, projectName, repositoryName, environment)

	return &VarsRepository{
		pathStateEnvironment: pathStateEnvironment,
	}, nil
}

func (r *VarsRepository) FindByStepName(stepName string) (map[string]string, error) {
	varsFilePath := r.getPathFileVars(stepName)

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

func (r *VarsRepository) Save(stepName string, vars map[string]string) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(vars); err != nil {
		return fmt.Errorf("error al serializar las variables a formato gob: %w", err)
	}

	varsFilePath := r.getPathFileVars(stepName)

	if err := os.MkdirAll(filepath.Dir(varsFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para las variables: %w", err)
	}

	return os.WriteFile(varsFilePath, buffer.Bytes(), 0644)
}

func (r *VarsRepository) getPathFileVars(stepName string) string {
	return filepath.Join(r.pathStateEnvironment, fmt.Sprintf("vars%s.gob", stepName))
}
