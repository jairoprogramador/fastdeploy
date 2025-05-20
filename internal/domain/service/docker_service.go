package service

import (
	"context"
	"deploy/internal/domain/model"
)

type DockerServiceInterface interface {
	ExistsContainer(ctx context.Context, variableStore *model.VariableStore) (bool, error)
	DockerComposeUpBuild(ctx context.Context, pathDockerCompose string, variableStore *model.VariableStore) (string, error)
	DockerComposeUp(ctx context.Context, pathDockerCompose string, variableStore *model.VariableStore) (string, error)
	DockerComposeDown(ctx context.Context, pathDockerCompose string) error
}
