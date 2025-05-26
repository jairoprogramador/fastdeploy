package adapter

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/template"
)

type DockerImage interface {
	CreateDockerfile() error
}

type DockerfileData struct {
	FileName      string
	CommitMessage string
	CommitHash    string
	CommitAuthor  string
	Team          string
	Organization  string
}

type localDockerImage struct {
	fileController FileController
	dockerTemplate DockerTemplate
	projectService service.ProjectService
	pathService    port.PathService
	store          *model.StoreEntity
}

func NewLocalDockerImage(
	fileController FileController,
	dockerTemplate DockerTemplate,
	projectService service.ProjectService,
	pathService port.PathService,
	store *model.StoreEntity,
) DockerImage {
	return &localDockerImage{
		fileController: fileController,
		dockerTemplate: dockerTemplate,
		projectService: projectService,
		pathService:    pathService,
		store:          store,
	}
}

func (docker *localDockerImage) CreateDockerfile() error {
	if err := docker.ensureTemplateExists(); err != nil {
		return err
	}

	pathDockerFile := docker.pathService.GetFullPathDockerfile()
	if err := docker.prepareDestinationFile(pathDockerFile); err != nil {
		return err
	}

	templateParams, err := docker.createTemplateParameters()
	if err != nil {
		return err
	}

	return docker.generateDockerfile(pathDockerFile, templateParams)
}

func (docker *localDockerImage) ensureTemplateExists() error {
	templatePath := docker.pathService.GetFullPathDockerfileTemplate()

	exists, err := docker.fileController.ExistsFile(templatePath)
	if err != nil {
		return err
	}

	if !exists {
		return docker.fileController.WriteFile(templatePath, template.DockerfileTemplate)
	}

	return nil
}

func (docker *localDockerImage) prepareDestinationFile(filePath string) error {
	exists, err := docker.fileController.ExistsFile(filePath)
	if err != nil {
		return err
	}

	if exists {
		return docker.fileController.DeleteFile(filePath)
	}

	return nil
}

func (docker *localDockerImage) createTemplateParameters() (DockerfileData, error) {
	resourcePath, err := docker.projectService.GetFullPathResource()
	if err != nil {
		return DockerfileData{}, err
	}

	relativePath := docker.pathService.GetRelativePathFromHome(resourcePath)

	return DockerfileData{
		FileName:      relativePath,
		CommitMessage: docker.store.Get(constant.VAR_COMMIT_MESSAGE),
		CommitHash:    docker.store.Get(constant.VAR_COMMIT_HASH),
		CommitAuthor:  docker.store.Get(constant.VAR_COMMIT_AUTHOR),
		Team:          docker.store.Get(constant.VAR_PROJECT_TEAM),
		Organization:  docker.store.Get(constant.VAR_PROJECT_ORGANIZATION),
	}, nil
}

func (docker *localDockerImage) generateDockerfile(destinationPath string, params DockerfileData) error {
	templatePath := docker.pathService.GetFullPathDockerfileTemplate()

	content, err := docker.dockerTemplate.GetContent(templatePath, params)
	if err != nil {
		return err
	}

	return docker.fileController.WriteFile(destinationPath, content)
}
