package adapter

import (
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/port"
	"deploy/internal/domain/service"
	"deploy/internal/domain/service/router"
	"deploy/internal/domain/template"
)

type DockerfileData struct {
	FileName      string
	CommitMessage string
	CommitHash    string
	CommitAuthor  string
	Team          string
	Organization  string
}

type localDockerImage struct {
	executorService port.ExecutorServiceInterface
	fileRepository  FileRepository
	dockerTemplate  port.DockerTemplate
	projectService  service.ProjectService
	router          *router.Router
	variables       *model.VariableStore
}

func NewLocalDockerImage(
	executorService port.ExecutorServiceInterface,
	fileRepository FileRepository,
	templateService port.DockerTemplate,
	projectService service.ProjectService,
	router *router.Router,
	variables *model.VariableStore,
) port.DockerImage {
	return &localDockerImage{
		executorService: executorService,
		fileRepository:  fileRepository,
		dockerTemplate:  templateService,
		projectService:  projectService,
		router:          router,
		variables:       variables,
	}
}

func (d *localDockerImage) CreateDockerfile() error {
	pathTemplateDockerfile := d.router.GetFullPathDockerfileTemplate()
	exists, err := d.fileRepository.ExistsFile(pathTemplateDockerfile)
	if err != nil {
		return err
	}
	if !exists {
		err := d.fileRepository.WriteFile(
			pathTemplateDockerfile, template.DockerfileTemplate)
		if err != nil {
			return err
		}
	}

	pathDockerFile := d.router.GetFullPathDockerfile()
	exists, err = d.fileRepository.ExistsFile(pathDockerFile)
	if err != nil {
		return err
	}
	if exists {
		if err := d.fileRepository.DeleteFile(pathDockerFile); err != nil {
			return err
		}
	}

	nameResource, err := d.projectService.GetFullPathResource()
	if err != nil {
		return err
	}

	nameResource = d.router.GetRelativePathFromHome(nameResource)

	params := DockerfileData{
		FileName:      nameResource,
		CommitMessage: d.variables.Get(constant.VAR_COMMIT_MESSAGE),
		CommitHash:    d.variables.Get(constant.VAR_COMMIT_HASH),
		CommitAuthor:  d.variables.Get(constant.VAR_COMMIT_AUTHOR),
		Team:          d.variables.Get(constant.VAR_PROJECT_TEAM),
		Organization:  d.variables.Get(constant.VAR_PROJECT_ORGANIZATION),
	}

	contentFile, err := d.dockerTemplate.GetContent(pathTemplateDockerfile, params)
	if err != nil {
		return err
	}

	if err = d.fileRepository.WriteFile(pathDockerFile, contentFile); err != nil {
		return err
	}
	return nil
}
