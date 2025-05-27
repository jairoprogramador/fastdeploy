package port

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

type RunCommand interface {
	Run(ctx context.Context, command string) result.InfraResult
}
