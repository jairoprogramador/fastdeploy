package common

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"os"
	"path/filepath"
)

type pathService struct{}

func NewPathService() port.PathService {
	return &pathService{}
}

func (st *pathService) GetFullPathDockerComposeTemplate() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeTemplateFileName)
}

func (st *pathService) GetFullPathDockerCompose() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeFileName)
}

func (st *pathService) GetFullPathDockerfileTemplate() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileTemplateFileName)
}

func (st *pathService) GetFullPathDockerfile() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileFileName)
}

func (st *pathService) GetPathDockerDirectory() string {
	projectName := st.getProjectName()
	return filepath.Join(constant.FastDeployRootDirectory,
		projectName, constant.DockerRootDirectory)
}

func (st *pathService) GetPathProjectFile() string {
	return filepath.Join(constant.ProjectRootDirectory, constant.ProjectFileName)
}

func (st *pathService) GetFullPathConfigFile() string {
	return filepath.Join(st.GetHomeDirectory(),
		constant.FastDeployRootDirectory, constant.ConfigFileName)
}

func (st *pathService) GetFullPathDeploymentFile() string {
	return filepath.Join(st.GetHomeDirectory(),
		constant.FastDeployRootDirectory, constant.DeploymentFileName)
}

func (st *pathService) GetFullPathLoggerFile() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(),
		constant.FastDeployRootDirectory, projectName, constant.LoggerFileName)
}

func (st *pathService) getProjectName() string {
	currentDir, err := st.GetProjectName()
	if err != nil {
		return ""
	}
	return currentDir
}

func (st *pathService) GetProjectName() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(currentDir), nil
}

func (st *pathService) GetHomeDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDirectory
}

func (st *pathService) GetRelativePathFromHome(absolutePath string) string {
	homeDir := st.GetHomeDirectory()
	if homeDir == "" {
		return absolutePath
	}

	relativePath, err := filepath.Rel(homeDir, absolutePath)
	if err != nil {
		return absolutePath
	}

	return relativePath
}
