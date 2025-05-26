package application

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/service"
)

func StartDeploy(projectService service.ProjectService,

) model.DomainResultEntity {
	return projectService.Start()
}
