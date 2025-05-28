package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/service"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

func StartDeploy(projectService service.ProjectService, fileLogger *logger.FileLogger) result.DomainResult {
	result := projectService.Start()
	fileLogger.WriteToFile()
	return result
}
