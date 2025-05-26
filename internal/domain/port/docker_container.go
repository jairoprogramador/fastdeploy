package port

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
)

type DockerContainer interface {
	Up(ctx context.Context) model.InfraResultEntity
	Exists(ctx context.Context, commitHash, version string) model.InfraResultEntity
	GetURLsUp(ctx context.Context, commitHash, version string) model.InfraResultEntity
	Start(ctx context.Context) model.InfraResultEntity
}
