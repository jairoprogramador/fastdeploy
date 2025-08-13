package project

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/git"
)

const (
	defaultVersion       = "1.0.0"
	defaultTechnology    = "springboot"
	defaultTeamName      = "itachi"
	defaultOrganization  = "FastDeploy"
	defaultRepositoryURL = "https://github.com/jairoprogramador/mydeploy.git"
)

type ProjectInitializeImpl struct {
	projectService ProjectService
	projectCreate  ProjectCreator
	configService  config.ConfigService
	gitService     git.GitService
}

func NewProjectInitialize(
	projectService ProjectService,
	projectCreate ProjectCreator,
	configService config.ConfigService,
	gitService git.GitService,
) ProjectInitialize {
	return &ProjectInitializeImpl{
		projectService: projectService,
		projectCreate:  projectCreate,
		configService:  configService,
		gitService:     gitService,
	}
}

func (pis *ProjectInitializeImpl) Initialize() (*ProjectEntity, error) {

	projectEntity, err := pis.projectCreate.Create()
	if err != nil {
		return nil, fmt.Errorf("error al crear el proyecto: %w", err)
	}

	configEntity, err := pis.configService.Load()
	if err != nil {
		return nil, fmt.Errorf("error al cargar la configuración: %w", err)
	}

	organization := defaultOrganization
	repository := defaultRepositoryURL
	teamName := defaultTeamName

	if configEntity.Organization != "" {
		organization = configEntity.Organization
	}
	if configEntity.Repository != "" {
		repository = configEntity.Repository
	}
	if configEntity.TeamName != "" {
		teamName = configEntity.TeamName
	}

	projectEntity.Version = defaultVersion
	projectEntity.Technology = defaultTechnology
	projectEntity.Organization = organization
	projectEntity.Repository = repository
	projectEntity.TeamName = teamName

	if err := pis.gitService.Clone(projectEntity.Repository); err != nil {
		return nil, fmt.Errorf("error al clonar el repositorio: %w", err)
	}

	if err := pis.projectService.Save(*projectEntity); err != nil {
		return nil, fmt.Errorf("error al guardar el proyecto: %w", err)
	}

	return projectEntity, nil
}

func (pis *ProjectInitializeImpl) IsInitialized() bool {
	if pis.projectService.Exists() {
		projectEntity, err := pis.projectService.Load()
		if err != nil {
			return false
		}
		return pis.gitService.IsCloned(projectEntity.Repository)
	}
	return false
}
