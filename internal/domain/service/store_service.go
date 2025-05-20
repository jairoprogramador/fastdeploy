package service

import (
	"context"
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/router"
	"fmt"
)

type StoreServiceInterface interface {
	GetVariablesGlobal(ctx context.Context, deployment *model.Deployment, project *model.Project) ([]model.Variable, error)
}

type StoreService struct {
	gitService GitServiceInterface
	router     *router.Router
	//muVariableService sync.RWMutex
}

func NewStoreService(
	gitService GitServiceInterface,
	router *router.Router,
) StoreServiceInterface {
	return &StoreService{
		gitService: gitService,
		router:     router,
	}
}

func (s *StoreService) GetVariablesGlobal(ctx context.Context, deployment *model.Deployment, project *model.Project) ([]model.Variable, error) {
	if project == nil {
		return nil, fmt.Errorf("el proyecto no puede ser nulo para GetVariablesGlobal")
	}

	commitHash, err := s.gitService.GetCommitHash(ctx)
	if err != nil {
		return []model.Variable{}, err
	}

	commitAuthor, err := s.gitService.GetCommitAuthor(ctx, commitHash)
	if err != nil {
		return []model.Variable{}, err
	}

	commitMessage, err := s.gitService.GetCommitMessage(ctx, commitHash)
	if err != nil {
		return []model.Variable{}, err
	}

	variables := []model.Variable{}

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_ORGANIZATION,
		Value: project.Organization,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_ID,
		Value: project.ProjectID,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_NAME,
		Value: project.Name,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_VERSION,
		Value: project.Version,
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PROJECT_TEAM,
		Value: project.TeamName,
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

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PATH_HOME_DIRECTORY,
		Value: s.router.GetHomeDirectory(),
	})

	variables = append(variables, model.Variable{
		Name:  constant.VAR_PATH_DOCKER_DIRECTORY,
		Value: s.router.GetPathDockerDirectory(),
	})

	for _, variable := range deployment.Variables.Global {
		variables = append(variables, model.Variable{
			Name:  variable.Name,
			Value: variable.Value,
		})
	}
	return variables, nil
}
