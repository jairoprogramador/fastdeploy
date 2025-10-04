package git

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git/mapper"
)

// StepVariableRepository implementa la interfaz ports.StepVariableRepository.
// Es un adaptador que carga las variables desde archivos YAML dentro de la estructura
// del repositorio de plantillas Git.
type StepVariableRepository struct {
	// Este repositorio necesita conocer la ruta local del repositorio ya clonado y verificado.
	// El servicio de aplicación se la proporcionará.
	repoPath string
}

// NewStepVariableRepository crea una nueva instancia del repositorio de variables.
func NewStepVariableRepository(repoPath string) ports.StepVariableRepository {
	return &StepVariableRepository{
		repoPath: repoPath,
	}
}

// Load busca y parsea el archivo de variables para una combinación específica de ambiente y paso.
func (r *StepVariableRepository) Load(environment string, stepName string) ([]vos.Variable, error) {
	// La convención de ruta es: <repo>/variables/<valor_ambiente>/<nombre_paso>.yaml
	varsPath := filepath.Join(
		r.repoPath,
		"variables",
		environment,
		fmt.Sprintf("%s.yaml", stepName),
	)

	// Es válido que un paso no tenga un archivo de variables.
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
