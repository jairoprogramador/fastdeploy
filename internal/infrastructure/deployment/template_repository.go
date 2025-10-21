package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"

	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"

	sharedvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"

	depAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/aggregates"
	depEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/entities"
	depPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/ports"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"

	iDepDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/deployment/dto"
	iDepMap "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/deployment/mapper"
)

type TemplateRepository struct {
	rootRepositoriesPath string
	environment string
	executor              appPor.CommandExecutor
}

func NewTemplateRepository(
	rootRepositoriesPath string,
	environment string,
	executor appPor.CommandExecutor) depPor.TemplateRepository {

	return &TemplateRepository{
		rootRepositoriesPath: rootRepositoriesPath,
		environment: environment,
		executor: executor,
	}
}

func (r *TemplateRepository) Load(ctx context.Context, source sharedvos.TemplateSource) (*depAgg.DeploymentTemplate, string, error) {
	repositoryPath, err := r.gitCloneRepository(ctx, source)
	if err != nil {
		return nil, "", err
	}

	template, err := r.loadFromSource(repositoryPath, source)
	if err != nil {
		return nil, "", err
	}

	return template, repositoryPath, nil
}

func (r *TemplateRepository) gitCloneRepository(ctx context.Context, source sharedvos.TemplateSource) (string, error) {
	repositoryPath := filepath.Join(r.rootRepositoriesPath, source.NameTemplate())

	if err := os.MkdirAll(repositoryPath, 0755); err != nil {
		return "", err
	}

	if _, err := os.Stat(filepath.Join(repositoryPath, ".git")); os.IsNotExist(err) {
		cloneCmd := fmt.Sprintf("git clone %s %s", source.Url(), repositoryPath)
		_, _, err := r.executor.Execute(ctx, r.rootRepositoriesPath, cloneCmd)
		if err != nil {
			return "", fmt.Errorf("error al clonar el repositorio '%s': %w", source.Url(), err)
		}
	} else {
		fetchCmd := "git fetch --all"
		_, _, err := r.executor.Execute(ctx, repositoryPath, fetchCmd)
		if err != nil {
			return "", fmt.Errorf("error al actualizar el repositorio '%s': %w", source.Url(), err)
		}
	}

	checkoutCmd := fmt.Sprintf("git checkout %s", source.Ref())
	if _, _, err := r.executor.Execute(ctx, repositoryPath, checkoutCmd); err != nil {
		return repositoryPath, fmt.Errorf("error al hacer checkout a la referencia '%s' en '%s': %w", source.Ref(), repositoryPath, err)
	}

	return repositoryPath, nil
}

func (r *TemplateRepository) loadFromSource(repositoryPath string, source sharedvos.TemplateSource) (*depAgg.DeploymentTemplate, error) {
	environments, err := r.loadEnvironmentsFromSource(repositoryPath)
	if err != nil {
		return nil, err
	}

	steps, err := r.loadStepsFromSource(repositoryPath)
	if err != nil {
		return nil, err
	}

	return depAgg.NewDeploymentTemplate(source, environments, steps)
}

func (r *TemplateRepository) loadEnvironmentsFromSource(repositoryPath string) ([]depVos.Environment, error) {
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

func (r *TemplateRepository) loadStepsFromSource(repositoryPath string) ([]depEnt.StepDefinition, error) {
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

		triggers, err := r.loadTriggersFromSource(directoryStepPath)
		if err != nil {
			return nil, err
		}

		commands, err := r.loadCommandsFromSource(directoryStepPath)
		if err != nil {
			return nil, err
		}

		variables, err := r.loadVariablesFromSource(repositoryPath, stepName)
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

func (r *TemplateRepository) loadTriggersFromSource(directoryStepPath string) ([]depVos.Trigger, error) {
	triggersPath := filepath.Join(directoryStepPath, "triggers.yaml")

	data, err := os.ReadFile(triggersPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []depVos.Trigger{}, nil
		}
		return nil, fmt.Errorf("no se pudo leer el archivo de triggers '%s': %w", triggersPath, err)
	}

	var triggerDTO iDepDto.TriggerDTO
	if err := yaml.Unmarshal(data, &triggerDTO); err != nil {
		return nil, fmt.Errorf("error al parsear archivo YAML en '%s': %w", triggersPath, err)
	}

	return iDepMap.TriggersToDomain(triggerDTO.Scopes), nil
}

func (r *TemplateRepository) loadCommandsFromSource(directoryStepPath string) ([]depVos.CommandDefinition, error) {
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

func (r *TemplateRepository) loadVariablesFromSource(repositoryPath string, stepName string) ([]depVos.Variable, error) {
	variablesPath := filepath.Join(repositoryPath,"variables",
		r.environment, fmt.Sprintf("%s.yaml", stepName))

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
