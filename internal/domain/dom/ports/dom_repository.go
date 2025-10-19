package ports

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
)

type DomRepository interface {
	Save(dom *aggregates.DeploymentObjectModel) error
	Load() (*aggregates.DeploymentObjectModel, error)
}
