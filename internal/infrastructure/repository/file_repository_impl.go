package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/variable"
	"deploy/internal/infrastructure/filesystem"
	"sync"
	"strings"
)

type fileRepositoryImpl struct {
	homeDirectory string
}

var (
	instanceFileRepository     repository.FileRepository
	instanceOnceFileRepository sync.Once
)

func GetFileRepository() repository.FileRepository {
	instanceOnceFileRepository.Do(func() {
		homeDirectory, _ := filesystem.GetHomeDirectory()
		instanceFileRepository = &fileRepositoryImpl{
			homeDirectory: homeDirectory,
		}
	})
	return instanceFileRepository
}

func (st *fileRepositoryImpl) GetFullPathDockerComposeTemplate(store variable.VariableStore) string {
	projectName := store.Get(constant.VAR_PROJECT_NAME)
	projectName = strings.ReplaceAll(projectName, " ", "")
	dockerRootDirectory := store.Get(constant.VAR_DOCKER_ROOT_DIRECTORY)
	dockerComposeTemplateFileName := store.Get(constant.VAR_DOCKERCOMPOSE_TEMPLATE_FILE_NAME)

	return filesystem.GetPath(st.homeDirectory, projectName, dockerRootDirectory, dockerComposeTemplateFileName)
}

func (st *fileRepositoryImpl) GetFullPathDockerCompose(store variable.VariableStore) string {
	projectName := store.Get(constant.VAR_PROJECT_NAME)
	projectName = strings.ReplaceAll(projectName, " ", "")
	dockerRootDirectory := store.Get(constant.VAR_DOCKER_ROOT_DIRECTORY)
	dockerComposeFileName := store.Get(constant.VAR_DOCKERCOMPOSE_FILE_NAME)

	return filesystem.GetPath(st.homeDirectory, projectName, dockerRootDirectory, dockerComposeFileName)
}

func (st *fileRepositoryImpl) GetFullPathDockerfileTemplate(store variable.VariableStore) string {
	projectName := store.Get(constant.VAR_PROJECT_NAME)
	projectName = strings.ReplaceAll(projectName, " ", "")
	dockerRootDirectory := store.Get(constant.VAR_DOCKER_ROOT_DIRECTORY)
	dockerfileTemplateFileName := store.Get(constant.VAR_DOCKERFILE_TEMPLATE_FILE_NAME)

	return filesystem.GetPath(st.homeDirectory, projectName, dockerRootDirectory, dockerfileTemplateFileName)
}

func (st *fileRepositoryImpl) GetFullPathDockerfile(store variable.VariableStore) string {
	projectName := store.Get(constant.VAR_PROJECT_NAME)
	projectName = strings.ReplaceAll(projectName, " ", "")
	dockerRootDirectory := store.Get(constant.VAR_DOCKER_ROOT_DIRECTORY)
	dockerfileFileName := store.Get(constant.VAR_DOCKERFILE_FILE_NAME)

	return filesystem.GetPath(st.homeDirectory, projectName, dockerRootDirectory, dockerfileFileName)
}

func (st *fileRepositoryImpl) ExistsFile(path string) bool {
	exists, _ := filesystem.ExistsFile(path)
	return exists
}

func (st *fileRepositoryImpl) DeleteFile(path string) error {
	return filesystem.RemoveFile(path)
}


