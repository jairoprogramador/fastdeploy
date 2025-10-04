package executionstate

import (
	"fmt"
	"path/filepath"
	"bytes"
	"encoding/gob"
	"os"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/executionstate/mapper"
)

type VarsRepository struct {
	pathProjectFastDeploy string
}

func NewVarsRepository(pathProjectsFastDeploy string, projectName string) (ports.VarsRepository, error) {
	if projectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if pathProjectsFastDeploy == "" {
		return nil, fmt.Errorf("base path is required")
	}

	basePathProject := filepath.Join(pathProjectsFastDeploy, projectName)

	return &VarsRepository{pathProjectFastDeploy: basePathProject}, nil
}

func (r *VarsRepository) Save(variables []vos.Variable, environment string) error {
	varsMap := mapper.VarsToDTO(variables)

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(varsMap); err != nil {
		return fmt.Errorf("error al serializar las variables a formato gob: %w", err)
	}

	varsFilePath := r.getPathFileEnvironment(environment)

	if err := os.MkdirAll(filepath.Dir(varsFilePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para las variables: %w", err)
	}

	return os.WriteFile(varsFilePath, buffer.Bytes(), 0644)
}

func (r *VarsRepository) GetStore(environment string) ([]vos.Variable, error) {
	filePath := r.getPathFileEnvironment(environment)

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

func (r *VarsRepository) getPathFileEnvironment(environmentName string) string {
	return filepath.Join(r.pathProjectFastDeploy, "environment", fmt.Sprintf("%s.vars", environmentName))
}