package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"

type Repository interface {
	Exists(projectName, environment string) (bool, error)
	Load(projectName, environment string) (service.Context, error)
	Save(projectName string, data service.Context) error
}
