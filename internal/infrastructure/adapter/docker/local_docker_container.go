package docker

import (
	"context"
	"fmt"
	engineModel "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/template"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	template2 "github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/template"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
	"net"
	"strings"
)

const (
	// Port range for container
	minPort = 2000
	maxPort = 3000

	// Docker command templates
	dockerComposeUpBuildCmd = "docker compose -f %s up -d --build"
	dockerComposeUpCmd      = "docker compose -f %s up -d"
	dockerComposeDownCmd    = "docker compose -f %s down --rmi local --remove-orphans -v"
	dockerPsIDsAllCmd       = "docker ps -aq --filter ancestor=%s:%s"
	dockerPsIDsUpCmd        = "docker ps -q --filter ancestor=%s:%s"
	dockerPortCmd           = "docker port %s"

	errContainerIDsNotFound  = "container IDs not found"
	errContainerPortNotFound = "container port not found"
	errComposeFileNotFound   = "compose file not found: %s"
	errPortNoAvailable       = "no available ports found between %d and %d"

	msgSuccessfullStart = "successfully started container"
)

type DockerComposeData struct {
	NameDelivery        string // Name of the delivery/project
	CommitHash          string // Git commit hash
	Port                string // Port to expose
	Version             string // Version of the application
	PathDockerDirectory string // Path to the docker directory
	PathHomeDirectory   string // Path to the home directory
}

type localDockerContainer struct {
	runCommand     port.RunCommand
	fileController file.FileController
	dockerTemplate template2.DockerTemplate
	dockerImage    DockerImage
	pathService    port.PathService
	store          *engineModel.StoreEntity
	fileLogger     *logger.FileLogger
}

func NewLocalDockerContainer(
	runCommand port.RunCommand,
	fileController file.FileController,
	dockerTemplate template2.DockerTemplate,
	dockerImage DockerImage,
	pathService port.PathService,
	store *engineModel.StoreEntity,
	fileLogger *logger.FileLogger,
) port.DockerContainer {
	return &localDockerContainer{
		runCommand:     runCommand,
		fileController: fileController,
		dockerTemplate: dockerTemplate,
		dockerImage:    dockerImage,
		pathService:    pathService,
		store:          store,
		fileLogger:     fileLogger,
	}
}

func (d *localDockerContainer) Up(ctx context.Context) result.InfraResult {
	pathDockerCompose := d.pathService.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpCmd, pathDockerCompose)
	return d.runCommand.Run(ctx, command)
}

func (d *localDockerContainer) Start(ctx context.Context) result.InfraResult {
	// Step 1: Stop any existing containers
	response := d.down(ctx)
	if response.Error != nil {
		return response
	}

	// Step 2: createContainer a Dockerfile
	if err := d.dockerImage.CreateDockerfile(); err != nil {
		return d.logError(err)
	}

	// Step 3: createContainer and start a new container
	if err := d.createContainer(ctx); err != nil {
		return d.logError(err)
	}

	return result.NewResult(msgSuccessfullStart)
}

func (d *localDockerContainer) Exists(ctx context.Context, commitHash, version string) result.InfraResult {
	command := fmt.Sprintf(dockerPsIDsAllCmd, commitHash, version)
	response := d.runCommand.Run(ctx, command)

	if !response.IsSuccess() {
		return response
	}

	containerId := response.Result.(string)
	return result.NewResult(len(containerId) > 0)
}

func (d *localDockerContainer) GetURLsUp(ctx context.Context, commitHash, version string) result.InfraResult {
	// Get IDs of running containers
	containerIDsResponse := d.getContainerIDsUp(ctx, commitHash, version)
	if !containerIDsResponse.IsSuccess() {
		return containerIDsResponse
	}

	containerIDs := containerIDsResponse.Result.([]string)
	if len(containerIDs) == 0 {
		return containerIDsResponse
	}

	var urls []string
	for _, containerID := range containerIDs {
		portResponse := d.getContainerPort(ctx, containerID)
		if !portResponse.IsSuccess() {
			return portResponse
		}

		port := portResponse.Result.(string)
		url := fmt.Sprintf("service available in: http://localhost:%s/", port)
		urls = append(urls, url)
	}
	return result.NewResult(urls)
}

func (d *localDockerContainer) upBuild(ctx context.Context) result.InfraResult {
	pathDockerCompose := d.pathService.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpBuildCmd, pathDockerCompose)
	return d.runCommand.Run(ctx, command)
}

