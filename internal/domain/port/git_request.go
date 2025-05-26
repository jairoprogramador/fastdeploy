package port

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
)

type GitRequest interface {
	GetHash(ctx context.Context) model.InfraResultEntity
	GetAuthor(ctx context.Context, commitHash string) model.InfraResultEntity
	GetMessage(ctx context.Context, commitHash string) model.InfraResultEntity
}
