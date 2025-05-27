package port

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

type DockerContainer interface {
	Up(ctx context.Context) result.InfraResult
	Exists(ctx context.Context, commitHash, version string) result.InfraResult
	GetURLsUp(ctx context.Context, commitHash, version string) result.InfraResult
	Start(ctx context.Context) result.InfraResult
}
