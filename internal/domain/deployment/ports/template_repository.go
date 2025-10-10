package ports

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

type TemplateRepository interface {
	GetTemplate(ctx context.Context, source vos.TemplateSource) (template *aggregates.DeploymentTemplate, repoLocalPath string, err error)
	GetRepositoryName(repoURL string) (string, error)
}
