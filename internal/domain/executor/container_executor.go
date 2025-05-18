package executor

import (
	"context"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/router"
	"deploy/internal/domain/service"
	"deploy/internal/domain/template"
	"fmt"
	"net"
)

type DockerfileData struct {
	FileName      string
	CommitMessage string
	CommitHash    string
	CommitAuthor  string
	Team          string
	Organization  string
}

type DockerComposeData struct {
	NameDelivery 		string
	CommitHash   		string
	Port                string
	Version             string
	PathDockerDirectory string
	PathHomeDirectory   string
}

type ContainerExecutor struct {
	baseExecutor        *BaseExecutor
	variables           *model.VariableStore
	dockerService       service.DockerServiceInterface
	containerRepository repository.ContainerRepository
	fileRepository      repository.FileRepository
	router              *router.Router
}

func GetContainerExecutor(
	containerRepository repository.ContainerRepository,
	fileRepository repository.FileRepository,
	variables *model.VariableStore) *ContainerExecutor {
	return &ContainerExecutor{
		baseExecutor:        GetBaseExecutor(),
		variables:           variables,
		dockerService:       service.GetDockerService(),
		containerRepository: containerRepository,
		fileRepository:      fileRepository,
		router:              router.GetRouter(),
	}
}

func (e *ContainerExecutor) Execute(ctx context.Context, step model.Step) (string, error) {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() (string, error) {
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		err := e.delete(ctx)
		if err != nil {
			return "", err
		}

		err = e.buildDockerfile()
		if err != nil {
			return "", err
		}

		return e.createContainer(ctx)
	})
}

func (e *ContainerExecutor) delete(ctx context.Context) error {
	pathDockerCompose := e.router.GetFullPathDockerCompose()
	if e.fileRepository.ExistsFile(pathDockerCompose) {
		if err := e.dockerService.DockerComposeDownLocal(ctx, pathDockerCompose); err != nil {
			return err
		}
	}
	return nil
}

func (e *ContainerExecutor) buildDockerfile() error {
	pathDockerfileTemplate := e.router.GetFullPathDockerfileTemplate()
	if !e.fileRepository.ExistsFile(pathDockerfileTemplate) {
		err := e.fileRepository.WriteFile(
			pathDockerfileTemplate, template.DockerfileTemplate)
		if err != nil {
			return err
		}
	}
	pathDockerfile := e.router.GetFullPathDockerfile()
	if e.fileRepository.ExistsFile(pathDockerfile) {
		err := e.fileRepository.DeleteFile(pathDockerfile)
		if err != nil {
			return err
		}
	}

	return e.createDockerfile(pathDockerfile, pathDockerfileTemplate)
}

func (e *ContainerExecutor) createDockerfile(pathDockerfile, pathDockerfileTemplate string) error {
	nameResource, err := e.containerRepository.GetFullPathResource()
	if err != nil {
		return err
	}

	nameResource = e.router.GetRelativePathFromHome(nameResource)

	params := DockerfileData{
		FileName:      nameResource,
		CommitMessage: e.variables.Get(constant.VAR_COMMIT_MESSAGE),
		CommitHash:    e.variables.Get(constant.VAR_COMMIT_HASH),
		CommitAuthor:  e.variables.Get(constant.VAR_COMMIT_AUTHOR),
		Team:          e.variables.Get(constant.VAR_PROJECT_TEAM),
		Organization:  e.variables.Get(constant.VAR_PROJECT_ORGANIZATION),
	}

	contentDockerfile, err := e.containerRepository.GetContentTemplate(pathDockerfileTemplate, params)
	if err != nil {
		return err
	}

	return e.fileRepository.WriteFile(pathDockerfile, contentDockerfile)
}

func (e *ContainerExecutor) createContainer(ctx context.Context) (string, error) {
	pathDockerComposeTemplate := e.router.GetFullPathDockerComposeTemplate()
	if !e.fileRepository.ExistsFile(pathDockerComposeTemplate) {
		err := e.fileRepository.WriteFile(pathDockerComposeTemplate, template.ComposeTemplate)
		if err != nil {
			return "", err
		}
	}
	pathDockerCompose := e.router.GetFullPathDockerCompose()
	if e.fileRepository.ExistsFile(pathDockerCompose) {
		err := e.fileRepository.DeleteFile(pathDockerCompose)
		if err != nil {
			return "", err
		}
	}

	err := e.createDockerCompose(pathDockerCompose, pathDockerComposeTemplate)
	if err != nil {
		return "", err
	}

	return e.dockerService.DockerComposeUp(ctx, pathDockerCompose, e.variables)
}

func (e *ContainerExecutor) createDockerCompose(pathDockerCompose, pathDockerComposeTemplate string) error {
	params := DockerComposeData{
		NameDelivery: e.variables.Get(constant.VAR_PROJECT_NAME),
		CommitHash:   e.variables.Get(constant.VAR_COMMIT_HASH),
		Version:      e.variables.Get(constant.VAR_PROJECT_VERSION),
		PathDockerDirectory: e.variables.Get(constant.VAR_PATH_DOCKER_DIRECTORY),
		PathHomeDirectory:   e.variables.Get(constant.VAR_PATH_HOME_DIRECTORY),
		Port:         e.getPort(),
	}

	contentDockerCompose, err := e.containerRepository.GetContentTemplate(pathDockerComposeTemplate, params)
	if err != nil {
		return err
	}

	return e.fileRepository.WriteFile(pathDockerCompose, contentDockerCompose)
}

func (e *ContainerExecutor) getPort() string {
	startPort := 2000
	endPort := 3000

	portFree := 2000

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", address)
		if err == nil {
			portFree = port
			ln.Close()
			break
		}
	}
	return fmt.Sprintf("%d", portFree)
}
