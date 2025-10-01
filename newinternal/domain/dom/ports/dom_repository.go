package ports

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/aggregates"
)

// DOMRepository define el contrato para la persistencia del agregado DeploymentObjectModel.
type DOMRepository interface {
	Save(ctx context.Context, dom *aggregates.DeploymentObjectModel) error
	Load(ctx context.Context) (*aggregates.DeploymentObjectModel, error)
}
