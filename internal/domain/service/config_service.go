package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/constant"
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/repository"
	"errors"
)

var (
	ErrConfigCanNotBeNil = errors.New(constant.MsgConfigCannoBeNil)
	ErrConfigNotComplete = errors.New(constant.MsgConfigNotComplete)
)

type ConfigService interface {
	Load() (*model.ConfigEntity, error)
	Save(configEntity *model.ConfigEntity) error
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

func (s *configService) Load() (*model.ConfigEntity, error) {
	result := s.configRepository.Load()
	if result.IsSuccess() {
		config := result.Result.(*model.ConfigEntity)
		if !config.IsComplete() {
			return &model.ConfigEntity{}, ErrConfigNotComplete
		}
	}
	return &model.ConfigEntity{}, result.Error
}

func (s *configService) Save(configEntity *model.ConfigEntity) error {
	if configEntity == nil {
		return ErrConfigCanNotBeNil
	}

	if !configEntity.IsComplete() {
		return ErrConfigNotComplete
	}
	result := s.configRepository.Save(configEntity)
	return result.Error
}
