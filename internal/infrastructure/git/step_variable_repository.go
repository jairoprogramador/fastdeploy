package git

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/mapper"
)

type StepVariableRepository struct {
	pathResposioryRootFastDeploy string
	environment string
}

func NewStepVariableRepository(pathResposioryRootFastDeploy string, environment string) ports.StepVariableRepository {
	return &StepVariableRepository{
		pathResposioryRootFastDeploy: pathResposioryRootFastDeploy,
		environment: environment,
	}
}

func (r *StepVariableRepository) Load(stepName string) ([]vos.Variable, error) {
	pathFile := r.getPathFile(stepName)

	if _, err := os.Stat(pathFile); os.IsNotExist(err) {
		return []vos.Variable{}, nil
	}

	data, err := os.ReadFile(pathFile)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo de variables '%s': %w", pathFile, err)
	}

	var dtos []dto.VariableDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("error al parsear YAML en '%s': %w", pathFile, err)
	}

	return mapper.VariablesToDomain(dtos)
}

func (r *StepVariableRepository) getPathFile(stepName string) string {
	return filepath.Join(
		r.pathResposioryRootFastDeploy,
		"variables",
		r.environment,
		fmt.Sprintf("%s.yaml", stepName),
	)
}