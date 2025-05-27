package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/service"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

func StartDeploy(projectService service.ProjectService) result.DomainResult {
	return projectService.Start()
}
