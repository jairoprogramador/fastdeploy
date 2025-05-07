package repository

import "context"

type DockerRepository interface {
	GetImageID(hashCommit string) (string, error)
	GetContainersID(imageID string) ([]string, error)
	SonarScanner(token, projectKey, projectName, projectPath, cacheDir, tmpDir, scannerWorkDir, sourcePath, testPath, binaryPath, testBinaryPath string) error
	GetPortContainer(containerId string) (string, error)
	BuildImage(hashCommit string, filePath string) error
	BuildContainer(filePath string) error
	StartContainer(containerID string) error
	RestartContainer(containerID string) error
	GetContainerStatus(containerID string) (string, error)
	StartContainerIfStopped(containerID string) error
	GetDockerfileContent(param map[string]string, filePath string) (string, error)
	GetComposeContent(param map[string]string, filePath string) (string, error)
	GetSonarqubeComposeContent(homeDir, templateData string) (string, error)
	GetUrlsContainer(containerIDs []string) (string, error)
	GetHostPort(containerID string) (string, error)
	GetContext() (context.Context, context.CancelFunc)
}
