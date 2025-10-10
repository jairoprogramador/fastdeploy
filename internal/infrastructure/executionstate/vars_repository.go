package executionstate

import (
	"fmt"
	"path/filepath"
	"bytes"
	"encoding/gob"
	"os"
	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate/mapper"
)

type VarsRepository struct {
	pathStateProjectFastDeploy string
}

func NewVarsRepository(
	pathStateRootFastDeploy string,
	projectName string,
	repositoryName string,
	environment string) (ports.VarsRepository, error) {

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

	pathStateProject := filepath.Join(pathStateRootFastDeploy, projectName, repositoryName, environment)

	return &VarsRepository{pathStateProjectFastDeploy: pathStateProject}, nil
}

func (r *VarsRepository) Save(variables []vos.Variable) error {
	varsMapDto := mapper.VarsToDTO(variables)

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(varsMapDto); err != nil {
		return fmt.Errorf("error al serializar las variables a formato gob: %w", err)
	}

	varsFilePath := r.getPathFileVars()

	if err := os.MkdirAll(filepath.Dir(varsFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para las variables: %w", err)
	}

	return os.WriteFile(varsFilePath, buffer.Bytes(), 0644)
}

func (r *VarsRepository) FindAll() ([]vos.Variable, error) {
	filePath := r.getPathFileVars()

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []vos.Variable{}, nil
		}
		return nil, fmt.Errorf("error al leer el archivo de variables: %w", err)
	}

	var varsMap map[string]string
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&varsMap); err != nil {
		return nil, fmt.Errorf("error al deserializar las variables: %w", err)
	}

	return mapper.VarsToDomain(varsMap), nil
}

func (r *VarsRepository) getPathFileVars() string {
	return filepath.Join(r.pathStateProjectFastDeploy, "vars.state")
}