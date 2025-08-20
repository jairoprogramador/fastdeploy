package project

import (
	"fmt"

	appConfig "github.com/jairoprogramador/fastdeploy/internal/application/configuration/ports"
	appProject "github.com/jairoprogramador/fastdeploy/internal/application/project/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/factories"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
)

type Initializer struct {
	readerConfig appConfig.Reader
	readerProject appProject.Reader
	writerProject appProject.Writer
	factory      factories.ProjectFactory
	git          ports.Git
	identifier   ports.Identifier
	name         ports.Name
}

func NewInitializer(
	readerConfig appConfig.Reader,
	readerProject appProject.Reader,
	writerProject appProject.Writer,
	factory factories.ProjectFactory,
	git ports.Git,
	identifier ports.Identifier,
	name ports.Name,
) *Initializer {
	return &Initializer{
		readerConfig: readerConfig,
		readerProject: readerProject,
		writerProject: writerProject,
		git:          git,
		identifier:   identifier,
		name:         name,
		factory:      factory,
	}
}

func (ps *Initializer) Initialize() (entities.Project, error) {

	isInitialized, err := ps.IsInitialized()
	if err != nil {
		return entities.Project{}, fmt.Errorf("initialize project failed, check if project is initialized error: %w", err)
	}

	if isInitialized {
		fmt.Println("El proyecto ya ha sido inicializado.")
		project, err := ps.readerProject.Read()
		if err != nil {
			return entities.Project{}, fmt.Errorf("initialize project failed, load project error: %w", err)
		}
		return project, nil
	}

	config, err := ps.readerConfig.Read()
	if err != nil {
		return entities.Project{}, fmt.Errorf("initialize project failed, load config error: %w", err)
	}

	nameProject, err := ps.name.GetName()
	if err != nil {
		return entities.Project{}, fmt.Errorf("initialize project failed, get project name error: %w", err)
	}

	idProject := ps.identifier.Generate(nameProject, config.GetNameOrganization().Value())

	project, err := ps.factory.Create(config, idProject, nameProject)
	if err != nil {
		return entities.Project{}, fmt.Errorf("initialize project failed, create project error: %w", err)
	}

	nameRepository, err := project.GetRepository().GetName()
	if err != nil {
		return entities.Project{}, fmt.Errorf("initialize project failed, get repository name error: %w", err)
	}

	if err = ps.git.Clone(project.GetRepository().GetURL().Value(), nameRepository.Value()); err != nil {
		return entities.Project{}, fmt.Errorf("initialize project failed, clone repository error: %w", err)
	}

	if err = ps.writerProject.Write(project); err != nil {
		return entities.Project{}, fmt.Errorf("initialize project failed, save project error: %w", err)
	}

	return project, nil
}

func (ps *Initializer) IsInitialized() (bool, error) {
	existsFile, err := ps.readerProject.ExistsFile()
	if err != nil {
		return false, fmt.Errorf("is initialized project failed, check if project exists error: %w", err)
	}

	if !existsFile {
		return false, nil
	}

	project, err := ps.readerProject.Read()
	if err != nil {
		return false, fmt.Errorf("is initialized project failed, load project error: %w", err)
	}

	nameRepository, err := project.GetRepository().GetName()
	if err != nil {
		return false, fmt.Errorf("is initialized project failed, get repository name error: %w", err)
	}

	isCloned, err := ps.git.IsCloned(nameRepository.Value())
	if err != nil {
		return false, fmt.Errorf("is initialized project failed, is cloned repository error: %w", err)
	}

	return isCloned, nil
}