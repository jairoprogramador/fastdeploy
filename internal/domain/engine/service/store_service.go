package service

import (
	"context"
	"errors"
	modelDeploy "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/model"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

const (
	errorProjectIsNil    = "project cannot be nil"
	errorDeploymentIsNil = "deployment cannot be nil"
)

type StoreServiceInterface interface {
	GetVariablesGlobal(ctx context.Context, deployment *modelDeploy.DeploymentEntity, project *model.ProjectEntity) ([]modelDeploy.Variable, error)
}

type StoreService struct {
	gitService  port.GitPort
	pathService port.PathPort
}

func NewStoreService(
	gitService port.GitPort,
	pathService port.PathPort,
) StoreServiceInterface {
	return &StoreService{
		gitService:  gitService,
		pathService: pathService,
	}
}

func (s *StoreService) GetVariablesGlobal(ctx context.Context, deployment *modelDeploy.DeploymentEntity, project *model.ProjectEntity) ([]modelDeploy.Variable, error) {
	if project == nil {
		return []modelDeploy.Variable{}, errors.New(errorProjectIsNil)
	}

	if deployment == nil {
		return []modelDeploy.Variable{}, errors.New(errorDeploymentIsNil)
	}

	response := s.gitService.GetHash(ctx)
	if !response.IsSuccess() {
		return []modelDeploy.Variable{}, response.Error
	}
	commitHash := response.Result.(string)

	response = s.gitService.GetAuthor(ctx, commitHash)
	if !response.IsSuccess() {
		return []modelDeploy.Variable{}, response.Error
	}
	commitAuthor := response.Result.(string)

	response = s.gitService.GetMessage(ctx, commitHash)
	if !response.IsSuccess() {
		return []modelDeploy.Variable{}, response.Error
	}
	commitMessage := response.Result.(string)

	var variables []modelDeploy.Variable

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyProjectOrganization,
		Value: project.Organization,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyProjectId,
		Value: project.ProjectID,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyProjectName,
		Value: project.Name,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyProjectVersion,
		Value: project.Version,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyProjectTeam,
		Value: project.TeamName,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyCommitHash,
		Value: commitHash,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyCommitAuthor,
		Value: commitAuthor,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyCommitMessage,
		Value: commitMessage,
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyPathHomeDirectory,
		Value: s.pathService.GetHomeDirectory(),
	})

	variables = append(variables, modelDeploy.Variable{
		Name:  constant.KeyPathDockerDirectory,
		Value: s.pathService.GetPathDockerDirectory(),
	})

	for _, variable := range deployment.Variables.Global {
		variables = append(variables, modelDeploy.Variable{
			Name:  variable.Name,
			Value: variable.Value,
		})
	}
	return variables, nil
}
