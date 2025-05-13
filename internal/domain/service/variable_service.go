package service

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/variable"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/repository"
	"sync"
)

type VariableServiceInterface interface {
	InitializeGlobalDefault(store *variable.VariableStore) *model.Response
}

type VariableService struct {
	projectService      ProjectServiceInterface
	variableRepository repository.VariableRepository
	muVariableService    sync.RWMutex
}

var (
	instanceVariableService     *VariableService
	instanceOnceVariableService sync.Once
)

func GetVariableService(projectService ProjectServiceInterface,
	variableRepository repository.VariableRepository,
	) VariableServiceInterface {
	instanceOnceVariableService.Do(func() {
		instanceVariableService = &VariableService {
			projectService: projectService,
			variableRepository: variableRepository,
		}
	})
	return instanceVariableService
}

func (s *VariableService) SetVariableService(projectService ProjectServiceInterface) {
	s.muVariableService.Lock()
	defer s.muVariableService.Unlock()
	s.projectService = projectService
}

func (s *VariableService) InitializeGlobalDefault(store *variable.VariableStore) *model.Response {
	project, err := s.projectService.Load()
	if err != nil {
		return model.GetNewResponseError(err)
	}
	store.AddVariableGlobal(constant.VAR_PROJECT_ORGANIZATION, project.Organization)
	store.AddVariableGlobal(constant.VAR_PROJECT_ID, project.ProjectID)
	store.AddVariableGlobal(constant.VAR_PROJECT_NAME, project.Name )
	store.AddVariableGlobal(constant.VAR_PROJECT_VERSION, project.Version)
	store.AddVariableGlobal(constant.VAR_PROJECT_TEAM, project.TeamName)
	store.AddVariableGlobal(constant.VAR_PROJECT_ROOT_DIRECTORY, constant.ProjectRootDirectory)
	store.AddVariableGlobal(constant.VAR_PROJECT_FILE_NAME, constant.ProjectFileName)
	store.AddVariableGlobal(constant.VAR_DOCKER_ROOT_DIRECTORY, constant.DockerRootDirectory)
	store.AddVariableGlobal(constant.VAR_DOCKERFILE_FILE_NAME, constant.DockerfileFileName)
	store.AddVariableGlobal(constant.VAR_DOCKERFILE_TEMPLATE_FILE_NAME, constant.DockerfileTemplateFileName)
	store.AddVariableGlobal(constant.VAR_DOCKERCOMPOSE_FILE_NAME, constant.DockerComposeFileName)
	store.AddVariableGlobal(constant.VAR_DOCKERCOMPOSE_TEMPLATE_FILE_NAME, constant.DockerComposeTemplateFileName)

	commitHash, err := s.variableRepository.GetCommitHash()
	if err != nil {
		return model.GetNewResponseError(err)
	}
	store.AddVariableGlobal(constant.VAR_COMMIT_HASH, commitHash)

	commitAuthor, err := s.variableRepository.GetCommitAuthor(commitHash)
	if err != nil {
		return model.GetNewResponseError(err)
	}
	store.AddVariableGlobal(constant.VAR_COMMIT_AUTHOR, commitAuthor)

	commitMessage, err := s.variableRepository.GetCommitMessage(commitHash)
	if err != nil {
		return model.GetNewResponseError(err)
	}
	store.AddVariableGlobal(constant.VAR_COMMIT_MESSAGE, commitMessage)	

	return model.GetNewResponse()
}
