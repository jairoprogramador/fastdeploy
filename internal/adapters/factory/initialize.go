package factory

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type InitializeFactory interface {
	CreateInitialize() project.ProjectInitialize
}
