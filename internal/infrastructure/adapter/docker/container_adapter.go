package docker

import (
	"context"
	"fmt"
	engineModel "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
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
	minPort = 2000
	maxPort = 3000

	dockerComposeUpBuildCmd = "docker compose -f %s up -d --build"
	dockerComposeUpCmd      = "docker compose -f %s up -d"
	dockerComposeDownCmd    = "docker compose -f %s down --rmi local --remove-orphans -v"
	dockerPsIDsUpCmd        = "docker ps -q --filter ancestor=%s:%s"
	dockerPsIDsAllCmd       = "docker ps -aq --filter ancestor=%s:%s"
	dockerPortCmd           = "docker port %s"
	dockerDeleteContainer   = "docker rm -f %s"
	dockerDeleteImage       = "docker rmi -f %s:%s"

	errContainerIDsNotFound  = "container IDs not found"
	errContainerPortNotFound = "container port not found"
	errPortNoAvailable       = "no available ports found between %d and %d"

	messageComposeFileNotFound = "compose file not found: %s"
	messageContainerDelete     = "the containers %s were removed"

	successStart = "successfully started container"
)

type ComposeData struct {
	NameDelivery        string
	CommitHash          string
	Port                string
	Version             string
	PathDockerDirectory string
	PathHomeDirectory   string
}

type containerAdapter struct {
	commandPort  port.CommandPort
	filePort     file.FilePort
	templatePort template2.DockerTemplatePort
	imagePort    ImagePort
	pathPort     port.PathPort
	store        *engineModel.StoreEntity
	fileLogger   *logger.FileLogger
}

func NewContainerAdapter(
	commandPort port.CommandPort,
	filePort file.FilePort,
	templatePort template2.DockerTemplatePort,
	imagePort ImagePort,
	pathPort port.PathPort,
	store *engineModel.StoreEntity,
	fileLogger *logger.FileLogger,
) port.ContainerPort {
	return &containerAdapter{
		commandPort:  commandPort,
		filePort:     filePort,
		templatePort: templatePort,
		imagePort:    imagePort,
		pathPort:     pathPort,
		store:        store,
		fileLogger:   fileLogger,
	}
}

func (d *containerAdapter) Up(ctx context.Context) result.InfraResult {
	pathDockerCompose := d.pathPort.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpCmd, pathDockerCompose)
	return d.commandPort.Run(ctx, command)
}

func (d *containerAdapter) Start(ctx context.Context, commitHash, version string) result.InfraResult {
	response := d.down(ctx, commitHash, version)
	if response.Error != nil {
		return response
	}

	if err := d.imagePort.CreateDockerfile(); err != nil {
		return d.logError(err)
	}

	if err := d.createContainer(ctx); err != nil {
		return d.logError(err)
	}

	return result.NewResult(successStart)
}

func (d *containerAdapter) ExistsFileCompose() result.InfraResult {
	pathCompose := d.pathPort.GetFullPathDockerCompose()
	exists, err := d.filePort.ExistsFile(pathCompose)
	if err != nil {
		return result.NewError(err)
	}
	return result.NewResult(exists)
}

func (d *containerAdapter) GetURLsUp(ctx context.Context, commitHash, version string) result.InfraResult {
	containerIDsResponse := d.getContainerIDsUp(ctx, commitHash, version)
	if !containerIDsResponse.IsSuccess() {
		return containerIDsResponse
	}

	containerIDs := containerIDsResponse.Result.([]string)
	if len(containerIDs) == 0 {
		return d.logError(fmt.Errorf(errContainerIDsNotFound))
	}

	var urls []string
	for _, containerID := range containerIDs {
		portResponse := d.getContainerPort(ctx, containerID)
		if !portResponse.IsSuccess() {
			return portResponse
		}

		port := portResponse.Result.(string)
		url := fmt.Sprintf("available in: http://localhost:%s/", port)
		urls = append(urls, url)
	}
	return result.NewResult(urls)
}

func (d *containerAdapter) Exists(ctx context.Context, commitHash, version string) result.InfraResult {
	response := d.getContainerIDsAll(ctx, commitHash, version)

	if response.IsSuccess() {
		containerIds := response.Result.([]string)
		return result.NewResult(len(containerIds) > 0)
	}
	return response
}

func (d *containerAdapter) upBuild(ctx context.Context) result.InfraResult {
	pathDockerCompose := d.pathPort.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpBuildCmd, pathDockerCompose)
	return d.commandPort.Run(ctx, command)
}

func (d *containerAdapter) down(ctx context.Context, commitHash, version string) result.InfraResult {
	pathDockerCompose := d.pathPort.GetFullPathDockerCompose()

	exists, err := d.filePort.ExistsFile(pathDockerCompose)
	if err != nil {
		return result.NewError(err)
	}

	if exists {
		command := fmt.Sprintf(dockerComposeDownCmd, pathDockerCompose)
		return d.commandPort.Run(ctx, command)
	} else {
		response := d.getContainerIDsAll(ctx, commitHash, version)
		containerIDs := response.Result.([]string)
		if len(containerIDs) > 0 {
			return d.deleteAllContainers(ctx, commitHash, version, containerIDs)
		}
		return result.NewResult(fmt.Sprintf(messageComposeFileNotFound, pathDockerCompose))
	}
}

