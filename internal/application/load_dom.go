package application

import (
	"errors"

	domAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	domPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"
)

type LoadConfigService struct {
	configRepository domPor.ConfigRepository
}

func NewLoadConfigService(domRepository domPor.ConfigRepository) *LoadConfigService {
	return &LoadConfigService{
		configRepository: domRepository,
	}
}

func (s *LoadConfigService) Load() (*domAgg.Config, error) {
	domModel, err := s.configRepository.Load()
	if err != nil {
		return &domAgg.Config{}, err
	}

	if domModel == nil {
		return &domAgg.Config{},
			errors.New("el archivo .fastdeploy/dom.yaml no existe. Por favor, ejecutar 'fd init' primero")
	}
	return domModel, nil
}
