package adapter

import (
	"deploy/internal/domain/constant"
	"deploy/internal/domain/engine/model"
	"deploy/internal/domain/service"
	"deploy/internal/domain/template"
)

// DockerImage defines operations for Docker image creation
type DockerImage interface {
	CreateDockerfile() error
}

// DockerfileData contains all parameters needed for Dockerfile template
type DockerfileData struct {
	FileName      string
	CommitMessage string
	CommitHash    string
	CommitAuthor  string
	Team          string
	Organization  string
}

// localDockerImage implements DockerImage interface
type localDockerImage struct {
	fileRepository FileController
	dockerTemplate DockerTemplate
	projectService service.ProjectService
	router         *service.PathService
	variables      *model.StoreEntity
}

// NewLocalDockerImage creates a new instance of DockerImage
func NewLocalDockerImage(
	fileRepository FileController,
	dockerTemplate DockerTemplate,
	projectService service.ProjectService,
	router *service.PathService,
	variables *model.StoreEntity,
) DockerImage {
	return &localDockerImage{
		fileRepository: fileRepository,
		dockerTemplate: dockerTemplate,
		projectService: projectService,
		router:         router,
		variables:      variables,
	}
}

// CreateDockerfile generates a Dockerfile based on templates and project configuration
func (docker *localDockerImage) CreateDockerfile() error {
	// Ensure template exists
	if err := docker.ensureTemplateExists(); err != nil {
		return err
	}

	// Prepare destination file
	pathDockerFile := docker.router.GetFullPathDockerfile()
	if err := docker.prepareDestinationFile(pathDockerFile); err != nil {
		return err
	}

	// Generate template parameters
	templateParams, err := docker.createTemplateParameters()
	if err != nil {
		return err
	}

	// Generate and write Dockerfile content
	return docker.generateDockerfile(pathDockerFile, templateParams)
}

// ensureTemplateExists checks if template exists and creates it if needed
func (docker *localDockerImage) ensureTemplateExists() error {
	templatePath := docker.router.GetFullPathDockerfileTemplate()

	exists, err := docker.fileRepository.ExistsFile(templatePath)
	if err != nil {
		return err
	}

	if !exists {
		return docker.fileRepository.WriteFile(templatePath, template.DockerfileTemplate)
	}

	return nil
}

// prepareDestinationFile ensures the destination file is ready for writing
func (docker *localDockerImage) prepareDestinationFile(filePath string) error {
	exists, err := docker.fileRepository.ExistsFile(filePath)
	if err != nil {
		return err
	}

	if exists {
		return docker.fileRepository.DeleteFile(filePath)
	}

	return nil
}

// createTemplateParameters builds the parameters needed for the Dockerfile template
func (docker *localDockerImage) createTemplateParameters() (DockerfileData, error) {
	resourcePath, err := docker.projectService.GetFullPathResource()
	if err != nil {
		return DockerfileData{}, err
	}

	relativePath := docker.router.GetRelativePathFromHome(resourcePath)

	return DockerfileData{
		FileName:      relativePath,
		CommitMessage: docker.variables.Get(constant.VAR_COMMIT_MESSAGE),
		CommitHash:    docker.variables.Get(constant.VAR_COMMIT_HASH),
		CommitAuthor:  docker.variables.Get(constant.VAR_COMMIT_AUTHOR),
		Team:          docker.variables.Get(constant.VAR_PROJECT_TEAM),
		Organization:  docker.variables.Get(constant.VAR_PROJECT_ORGANIZATION),
	}, nil
}

// generateDockerfile creates the final Dockerfile from template and parameters
func (docker *localDockerImage) generateDockerfile(destinationPath string, params DockerfileData) error {
	templatePath := docker.router.GetFullPathDockerfileTemplate()

	content, err := docker.dockerTemplate.GetContent(templatePath, params)
	if err != nil {
		return err
	}

	return docker.fileRepository.WriteFile(destinationPath, content)
}
