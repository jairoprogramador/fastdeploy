package service

import (
	"context"
	"deploy/internal/domain/constant"
	model2 "deploy/internal/domain/engine/model"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/port"
)

const (
	prefix = "engine"
)

type StoreServiceInterface interface {
	GetVariablesGlobal(ctx context.Context, deployment *model2.DeploymentEntity, project *model.ProjectEntity) ([]model2.Variable, error)
}

type StoreService struct {
	gitService port.GitRequest
	router     *PathService
	logger     *logger.Logger
}

func NewStoreService(
	logger *logger.Logger,
	gitService port.GitRequest,
	router *PathService,
) StoreServiceInterface {
	return &StoreService{
		gitService: gitService,
		router:     router,
		logger:     logger,
	}
}

func (s *StoreService) GetVariablesGlobal(ctx context.Context, deployment *model2.DeploymentEntity, project *model.ProjectEntity) ([]model2.Variable, error) {
	if project == nil {
		return []model2.Variable{}, s.logger.NewError("data project cannot be nil")
	}

	response := s.gitService.GetHash(ctx)
	if !response.IsSuccess() {
		s.setError(response.Details, response.Error)
		return []model2.Variable{}, response.Error
	}
	commitHash := response.Result.(string)

	response = s.gitService.GetAuthor(ctx, commitHash)
	if !response.IsSuccess() {
		s.setError(response.Details, response.Error)
		return []model2.Variable{}, response.Error
	}
	commitAuthor := response.Result.(string)

	response = s.gitService.GetMessage(ctx, commitHash)
	if !response.IsSuccess() {
		s.setError(response.Details, response.Error)
		return []model2.Variable{}, response.Error
	}
	commitMessage := response.Result.(string)

	variables := []model2.Variable{}

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PROJECT_ORGANIZATION,
		Value: project.Organization,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PROJECT_ID,
		Value: project.ProjectID,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PROJECT_NAME,
		Value: project.Name,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PROJECT_VERSION,
		Value: project.Version,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PROJECT_TEAM,
		Value: project.TeamName,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_COMMIT_HASH,
		Value: commitHash,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_COMMIT_AUTHOR,
		Value: commitAuthor,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_COMMIT_MESSAGE,
		Value: commitMessage,
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PATH_HOME_DIRECTORY,
		Value: s.router.GetHomeDirectory(),
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PATH_DOCKER_DIRECTORY,
		Value: s.router.GetPathDockerDirectory(),
	})

	for _, variable := range deployment.Variables.Global {
		variables = append(variables, model2.Variable{
			Name:  variable.Name,
			Value: variable.Value,
		})
	}
	return variables, nil
}

func (s *StoreService) setError(message string, err error) {
	s.logger.ErrorSystemMessage(message, err)
	s.logger.Error(err)
}