func (d *localDockerContainer) down(ctx context.Context) result.InfraResult {
	pathDockerCompose := d.pathService.GetFullPathDockerCompose()

	exists, err := d.fileController.ExistsFile(pathDockerCompose)
	if err != nil {
		return result.NewError(err)
	}

	if !exists {
		return d.logError(fmt.Errorf(errComposeFileNotFound, pathDockerCompose))
	}

	command := fmt.Sprintf(dockerComposeDownCmd, pathDockerCompose)
	return d.runCommand.Run(ctx, command)
}

func (d *localDockerContainer) getContainerIDsUp(ctx context.Context, commitHash, version string) result.InfraResult {
	command := fmt.Sprintf(dockerPsIDsUpCmd, commitHash, version)
	containerIDsResponse := d.runCommand.Run(ctx, command)
	if !containerIDsResponse.IsSuccess() {
		return containerIDsResponse
	}

	containerIDs := containerIDsResponse.Result.(string)
	containerIDs = strings.TrimSpace(containerIDs)
	if containerIDs == "" {
		return d.logError(fmt.Errorf(errContainerIDsNotFound))
	}

	return result.NewResult(strings.Split(containerIDs, "\n"))
}

func (d *localDockerContainer) getContainerPort(ctx context.Context, containerID string) result.InfraResult {
	command := fmt.Sprintf(dockerPortCmd, containerID)
	response := d.runCommand.Run(ctx, command)
	if !response.IsSuccess() {
		return response
	}

	output := response.Result.(string)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		portHost := strings.TrimSpace(parts[1])
		if portHost != "" {
			return result.NewResult(portHost)
		}
	}
	return d.logError(fmt.Errorf(errContainerPortNotFound))
}

func (d *localDockerContainer) createContainer(ctx context.Context) error {
	pathComposeTemplate := d.pathService.GetFullPathDockerComposeTemplate()
	pathCompose := d.pathService.GetFullPathDockerCompose()

	composeData, err := d.prepareComposeData()
	if err != nil {
		return err
	}

	if err := d.ensureTemplateExists(pathComposeTemplate); err != nil {
		return err
	}

	if err := d.removeExistingComposeFile(pathCompose); err != nil {
		return err
	}

	if err := d.generateComposeFile(pathComposeTemplate, pathCompose, composeData); err != nil {
		return err
	}

	return d.startContainer(ctx)
}

func (d *localDockerContainer) prepareComposeData() (DockerComposeData, error) {
	port, err := d.getPort()
	if err != nil {
		return DockerComposeData{}, err
	}
	return DockerComposeData{
		NameDelivery:        d.store.Get(constant.KeyProjectName),
		CommitHash:          d.store.Get(constant.KeyCommitHash),
		Version:             d.store.Get(constant.KeyProjectVersion),
		PathDockerDirectory: d.store.Get(constant.KeyPathDockerDirectory),
		PathHomeDirectory:   d.store.Get(constant.KeyPathHomeDirectory),
		Port:                port,
	}, nil
}

func (d *localDockerContainer) ensureTemplateExists(pathTemplate string) error {
	exists, err := d.fileController.ExistsFile(pathTemplate)
	if err != nil {
		return err
	}

	if !exists {
		if err := d.fileController.WriteFile(pathTemplate, template.ComposeTemplate); err != nil {
			return err
		}
	}

	return nil
}

func (d *localDockerContainer) removeExistingComposeFile(pathDockerCompose string) error {
	exists, err := d.fileController.ExistsFile(pathDockerCompose)
	if err != nil {
		return err
	}

	if exists {
		if err := d.fileController.DeleteFile(pathDockerCompose); err != nil {
			return err
		}
	}

	return nil
}

func (d *localDockerContainer) generateComposeFile(pathTemplate, pathDockerCompose string, composeData DockerComposeData) error {
	contentFile, err := d.dockerTemplate.GetContent(pathTemplate, composeData)
	if err != nil {
		return err
	}

	if err = d.fileController.WriteFile(pathDockerCompose, contentFile); err != nil {
		return err
	}

	return nil
}

func (d *localDockerContainer) startContainer(ctx context.Context) error {
	response := d.upBuild(ctx)
	if !response.IsSuccess() {
		return response.Error
	}
	return nil
}

func (d *localDockerContainer) getPort() (string, error) {
	for port := minPort; port <= maxPort; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)

		if err == nil {
			listener.Close()
			return fmt.Sprintf("%d", port), nil
		} else {
			d.fileLogger.Error(err)
		}
		continue
	}
	err := fmt.Errorf(errPortNoAvailable, minPort, maxPort)
	d.fileLogger.Error(err)
	return "", err
}

func (d *localDockerContainer) logError(err error) result.InfraResult {
	if err != nil {
		d.fileLogger.Error(err)
	}
	return result.NewError(err)
}
