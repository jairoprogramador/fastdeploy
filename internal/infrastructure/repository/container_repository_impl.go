package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"deploy/internal/infrastructure/tools"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/variable"
	"text/template"
	"fmt"
	"net"
	"sync"
	"strings"
)

type DockerfileData struct {
	FileName      string
	CommitMessage string
	CommitHash    string
	CommitAuthor  string
	Team          string
	Organization  string
}

type DockerComposeData struct {
	NameDelivery string
	CommitHash   string
	Port         string
}

type containerRepositoryImpl struct{}

var (
	instanceContainerRepository     repository.ContainerRepository
	instanceOnceContainerRepository sync.Once
)

func GetContainerRepository() repository.ContainerRepository {
	instanceOnceContainerRepository.Do(func() {
		instanceContainerRepository = &containerRepositoryImpl{}
	})
	return instanceContainerRepository
}

func (st *containerRepositoryImpl) CreateFile(pathFile string, content string) error {
	err := filesystem.WriteFile(pathFile, content)
	if err != nil {
		return err
	}
	return nil
}

func (st *containerRepositoryImpl) CreateDockerfile(pathFile, pathTemplate string, store variable.VariableStore) error {
	directoryTarget := "target"
	exists, err := filesystem.ExistsDirectory(directoryTarget)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no se encontró el directorio target")
	}

	fullPathJarFiles, err := tools.GetFullPathFiles(directoryTarget)
	if err != nil {
		return err
	}

	dockerfileTemplate, err := template.ParseFiles(pathTemplate)
	if err != nil {
		return err
	}

	params := DockerfileData{
		FileName:      fullPathJarFiles[0],
		CommitMessage: store.Get(constant.VAR_COMMIT_MESSAGE),
		CommitHash:    store.Get(constant.VAR_COMMIT_HASH),
		CommitAuthor:  store.Get(constant.VAR_COMMIT_AUTHOR),
		Team:          store.Get(constant.VAR_PROJECT_TEAM),
		Organization:  store.Get(constant.VAR_PROJECT_ORGANIZATION),
	}

	var result strings.Builder
	err = dockerfileTemplate.Execute(&result, params)
	if err != nil {
		return err
	}
	
	return st.CreateFile(pathFile, result.String())
}

func (st *containerRepositoryImpl) CreateDockerCompose(pathFile, pathTemplate string, store variable.VariableStore) error { 
	dockerComposeTemplate, err := template.ParseFiles(pathTemplate)
	if err != nil {
		return err
	}

	params := DockerComposeData{
		NameDelivery: store.Get(constant.VAR_PROJECT_NAME),
		CommitHash:   store.Get(constant.VAR_COMMIT_HASH),
		Port:         st.getPort(),
	}

	var result strings.Builder
	err = dockerComposeTemplate.Execute(&result, params)
	if err != nil {
		return err
	}

	return st.CreateFile(pathFile, result.String())
}

func (st *containerRepositoryImpl) getPort() string {
	startPort := 2000
	endPort := 3000

	portFree := 2000

	for port := startPort; port <= endPort; port++ {
		address := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", address)
		if err == nil {
			portFree = port
			ln.Close()
			break
		}
	}
	return fmt.Sprintf("%d", portFree)
}
