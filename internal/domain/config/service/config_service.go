package service

import (
	"errors"
	"github.com/jairoprogramador/fastdeploy/internal/domain/config/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/config/repository"
	"github.com/jairoprogramador/fastdeploy/pkg/constant"
)

var (
	ErrConfigCanNotBeNil = errors.New(constant.ErrorConfigCannotBeNil)
	ErrConfigNotComplete = errors.New(constant.ErrorConfigNotComplete)
)

type ConfigService interface {
	Load() (*entity.ConfigEntity, error)
	Save(configEntity *entity.ConfigEntity) error
}

type configService struct {
	configRepository repository.ConfigRepository
}

func NewConfigService(
	configRepository repository.ConfigRepository,
) ConfigService {
	return &configService{
		configRepository: configRepository,
	}
}

func (s *configService) Load() (*entity.ConfigEntity, error) {
	result := s.configRepository.Load()
	if result.IsSuccess() {
		config := result.Result.(*entity.ConfigEntity)
		if !config.IsComplete() {
			return &entity.ConfigEntity{}, ErrConfigNotComplete
		}
	}
	return &entity.ConfigEntity{}, result.Error
}

func (s *configService) Save(configEntity *entity.ConfigEntity) error {
	if configEntity == nil {
		return ErrConfigCanNotBeNil
	}

	if !configEntity.IsComplete() {
		return ErrConfigNotComplete
	}
	result := s.configRepository.Save(configEntity)
	return result.Error
}
