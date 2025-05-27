package docker

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/template"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	template2 "github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/template"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
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
	fileController file.FileController
	dockerTemplate template2.DockerTemplate
	projectService service.ProjectService
	pathService    port.PathService
	store          *entity.StoreEntity
}

func NewLocalDockerImage(
	fileController file.FileController,
	dockerTemplate template2.DockerTemplate,
	projectService service.ProjectService,
	pathService port.PathService,
	store *entity.StoreEntity,
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
		CommitMessage: docker.store.Get(constant.KeyCommitMessage),
		CommitHash:    docker.store.Get(constant.KeyCommitHash),
		CommitAuthor:  docker.store.Get(constant.KeyCommitAuthor),
		Team:          docker.store.Get(constant.KeyProjectTeam),
		Organization:  docker.store.Get(constant.KeyProjectOrganization),
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
