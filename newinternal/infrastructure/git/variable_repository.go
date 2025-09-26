package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"

	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git/mapper"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
)

// VariableRepository implementa la interfaz ports.StepVariableRepository.
// Es un adaptador que carga las variables desde archivos YAML dentro de la estructura
// del repositorio de plantillas Git.
type VariableRepository struct {
	// Este repositorio necesita conocer la ruta local del repositorio ya clonado y verificado.
	// El servicio de aplicación se la proporcionará.
	repoPath string
}

// NewVariableRepository crea una nueva instancia del repositorio de variables.
func NewVariableRepository(repoPath string) *VariableRepository {
	return &VariableRepository{
		repoPath: repoPath,
	}
}

// Load busca y parsea el archivo de variables para una combinación específica de ambiente y paso.
func (r *VariableRepository) Load(
	_ context.Context,
	environment deploymentvos.Environment,
	stepDefinition deploymententities.StepDefinition,
) ([]vos.Variable, error) {
	// La convención de ruta es: <repo>/variables/<valor_ambiente>/<nombre_paso>.yaml
	varsPath := filepath.Join(
		r.repoPath,
		"variables",
		environment.Value(),
		fmt.Sprintf("%s.yaml", stepDefinition.Name()),
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
