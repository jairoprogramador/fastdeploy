package adapter

import (
	"context"
	"deploy/internal/domain/constant"
	model2 "deploy/internal/domain/engine/model"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/port"
	"deploy/internal/domain/service"
	"deploy/internal/domain/template"
	"fmt"
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

	errComposeFileNotFound = "file compose not found in %s"
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
	fileRepository FileController
	dockerTemplate DockerTemplate
	dockerImage    DockerImage
	router         *service.PathService
	variables      *model2.VariableStore
	logger         *logger.Logger
}

func NewLocalDockerContainer(
	runCommand port.RunCommand,
	fileRepository FileController,
	dockerTemplate DockerTemplate,
	dockerImage DockerImage,
	router *service.PathService,
	variables *model2.VariableStore,
	logger *logger.Logger,
) port.DockerContainer {
	return &localDockerContainer{
		runCommand:     runCommand,
		fileRepository: fileRepository,
		dockerTemplate: dockerTemplate,
		dockerImage:    dockerImage,
		router:         router,
		variables:      variables,
		logger:         logger,
	}
}

func (d *localDockerContainer) Up(ctx context.Context) model.InfrastructureResponse {
	pathDockerCompose := d.router.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpCmd, pathDockerCompose)
	return d.runCommand.Run(ctx, command)
}

func (d *localDockerContainer) Start(ctx context.Context) error {
	// Step 1: Stop any existing containers
	result := d.down(ctx)
	if result.Error != nil {
		d.logger.Error(result.Error)
		return result.Error
	}

	// Step 2: createContainer a Dockerfile
	if err := d.dockerImage.CreateDockerfile(); err != nil {
		d.logger.Error(err)
		return err
	}

	// Step 3: createContainer and start a new container
	if err := d.createContainer(ctx); err != nil {
		d.logger.Error(err)
		return err
	}

	return nil
}

func (d *localDockerContainer) Exists(ctx context.Context, commitHash, version string) model.InfrastructureResponse {
	command := fmt.Sprintf(dockerPsIDsAllCmd, commitHash, version)
	response := d.runCommand.Run(ctx, command)

	if !response.IsSuccess() {
		return response
	}

	containerId := response.Result.(string)
	return model.NewResponseWithDetails(len(containerId) > 0, response.Details)
}

func (d *localDockerContainer) GetURLsUp(ctx context.Context, commitHash, version string) model.InfrastructureResponse {
	// Get IDs of running containers
	containerIDsResponse := d.getContainerIDsUp(ctx, commitHash, version)
	if !containerIDsResponse.IsSuccess() {
		return containerIDsResponse
	}

	containerIDs := containerIDsResponse.Result.([]string)
	if len(containerIDs) == 0 {
		return containerIDsResponse
	}

	// Build details and collect URLs
	var details strings.Builder
	details.WriteString(containerIDsResponse.Details)

	var urls []string
	for _, containerID := range containerIDs {
		portResponse := d.getContainerPort(ctx, containerID)
		if !portResponse.IsSuccess() {
			return portResponse
		}

		port := portResponse.Result.(string)
		url := fmt.Sprintf("service available in: http://localhost:%s/", port)
		urls = append(urls, url)
		details.WriteString(fmt.Sprintf("\n%s", portResponse.Details))
	}

	return model.NewResponseWithDetails(urls, details.String())
}

func (d *localDockerContainer) upBuild(ctx context.Context) model.InfrastructureResponse {
	pathDockerCompose := d.router.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpBuildCmd, pathDockerCompose)
	return d.runCommand.Run(ctx, command)
}

func (d *localDockerContainer) down(ctx context.Context) model.InfrastructureResponse {
	pathDockerCompose := d.router.GetFullPathDockerCompose()

	exists, err := d.fileRepository.ExistsFile(pathDockerCompose)
	if err != nil {
		return model.NewErrorResponse(err)
	}

	if !exists {
		err = fmt.Errorf(errComposeFileNotFound, pathDockerCompose)
		return model.NewErrorResponse(err)
	}

	command := fmt.Sprintf(dockerComposeDownCmd, pathDockerCompose)
	return d.runCommand.Run(ctx, command)
}

