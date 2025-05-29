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

type StoreServicePort interface {
	InitStore(ctx context.Context) error
	AddDataProject(project *model.ProjectEntity) error
	AddDataDeployment(deployment *modelDeploy.DeploymentEntity) error
	GetStore() *modelDeploy.StoreEntity
}

type StoreService struct {
	gitService  port.GitPort
	pathService port.PathPort
	store       *modelDeploy.StoreEntity
}

func NewStoreService(
	gitService port.GitPort,
	pathService port.PathPort,
	store *modelDeploy.StoreEntity,
) StoreServicePort {
	return &StoreService{
		gitService:  gitService,
		pathService: pathService,
		store:       store,
	}
}

func (s *StoreService) GetStore() *modelDeploy.StoreEntity {
	return s.store
}

func (s *StoreService) InitStore(ctx context.Context) error {

	response := s.gitService.GetHash(ctx)
	if !response.IsSuccess() {
		return response.Error
	}
	commitHash := response.Result.(string)

	response = s.gitService.GetAuthor(ctx, commitHash)
	if !response.IsSuccess() {
		return response.Error
	}
	commitAuthor := response.Result.(string)

	response = s.gitService.GetMessage(ctx, commitHash)
	if !response.IsSuccess() {
		return response.Error
	}
	commitMessage := response.Result.(string)

	var variables []modelDeploy.Variable

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

	s.store.Initialize(variables)
	return nil
}

func (s *StoreService) AddDataProject(project *model.ProjectEntity) error {
	if project == nil {
		return errors.New(errorProjectIsNil)
	}

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

	s.store.AddVariables(variables)
	return nil
}

func (s *StoreService) AddDataDeployment(deployment *modelDeploy.DeploymentEntity) error {
	if deployment == nil {
		return errors.New(errorDeploymentIsNil)
	}

	if deployment.Variables.Global != nil {
		var variables []modelDeploy.Variable

		for _, variable := range deployment.Variables.Global {
			variables = append(variables, modelDeploy.Variable{
				Name:  variable.Name,
				Value: variable.Value,
			})
		}

		s.store.AddVariables(variables)
	}

	return nil
}
