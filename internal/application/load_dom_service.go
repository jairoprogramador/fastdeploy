package application

import (
	"errors"
	"context"

	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"
)

type LoadDOMService struct {
	domRepository ports.DOMRepository
}

func NewLoadDOMService(domRepository ports.DOMRepository) *LoadDOMService {
	return &LoadDOMService{
		domRepository: domRepository,
	}
}

func (s *LoadDOMService) Load(ctx context.Context) (*aggregates.DeploymentObjectModel, error) {
	domModel, err := s.domRepository.Load(ctx)
	if err != nil {
		return &aggregates.DeploymentObjectModel{}, err
	}

	if domModel == nil {
		return &aggregates.DeploymentObjectModel{},
		errors.New("el archivo .fastdeploy/dom.yaml no existe. Por favor, ejecutar 'fd init' primero")
	}

	return domModel, nil
}