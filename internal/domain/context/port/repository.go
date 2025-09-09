package port

import "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"

type Repository interface {
	Exists(projectName string) (bool, error)
	Load(projectName string) (service.Context, error)
	Save(projectName string, data service.Context) error
}