func (d *localDockerContainer) getContainerIDsUp(ctx context.Context, commitHash, version string) model.InfrastructureResponse {
	command := fmt.Sprintf(dockerPsIDsUpCmd, commitHash, version)
	containerIDsResponse := d.runCommand.Run(ctx, command)
	if !containerIDsResponse.IsSuccess() {
		return containerIDsResponse
	}

	containerIDs := containerIDsResponse.Result.(string)
	containerIDs = strings.TrimSpace(containerIDs)
	if containerIDs == "" {
		return model.NewResponseWithDetails([]string{}, containerIDsResponse.Details)
	}

	return model.NewResponseWithDetails(strings.Split(containerIDs, "\n"), containerIDsResponse.Details)
}

func (d *localDockerContainer) getContainerPort(ctx context.Context, containerID string) model.InfrastructureResponse {
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
			return model.NewResponseWithDetails(portHost, response.Details)
		}
	}

	return model.NewErrorResponseWithDetails(fmt.Errorf(constant.MsgErrorNoPortHost), response.Details)
}

func (d *localDockerContainer) createContainer(ctx context.Context) error {
	// Get paths
	pathComposeTemplate := d.router.GetFullPathDockerComposeTemplate()
	pathCompose := d.router.GetFullPathDockerCompose()

	// Prepare compose data
	composeData := d.prepareComposeData()

	// Ensure template exists
	if err := d.ensureTemplateExists(pathComposeTemplate); err != nil {
		return err
	}

	// Remove existing compose file if it exists
	if err := d.removeExistingComposeFile(pathCompose); err != nil {
		return err
	}

	// Generate compose file from template
	if err := d.generateComposeFile(pathComposeTemplate, pathCompose, composeData); err != nil {
		return err
	}

	// Start the container
	return d.startContainer(ctx)
}

func (d *localDockerContainer) prepareComposeData() DockerComposeData {
	return DockerComposeData{
		NameDelivery:        d.variables.Get(constant.VAR_PROJECT_NAME),
		CommitHash:          d.variables.Get(constant.VAR_COMMIT_HASH),
		Version:             d.variables.Get(constant.VAR_PROJECT_VERSION),
		PathDockerDirectory: d.variables.Get(constant.VAR_PATH_DOCKER_DIRECTORY),
		PathHomeDirectory:   d.variables.Get(constant.VAR_PATH_HOME_DIRECTORY),
		Port:                d.getPort(),
	}
}

func (d *localDockerContainer) ensureTemplateExists(pathTemplate string) error {
	exists, err := d.fileRepository.ExistsFile(pathTemplate)
	if err != nil {
		return err
	}

	if !exists {
		if err := d.fileRepository.WriteFile(pathTemplate, template.ComposeTemplate); err != nil {
			return err
		}
	}

	return nil
}

func (d *localDockerContainer) removeExistingComposeFile(pathDockerCompose string) error {
	exists, err := d.fileRepository.ExistsFile(pathDockerCompose)
	if err != nil {
		return err
	}

	if exists {
		if err := d.fileRepository.DeleteFile(pathDockerCompose); err != nil {
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

	if err = d.fileRepository.WriteFile(pathDockerCompose, contentFile); err != nil {
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

func (d *localDockerContainer) getPort() string {
	defaultPort := minPort

	// Try to find an available port
	for port := minPort; port <= maxPort; port++ {
		address := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", address)

		if err == nil {
			// Found an available port
			listener.Close()
			return fmt.Sprintf("%d", port)
		}
	}

	// If no port is available, return the default
	d.logger.Info(fmt.Sprintf("No available ports found between %d and %d, using default: %d",
		minPort, maxPort, defaultPort))
	return fmt.Sprintf("%d", defaultPort)
}
