package tools

import (
	"context"
	constants "deploy/internal/domain"
	"deploy/internal/domain/repository"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"
)

const (
	dockerCmd        = "docker"
	dockerComposeCmd = "docker compose"
)

var (
	ErrImageNotFound     = errors.New("imagen docker no encontrada")
	ErrContainerNotFound = errors.New("contenedor no encontrado")
	ErrInvalidID         = errors.New("ID inválido o vacío")
)

var (
	dockerfileTemplate *template.Template
	composeTemplate    *template.Template
)

type DockerfileParams struct {
	FileName      string
	CommitMessage string
	CommitHash    string
	CommitAuthor  string
	Team          string
	Organization  string
}

type DockerComposeParams struct {
	NameDelivery string
	CommitHash   string
	Port         string
}

// DockerService implementa la interfaz repository.DockerRepository
type DockerService struct{}

// NewDockerService crea una nueva instancia de DockerService
func NewDockerService() repository.DockerRepository {
	return &DockerService{}
}

// GetImageID obtiene el ID de una imagen Docker basada en el hash del commit
func (d *DockerService) GetImageID(hashCommit string) (string, error) {
	if hashCommit == "" {
		return "", ErrInvalidID
	}

	ctx, cancel := d.GetContext()
	defer cancel()
	imageId, err := ExecuteCommandWithContext(ctx, dockerCmd, "images", "-q", hashCommit)
	if err != nil {
		return "", ErrImageNotFound
	}

	imageId = strings.TrimSpace(imageId)
	if imageId == "" {
		return "", ErrImageNotFound
	}

	return strings.TrimSpace(imageId), nil
}

func (d *DockerService) GetContainersID(imageID string) ([]string, error) {
	if imageID == "" {
		return []string{}, ErrInvalidID
	}

	ctx, cancel := d.GetContext()
	defer cancel()
	ancestor := "ancestor=" + imageID
	containerIds, err := ExecuteCommandWithContext(ctx, dockerCmd, "ps", "-qa", "--filter", ancestor)
	if err != nil {
		return []string{}, ErrContainerNotFound
	}

	containerIds = strings.TrimSpace(containerIds)
	if containerIds == "" {
		return []string{}, ErrContainerNotFound
	}

	return getArray(containerIds), nil
}

// SonarScanner ejecuta el análisis de SonarQube en un contenedor Docker
func (d *DockerService) SonarScanner(token, projectKey, projectName, projectPath, cacheDir, tmpDir, scannerWorkDir, sourcePath, testPath, binaryPath, testBinaryPath string) error {
	args := []string{
		"run",
		"--rm",
		"--network=host",
		"-v", tmpDir + ":/opt/sonar-scanner/.sonar/_tmp",
		"-e", "SONAR_HOST_URL=http://localhost:9000",
		"-e", "SONAR_SCANNER_OPTS=-Xmx1024m -Djava.io.tmpdir=/opt/sonar-scanner/.sonar/_tmp",
		"-v", projectPath + ":/usr/src",
		"-v", cacheDir + ":/opt/sonar-scanner/.sonar/cache",
		"-v", scannerWorkDir + ":/opt/sonar-scanner/.scannerwork",
		"sonarsource/sonar-scanner-cli:latest",
		"-Dsonar.token=" + token,
		"-Dsonar.projectKey=" + projectKey,
		"-Dsonar.projectName=" + projectName,
		"-Dsonar.sources=" + sourcePath,
		"-Dsonar.tests=" + testPath,
		"-Dsonar.java.binaries=" + binaryPath,          //solo para java
		"-Dsonar.java.test.binaries=" + testBinaryPath, //solo para java
		"-Dsonar.sourceEncoding=UTF-8",
		"-Dsonar.scm.provider=git",
		"-Dsonar.tempFolder=/opt/sonar-scanner/.sonar/_tmp",
		"-Dsonar.working.directory=/opt/sonar-scanner/.scannerwork",
	}
	_, err := ExecuteCommand("docker", args...)
	return err
}

// GetPortContainer obtiene el puerto de un contenedor
func (d *DockerService) GetPortContainer(containerId string) (string, error) {
	return ExecuteCommand(dockerCmd, "port", containerId)
}

// BuildImage construye una imagen Docker
func (d *DockerService) BuildImage(hashCommit string, filePath string) error {
	_, err := ExecuteCommand(dockerCmd, "build", "-t", hashCommit, "-f", filePath, ".")
	return err
}

// BuildContainer construye un contenedor usando docker-compose
func (d *DockerService) BuildContainer(filePath string) error {
	_, err := ExecuteCommand(dockerComposeCmd, "-f", filePath, "up", "-d")
	return err
}

