package port

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
)

type RunCommand interface {
	Run(ctx context.Context, command string) model.InfraResultEntity
}
