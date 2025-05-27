package store

import (
	"context"
	"errors"
	entity2 "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

const (
	prefix = "engine"
)

type StoreServiceInterface interface {
	GetVariablesGlobal(ctx context.Context, deployment *entity2.DeploymentEntity, project *entity.ProjectEntity) ([]entity2.Variable, error)
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

func (s *StoreService) GetVariablesGlobal(ctx context.Context, deployment *entity2.DeploymentEntity, project *entity.ProjectEntity) ([]entity2.Variable, error) {
	if project == nil {
		return []entity2.Variable{}, errors.New("data project cannot be nil")
	}

	response := s.gitService.GetHash(ctx)
	if !response.IsSuccess() {
		return []entity2.Variable{}, response.Error
	}
	commitHash := response.Result.(string)

	response = s.gitService.GetAuthor(ctx, commitHash)
	if !response.IsSuccess() {
		return []entity2.Variable{}, response.Error
	}
	commitAuthor := response.Result.(string)

	response = s.gitService.GetMessage(ctx, commitHash)
	if !response.IsSuccess() {
		return []entity2.Variable{}, response.Error
	}
	commitMessage := response.Result.(string)

	variables := []entity2.Variable{}

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyProjectOrganization,
		Value: project.Organization,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyProjectId,
		Value: project.ProjectID,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyProjectName,
		Value: project.Name,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyProjectVersion,
		Value: project.Version,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyProjectTeam,
		Value: project.TeamName,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyCommitHash,
		Value: commitHash,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyCommitAuthor,
		Value: commitAuthor,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyCommitMessage,
		Value: commitMessage,
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyPathHomeDirectory,
		Value: s.pathService.GetHomeDirectory(),
	})

	variables = append(variables, entity2.Variable{
		Name:  constant.KeyPathDockerDirectory,
		Value: s.pathService.GetPathDockerDirectory(),
	})

	for _, variable := range deployment.Variables.Global {
		variables = append(variables, entity2.Variable{
			Name:  variable.Name,
			Value: variable.Value,
		})
	}
	return variables, nil
}
