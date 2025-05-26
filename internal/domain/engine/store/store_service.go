package store

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/constant"
	model2 "github.com/jairoprogramador/fastdeploy/internal/domain/engine/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"errors"
)

const (
	prefix = "engine"
)

type StoreServiceInterface interface {
	GetVariablesGlobal(ctx context.Context, deployment *model2.DeploymentEntity, project *model.ProjectEntity) ([]model2.Variable, error)
}

type StoreService struct {
	gitService  port.GitRequest
	pathService port.PathService
}

func NewStoreService(
	gitService port.GitRequest,
	pathService port.PathService,
) StoreServiceInterface {
	return &StoreService{
		gitService:  gitService,
		pathService: pathService,
	}
}

func (s *StoreService) GetVariablesGlobal(ctx context.Context, deployment *model2.DeploymentEntity, project *model.ProjectEntity) ([]model2.Variable, error) {
	if project == nil {
		return []model2.Variable{}, errors.New("data project cannot be nil")
	}

	response := s.gitService.GetHash(ctx)
	if !response.IsSuccess() {
		return []model2.Variable{}, response.Error
	}
	commitHash := response.Result.(string)

	response = s.gitService.GetAuthor(ctx, commitHash)
	if !response.IsSuccess() {
		return []model2.Variable{}, response.Error
	}
	commitAuthor := response.Result.(string)

	response = s.gitService.GetMessage(ctx, commitHash)
	if !response.IsSuccess() {
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
		Value: s.pathService.GetHomeDirectory(),
	})

	variables = append(variables, model2.Variable{
		Name:  constant.VAR_PATH_DOCKER_DIRECTORY,
		Value: s.pathService.GetPathDockerDirectory(),
	})

	for _, variable := range deployment.Variables.Global {
		variables = append(variables, model2.Variable{
			Name:  variable.Name,
			Value: variable.Value,
		})
	}
	return variables, nil
}
