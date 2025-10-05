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
	repoPath string
}

func NewStepVariableRepository(repoPath string) ports.StepVariableRepository {
	return &StepVariableRepository{
		repoPath: repoPath,
	}
}

func (r *StepVariableRepository) Load(environment string, stepName string) ([]vos.Variable, error) {
	varsPath := filepath.Join(
		r.repoPath,
		"variables",
		environment,
		fmt.Sprintf("%s.yaml", stepName),
	)

	if _, err := os.Stat(varsPath); os.IsNotExist(err) {
		return []vos.Variable{}, nil
	}

	data, err := os.ReadFile(varsPath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo de variables '%s': %w", varsPath, err)
	}

	var dtos []dto.VariableDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("error al parsear YAML en '%s': %w", varsPath, err)
	}

	return mapper.VariablesToDomain(dtos)
}
