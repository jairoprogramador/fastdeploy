package config

import (
	"fmt"
)

type ConfigServiceImpl struct {
	repository ConfigRepository
	validator  ConfigValidator
}

func NewConfigService(
	repository ConfigRepository,
	validator ConfigValidator,
) ConfigService {
	return &ConfigServiceImpl{
		repository: repository,
		validator:  validator,
	}
}

func (cs *ConfigServiceImpl) Save(configEntity ConfigEntity) error {
	if err := cs.validate(configEntity); err != nil {
		return fmt.Errorf("configuración inválida: %w", err)
	}

	return cs.repository.Save(configEntity)
}

func (cs *ConfigServiceImpl) Load() (*ConfigEntity, error) {
	return cs.repository.Load()
}

func (cs *ConfigServiceImpl) Exists() bool {
	return cs.repository.Exists()
}

func (cs *ConfigServiceImpl) Delete() error {
	return cs.repository.Delete()
}

func (cs *ConfigServiceImpl) validate(config ConfigEntity) error {
	return cs.validator.Validate(config)
}
