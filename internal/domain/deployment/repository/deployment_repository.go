package repository

import (
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

type DeploymentRepository interface {
	Load() result.InfraResult
}
