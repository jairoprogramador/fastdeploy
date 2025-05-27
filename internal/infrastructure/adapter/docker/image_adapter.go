package docker

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/repository"
	"github.com/jairoprogramador/fastdeploy/internal/domain/template"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	template2 "github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/template"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

type ImagePort interface {
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

type imageAdapter struct {
	filePort     file.FilePort
	templatePort template2.DockerTemplatePort
	projectRepo  repository.ProjectRepository
	pathPort     port.PathPort
	store        *model.StoreEntity
}

func NewImageAdapter(
	filePort file.FilePort,
	templatePort template2.DockerTemplatePort,
	projectRepo repository.ProjectRepository,
	pathPort port.PathPort,
	store *model.StoreEntity,
) ImagePort {
	return &imageAdapter{
		filePort:     filePort,
		templatePort: templatePort,
		projectRepo:  projectRepo,
		pathPort:     pathPort,
		store:        store,
	}
}

func (docker *imageAdapter) CreateDockerfile() error {
	if err := docker.ensureTemplateExists(); err != nil {
		return err
	}

	pathDockerFile := docker.pathPort.GetFullPathDockerfile()
	if err := docker.prepareDestinationFile(pathDockerFile); err != nil {
		return err
	}

	templateParams, err := docker.createTemplateParameters()
	if err != nil {
		return err
	}

	return docker.generateDockerfile(pathDockerFile, templateParams)
}

func (docker *imageAdapter) ensureTemplateExists() error {
	templatePath := docker.pathPort.GetFullPathDockerfileTemplate()

	exists, err := docker.filePort.ExistsFile(templatePath)
	if err != nil {
		return err
	}

	if !exists {
		return docker.filePort.WriteFile(templatePath, template.DockerfileTemplate)
	}

	return nil
}

func (docker *imageAdapter) prepareDestinationFile(filePath string) error {
	exists, err := docker.filePort.ExistsFile(filePath)
	if err != nil {
		return err
	}

	if exists {
		return docker.filePort.DeleteFile(filePath)
	}

	return nil
}

func (docker *imageAdapter) createTemplateParameters() (DockerfileData, error) {
	response := docker.projectRepo.GetFullPathResource()
	if !response.IsSuccess() {
		return DockerfileData{}, response.Error
	}
	resourcePath := response.Result.(string)

	relativePath := docker.pathPort.GetRelativePathFromHome(resourcePath)

	return DockerfileData{
		FileName:      relativePath,
		CommitMessage: docker.store.Get(constant.KeyCommitMessage),
		CommitHash:    docker.store.Get(constant.KeyCommitHash),
		CommitAuthor:  docker.store.Get(constant.KeyCommitAuthor),
		Team:          docker.store.Get(constant.KeyProjectTeam),
		Organization:  docker.store.Get(constant.KeyProjectOrganization),
	}, nil
}

func (docker *imageAdapter) generateDockerfile(destinationPath string, params DockerfileData) error {
	templatePath := docker.pathPort.GetFullPathDockerfileTemplate()

	content, err := docker.templatePort.GetContent(templatePath, params)
	if err != nil {
		return err
	}

	return docker.filePort.WriteFile(destinationPath, content)
}
