package service

import (
	"deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"errors"
)

var (
	ErrConfigNotFound    = errors.New(constant.MsgConfigNotFound)
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
	config, err := s.configRepository.Load()
	if err == nil && config != nil {
		if !config.IsComplete() {
			return &model.ConfigEntity{}, ErrConfigNotComplete
		}
	}
	return config, err
}

func (s *configService) Save(configEntity *model.ConfigEntity) error {
	if configEntity == nil {
		return ErrConfigCanNotBeNil
	}

	if !configEntity.IsComplete() {
		return ErrConfigNotFound
	}
	return s.configRepository.Save(configEntity)
}
