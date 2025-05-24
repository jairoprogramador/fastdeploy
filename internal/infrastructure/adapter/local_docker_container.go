package adapter

import (
	"context"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/port"
	"deploy/internal/domain/service/router"
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

	// Error messages
	errComposeFileNotFound = "file compose not found in %s"
)

// DockerComposeData holds the data needed to generate a docker-compose file
type DockerComposeData struct {
	NameDelivery        string // Name of the delivery/project
	CommitHash          string // Git commit hash
	Port                string // Port to expose
	Version             string // Version of the application
	PathDockerDirectory string // Path to the docker directory
	PathHomeDirectory   string // Path to the home directory
}

// localDockerContainer implements the port.DockerContainer interface
type localDockerContainer struct {
	executorService port.ExecutorServiceInterface
	fileRepository  FileRepository
	dockerTemplate  port.DockerTemplate
	dockerImage     port.DockerImage
	router          *router.Router
	variables       *model.VariableStore
	logger          *logger.Logger
}

// NewLocalDockerContainer creates a new instance of localDockerContainer
func NewLocalDockerContainer(
	executorService port.ExecutorServiceInterface,
	fileRepository FileRepository,
	templateService port.DockerTemplate,
	dockerImage port.DockerImage,
	router *router.Router,
	variables *model.VariableStore,
	logger *logger.Logger,
) port.DockerContainer {
	return &localDockerContainer{
		executorService: executorService,
		fileRepository:  fileRepository,
		dockerTemplate:  templateService,
		dockerImage:     dockerImage,
		router:          router,
		variables:       variables,
		logger:          logger,
	}
}

// UpBuild starts containers with docker-compose up --build
func (d *localDockerContainer) UpBuild(ctx context.Context) model.InfrastructureResponse {
	pathDockerCompose := d.router.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpBuildCmd, pathDockerCompose)
	return d.executorService.Run(ctx, command)
}

// Up starts containers with docker-compose up
func (d *localDockerContainer) Up(ctx context.Context) model.InfrastructureResponse {
	pathDockerCompose := d.router.GetFullPathDockerCompose()
	command := fmt.Sprintf(dockerComposeUpCmd, pathDockerCompose)
	return d.executorService.Run(ctx, command)
}

// Down stops and removes containers with docker-compose down
func (d *localDockerContainer) Down(ctx context.Context) model.InfrastructureResponse {
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
	return d.executorService.Run(ctx, command)
}

// Exists checks if a container with the given commit hash and version exists
func (d *localDockerContainer) Exists(ctx context.Context, commitHash, version string) model.InfrastructureResponse {
	command := fmt.Sprintf(dockerPsIDsAllCmd, commitHash, version)
	response := d.executorService.Run(ctx, command)

	if !response.IsSuccess() {
		return response
	}

	containerId := response.Result.(string)
	return model.NewResponseWithDetails(len(containerId) > 0, response.Details)
}

// GetURLs returns the URLs for all running containers with the given commit hash and version
func (d *localDockerContainer) GetURLs(ctx context.Context, commitHash, version string) model.InfrastructureResponse {
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

// getContainerIDsUp returns the IDs of all running containers with the given commit hash and version
func (d *localDockerContainer) getContainerIDsUp(ctx context.Context, commitHash, version string) model.InfrastructureResponse {
	command := fmt.Sprintf(dockerPsIDsUpCmd, commitHash, version)
	containerIDsResponse := d.executorService.Run(ctx, command)
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

// getContainerPort returns the host port for a container
func (d *localDockerContainer) getContainerPort(ctx context.Context, containerID string) model.InfrastructureResponse {
	command := fmt.Sprintf(dockerPortCmd, containerID)
	response := d.executorService.Run(ctx, command)
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

// Create implements the Create method of the DockerContainer interface
// It creates and starts a Docker container using a docker-compose file
func (d *localDockerContainer) Create(ctx context.Context) error {
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

// prepareComposeData creates a DockerComposeData struct with values from variables
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

// ensureTemplateExists checks if the template file exists and creates it if not
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

// removeExistingComposeFile removes the compose file if it exists
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

// generateComposeFile creates a docker-compose file from a template
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

// startContainer builds and starts the container
func (d *localDockerContainer) startContainer(ctx context.Context) error {
	response := d.UpBuild(ctx)
	if !response.IsSuccess() {
		return response.Error
	}

	return nil
}

// Start implements the Start method of the DockerContainer interface
// It stops any existing containers, creates a Dockerfile, and starts a new container
func (d *localDockerContainer) Start(ctx context.Context) error {
	// Step 1: Stop any existing containers
	result := d.Down(ctx)
	if result.Error != nil {
		d.logger.Error(result.Error)
		return result.Error
	}

	// Step 2: Create a Dockerfile
	if err := d.dockerImage.CreateDockerfile(); err != nil {
		d.logger.Error(err)
		return err
	}

	// Step 3: Create and start a new container
	if err := d.Create(ctx); err != nil {
		d.logger.Error(err)
		return err
	}

	return nil
}

// getPort finds an available port between minPort and maxPort
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
