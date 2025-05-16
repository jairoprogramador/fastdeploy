package service

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/repository"
	"sync"
	"context"
)

type StoreServiceInterface interface {
	GetVariablesGlobal(ctx context.Context, deployment *model.Deployment) ([]model.Variable, error) 
}

type StoreService struct {
	project *model.Project
	gitRepository repository.GitRepository
	muVariableService    sync.RWMutex
}

var (
	instanceStoreService     *StoreService
	instanceOnceStoreService sync.Once
)

func GetStoreService(project *model.Project,
	gitRepository repository.GitRepository) StoreServiceInterface {
	instanceOnceStoreService.Do(func() {
		instanceStoreService = &StoreService {
			project: project,
			gitRepository: gitRepository,
		}
	})
	return instanceStoreService
}

func (s *StoreService) SetProjectModel(project *model.Project) {
	s.muVariableService.Lock()
	defer s.muVariableService.Unlock()
	s.project = project
}

func (s *StoreService) GetVariablesGlobal(ctx context.Context, deployment *model.Deployment) ([]model.Variable, error) {
	commitHash, err := s.gitRepository.GetCommitHash(ctx)
	if err != nil {
		return []model.Variable{}, err
	}

	commitAuthor, err := s.gitRepository.GetCommitAuthor(ctx, commitHash)
	if err != nil {
		return []model.Variable{}, err
	}

	commitMessage, err := s.gitRepository.GetCommitMessage(ctx, commitHash)
	if err != nil {
		return []model.Variable{}, err
	}

	variables := []model.Variable{}

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_ORGANIZATION,
		Value: s.project.Organization,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_ID,
		Value: s.project.ProjectID,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_NAME,
		Value: s.project.Name,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_VERSION,
		Value: s.project.Version,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_TEAM,
		Value: s.project.TeamName,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_COMMIT_HASH,
		Value: commitHash,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_COMMIT_AUTHOR,
		Value: commitAuthor,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_COMMIT_MESSAGE,
		Value: commitMessage,
	})

	for _, variable := range deployment.Variables.Global {
		variables = append(variables, model.Variable{
			Name:  variable.Name,
			Value: variable.Value,
		})
	}
	return variables, nil
}
