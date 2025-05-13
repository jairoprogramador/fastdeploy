package tools

import (
	"context"
	"deploy/internal/domain/constant"
	"strings"
	"text/template"
)

// DockerfileParams contiene los parámetros para generar un Dockerfile
type DockerfileParamsd struct {
	FileName      string
	CommitMessage string
	CommitHash    string
	CommitAuthor  string
	Team          string
	Organization  string
}

// GetDockerfileContent genera el contenido de un Dockerfile a partir de una plantilla
func GetDockerfileContents(param map[string]string, filePath string) (string, error) {
	// Cargar la plantilla si aún no está cargada
	var err error
	if dockerfileTemplate == nil {
		dockerfileTemplate, err = template.ParseFiles(filePath)
		if err != nil {
			return "", err
		}
	}

	params := DockerfileParams{
		FileName:      param[constant.FileNameKey],
		CommitMessage: param[constant.CommitMessageKey],
		CommitHash:    param[constant.CommitHashKey],
		CommitAuthor:  param[constant.CommitAuthorKey],
		Team:          param[constant.TeamKey],
		Organization:  param[constant.OrganizationKey],
	}

	var result strings.Builder
	err = dockerfileTemplate.Execute(&result, params)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

// GetComposeContent contiene el contenido de un archivo docker-compose
func GetComposeContents(param map[string]string, filePath string) (string, error) {

	type DockerParams struct {
		NameDelivery string
		CommitHash   string
		Port         string
	}

	params := DockerParams{
		NameDelivery: param[constant.NameDeliveryKey],
		CommitHash:   param[constant.CommitHashKey],
		Port:         param[constant.PortKey],
	}

	tmpl, err := template.ParseFiles(filePath)
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

// GetSonarqubeComposeContent contiene el contenido de un archivo docker-compose para SonarQube
func GetSonarqubeComposeContents(homeDir, templateData string) (string, error) {

	type ComposeParams struct {
		HomeDir string
	}

	params := ComposeParams{
		HomeDir: homeDir,
	}

	//tmpl, err := template.ParseFiles(filePath)
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

// DockerClient define una interfaz para operaciones con Docker
type DockerClient interface {
	GetImageID(ctx context.Context, hashCommit string) (string, error)
	GetContainersID(ctx context.Context, imageID string) ([]string, error)
	StartContainer(ctx context.Context, containerID string) error
	RestartContainer(ctx context.Context, containerID string) error
	GetContainerStatus(ctx context.Context, containerID string) (string, error)
	BuildImage(ctx context.Context, hashCommit, filePath string) error
	// ... otros métodos...
}

// DefaultDockerClient implementa DockerClient usando comandos de shell
type DefaultDockerClient struct{}

// NewDockerClient crea una nueva instancia de DockerClient
func NewDockerClient() DockerClient {
	return &DefaultDockerClient{}
}

// Implementa los métodos de DockerClient en DefaultDockerClient
// ... (implementaciones de los métodos utilizando las funciones que ya teníamos)

// GetImageID devuelve el ID de la imagen Docker correspondiente al hash de commit especificado.
// Si la imagen no existe, devuelve un string vacío y ningún error.
func (c *DefaultDockerClient) GetImageID(ctx context.Context, hashCommit string) (string, error) {
	if hashCommit == "" {
		return "", ErrInvalidID
	}

	imageID, err := ExecuteCommandWithContext(ctx, dockerCmd, "images", "-q", hashCommit)
	if err != nil {
		return "", err
	}

	imageID = strings.TrimSpace(imageID)
	if imageID == "" {
		return "", ErrImageNotFound
	}

	return imageID, nil
}

// GetContainersID devuelve los IDs de contenedores que usan la imagen especificada
func (c *DefaultDockerClient) GetContainersID(ctx context.Context, imageID string) ([]string, error) {
	if imageID == "" {
		return nil, ErrInvalidID
	}

	ancestor := "ancestor=" + imageID
	containerIDs, err := ExecuteCommandWithContext(ctx, dockerCmd, "ps", "-qa", "--filter", ancestor)
	if err != nil {
		return nil, err
	}

	if containerIDs == "" {
		return nil, nil
	}

	return strings.Split(strings.TrimSpace(containerIDs), "\n"), nil
}

// StartContainer inicia un contenedor Docker
func (c *DefaultDockerClient) StartContainer(ctx context.Context, containerID string) error {
	if containerID == "" {
		return ErrInvalidID
	}
	_, err := ExecuteCommandWithContext(ctx, dockerCmd, "start", containerID)
	return err
}

// RestartContainer reinicia un contenedor Docker
func (c *DefaultDockerClient) RestartContainer(ctx context.Context, containerID string) error {
	if containerID == "" {
		return ErrInvalidID
	}
	_, err := ExecuteCommandWithContext(ctx, dockerCmd, "restart", containerID)
	return err
}

// GetContainerStatus devuelve el estado actual de un contenedor Docker
func (c *DefaultDockerClient) GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	if containerID == "" {
		return "", ErrInvalidID
	}

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

// BuildImage construye una imagen Docker a partir de un Dockerfile
func (c *DefaultDockerClient) BuildImage(ctx context.Context, hashCommit, filePath string) error {
	_, err := ExecuteCommand("docker", "build", "-t", hashCommit, "-f", filePath, ".")
	return err
}

func (c *DefaultDockerClient) GetContainersId(imageId string) ([]string, error) {
	ancestor := "ancestor=" + imageId
	containerIds, err := ExecuteCommand("docker", "ps", "-qa", "--filter", ancestor)
	if err != nil {
		return []string{}, err
	}

	if containerIds == "" {
		return []string{}, nil
	}

	return getArray(containerIds), nil
}

func (c *DefaultDockerClient)GetPortContainer(containerId string) (string, error) {
	return ExecuteCommand("docker", "port", containerId)
}

func (c *DefaultDockerClient)BuildContainer(filePath string) error {
	_, err := ExecuteCommand("docker", "compose", "-f", filePath, "up", "-d")
	return err
}

func (c *DefaultDockerClient) StartContainerIfStopped(ctx context.Context, containerID string) error {
	if containerID == "" {
		return ErrInvalidID
	}

	status, err := c.GetContainerStatus(ctx, containerID)
	if err != nil {
		return err
	}

	if status != "running" {
		return c.StartContainer(ctx, containerID)
	}

	return nil
}

func getArrays(data string) []string {
	return strings.Split(strings.TrimSpace(data), "\n")
}
