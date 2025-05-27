package path

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
	"os"
	"path/filepath"
)

type pathAdapter struct{}

func NewPathAdapter() port.PathPort {
	return &pathAdapter{}
}

func (s *pathAdapter) GetFullPathDockerComposeTemplate() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeTemplateFileName)
}

func (s *pathAdapter) GetFullPathDockerCompose() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeFileName)
}

func (s *pathAdapter) GetFullPathDockerfileTemplate() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileTemplateFileName)
}

func (s *pathAdapter) GetFullPathDockerfile() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileFileName)
}

func (s *pathAdapter) GetPathDockerDirectory() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory)
}

func (s *pathAdapter) GetPathProjectFile() string {
	return filepath.Join(constant.ProjectRootDirectory, constant.ProjectFileName)
}

func (s *pathAdapter) GetFullPathConfigFile() string {
	return filepath.Join(s.GetHomeDirectory(),
		constant.FastDeployRootDirectory, constant.ConfigFileName)
}

func (s *pathAdapter) GetFullPathDeploymentFile() string {
	return filepath.Join(s.GetHomeDirectory(),
		constant.FastDeployRootDirectory, constant.DeploymentFileName)
}

func (s *pathAdapter) GetFullPathLoggerFile() string {
	projectName, err := s.GetProjectName()
	if err != nil {
		return ""
	}
	return filepath.Join(s.GetHomeDirectory(),
		constant.FastDeployRootDirectory, projectName, constant.LoggerFileName)
}

func (s *pathAdapter) GetProjectName() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(currentDir), nil
}

func (s *pathAdapter) GetHomeDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDirectory
}

func (s *pathAdapter) GetRelativePathFromHome(absolutePath string) string {
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
