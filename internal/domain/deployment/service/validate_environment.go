package service

import (
	"slices"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/port"
)

type ValidateEnvironment interface {
	IsValidEnvironment(repositoryName string, environment string) (bool, error)
}

type ValidateEnvironmentImpl struct {
	environmentRepository port.EnvironmentRepository
}

func NewValidateEnvironment(environmentRepository port.EnvironmentRepository) ValidateEnvironment {
	return &ValidateEnvironmentImpl{environmentRepository: environmentRepository}
}

func (v *ValidateEnvironmentImpl) IsValidEnvironment(repositoryName string, environment string) (bool, error) {
	if environment == "local" {
		return true, nil
	}
	environments, err := v.environmentRepository.GetEnvironments(repositoryName)
	if err != nil {
		return false, err
	}
	return slices.Contains(environments, entity.NewEnvironment(environment)), nil
}
