package git

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/newinternal/application/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git/mapper"
)

/* type stepDefinitionDTO struct {
	Name     string                 `yaml:"name"`
	Commands []dto.CommandDefinitionDTO `yaml:"commands"`
} */

// TemplateRepository implementa la interfaz ports.TemplateRepository.
// Es un adaptador que obtiene la definición de un despliegue desde un repositorio Git.
type TemplateRepository struct {
	reposBasePath     string
	//relativeStepsPath string
	executor          ports.CommandExecutor
}

// NewTemplateRepository crea una nueva instancia del repositorio de plantillas Git.
func NewTemplateRepository(reposBasePath string, executor ports.CommandExecutor) *TemplateRepository {
	return &TemplateRepository{
		reposBasePath:     reposBasePath,
	//	relativeStepsPath: relativeStepsPath,
		executor:          executor,
	}
}

// GetTemplate orquesta la clonación/actualización de un repo Git, el checkout a una
// versión específica, y el parsing de los archivos para construir el agregado DeploymentTemplate.
func (r *TemplateRepository) GetTemplate(ctx context.Context, source vos.TemplateSource) (*aggregates.DeploymentTemplate, error) {
	repoPath, err := r.ensureRepo(ctx, source.RepoURL())
	if err != nil {
		return nil, err
	}

	// Checkout a la referencia específica para asegurar una ejecución reproducible.
	checkoutCmd := fmt.Sprintf("git checkout %s", source.Ref())
	if _, _, err := r.executor.Execute(ctx, repoPath, checkoutCmd); err != nil {
		return nil, fmt.Errorf("error al hacer checkout a la referencia '%s' en '%s': %w", source.Ref(), repoPath, err)
	}

	// Leer y construir el agregado desde los archivos.
	return r.buildTemplateFromFile(source, repoPath)
}

func (r *TemplateRepository) ensureRepo(ctx context.Context, repoURL string) (string, error) {
	repoPath, err := repoPathFromURL(r.reposBasePath, repoURL)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// Clonar si el repo no existe localmente.
		cloneCmd := fmt.Sprintf("git clone %s %s", repoURL, repoPath)
		_, _, err := r.executor.Execute(ctx, r.reposBasePath, cloneCmd)
		if err != nil {
			return "", fmt.Errorf("error al clonar el repositorio '%s': %w", repoURL, err)
		}
	} else {
		// Actualizar si ya existe.
		fetchCmd := "git fetch --all"
		_, _, err := r.executor.Execute(ctx, repoPath, fetchCmd)
		if err != nil {
			return "", fmt.Errorf("error al actualizar el repositorio '%s': %w", repoURL, err)
		}
	}
	return repoPath, nil
}

func (r *TemplateRepository) buildTemplateFromFile(source vos.TemplateSource, repoPath string) (*aggregates.DeploymentTemplate, error) {
	// Leer environments.yaml
	environments, err := r.parseEnvironments(filepath.Join(repoPath, "environments.yaml"))
	if err != nil {
		return nil, err
	}

	// Leer steps.yaml
	steps, err := r.parseSteps(filepath.Join(repoPath, "steps"))
	if err != nil {
		return nil, err
	}

	// Usar el constructor del agregado para crear una instancia válida.
	return aggregates.NewDeploymentTemplate(source, environments, steps)
}

// parseEnvironments lee y convierte el DTO de environments a objetos de valor del dominio.
func (r *TemplateRepository) parseEnvironments(filePath string) ([]vos.Environment, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo de ambientes: %w", err)
	}

	var dtos []dto.EnvironmentDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("error al parsear el YAML de ambientes: %w", err)
	}

	return mapper.EnvironmentsToDomain(dtos)
}

// parseSteps implementa la lógica de descubrimiento de pasos basada en la convención
// de nomenclatura de directorios con prefijo numérico (e.g., "01-test", "02-supply").
func (r *TemplateRepository) parseSteps(repoPath string) ([]entities.StepDefinition, error) {
	stepsRootPath := filepath.Join(repoPath, "steps")

	entries, err := os.ReadDir(stepsRootPath)
	if err != nil {
		// Si el directorio de pasos no existe, lo tratamos como si no hubiera pasos definidos.
		if os.IsNotExist(err) {
			return []entities.StepDefinition{}, fmt.Errorf("no se pudo leer el directorio de pasos '%s': %w", stepsRootPath, err)
		}
		return nil, fmt.Errorf("no se pudo leer el directorio de pasos '%s': %w", stepsRootPath, err)
	}

	var dirNames []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirNames = append(dirNames, entry.Name())
		}
	}

	// Ordenar alfabéticamente asegura el orden de ejecución correcto (01-test, 02-supply, etc.).
	// sort.Strings(dirNames) // os.ReadDir ya devuelve los resultados ordenados por nombre.

	var stepsDefinitions []entities.StepDefinition
	for _, dirName := range dirNames {
		stepName, err := extractStepName(dirName)
		if err != nil {
			// Ignorar directorios que no siguen la convención de nomenclatura.
			continue
		}

		commandsPath := filepath.Join(stepsRootPath, dirName, "commands.yaml")
		data, err := os.ReadFile(commandsPath)
		if err != nil {
			// Permitir que un paso no tenga un commands.yaml si no tiene comandos.
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("no se pudo leer el archivo de comandos '%s': %w", commandsPath, err)
		}

		var dtos []dto.CommandDefinitionDTO
		if err := yaml.Unmarshal(data, &dtos); err != nil {
			return nil, fmt.Errorf("error al parsear YAML en '%s': %w", commandsPath, err)
		}

		commands, err := mapper.CommandsToDomain(dtos)
		if err != nil {
			return nil, fmt.Errorf("error al mapear comandos desde '%s': %w", commandsPath, err)
		}

		// Si un paso no tiene comandos, no lo añadimos.
		if len(commands) == 0 {
			continue
		}

		stepDefinition, err := entities.NewStepDefinition(stepName, commands)
		if err != nil {
			return nil, fmt.Errorf("error al crear la definición del paso '%s': %w", stepName, err)
		}
		stepsDefinitions = append(stepsDefinitions, stepDefinition)
	}

	return stepsDefinitions, nil
}

var stepNameRegex = regexp.MustCompile(`^\d+-(.*)$`)

// extractStepName extrae el nombre limpio de un paso desde el nombre del directorio,
// validando que siga la convención "NN-nombre".
func extractStepName(dirName string) (string, error) {
	matches := stepNameRegex.FindStringSubmatch(dirName)
	if len(matches) < 2 {
		return "", fmt.Errorf("el nombre del directorio '%s' no sigue la convención 'NN-nombre'", dirName)
	}
	return matches[1], nil
}

// repoPathFromURL genera un nombre de directorio local seguro a partir de una URL de repo.
func repoPathFromURL(basePath, repoURL string) (string, error) {
	parsed, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("URL de repositorio no válida: %w", err)
	}
	// Limpiar el path para usarlo como nombre de directorio. e.g., /user/repo.git -> user_repo
	safePath := strings.Trim(parsed.Path, "/")
	safePath = strings.TrimSuffix(safePath, ".git")
	safePath = strings.ReplaceAll(safePath, "/", "_")
	return filepath.Join(basePath, safePath), nil
}
