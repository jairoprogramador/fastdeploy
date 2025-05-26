package application

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/service"
)

func InitApp(projectService service.ProjectService) model.DomainResultEntity {
	return projectService.Initialize()
}
