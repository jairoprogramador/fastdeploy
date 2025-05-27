package repository

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/config/entity"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

type ConfigRepository interface {
	Load() result.InfraResult
	Save(config *entity.ConfigEntity) result.InfraResult
}
