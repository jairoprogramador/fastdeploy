package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/cli/presenter"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/service"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"github.com/jairoprogramador/fastdeploy/pkg/common/result"
)

func InitApp(projectService service.ProjectService, fileLogger *logger.FileLogger) result.DomainResult {
	spinner := presenter.ShowMessageSpinner("implementation in progress")

	result := projectService.Initialize()
	fileLogger.WriteToFile()

	defer spinner.Stop()
	return result
}
