package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"

	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	proVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/vos"

	depAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/aggregates"
	depEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/entities"
	depPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/ports"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"

	iDepDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/template/dto"
	iDepMap "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/template/mapper"
)

type TemplateRepository struct {
	pathRepositoriesRoot string
	executor appPor.CommandExecutor
}

func NewTemplateRepository(
	rootRepositoriesPath string,
	executor appPor.CommandExecutor) depPor.TemplateRepository {

	return &TemplateRepository{
		pathRepositoriesRoot: rootRepositoriesPath,
		executor: executor,
	}
}

func (r *TemplateRepository) PathLocal(source proVos.Template) string {
	return filepath.Join(r.pathRepositoriesRoot, source.NameTemplate())
}


func (r *TemplateRepository) LoadEnvironments(ctx context.Context, source proVos.Template) ([]depVos.Environment, error) {
	pathRepository := r.PathLocal(source)

	err := r.cloneRepository(ctx, source, pathRepository)
	if err != nil {
		return nil, err
	}

	environments, err := r.loadEnvironments(pathRepository)
	if err != nil {
		return nil, err
	}

	return environments, nil
}

func (r *TemplateRepository) LoadDeployment(ctx context.Context, source proVos.Template, environment string) (*depAgg.Deployment, error) {
	pathRepository := r.PathLocal(source)

	err := r.cloneRepository(ctx, source, pathRepository)
	if err != nil {
		return nil, err
	}

	template, err := r.loadFromSource(environment, pathRepository)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func (r *TemplateRepository) cloneRepository(ctx context.Context, source proVos.Template, pathRepository string) error {

	if err := os.MkdirAll(pathRepository, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(pathRepository, ".git")); os.IsNotExist(err) {
		cloneCmd := fmt.Sprintf("git clone %s %s", source.URL(), pathRepository)
		_, _, err := r.executor.Execute(ctx, r.pathRepositoriesRoot, cloneCmd)
		if err != nil {
			return fmt.Errorf("error al clonar el repositorio '%s': %w", source.URL(), err)
		}
	} else {
		fetchCmd := "git fetch --all"
		_, _, err := r.executor.Execute(ctx, pathRepository, fetchCmd)
		if err != nil {
			return fmt.Errorf("error al actualizar el repositorio '%s': %w", source.URL(), err)
		}
	}

	checkoutCmd := fmt.Sprintf("git checkout %s", source.Ref())
	if _, _, err := r.executor.Execute(ctx, pathRepository, checkoutCmd); err != nil {
		return fmt.Errorf("error al hacer checkout a la referencia '%s' en '%s': %w", source.Ref(), pathRepository, err)
	}

	return nil
}

func (r *TemplateRepository) loadFromSource(environment, pathRepository string) (*depAgg.Deployment, error) {
	environments, err := r.loadEnvironments(pathRepository)
	if err != nil {
		return nil, err
	}

	steps, err := r.loadSteps(pathRepository, environment)
	if err != nil {
		return nil, err
	}

	return depAgg.NewDeployment(environments, steps)
}

func (r *TemplateRepository) loadEnvironments(repositoryPath string) ([]depVos.Environment, error) {
	environmentsPath := filepath.Join(repositoryPath, "environments.yaml")

	data, err := os.ReadFile(environmentsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []depVos.Environment{}, nil
		}
		return []depVos.Environment{}, fmt.Errorf("no se pudo leer el archivo de ambientes: %w", err)
	}

	var dtos []iDepDto.EnvironmentDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return []depVos.Environment{}, fmt.Errorf("error al parsear el archivo YAML de ambientes: %w", err)
	}

	return iDepMap.EnvironmentsToDomain(dtos)
}

func (r *TemplateRepository) loadSteps(repositoryPath string, environment string) ([]depEnt.StepDefinition, error) {
	stepsPath := filepath.Join(repositoryPath, "steps")

	directoriesSteps, err := os.ReadDir(stepsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []depEnt.StepDefinition{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el directorio de pasos '%s': %w", stepsPath, err)
	}

	var directoriesStepsNames []string
	for _, entry := range directoriesSteps {
		if entry.IsDir() {
			directoriesStepsNames = append(directoriesStepsNames, entry.Name())
		}
	}

	var stepsDefinitions []depEnt.StepDefinition
	for _, directoryStepName := range directoriesStepsNames {
		stepName, err := r.extractStepName(directoryStepName)
		if err != nil {
			continue
		}

		directoryStepPath := filepath.Join(stepsPath, directoryStepName)

		triggers, err := r.loadTriggers(directoryStepPath)
		if err != nil {
			return nil, err
		}

		commands, err := r.loadCommands(directoryStepPath)
		if err != nil {
			return nil, err
		}

		variables, err := r.loadVariables(repositoryPath, stepName, environment)
		if err != nil {
			return nil, err
		}

		stepDefinition, err := depEnt.NewStepDefinition(stepName, triggers, commands, variables)
		if err != nil {
			return nil, fmt.Errorf("error al crear la definición del paso '%s': %w", stepName, err)
		}
		stepsDefinitions = append(stepsDefinitions, stepDefinition)
	}

	return stepsDefinitions, nil
}

func (r *TemplateRepository) loadTriggers(directoryStepPath string) ([]depVos.Trigger, error) {
	triggersPath := filepath.Join(directoryStepPath, "triggers.yaml")

	data, err := os.ReadFile(triggersPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []depVos.Trigger{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el archivo de triggers '%s': %w", triggersPath, err)
	}

	var scopes []string
	if err := yaml.Unmarshal(data, &scopes); err != nil {
		return nil, fmt.Errorf("error al parsear archivo YAML en '%s': %w", triggersPath, err)
	}

	return iDepMap.TriggersToDomain(scopes), nil
}

func (r *TemplateRepository) loadCommands(directoryStepPath string) ([]depVos.CommandDefinition, error) {
	commandsPath := filepath.Join(directoryStepPath, "commands.yaml")

	data, err := os.ReadFile(commandsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []depVos.CommandDefinition{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el archivo de comandos '%s': %w", commandsPath, err)
	}

	var dtos []iDepDto.CommandDefinitionDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("error al parsear archivo YAML en '%s': %w", commandsPath, err)
	}

	commands, err := iDepMap.CommandsToDomain(dtos)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (r *TemplateRepository) loadVariables(repositoryPath string, stepName, environment string) ([]depVos.Variable, error) {
	variablesPath := filepath.Join(repositoryPath, "variables",
		environment, fmt.Sprintf("%s.yaml", stepName))

	data, err := os.ReadFile(variablesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []depVos.Variable{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el archivo de variables '%s': %w", variablesPath, err)
	}

	var dtos []iDepDto.VariableDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("error al parsear archivo YAML en '%s': %w", variablesPath, err)
	}

	return iDepMap.VariablesToDomain(dtos)
}

var stepNameRegex = regexp.MustCompile(`^\d+-(.*)$`)

func (r *TemplateRepository) extractStepName(dirName string) (string, error) {
	matches := stepNameRegex.FindStringSubmatch(dirName)
	if len(matches) < 2 {
		return "", fmt.Errorf("el nombre del directorio '%s' no sigue la convención 'NN-nombre'", dirName)
	}
	return matches[1], nil
}
