package ports

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
)

type DomRepository interface {
	Save(dom *aggregates.DeploymentObjectModel) error
	Load() (*aggregates.DeploymentObjectModel, error)
}
