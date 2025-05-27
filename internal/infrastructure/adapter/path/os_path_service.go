// Package adapter provides implementations of domain interfaces for external dependencies.
package path

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
	"os"
	"path/filepath"
)

// OsPathService implements the port.PathService interface using the os package.
type OsPathService struct{}

// NewOsPathService creates a new OsPathService instance.
func NewOsPathService() port.PathService {
	return &OsPathService{}
}

// GetFullPathDockerComposeTemplate returns the full path to the Docker Compose template file.
func (s *OsPathService) GetFullPathDockerComposeTemplate() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeTemplateFileName)
}

// GetFullPathDockerCompose returns the full path to the Docker Compose file.
func (s *OsPathService) GetFullPathDockerCompose() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeFileName)
}

// GetFullPathDockerfileTemplate returns the full path to the Dockerfile template.
func (s *OsPathService) GetFullPathDockerfileTemplate() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileTemplateFileName)
}

// GetFullPathDockerfile returns the full path to the Dockerfile.
func (s *OsPathService) GetFullPathDockerfile() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileFileName)
}

// GetPathDockerDirectory returns the path to the Docker directory.
func (s *OsPathService) GetPathDockerDirectory() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory)
}

// GetPathProjectFile returns the path to the project file.
func (s *OsPathService) GetPathProjectFile() string {
	return filepath.Join(constant.ProjectRootDirectory, constant.ProjectFileName)
}

// GetFullPathConfigFile returns the full path to the configuration file.
func (s *OsPathService) GetFullPathConfigFile() string {
	return filepath.Join(s.GetHomeDirectory(),
		constant.FastDeployRootDirectory, constant.ConfigFileName)
}

// GetFullPathDeploymentFile returns the full path to the deployment file.
func (s *OsPathService) GetFullPathDeploymentFile() string {
	return filepath.Join(s.GetHomeDirectory(),
		constant.FastDeployRootDirectory, constant.DeploymentFileName)
}

// GetFullPathLoggerFile returns the full path to the logger file.
func (s *OsPathService) GetFullPathLoggerFile() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(),
		constant.FastDeployRootDirectory, projectName, constant.LoggerFileName)
}

// GetProjectName returns the name of the current project.
func (s *OsPathService) GetProjectName() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(currentDir), nil
}

// GetHomeDirectory returns the home directory of the current user.
func (s *OsPathService) GetHomeDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDirectory
}

// GetRelativePathFromHome returns the relative path from the home directory.
func (s *OsPathService) GetRelativePathFromHome(absolutePath string) string {
	homeDir := s.GetHomeDirectory()
	if homeDir == "" {
		return absolutePath
	}

	relativePath, err := filepath.Rel(homeDir, absolutePath)
	if err != nil {
		return absolutePath
	}

	return relativePath
}
