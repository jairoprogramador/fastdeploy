package application

import (
	"errors"

	domAgg "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domPor "github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"
)

type LoadDOMService struct {
	domRepository domPor.DomRepository
}

func NewLoadDOMService(domRepository domPor.DomRepository) *LoadDOMService {
	return &LoadDOMService{
		domRepository: domRepository,
	}
}

func (s *LoadDOMService) Load() (*domAgg.DeploymentObjectModel, error) {
	domModel, err := s.domRepository.Load()
	if err != nil {
		return &domAgg.DeploymentObjectModel{}, err
	}

	if domModel == nil {
		return &domAgg.DeploymentObjectModel{},
			errors.New("el archivo .fastdeploy/dom.yaml no existe. Por favor, ejecutar 'fd init' primero")
	}
	return domModel, nil
}
