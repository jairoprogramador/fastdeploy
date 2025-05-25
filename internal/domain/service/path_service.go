package service

import (
	"deploy/internal/domain/constant"
	"os"
	"path/filepath"
)

type PathService struct{}

func NewPathService() *PathService {
	return &PathService{}
}

func (st *PathService) GetFullPathDockerComposeTemplate() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeTemplateFileName)
}

func (st *PathService) GetFullPathDockerCompose() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeFileName)
}

func (st *PathService) GetFullPathDockerfileTemplate() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileTemplateFileName)
}

func (st *PathService) GetFullPathDockerfile() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileFileName)
}

func (st *PathService) GetPathDockerDirectory() string {
	projectName := st.getProjectName()
	return filepath.Join(constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory)
}

func (st *PathService) GetPathProjectFile() string {
	return filepath.Join(constant.ProjectRootDirectory, constant.ProjectFileName)
}

func (st *PathService) GetFullPathGlobalConfigFile() string {
	return filepath.Join(st.GetHomeDirectory(),
		constant.FastdeployRootDirectory, constant.GlobalConfigFileName)
}

func (st *PathService) GetFullPathDeploymentFile() string {
	return filepath.Join(st.GetHomeDirectory(),
		constant.FastdeployRootDirectory, constant.DeploymentFileName)
}

func (st *PathService) GetFullPathLoggerFile() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(),
		constant.FastdeployRootDirectory, projectName, constant.LoggerFileName)
}

func (st *PathService) getProjectName() string {
	currentDir, err := st.GetProjectName()
	if err != nil {
		return ""
	}
	return currentDir
}

func (st *PathService) GetProjectName() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(currentDir), nil
}

func (st *PathService) GetHomeDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDirectory
}

func (st *PathService) GetRelativePathFromHome(absolutePath string) string {
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
