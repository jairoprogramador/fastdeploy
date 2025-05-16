package router

import (
	"sync"
	"os"
	"path/filepath"
	"deploy/internal/domain/constant"
)

type Router struct {}

var (
	instanceRouter     *Router
	instanceOnceRouter sync.Once
)

func GetRouter() *Router {
	instanceOnceRouter.Do(func() {
		instanceRouter = &Router {}
	})
	return instanceRouter
}

func (st *Router) GetFullPathDockerComposeTemplate() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeTemplateFileName)
}

func (st *Router) GetFullPathDockerCompose() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerComposeFileName)
}

func (st *Router) GetFullPathDockerfileTemplate() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory, 
		projectName, constant.DockerRootDirectory, constant.DockerfileTemplateFileName)
}

func (st *Router) GetFullPathDockerfile() string {
	projectName := st.getProjectName()
	return filepath.Join(st.GetHomeDirectory(), constant.FastdeployRootDirectory,
		projectName, constant.DockerRootDirectory, constant.DockerfileFileName)
}

func (st *Router) GetPathProjectFile() string {
	return filepath.Join(constant.ProjectRootDirectory, constant.ProjectFileName)
}

func (st *Router) GetFullPathGlobalConfigFile() string {
	return filepath.Join(st.GetHomeDirectory(), 
	constant.FastdeployRootDirectory, constant.GlobalConfigFileName)
}

func (st *Router) GetFullPathDeploymentFile() string {
	return filepath.Join(st.GetHomeDirectory(), 
	constant.FastdeployRootDirectory, constant.DeploymentFileName)
}

func (st *Router) getProjectName() string {
	currentDir, err := st.GetProjectName()
	if err != nil {
		return ""
	}
	return currentDir
}


func (st *Router) GetProjectName() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(currentDir), nil
}

func (st *Router) GetHomeDirectory() string {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return homeDirectory
}