// StartContainer inicia un contenedor Docker
func (d *DockerService) StartContainer(containerID string) error {
	if containerID == "" {
		return ErrInvalidID
	}
	ctx, cancel := d.GetContext()
	defer cancel()
	_, err := ExecuteCommandWithContext(ctx, dockerCmd, "start", containerID)
	return err
}

// RestartContainer reinicia un contenedor Docker
func (d *DockerService) RestartContainer(containerID string) error {
	if containerID == "" {
		return ErrInvalidID
	}
	ctx, cancel := d.GetContext()
	defer cancel()
	_, err := ExecuteCommandWithContext(ctx, dockerCmd, "restart", containerID)
	return err
}

// GetContainerStatus obtiene el estado de un contenedor Docker
func (d *DockerService) GetContainerStatus(containerID string) (string, error) {
	if containerID == "" {
		return "", ErrInvalidID
	}
	ctx, cancel := d.GetContext()
	defer cancel()
	status, err := ExecuteCommandWithContext(ctx, dockerCmd, "inspect", "--format", "{{.State.Status}}", containerID)
	if err != nil {
		return "", err
	}

	status = strings.TrimSpace(status)
	if status == "" {
		return "", ErrContainerNotFound
	}

	return status, nil
}

// StartContainerIfStopped inicia un contenedor Docker si está detenido
func (d *DockerService) StartContainerIfStopped(containerID string) error {
	if containerID == "" {
		return ErrInvalidID
	}

	status, err := d.GetContainerStatus(containerID)
	if err != nil {
		return err
	}
	if status != "running" {
		if err := d.StartContainer(containerID); err != nil {
			return err
		}
	}
	return nil
}

// GetDockerfileContent obtiene el contenido de un Dockerfile a partir de una plantilla
func (d *DockerService) GetDockerfileContent(param map[string]string, filePath string) (string, error) {
	var err error
	if dockerfileTemplate == nil {
		dockerfileTemplate, err = template.ParseFiles(filePath)
		if err != nil {
			return "", err
		}
	}

	params := DockerfileParams{
		FileName:      param[constants.FileNameKey],
		CommitMessage: param[constants.CommitMessageKey],
		CommitHash:    param[constants.CommitHashKey],
		CommitAuthor:  param[constants.CommitAuthorKey],
		Team:          param[constants.TeamKey],
		Organization:  param[constants.OrganizationKey],
	}

	var result strings.Builder
	err = dockerfileTemplate.Execute(&result, params)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// GetComposeContent obtiene el contenido de un archivo docker-compose a partir de una plantilla
func (d *DockerService) GetComposeContent(param map[string]string, filePath string) (string, error) {
	var err error
	if composeTemplate == nil {
		composeTemplate, err = template.ParseFiles(filePath)
		if err != nil {
			return "", err
		}
	}

	params := DockerComposeParams{
		NameDelivery: param[constants.NameDeliveryKey],
		CommitHash:   param[constants.CommitHashKey],
		Port:         param[constants.PortKey],
	}

	var result strings.Builder
	err = composeTemplate.Execute(&result, params)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// GetSonarqubeComposeContent obtiene el contenido de un archivo docker-compose para SonarQube
func (d *DockerService) GetSonarqubeComposeContent(homeDir, templateData string) (string, error) {
	type ComposeParams struct {
		HomeDir string
	}

	params := ComposeParams{
		HomeDir: homeDir,
	}

	tmpl, err := template.New("compose").Parse(templateData)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, params)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// GetContext obtiene un contexto con tiempo de espera
func (d *DockerService) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Minute)
}

func (s *DockerService) GetUrlsContainer(containerIDs []string) (string, error) {
	var result strings.Builder
	for _, id := range containerIDs {
		port, err := s.GetHostPort(id)
		if err != nil {
			return "", err
		}
		url := fmt.Sprintf(constants.MessageSuccessPublish, port)
		result.WriteString(url)
	}
	return result.String(), nil
}

func (s *DockerService) GetHostPort(containerID string) (string, error) {
	ports, err := s.GetPortContainer(containerID)
	if err != nil {
		return "", err
	}

	lines := strings.Split(ports, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Split(line, "->")
		if len(parts) == 2 {
			hostPart := strings.TrimSpace(parts[1])
			if strings.Contains(hostPart, ":") {
				hostPort := strings.Split(hostPart, ":")[1]
				return hostPort, nil
			}
		}
	}
	return "", fmt.Errorf(constants.MessageErrorNoPortHost)
}

func getArray(data string) []string {
	return strings.Split(strings.TrimSpace(data), "\n")
}
