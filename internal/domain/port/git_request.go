package port

import (
	"context"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

type GitRequest interface {
	GetHash(ctx context.Context) result.InfraResult
	GetAuthor(ctx context.Context, commitHash string) result.InfraResult
	GetMessage(ctx context.Context, commitHash string) result.InfraResult
}
