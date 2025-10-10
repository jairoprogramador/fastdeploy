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

	applicationports "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/dto"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/mapper"
	deploymentports "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/ports"
)

type TemplateRepository struct {
	reposBasePath string
	executor      applicationports.CommandExecutor
}

func NewTemplateRepository(
	reposBasePath string, executor applicationports.CommandExecutor) deploymentports.TemplateRepository {
	return &TemplateRepository{
		reposBasePath: reposBasePath,
		executor:      executor,
	}
}

func (r *TemplateRepository) GetRepositoryName(repoURL string) (string, error) {
	parsed, err := url.Parse(repoURL)
	if err != nil {
		return "", fmt.Errorf("URL de repositorio no válida: %w", err)
	}

	safePath := strings.Split(parsed.Path, "/")
	lastPart := safePath[len(safePath)-1]
	repositoryName := strings.TrimSuffix(lastPart, ".git")

	return repositoryName, nil
}

func (r *TemplateRepository) GetTemplate(ctx context.Context, source vos.TemplateSource) (*aggregates.DeploymentTemplate, string, error) {
	repoPath, err := r.ensureRepo(ctx, source.RepoURL())
	if err != nil {
		return nil, "", err
	}

	checkoutCmd := fmt.Sprintf("git checkout %s", source.Ref())
	if _, _, err := r.executor.Execute(ctx, repoPath, checkoutCmd); err != nil {
		return nil, "", fmt.Errorf("error al hacer checkout a la referencia '%s' en '%s': %w", source.Ref(), repoPath, err)
	}

	template, err := r.buildTemplateFromFile(source, repoPath)
	if err != nil {
		return nil, "", err
	}

	return template, repoPath, nil
}

func (r *TemplateRepository) ensureRepo(ctx context.Context, repoURL string) (string, error) {
	repoPath, err := r.repoPathFromURL(repoURL)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(repoPath, 0755); err != nil {
		return "", err
	}

	if _, err := os.Stat(filepath.Join(repoPath, ".git")); os.IsNotExist(err) {
		cloneCmd := fmt.Sprintf("git clone %s %s", repoURL, repoPath)
		_, _, err := r.executor.Execute(ctx, r.reposBasePath, cloneCmd)
		if err != nil {
			return "", fmt.Errorf("error al clonar el repositorio '%s': %w", repoURL, err)
		}
	} else {
		fetchCmd := "git fetch --all"
		_, _, err := r.executor.Execute(ctx, repoPath, fetchCmd)
		if err != nil {
			return "", fmt.Errorf("error al actualizar el repositorio '%s': %w", repoURL, err)
		}
	}
	return repoPath, nil
}

func (r *TemplateRepository) buildTemplateFromFile(source vos.TemplateSource, repoPath string) (*aggregates.DeploymentTemplate, error) {
	environments, err := r.parseEnvironments(repoPath)
	if err != nil {
		return nil, err
	}

	steps, err := r.parseSteps(repoPath)
	if err != nil {
		return nil, err
	}

	return aggregates.NewDeploymentTemplate(source, environments, steps)
}

func (r *TemplateRepository) parseEnvironments(repoPath string) ([]vos.Environment, error) {
	filePath := filepath.Join(repoPath, "environments.yaml")
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

func (r *TemplateRepository) parseSteps(repoPath string) ([]entities.StepDefinition, error) {
	stepsRootPath := filepath.Join(repoPath, "steps")

	entries, err := os.ReadDir(stepsRootPath)
	if err != nil {
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

	var stepsDefinitions []entities.StepDefinition
	for _, dirName := range dirNames {
		stepName, err := extractStepName(dirName)
		if err != nil {
			continue // Ignorar directorios que no siguen la convención
		}

		stepDirPath := filepath.Join(stepsRootPath, dirName)

		var verifications []vos.VerificationType
		metaPath := filepath.Join(stepDirPath, "verifications.yaml")
		if _, err := os.Stat(metaPath); !os.IsNotExist(err) {
			data, err := os.ReadFile(metaPath)
			if err != nil {
				return nil, fmt.Errorf("no se pudo leer el archivo de metadatos '%s': %w", metaPath, err)
			}
			var stepDTO dto.StepDTO
			if err := yaml.Unmarshal(data, &stepDTO); err != nil {
				return nil, fmt.Errorf("error al parsear YAML en '%s': %w", metaPath, err)
			}
			verifications, err = mapper.VerificationsToDomain(stepDTO.Verifications)
			if err != nil {
				return nil, fmt.Errorf("error en los metadatos de '%s': %w", metaPath, err)
			}
		}

		var commands []vos.CommandDefinition
		commandsPath := filepath.Join(stepDirPath, "commands.yaml")
		if _, err := os.Stat(commandsPath); !os.IsNotExist(err) {
			data, err := os.ReadFile(commandsPath)
			if err != nil {
				return nil, fmt.Errorf("no se pudo leer el archivo de comandos '%s': %w", commandsPath, err)
			}
			var dtos []dto.CommandDefinitionDTO
			if err := yaml.Unmarshal(data, &dtos); err != nil {
				return nil, fmt.Errorf("error al parsear YAML en '%s': %w", commandsPath, err)
			}
			commands, err = mapper.CommandsToDomain(dtos)
			if err != nil {
				return nil, fmt.Errorf("error al mapear comandos desde '%s': %w", commandsPath, err)
			}
		}

		stepDefinition, err := entities.NewStepDefinition(stepName, verifications, commands)
		if err != nil {
			return nil, fmt.Errorf("error al crear la definición del paso '%s': %w", stepName, err)
		}
		stepsDefinitions = append(stepsDefinitions, stepDefinition)
	}

	return stepsDefinitions, nil
}

var stepNameRegex = regexp.MustCompile(`^\d+-(.*)$`)

func extractStepName(dirName string) (string, error) {
	matches := stepNameRegex.FindStringSubmatch(dirName)
	if len(matches) < 2 {
		return "", fmt.Errorf("el nombre del directorio '%s' no sigue la convención 'NN-nombre'", dirName)
	}
	return matches[1], nil
}

func (r *TemplateRepository) repoPathFromURL(repoURL string) (string, error) {
	repositoryName, err := r.GetRepositoryName(repoURL)
	if err != nil {
		return "", err
	}

	return filepath.Join(r.reposBasePath, repositoryName), nil
}
