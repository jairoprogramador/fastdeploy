package repository

import "deploy/internal/domain/model"

type GlobalConfigRepository interface {
	Create(globalConfig *model.GlobalConfig) error
	Load() (model.GlobalConfig, error)
	ExistsFile() bool
	RemoveFile() error
}
