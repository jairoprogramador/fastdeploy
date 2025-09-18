package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"

type EnvironmentRepository interface {
	GetEnvironments(repositoryName string) ([]entity.Environment, error)
}