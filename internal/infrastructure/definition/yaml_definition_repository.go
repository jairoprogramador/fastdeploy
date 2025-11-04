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

	defAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/aggregates"
	defEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/entities"
	defPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/ports"
	defVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"

	iDefDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition/dto"
	iDefMap "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/definition/mapper"
)

type YamlTemplateRepository struct {
	pathRepositoriesRoot string
	executor             appPor.CommandService
}

func NewYamlTemplateRepository(
	rootRepositoriesPath string,
	executor appPor.CommandService) defPor.DefinitionRepository {

	return &YamlTemplateRepository{
		pathRepositoriesRoot: rootRepositoriesPath,
		executor:             executor,
	}
}

func (r *YamlTemplateRepository) PathLocal(source proVos.Template) string {
	return filepath.Join(r.pathRepositoriesRoot, source.NameTemplate())
}

func (r *YamlTemplateRepository) LoadEnvironments(ctx context.Context, source proVos.Template) ([]defVos.EnvironmentDefinition, error) {
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

func (r *YamlTemplateRepository) LoadDeployment(ctx context.Context, source proVos.Template, environment string) (*defAgg.Deployment, error) {
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

func (r *YamlTemplateRepository) cloneRepository(ctx context.Context, source proVos.Template, pathRepository string) error {

	if err := os.MkdirAll(pathRepository, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(pathRepository, ".git")); os.IsNotExist(err) {
		cloneCmd := fmt.Sprintf("git clone %s %s", source.URL(), pathRepository)
		_, exitCode, err := r.executor.Run(ctx, r.pathRepositoriesRoot, cloneCmd)
		if err != nil {
			return fmt.Errorf("error al clonar el repositorio '%s': %w", source.URL(), err)
		}
		if exitCode != 0 {
			return fmt.Errorf("error al clonar el repositorio '%s': exit code %d", source.URL(), exitCode)
		}
	} else {
		fetchCmd := "git fetch --all"
		_, exitCode, err := r.executor.Run(ctx, pathRepository, fetchCmd)
		if err != nil {
			return fmt.Errorf("error al actualizar el repositorio '%s': %w", source.URL(), err)
		}
		if exitCode != 0 {
			return fmt.Errorf("error al actualizar el repositorio '%s': exit code %d", source.URL(), exitCode)
		}
	}

	checkoutCmd := fmt.Sprintf("git checkout %s", source.Ref())
	_, exitCode, err := r.executor.Run(ctx, pathRepository, checkoutCmd)
	if  err != nil {
		return fmt.Errorf("error al hacer checkout a la referencia '%s' en '%s': %w", source.Ref(), pathRepository, err)
	}
	if exitCode != 0 {
		return fmt.Errorf("error al hacer checkout a la referencia '%s' en '%s': exit code %d", source.Ref(), pathRepository, exitCode)
	}

	return nil
}

func (r *YamlTemplateRepository) loadFromSource(environment, pathRepository string) (*defAgg.Deployment, error) {
	environments, err := r.loadEnvironments(pathRepository)
	if err != nil {
		return nil, err
	}

	steps, err := r.loadSteps(pathRepository, environment)
	if err != nil {
		return nil, err
	}

	return defAgg.NewDeployment(environments, steps)
}

func (r *YamlTemplateRepository) loadEnvironments(repositoryPath string) ([]defVos.EnvironmentDefinition, error) {
	environmentsPath := filepath.Join(repositoryPath, "environments.yaml")

	data, err := os.ReadFile(environmentsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []defVos.EnvironmentDefinition{}, nil
		}
		return []defVos.EnvironmentDefinition{}, fmt.Errorf("no se pudo leer el archivo de ambientes: %w", err)
	}

	var dtos []iDefDto.EnvironmentDefinitionDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return []defVos.EnvironmentDefinition{}, fmt.Errorf("error al parsear el archivo YAML de ambientes: %w", err)
	}

	return iDefMap.EnvironmentsToDomain(dtos)
}

func (r *YamlTemplateRepository) loadSteps(repositoryPath string, environment string) ([]defEnt.StepDefinition, error) {
	stepsPath := filepath.Join(repositoryPath, "steps")

	directoriesSteps, err := os.ReadDir(stepsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []defEnt.StepDefinition{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el directorio de pasos '%s': %w", stepsPath, err)
	}

	var directoriesStepsNames []string
	for _, entry := range directoriesSteps {
		if entry.IsDir() {
			directoriesStepsNames = append(directoriesStepsNames, entry.Name())
		}
	}

	var stepsDefinitions []defEnt.StepDefinition
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

		stepDefinition, err := defEnt.NewStepDefinition(stepName, triggers, commands, variables)
		if err != nil {
			return nil, fmt.Errorf("error al crear la definición del paso '%s': %w", stepName, err)
		}
		stepsDefinitions = append(stepsDefinitions, stepDefinition)
	}

	return stepsDefinitions, nil
}

func (r *YamlTemplateRepository) loadTriggers(directoryStepPath string) ([]defVos.TriggerDefinition, error) {
	triggersPath := filepath.Join(directoryStepPath, "triggers.yaml")

	data, err := os.ReadFile(triggersPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []defVos.TriggerDefinition{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el archivo de triggers '%s': %w", triggersPath, err)
	}

	var scopes []string
	if err := yaml.Unmarshal(data, &scopes); err != nil {
		return nil, fmt.Errorf("error al parsear archivo YAML en '%s': %w", triggersPath, err)
	}

	return iDefMap.TriggersToDomain(scopes), nil
}

func (r *YamlTemplateRepository) loadCommands(directoryStepPath string) ([]defVos.CommandDefinition, error) {
	commandsPath := filepath.Join(directoryStepPath, "commands.yaml")

	data, err := os.ReadFile(commandsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []defVos.CommandDefinition{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el archivo de comandos '%s': %w", commandsPath, err)
	}

	var dtos []iDefDto.CommandDefinitionDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("error al parsear archivo YAML en '%s': %w", commandsPath, err)
	}

	commands, err := iDefMap.CommandsToDomain(dtos)
	if err != nil {
		return nil, err
	}

	return commands, nil
}

func (r *YamlTemplateRepository) loadVariables(repositoryPath string, stepName, environment string) ([]defVos.VariableDefinition, error) {
	variablesPath := filepath.Join(repositoryPath, "variables",
		environment, fmt.Sprintf("%s.yaml", stepName))

	data, err := os.ReadFile(variablesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []defVos.VariableDefinition{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el archivo de variables '%s': %w", variablesPath, err)
	}

	var dtos []iDefDto.VariableDefinitionDTO
	if err := yaml.Unmarshal(data, &dtos); err != nil {
		return nil, fmt.Errorf("error al parsear archivo YAML en '%s': %w", variablesPath, err)
	}

	return iDefMap.VariablesToDomain(dtos)
}

var stepNameRegex = regexp.MustCompile(`^\d+-(.*)$`)

func (r *YamlTemplateRepository) extractStepName(dirName string) (string, error) {
	matches := stepNameRegex.FindStringSubmatch(dirName)
	if len(matches) < 2 {
		return "", fmt.Errorf("el nombre del directorio '%s' no sigue la convención 'NN-nombre'", dirName)
	}
	return matches[1], nil
}