func (d *containerAdapter) deleteAllContainers(ctx context.Context, commitHash string, version string, containerIDs []string) result.InfraResult {
	for _, containerId := range containerIDs {
		respDeleteContainer := d.deleteContainer(ctx, containerId)
		if !respDeleteContainer.IsSuccess() {
			return result.NewError(respDeleteContainer.Error)
		}
	}

	respDeleteImage := d.deleteImage(ctx, commitHash, version)
	if !respDeleteImage.IsSuccess() {
		return result.NewError(respDeleteImage.Error)
	}
	return result.NewResult(fmt.Sprintf(messageContainerDelete, commitHash))
}

func (d *containerAdapter) getContainerIDsAll(ctx context.Context, commitHash, version string) result.InfraResult {
	command := fmt.Sprintf(dockerPsIDsAllCmd, commitHash, version)
	return d.getContainerIDs(ctx, command)
}

func (d *containerAdapter) getContainerIDsUp(ctx context.Context, commitHash, version string) result.InfraResult {
	command := fmt.Sprintf(dockerPsIDsUpCmd, commitHash, version)
	return d.getContainerIDs(ctx, command)
}

func (d *containerAdapter) getContainerIDs(ctx context.Context, command string) result.InfraResult {
	response := d.commandPort.Run(ctx, command)
	if !response.IsSuccess() {
		return response
	}

	containerIDs := response.Result.(string)
	return result.NewResult(d.getArrayContainerIDs(containerIDs))
}

func (d *containerAdapter) getArrayContainerIDs(containerIDs string) []string {
	containerIDs = strings.TrimSpace(containerIDs)
	if containerIDs == "" {
		return []string{}
	}
	return strings.Split(containerIDs, "\n")
}

func (d *containerAdapter) deleteContainer(ctx context.Context, containerID string) result.InfraResult {
	command := fmt.Sprintf(dockerDeleteContainer, containerID)
	return d.commandPort.Run(ctx, command)
}

func (d *containerAdapter) deleteImage(ctx context.Context, commitHash, version string) result.InfraResult {
	command := fmt.Sprintf(dockerDeleteImage, commitHash, version)
	return d.commandPort.Run(ctx, command)
}

func (d *containerAdapter) getContainerPort(ctx context.Context, containerID string) result.InfraResult {
	command := fmt.Sprintf(dockerPortCmd, containerID)
	response := d.commandPort.Run(ctx, command)
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

func (d *containerAdapter) createContainer(ctx context.Context) error {
	pathComposeTemplate := d.pathPort.GetFullPathDockerComposeTemplate()
	pathCompose := d.pathPort.GetFullPathDockerCompose()

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

func (d *containerAdapter) prepareComposeData() (ComposeData, error) {
	port, err := d.getPort()
	if err != nil {
		return ComposeData{}, err
	}
	return ComposeData{
		NameDelivery:        d.store.Get(constant.KeyProjectName),
		CommitHash:          d.store.Get(constant.KeyCommitHash),
		Version:             d.store.Get(constant.KeyProjectVersion),
		PathDockerDirectory: d.store.Get(constant.KeyPathDockerDirectory),
		PathHomeDirectory:   d.store.Get(constant.KeyPathHomeDirectory),
		Port:                port,
	}, nil
}

func (d *containerAdapter) ensureTemplateExists(pathTemplate string) error {
	exists, err := d.filePort.ExistsFile(pathTemplate)
	if err != nil {
		return err
	}

	if !exists {
		if err := d.filePort.WriteFile(pathTemplate, template.ComposeTemplate); err != nil {
			return err
		}
	}

	return nil
}

func (d *containerAdapter) removeExistingComposeFile(pathDockerCompose string) error {
	exists, err := d.filePort.ExistsFile(pathDockerCompose)
	if err != nil {
		return err
	}

	if exists {
		if err := d.filePort.DeleteFile(pathDockerCompose); err != nil {
			return err
		}
	}

	return nil
}

func (d *containerAdapter) generateComposeFile(pathTemplate, pathDockerCompose string, composeData ComposeData) error {
	contentFile, err := d.templatePort.GetContent(pathTemplate, composeData)
	if err != nil {
		return err
	}

	if err = d.filePort.WriteFile(pathDockerCompose, contentFile); err != nil {
		return err
	}

	return nil
}

func (d *containerAdapter) startContainer(ctx context.Context) error {
	response := d.upBuild(ctx)
	if !response.IsSuccess() {
		return response.Error
	}
	return nil
}

func (d *containerAdapter) getPort() (string, error) {
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

func (d *containerAdapter) logError(err error) result.InfraResult {
	if err != nil {
		d.fileLogger.Error(err)
	}
	return result.NewError(err)
}
