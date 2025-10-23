package services

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"
)

type GeneratorID interface {
	ProjectID(config *aggregates.Config) vos.ProjectID
}
